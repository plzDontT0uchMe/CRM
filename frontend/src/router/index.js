import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user.js'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'main',
            component: () => import('../views/MainView.vue')
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
            component: () => import('../views/SettingsView.vue')
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
//todo перенести проверку checkAuth сюда!

export default router
