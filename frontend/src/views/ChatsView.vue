<template>
    <div class="chats-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <div class="mb-8">
                    <div class="flex justify-between items-center mb-3">
                        <h3 class="text-xl font-semibold text-nordic-dark">Privat chats</h3>
                        <span v-if="isConnected"
                            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            Online
                        </span>
                        <span v-else
                            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                            Offline
                        </span>
                    </div>
                    <ul v-if="chats.length > 0" class="space-y-2">
                        <li v-for="chat in chats" :key="chat.id" @click="select(chat)" :class="[
                            'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 p-2 rounded-lg',
                            selectedChat && selectedChat.id === chat.id ? 'font-semibold text-nordic-primary-accent bg-gray-100' : ''
                        ]">
                            <div class="flex justify-between">
                                <span>{{ chat.name }}</span>
                                <span class="text-xs text-gray-500" v-if="getLastMessage(chat)">
                                    {{ formatTime(getLastMessage(chat).timestamp) }}
                                </span>
                            </div>
                            <div class="text-sm text-gray-500 truncate" v-if="getLastMessage(chat)">
                                {{ getLastMessage(chat).sender }}: {{ getLastMessage(chat).text }}
                            </div>
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">No chats yet.</p>
                </div>

                <div class="mb-8">
                         <h4 class="text-xl font-semibold text-nordic-dark">Start a conversation: </h4>
                    <h4 class="text-lg font-medium text-nordic-dark mb-2">Following</h4>
                    <ul v-if="followed.length > 0" class="space-y-2">
                        <li v-for="user in followed" :key="user.id">
                            <button @click="startChat(user)"
                                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
                                {{ user.name }}
                            </button>
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">You follow no one yet.</p>
                </div>

                <div>
                    <h4 class="text-lg font-medium text-nordic-dark mb-2">Followers</h4>
                    <ul v-if="followers.length > 0" class="space-y-2">
                        <li v-for="user in followers" :key="user.id">
                            <button @click="startChat(user)"
                                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
                                {{ user.name }}
                            </button>
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">No followers yet.</p>
                </div>
            </template>

            <template #main>
                <div v-if="selectedChat" class="mb-6">
                    <div class="flex justify-between items-center">
                        <h2 class="text-3xl font-bold text-nordic-dark">{{ selectedChat.name }}</h2>
                        <button v-if="!isConnected" @click="connectWebSocket"
                            class="px-4 py-2 bg-nordic-primary-accent text-white rounded-lg hover:bg-nordic-dark transition-colors">
                            Reconnect
                        </button>
                    </div>
                </div>
                <div v-else class="mb-6">
                    <h2 class="text-3xl font-bold text-nordic-dark">Chat</h2>
                    <p class="text-nordic-light">Select a conversation or start a new one</p>
                </div>

                <ChatBox :chat="selectedChat" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import TopBar from '@/components/TopBar.vue'
import ChatBox from '@/components/ChatBox.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const websocketStore = useWebSocketStore()
const isConnected = computed(() => websocketStore.isConnected)
const wsUrl = import.meta.env.VITE_WS_URL

//temporary data for testing
const chats = ref([
    {
        id: 1,
        name: 'Omar',
        userId: '101', // Keep user IDs as strings to avoid Go JSON parsing issues
        messages: [
            { id: 1, text: 'Hey!', sender: 'Omar', timestamp: new Date(Date.now() - 3600000) },
            { id: 2, text: 'Hi there!', sender: 'You', timestamp: new Date(Date.now() - 3500000) }
        ]
    }
])

//temporary data for testing
const selectedChat = ref(null)
const followed = ref([{ id: 2, name: 'Dolgorsürengiin', userId: '102' }])
const followers = ref([{ id: 3, name: 'Alex', userId: '103' }])


onMounted(() => {
  
    // when this ChatsView.vue page is rendered /mounted:
    // fetch chats list from API
    // display chat list
    // open the first chat

    // Set default selected_chat
    if (chats.value.length > 0) {
        selectedChat.value = chats.value[0]
    }

    connectWebSocket()
    listenForNewChats()
})

function connectWebSocket() {
    websocketStore.connect(websocketStore.connect(`${wsUrl}/ws`))
}

//listens for new converasation starters by other users
function listenForNewChats() {
    // Store the current user ID for comparison
    const currentUserId = localStorage.getItem('userId') || '0';

   //listen for new webSocket chat_messages:
    websocketStore.$subscribe((mutation, state) => {
        const message = state.message;
        if (!message || message.type !== 'chat_message') return;

        //check that the message is not from the current user itself
        if (message.receiver_id !== currentUserId) return;

        // check if current users chats allready include one with the same sender
        const existingChat = chats.value.find(c => c.userId === message.sender_id);

        //if not existing, create a new one
        if (!existingChat && message.sender_id) {

            //fetch sender info based on sender_id
            //replace this with a fetch from API
            fetchUserInfo(message.sender_id).then(senderData => {
                // Create a new chat with this sender
                const newChat = {
                    chatId: Date.now(),
                    chatPartnerName: senderData.name,
                    userId: message.sender_id,
                    messages: [{
                        id: Date.now(),
                        text: message.content,
                        sender: senderData.name
                    }]
                };
                chats.value.push(newChat);
            });
        }
    });
}

// Placeholder function 
async function fetchUserInfo(sender_id) {
    // replace this with a real API call
    console.log(`Fetching user info for ID: ${sender_id}`);

    return {
        id: userId,
        name: `User ${userId}` // Placeholder name
    };
}

function select(chat) {
    selectedChat.value = chat
}

function startChat(user) {
    const existingChat = chats.value.find(chat => chat.userId === user.userId)

    if (existingChat) {
        selectedChat.value = existingChat
    } else {
        const newChat = {
            id: Date.now(),
            chatPartnerName: user.name,
            userId: user.userId,
            messages: []
        }
        chats.value.push(newChat)
        selectedChat.value = newChat
    }
}

function getLastMessage(chat) {
    if (!chat.messages || chat.messages.length === 0) return null
    return chat.messages[chat.messages.length - 1]
}

function formatTime(timestamp) {
    if (!timestamp) return ''

    const date = new Date(timestamp)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.chats-page-wrapper {
    min-height: 100vh;
}
</style>