import { createRouter, createWebHistory } from 'vue-router'
// import Products from '../views/Products.vue'
// import Orders from '../views/Orders.vue'
// import PackingList from '../views/PackingList.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Товары',
      // component: Products,
      meta: {
        requiresAuth: false,
      },
      props: true,
      component: () => import('../views/Products.vue'),
    },
    {
      path: '/orders',
      name: 'Заказы',
      // component: Orders,
      meta: {
        requiresAuth: false,
      },
      props: true,
      component: () => import('../views/Orders.vue'),
    },
    {
      path: '/orders-products',
      name: 'Сборочный лист',
      meta: {
        requiresAuth: true,
      },
      props: true,
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/PackingList.vue')
    }/*,
    { 
      path: '/:pathMatch(.*)*', 
      name: 'not-found', 
      component: NotFound 
    }*/
  ]
})

//router.beforeEach(checkAuth());

export default router
