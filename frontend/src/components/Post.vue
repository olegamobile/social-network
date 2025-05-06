<template>
    <div class="post-card">
        <p class="post-content">{{ post.content }}</p>
        <small class="post-date">
            Posted by <strong>{{ post.username }}</strong> on {{ formattedDate }}
        </small>
    </div>
</template>

<script setup>
import { computed } from 'vue'

const { post } = defineProps({
    post: {
        type: Object,
        required: true
    }
})

const formattedDate = computed(() => {
    const date = new Date(post.created_at)
    return date.toLocaleString(undefined, {
        dateStyle: 'medium',
        timeStyle: 'short'
    })
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
</style>