<template>
    <div>
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Requested Groups</h3>
        <ul v-if="groups && groups.length > 0" class="space-y-2 mb-5">
            <li v-for="group in groups" :key="group.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                <RouterLink :to="`/group/${group.id}`">
                    {{ group.title }}
                </RouterLink>
            </li>
        </ul>
        <p v-else class="text-nordic-light italic mb-5">No pending group requests</p>
    </div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useErrorStore } from '@/stores/error'

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()

const groups = ref([]);


async function getGroups() {
    try {
        console.log(`${apiUrl}/api/groups/requested`)
        const res = await fetch(`${apiUrl}/api/groups/requested`, {
            credentials: 'include' // This sends the session cookie with the request
        });

        if (res.status === 401) {
            // Session is invalid — logout and redirect
            errorStore.setError('Session Expired', 'Your session has expired. Please log in again.');
            logout(); // your logout function
            router.push('/login');
            return;
        }

        if (!res.ok) {
            // Handle other non-successful HTTP statuses (e.g., 400, 404, 500)
            const errorData = await res.json().catch(() => ({ message: 'Failed to fetch groups and parse error.' }));
            throw new Error(errorData.message || `HTTP error ${res.status}`);
        }
        groups.value = await res.json()
    } catch (error) {
        errorStore.setError('Error Loading Groups', error.message || 'An unexpected error occurred while trying to load groups. Please try again later.');
        router.push('/error')
        return
    }
}
onMounted(() => {
    getGroups()
})
</script>