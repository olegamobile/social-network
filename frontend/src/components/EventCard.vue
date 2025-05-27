<template>
    <div class="event-card">
        <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-semibold text-nordic-dark mb-4">{{ event.title }}</h2>

        <div class="space-y-2">
          
        <p>
           <span class="font-semibold text-nordic-dark">Group:</span>
           <span class="text-nordic-light ml-1">{{ event.group }}</span> 
        </p>
        <p>
            <span class="font-semibold text-nordic-dark">Time:</span> 
            <span class="text-nordic-light ml-1">{{ event.time }}</span>
        </p>
        <p>
            <span class="font-semibold text-nordic-dark">Description:</span> 
            <span class="text-nordic-light ml-1">{{ event.description }}</span>
        </p>
        <p>
            <span class="font-semibold text-nordic-dark">Going:</span> 
            <span class="text-nordic-light ml-1">{{ event.going.join(', ') }}</span>
        </p>
        </div>
        <!-- <p><strong>Not Going:</strong> {{ event.notGoing.join(', ') }}</p>
        <p><strong>Pending:</strong> {{ event.pending.join(', ') }}</p> -->

        <div class="flex flex-row gap-4 mt-4">
        <button type="button" @click="$emit('edit', event.id)"
        class="bg-nordic-primary-accent text-white font-medium px-5 py-2 rounded-md hover:bg-nordic-secondary-accent transition">Going</button>
        <button type="button" @click="$emit('notGoing', event.id)"
        class="bg-nordic-secondary-bg text-nordic-dark font-medium px-5 py-2 rounded-md hover:bg-nordic-secondary-accent transition">Not Going</button>
        </div>
    </div>
    </div>
</template>

<script setup>
defineProps({ event: Object })
</script>
