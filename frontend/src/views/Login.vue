<template>
  <div class="max-w-md mx-auto my-12 p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md">
    <h1 class="text-2xl font-bold text-center mb-6 text-gray-900 dark:text-white">Log In</h1>
    
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
          placeholder="Enter your username"
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
          placeholder="Enter your password"
        />
      </div>
      
      <!-- Submit Button -->
      <div>
        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2 px-4 bg-primary-600 hover:bg-primary-700 text-white font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition duration-300 disabled:opacity-50"
        >
          <span v-if="loading">Logging in...</span>
          <span v-else>Log In</span>
        </button>
      </div>
    </form>
    
    <div class="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
      Don't have an account?
      <router-link to="/register" class="text-primary-600 hover:text-primary-500 font-medium">
        Register
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');

const handleSubmit = async () => {
  error.value = '';
  loading.value = true;
  
  try {
    await authStore.login(username.value, password.value);
    
    // Redirect to dashboard or the page they were trying to access
    const redirectPath = route.query.redirect || '/dashboard';
    router.push(redirectPath);
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to log in. Please check your credentials.';
  } finally {
    loading.value = false;
  }
};
</script>
