<template>
  <div class="error-view">
    <TopBar />
    <TwoColumnLayout>
      <template #main>
        <h1>{{ title }}</h1>
        <p>{{ message }}</p>
      </template>
    </TwoColumnLayout>
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