import { computed } from 'vue';

export const useFormats = () => {
    const currentTime = computed(() => {
        const now = new Date();
        now.setSeconds(0, 0); // Remove seconds and milliseconds
        return now.toISOString().slice(0, 16);
    })

    const finnishTime = (isoString) => {

    }

    return {
        currentTime,
        finnishTime,
    }
}