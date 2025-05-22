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
        
        <button 
            @click="toggleComments" 
            class="mt-2 mr-2 text-blue-500 hover:underline"
            :disabled="commentCount === 0"
            :class="{ 'text-gray-400 hover:no-underline cursor-default': commentCount === 0 }"
        >
            {{ commentLabel }}
        </button>

        <button @click="showNewCommentForm" class="mt-2 text-blue-500 hover:underline">
            {{ newComment ? 'Hide form' : 'Add comment' }}
        </button>
        
        <div v-if="newComment" class="space-y-4 bg-white p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
            <h2 class="text-xl font-semibold text-gray-800">Create a new comment</h2>
            <textarea
                v-model="content"
                rows="4"
                placeholder="Enter your comment"
                class="w-full p-3 border border-gray-300 rounded-md text-gray-700 placeholder-gray-400 
                focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y"
            ></textarea>
            <button
                @click="submitComment"
                :disabled="!content.trim()"
                class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md 
                bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none 
                focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium 
                disabled:opacity-50 disabled:cursor-not-allowed"
            >
                submit comment
            </button>
        </div>
        
        <div v-if="showComments" class="mt-2">
            <div v-if="loadingComments">Loading comments...</div>
            <div v-else-if="commentCount === 0">No comments yet.</div>
            <ul v-else>
                <li v-for="comment in comments" :key="comment.id" class="mt-1 border-t pt-1">
                    <strong>{{ comment.user.first_name }} {{ comment.user.last_name }}</strong>: {{ comment.content }}
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useUserStore } from '@/stores/user';
import { storeToRefs } from 'pinia';

const { post } = defineProps({
    post: {
        type: Object,
        required: true
    }
})

const comments = ref([]);
const content = ref('');
const showComments = ref(false);
const newComment = ref(false);
const loadingComments = ref(false);
const apiUrl = import.meta.env.VITE_API_URL || '/api';
const userStore = useUserStore();
const { user } = storeToRefs(userStore);

// Track comment count separately
const commentCount = ref(post.numberOfComments || 0);

// Update comment count when comments are loaded
watch(comments, (newComments) => {
    commentCount.value = newComments.length;
});

const commentLabel = computed(() => {
    if (commentCount.value === 0) {
        return 'No comments yet';
    }
    return showComments.value 
        ? `Hide Comments (${commentCount.value})`
        : `Show Comments (${commentCount.value})`;
});

const formattedDate = computed(() => {
    const date = new Date(post.created_at)
    return date.toLocaleString("fi-FI", {
        dateStyle: 'medium',
        timeStyle: 'medium'
    }).replace("klo ", "")
})

const showNewCommentForm = async() => {
    newComment.value = !newComment.value; 
};

const toggleComments = async () => {
    if (commentCount.value === 0) return;
    
    showComments.value = !showComments.value;
    if (showComments.value) {
        await loadComments();
    }
};

const loadComments = async () => {
    loadingComments.value = true;
    try {
        const res = await fetch(`${apiUrl}/api/comments/show?post_id=${post.id}&user_id=${user.value.id}&type=${post.postType}`, {
            credentials: 'include',
        });
        comments.value = await res.json();
        // Update the count after loading
        commentCount.value = comments.value.length;
    } catch (error) {
        comments.value = [];
        console.error("Failed to load comments", error);
    } finally {
        loadingComments.value = false;
    }
};

const submitComment = async () => {
    if (!content.value.trim()) return;
    
    newComment.value = false;
    try {
        const res = await fetch(`${apiUrl}/api/comments/create?post_id=${post.id}&user_id=${user.value.id}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ content: content.value, type: post.postType })
        });
        const result = await res.json();
        if (result.success) {
            content.value = '';
            // Increment comment count
            commentCount.value++;
            // Reload comments if they're shown
            if (showComments.value) {
                await loadComments();
            }
        }
    } catch (error) {
        console.log("error: ", error);
    }
};
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