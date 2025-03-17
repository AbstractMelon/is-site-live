import { defineStore } from 'pinia';
import api from '../services/api';

export const useSitesStore = defineStore('sites', {
  state: () => ({
    sites: [],
    currentSite: null,
    loading: false,
    error: null
  }),
  
  getters: {
    getSites: (state) => state.sites,
    getCurrentSite: (state) => state.currentSite,
    isLoading: (state) => state.loading,
    getError: (state) => state.error
  },
  
  actions: {
    async fetchSites() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.sites.getUserSites();
        this.sites = response.data;
        return this.sites;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch sites';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchSite(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.sites.getSite(id);
        this.currentSite = response.data;
        return this.currentSite;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch site';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchSiteStats(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.sites.getSiteStats(id);
        return response.data;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch site stats';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async createSite(siteData) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.sites.createSite(siteData);
        this.sites.push(response.data);
        return response.data;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create site';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async updateSite(id, siteData) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.sites.updateSite(id, siteData);
        
        // Update the site in the sites array
        const index = this.sites.findIndex(site => site.id === id);
        if (index !== -1) {
          this.sites[index] = response.data;
        }
        
        // Update currentSite if it's the one being updated
        if (this.currentSite && this.currentSite.id === id) {
          this.currentSite = response.data;
        }
        
        return response.data;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update site';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteSite(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await api.sites.deleteSite(id);
        
        // Remove the site from the sites array
        this.sites = this.sites.filter(site => site.id !== id);
        
        // Clear currentSite if it's the one being deleted
        if (this.currentSite && this.currentSite.id === id) {
          this.currentSite = null;
        }
        
        return true;
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete site';
        throw error;
      } finally {
        this.loading = false;
      }
    }
  }
});
