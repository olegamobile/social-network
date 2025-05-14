<template>
    <div class="chat-box">
        <div class="messages">
            <div v-for="msg in chat.messages" :key="msg.id" class="message">
                <strong>{{ msg.sender }}:</strong> {{ msg.text }}
            </div>
        </div>

        <div class="max-w-lg w-full lg:max-w-xs mt-4">
            <input type="text"
                class="block w-full py-2 px-3 border border-nordic-light rounded-md bg-white 
                text-nordic-dark placeholder-nordic-light focus:outline-none focus:ring-2 
                focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
                v-model="newMessage" @keyup.enter="sendMessage" placeholder="Type a message..." />
        </div>


    </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({ chat: Object })
const newMessage = ref('')

function sendMessage() {
    if (!newMessage.value.trim()) return
    props.chat.messages.push({ id: Date.now(), text: newMessage.value, sender: 'You' })
    newMessage.value = ''
}
</script>
