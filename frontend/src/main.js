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

try {
    const res = await fetch(`${apiUrl}/api/me`, {
        credentials: 'include'
    });

    if (res.status !== 200) {
        throw new Error('Not logged in');
    }

    const user = await res.json();
    userStore.setUser(user);
    app.mount('#app'); // mount app after successful credentials check
} catch (error) {
    // If not logged in and not already on login page, redirect to login
    if (router.currentRoute.value.path !== '/login' && router.currentRoute.value.path !== '/register') {
        if (router.currentRoute.value.path === '/error') {  // don't redirect to error page
            router.replace('/login')
        } else {
            router.replace({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } });
        }
    }
    app.mount('#app');
}
