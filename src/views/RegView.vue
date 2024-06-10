<template>
  <div class="register-container">
    <h1>Регистрация</h1>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="username">Login:</label>
        <input type="text" id="username" v-model="username" required />
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
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faEye, faEyeSlash)

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const acceptTerms = ref(false)
const errorMessage = ref('')
const passwordFieldType = ref('password')

const handleSubmit = async () => {
  if (username.value && password.value && confirmPassword.value) {
    if (password.value !== confirmPassword.value) {
      errorMessage.value = 'Passwords do not match'
      return
    }
    if (!acceptTerms.value) {
      errorMessage.value = 'You must accept the terms and conditions'
      return
    }
    try {
      const response = await fakeRegister(username.value, password.value)
      if (response.success) {
        alert('Registration successful!')
        // Дополнительные действия при успешной регистрации
      } else {
        errorMessage.value = 'Registration failed'
      }
    } catch (error) {
      errorMessage.value = 'An error occurred. Please try again.'
    }
  } else {
    errorMessage.value = 'Please fill in all fields.'
  }
}

const togglePasswordVisibility = () => {
  passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password'
}

const fakeRegister = (username, password) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      if (username && password) {
        resolve({ success: true })
      } else {
        resolve({ success: false })
      }
    }, 1000)
  })
}

const login = () => {
  alert('Login link clicked')
}
</script>

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
