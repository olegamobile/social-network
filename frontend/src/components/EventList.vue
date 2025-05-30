<template>
    <ul v-if="events && events.length > 0" class="space-y-2">
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Events</h3>
        <li v-for="event in events" :key="event.id" class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
            <RouterLink :to="`/events/${event.id}`">
            <p><strong>{{ event.title }}</strong><br>{{ formatDate(event.event_datetime) }}</p>
            </RouterLink>
        </li>        
    </ul>
</template>

<script setup>
const { events } = defineProps({
    events: {
        type: Array,
        required: true
    }
});

const formatDate = (isoString) => {
    const date = new Date(isoString)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'short'
    }).replace("klo ", "")
}

</script>
