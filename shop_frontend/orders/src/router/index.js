import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        requiresAuth: false,
      },
      props: true,
    },
    {
      path: '/about',
      name: 'about',
      meta: {
        requiresAuth: true,
      },
      props: true,
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    { 
      path: '/:pathMatch(.*)*', 
      name: 'not-found', 
      component: NotFound 
    }
  ]
})

router.beforeEach(checkAuth());

export default router
