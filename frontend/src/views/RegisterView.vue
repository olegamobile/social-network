<template>
    <div class="register">
        <TopBar />

        <h2>Register</h2>
        <form @submit.prevent="register">
            <input v-model="form.email" type="email" placeholder="Email" required autocomplete="email" />
            <input v-model="form.password" type="password" placeholder="Password" required
                autocomplete="new-password" />
            <input v-model="form.firstName" type="text" placeholder="First Name" required />
            <input v-model="form.lastName" type="text" placeholder="Last Name" required />
            <input v-model="form.dob" type="date" placeholder="Date of Birth" required />

            <input v-model="form.nickname" type="text" placeholder="Nickname (Optional)" />
            <textarea v-model="form.about" placeholder="About Me (Optional)"></textarea>

            <label>
                Avatar/Image (Optional):
                <input type="file" @change="handleFileUpload" accept="image/*" />
            </label>

            <button type="submit">Register</button>
            <p v-if="error">{{ error }}</p>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import TopBar from '@/components/TopBar.vue'; // Assuming TopBar is in components

const form = ref({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    nickname: '',
    about: '',
    avatar: null,
});

const error = ref(null);
const router = useRouter();
const apiUrl = import.meta.env.VITE_API_URL

const handleFileUpload = (event) => {
    form.value.avatar = event.target.files[0];
};

const register = async () => {
    try {
        const payload = {
            email: form.value.email,
            password: form.value.password,
            firstName: form.value.firstName,
            lastName: form.value.lastName,
            dob: form.value.dob,
            nickname: form.value.nickname || null,
            about: form.value.about || null,
            // Avatar will likely need to be sent as FormData
        };

        const formData = new FormData();
        for (const [key, value] of Object.entries(payload)) {
            formData.append(key, value);
        }
        if (form.value.avatar) {
            formData.append('avatar', form.value.avatar);
        }

        // Replace this with your actual API endpoint:
        const response = await fetch(`${apiUrl}/api/register`, {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Registration failed');
        }

        // Handle success (e.g., redirect to login)
        router.push('/login');
    } catch (err) {
        error.value = err.message || 'An error occurred';
    }
};
</script>

<style scoped>
/* Add your styles here */
.register {
    display: flex;
    flex-direction: column;
}

h2 {
    padding-left: 2rem;
}

form {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.4rem;
    padding: 0 2rem;
}
</style>