<template>
    <div class="event-card">
        <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-2xl font-semibold text-nordic-dark mb-4">{{ event.title }}</h2>

            <div class="space-y-2">

                <p>
                    <span class="font-semibold text-nordic-dark">Creator:</span>
                    <span class="text-nordic-light ml-1">{{ event.creator?.first_name + ' ' + event.creator?.last_name
                    }}</span>
                </p>
                <p>
                    <span class="font-semibold text-nordic-dark">Group:</span>
                    <span class="text-nordic-light ml-1">{{ event.group }}</span>
                </p>
                <p>
                    <span class="font-semibold text-nordic-dark">Time:</span>
                    <span class="text-nordic-light ml-1">{{ formatDate(event.event_datetime) }}</span>
                </p>
                <p>
                    <span class="font-semibold text-nordic-dark">Description:</span>
                    <span class="text-nordic-light ml-1">{{ event.description }}</span>
                </p>

                <br>

                <!-- people going -->
                <p>
                    <span class="font-semibold text-nordic-dark">Going:</span>
                <ul v-if="event.going && event.going.length > 0">
                    <li v-for="user in event.going" :key="user.id"
                        class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                        <RouterLink :to="`/profile/${user.id}`" class="flex items-center gap-2">
                            {{ user.first_name }} {{ user.last_name }}
                            <div v-if="user.avatar_url"
                                class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1">
                                <img :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar"
                                    class="w-full h-full object-cover" />
                            </div>
                        </RouterLink>
                    </li>
                </ul>
                </p>

                <!-- people not going -->
                <p>
                    <span class="font-semibold text-nordic-dark">Not Going:</span>
                <ul v-if="event.not_going && event.not_going.length > 0">
                    <li v-for="user in event.not_going" :key="user.id"
                        class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                        <RouterLink :to="`/profile/${user.id}`" class="flex items-center gap-2">
                            {{ user.first_name }} {{ user.last_name }}
                            <div v-if="user.avatar_url"
                                class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1">
                                <img :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar"
                                    class="w-full h-full object-cover" />
                            </div>
                        </RouterLink>
                    </li>
                </ul>
                </p>

                <!-- people who haven't answered -->
                <p>
                    <span class="font-semibold text-nordic-dark">No Response:</span>
                    <!-- <span class="text-nordic-light ml-1">{{ event.going.join(', ') }}</span> -->
                <ul v-if="event.no_response && event.no_response.length > 0">
                    <li v-for="user in event.no_response" :key="user.id"
                        class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                        <RouterLink :to="`/profile/${user.id}`" class="flex items-center gap-2">
                            {{ user.first_name }} {{ user.last_name }}
                            <div v-if="user.avatar_url"
                                class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1">
                                <img :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar"
                                    class="w-full h-full object-cover" />
                            </div>
                        </RouterLink>
                    </li>
                </ul>
                </p>
            </div>

            <div class="flex flex-row gap-4 mt-4">
                <button type="button" @click="$emit('edit', event.id)"
                    class="bg-nordic-primary-accent text-white font-medium px-5 py-2 rounded-md hover:bg-nordic-secondary-accent transition">Going</button>
                <button type="button" @click="$emit('notGoing', event.id)"
                    class="bg-nordic-secondary-bg text-nordic-dark font-medium px-5 py-2 rounded-md hover:bg-nordic-secondary-accent transition">Not
                    Going</button>
            </div>
        </div>
    </div>
</template>

<script setup>
const apiUrl = import.meta.env.VITE_API_URL || '/api'

defineProps({ event: Object })

const formatDate = (isoString) => {
    const date = new Date(isoString)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'short'
    }).replace("klo ", "")
}
</script>
