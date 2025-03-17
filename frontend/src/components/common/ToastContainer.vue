<template>
  <div class="fixed top-4 right-4 z-50 space-y-4 w-80">
    <transition-group name="toast">
      <div 
        v-for="toast in toasts" 
        :key="toast.id" 
        class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden"
      >
        <div class="flex p-4">
          <!-- Icon -->
          <div class="flex-shrink-0">
            <!-- Success icon -->
            <svg v-if="toast.type === 'success'" class="h-6 w-6 text-green-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            
            <!-- Error icon -->
            <svg v-else-if="toast.type === 'error'" class="h-6 w-6 text-red-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            
            <!-- Warning icon -->
            <svg v-else-if="toast.type === 'warning'" class="h-6 w-6 text-yellow-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            
            <!-- Info icon -->
            <svg v-else class="h-6 w-6 text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          
          <!-- Content -->
          <div class="ml-3 flex-1">
            <p class="text-sm font-medium text-gray-900 dark:text-white">
              {{ toast.title }}
            </p>
            <p v-if="toast.message" class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ toast.message }}
            </p>
          </div>
          
          <!-- Close button -->
          <div class="ml-4 flex-shrink-0 flex">
            <button 
              @click="removeToast(toast.id)" 
              class="inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:text-gray-500"
            >
              <span class="sr-only">Close</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
        
        <!-- Progress bar -->
        <div 
          v-if="toast.duration > 0" 
          class="h-1 bg-gray-200 dark:bg-gray-700"
        >
          <div 
            class="h-full transition-all duration-linear"
            :class="getProgressBarColor(toast.type)"
            :style="{ width: `${getProgressWidth(toast)}%` }"
          ></div>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue';

// State
const toasts = ref([]);
const toastTimeouts = ref({});

// Methods
const addToast = (toast) => {
  const id = Date.now().toString();
  const newToast = {
    id,
    title: toast.title || 'Notification',
    message: toast.message || '',
    type: toast.type || 'info',
    duration: toast.duration !== undefined ? toast.duration : 5000, // Default 5 seconds, 0 for persistent
    createdAt: Date.now()
  };
  
  toasts.value.push(newToast);
  
  // Auto-remove toast after duration (if not persistent)
  if (newToast.duration > 0) {
    toastTimeouts.value[id] = setTimeout(() => {
      removeToast(id);
    }, newToast.duration);
  }
  
  return id;
};

const removeToast = (id) => {
  const index = toasts.value.findIndex(toast => toast.id === id);
  if (index !== -1) {
    toasts.value.splice(index, 1);
    
    // Clear timeout if exists
    if (toastTimeouts.value[id]) {
      clearTimeout(toastTimeouts.value[id]);
      delete toastTimeouts.value[id];
    }
  }
};

const getProgressWidth = (toast) => {
  if (toast.duration <= 0) return 100;
  
  const elapsed = Date.now() - toast.createdAt;
  const remaining = Math.max(0, toast.duration - elapsed);
  return (remaining / toast.duration) * 100;
};

const getProgressBarColor = (type) => {
  switch (type) {
    case 'success':
      return 'bg-green-500';
    case 'error':
      return 'bg-red-500';
    case 'warning':
      return 'bg-yellow-500';
    case 'info':
    default:
      return 'bg-blue-500';
  }
};

// Clean up timeouts on component unmount
onUnmounted(() => {
  Object.values(toastTimeouts.value).forEach(timeout => {
    clearTimeout(timeout);
  });
});

// Expose methods to parent components
defineExpose({
  addToast,
  removeToast
});
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.toast-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>
