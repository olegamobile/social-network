<template>
    <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
        <div
            class="bg-[var(--nordic-primary-bg)] border border-[var(--nordic-border-light)] p-6 rounded-lg shadow-lg max-w-sm w-full">
            <h3 class="text-[var(--nordic-text-dark)] text-lg font-semibold mb-3">{{ title }}</h3>
            <p class="text-[var(--nordic-text-light)] mb-5">{{ message }}</p>
            <div class="flex justify-end space-x-3">
                <button @click="cancel"
                    class="px-4 py-2 rounded border border-[var(--nordic-border-light)] text-[var(--nordic-text-light)] bg-[var(--nordic-secondary-bg)] hover:bg-[var(--nordic-border-light)] transition">
                    Cancel
                </button>
                <button @click="confirm"
                    class="px-4 py-2 rounded text-white bg-[var(--nordic-primary-accent)] hover:bg-[var(--nordic-secondary-accent)] transition">
                    Confirm
                </button>
            </div>
        </div>
    </div>
</template>


<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
    visible: Boolean,
    title: String,
    message: String
})

const emit = defineEmits(['confirm', 'cancel'])

function confirm() {
    emit('confirm')
}
function cancel() {
    emit('cancel')
}
</script>
