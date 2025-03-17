<template>
  <div>
    <!-- Hero Section -->
    <section class="py-12 md:py-20 bg-gradient-to-br from-primary-500 to-primary-700 text-white">
      <div class="container mx-auto px-4">
        <div class="flex flex-col md:flex-row items-center">
          <div class="md:w-1/2 mb-8 md:mb-0">
            <h1 class="text-4xl md:text-5xl font-bold mb-4">
              Know when your sites are down before your users do
            </h1>
            <p class="text-xl mb-8 text-white/90">
              Is It Live is a self-hosted uptime monitoring platform that helps you keep track of your websites and services.
            </p>
            <div class="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
              <router-link to="/register" class="btn-primary">
                Get Started
              </router-link>
              <a href="https://github.com/abstractmelon/is-site-live" target="_blank" class="btn-secondary">
                View on GitHub
              </a>
            </div>
          </div>
          <div class="md:w-1/2 flex justify-center">
            <img src="/images/dashboard-preview.png" alt="Dashboard Preview" class="rounded-lg shadow-xl max-w-full h-auto" />
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16 bg-white dark:bg-gray-900">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl font-bold text-center mb-12 text-gray-900 dark:text-white">
          Everything you need for uptime monitoring
        </h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="bg-gray-50 dark:bg-gray-800 p-6 rounded-lg shadow-md">
            <div class="text-primary-500 text-3xl mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">Real-time Monitoring</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Monitor your websites and services with 60-second checks. Get notified when something goes down.
            </p>
          </div>
          
          <!-- Feature 2 -->
          <div class="bg-gray-50 dark:bg-gray-800 p-6 rounded-lg shadow-md">
            <div class="text-primary-500 text-3xl mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">Detailed Statistics</h3>
            <p class="text-gray-600 dark:text-gray-400">
              View detailed uptime statistics, response times, and historical data with interactive charts.
            </p>
          </div>
          
          <!-- Feature 3 -->
          <div class="bg-gray-50 dark:bg-gray-800 p-6 rounded-lg shadow-md">
            <div class="text-primary-500 text-3xl mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">Self-Hosted</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Host it yourself with Docker. You own your data and have complete control over your monitoring setup.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="py-16 bg-gray-50 dark:bg-gray-800">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl font-bold text-center mb-12 text-gray-900 dark:text-white">
          How It Works
        </h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Step 1 -->
          <div class="flex flex-col items-center text-center">
            <div class="bg-primary-500 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">
              1
            </div>
            <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">Add Your Sites</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Register and add the websites or services you want to monitor.
            </p>
          </div>
          
          <!-- Step 2 -->
          <div class="flex flex-col items-center text-center">
            <div class="bg-primary-500 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">
              2
            </div>
            <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">Automatic Monitoring</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Our system checks your sites every 60 seconds and records their status.
            </p>
          </div>
          
          <!-- Step 3 -->
          <div class="flex flex-col items-center text-center">
            <div class="bg-primary-500 text-white rounded-full w-12 h-12 flex items-center justify-center text-xl font-bold mb-4">
              3
            </div>
            <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">Get Notified</h3>
            <p class="text-gray-600 dark:text-gray-400">
              Receive email alerts when your sites go down and when they recover.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-16 bg-primary-500 text-white">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl font-bold mb-4">Ready to start monitoring your sites?</h2>
        <p class="text-xl mb-8 max-w-2xl mx-auto">
          Join today and never be caught off guard by downtime again.
        </p>
        <div class="flex flex-col sm:flex-row justify-center space-y-4 sm:space-y-0 sm:space-x-4">
          <router-link to="/register" class="btn-white">
            Sign Up Now
          </router-link>
          <router-link to="/login" class="btn-outline-white">
            Log In
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
// No additional script needed
</script>

<style scoped>
.btn-primary {
  @apply bg-white text-primary-700 font-medium py-2 px-6 rounded-lg shadow-md hover:bg-gray-100 transition duration-300 text-center;
}

.btn-secondary {
  @apply bg-transparent border-2 border-white text-white font-medium py-2 px-6 rounded-lg hover:bg-white/10 transition duration-300 text-center;
}

.btn-white {
  @apply bg-white text-primary-700 font-medium py-2 px-6 rounded-lg shadow-md hover:bg-gray-100 transition duration-300 text-center;
}

.btn-outline-white {
  @apply bg-transparent border-2 border-white text-white font-medium py-2 px-6 rounded-lg hover:bg-white/10 transition duration-300 text-center;
}
</style>
