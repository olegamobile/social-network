<template>
    <ul class="notifications-list">
        <li v-for="notification in notifications" :key="notification.id" class="notification-item">
            {{ notification.text }}
        </li>
    </ul>
</template>

<script setup>
defineProps({
    notifications: {
        type: Array,
        required: true
    }
})
</script>

<style scoped>
.notifications-list {
    list-style: none;
    padding: 0;
}

.notification-item {
    padding: 1rem;
    margin-bottom: 0.5rem;
    background-color: #f9f9f9;
    border-left: 5px solid #888;
    border-radius: 4px;
}
</style>
