<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Settings</h1>
    
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <!-- Settings tabs -->
      <div class="flex border-b border-gray-200 dark:border-gray-700">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          @click="activeTab = tab.id"
          class="px-4 py-3 text-sm font-medium transition-colors duration-200"
          :class="[
            activeTab === tab.id 
              ? 'text-primary-600 border-b-2 border-primary-600 dark:text-primary-400 dark:border-primary-400' 
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
          ]"
        >
          {{ tab.name }}
        </button>
      </div>
      
      <!-- Account settings tab -->
      <div v-if="activeTab === 'account'" class="p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Account Settings</h2>
        
        <div class="space-y-6">
          <!-- API Key -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              API Key
            </label>
            <div class="flex">
              <input
                type="text"
                readonly
                :value="apiKey"
                class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-l-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
              />
              <button 
                @click="regenerateApiKey"
                class="px-4 py-2 border border-gray-300 dark:border-gray-600 border-l-0 rounded-r-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                Regenerate
              </button>
            </div>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Use this API key to access the Is It Live API programmatically
            </p>
          </div>
          
          <!-- Delete account -->
          <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
            <h3 class="text-md font-semibold text-red-600 dark:text-red-400 mb-2">Danger Zone</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
              Once you delete your account, there is no going back. Please be certain.
            </p>
            <button 
              @click="showDeleteAccountModal = true"
              class="px-4 py-2 border border-red-300 dark:border-red-700 rounded-md shadow-sm text-sm font-medium text-red-700 dark:text-red-400 bg-white dark:bg-gray-800 hover:bg-red-50 dark:hover:bg-red-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              Delete Account
            </button>
          </div>
        </div>
      </div>
      
      <!-- Appearance tab -->
      <div v-if="activeTab === 'appearance'" class="p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Appearance</h2>
        
        <div class="space-y-6">
          <!-- Theme selection -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              Theme
            </label>
            <div class="grid grid-cols-3 gap-4">
              <button 
                @click="setTheme('light')"
                class="p-4 border rounded-md flex flex-col items-center"
                :class="[
                  theme === 'light' 
                    ? 'border-primary-500 bg-primary-50 dark:bg-primary-900 dark:border-primary-400' 
                    : 'border-gray-200 dark:border-gray-700'
                ]"
              >
                <div class="w-12 h-12 bg-white border border-gray-200 rounded-md mb-2"></div>
                <span class="text-sm font-medium text-gray-900 dark:text-white">Light</span>
              </button>
              
              <button 
                @click="setTheme('dark')"
                class="p-4 border rounded-md flex flex-col items-center"
                :class="[
                  theme === 'dark' 
                    ? 'border-primary-500 bg-primary-50 dark:bg-primary-900 dark:border-primary-400' 
                    : 'border-gray-200 dark:border-gray-700'
                ]"
              >
                <div class="w-12 h-12 bg-gray-900 border border-gray-700 rounded-md mb-2"></div>
                <span class="text-sm font-medium text-gray-900 dark:text-white">Dark</span>
              </button>
              
              <button 
                @click="setTheme('system')"
                class="p-4 border rounded-md flex flex-col items-center"
                :class="[
                  theme === 'system' 
                    ? 'border-primary-500 bg-primary-50 dark:bg-primary-900 dark:border-primary-400' 
                    : 'border-gray-200 dark:border-gray-700'
                ]"
              >
                <div class="w-12 h-12 bg-gradient-to-r from-white to-gray-900 border border-gray-200 rounded-md mb-2"></div>
                <span class="text-sm font-medium text-gray-900 dark:text-white">System</span>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Notifications tab -->
      <div v-if="activeTab === 'notifications'" class="p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Notification Settings</h2>
        
        <div class="space-y-6">
          <!-- Email notifications -->
          <div>
            <h3 class="text-md font-medium text-gray-900 dark:text-white mb-3">Email Notifications</h3>
            
            <div class="space-y-3">
              <div class="flex items-center">
                <input
                  id="emailDowntime"
                  v-model="notificationSettings.emailDowntime"
                  type="checkbox"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label for="emailDowntime" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                  Site downtime alerts
                </label>
              </div>
              
              <div class="flex items-center">
                <input
                  id="emailRecovery"
                  v-model="notificationSettings.emailRecovery"
                  type="checkbox"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label for="emailRecovery" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                  Site recovery notifications
                </label>
              </div>
              
              <div class="flex items-center">
                <input
                  id="emailSummary"
                  v-model="notificationSettings.emailSummary"
                  type="checkbox"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label for="emailSummary" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                  Weekly summary reports
                </label>
              </div>
            </div>
          </div>
          
          <!-- Webhook notifications -->
          <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
            <h3 class="text-md font-medium text-gray-900 dark:text-white mb-3">Webhook Notifications</h3>
            
            <div class="mb-4">
              <label for="webhookUrl" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Webhook URL
              </label>
              <input
                id="webhookUrl"
                v-model="webhookUrl"
                type="url"
                placeholder="https://example.com/webhook"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
              />
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                We'll send a POST request to this URL when your site status changes
              </p>
            </div>
            
            <div class="space-y-3">
              <div class="flex items-center">
                <input
                  id="webhookDowntime"
                  v-model="notificationSettings.webhookDowntime"
                  type="checkbox"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label for="webhookDowntime" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                  Site downtime alerts
                </label>
              </div>
              
              <div class="flex items-center">
                <input
                  id="webhookRecovery"
                  v-model="notificationSettings.webhookRecovery"
                  type="checkbox"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label for="webhookRecovery" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                  Site recovery notifications
                </label>
              </div>
            </div>
          </div>
          
          <!-- Save button -->
          <div class="pt-4">
            <button 
              @click="saveNotificationSettings"
              :disabled="saving"
              class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
            >
              <span v-if="saving">Saving...</span>
              <span v-else>Save Settings</span>
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Delete Account Confirmation Modal -->
    <ConfirmModal
      v-if="showDeleteAccountModal"
      title="Delete Account"
      message="Are you sure you want to delete your account? All of your data will be permanently removed. This action cannot be undone."
      confirm-text="Delete Account"
      cancel-text="Cancel"
      confirm-variant="danger"
      @confirm="deleteAccount"
      @cancel="showDeleteAccountModal = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import ConfirmModal from '../components/common/ConfirmModal.vue';

