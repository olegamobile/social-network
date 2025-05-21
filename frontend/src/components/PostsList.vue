<template>
  <div class="posts-container max-w-screen-lg">
    <Post v-for="p in posts" :key="p.id" :post="p" />

    <!-- Sentinel div for triggering infinite scroll -->
    <div ref="scrollObserver" class="h-10"></div>
  </div>
</template>

<script setup>
import Post from "./Post.vue";
import { ref, onMounted } from 'vue';

const { posts } = defineProps({
  posts: {
    type: Array,
    required: true
  }
});

const scrollObserver = ref(null);
defineExpose({ scrollObserver }); // Let HomeView.vue access it
</script>