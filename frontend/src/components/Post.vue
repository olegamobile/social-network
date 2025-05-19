<template>
    <div class="post-card">
        <p class="post-content">{{ post.content }}</p>
        <small class="post-date flex items-center">
            <RouterLink :to="`/profile/${post.user_id}`" class="post-user flex items-center mr-1">
                <div v-if="post.avatar_url" class="post-user-avatar w-6 h-6 rounded-full overflow-hidden mr-1">
                    <img :src="`${apiUrl}/${post.avatar_url}`" alt="User Avatar" class="w-full h-full object-cover" />
                </div>
                {{ post.username }}
            </RouterLink>
            <span v-if="post.group_id">
                in
                <RouterLink :to="`/groups/${post.group_id}`" class="text-blue-500 hover:underline mr-1">
                    {{ post.group_name }}
                </RouterLink>
            </span>
            on {{ formattedDate }}
        </small>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const apiUrl = import.meta.env.VITE_API_URL

const { post } = defineProps({
    post: {
        type: Object,
        required: true
    }
})

const formattedDate = computed(() => {
    const date = new Date(post.created_at)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'medium'
    }).replace("klo ", "")
})

</script>


<style scoped>
.post-card {
    border: 1px solid #ccc;
    padding: 1rem;
    margin-bottom: 1rem;
    border-radius: 8px;
    background-color: #fafafa;
}

.post-content {
    font-size: 1.1rem;
    margin-bottom: 0.5rem;
}

.post-date {
    color: #666;
    font-size: 0.85rem;
}

.post-user {
    color: #0077cc;
    text-decoration: none;
    font-weight: bold;
}

.post-user:hover {
    text-decoration: underline;
}
</style>