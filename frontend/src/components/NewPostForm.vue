<template>
  <div class="new-post-form">
    <h3>Create a New Post</h3>
    <textarea v-model="content" placeholder="What's on your mind?" rows="4"></textarea>
    <button @click="submitPost" :disabled="!content.trim()">Post</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['post-submitted'])
const content = ref('')

const submitPost = async () => {
  const res = await fetch(`${apiUrl}/api/posts/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include', // to send the session cookie
    body: JSON.stringify({ content: content.value })
  })

  if (res.ok) {
    const newPost = await res.json()
    emit('post-submitted', newPost)
    content.value = ''
  } else {
    alert('Failed to post. Are you logged in?')
  }
}
</script>

<style scoped>
.new-post-form {
  margin: 1rem 0;
}

textarea {
  width: 100%;
  padding: 0.5rem;
  resize: vertical;
}

button {
  margin-top: 0.5rem;
  padding: 0.5rem 1rem;
}
</style>