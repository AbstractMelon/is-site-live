<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Your Monitored Sites</h1>
      <button 
        @click="showAddSiteModal = true" 
        class="bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-4 rounded-md shadow-sm transition duration-300"
      >
        Add New Site
      </button>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center my-12">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-600"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-md mb-6">
      {{ error }}
      <button @click="fetchSites" class="underline ml-2">Try again</button>
    </div>

    <!-- Empty state -->
    <div v-else-if="sites.length === 0" class="bg-gray-50 dark:bg-gray-800 rounded-lg p-8 text-center my-8">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
      <h3 class="text-xl font-semibold mb-2 text-gray-900 dark:text-white">No sites monitored yet</h3>
      <p class="text-gray-600 dark:text-gray-400 mb-4">
        Add your first site to start monitoring its uptime.
      </p>
      <button 
        @click="showAddSiteModal = true" 
        class="bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-4 rounded-md shadow-sm transition duration-300"
      >
        Add Your First Site
      </button>
    </div>

    <!-- Sites grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <SiteCard 
        v-for="site in sites" 
        :key="site.id" 
        :site="site" 
        @edit="editSite = site; showEditSiteModal = true"
        @delete="confirmDeleteSite = site; showDeleteConfirmModal = true"
      />
    </div>

    <!-- Add Site Modal -->
    <AddSiteModal 
      v-if="showAddSiteModal" 
      @close="showAddSiteModal = false"
      @site-added="onSiteAdded"
    />

    <!-- Edit Site Modal -->
    <EditSiteModal 
      v-if="showEditSiteModal && editSite" 
      :site="editSite"
      @close="showEditSiteModal = false; editSite = null"
      @site-updated="onSiteUpdated"
    />

    <!-- Delete Confirmation Modal -->
    <ConfirmModal
      v-if="showDeleteConfirmModal && confirmDeleteSite"
      title="Delete Site"
      :message="`Are you sure you want to delete '${confirmDeleteSite.name}'? This action cannot be undone.`"
      confirm-text="Delete"
      cancel-text="Cancel"
      confirm-variant="danger"
      @confirm="deleteSite(confirmDeleteSite.id)"
      @cancel="showDeleteConfirmModal = false; confirmDeleteSite = null"
    />

    <!-- Public Dashboard URL Section -->
    <div class="mt-12 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
      <h2 class="text-xl font-semibold mb-4 text-gray-900 dark:text-white">Your Public Dashboard</h2>
      <p class="text-gray-600 dark:text-gray-400 mb-4">
        Share this URL to let others see the status of your monitored sites:
      </p>
      <div class="flex items-center">
        <input 
          type="text" 
          readonly 
          :value="publicDashboardUrl" 
          class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:text-white"
        />
        <button 
          @click="copyToClipboard(publicDashboardUrl)" 
          class="ml-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 font-medium py-2 px-4 rounded-md transition duration-300"
        >
          Copy
        </button>
      </div>
    </div>

    <!-- Custom Domains Section -->
    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Custom Domains</h2>
        <button 
          @click="showAddDomainModal = true" 
          class="bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-3 rounded-md shadow-sm text-sm transition duration-300"
        >
          Add Domain
        </button>
      </div>

      <div v-if="loadingDomains" class="flex justify-center my-6">
        <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary-600"></div>
      </div>

      <div v-else-if="domains.length === 0" class="text-gray-600 dark:text-gray-400 mb-4">
        You haven't added any custom domains yet. Add a domain to create a branded dashboard.
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="domain in domains" 
          :key="domain.id" 
          class="border border-gray-200 dark:border-gray-700 rounded-md p-4"
        >
          <div class="flex justify-between items-center">
            <div>
              <div class="flex items-center">
                <h3 class="font-medium text-gray-900 dark:text-white">{{ domain.domain }}</h3>
                <span 
                  v-if="domain.verified" 
                  class="ml-2 px-2 py-0.5 bg-green-100 text-green-800 text-xs rounded-full"
                >
                  Verified
                </span>
                <span 
                  v-else 
                  class="ml-2 px-2 py-0.5 bg-yellow-100 text-yellow-800 text-xs rounded-full"
                >
                  Unverified
                </span>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Added on {{ new Date(domain.created_at).toLocaleDateString() }}
              </p>
            </div>
            <div class="flex space-x-2">
              <button 
                v-if="!domain.verified"
                @click="verifyDomain(domain.id)" 
                class="text-primary-600 hover:text-primary-700 font-medium text-sm"
              >
                Verify
              </button>
              <button 
                @click="confirmDeleteDomain = domain; showDeleteDomainModal = true" 
                class="text-red-600 hover:text-red-700 font-medium text-sm"
              >
                Delete
              </button>
            </div>
          </div>
          <div v-if="!domain.verified" class="mt-3 text-sm text-gray-600 dark:text-gray-400">
            <p>To verify this domain, add a CNAME record pointing to <code class="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{{ window.location.hostname }}</code></p>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Domain Modal -->
    <AddDomainModal 
      v-if="showAddDomainModal" 
      @close="showAddDomainModal = false"
      @domain-added="onDomainAdded"
    />

    <!-- Delete Domain Confirmation Modal -->
    <ConfirmModal
      v-if="showDeleteDomainModal && confirmDeleteDomain"
      title="Delete Custom Domain"
      :message="`Are you sure you want to delete '${confirmDeleteDomain.domain}'? This action cannot be undone.`"
      confirm-text="Delete"
      cancel-text="Cancel"
      confirm-variant="danger"
      @confirm="deleteDomain(confirmDeleteDomain.id)"
      @cancel="showDeleteDomainModal = false; confirmDeleteDomain = null"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { useSitesStore } from '../store/sites';
