<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-md mx-4">
      <div class="p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Add Custom Domain</h2>
          <button @click="$emit('close')" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Error message -->
        <div v-if="error" class="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
          {{ error }}
        </div>
        
        <form @submit.prevent="handleSubmit">
          <!-- Domain Name -->
          <div class="mb-4">
            <label for="domain" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Domain Name
            </label>
            <input
              id="domain"
              v-model="domain"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
              placeholder="status.example.com"
            />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Enter your domain without http:// or https://
            </p>
          </div>
          
          <!-- Instructions -->
          <div class="mb-6 p-3 bg-gray-50 dark:bg-gray-700 rounded-md">
            <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              How to set up your custom domain:
            </h3>
            <ol class="text-xs text-gray-600 dark:text-gray-400 list-decimal pl-4 space-y-1">
              <li>Add a CNAME record in your DNS settings</li>
              <li>Point it to <code class="bg-gray-200 dark:bg-gray-600 px-1 py-0.5 rounded">{{ window.location.hostname }}</code></li>
              <li>After adding the domain here, click "Verify" to confirm ownership</li>
            </ol>
          </div>
          
          <!-- Submit and Cancel buttons -->
          <div class="flex justify-end space-x-3">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
            >
              <span v-if="loading">Adding...</span>
              <span v-else>Add Domain</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useDomainsStore } from '../../store/domains';

const domainsStore = useDomainsStore();

const domain = ref('');
const loading = ref(false);
const error = ref('');

const handleSubmit = async () => {
  // Validate domain format
  if (!isValidDomain(domain.value)) {
    error.value = 'Please enter a valid domain name (e.g., status.example.com)';
    return;
  }
  
  error.value = '';
  loading.value = true;
  
  try {
    const newDomain = await domainsStore.createDomain({
      domain: domain.value
    });
    
    // Emit event to parent component
    emit('domain-added', newDomain);
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to add domain. Please try again.';
  } finally {
    loading.value = false;
  }
};

const isValidDomain = (domain) => {
  // Simple domain validation
  const domainRegex = /^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/;
  return domainRegex.test(domain);
};

const emit = defineEmits(['close', 'domain-added']);
</script>
