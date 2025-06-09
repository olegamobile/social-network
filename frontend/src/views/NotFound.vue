<template>
    <div class="not-found">
        <TopBar />
        <h1 class="mt-20">404 - Page Not Found</h1>
        <p>The page you're looking for doesn't exist.</p>
    </div>
</template>

<script setup>
import TopBar from '@/components/TopBar.vue'

// No actual 404 status, that should come from the server.
</script>

<style scoped>
.not-found {
    text-align: center;
}
</style>