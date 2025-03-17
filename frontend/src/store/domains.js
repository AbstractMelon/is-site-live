import { defineStore } from 'pinia';
import api from '../services/api';

export const useDomainsStore = defineStore('domains', {
  state: () => ({
    domains: [],
    loading: false,
    error: null
  }),
  
  getters: {
    getDomains: (state) => state.domains,
    isLoading: (state) => state.loading,
    getError: (state) => state.error
  },
  
  actions: {
    async fetchDomains() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.domains.getUserDomains();
        this.domains = response.data;
        return this.domains;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch domains';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async createDomain(domainData) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.domains.createCustomDomain(domainData);
        this.domains.push(response.data);
        return response.data;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create domain';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async verifyDomain(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.domains.verifyCustomDomain(id);
        
        // Update the domain in the domains array
        const index = this.domains.findIndex(domain => domain.id === id);
        if (index !== -1) {
          this.domains[index] = {
            ...this.domains[index],
            verified: response.data.verified,
            verificationStatus: response.data.verificationStatus
          };
        }
        
        return response.data;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to verify domain';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteDomain(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await api.domains.deleteCustomDomain(id);
        
        // Remove the domain from the domains array
        this.domains = this.domains.filter(domain => domain.id !== id);
        
        return true;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete domain';
        throw error;
      } finally {
        this.loading = false;
      }
    }
  }
});
