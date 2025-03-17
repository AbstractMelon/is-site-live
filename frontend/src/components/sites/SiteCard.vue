<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
    <!-- Status header -->
    <div 
      class="px-4 py-3 flex justify-between items-center"
      :class="{
        'bg-green-100 dark:bg-green-900': siteStatus === 'up',
        'bg-red-100 dark:bg-red-900': siteStatus === 'down',
        'bg-yellow-100 dark:bg-yellow-900': siteStatus === 'checking'
      }"
    >
      <div class="flex items-center">
        <div 
          class="w-3 h-3 rounded-full mr-2"
          :class="{
            'bg-green-500': siteStatus === 'up',
            'bg-red-500': siteStatus === 'down',
            'bg-yellow-500': siteStatus === 'checking'
          }"
        ></div>
        <span 
          class="font-medium"
          :class="{
            'text-green-800 dark:text-green-200': siteStatus === 'up',
            'text-red-800 dark:text-red-200': siteStatus === 'down',
            'text-yellow-800 dark:text-yellow-200': siteStatus === 'checking'
          }"
        >
          {{ statusText }}
        </span>
      </div>
      <div class="flex space-x-2">
        <button 
          @click="$emit('edit')" 
          class="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white"
          title="Edit site"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button 
          @click="$emit('delete')" 
          class="text-gray-600 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400"
          title="Delete site"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
    
    <!-- Site info -->
    <div class="p-4">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">{{ site.name }}</h3>
      <a 
        :href="site.url" 
        target="_blank" 
        class="text-sm text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 break-all"
      >
        {{ site.url }}
      </a>
      
      <!-- Stats -->
      <div class="mt-4 grid grid-cols-2 gap-4">
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">Uptime</p>
          <p class="text-xl font-semibold text-gray-900 dark:text-white">
            {{ stats.uptime ? `${stats.uptime}%` : 'N/A' }}
          </p>
        </div>
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">Response Time</p>
          <p class="text-xl font-semibold text-gray-900 dark:text-white">
            {{ stats.responseTime ? `${stats.responseTime}ms` : 'N/A' }}
          </p>
        </div>
      </div>
      
      <!-- View details button -->
      <div class="mt-4">
        <router-link 
          :to="`/site/${site.id}`" 
          class="block w-full text-center bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 font-medium py-2 px-4 rounded-md transition duration-300"
        >
          View Details
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import axios from 'axios';

const props = defineProps({
  site: {
    type: Object,
    required: true
  }
});

defineEmits(['edit', 'delete']);

// Site status (simulated for now, will be replaced with real data)
const siteStatus = ref('checking');
const stats = ref({
  uptime: null,
  responseTime: null
});

// Computed properties
const statusText = computed(() => {
  switch (siteStatus.value) {
    case 'up':
      return 'Online';
    case 'down':
      return 'Offline';
    case 'checking':
      return 'Checking...';
    default:
      return 'Unknown';
  }
});

// Methods
const fetchSiteStatus = async () => {
  try {
    const response = await axios.get(`/api/sites/${props.site.id}/status`);
    siteStatus.value = response.data.status;
    stats.value = {
      uptime: response.data.uptime_percentage,
      responseTime: response.data.avg_response_time
    };
  } catch (error) {
    console.error('Failed to fetch site status:', error);
    siteStatus.value = 'unknown';
  }
};

// For demo purposes, we'll simulate a status
const simulateStatus = () => {
  // 90% chance of being up, 10% chance of being down
  siteStatus.value = Math.random() > 0.1 ? 'up' : 'down';
  stats.value = {
    uptime: siteStatus.value === 'up' ? (95 + Math.random() * 5).toFixed(2) : (85 + Math.random() * 10).toFixed(2),
    responseTime: Math.floor(100 + Math.random() * 400)
  };
};

// Polling interval
let statusInterval;

// Lifecycle hooks
onMounted(() => {
  // In a real app, we would fetch the status from the API
  // fetchSiteStatus();
  
  // For demo purposes, we'll simulate a status
  setTimeout(simulateStatus, 1000);
  
  // Set up polling (every 60 seconds in a real app, but faster for demo)
  statusInterval = setInterval(() => {
    // fetchSiteStatus();
    simulateStatus();
  }, 30000); // Every 30 seconds for demo
});

onUnmounted(() => {
  // Clean up interval
  if (statusInterval) {
    clearInterval(statusInterval);
  }
});
</script>
