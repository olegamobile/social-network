import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

createApp(App).use(router).mount('#app')

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
        storeUser(user)  // Set to your Vue store or global state
        router.push('/') // Go to home
    })
    .catch(() => {
        router.push('/login')
    })
