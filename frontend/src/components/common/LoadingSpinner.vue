<template>
  <div class="flex justify-center items-center" :class="containerClass">
    <div 
      class="animate-spin rounded-full border-t-2 border-b-2" 
      :class="[
        sizeClasses, 
        colorClass || 'border-primary-600 dark:border-primary-400'
      ]"
    ></div>
    <span v-if="text" class="ml-3 text-gray-700 dark:text-gray-300">{{ text }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl'].includes(value)
  },
  text: {
    type: String,
    default: ''
  },
  color: {
    type: String,
    default: 'primary'
  },
  containerClass: {
    type: String,
    default: ''
  }
});

// Computed properties for dynamic classes
const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'h-4 w-4';
    case 'md':
      return 'h-8 w-8';
    case 'lg':
      return 'h-12 w-12';
    case 'xl':
      return 'h-16 w-16';
    default:
      return 'h-8 w-8';
  }
});

const colorClass = computed(() => {
  switch (props.color) {
    case 'primary':
      return 'border-primary-600 dark:border-primary-400';
    case 'secondary':
      return 'border-gray-600 dark:border-gray-400';
    case 'success':
      return 'border-green-600 dark:border-green-400';
    case 'danger':
      return 'border-red-600 dark:border-red-400';
    case 'warning':
      return 'border-yellow-600 dark:border-yellow-400';
    case 'info':
      return 'border-blue-600 dark:border-blue-400';
    default:
      return 'border-primary-600 dark:border-primary-400';
  }
});
</script>
