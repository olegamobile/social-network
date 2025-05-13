<template>
    <div class="notifications-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Notifications</h3>
                <button @click="showCurrent = true" :class="{ active: showCurrent }">Current</button>
                <button @click="showCurrent = false" :class="{ active: !showCurrent }">Old</button>
            </template>

            <template #main>
                <h2 class="text-2xl font-bold mb-4">{{ showCurrent ? 'Current' : 'Old' }} Notifications</h2>
                <NotificationsList :notifications="filteredNotifications" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import NotificationsList from '@/components/NotificationsList.vue'

const showCurrent = ref(true)

const mockNotifications = ref([
    { id: 1, text: "Omar accepted your follow request", isOld: false },
    { id: 2, text: "There's a new event 'Graduation Party' in your group 'Batch of 2024", isOld: false },
    { id: 3, text: "Dolgorsürengiin accepted your follow request", isOld: false },
    { id: 4, text: "You were removed from the group 'Old Friends'", isOld: true },
    { id: 5, text: "New message from Alex", isOld: true }
])

const filteredNotifications = computed(() =>
    mockNotifications.value.filter(n => n.isOld !== showCurrent.value)
)
</script>

<style scoped>
.notifications-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}

button {
    display: block;
    margin: 0.5rem 0;
    padding: 0.5rem 1rem;
    background-color: #eee;
    border: none;
    cursor: pointer;
    border-radius: 5px;
}

button.active {
    background-color: #ccc;
    font-weight: bold;
}
</style>
