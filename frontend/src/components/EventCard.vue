<template>
    <div class="event-card">
        <h2 class="text-2xl font-bold mb-4">{{ event.title }}</h2>
        <p><strong>Group:</strong> {{ event.group }}</p>
        <p><strong>Time:</strong> {{ event.time }}</p>
        <p><strong>Description:</strong> {{ event.description }}</p>
        <p><strong>Going:</strong> {{ event.going.join(', ') }}</p>
        <button>Going</button>
        <button>Not Going</button>
    </div>
</template>

<script setup>
defineProps({ event: Object })
</script>
