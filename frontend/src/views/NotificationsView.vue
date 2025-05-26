<template>
    <div class="notifications-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Notifications</h3>
                <button @click="showCurrent = true" :class="{ active: showCurrent }">
                    Current ({{fetchedNotifications ? fetchedNotifications.filter(n => !n.is_read).length : 0}})
                </button>
                <button @click="showCurrent = false" :class="{ active: !showCurrent }">
                    Old ({{fetchedNotifications ? fetchedNotifications.filter(n => n.is_read).length : 0}})
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
import { useRouter } from 'vue-router'
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
    (fetchedNotifications.value || []).filter(n => n.is_read !== showCurrent.value)
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
        console.log("notifications fetch error:", err)
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

async function approveFollowRequest(id, action) {
    try {
        const res = await fetch(`${apiUrl}/api/follow/requests/${id}/${action}`, {
            credentials: 'include'
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to accept/decline follow request: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while accepting/declining follow request`)
        router.push('/error')
    }
}

async function approveGroupRequest(groupID, senderID, action) {
    try {
        const res = await fetch(`${apiUrl}/api/group/requests/${action}`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                group_id: groupID,
                requester_id: senderID
            })
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to accept/decline group request: ${res.status}`)

    } catch (err) {
        errorStore.setError('Error', `Error while accepting/declining group request`)
        router.push('/error')
    }
}


onMounted(() => {
    fetchNotifications()
})


function handleClose(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (n) n.is_read = true
    readNotification(n.id)
}

async function handleAccept(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (!n.is_read) n.is_read = true
    readNotification(n.id)

    switch (n.type) {
        case 'follow_request':
            await approveFollowRequest(n.follow_req_id, 'accept');
            readNotification(n.id)
            fetchNotifications();
            break;
        case 'group_join_request':
            await approveGroupRequest(n.group_id, n.sender_id, 'accepted');
            readNotification(n.id)
            fetchNotifications();
            break;
    }

}

function handleDecline(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (!n.is_read) n.is_read = true
    readNotification(n.id)

    switch (n.type) {
        case 'follow_request':
            approveFollowRequest(n.follow_req_id, 'decline');
            break;
        case 'group_join_request':
            approveGroupRequest(n.group_id, n.sender_id, 'declined');
            break;
    }
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
