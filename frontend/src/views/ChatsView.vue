<template>
    <div class="chats-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Chats</h3>
                <ul>
                    <li v-for="chat in chats" :key="chat.id" :class="{ active: selectedChat.id === chat.id }"
                        @click="select(chat)">
                        {{ chat.name }}
                    </li>
                </ul>

                <h4>Followed Users</h4>
                <ul>
                    <li v-for="user in followed" :key="user.id">
                        <button @click="startChat(user)">{{ user.name }}</button>
                    </li>
                </ul>

                <h4>Followers</h4>
                <ul>
                    <li v-for="user in followers" :key="user.id">
                        <button @click="startChat(user)">{{ user.name }}</button>
                    </li>
                </ul>
            </template>

            <template #main>
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
    { id: 1, name: 'Omar', messages: [{ id: 1, text: 'Hey!', sender: 'Omar' }, { id: 2, text: 'Hi there!', sender: 'You' }] }
])

const selectedChat = ref(chats.value[0])

const followed = ref([{ id: 2, name: 'Dolgorsürengiin' }])
const followers = ref([{ id: 3, name: 'Alex' }])

function select(chat) {
    selectedChat.value = chat
}

function startChat(user) {
    const newChat = { id: Date.now(), name: user.name, messages: [] }
    chats.value.push(newChat)
    selectedChat.value = newChat
}
</script>
