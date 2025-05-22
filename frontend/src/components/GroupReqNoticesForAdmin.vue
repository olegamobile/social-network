<template>

    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Requests to Join</h3>
    <NotificationsList :notifications="filteredNotifications" @close="handleClose" @accept="handleAccept"
        @decline="handleDecline" />

</template>

<script setup>

import NotificationsList from '@/components/NotificationsList.vue'

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


onMounted(() => {
    fetchNotifications()
})


function handleClose(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (n) n.is_read = true
    readNotification(n.id)
}

function handleAccept(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (!n.is_read) n.is_read = true
    readNotification(n.id)
    approveFollowRequest(n.follow_req_id, 'accept')
}

function handleDecline(id) {
    const n = fetchedNotifications.value.find(n => n.id === id)
    if (!n.is_read) n.is_read = true
    readNotification(n.id)
    approveFollowRequest(n.follow_req_id, 'decline')
}

</script>
