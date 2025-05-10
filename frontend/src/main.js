import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import { useUserStore } from './stores/user'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

const userStore = useUserStore()
const apiUrl = import.meta.env.VITE_API_URL || '/api'

fetch(`${apiUrl}/api/me`, {
    credentials: 'include'
})
    .then(res => {
        if (res.status === 200) return res.json()
        throw new Error('Not logged in')
    })
    .then(user => {
        userStore.setUser(user)
        app.mount('#app')   // mount app after credentials check
        //router.push('/') // Go to home
    })
    .catch(() => {
        // If not logged in and not already on login page, redirect to login
        if (router.currentRoute.value.path !== '/login') {
            router.replace({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } })
        }
        app.mount('#app')
    })
