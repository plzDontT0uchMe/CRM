<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import ShowIcons from '@/components/icons/ShowIcons.vue'
import HideIcons from '@/components/icons/HideIcons.vue'

const Password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const router = useRouter()

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}
const toggleConfirmPasswordVisibility = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}
const navigateToAuth = () => {
  router.push('/auth')
}
</script>

<template>
  <div class="main">
    <div class="container-auth">
      <div class="container">
        <div class="logo">
          <p>Регистрация</p>
        </div>

        <div class="field-input">
          <label class="input input-bordered flex items-center gap-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="w-4 h-4 opacity-70"
            >
              <path
                d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z"
              />
            </svg>
            <input type="text" class="grow" placeholder="Login" />
          </label>
          <label class="input input-bordered flex items-center gap-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="w-4 h-4 opacity-70"
            >
              <path
                fill-rule="evenodd"
                d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                clip-rule="evenodd"
              />
            </svg>
            <input
              class="grow"
              v-model="Password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Password"
            />
          </label>
          <button class="showHide" @click="togglePasswordVisibility">
            <HideIcons v-if="showPassword" />
            <ShowIcons v-else />
          </button>

          <label class="input input-bordered flex items-center gap-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="w-4 h-4 opacity-70"
            >
              <path
                fill-rule="evenodd"
                d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                clip-rule="evenodd"
              />
            </svg>
            <input
              class="grow"
              v-model="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              placeholder="Confirm password"
            />
          </label>
          <button class="showHideConfirm" @click="toggleConfirmPasswordVisibility">
            <HideIcons class="icon" v-if="showConfirmPassword" />
            <ShowIcons class="icon" v-else />
          </button>
        </div>

        <div class="link-reg">
          <p>Have account?</p>
          <a @click="navigateToAuth"> Click here</a>
        </div>
        <div class="button">
          <button class="button-auth btn btn-outline btn-secondary">Зарегистрироваться</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
* {
  font-family: Inter;
}
.main {
  background-color: #16141c;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100vw;
  height: 100vh;
  position: relative;
}

.container-auth {
  width: 50%;
  padding: 5px;
  border: 2px solid #fff;
  border-radius: 10%;
  display: block;
  overflow: hidden;
}
.container {
  margin: 5px 5px;
  width: 100%;
  position: relative;
}
.logo {
  align-items: center;
  overflow: hidden;
}
p {
  color: white;
  text-align: center;
  font-size: 1.5em;
  margin-bottom: 2%;
}
.field-input {
  width: 96%;
  height: 100%;
  position: relative;
  align-items: center;
}
.input {
  margin-top: 10px;
  width: 100%;
}

.showHide {
  outline: none;
  border: none;
  background: transparent;
  position: absolute;
  top: 50%;
  right: 10px;
  transform: translateY(-50%);
}
.showHideConfirm {
  outline: none;
  border: none;
  background: transparent;
  position: absolute;
  top: 85%;
  right: 10px;
  transform: translateY(-50%);
}
.forgot-pass {
  width: 100%;
  margin: 8px auto;
  display: flex;
}
.forgot-pass a {
  text-align: right;
  font-size: 0.8em;
}

.for-got {
  position: absolute;
  top: 75%;
  right: 10px;
}

.container input {
  outline: none;
  border: none;
  font-size: 1em;
  color: white;
}
input::placeholder {
  justify-content: center;
  color: white;
  font-size: 1em;
  vertical-align: middle;
}
input:focus {
  color: white;
}
a {
  display: block;
  text-align: left;
  font-size: 1em;
  text-decoration: none;
  color: #fff;
}
a:hover {
  color: #6e52f9;
}
.input-checkbox {
  margin-right: 7px;
  transform: scale(1.1);
  justify-content: center;
}
.input-checkbox + label {
  display: inline-block;
  font-size: 0.8em;
  color: #fff;
  cursor: pointer;
}
.link-reg {
  margin-top: 2%;
  display: flex;
}
.link-reg p {
  font-size: 1em;
}
.button {
  padding: 0 auto;
  width: 100%;
  margin: 0 auto;
  margin-top: 20px;
  position: relative;
}
.button-auth {
  color: #ffffff;
  width: 70%;
  position: absolute;
  top: 100%;
  left: 13%;
  padding: 0.5em 2em;
  margin: 0 auto;
  font-size: 1em;
  transform: translate(-50% -50%);
}
.button-auth:hover {
  background-color: #6e52f9;
  border-color: #6e52f9;
  color: #fff;
}

@media (max-width: 1920px) {
  .container-auth {
    height: 43vh;
    width: 20%;
  }
}

@media (max-width: 1500px) {
  .container-auth {
    height: 43vh;
    width: 25%;
  }
}
@media (max-width: 1200px) {
  .container-auth {
    height: 43vh;
    width: 30%;
  }
}
@media (max-width: 900px) {
  .container-auth {
    height: 43vh;
    width: 35%;
  }
}
@media (max-width: 660px) {
  .container-auth {
    height: 43vh;
    width: 45%;
  }
}
@media (max-width: 530px) {
  .container-auth {
    height: 43vh;
    width: 60%;
  }
}
@media (max-width: 480px) {
  .container-auth {
    height: 43vh;
    width: 70%;
  }
}
</style>
