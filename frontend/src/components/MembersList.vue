<template>
    <ul v-if="members && members.length > 0">
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Members</h3>
        <li v-for="member in members" :key="member.id"
            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer mb-1">
            <RouterLink :to="`/profile/${member.id}`" class="flex items-center gap-2">
                <span class="break-all">{{ member.username }}</span>
                <div v-if="member.avatar_url" class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1">
                    <img :src="`${apiUrl}/${member.avatar_url}`" alt="User Avatar" class="w-full h-full object-cover" />
                </div>
            </RouterLink>
        </li>
    </ul>
</template>

<script setup>

const apiUrl = import.meta.env.VITE_API_URL

const { members } = defineProps({
    members: {
        type: Array,
        required: true
    }
});
</script>
