<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-md mb-6">
      {{ error }}
      <button @click="fetchSite" class="underline ml-2">Try again</button>
    </div>

    <!-- Site details -->
    <div v-else-if="site" class="space-y-8">
      <!-- Header with site info and status -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div class="flex flex-col md:flex-row md:justify-between md:items-center">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ site.name }}</h1>
            <a 
              :href="site.url" 
              target="_blank" 
              class="text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400"
            >
              {{ site.url }}
            </a>
          </div>
          
          <div class="mt-4 md:mt-0 flex items-center">
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
                'text-green-600 dark:text-green-400': siteStatus === 'up',
                'text-red-600 dark:text-red-400': siteStatus === 'down',
                'text-yellow-600 dark:text-yellow-400': siteStatus === 'checking'
              }"
            >
              {{ statusText }}
            </span>
          </div>
        </div>
      </div>

      <!-- Uptime statistics -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Uptime Statistics</h2>
        
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">Today</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ stats.today }}%</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">7 Days</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ stats.week }}%</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">30 Days</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ stats.month }}%</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">All Time</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ stats.allTime }}%</p>
          </div>
        </div>
        
        <!-- Uptime chart -->
        <div class="h-64 bg-gray-50 dark:bg-gray-700 rounded-md p-4">
          <canvas ref="uptimeChart"></canvas>
        </div>
      </div>

      <!-- Response time -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Response Time</h2>
        
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">Current</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ responseTime.current }}ms</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">Average</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ responseTime.average }}ms</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">Minimum</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ responseTime.min }}ms</p>
          </div>
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-md">
            <p class="text-sm text-gray-500 dark:text-gray-400">Maximum</p>
            <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ responseTime.max }}ms</p>
          </div>
        </div>
        
        <!-- Response time chart -->
        <div class="h-64 bg-gray-50 dark:bg-gray-700 rounded-md p-4">
          <canvas ref="responseTimeChart"></canvas>
        </div>
      </div>

      <!-- Recent incidents -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Recent Incidents</h2>
        
        <div v-if="incidents.length === 0" class="text-gray-600 dark:text-gray-400 text-center py-8">
          No incidents recorded in the last 30 days.
        </div>
        
        <div v-else class="space-y-4">
          <div 
            v-for="(incident, index) in incidents" 
            :key="index" 
            class="border-l-4 border-red-500 pl-4 py-2"
          >
            <div class="flex justify-between items-start">
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">{{ incident.statusCode }} Error</h3>
                <p class="text-sm text-gray-600 dark:text-gray-400">{{ incident.message }}</p>
              </div>
              <div class="text-right">
                <p class="text-sm text-gray-900 dark:text-white">{{ formatDate(incident.startTime) }}</p>
                <p class="text-xs text-gray-600 dark:text-gray-400">Duration: {{ incident.duration }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex space-x-4">
        <button 
          @click="showEditSiteModal = true" 
          class="bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-4 rounded-md shadow-sm transition duration-300"
        >
          Edit Site
        </button>
        <button 
          @click="showDeleteConfirmModal = true" 
          class="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-md shadow-sm transition duration-300"
        >
          Delete Site
        </button>
      </div>
    </div>

    <!-- Edit Site Modal -->
    <EditSiteModal 
      v-if="showEditSiteModal && site" 
      :site="site"
      @close="showEditSiteModal = false"
      @site-updated="onSiteUpdated"
    />

    <!-- Delete Confirmation Modal -->
    <ConfirmModal
      v-if="showDeleteConfirmModal"
      title="Delete Site"
      :message="`Are you sure you want to delete '${site?.name}'? This action cannot be undone.`"
      confirm-text="Delete"
      cancel-text="Cancel"
      confirm-variant="danger"
      @confirm="deleteSite"
      @cancel="showDeleteConfirmModal = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useSitesStore } from '../store/sites';
import EditSiteModal from '../components/sites/EditSiteModal.vue';
import ConfirmModal from '../components/common/ConfirmModal.vue';
import Chart from 'chart.js/auto';

const route = useRoute();
const router = useRouter();
const sitesStore = useSitesStore();

// State
const loading = ref(true);
const error = ref(null);
const site = computed(() => sitesStore.getCurrentSite);
const siteStatus = ref('checking');
const showEditSiteModal = ref(false);
const showDeleteConfirmModal = ref(false);
const uptimeChart = ref(null);
const responseTimeChart = ref(null);
const uptimeChartInstance = ref(null);
const responseTimeChartInstance = ref(null);

// Mock data for demo (would be replaced with real API data)
const stats = ref({
  today: 100,
  week: 99.8,
  month: 99.5,
  allTime: 99.7
});

const responseTime = ref({
  current: 187,
  average: 210,
  min: 150,
  max: 450
});

const incidents = ref([
  {
    statusCode: 503,
    message: 'Service Unavailable',
    startTime: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), // 7 days ago
    duration: '15 minutes'
  },
  {
    statusCode: 500,
    message: 'Internal Server Error',
    startTime: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000), // 15 days ago
    duration: '8 minutes'
  }
]);

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
const fetchSite = async () => {
  const siteId = route.params.id;
  
  if (!siteId) {
    router.push('/dashboard');
    return;
  }
  
  loading.value = true;
  error.value = null;
  
  try {
    await sitesStore.fetchSite(siteId);
    
    // Simulate status check (would be replaced with real API call)
    setTimeout(() => {
      siteStatus.value = Math.random() > 0.1 ? 'up' : 'down';
      loading.value = false;
    }, 1000);
  } catch (err) {
    error.value = 'Failed to load site details. Please try again.';
    loading.value = false;
  }
};

const onSiteUpdated = () => {
  showEditSiteModal.value = false;
  fetchSite();
};

const deleteSite = async () => {
  try {
    await sitesStore.deleteSite(site.value.id);
    showDeleteConfirmModal.value = false;
    router.push('/dashboard');
  } catch (err) {
    error.value = 'Failed to delete site. Please try again.';
  }
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const initCharts = () => {
  // Generate random data for charts
  const days = 30;
  const labels = Array.from({ length: days }, (_, i) => {
    const date = new Date();
    date.setDate(date.getDate() - (days - i - 1));
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  });
  
  // Uptime data (95-100%)
  const uptimeData = Array.from({ length: days }, () => 95 + Math.random() * 5);
  
  // Response time data (100-500ms)
  const responseTimeData = Array.from({ length: days }, () => 100 + Math.random() * 400);
  
  // Destroy existing charts if they exist
  if (uptimeChartInstance.value) {
    uptimeChartInstance.value.destroy();
  }
  
  if (responseTimeChartInstance.value) {
    responseTimeChartInstance.value.destroy();
  }
  
  // Create uptime chart
  uptimeChartInstance.value = new Chart(uptimeChart.value, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Uptime %',
        data: uptimeData,
        borderColor: '#10B981', // Green
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        tension: 0.3,
        fill: true
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          min: 90,
          max: 100
        }
      }
    }
  });
  
  // Create response time chart
  responseTimeChartInstance.value = new Chart(responseTimeChart.value, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Response Time (ms)',
        data: responseTimeData,
        borderColor: '#6366F1', // Indigo
        backgroundColor: 'rgba(99, 102, 241, 0.1)',
        tension: 0.3,
        fill: true
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false
    }
  });
};

// Lifecycle hooks
onMounted(() => {
  fetchSite();
});

// Initialize charts when site data is loaded
watch(() => loading.value, (newVal) => {
  if (!newVal && site.value) {
    // Wait for DOM to update
    setTimeout(() => {
      initCharts();
    }, 100);
  }
});
</script>
