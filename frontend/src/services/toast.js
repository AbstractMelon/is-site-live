// Toast service for managing toast notifications
import { ref } from 'vue';

// Create a reactive reference to store the toast container component
const toastContainer = ref(null);

// Register the toast container component
const registerContainer = (container) => {
  toastContainer.value = container;
};

// Show a toast notification
const showToast = (options) => {
  if (!toastContainer.value) {
    console.warn('Toast container not registered. Make sure to include ToastContainer component in your app.');
    return null;
  }
  
  return toastContainer.value.addToast(options);
};

// Helper methods for common toast types
const success = (title, message = '', duration = 5000) => {
  return showToast({ title, message, type: 'success', duration });
};

const error = (title, message = '', duration = 5000) => {
  return showToast({ title, message, type: 'error', duration });
};

const warning = (title, message = '', duration = 5000) => {
  return showToast({ title, message, type: 'warning', duration });
};

const info = (title, message = '', duration = 5000) => {
  return showToast({ title, message, type: 'info', duration });
};

// Remove a specific toast by ID
const removeToast = (id) => {
  if (!toastContainer.value) {
    console.warn('Toast container not registered.');
    return;
  }
  
  toastContainer.value.removeToast(id);
};

export default {
  registerContainer,
  showToast,
  success,
  error,
  warning,
  info,
  removeToast
};
