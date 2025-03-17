<template>
  <footer class="bg-white dark:bg-gray-800 shadow-md mt-12">
    <div class="container mx-auto px-4 py-6">
      <div class="flex flex-col md:flex-row justify-between items-center">
        <div class="mb-4 md:mb-0">
          <div class="flex items-center">
            <span class="text-primary-500 text-xl font-bold mr-2">⚡</span>
            <span class="text-gray-900 dark:text-white text-lg font-semibold">Is It Live</span>
          </div>
          <p class="text-gray-600 dark:text-gray-400 text-sm mt-2">
            A self-hosted uptime monitoring platform
          </p>
        </div>
        
        <div class="flex flex-col md:flex-row md:space-x-8">
          <div class="mb-4 md:mb-0">
            <h3 class="text-gray-900 dark:text-white font-medium mb-2">Links</h3>
            <ul class="text-gray-600 dark:text-gray-400 text-sm space-y-2">
              <li>
                <router-link to="/" class="hover:text-primary-500 dark:hover:text-primary-400">
                  Home
                </router-link>
              </li>
              <li>
                <router-link to="/dashboard" class="hover:text-primary-500 dark:hover:text-primary-400">
                  Dashboard
                </router-link>
              </li>
              <li>
                <a href="https://github.com/abstractmelon/is-site-live" target="_blank" class="hover:text-primary-500 dark:hover:text-primary-400">
                  GitHub
                </a>
              </li>
            </ul>
          </div>
          
          <div>
            <h3 class="text-gray-900 dark:text-white font-medium mb-2">Resources</h3>
            <ul class="text-gray-600 dark:text-gray-400 text-sm space-y-2">
              <li>
                <a href="https://github.com/abstractmelon/is-site-live/blob/main/README.md" target="_blank" class="hover:text-primary-500 dark:hover:text-primary-400">
                  Documentation
                </a>
              </li>
              <li>
                <a href="https://github.com/abstractmelon/is-site-live/issues" target="_blank" class="hover:text-primary-500 dark:hover:text-primary-400">
                  Report an Issue
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>
      
      <div class="border-t border-gray-200 dark:border-gray-700 mt-6 pt-6 text-center text-gray-600 dark:text-gray-400 text-sm">
        <p>&copy; {{ new Date().getFullYear() }} Is It Live. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
// No additional script needed
</script>
