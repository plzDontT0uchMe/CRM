import axios from 'axios'
import router from '@/router/index.js'
import { useUserStore } from '@/stores/user.js'

const customAxios = axios.create({
    baseURL: import.meta.env.MODE == 'development' ? 'http://localhost:3000' : 'https://dev.crm-fitness.ru',
    withCredentials: true,
    validateStatus: () => true
})

customAxios.interceptors.request.use(function(config) {
    //console.log(config)
    return config
}, function(error) {
    console.log(error)
    return Promise.reject(error)
})

customAxios.interceptors.response.use(async function(response) {
    if (response.config.url == '/api/auth' || response.config.url == '/api/reg') {
        return response
    }
    if (!response.data.successfully && response.data.message == 'access token not found in cookie, authorization failed') {
        const { data } = await axios.get(customAxios.defaults.baseURL + '/api/updateToken', { withCredentials: true })
        if (data.successfully) {
            response.config._retryAttempt = true
            return customAxios(response.config)
        } else {
            useUserStore().data = null
            await router.push({ name: 'auth' })
        }
    }
    return response
}, function(error) {
    console.log(error)
    return Promise.reject(error)
})

export default customAxios