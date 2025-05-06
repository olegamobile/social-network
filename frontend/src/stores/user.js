// stores/user.js
import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
    const user = ref(null)

    const setUser = (userData) => {
        user.value = userData
    }

    const clearUser = () => {
        user.value = null
    }

    const isLoggedIn = computed(() => !!user.value)

    return {
        user,
        setUser,
        clearUser,
        isLoggedIn
    }
})