const router = useRouter();
const authStore = useAuthStore();

// State
const activeTab = ref('account');
const apiKey = ref('api_12345abcdef67890ghijklmnopqrstuvwxyz');
const theme = ref('system');
const webhookUrl = ref('');
const saving = ref(false);
const showDeleteAccountModal = ref(false);

// Tabs
const tabs = [
  { id: 'account', name: 'Account' },
  { id: 'appearance', name: 'Appearance' },
  { id: 'notifications', name: 'Notifications' }
];

// Notification settings
const notificationSettings = ref({
  emailDowntime: true,
  emailRecovery: true,
  emailSummary: false,
  webhookDowntime: false,
  webhookRecovery: false
});

// Methods
const regenerateApiKey = () => {
  // In a real app, we would make an API call to regenerate the API key
  // For demo purposes, we'll just simulate a new API key
  const characters = 'abcdefghijklmnopqrstuvwxyz0123456789';
  let newKey = 'api_';
  
  for (let i = 0; i < 36; i++) {
    newKey += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  
  apiKey.value = newKey;
};

const setTheme = (newTheme) => {
  theme.value = newTheme;
  
  // In a real app, we would save this preference to localStorage or the API
  // For demo purposes, we'll just update the theme class on the document
  if (newTheme === 'dark' || (newTheme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
};

const saveNotificationSettings = () => {
  saving.value = true;
  
  // In a real app, we would save these settings to the API
  // For demo purposes, we'll just simulate a successful save
  setTimeout(() => {
    saving.value = false;
  }, 1000);
};

const deleteAccount = () => {
  // In a real app, we would make an API call to delete the account
  // For demo purposes, we'll just log out and redirect to the home page
  authStore.logout();
  router.push('/');
};

// Lifecycle hooks
onMounted(() => {
  // In a real app, we would fetch the user's settings from the API
  // For demo purposes, we'll just use default values
  
  // Check system theme preference
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    document.documentElement.classList.add('dark');
  }
});
</script>
