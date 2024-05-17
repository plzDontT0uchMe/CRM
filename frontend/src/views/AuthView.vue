<script setup>
import axios from '@/axios/index.js'
import { onMounted, ref } from 'vue'
import { useCookies } from "vue3-cookies";
import router from '@/router/index.js'

const { cookies } = useCookies();

const login = ref('')
const password = ref('')

const auth = async () => {
    try {
        const { data } = await axios.post('/api/auth', {
            login: login.value,
            password: password.value
        })
        console.log('OK')
        console.log(data)
        if (data.successfully) {
            cookies.set('access_token', data.session.access_token)
            cookies.set('refresh_token', data.session.refresh_token)
            router.push({ name: 'main' })
        }
    } catch (error) {
        console.log('123')
    }
}

const getCourse = async () => {
    try {
        const resp = await axios.post('/api/getList')
        console.log(resp.data)
    }
    catch (err) {
        console.log(err)
    }
}

const checkAuth = async () => {
    try {
        const resp = await axios.post('/api/checkAuth', {
            access_token: cookies.get('access_token'),
        })
        console.log(resp.data)
        if (resp.data.successfully) {
            if (resp.data.flag == 'getRefreshToken') {
                try {
                    const resp1 = await axios.post('/api/checkAuth', {
                        refresh_token: cookies.get('refresh_token'),
                    })
                    console.log(resp1.data)
                    console.log('OK')
                    if(resp1.data.flag == 'newAccessToken') {
                        cookies.set('access_token', resp1.data.message)
                        router.push({ name: 'main' })
                    }
                    if (resp1.data.flag == 'authorizationFailed') {
                        console.log('Ошибка авторизации!')
                    }
                }
                catch (error) {
                    console.log('321')
                }
            }
            if(resp.data.flag != 'getRefreshToken') {
                router.push({ name: 'main' })
            }
        }
    } catch (error) {
        console.log('123')
    }
}

const checkAuthTest = async () => {
    const resp = await axios.get('/hello')
    console.log(resp.data)
}

onMounted(() => {
    checkAuthTest()
    //checkAuth()
})

/*setInterval(() => {
    checkAuth()
}, 5000)*/

</script>

<template>
    <div class="w-screen h-screen flex flex-col space-y-4 justify-center items-center">
        <label class="form-control w-full max-w-xs">
            <div class="label">
                <span class="label-text pl-2">Enter the login</span>
            </div>
            <label class="input input-bordered flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
                     class="w-4 h-4 opacity-70">
                    <path
                        d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
                </svg>
                <input type="text" class="grow" placeholder="login" v-model="login" />
            </label>
        </label>
        <label class="form-control w-full max-w-xs">
            <div class="label">
                <span class="label-text pl-2">Enter the password</span>
            </div>
            <label class="input input-bordered flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
                     class="w-4 h-4 opacity-70">
                    <path fill-rule="evenodd"
                          d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                          clip-rule="evenodd" />
                </svg>
                <input type="text" class="grow" placeholder="password" v-model="password" />
            </label>
        </label>
        <button class="btn btn-outline" @click="auth">Submit</button>
    </div>
</template>

<style scoped>
</style>