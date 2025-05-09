<template>
    <div class="chat-box">
        <div class="messages">
            <div v-for="msg in chat.messages" :key="msg.id" class="message">
                <strong>{{ msg.sender }}:</strong> {{ msg.text }}
            </div>
        </div>
        <input v-model="newMessage" @keyup.enter="sendMessage" placeholder="Type a message..." />
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
