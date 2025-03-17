<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-md mx-4">
      <div class="p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Edit Site</h2>
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
          <!-- Site Name -->
          <div class="mb-4">
            <label for="siteName" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Site Name
            </label>
            <input
              id="siteName"
              v-model="siteName"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          
          <!-- Site URL -->
          <div class="mb-6">
            <label for="siteUrl" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Site URL
            </label>
            <input
              id="siteUrl"
              v-model="siteUrl"
              type="url"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Must include http:// or https://
            </p>
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
              <span v-if="loading">Saving...</span>
              <span v-else>Save Changes</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useSitesStore } from '../../store/sites';

const props = defineProps({
  site: {
    type: Object,
    required: true
  }
});

const sitesStore = useSitesStore();

const siteName = ref('');
const siteUrl = ref('');
const loading = ref(false);
const error = ref('');

onMounted(() => {
  // Initialize form with site data
  siteName.value = props.site.name;
  siteUrl.value = props.site.url;
});

const handleSubmit = async () => {
  // Validate URL format
  if (!isValidUrl(siteUrl.value)) {
    error.value = 'Please enter a valid URL including http:// or https://';
    return;
  }
  
  error.value = '';
  loading.value = true;
  
  try {
    const updatedSite = await sitesStore.updateSite(props.site.id, {
      name: siteName.value,
      url: siteUrl.value
    });
    
    // Emit event to parent component
    emit('site-updated', updatedSite);
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to update site. Please try again.';
  } finally {
    loading.value = false;
  }
};

const isValidUrl = (url) => {
  try {
    new URL(url);
    return true;
  } catch (e) {
    return false;
  }
};

const emit = defineEmits(['close', 'site-updated']);
</script>
