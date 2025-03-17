<template>
  <div class="max-w-md mx-auto my-12 p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md">
    <h1 class="text-2xl font-bold text-center mb-6 text-gray-900 dark:text-white">Create an Account</h1>
    
    <!-- Alert for errors -->
    <div v-if="error" class="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
      {{ error }}
    </div>
    
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <!-- Username -->
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Username
        </label>
        <input
          id="username"
          v-model="username"
          type="text"
          required
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
          placeholder="Choose a username"
        />
      </div>
      
      <!-- Email (Optional) -->
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Email (Optional for downtime alerts)
        </label>
        <input
          id="email"
          v-model="email"
          type="email"
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
          placeholder="Enter your email (optional)"
        />
      </div>
      
      <!-- Password -->
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Password
        </label>
        <input
          id="password"
          v-model="password"
          type="password"
          required
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
          placeholder="Choose a password"
        />
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          Password must be at least 8 characters long
        </p>
      </div>
      
      <!-- Confirm Password -->
      <div>
        <label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Confirm Password
        </label>
        <input
          id="confirmPassword"
          v-model="confirmPassword"
          type="password"
          required
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
          placeholder="Confirm your password"
        />
      </div>
      
      <!-- Submit Button -->
      <div>
        <button
          type="submit"
          :disabled="loading || !isFormValid"
          class="w-full py-2 px-4 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition duration-300 disabled:opacity-50"
        >
          <span v-if="loading">Creating account...</span>
          <span v-else>Register</span>
        </button>
      </div>
    </form>
    
    <div class="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
      Already have an account?
      <router-link to="/login" class="text-primary-600 hover:text-primary-500 font-medium">
        Log In
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const error = ref('');

const isFormValid = computed(() => {
  return (
    username.value.trim().length > 0 &&
    password.value.length >= 8 &&
    password.value === confirmPassword.value
  );
});

const handleSubmit = async () => {
  if (!isFormValid.value) {
    if (password.value.length < 8) {
      error.value = 'Password must be at least 8 characters long';
    } else if (password.value !== confirmPassword.value) {
      error.value = 'Passwords do not match';
    }
    return;
  }
  
  error.value = '';
  loading.value = true;
  
  try {
    await authStore.register(
      username.value,
      password.value,
      email.value || null // Send null if email is empty
    );
    
    // Redirect to dashboard
    router.push('/dashboard');
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to create account. Please try again.';
  } finally {
    loading.value = false;
  }
};
</script>
