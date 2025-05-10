<template>
    <div class="login">
        <TopBar />

        <h2>Login</h2>
        <form @submit.prevent="login">
            <input v-model="email" type="email" placeholder="Email" required autocomplete="email" />
            <input v-model="password" type="password" placeholder="Password" required autocomplete="current-password" />
            <button type="submit">Login</button>
            <p v-if="error">{{ error }}</p>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import TopBar from '@/components/TopBar.vue'

const email = ref('')
const password = ref('')
const error = ref('')
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const apiUrl = import.meta.env.VITE_API_URL

async function login() {

    try {
        const res = await fetch(`${apiUrl}/api/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ email: email.value, password: password.value }),
        })
        if (res.ok) {
            const data = await res.json()
            userStore.setUser(data.user)

            // Navigate to what the user wanted or home 
            const redirectTo = route.query.redirect || '/'
            router.push(redirectTo)
        } else {
            const msg = await res.text()
            error.value = msg || 'Login failed'
        }

    } catch (error) {
        errorStore.setError('Unexpected Error', 'Something went wrong while logging in.')
        router.push('/error')
        return
    }

}
</script>

<style scoped>
.login {
    display: flex;
    flex-direction: column;
    width: 100%;
    align-items: center;
}
</style>