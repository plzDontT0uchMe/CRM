<script setup>
import axios from '@/axios/index.js'
import { ref } from 'vue'
import router from '@/router/index.js'
import { useUserStore } from '@/stores/user.js'
import { useToastStore } from '@/stores/toast.js'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faEye, faEyeSlash)

const userStore = useUserStore()
const toastStore = useToastStore()

const login = ref('')
const password = ref('')

const reg = async () => {
    const notifyId = toastStore.startToast('loading', 'Выполняется регистрация... 🚀', 'top-center')
    try {
        const { data } = await axios.post('/api/reg', {
            login: login.value,
            password: password.value
        })
        if (data.successfully) {
            await userStore.fetchUser()
            toastStore.stopToast(notifyId, data.message, userStore.data ? 'success' : 'error')
            setTimeout(async () => {
                await router.push({ name: 'main' })
            }, 2000)
        } else {
            toastStore.stopToast(notifyId, data.message, 'error')
        }
    } catch (err) {
        console.log(err)
    }
}


const acceptTerms = ref(false)
const confirmPassword = ref('')
const rememberMe = ref(false)
const passwordFieldType = ref('password')
const togglePasswordVisibility = () => {
    passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password'
}

</script>

<template>
    <div class="register-container">
        <h1>Регистрация</h1>
        <form @submit.prevent="reg">
            <div class="form-group">
                <label for="username">Login:</label>
                <input type="text" id="username" v-model="login" required />
            </div>
            <div class="form-group">
                <label for="password">Password:</label>
                <div class="password-input">
                    <input :type="passwordFieldType" id="password" v-model="password" required />
                    <font-awesome-icon
                        :icon="passwordFieldType === 'password' ? 'eye' : 'eye-slash'"
                        @click="togglePasswordVisibility"
                        class="password-icon"
                    />
                </div>
            </div>
            <div class="form-group">
                <label for="confirm-password">Confirm Password:</label>
                <div class="password-input">
                    <input
                        :type="passwordFieldType"
                        id="confirm-password"
                        v-model="confirmPassword"
                        required
                    />
                    <font-awesome-icon
                        :icon="passwordFieldType === 'password' ? 'eye' : 'eye-slash'"
                        @click="togglePasswordVisibility"
                        class="password-icon"
                    />
                </div>
            </div>
            <div class="form-group remember-me">
                <input type="checkbox" id="terms" v-model="acceptTerms" required />
                <label for="terms">I accept the terms and conditions</label>
            </div>
            <button type="submit">Register</button>
        </form>
        <div class="links">
            <router-link to="/auth">Already have an account? Login here</router-link>
        </div>
<!--        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>-->
    </div>


    <!--    <div class="w-screen h-screen flex flex-col space-y-4 justify-center items-center">
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
            <button class="btn btn-outline" @click="reg">Submit</button>
        </div>-->
</template>

<style scoped>
.register-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100vh;
    padding: 20px;
    box-sizing: border-box;
}

h1 {
    text-align: center;
    margin-bottom: 20px;
}

form {
    width: 100%;
    max-width: 400px;
}

.form-group {
    margin-bottom: 15px;
}

label {
    display: block;
    margin-bottom: 5px;
}

input {
    width: 100%;
    padding: 8px;
    box-sizing: border-box;
    border: 1px solid #ccc;
    border-radius: 4px;
}

.password-input {
    display: flex;
    align-items: center;
    border: 1px solid #ccc;
    border-radius: 4px;
}

.password-input input {
    width: calc(100% - 40px);
    border: none;
    border-radius: 4px 0 0 4px;
    padding-right: 40px;
}

.password-icon {
    padding: 8px;
    cursor: pointer;
}

button[type='submit'] {
    width: 100%;
    padding: 10px;
    background-color: var(--button-background-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

button[type='submit']:hover {
    background-color: var(--button-hover-background-color);
}

.remember-me {
    display: flex;
    align-items: center;
    margin-bottom: 15px;
}

.remember-me input {
    width: 20px;
    height: 20px;
    margin-right: 10px;
}

.links {
    text-align: center;
    margin-top: 20px;
}

.links a {
    display: block;
    margin-bottom: 5px;
    color: #007bff;
    text-decoration: none;
}

.links a:hover {
    text-decoration: underline;
}

.error {
    color: red;
    margin-top: 10px;
    text-align: center;
}
</style>