import { useDomainsStore } from '../store/domains';
import SiteCard from '../components/sites/SiteCard.vue';
import AddSiteModal from '../components/sites/AddSiteModal.vue';
import EditSiteModal from '../components/sites/EditSiteModal.vue';
import AddDomainModal from '../components/domains/AddDomainModal.vue';
import ConfirmModal from '../components/common/ConfirmModal.vue';

const router = useRouter();
const authStore = useAuthStore();
const sitesStore = useSitesStore();
const domainsStore = useDomainsStore();

// Sites state
const loading = computed(() => sitesStore.isLoading);
const error = computed(() => sitesStore.getError);
const sites = computed(() => sitesStore.getSites);

// Domains state
const loadingDomains = computed(() => domainsStore.isLoading);
const domains = computed(() => domainsStore.getDomains);

// Modal state
const showAddSiteModal = ref(false);
const showEditSiteModal = ref(false);
const showDeleteConfirmModal = ref(false);
const showAddDomainModal = ref(false);
const showDeleteDomainModal = ref(false);
const editSite = ref(null);
const confirmDeleteSite = ref(null);
const confirmDeleteDomain = ref(null);

// Computed properties
const publicDashboardUrl = computed(() => {
  const user = authStore.getUser;
  if (!user) return '';
  return `${window.location.origin}/user/${user.username}`;
});

// Methods
const fetchSites = async () => {
  try {
    await sitesStore.fetchSites();
  } catch (error) {
    console.error('Failed to fetch sites:', error);
  }
};

const fetchDomains = async () => {
  try {
    await domainsStore.fetchDomains();
  } catch (error) {
    console.error('Failed to fetch domains:', error);
  }
};

const onSiteAdded = (site) => {
  showAddSiteModal.value = false;
  fetchSites();
};

const onSiteUpdated = (site) => {
  showEditSiteModal.value = false;
  editSite.value = null;
  fetchSites();
};

const deleteSite = async (id) => {
  try {
    await sitesStore.deleteSite(id);
    showDeleteConfirmModal.value = false;
    confirmDeleteSite.value = null;
  } catch (error) {
    console.error('Failed to delete site:', error);
  }
};

const onDomainAdded = (domain) => {
  showAddDomainModal.value = false;
  fetchDomains();
};

const verifyDomain = async (id) => {
  try {
    await domainsStore.verifyDomain(id);
  } catch (error) {
    console.error('Failed to verify domain:', error);
  }
};

const deleteDomain = async (id) => {
  try {
    await domainsStore.deleteDomain(id);
    showDeleteDomainModal.value = false;
    confirmDeleteDomain.value = null;
  } catch (error) {
    console.error('Failed to delete domain:', error);
  }
};

const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text)
    .then(() => {
      // Could add a toast notification here
      console.log('URL copied to clipboard');
    })
    .catch(err => {
      console.error('Failed to copy URL:', err);
    });
};

// Lifecycle hooks
onMounted(() => {
  fetchSites();
  fetchDomains();
});
</script>
