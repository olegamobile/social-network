
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useErrorStore = defineStore('error', () => {
    const title = ref('')
    const message = ref('')

    function setError(newTitle, newMessage) {
        title.value = newTitle
        message.value = newMessage
    }

    function clearError() {
        title.value = ''
        message.value = ''
    }

    return { title, message, setError, clearError }
})
