<template>
    <div class="chats-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <div class="mb-8">
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Chats</h3>
                    <ul v-if="chats.length > 0" class="space-y-2">
                        <li v-for="chat in chats" :key="chat.id" @click="select(chat)" :class="[
                            'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150',
                            selectedChat.id === chat.id ? 'font-semibold text-nordic-primary-accent' : ''
                        ]">
                            {{ chat.name }}
                        </li>
                    </ul>
                    <p v-else class="italic text-nordic-light">No chats yet.</p>
                </div>

                <div class="mb-8">
                    <h4 class="text-lg font-medium text-nordic-dark mb-2">Followed Users</h4>
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
                <h2 class="text-3xl font-bold text-nordic-dark mb-6">Conversation</h2>
                <ChatBox :chat="selectedChat" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import TopBar from '@/components/TopBar.vue'
import ChatBox from '@/components/ChatBox.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const chats = ref([
    {
        id: 1,
        name: 'Omar',
        messages: [
            { id: 1, text: 'Hey!', sender: 'Omar' },
            { id: 2, text: 'Hi there!', sender: 'You' }
        ]
    }
])

const selectedChat = ref(chats.value[0])
const followed = ref([{ id: 2, name: 'Dolgorsürengiin' }])
const followers = ref([{ id: 3, name: 'Alex' }])

function select(chat) {
    selectedChat.value = chat
}

function startChat(user) {
    const newChat = { id: Date.now(), name: user.name, messages: [] }
    const existingChat = chats.value.find(chat => chat.name === user.name)
    if (!existingChat) {
        chats.value.push(newChat)
        selectedChat.value = newChat
    }
}
</script>

<style scoped>
.chats-page-wrapper {
    min-height: 100vh;
}
</style>
