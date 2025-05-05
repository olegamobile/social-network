<template>
    <div class="login">
        <h2>Login</h2>
        <form @submit.prevent="login">
            <input v-model="email" type="email" placeholder="Email" required />
            <input v-model="password" type="password" placeholder="Password" required />
            <button type="submit">Login</button>
            <p v-if="error">{{ error }}</p>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const error = ref('')
const router = useRouter()

async function login() {
    const res = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email: email.value, password: password.value }),
    })

    if (res.ok) {
        const data = await res.json()
        router.push(`/profile/${data.user.id}`)
    } else {
        const msg = await res.text()
        error.value = msg || 'Login failed'
    }
}
</script>

<style scoped>
.login {
    max-width: 400px;
    margin: 2rem auto;
}
</style>