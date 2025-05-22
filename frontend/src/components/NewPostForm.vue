<template>
  <div class="space-y-4 bg-white p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
    <h2 class="text-xl font-semibold text-gray-800">Create a New Post</h2>

    <!-- content -->
    <textarea v-model="content" rows="4" placeholder="Enter your message" class="w-full p-3 border border-gray-300 rounded-md text-gray-700 placeholder-gray-400 
      focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y"></textarea>

    <!-- privacy setting -->
    <div>
      <label for="select" class="block text-sm font-medium text-gray-700 mb-1">
        Select privacy level
      </label>
      <select id="select" v-model="privacy_level" class="block w-full p-2.5 border border-gray-300 rounded-md bg-white text-gray-700 
        focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
        <option value="public">Public</option>
        <option value="almost_private">Almost Private</option>
        <option value="private">Private</option>
      </select>
    </div>

    <!-- upload image -->
    <div class="space-y-2">
      <label class="block text-sm font-medium text-nordic-light">Image (Optional):</label>
      <input type="file" @change="handleFileUpload" accept="image/*"
        class="block w-full text-sm text-nordic-light file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-gray-100 file:text-nordic-light hover:file:bg-gray-200" />
    </div>

    <!-- submit button -->
    <button @click="submitPost" :disabled="!content.trim()" class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md 
    bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none 
    focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium 
    disabled:opacity-50 disabled:cursor-not-allowed">
      Post
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useImageProcessor } from '@/composables/useImageProcessor';

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['post-submitted'])
const content = ref('')
const image = ref(null)
const privacy_level = ref('public') // Default to 'public'

const handleFileUpload = (event) => {
  image.value = event.target.files[0];
};

const submitPost = async () => {
  const { processPostImage } = useImageProcessor();

  try {

    const payload = {
      content: content.value,
      privacy_level: privacy_level.value
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

/*     const res = await fetch(`${apiUrl}/api/posts/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include', // to send the session cookie
      body: JSON.stringify({ content: content.value, privacy_level: privacy_level.value })
    }) */

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
</script>