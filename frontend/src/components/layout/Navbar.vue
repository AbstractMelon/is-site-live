<template>
  <nav class="bg-white dark:bg-gray-800 shadow-md">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and site name -->
        <div class="flex items-center">
          <router-link to="/" class="flex items-center">
            <span class="text-primary-500 text-2xl font-bold mr-2">⚡</span>
            <span class="text-gray-900 dark:text-white text-xl font-semibold">Is It Live</span>
          </router-link>
        </div>
        
        <!-- Navigation links -->
        <div class="hidden md:flex space-x-6">
          <router-link to="/" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
            Home
          </router-link>
          
          <template v-if="isAuthenticated">
            <router-link to="/dashboard" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Dashboard
            </router-link>
            <router-link to="/settings" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Settings
            </router-link>
          </template>
          
          <template v-else>
            <router-link to="/login" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Login
            </router-link>
            <router-link to="/register" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Register
            </router-link>
          </template>
        </div>
        
        <!-- User menu and theme toggle -->
        <div class="flex items-center space-x-4">
          <!-- Theme toggle -->
          <button @click="toggleDarkMode" class="p-2 rounded-full text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700">
            <svg v-if="isDarkMode" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          </button>
          
          <!-- User dropdown (if authenticated) -->
          <div v-if="isAuthenticated" class="relative">
            <button @click="toggleUserMenu" class="flex items-center text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              <span class="mr-1">{{ user?.username }}</span>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            
            <!-- Dropdown menu -->
            <div v-if="userMenuOpen" class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 z-10">
              <router-link :to="`/user/${user?.username}`" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                Profile
              </router-link>
              <router-link to="/settings" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                Settings
              </router-link>
              <button @click="logout" class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                Logout
              </button>
            </div>
          </div>
        </div>
        
        <!-- Mobile menu button -->
        <div class="md:hidden">
          <button @click="toggleMobileMenu" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
            <svg v-if="!mobileMenuOpen" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      
      <!-- Mobile menu -->
      <div v-if="mobileMenuOpen" class="md:hidden py-4">
        <div class="flex flex-col space-y-4">
          <router-link to="/" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
            Home
          </router-link>
          
          <template v-if="isAuthenticated">
            <router-link to="/dashboard" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Dashboard
            </router-link>
            <router-link :to="`/user/${user?.username}`" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Profile
            </router-link>
            <router-link to="/settings" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Settings
            </router-link>
            <button @click="logout" class="text-left text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Logout
            </button>
          </template>
          
          <template v-else>
            <router-link to="/login" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Login
            </router-link>
            <router-link to="/register" class="text-gray-700 dark:text-gray-300 hover:text-primary-500 dark:hover:text-primary-400">
              Register
            </router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../store/auth';

const router = useRouter();
const authStore = useAuthStore();

const mobileMenuOpen = ref(false);
const userMenuOpen = ref(false);
const isDarkMode = ref(document.documentElement.classList.contains('dark'));

const isAuthenticated = computed(() => authStore.isAuthenticated);
const user = computed(() => authStore.getUser);

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
  
  // Close user menu if open
  if (userMenuOpen.value) {
    userMenuOpen.value = false;
  }
};

const toggleUserMenu = () => {
  userMenuOpen.value = !userMenuOpen.value;
};

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value;
  
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark');
    localStorage.setItem('darkMode', 'true');
  } else {
    document.documentElement.classList.remove('dark');
    localStorage.setItem('darkMode', 'false');
  }
};

const logout = () => {
  authStore.logout();
  router.push('/login');
  
  // Close menus
  userMenuOpen.value = false;
  mobileMenuOpen.value = false;
};

// Close user menu when clicking outside
window.addEventListener('click', (event) => {
  if (userMenuOpen.value && !event.target.closest('.user-menu')) {
    userMenuOpen.value = false;
  }
});
</script>
