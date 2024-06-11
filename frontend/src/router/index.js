import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user.js'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'main',
            component: () => import('../views/MainView.vue'),
            children: [
                {
                    path: 'subs',
                    name: 'subscriptions',
                    component: () => import('../views/SubscriptionsView.vue')
                }
            ]
        },
        {
            path: '/auth',
            name: 'auth',
            component: () => import('../views/AuthView.vue')
        },
        {
            path: '/reg',
            name: 'reg',
            component: () => import('../views/RegView.vue')
        },
        {
            path: '/profile/:id?',
            name: 'profile',
            component: () => import('../views/ProfileView.vue')
        },
        {
            path: '/settings',
            name: 'settings',
            component: () => import('../views/SettingsView.vue'),
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
        },
        {
            path: "/:pathMatch(.*)*",
            component: () => import('../views/NotFoundView.vue')
        }
    ]
})

router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    if((to.name == "auth" || to.name == "reg") && userStore?.data?.id) {
        return next({
            name: 'main'
        })
    }
    if(!userStore?.data && (to.name != "auth" && to.name != "reg")) {
        return next({
            name: 'auth'
        })
    }
    next()
})

export default router
