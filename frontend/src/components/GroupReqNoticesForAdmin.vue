<template>

    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Requests to Join</h3>
    <GroupRequestsList :notifications="filteredNotifications" @close="handleClose" @accept="handleAccept"
        @decline="handleDecline" />

</template>

<script setup>

import GroupRequestsList from '@/components/GroupRequestsList.vue'

import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'

const showCurrent = ref(true)
const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const route = useRoute()
const router = useRouter()
const { logout } = useAuth()

const fetchedNotifications = ref([])

const filteredNotifications = computed(() =>
    (fetchedNotifications.value || []).filter(n => n.is_read !== showCurrent.value)
)


async function fetchNotifications() {
    try {
        const res = await fetch(`${apiUrl}/api/notifications/${route.params.id}/joingroup`, {
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
    await approveGroupRequest(n.group_id, n.sender_id, 'accepted');
    await readNotification(n.id)
    await fetchNotifications();
}

async function handleDecline(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (!n.is_read) n.is_read = true
    await approveGroupRequest(n.group_id, n.sender_id, 'declined');
    await readNotification(n.id)
    await fetchNotifications();
}

</script>
