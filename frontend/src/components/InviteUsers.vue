<template>
    <div class="space-y-4 bg-white p-4 border rounded-md shadow-sm w-full max-w-screen-lg">
        <h2 class="text-xl font-semibold text-gray-800">Invite Users</h2>

        <SearchBox @results="handleResults" />
        <SearchResults :results="searchResults" :searchInitiated="searchInitiated" />
    </div>
</template>

<script setup>
import { ref } from 'vue'
import SearchBox from './SearchBox.vue'
import SearchResults from './SearchResults.vue'

const emit = defineEmits(['post-submitted'])
const searchResults = ref([])
const searchInitiated = ref(false)

const { members } = defineProps({
    members: {
        type: Array,
        required: true
    }
});

function handleResults(results) {
    searchInitiated.value = true
    searchResults.value = results.filter(user => {
        const isMember = members.some(member => member.id === user.id)
        return !isMember
    });
    //console.log("members:", members)
    //console.log("search res:", searchResults.value)
}

</script>