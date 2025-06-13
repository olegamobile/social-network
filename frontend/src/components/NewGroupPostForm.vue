<template>
    <div
        class="space-y-4 bg-[var(--nordic-primary-bg)] p-4 border border-[var(--nordic-border-light)] rounded-md shadow-sm w-full max-w-screen-lg">
        <!-- header -->
        <h2 class="text-xl font-semibold text-[var(--nordic-text-dark)]">Create a New Post</h2>

        <!-- content input -->
        <textarea v-model="content" rows="4" placeholder="Enter your message"
            class="w-full p-3 border border-[var(--nordic-border-light)] rounded-md 
             text-[var(--nordic-text-dark)] placeholder-[var(--nordic-text-light)] 
             bg-white
             focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] focus:border-[var(--nordic-secondary-accent)] resize-y"></textarea>

        <!-- upload image -->
        <div class="space-y-2">
            <label class="block text-sm font-medium text-[var(--nordic-text-light)]">Image (Optional):</label>
            <input type="file" @change="handleFileUpload" accept=".jpg, .jpeg, .gif, .png, .webp" class="block w-full text-sm text-[var(--nordic-text-light)]
               file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 
               file:text-sm file:font-semibold 
               file:bg-[var(--nordic-secondary-bg)] file:text-[var(--nordic-text-dark)] 
               hover:file:bg-[var(--nordic-border-light)]" />
        </div>

        <!-- submit button -->
        <button @click="submitPost" :disabled="!content.trim()" class="inline-flex items-center py-2 px-4 border border-nordic-light rounded-md 
    bg-nordic-primary-accent text-white hover:bg-nordic-secondary-accent focus:outline-none 
    focus:ring-2 focus:ring-nordic-secondary-accent transition font-medium 
    disabled:opacity-50 disabled:cursor-not-allowed">
            Post
        </button>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useImageProcessor } from '@/composables/useImageProcessor';

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const emit = defineEmits(['post-submitted'])
const content = ref('')
const image = ref(null)

const { group_id } = defineProps({
    group_id: {
        type: Number,
        required: true
    }

});

const handleFileUpload = (event) => {
    image.value = event.target.files[0];
};

const submitPost = async () => {
    const { processPostImage } = useImageProcessor();

    try {
        const payload = {
            content: content.value,
            group_id: group_id
        };

        const formData = new FormData();
        for (const [key, value] of Object.entries(payload)) {
            formData.append(key, value);
        }

        if (image.value) {
            const processedImg = await processPostImage(image.value);
            formData.append('image', processedImg);
        }

        const res = await fetch(`${apiUrl}/api/group-posts/create`, {
            method: 'POST',
            body: formData,
            credentials: 'include',
        })

        if (res.ok) {
            const newPost = await res.json()
            emit('post-submitted', newPost)
            content.value = ''
            image.value = null
        } else {
            alert('Failed to post. Are you logged in?')
        }
    } catch (err) {
        console.error(err)
    }
}
</script>