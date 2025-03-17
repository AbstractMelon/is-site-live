import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../store/auth';

// Import views
import Home from '../views/Home.vue';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Dashboard from '../views/Dashboard.vue';
import SiteDetail from '../views/SiteDetail.vue';
import UserProfile from '../views/UserProfile.vue';
import Settings from '../views/Settings.vue';
import NotFound from '../views/NotFound.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { title: 'Home' }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: 'Login', guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { title: 'Register', guest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { title: 'Dashboard', requiresAuth: true }
  },
  {
    path: '/site/:id',
    name: 'SiteDetail',
    component: SiteDetail,
    meta: { title: 'Site Details' }
  },
  {
    path: '/user/:username',
    name: 'UserProfile',
    component: UserProfile,
    meta: { title: 'User Profile' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: { title: 'Settings', requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: { title: '404 Not Found' }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    // Always scroll to top
    return { top: 0 };
  }
});

// Navigation guards
router.beforeEach((to, from, next) => {
  // Set page title
  document.title = `${to.meta.title || to.name} | Is It Live`;
  
  const authStore = useAuthStore();
  const isLoggedIn = authStore.isAuthenticated;
  
  // Check if the route requires authentication
  if (to.meta.requiresAuth && !isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } });
    return;
  }
  
  // Check if the route is for guests only (like login/register)
  if (to.meta.guest && isLoggedIn) {
    next({ name: 'Dashboard' });
    return;
  }
  
  next();
});

export default router;
