<template>
    <li class="notification-item">
        <div class="flex justify-between space-x-1 items-center">
            <div class="flex-grow">
                <router-link :to="'/profile/' + notification.sender_id" class="font-bold">{{ notification.sender_name
                    }}</router-link>
            </div>
            <!-- Accept/Decline Buttons -->
            <div class="flex justify-end space-x-1 items-center">
                <button @click="acceptAction(notification.id)"
                    class="w-8 h-8 flex items-center justify-center rounded text-white bg-[var(--nordic-primary-accent)] hover:bg-[var(--nordic-secondary-accent)] transition">
                    ✓
                </button>
                <button @click="declineAction(notification.id)"
                    class="w-8 h-8 flex items-center justify-center rounded border border-[var(--nordic-border-light)] text-[var(--nordic-text-light)] bg-[var(--nordic-secondary-bg)] hover:bg-[var(--nordic-border-light)] transition">
                    ×
                </button>
            </div>
        </div>
    </li>
</template>

<script setup>


// Accept props
defineProps({
    notification: {
        type: Object,
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
.notification-item {
    position: relative;
    padding: 1rem;
    margin-bottom: 0.5rem;
    background-color: #f9f9f9;
    border-left: 5px solid #888;
    border-radius: 4px;
}

</style>