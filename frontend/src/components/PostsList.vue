<template>
  <div class="posts-container max-w-screen-lg">
    <Post v-for="p in posts" :key="p.id" :post="p" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import Post from "./Post.vue";

const { posts } = defineProps({
  posts: {
    type: Array,
    required: true
  }
});

</script>