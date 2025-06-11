<template>
    <div v-if="members && members.length > 0">
        <!-- Admin section -->
        <div v-if="adminMember">
            <h3 class="text-xl font-semibold text-nordic-dark mb-3">Admin</h3>
            <div
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer mb-1">
                <RouterLink :to="`/profile/${adminMember.id}`" class="flex items-center gap-2">
                    <!-- Avatar first -->
                    <div v-if="adminMember.avatar_url" class="post-user-avatar w-6 h-6 rounded-full overflow-hidden">
                        <img :src="`${apiUrl}/${adminMember.avatar_url}`" alt="User Avatar" class="w-full h-full object-cover" />
                    </div>
                    <!-- Name -->
                    <span class="break-all">{{ adminMember.username }}</span>
                </RouterLink>
            </div>
        </div>

        <!-- Members section -->
        <h3 class="text-xl font-semibold text-nordic-dark mt-5 mb-3">Members</h3>
        <ul>
            <li v-for="member in regularMembers" :key="member.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer mb-1">
                <RouterLink :to="`/profile/${member.id}`" class="flex items-center gap-2">
                    <!-- Avatar first -->
                    <div v-if="member.avatar_url" class="post-user-avatar w-6 h-6 rounded-full overflow-hidden">
                        <img :src="`${apiUrl}/${member.avatar_url}`" alt="User Avatar" class="w-full h-full object-cover" />
                    </div>
                    <!-- Name -->
                    <span class="break-all">{{ member.username }}</span>
                </RouterLink>
            </li>
        </ul>
    </div>
</template>

<script setup>
const apiUrl = import.meta.env.VITE_API_URL
import { computed } from 'vue'

const { members } = defineProps({
    members: {
        type: Array,
        required: true
    }
});

// Split members into admin + regular members
const adminMember = computed(() => members.find(m => m.is_admin));
const regularMembers = computed(() => members.filter(m => !m.is_admin));
</script>
