import 'devextreme/dist/css/dx.light.css'
import './assets/main.css'
import store from './store'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import i18n from './i18n'

import App from './App.vue'

import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(store)
app.use(i18n)

app.mount('#app')
