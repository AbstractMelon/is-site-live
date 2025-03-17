<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-md mx-4">
      <div class="p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ title }}</h2>
          <button @click="$emit('cancel')" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="mb-6 text-gray-700 dark:text-gray-300">
          {{ message }}
        </div>
        
        <div class="flex justify-end space-x-3">
          <button
            type="button"
            @click="$emit('cancel')"
            class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            {{ cancelText }}
          </button>
          <button
            type="button"
            @click="$emit('confirm')"
            :class="[
              'px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2',
              confirmVariant === 'danger' 
                ? 'bg-red-600 hover:bg-red-700 focus:ring-red-500' 
                : 'bg-primary-600 hover:bg-primary-700 focus:ring-primary-500'
            ]"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  cancelText: {
    type: String,
    default: 'Cancel'
  },
  confirmVariant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'danger'].includes(value)
  }
});

defineEmits(['confirm', 'cancel']);
</script>
