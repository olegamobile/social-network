<template>
    <div class="space-y-4 bg-white p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
        <h2 class="text-xl font-semibold text-gray-800">Create a New Group</h2>

        <!-- title -->
        <input type="text" v-model="title" placeholder="Enter title" class="w-full p-3 border border-gray-300 rounded-md text-gray-700 placeholder-gray-400 
      focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y"></input>

        <!-- description -->
        <textarea v-model="description" rows="4" placeholder="Enter description" class="w-full p-3 border border-gray-300 rounded-md text-gray-700 placeholder-gray-400 
      focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y"></textarea>

        <!-- submit button -->
        <button @click="createGroup" :disabled="!title.trim() || !description.trim()" class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md 
    bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none 
    focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium 
    disabled:opacity-50 disabled:cursor-not-allowed">
            Create
        </button>
    </div>
</template>

<script setup>
import { ref } from 'vue'

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['group-created'])
const title = ref('')
const description = ref('')


const createGroup = async () => {
    try {
        const res = await fetch(`${apiUrl}/api/groups/create`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ title: title.value, description: description.value }),
        })

        if (res.ok) {
            const group = await res.json()
            emit('group-created', group.id)
            title.value = ''
            description.value = ''
        } else {
            alert('Failed to create group. Are you logged in?')
        }
    } catch (error) {
        console.error(error)
    }
}
</script>