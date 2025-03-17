<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <Navbar />
    <main class="container mx-auto px-4 py-8">
      <router-view />
    </main>
    <Footer />
    <ToastContainer ref="toastContainerRef" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from './store/auth';
import Navbar from './components/layout/Navbar.vue';
import Footer from './components/layout/Footer.vue';
import ToastContainer from './components/common/ToastContainer.vue';
import toastService from './services/toast';

const router = useRouter();
const authStore = useAuthStore();
const toastContainerRef = ref(null);

onMounted(async () => {
  // Check if the user is already logged in (token in localStorage)
  await authStore.checkAuth();
  
  // If dark mode is enabled in localStorage, apply it
  if (localStorage.getItem('darkMode') === 'true') {
    document.documentElement.classList.add('dark');
  }
  
  // Register toast container with the toast service
  toastService.registerContainer(toastContainerRef.value);
});
</script>
