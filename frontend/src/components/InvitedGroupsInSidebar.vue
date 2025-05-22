<template>
    <div>
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Group Invitations</h3>
        <ul v-if="groups && groups.length > 0" class="space-y-2 mb-5">
            <li v-for="group in groups" :key="group.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                <RouterLink :to="`/group/${group.id}`">
                    {{ group.title }}
                </RouterLink>
            </li>
        </ul>
        <p v-else class="text-nordic-light italic mb-5">No pending group invitations</p>
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


async function getInvites() {
    try {
        const res = await fetch(`${apiUrl}/api/groups/invitations`, {
            credentials: 'include'
        });

        if (res.status === 401) {
            errorStore.setError('Session Expired', 'Your session has expired. Please log in again.');
            logout(); // your logout function
            router.push('/login');
            return;
        }

        if (!res.ok) {
            const errorData = await res.json().catch(() => ({ message: 'Failed to fetch invitations and parse error.' }));
            throw new Error(errorData.message || `HTTP error ${res.status}`);
        }
        groups.value = await res.json()
    } catch (error) {
        errorStore.setError('Error Loading Invitations', error.message || 'An unexpected error occurred while trying to load invitations. Please try again later.');
        router.push('/error')
        return
    }
}
onMounted(() => {
    getInvites()
})
</script>