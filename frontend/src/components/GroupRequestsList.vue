<template>
    <ul class="notifications-list">
        <GroupRequest 
            v-for="notification in notifications" 
            :key="notification.id" 
            :notification="notification"
            @close="closeNotification"
            @accept="acceptAction"
            @decline="declineAction"
        />
    </ul>
</template>

<script setup>
import GroupRequest from './GroupRequest.vue'

// Accept props
defineProps({
    notifications: {
        type: Array,
        required: true
    }
})

// Emit events to the parent component
const emit = defineEmits(['close', 'accept', 'decline'])

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
</style>