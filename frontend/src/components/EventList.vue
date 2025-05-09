<template>
    <ul class="event-list">
        <li v-for="event in events" :key="event.id">
            {{ event.title }} — {{ event.time }}
        </li>
    </ul>
</template>

<script setup>
defineProps({ events: Array })
</script>
