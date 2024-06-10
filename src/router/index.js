import { createRouter, createWebHistory } from 'vue-router'

const routes = [
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
    component: () => import('@/views/MainView.vue')
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/ProfileView.vue')
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/views/SettingView.vue'),
    children: [
      {
        path: 'account',
        name: 'account',
        component: () => import('@/components/SettingsForm.vue')
      },
      {
        path: 'security',
        name: 'security',
        component: () => import('@/components/SecuritySettings.vue')
      },
      {
        path: 'info',
        name: 'info',
        component: () => import('@/components/InfoSettings.vue')
      },
      {
        path: 'billing',
        name: 'billing',
        component: () => import('@/components/BillingSettings.vue')
      }
    ]
  },
  {
    path: '/membership',
    name: 'membership',
    component: () => import('@/views/MemberShipView.vue')
  },
  {
    path: '/exercise',
    name: 'exercise',
    component: () => import('@/views/ExerciseView.vue')
  },
  {
    path: '/workout-plans',
    name: 'workout-plans',
    component: () => import('@/views/WorkoutPlansView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
