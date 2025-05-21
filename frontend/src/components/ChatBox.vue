<template>
    <div class="chat-box">
        <div v-if="chat" class="messages-container mb-4 p-4 bg-gray-50 rounded-lg h-96 overflow-y-auto">
            <div v-if="chat.messages && chat.messages.length > 0">
                <div v-for="msg in chat.messages" :key="msg.id" :class="[
                    'mb-3 p-3 rounded-lg max-w-xs',
                    msg.sender === 'You' ? 'ml-auto bg-nordic-primary-accent text-white' : 'bg-gray-200'
                ]">
                    <p class="text-xs font-semibold mb-1">{{ msg.sender }}</p>
                    <p>{{ msg.text }}</p>
                    <span class="text-xs opacity-70 block text-right">{{ formatTime(msg.timestamp) }}</span>
                </div>
            </div>
            <div v-else class="h-full flex items-center justify-center">
                <p class="text-nordic-light italic">Start a conversation...</p>
            </div>
        </div>

        <div v-if="chat" class="message-input-container">
            <form @submit.prevent="sendMessage" class="flex gap-2">
                <input id="message-input" v-model="newMessage" type="text" placeholder="Type a message..."
                    class="flex-grow p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent"
                    :disabled="!isConnected" />
                <button type="submit"
                    class="px-4 py-2 bg-nordic-primary-accent text-white rounded-lg hover:bg-nordic-dark transition-colors"
                    :disabled="!newMessage.trim() || !isConnected">
                    Send
                </button>
            </form>
            <div v-if="!isConnected" class="text-red-500 text-sm mt-2">
                Not connected to chat server. Please reconnect. {{ isConnected }}
            </div>
        </div>

        <div v-else class="h-64 flex items-center justify-center">
            <p class="text-nordic-light italic">Select a chat or start a new conversation</p>
        </div>
    </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/user';

const props = defineProps({
    chat: {
        type: Object,
        default: null
    }
})

const newMessage = ref('')
const websocketStore = useWebSocketStore()
//const isConnected = ref(websocketStore.isConnected)
const { isConnected } = storeToRefs(websocketStore)
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

watch(() => websocketStore.message, (newMessage) => {
    if (newMessage && newMessage.type === 'chat_message' && props.chat) {
        // Message belongs to this chat if:
        // - We're the sender and the receiver is the chat partner, OR
        // - We're the receiver and the sender is the chat partner
        if ((newMessage.sender_id === user.id && newMessage.receiver_id === props.chat.userId) ||
            (newMessage.receiver_id === user.id && newMessage.sender_id === props.chat.userId)) {

            // Add the message to the current chat
            props.chat.messages.push({
                id: Date.now(),
                text: newMessage.content,
                sender: newMessage.sender_id === user.id ? 'You' : props.chat.name,
                timestamp: new Date()
            })
        }
    }
})

function sendMessage() {
    if (!newMessage.value.trim() || !props.chat || !isConnected.value) return

    // Add message to local chat first
    const message = {
        id: Date.now(),
        text: newMessage.value,
        sender: 'You',
        timestamp: new Date()
    }
    props.chat.messages.push(message)

    // Send via websocket - using string values for all fields to avoid numeric parsing issues
    websocketStore.send({
        type: 'chat_message',
        sender_id: user.id,
        receiver_id: props.chat.userId || '0', // Use the user ID from the chat object
        content: newMessage.value
    })

    // Clear input
    newMessage.value = ''
}

function formatTime(timestamp) {
    if (!timestamp) return ''

    const date = new Date(timestamp)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
</script>