// stores/user.js
import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
    // state
    const user = ref(null)

    // actions
    const setUser = (userData) => {
        user.value = userData
    }
    const clearUser = () => {
        user.value = null
    }

    // getters
    const isLoggedIn = computed(() => !!user.value)
    const userId = computed(() => user.value?.id ?? null)
    const username = computed(() => user.value?.username ?? '')

    return {
        // state
        user,

        //actions
        setUser,
        clearUser,

        // getters
        isLoggedIn,
        userId,
        username
    }
})
