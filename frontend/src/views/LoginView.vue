<template>
    <div>
        <TopBar />
        
        <div class="px-4 mt-20 flex justify-center">
        <div class="w-full max-w-md bg-white p-8 rounded-lg shadow-md border border-nordic-light">
        <h2 class="text-2xl font-bold text-center text-nordic-light mb-6">Login</h2>

        <form @submit.prevent="login" class="flex flex-col gap-4">
            <input v-model="email" type="email" placeholder="Email" required autocomplete="email" class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />
            <input v-model="password" type="password" placeholder="Password" required autocomplete="current-password" class="p-3 border border-nordic-light rounded-md focus:outline-none focus:ring-2 focus:ring-nordic-primary-accent" />
            <button type="submit" class="bg-nordic-primary-accent text-white font-medium py-2 rounded-md hover:bg-nordic-secondary-accent transition">Login</button>
            <p v-if="error">{{ error }}</p>
        </form>
    </div>
</div>
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
            //console.log("data to store at login:", data.user)
            userStore.setUser(data.user)

            // Navigate to what the user wanted or home 
            let redirectTo = route.query.redirect || '/'
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
