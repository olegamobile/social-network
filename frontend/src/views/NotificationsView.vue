<template>
    <div class="notifications-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Notifications</h3>
                <button @click="showCurrent = true" :class="{ active: showCurrent }">
                    Current ({{fetchedNotifications.filter(n => !n.is_read).length}})
                </button>
                <button @click="showCurrent = false" :class="{ active: !showCurrent }">
                    Old ({{fetchedNotifications.filter(n => n.is_read).length}})
                </button>
            </template>

            <template #main>
                <h2 class="text-2xl font-bold mb-4">{{ showCurrent ? 'Current' : 'Old' }} Notifications</h2>
                <NotificationsList :notifications="filteredNotifications" @close="handleClose" @accept="handleAccept"
                    @decline="handleDecline" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import NotificationsList from '@/components/NotificationsList.vue'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'

const showCurrent = ref(true)
const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const { logout } = useAuth()

const fetchedNotifications = ref([])

const filteredNotifications = computed(() =>
    fetchedNotifications.value.filter(n => n.is_read !== showCurrent.value)
)


async function fetchNotifications() {
    try {
        const res = await fetch(`${apiUrl}/api/notifications`, {
            credentials: 'include'
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to fetch notifications: ${res.status}`)

        fetchedNotifications.value = await res.json()

    } catch (err) {
        errorStore.setError('Error', 'Something went wrong while fetching notifications.')
        router.push('/error')
    }
}

async function readNotification(id) {
    try {
        const res = await fetch(`${apiUrl}/api/notifications/${id}/read`, {
            credentials: 'include'
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to mark notification ${id} as read: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while marking notification ${id} as read`)
        router.push('/error')
    }
}


onMounted(() => {
    fetchNotifications()
})


function handleClose(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (n) n.is_read = true // or false, depending on your logic
    readNotification(n.id)
}

function handleAccept(id) {
    // Accept logic (e.g., call API, then remove)
}

function handleDecline(id) {
    // Decline logic (e.g., call API, then remove)
}
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
