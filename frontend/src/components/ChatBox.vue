<template>
    <div class="chat-box">

        <!-- Listing messages -->
        <div v-if="chat" ref="messagesContainer"
            class="messages-container mb-4 p-4 bg-gray-50 rounded-lg h-96 overflow-y-auto">
            <div v-if="chat.messages && chat.messages.length > 0">
                <div v-for="msg in chat.messages" :key="msg.id" :class="[
                    'mb-3 p-3 rounded-lg max-w-xs',
                    msg.sender_id === user.id ? 'ml-auto bg-nordic-primary-accent text-white' : 'bg-gray-200'
                ]">
                    <p class="text-xs font-semibold mb-1">{{ msg.sender_name }}</p>
                    <p>{{ msg.content }}</p>
                    <span class="text-xs opacity-70 block text-right">{{ formatTime(msg.created_at) }}</span>
                </div>
            </div>
            <div v-else class="h-full flex items-center justify-center">
                <p class="text-nordic-light italic">Start a conversation...</p>
            </div>
        </div>

        <!-- New message form -->
        <div v-if="chat && chat.is_active" class="message-input-container">
            <form @submit.prevent="sendMessage" class="flex gap-2">
                <input id="message-input" v-model="newMessage" type="text" placeholder="Type a message..."
                    class="flex-grow p-3 bg-[var(--nordic-secondary-bg)] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent"
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

        <div v-else-if="chat" class="flex items-center justify-center">
            <p class="text-nordic-light italic">Follow or be followed by {{ chat.name }} to send and receive messages
            </p>
        </div>

        <div v-else class="h-64 flex items-center justify-center">
            <p class="text-nordic-light italic">Select a chat or start a new conversation</p>
        </div>
    </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useWebSocketStore } from '@/stores/websocket'
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/user';

const props = defineProps({
    chat: {
        type: Object,
        default: null
    },
    groupString: {      // put "group" in front of message type when in group chat
        type: String,
        default: ''
    }
})

const route = useRoute()
const newMessage = ref('')
const websocketStore = useWebSocketStore()
const { isConnected } = storeToRefs(websocketStore)
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const messagesContainer = ref(null);

watch(() => websocketStore.message, (newMessage) => {
    if (newMessage && newMessage.type === `${props.groupString}chat_message` && props.chat) {

        //console.log("New message detected in chat box:", newMessage)

        // Add the message to the current chat
        if (!props.chat.messages) {
            props.chat.messages = []
        }

        // own messages
        if (newMessage.from == user.value.id && newMessage.receiver_id == props.chat.user_id) { // soft equal, string and number
            props.chat.messages.push({
                content: newMessage.content,
                sender_name: user.value.first_name,
                created_at: new Date()
            })
        }

        // in private chat
        if (newMessage.type === 'chat_message' && props.groupString === '' && newMessage.receiver_id == user.value.id && newMessage.from == props.chat.user_id) {
            props.chat.messages.push({
                content: newMessage.content,
                //sender_name: String(props.chat.name).split(" ")[0],
                sender_name: newMessage.from_name,
                created_at: new Date()
            })
        }

        // in group chat
        if (newMessage.type === 'groupchat_message' && props.groupString === 'group' && newMessage.receiver_id == props.chat.user_id) {
            props.chat.messages.push({
                content: newMessage.content,
                sender_name: newMessage.from_name,
                created_at: new Date()
            })
        }

    }
})

function sendMessage() {
    if (!newMessage.value.trim() || !props.chat || !isConnected.value) return

    // Add message to local chat first
    const message = {
        created_at: Date.now(),
        content: newMessage.value,
        sender_name: user.value.first_name,
        sender_id: user.value.id
    }

    if (props.chat.messages) {
        props.chat.messages.push(message)
    } else {
        props.chat.messages = [message]
    }

    // Send via websocket - using string values for all fields to avoid numeric parsing issues
    websocketStore.send({
        type: `${props.groupString}chat_message`,
        from_name: user.value.first_name,
        receiver_id: props.chat.user_id || '0', // Use the user ID from the chat object, means group id when group chat
        content: newMessage.value
    })

    // Clear input
    newMessage.value = ''
}

function formatTime(isoString) {
    const date = new Date(isoString)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'short'
    }).replace("klo ", "")
}

function defaultChat() {
    return {
        "is_active": true,
        "name": "",
        "user_id": route.params.id,
        "messages": []
    }
}

// Watch for changes in chat.messages. Syntax: watch(source, callback, options?)
watch(
    () => props.chat?.messages,
    () => {
        nextTick(() => {    // nextTick ensures the DOM is updated before scrolling
            if (messagesContainer.value) {
                messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
            }
        })
    },
    { deep: true }  // detects changes inside the messages array
)

// Scroll to bottom on mount (in case there are already messages)
onMounted(() => {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
        }
    })
    if (!props.chat) {
        props.chat = defaultChat()
    }
})
</script>