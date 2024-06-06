import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { useUserStore } from './stores/user'
import Vue3Toastify from 'vue3-toastify'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(Vue3Toastify)
app.use(createPinia())
const userStore = useUserStore()
await userStore.fetchUser()
app.use(router)
//await router.isReady()

app.mount('#app')