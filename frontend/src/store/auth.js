import { defineStore } from 'pinia';
import api from '../services/api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
    token: localStorage.getItem('token') || null,
    loading: false,
    error: null
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token,
    getUser: (state) => state.user,
    getToken: (state) => state.token,
    isLoading: (state) => state.loading,
    getError: (state) => state.error
  },
  
  actions: {
    async register(username, password, email) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.auth.register({ username, password, email });
        
        const { user, token } = response.data;
        
        this.setUser(user);
        this.setToken(token);
        
        return user;
      } catch (error) {
        this.error = error.response?.data?.error || 'Registration failed';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async login(username, password) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.auth.login({ username, password });
        
        const { user, token } = response.data;
        
        this.setUser(user);
        this.setToken(token);
        
        return user;
      } catch (error) {
        this.error = error.response?.data?.error || 'Login failed';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async checkAuth() {
      if (!this.token) {
        return null;
      }
      
      this.loading = true;
      
      try {
        const response = await api.auth.getCurrentUser();
        
        this.setUser(response.data);
        
        return this.user;
      } catch (error) {
        // If authentication fails, clear the token and user
        this.logout();
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async updateUser(userData) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.auth.updateUser(userData);
        
        this.setUser(response.data);
        
        return this.user;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update user';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    setUser(user) {
      this.user = user;
      localStorage.setItem('user', JSON.stringify(user));
    },
    
    setToken(token) {
      this.token = token;
      if (token) {
        localStorage.setItem('token', token);
      } else {
        localStorage.removeItem('token');
      }
    },
    
    logout() {
      this.user = null;
      this.token = null;
      localStorage.removeItem('user');
      localStorage.removeItem('token');
    }
  }
});
