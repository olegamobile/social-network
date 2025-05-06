import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import { useUserStore } from './stores/user'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')

const userStore = useUserStore()

fetch('http://localhost:8080/api/me', {
    credentials: 'include'
})
    .then(res => {
        if (res.status === 200) {
            return res.json()
        } else {
            throw new Error('Not logged in')
        }
    })
    .then(user => {
        console.log("User in main.js:",user)
        userStore.setUser(user)
        router.push('/') // Go to home
    })
    .catch(() => {
        router.push('/login')
    })
