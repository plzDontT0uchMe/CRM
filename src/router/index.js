import { createRouter, createWebHistory } from 'vue-router'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // {
    //   path: '/',
    //   name: 'main',
    //   component: () => import('@/views/MainView.vue')
    // },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/views/AuthView.vue')
    },
    {
      path: '/reg',
      name: 'reg',
      component: () => import('@/views/RegView.vue')
    },
    {
      path: '/',
      name: 'main',
      component:()=>import('@/views/MainView.vue')

    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminPanel.vue')
    },
    {
      path: '/profile',
      name: 'profile',
      component:()=>import('@/views/ProfileView.vue')

    },


    
  ]
})

export default router
