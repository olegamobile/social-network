import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
    state: () => ({
        user: null
    }),
    actions: {
        setUser(userData) {
            console.log("This user:", userData)
            this.user = userData
        },
        clearUser() {
            this.user = null
        }
    },
    getters: {
        isLoggedIn: (state) => !!state.user
    }
})
