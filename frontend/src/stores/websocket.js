import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useWebSocketStore = defineStore('websocket', () => {
    const socket = ref(null)
    const isConnected = ref(false)
    const message = ref(null)
    const reconnectAttempts = ref(0)
    const maxReconnectAttempts = 5
    const reconnectTimeout = ref(null)

    function connect(url) {
        if (socket.value) return // already connected or connectingß

        clearTimeout(reconnectTimeout.value)
        socket.value = new WebSocket(url)

        socket.value.onopen = () => {
            isConnected.value = true
            reconnectAttempts.value = 0
            console.log('WebSocket connected')

            // Send a ping to keep the connection alive
            send({
                type: 'ping',
                content: 'ping',
                sender_id: localStorage.getItem('userId') || '0',
                receiver_id: '0'
            })
        }

        socket.value.onmessage = (event) => {
            try {
                // Parse the incoming message
                message.value = JSON.parse(event.data)

                // Handle ping response if needed
                if (message.value.type === 'pong') {
                    console.log('Received pong from server')
                }

            } catch (error) {
                console.error('Error parsing message:', error)
                message.value = event.data
            }
        }

        socket.value.onclose = () => {
            isConnected.value = false
            socket.value = null
            console.log('WebSocket closed')

            // Try to reconnect if not maxed out
            if (reconnectAttempts.value < maxReconnectAttempts) {
                reconnectAttempts.value++
                const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.value), 30000) // Exponential backoff with 30s max
                console.log(`Attempting to reconnect in ${delay}ms (attempt ${reconnectAttempts.value})`)

                reconnectTimeout.value = setTimeout(() => {
                    connect(url)
                }, delay)
            }
        }

        socket.value.onerror = (error) => {
            console.error('WebSocket error:', error)
        }
    }

    function send(data) {
        if (!isConnected.value || !socket.value) {
            console.warn('Cannot send message: WebSocket not connected')
            return false
        }

        try {
            // Ensure all numeric values are sent as strings to avoid parsing issues in Go
            const processedData = Object.entries(data).reduce((acc, [key, value]) => {
                // Convert any numbers to strings
                if (typeof value === 'number') {
                    acc[key] = value.toString()
                } else {
                    acc[key] = value
                }
                return acc
            }, {})

            socket.value.send(JSON.stringify(processedData))
            return true
        } catch (error) {
            console.error('Error sending message:', error)
            return false
        }
    }

    function disconnect() {
        if (socket.value && isConnected.value) {
            socket.value.close()
        }

        // Clear any pending reconnect attempts
        clearTimeout(reconnectTimeout.value)
        reconnectAttempts.value = maxReconnectAttempts // Prevent further reconnects
    }

    return { connect, send, disconnect, isConnected, message }
})