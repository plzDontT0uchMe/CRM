<script setup>

import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast.js'
import router from '@/router/index.js'
import axios from '@/axios/index.js'
import WomenIcon from '@/components/icons/WomenIcon.vue'
import UserIcon from '@/components/icons/UserIcon.vue'
import ManIcon from '@/components/icons/ManIcon.vue'


const userStore = useUserStore()
const toastStore = useToastStore()

const route = useRoute()

const logout = async () => {
    const notifyId = toastStore.startToast('loading', 'Выполняется выход... 🚀', 'top-center')
    try {
        const { data } = await axios.post('/api/logout')
        if (data.successfully) {
            toastStore.stopToast(notifyId, "Выход выполнен успешно", "success")
            setTimeout(async () => {
                userStore.data = null
                await router.push({ name: 'auth' })
            }, 2000)
        }
        toastStore.stopToast(notifyId, data.message, "error")
    }
    catch (err) {
        console.log(err)
    }
}

</script>

<template>
    <div class="navbar bg-base-100"
         v-if="!(route.name === 'auth' || route.name === 'reg' || route.name === undefined)">
        <div class="flex-1">
            <router-link :to="{ name: 'main'}" class="justify-between">
                <span class="btn btn-ghost text-xl">CRM</span>
            </router-link>
        </div>
        <div class="flex-none">
            <div class="dropdown dropdown-end">
                <div tabindex="0" role="button" class="btn btn-ghost btn-circle">
                    <div class="indicator">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                             stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                        <span class="badge badge-sm indicator-item">8</span>
                    </div>
                </div>
                <div tabindex="0" class="mt-3 z-[1] card card-compact dropdown-content w-52 bg-base-100 shadow">
                    <div class="card-body">
                        <span class="font-bold text-lg">8 Items</span>
                        <span class="text-info">Subtotal: $999</span>
                        <div class="card-actions">
                            <button class="btn btn-primary btn-block">View cart</button>
                        </div>
                    </div>
                </div>
            </div>
            <div class="dropdown dropdown-end">
                <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
                    <div class="w-10 rounded-full">
                        <img v-if="userStore?.data?.image" alt="Tailwind CSS Navbar component"
                             :src="axios.defaults.baseURL + '/api/getImage/' + userStore?.data?.image" />
                        <div v-else>
                            <ManIcon v-if="userStore?.data?.gender === 1"/>
                            <WomenIcon v-else-if="userStore?.data?.gender === 2"/>
                            <UserIcon v-else />
                        </div>
                    </div>
                </div>
                <ul tabindex="0"
                    class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
                    <li>
                        <router-link :to="{ name: 'profile', params: { id: userStore?.data?.id || null } }"
                                     class="justify-between">
                            Profile
                        </router-link>
                    </li>
                    <li>
                        <router-link :to="{ name: 'settings' }"
                                     class="justify-between">
                            Settings
                        </router-link>
                    </li>
                    <li>
                        <a @click="logout">
                            Logout
                        </a>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<style scoped>

</style>