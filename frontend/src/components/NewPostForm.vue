<template>
  <div
    class="space-y-4 bg-[var(--nordic-primary-bg)] p-4 border border-[var(--nordic-border-light)] rounded-md shadow-sm w-full max-w-screen-lg">
    <h2 class="text-xl font-semibold text-[var(--nordic-text-dark)]">Create a New Post</h2>

    <!-- content -->
    <textarea v-model="content" rows="4" placeholder="Enter your message"
      class="w-full p-3 border border-[var(--nordic-border-light)] rounded-md 
             text-[var(--nordic-text-dark)] placeholder-[var(--nordic-text-light)] 
             bg-white
             focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] focus:border-[var(--nordic-secondary-accent)] resize-y"></textarea>

    <!-- privacy setting -->
    <label for="select" class="block text-sm font-medium text-[var(--nordic-text-light)]">
      Select privacy level
    </label>
    <select id="select" v-model="privacy_level"
      class="block w-full p-2.5 border border-[var(--nordic-border-light)] rounded-md 
               bg-white text-[var(--nordic-text-dark)] 
               focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] focus:border-[var(--nordic-secondary-accent)]">
      <option value="public">Public</option>
      <option value="almost_private">Almost Private</option>
      <option value="private">Private</option>
    </select>

    <!-- select who gets to see private post -->
    <div v-if="followers && privacy_level === 'private'">
      <h3 class="block text-sm font-medium text-[var(--nordic-text-light)]">Select who can see private post</h3>

      <ul v-if="followers.length > 0" class="m-1">
        <li v-for="user in followers" :key="user.id"
          class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer ml-2">
          <div class="flex items-center gap-2">

            <!-- Checkbox -->
            <input type="checkbox" :value="user.id" :checked="selectedViewers.includes(user.id)"
              @change="toggleUserSelection(user.id)" />

            <!-- Link and user info -->
            <RouterLink :to="`/profile/${user.id}`" class="flex items-center gap-2">
              {{ user.first_name }} {{ user.last_name }}
              <span v-if="user.username"> - {{ user.username }}</span>
              <div v-if="user.avatar_url" class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1">
                <img :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar" class="w-full h-full object-cover" />
              </div>
            </RouterLink>
          </div>
        </li>
      </ul>

      <p v-else class="text-nordic-light italic mb-5">No one following</p>
    </div>


    <!-- upload image -->
    <div class="space-y-2">
      <label class="block text-sm font-medium text-[var(--nordic-text-light)]">Image (Optional):</label>
      <input type="file" @change="handleFileUpload" accept="image/*" class="block w-full text-sm text-[var(--nordic-text-light)]
               file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 
               file:text-sm file:font-semibold 
               file:bg-[var(--nordic-secondary-bg)] file:text-[var(--nordic-text-dark)] 
               hover:file:bg-[var(--nordic-border-light)]" />
    </div>

    <!-- submit button -->
    <button @click="submitPost" :disabled="!content.trim()" class="inline-flex items-center py-2 px-4 border border-[var(--nordic-border-light)] rounded-md 
             bg-[var(--nordic-primary-accent)] text-white hover:bg-[var(--nordic-secondary-accent)] 
             focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] transition font-medium 
             disabled:opacity-50 disabled:cursor-not-allowed">
      Post
    </button>
  </div>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import { useImageProcessor } from '@/composables/useImageProcessor';
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['post-submitted'])
const content = ref('')
const image = ref(null)
const privacy_level = ref('public') // Default to 'public'
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const followers = ref([])
const selectedViewers = ref([])

const handleFileUpload = (event) => {
  image.value = event.target.files[0];
};

function toggleUserSelection(userId) {
  const index = selectedViewers.value.indexOf(userId)
  if (index === -1) {
    selectedViewers.value.push(userId)
  } else {
    selectedViewers.value.splice(index, 1)
  }
  //console.log("checked user ids:", selectedViewers.value)
}

const submitPost = async () => {
  const { processPostImage } = useImageProcessor();

  try {

    const payload = {
      content: content.value,
      privacy_level: privacy_level.value,
      selected_viewers: JSON.stringify(selectedViewers.value)
    };

    const formData = new FormData();
    for (const [key, value] of Object.entries(payload)) {
      formData.append(key, value);
    }

    if (image.value) {
      const processedImg = await processPostImage(image.value);
      formData.append('image', processedImg);
    }

    const res = await fetch(`${apiUrl}/api/posts/create`, {
      method: 'POST',
      body: formData,
      credentials: 'include',
    })

    if (res.ok) {
      const newPost = await res.json()
      emit('post-submitted', newPost)
      content.value = ''
      privacy_level.value = 'public'
    } else {
      alert('Failed to post. Are you logged in?')
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchFollowers() {
  try {
    const followerResp = await fetch(`${apiUrl}/api/followers/${user.value.id}`, { credentials: 'include' })
    if (followerResp.status === 401) {
      logout()
      router.push('/login')
      return
    }
    if (!followerResp.ok) throw new Error('Failed to fetch follow data')

    const followersJson = await followerResp.json()
    if (followersJson) followers.value = followersJson
  } catch (err) {
    errorStore.setError('Error', 'Failed to load follow data.')
    router.push('/error')
  }
}

onMounted(async () => {
  if (!user.value.is_public) privacy_level.value = 'almost_private'  // default post to almost_private for private profiles
  await fetchFollowers()
  //console.log("followers at new post:", followers.value)
  //console.log("active user is public?", user.value.is_public)
})

</script>