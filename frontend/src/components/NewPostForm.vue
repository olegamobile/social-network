<template>
  <div class="new-post-form">

    <div class="mb-0">
      <h3 for="textarea" class="text-lg font-semibold">Create a New Post</h3>
      <textarea id="textarea" v-model="content" rows=4 cols=100
        class="block w-full py-2 px-3 border border-nordic-light rounded-md bg-white text-nordic-dark placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm h-32 resize-y"
        placeholder="Enter your message"></textarea>
    </div>

    <div class="mb-0">
      <label for="select" class="block text-sm font-medium text-nordic-dark mb-1">Select privacy level</label>
      <select id="select" v-model="privacy_level"
        class="block w-full py-2 px-3 border border-nordic-light rounded-md bg-white 
        text-nordic-dark focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent 
        focus:border-nordic-secondary-accent sm:text-sm">
        <option value="public" selected>Public</option>
        <option value="almost_private">Almost Private</option>
        <option value="private">Private</option>
      </select>
    </div>

    <button @click="submitPost" :disabled="!content.trim()" 
    class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md 
    bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none 
    focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium">Post</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['post-submitted'])
const content = ref('')
const privacy_level = ref('public') // Default to 'public'

const submitPost = async () => {
  const res = await fetch(`${apiUrl}/api/posts/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include', // to send the session cookie
    body: JSON.stringify({ content: content.value, privacy_level: privacy_level.value })
  })

  if (res.ok) {
    const newPost = await res.json()
    emit('post-submitted', newPost)
    content.value = ''
    privacy_level.value = 'public'
  } else {
    alert('Failed to post. Are you logged in?')
  }
}
</script>

<style scoped>
.new-post-form {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
  margin: 1rem 0 2rem;
}

textarea {
  width: 100%;
  resize: vertical;
}

button {
  margin-top: 0.5rem;
  padding: 0.5rem 1rem;
}
</style>