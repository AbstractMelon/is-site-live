<template>
  <div 
    v-if="show"
    class="rounded-md p-4 mb-4"
    :class="typeClasses"
    role="alert"
  >
    <div class="flex">
      <div class="flex-shrink-0">
        <!-- Success icon -->
        <svg v-if="type === 'success'" class="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
        
        <!-- Error icon -->
        <svg v-else-if="type === 'error'" class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        
        <!-- Warning icon -->
        <svg v-else-if="type === 'warning'" class="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        
        <!-- Info icon -->
        <svg v-else class="h-5 w-5 text-blue-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
        </svg>
      </div>
      
      <div class="ml-3">
        <h3 v-if="title" class="text-sm font-medium" :class="textColorClass">
          {{ title }}
        </h3>
        <div class="text-sm" :class="[title ? 'mt-2' : '', textColorClass]">
          <p>{{ message }}</p>
        </div>
      </div>
      
      <div class="ml-auto pl-3">
        <div class="-mx-1.5 -my-1.5">
          <button 
            v-if="dismissible"
            @click="dismiss"
            class="inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2"
            :class="dismissButtonClass"
          >
            <span class="sr-only">Dismiss</span>
            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue';

const props = defineProps({
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['info', 'success', 'warning', 'error'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    required: true
  },
  dismissible: {
    type: Boolean,
    default: true
  },
  autoClose: {
    type: Boolean,
    default: false
  },
  duration: {
    type: Number,
    default: 5000 // 5 seconds
  }
});

const emit = defineEmits(['dismiss']);

// State
const show = ref(true);
let autoCloseTimeout = null;

// Computed properties
const typeClasses = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-50 dark:bg-green-900 border border-green-100 dark:border-green-800';
    case 'error':
      return 'bg-red-50 dark:bg-red-900 border border-red-100 dark:border-red-800';
    case 'warning':
      return 'bg-yellow-50 dark:bg-yellow-900 border border-yellow-100 dark:border-yellow-800';
    case 'info':
    default:
      return 'bg-blue-50 dark:bg-blue-900 border border-blue-100 dark:border-blue-800';
  }
});

const textColorClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-800 dark:text-green-200';
    case 'error':
      return 'text-red-800 dark:text-red-200';
    case 'warning':
      return 'text-yellow-800 dark:text-yellow-200';
    case 'info':
    default:
      return 'text-blue-800 dark:text-blue-200';
  }
});

const dismissButtonClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-50 dark:bg-green-900 text-green-500 hover:bg-green-100 dark:hover:bg-green-800 focus:ring-green-500 focus:ring-offset-green-50 dark:focus:ring-offset-green-900';
    case 'error':
      return 'bg-red-50 dark:bg-red-900 text-red-500 hover:bg-red-100 dark:hover:bg-red-800 focus:ring-red-500 focus:ring-offset-red-50 dark:focus:ring-offset-red-900';
    case 'warning':
      return 'bg-yellow-50 dark:bg-yellow-900 text-yellow-500 hover:bg-yellow-100 dark:hover:bg-yellow-800 focus:ring-yellow-500 focus:ring-offset-yellow-50 dark:focus:ring-offset-yellow-900';
    case 'info':
    default:
      return 'bg-blue-50 dark:bg-blue-900 text-blue-500 hover:bg-blue-100 dark:hover:bg-blue-800 focus:ring-blue-500 focus:ring-offset-blue-50 dark:focus:ring-offset-blue-900';
  }
});

// Methods
const dismiss = () => {
  show.value = false;
  emit('dismiss');
};

// Lifecycle hooks
onMounted(() => {
  if (props.autoClose) {
    autoCloseTimeout = setTimeout(() => {
      dismiss();
    }, props.duration);
  }
});

onUnmounted(() => {
  if (autoCloseTimeout) {
    clearTimeout(autoCloseTimeout);
  }
});
</script>
