import 'devextreme/dist/css/dx.light.css'
import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { useUserStore } from './stores/user'
import Vue3Toastify from 'vue3-toastify'
import i18n from '@/i18n/index.js'

import App from './App.vue'
import router from './router'
import store from '@/store/index.js'

const app = createApp(App)

app.use(Vue3Toastify)
app.use(createPinia())
const userStore = useUserStore()
await userStore.fetchUser()
app.use(store)
app.use(router)
app.use(i18n)
//await router.isReady()

app.mount('#app')