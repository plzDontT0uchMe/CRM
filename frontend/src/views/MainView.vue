<script setup>
import axios from '@/axios/index.js'
import { onMounted, ref } from 'vue'

const menuItems = ref([
    {
        name: 'subscriptions',
        route: 'subscriptions'
    }
])
</script>

<template>
    <div class="flex justify-between items-center">
        <div class="w-[15%] h-96 bg-white">
            <ul class="menu flex items-center w-full h-full bg-navbar rounded-xl">
                <li class="w-full p-1" v-for="(item, index) in menuItems" :key="index">
                    <router-link :to="{ name: item.route }" class="w-full h-full flex items-center justify-center">
                        {{ item.name }}
                    </router-link>
                    <ul v-if="item.children" class="w-full h-full bg-white rounded-xl">
                        <li class="w-full p-1" v-for="(child, index) in item.children" :key="index">
                            <router-link :to="{ name: child.route }" class="w-full h-full flex items-center justify-center">
                                {{ child.name }}
                            </router-link>
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
        <div class="w-[83%] h-96 bg-white">
            <RouterView />
        </div>
    </div>
</template>

<style scoped>

</style>