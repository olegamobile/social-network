<template>

    <div class="max-w-xl px-4 py-8 bg-nordic-secondary-bg shadow-md rounded-xl mb-10 border-nordic-light">
        <h2 class="text-2xl font-semibold mb-6 text-nordic-dark">Edit Profile</h2>

        <form @submit.prevent="updateProfile" class="space-y-4">
            <input v-model="form.firstName" type="text" placeholder="First Name" required
                class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />

            <input v-model="form.lastName" type="text" placeholder="Last Name" required
                class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />

            <input v-model="form.dob" type="date" required
                class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />

            <input v-model="form.nickname" type="text" placeholder="Nickname (Optional)"
                class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark" />

            <textarea v-model="form.about" placeholder="About Me (Optional)" rows="3"
                class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 resize-none text-nordic-dark"></textarea>

            <div>
                <label for="privacy" class="block text-sm font-medium text-nordic-light mb-1">Profile Privacy</label>
                <select v-model="form.isPublic" id="privacy"
                    class="w-full px-4 py-2 border border-nordic-light rounded-md focus:outline-none focus:ring focus:ring-gray-200 text-nordic-dark">
                    <option :value="true">Public</option>
                    <option :value="false">Private</option>
                </select>
            </div>

            <div class="space-y-2">
                <label class="block text-sm font-medium text-nordic-light">Avatar/Image (Optional):</label>
                <input type="file" @change="handleFileUpload" accept="image/*"
                    class="block w-full text-sm text-nordic-light file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-gray-100 file:text-nordic-light hover:file:bg-gray-200" />

                <div v-if="form.avatarUrl" class="mt-2">
                    <p class="text-sm text-nordic-light">Current Avatar:</p>
                    <img :src="`${apiUrl}/${form.avatarUrl}`" alt="Avatar"
                        class="h-24 w-24 object-cover rounded-md border border-nordic-light" />
                    <button type="button" @click="deleteAvatar"
                        class="mt-2 text-sm text-red-600 hover:underline text-nordic-primary-accent hover:text-nordic-secondary-accent">Delete
                        Avatar</button>
                </div>
                <div v-else>
                    <span class="block text-sm font-medium text-nordic-light">No current avatar</span>
                </div>
            </div>




            <button type="submit"
                class="w-full bg-nordic-primary-accent text-white py-2 px-4 rounded-md hover:bg-nordic-secondary-accent transition">Save
                Changes</button>

            <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
            <p v-if="success" class="text-green-600 text-sm">{{ success }}</p>
        </form>
    </div>

</template>


<script setup>
import { ref, onMounted } from 'vue';
import { useUserStore } from '@/stores/user'
import { useImageProcessor } from '@/composables/useImageProcessor';

const form = ref({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    nickname: '',
    about: '',
    avatar: null,
    avatarUrl: null,
    isPublic: false,
});

const error = ref(null);
const success = ref(null);
const apiUrl = import.meta.env.VITE_API_URL;
const userStore = useUserStore()

const handleFileUpload = (event) => {
    form.value.avatar = event.target.files[0];
};

const deleteAvatar = () => {
    form.value.avatar = null;
    form.value.avatarUrl = null; // Clear preview — actual delete happens on backend
};

const fetchProfile = async () => {
    try {
        const response = await fetch(`${apiUrl}/api/me`, {
            credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch profile');
        const data = await response.json();

        form.value = {
            ...form.value,
            email: data.email,
            firstName: data.first_name,
            lastName: data.last_name,
            dob: data.birthday ? data.birthday.slice(0, 10) : '', // remove time from date & time
            nickname: data.username || '',
            about: data.about_me || '',
            avatarUrl: data.avatar_url || null,
            isPublic: data.is_public ?? false,
            password: '', // Don't pre-fill password
        };

        //console.log("data from api/me in fetchP in Edit P:", data)
    } catch (err) {
        error.value = err.message;
    }
};

const updateProfile = async () => {
    const { processAvatarImage } = useImageProcessor();

    try {
        const payload = {
            email: form.value.email,
            password: form.value.password,
            firstName: form.value.firstName,
            lastName: form.value.lastName,
            dob: form.value.dob,
            nickname: form.value.nickname || null,
            about: form.value.about || null,
            is_public: form.value.isPublic,
        };

        const formData = new FormData();
        for (const [key, value] of Object.entries(payload)) {
            formData.append(key, value);
        }

        if (form.value.avatar) {
            const processedAvatar = await processAvatarImage(form.value.avatar);
            formData.append('avatar', processedAvatar);
        } else if (form.value.avatarUrl === null) {
            formData.append('delete_avatar', 'true');
        }

        const response = await fetch(`${apiUrl}/api/me/update`, {
            method: 'POST',
            body: formData,
            credentials: 'include',
        });

        if (!response.ok) {
            const errData = await response.json();
            throw new Error(errData.message || 'Update failed');
        }
        const data = await response.json();

        userStore.setUser(data);

        success.value = 'Profile updated successfully';
        error.value = null;
    } catch (err) {
        error.value = err.message || 'An error occurred';
        success.value = null;
    }
};

onMounted(fetchProfile);
</script>