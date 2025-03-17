<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-md mb-6">
      {{ error }}
      <button @click="fetchUserProfile" class="underline ml-2">Try again</button>
    </div>

    <!-- User profile -->
    <div v-else-if="user" class="space-y-8">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div class="flex flex-col md:flex-row md:items-center">
          <div class="flex-shrink-0 mb-4 md:mb-0 md:mr-6">
            <div class="w-24 h-24 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center text-gray-600 dark:text-gray-400 text-2xl font-semibold">
              {{ userInitials }}
            </div>
          </div>
          
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-1">{{ user.username }}</h1>
            <p v-if="user.email" class="text-gray-600 dark:text-gray-400 mb-2">{{ user.email }}</p>
            <p class="text-sm text-gray-500 dark:text-gray-500">
              Member since {{ formatDate(user.created_at) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Account settings -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-6">Account Settings</h2>
        
        <form @submit.prevent="updateProfile">
          <!-- Username -->
          <div class="mb-4">
            <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Username
            </label>
            <input
              id="username"
              v-model="formData.username"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          
          <!-- Email -->
          <div class="mb-4">
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Email Address
            </label>
            <input
              id="email"
              v-model="formData.email"
              type="email"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Used for notifications and password recovery
            </p>
          </div>
          
          <!-- Save button -->
          <div class="mt-6">
            <button
              type="submit"
              :disabled="updateLoading"
              class="w-full md:w-auto px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
            >
              <span v-if="updateLoading">Saving...</span>
              <span v-else>Save Changes</span>
            </button>
          </div>
        </form>
      </div>

      <!-- Change password -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-6">Change Password</h2>
        
        <form @submit.prevent="changePassword">
          <!-- Current password -->
          <div class="mb-4">
            <label for="currentPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Current Password
            </label>
            <input
              id="currentPassword"
              v-model="passwordData.currentPassword"
              type="password"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          
          <!-- New password -->
          <div class="mb-4">
            <label for="newPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              New Password
            </label>
            <input
              id="newPassword"
              v-model="passwordData.newPassword"
              type="password"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          
          <!-- Confirm new password -->
          <div class="mb-4">
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Confirm New Password
            </label>
            <input
              id="confirmPassword"
              v-model="passwordData.confirmPassword"
              type="password"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          
          <!-- Error message -->
          <div v-if="passwordError" class="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
            {{ passwordError }}
          </div>
          
          <!-- Save button -->
          <div class="mt-6">
            <button
              type="submit"
              :disabled="passwordLoading"
              class="w-full md:w-auto px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
            >
              <span v-if="passwordLoading">Updating...</span>
              <span v-else>Change Password</span>
            </button>
          </div>
        </form>
      </div>

      <!-- Notification settings -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-6">Notification Settings</h2>
        
        <div class="space-y-4">
          <div class="flex items-center">
            <input
              id="emailDowntime"
              v-model="notificationSettings.emailDowntime"
              type="checkbox"
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
            />
            <label for="emailDowntime" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
              Email me when a site goes down
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
              Email me when a site recovers
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
              Send me weekly summary reports
            </label>
          </div>
        </div>
        
        <!-- Save button -->
        <div class="mt-6">
          <button
            @click="saveNotificationSettings"
            :disabled="notificationLoading"
            class="w-full md:w-auto px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
          >
            <span v-if="notificationLoading">Saving...</span>
            <span v-else>Save Preferences</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../store/auth';

const authStore = useAuthStore();

// State
const loading = ref(true);
const error = ref(null);
const user = ref(null);
const updateLoading = ref(false);
const passwordLoading = ref(false);
const notificationLoading = ref(false);
const passwordError = ref(null);

// Form data
const formData = ref({
  username: '',
  email: ''
});

const passwordData = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const notificationSettings = ref({
  emailDowntime: true,
  emailRecovery: true,
  emailSummary: false
});

// Computed properties
const userInitials = computed(() => {
  if (!user.value) return '';
  
  const username = user.value.username;
  if (!username) return '';
  
  return username.substring(0, 2).toUpperCase();
});

// Methods
const fetchUserProfile = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    // In a real app, we would fetch the user profile from the API
    // For demo purposes, we'll use the current user from the auth store
    const currentUser = authStore.getUser;
    
    if (!currentUser) {
      error.value = 'User not found';
      loading.value = false;
      return;
    }
    
    // Simulate API call
    setTimeout(() => {
      user.value = {
        ...currentUser,
        created_at: new Date().toISOString() // Mock creation date
      };
      
      // Initialize form data
      formData.value = {
        username: user.value.username,
        email: user.value.email || ''
      };
      
      loading.value = false;
    }, 1000);
  } catch (err) {
    error.value = 'Failed to load user profile. Please try again.';
    loading.value = false;
  }
};

const updateProfile = async () => {
  updateLoading.value = true;
  
  try {
    // In a real app, we would send the updated profile to the API
    // For demo purposes, we'll just simulate a successful update
    setTimeout(() => {
      user.value = {
        ...user.value,
        username: formData.value.username,
        email: formData.value.email
      };
      
      // Update the user in the auth store
      authStore.updateUser({
        username: formData.value.username,
        email: formData.value.email
      });
      
      updateLoading.value = false;
    }, 1000);
  } catch (err) {
    error.value = 'Failed to update profile. Please try again.';
    updateLoading.value = false;
  }
};

const changePassword = async () => {
  passwordLoading.value = true;
  passwordError.value = null;
  
  // Validate passwords
  if (passwordData.value.newPassword !== passwordData.value.confirmPassword) {
    passwordError.value = 'New passwords do not match';
    passwordLoading.value = false;
    return;
  }
  
  if (passwordData.value.newPassword.length < 8) {
    passwordError.value = 'Password must be at least 8 characters long';
    passwordLoading.value = false;
    return;
  }
  
  try {
    // In a real app, we would send the password change request to the API
    // For demo purposes, we'll just simulate a successful update
    setTimeout(() => {
      // Reset form
      passwordData.value = {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      };
      
      passwordLoading.value = false;
    }, 1000);
  } catch (err) {
    passwordError.value = 'Failed to change password. Please check your current password and try again.';
    passwordLoading.value = false;
  }
};

const saveNotificationSettings = async () => {
  notificationLoading.value = true;
  
  try {
    // In a real app, we would send the notification settings to the API
    // For demo purposes, we'll just simulate a successful update
    setTimeout(() => {
      notificationLoading.value = false;
    }, 1000);
  } catch (err) {
    error.value = 'Failed to save notification preferences. Please try again.';
    notificationLoading.value = false;
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

// Lifecycle hooks
onMounted(() => {
  fetchUserProfile();
});
</script>
