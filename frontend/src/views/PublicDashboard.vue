<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-md mb-6">
      {{ error }}
      <button @click="fetchDashboard" class="underline ml-2">Try again</button>
    </div>

    <!-- Dashboard content -->
    <div v-else-if="dashboard" class="space-y-8">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div class="flex flex-col md:flex-row md:justify-between md:items-center">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ dashboard.title }}</h1>
            <p class="text-gray-600 dark:text-gray-400">
              Last updated: {{ formatDateTime(dashboard.lastUpdated) }}
            </p>
          </div>
          
          <div class="mt-4 md:mt-0">
            <div class="flex items-center">
              <div 
                class="w-3 h-3 rounded-full mr-2"
                :class="{
                  'bg-green-500': dashboard.overallStatus === 'operational',
                  'bg-red-500': dashboard.overallStatus === 'major_outage',
                  'bg-yellow-500': dashboard.overallStatus === 'partial_outage'
                }"
              ></div>
              <span 
                class="font-medium"
                :class="{
                  'text-green-600 dark:text-green-400': dashboard.overallStatus === 'operational',
                  'text-red-600 dark:text-red-400': dashboard.overallStatus === 'major_outage',
                  'text-yellow-600 dark:text-yellow-400': dashboard.overallStatus === 'partial_outage'
                }"
              >
                {{ statusText }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Sites grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="site in dashboard.sites" 
          :key="site.id" 
          class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden"
        >
          <!-- Status header -->
          <div 
            class="px-4 py-3"
            :class="{
              'bg-green-100 dark:bg-green-900': site.status === 'up',
              'bg-red-100 dark:bg-red-900': site.status === 'down',
              'bg-yellow-100 dark:bg-yellow-900': site.status === 'degraded'
            }"
          >
            <div class="flex items-center">
              <div 
                class="w-3 h-3 rounded-full mr-2"
                :class="{
                  'bg-green-500': site.status === 'up',
                  'bg-red-500': site.status === 'down',
                  'bg-yellow-500': site.status === 'degraded'
                }"
              ></div>
              <span 
                class="font-medium"
                :class="{
                  'text-green-800 dark:text-green-200': site.status === 'up',
                  'text-red-800 dark:text-red-200': site.status === 'down',
                  'text-yellow-800 dark:text-yellow-200': site.status === 'degraded'
                }"
              >
                {{ getSiteStatusText(site.status) }}
              </span>
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
                  {{ site.uptime }}%
                </p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Response Time</p>
                <p class="text-xl font-semibold text-gray-900 dark:text-white">
                  {{ site.responseTime }}ms
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Incident history -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Recent Incidents</h2>
        
        <div v-if="dashboard.incidents.length === 0" class="text-gray-600 dark:text-gray-400 text-center py-8">
          No incidents in the last 30 days.
        </div>
        
        <div v-else class="space-y-6">
          <div 
            v-for="(incident, index) in dashboard.incidents" 
            :key="index" 
            class="border-l-4 pl-4 py-2"
            :class="{
              'border-red-500': incident.severity === 'critical',
              'border-yellow-500': incident.severity === 'major',
              'border-blue-500': incident.severity === 'minor'
            }"
          >
            <div class="flex flex-col md:flex-row md:justify-between md:items-start">
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">{{ incident.title }}</h3>
                <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">{{ incident.description }}</p>
                <div class="mt-2">
                  <span 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="{
                      'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': incident.severity === 'critical',
                      'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': incident.severity === 'major',
                      'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200': incident.severity === 'minor'
                    }"
                  >
                    {{ getSeverityText(incident.severity) }}
                  </span>
                  <span 
                    class="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="{
                      'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200': incident.status === 'resolved',
                      'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200': incident.status === 'monitoring',
                      'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200': incident.status === 'investigating'
                    }"
                  >
                    {{ getIncidentStatusText(incident.status) }}
                  </span>
                </div>
              </div>
              <div class="mt-2 md:mt-0 md:ml-4 text-right">
                <p class="text-sm text-gray-900 dark:text-white">{{ formatDate(incident.startTime) }}</p>
                <p class="text-xs text-gray-600 dark:text-gray-400 mt-1">
                  {{ incident.status === 'resolved' ? `Duration: ${incident.duration}` : 'Ongoing' }}
                </p>
              </div>
            </div>
            
            <!-- Incident updates -->
            <div v-if="incident.updates && incident.updates.length > 0" class="mt-4 pl-4 border-l border-gray-200 dark:border-gray-700">
              <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Updates</h4>
              <div 
                v-for="(update, updateIndex) in incident.updates" 
                :key="updateIndex"
                class="mb-3 text-sm"
              >
                <p class="text-gray-900 dark:text-white">{{ update.message }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  {{ formatDateTime(update.timestamp) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="text-center text-sm text-gray-500 dark:text-gray-400 mt-8">
        <p>Powered by <a href="/" class="text-primary-600 dark:text-primary-400 hover:underline">Is It Live</a></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

// State
const loading = ref(true);
const error = ref(null);
const dashboard = ref(null);

// Computed properties
const statusText = computed(() => {
  if (!dashboard.value) return '';
  
  switch (dashboard.value.overallStatus) {
    case 'operational':
      return 'All Systems Operational';
    case 'partial_outage':
      return 'Partial System Outage';
    case 'major_outage':
      return 'Major System Outage';
    default:
      return 'Unknown Status';
  }
});

// Methods
const fetchDashboard = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    // In a real app, we would fetch the dashboard data from the API
    // For demo purposes, we'll use mock data
    setTimeout(() => {
      // Get username from route params
      const username = route.params.username;
      
      if (!username) {
        error.value = 'Invalid dashboard URL';
        loading.value = false;
        return;
      }
      
      // Mock dashboard data
      dashboard.value = {
        title: `${username}'s Status Dashboard`,
        lastUpdated: new Date().toISOString(),
        overallStatus: 'operational', // operational, partial_outage, major_outage
        sites: [
          {
            id: 1,
            name: 'Main Website',
            url: 'https://example.com',
            status: 'up',
            uptime: 99.98,
            responseTime: 187
          },
          {
            id: 2,
            name: 'API Service',
            url: 'https://api.example.com',
            status: 'up',
            uptime: 99.95,
            responseTime: 210
          },
          {
            id: 3,
            name: 'Documentation',
            url: 'https://docs.example.com',
            status: 'up',
            uptime: 100,
            responseTime: 156
          },
          {
            id: 4,
            name: 'Customer Portal',
            url: 'https://portal.example.com',
            status: 'up',
            uptime: 99.99,
            responseTime: 230
          }
        ],
        incidents: [
          {
            title: 'API Service Degraded Performance',
            description: 'We are investigating reports of slow response times from our API service.',
            severity: 'major', // critical, major, minor
            status: 'resolved', // investigating, monitoring, resolved
            startTime: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(), // 7 days ago
            duration: '45 minutes',
            updates: [
              {
                message: 'We are investigating reports of slow response times.',
                timestamp: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000 - 45 * 60 * 1000).toISOString()
              },
              {
                message: 'We have identified the issue as a database connection bottleneck.',
                timestamp: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000 - 30 * 60 * 1000).toISOString()
              },
              {
                message: 'The issue has been resolved and service has been fully restored.',
                timestamp: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString()
              }
            ]
          }
        ]
      };
      
      loading.value = false;
    }, 1000);
  } catch (err) {
    error.value = 'Failed to load dashboard. Please try again.';
    loading.value = false;
  }
};

const getSiteStatusText = (status) => {
  switch (status) {
    case 'up':
      return 'Operational';
    case 'down':
      return 'Down';
    case 'degraded':
      return 'Degraded';
    default:
      return 'Unknown';
  }
};

const getSeverityText = (severity) => {
  switch (severity) {
    case 'critical':
      return 'Critical';
    case 'major':
      return 'Major';
    case 'minor':
      return 'Minor';
    default:
      return 'Unknown';
  }
};

const getIncidentStatusText = (status) => {
  switch (status) {
    case 'investigating':
      return 'Investigating';
    case 'monitoring':
      return 'Monitoring';
    case 'resolved':
      return 'Resolved';
    default:
      return 'Unknown';
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
};

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// Lifecycle hooks
onMounted(() => {
  fetchDashboard();
});
</script>
