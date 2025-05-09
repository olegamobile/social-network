import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user' // adjust path if needed

const apiUrl = import.meta.env.VITE_API_URL

export function useAuth() {
    const router = useRouter()
    const userStore = useUserStore()

    async function logout() {
        await fetch(`${apiUrl}/api/logout`, {
            method: 'POST',
            credentials: 'include',
        })
        userStore.clearUser()
        router.push('/login')
    }

    return { logout }
}
