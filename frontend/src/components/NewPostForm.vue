<template>
  <div class="space-y-4 bg-white p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
    <h2 class="text-xl font-semibold text-gray-800">Create a New Post</h2>

    <textarea
      v-model="content"
      rows="4"
      placeholder="Enter your message"
      class="w-full p-3 border border-gray-300 rounded-md text-gray-700 placeholder-gray-400 
      focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y"
    ></textarea>

    <div>
      <label for="select" class="block text-sm font-medium text-gray-700 mb-1">
        Select privacy level
      </label>
      <select
        id="select"
        v-model="privacy_level"
        class="block w-full p-2.5 border border-gray-300 rounded-md bg-white text-gray-700 
        focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
      >
        <option value="public">Public</option>
        <option value="almost_private">Almost Private</option>
        <option value="private">Private</option>
      </select>
    </div>

    <button
      @click="submitPost"
      :disabled="!content.trim()"
      class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md 
    bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none 
    focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium 
    disabled:opacity-50 disabled:cursor-not-allowed"
    >
      Post
    </button>
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