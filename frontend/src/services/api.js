import axios from 'axios';

// Create axios instance with default config
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to add the auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor to handle common errors
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Handle 401 Unauthorized errors (token expired or invalid)
    if (error.response && error.response.status === 401) {
      // Clear token and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
const authAPI = {
  register: (userData) => api.post('/auth/register', userData),
  login: (credentials) => api.post('/auth/login', credentials),
  getCurrentUser: () => api.get('/user'),
  updateUser: (userData) => api.put('/user', userData),
};

// Sites API
const sitesAPI = {
  createSite: (siteData) => api.post('/sites', siteData),
  getUserSites: () => api.get('/sites'),
  getSite: (id) => api.get(`/sites/${id}`),
  updateSite: (id, siteData) => api.put(`/sites/${id}`, siteData),
  deleteSite: (id) => api.delete(`/sites/${id}`),
  getSiteStats: (id) => api.get(`/site/${id}/stats`),
};

// Domains API
const domainsAPI = {
  createCustomDomain: (domainData) => api.post('/domains', domainData),
  getUserDomains: () => api.get('/domains'),
  deleteCustomDomain: (id) => api.delete(`/domains/${id}`),
  verifyCustomDomain: (id) => api.get(`/domains/${id}/verify`),
  getDomainDashboard: (domain) => api.get(`/domain/${domain}`),
};

// Public API
const publicAPI = {
  getUserProfile: (username) => api.get(`/user/${username}`),
};

export default {
  auth: authAPI,
  sites: sitesAPI,
  domains: domainsAPI,
  public: publicAPI,
};
