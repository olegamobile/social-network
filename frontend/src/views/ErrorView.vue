<template>
  <div class="error-view">
    <TopBar />
    <h1 class="mt-20">{{ title }}</h1>
    <p>{{ message }}</p>
  </div>
</template>

<script setup>
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useErrorStore } from '../stores/error'

const { title, message } = useErrorStore()
</script>

<style scoped>
.error-view {
  text-align: center;
}
</style>