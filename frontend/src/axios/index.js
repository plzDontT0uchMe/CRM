import axios from 'axios'
import router from '@/router/index.js'
//import { useCookies } from "vue3-cookies";

//const { cookies } = useCookies();

const customAxios = axios.create({
    baseURL: import.meta.env.MODE == 'development' ? 'http://127.0.0.1:3000' : 'http://45.156.23.236:3000',
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
    //console.log(response)
    if (response.config.url == '/api/auth' || response.config.url == '/api/reg') {
        return response
    }
    if (!response.data.successfully) {
        const {data} = await axios.get(customAxios.defaults.baseURL + '/api/updateToken', { withCredentials: true })
        if (data.successfully) {
            response.config._retryAttempt = true;
            return customAxios(response.config);
        } else {
            if (response.config.url != '/api/checkAuth') {
                router.push('/suck')
            }
        }
    }
    return response
}, function(error) {
    console.log(error)
    return Promise.reject(error)
})

export default customAxios