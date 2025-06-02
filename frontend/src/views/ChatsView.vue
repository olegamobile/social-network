<template>
    <div class="chats-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>

                <!-- Existing chats -->
                <div class="mb-8">

                    <!-- Chats title -->
                    <div class="flex justify-between items-center mb-3">
                        <h3 class="text-xl font-semibold text-nordic-dark">Chats</h3>
                        <span v-if="isConnected"
                            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            Online
                        </span>
                        <span v-else
                            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                            Offline
                        </span>
                    </div>

                    <!-- Links to existing chats -->
                    <ul v-if="chats.length > 0" class="space-y-2">
                        <li v-for="chat in chats" :key="chat.id" @click="select(chat)" :class="[
                            'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 p-2 rounded-lg',
                            selectedChat && selectedChat.name === chat.name ? 'font-semibold text-nordic-primary-accent bg-gray-100' : ''
                        ]">
                            <span>{{ chat.name }}</span>
                            <div class="text-sm text-gray-500 truncate" v-if="getLastMessage(chat)">
                                {{ getLastMessage(chat).sender_name }}: {{ getLastMessage(chat).content }}
                            </div>
                            <span class="text-xs text-gray-500" v-if="getLastMessage(chat)">
                                {{ formatDate(getLastMessage(chat).created_at) }}
                            </span>

                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">No chats yet.</p>
                </div>


                <!-- Start a new chat -->
                <h4 class="text-xl font-semibold text-nordic-dark mb-3">Start a new chat</h4>

                <!-- Following -->
                <div class="mb-4">
                    <h4 class="text-lg font-medium text-nordic-dark">Following</h4>
                    <ul v-if="followedUsers.length > 0" class="space-y-2">
                        <li v-for="user in followedUsers" :key="user.id">
                            <button @click="startChat(user)"
                                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
                                {{ user.first_name }} {{ user.last_name }}
                            </button>
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">You follow no one yet.</p>
                </div>

                <!-- Followers -->
                <div>
                    <h4 class="text-lg font-medium text-nordic-dark">Followers</h4>
                    <ul v-if="followers.length > 0" class="space-y-2">
                        <li v-for="user in followers" :key="user.id">
                            <button @click="startChat(user)"
                                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
                                {{ user.first_name }} {{ user.last_name }}
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
                <!--                 <div v-else class="mb-6">
                    <h2 class="text-3xl font-bold text-nordic-dark">Chat</h2>
                    <p class="text-nordic-light">Select a conversation or start a new one</p>
                </div> -->

                <ChatBox :chat="selectedChat" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, toRefs, watch } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import TopBar from '@/components/TopBar.vue'
import ChatBox from '@/components/ChatBox.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';
import { useAuth } from '@/composables/useAuth'
import { useErrorStore } from '@/stores/error'
import { useRouter } from 'vue-router'

const websocketStore = useWebSocketStore()
const isConnected = computed(() => websocketStore.isConnected)
const wsUrl = import.meta.env.VITE_WS_URL
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const followedUsers = ref([])
const followers = ref([])
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()
const apiUrl = import.meta.env.VITE_API_URL

//temporary data for testing
const selectedChat = ref(null)

//temporary data for testing
const chats = ref([])


onMounted(async () => {

    // when this ChatsView.vue page is rendered /mounted:
    // fetch chats list from API
    // display chat list
    // open the first chat

    fetchFollowData()
    //connectWebSocket()    // Should be connected already by main.js
    //listenForNewChats()   // use watch() instead
    testConnection()
    await getAllUserChats()

    // Set default selected_chat
    if (chats.value.length > 0) {
        selectedChat.value = chats.value[0]
    }
})

function connectWebSocket() {
    websocketStore.connect(websocketStore.connect(`${wsUrl}/ws`))
}

function testConnection() {
    const nuMsg = {
        type: "ping",
        from: user.value.id,
    }
    websocketStore.send(nuMsg)
}

//listen for new webSocket chat_messages:
watch(() => websocketStore.message, (message) => {

    //console.log("New message gotten:", message)

    if (!message || message.type !== 'chat_message') return;

    //check if the message is not for the current user
    if (message.receiver_id != user.value.id) {     // soft equal, string and number
        console.log("wrong id:", message.receiver_id, user.value.id)
        return
    }

    // check if current users chats already include one with the same sender
    const existingChat = chats.value.find(c => c.user_id === message.sender_id);

    //if not existing, create a new one
    if (!existingChat && message.sender_id) {

        //fetch sender info based on sender_id 
        fetch(`/api/users/${message.sender_id}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch user info');
                }
                return response.json();
            })
            .then(senderData => {
                const newChat = {
                    is_active: true,
                    name: senderData.first_name + " " + senderData.last_name,
                    userId: message.sender_id,
                    /*                         messages: [{
                                                id: Date.now(),
                                                text: message.content,
                                                sender: senderData.name
                                            }] */
                    messages: [message.content]
                };
                chats.value.push(newChat);
            })
            .catch(error => {
                console.error('Error fetching user info:', error);
            });
    } else {
        console.log("chat exists, not creating a new one")
    }
})

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

const formatDate = (isoString) => {
    const date = new Date(isoString)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'short'
    }).replace("klo ", "")
}

async function fetchFollowData() {
    try {

        const [res1, res2] = await Promise.all([
            fetch(`${apiUrl}/api/followed/0`, { credentials: 'include' }),
            fetch(`${apiUrl}/api/followers/0`, { credentials: 'include' })
        ])

        if (res1.status === 401 || res2.status === 401) {
            logout()
            router.push('/login')
            return
        }

        if (!res1.ok || !res2.ok) throw new Error('Failed to fetch follow data')

        const [followedJson, followersJson] = await Promise.all([
            await res1.json(),
            await res2.json()
        ])

        if (followedJson) followedUsers.value = followedJson
        if (followersJson) followers.value = followersJson

    } catch (err) {
        errorStore.setError('Error', 'Failed to load follow data.')
        router.push('/error')
        console.log('error fetching in chatsVue')
    }
}

async function getAllUserChats() {
    try {
        const chatResp = await fetch(`${apiUrl}/api/chat/messages`, { credentials: 'include' })

        if (chatResp.status === 401) {
            logout()
            router.push('/login')
            return
        }

        if (!chatResp.ok) {
            const errorData = await res.json().catch(() => ({ message: 'Failed to fetch posts and parse error.' }));
            isLoading.value = false
            throw new Error(errorData.message || `HTTP error ${res.status}`)
        }

        chats.value = await chatResp.json()

        console.log("user chats gotten:", chats.value)

    } catch (error) {
        errorStore.setError('Error Loading Posts', error.message || 'An unexpected error occurred while trying to load posts. Please try again later.');
        router.push('/error')
        return
    }
}


</script>

<style scoped>
.chats-page-wrapper {
    min-height: 100vh;
}
</style>
