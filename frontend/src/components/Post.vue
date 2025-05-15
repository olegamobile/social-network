<template>
    <div class="post-card">
        <p class="post-content">{{ post.content }}</p>
        <small class="post-date">
            Posted by
            <RouterLink :to="`/profile/${post.user_id}`" class="post-user">
                {{ post.username }}
            </RouterLink>
            on {{ formattedDate }}
        </small>
        <button @click="toggleComments" class="mt-2 text-blue-500 hover:underline">
      {{ showComments ? 'Hide Comments' : 'Show Comments' }}
    </button>

    <div v-if="showComments" class="mt-2">
      <div v-if="loadingComments">Loading comments...</div>
      <div v-else-if="Array.isArray(comments) && comments.length === 0">No comments yet.</div>
      <ul v-else>
        <li v-for="comment in comments" :key="comment.id" class="mt-1 border-t pt-1">
          <strong>{{ comment.user.first_name }} {{ comment.user.last_name }}</strong>: {{ comment.content }}
        </li>
      </ul>
    </div>
    </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useUserStore } from '@/stores/user';
import { storeToRefs } from 'pinia';

const { post } = defineProps({
    post: {
        type: Object,
        required: true
    }
})

const formattedDate = computed(() => {
    const date = new Date(post.created_at)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'medium'
    }).replace("klo ", "")
})


const comments = ref([]);
const showComments = ref(false);
const loadingComments = ref(false);
const apiUrl = import.meta.env.VITE_API_URL || '/api';
const userStore = useUserStore();
const { user } = storeToRefs(userStore)

const toggleComments = async () => {
  showComments.value = !showComments.value;
  if (showComments.value && comments.value.length === 0) {
    loadingComments.value = true;
    if (!user.value || !user.value.id) {
  console.error("User not loaded, cannot fetch comments");
  return;
}
    try {
        
        const res = await fetch(`${apiUrl}/api/comments/show?post_id=${post.id}&user_id=${user.value.id}`, {
        credentials: 'include'
        });
      comments.value = await res.json();
    } catch (error) {
      comments.value = [];
    } finally {
      loadingComments.value = false;
    }
  }
};
</script>


<style scoped>
.post-card {
    border: 1px solid #ccc;
    padding: 1rem;
    margin-bottom: 1rem;
    border-radius: 8px;
    background-color: #fafafa;
}

.post-content {
    font-size: 1.1rem;
    margin-bottom: 0.5rem;
}

.post-date {
    color: #666;
    font-size: 0.85rem;
}

.post-user {
    color: #0077cc;
    text-decoration: none;
    font-weight: bold;
}

.post-user:hover {
    text-decoration: underline;
}
</style>