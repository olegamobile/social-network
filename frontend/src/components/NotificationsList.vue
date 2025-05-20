<template>
    <ul class="notifications-list">
        <li v-for="notification in notifications" :key="notification.id" class="notification-item">
            <div class="notification-content">
                <span class="notification-text">
                    {{ notification.content }}
                </span>

                <!-- Close Button -->
                <button class="close-button" @click="closeNotification(notification.id)" v-if="!notification.is_read">
                    &times;
                </button>
            </div>

            <!-- Accept/Decline Buttons -->
            <div class="action-buttons" v-if="notification.pending">
                <button class="accept-button" @click="acceptAction(notification.id)">
                    {{ notification.type === 'event_creation' ? 'Going' : 'Accept' }}
                </button>
                <button class="decline-button" @click="declineAction(notification.id)">
                    {{ notification.type === 'event_creation' ? 'Not going' : 'Decline' }}
                </button>
            </div>
        </li>
    </ul>
</template>


<script setup>
import { defineEmits } from 'vue'

// Accept props
defineProps({
  notifications: {
    type: Array,
    required: true
  }
})

// Emit events to the parent component
const emit = defineEmits(['close', 'accept', 'decline'])

// These functions emit events to the parent
function closeNotification(id) {
  emit('close', id)
}

function acceptAction(id) {
  emit('accept', id)
}

function declineAction(id) {
  emit('decline', id)
}
</script>

<style scoped>
.notifications-list {
    list-style: none;
    padding: 0;
}

.notification-item {
    position: relative;
    padding: 1rem;
    margin-bottom: 0.5rem;
    background-color: #f9f9f9;
    border-left: 5px solid #888;
    border-radius: 4px;
}

.notification-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.notification-text {
    flex: 1;
    color: var(--nordic-text-dark);
}

/* Close Button (X) */
.close-button {
    position: absolute;
    top: 0.5rem;
    right: 0.75rem;
    background: none;
    border: none;
    font-size: 1.5rem;
    color: var(--nordic-text-light);
    cursor: pointer;
    transition: color 0.2s ease-in-out;
}

.close-button:hover {
    color: var(--nordic-primary-accent);
}

/* Invitation Action Buttons */
.action-buttons {
    margin-top: 0.75rem;
    display: flex;
    gap: 0.5rem;
}

.accept-button,
.decline-button {
    padding: 0.4rem 1rem;
    font-size: 0.9rem;
    border-radius: 20px;
    border: none;
    cursor: pointer;
    transition: background-color 0.2s ease-in-out;
}

.accept-button {
    background-color: var(--nordic-primary-accent);
    color: white;
}

.accept-button:hover {
    background-color: var(--nordic-secondary-accent);
}

.decline-button {
    background-color: var(--nordic-secondary-bg);
    color: var(--nordic-text-dark);
}

.decline-button:hover {
    background-color: #dbe4ec;
}
</style>
