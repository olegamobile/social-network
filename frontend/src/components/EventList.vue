<template>
    <ul class="event-list">
        <h4>Events</h4>
        <li v-for="event in events" :key="event.id">
            {{ event.title }} — {{ formatDate(event.datetime) }}
        </li>
    </ul>
</template>

<script setup>
defineProps({ events: Array })

const formatDate = (isoString) => {
    const date = new Date(isoString)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'short'
    }).replace("klo ", "")
}

</script>
