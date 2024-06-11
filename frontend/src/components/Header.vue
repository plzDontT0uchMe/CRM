<template>
  <header class="header">
    <div class="header-right">
      <!-- Компонент смены темы -->
      <div v-if="theme === 'light'">
        <MoonIcons />
      </div>
      <div v-else>
        <SunIcons />
      </div>

      <!-- Кнопка перевода текста -->
      <button @click="toggleLanguage" class="icon-button">
        <img
          v-if="currentLanguage === 'en'"
          src="./icons/eng.png"
          alt="English"
          class="flag-icon"
        />
        <img v-else src="./icons/rus.png" alt="Russian" class="flag-icon" />
      </button>

      <!-- Иконка профиля -->
        <div class="dropdown dropdown-end">
            <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
                <div class="w-10 rounded-full">
                    <img v-if="userStore?.data?.image" alt="Tailwind CSS Navbar component"
                         :src="axios.defaults.baseURL + '/api/getImage/' + userStore?.data?.image" />
                    <div v-else>
                        <ManIcon v-if="userStore?.data?.gender === 1"/>
                        <WomenIcon v-else-if="userStore?.data?.gender === 2"/>
                        <UserIcon v-else />
                    </div>
                </div>
            </div>
            <ul tabindex="0"
                class="menu menu-sm profile-menu dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
                <li>
                    <router-link :to="{ name: 'profile', params: { id: userStore?.data?.id || null } }"
                                 class="justify-between">
                        Profile
                    </router-link>
                </li>
                <li>
                    <router-link :to="{ name: 'settings' }"
                                 class="justify-between">
                        Settings
                    </router-link>
                </li>
                <li>
                    <a @click="logout">
                        Logout
                    </a>
                </li>
            </ul>
        </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
import SunIcons from './icons/SunIcons.vue'
import MoonIcons from './icons/MoonIcons.vue'
import { useUserStore } from '@/stores/user.js'
import { useToastStore } from '@/stores/toast.js'
import axios from '@/axios/index.js'
import ManIcon from '@/components/icons/ManIcon.vue'
import UserIcon from '@/components/icons/UserIcon.vue'
import WomenIcon from '@/components/icons/WomenIcon.vue'

const store = useStore()
const { t, locale } = useI18n()
const theme = computed(() => store.state.theme)

const userStore = useUserStore()
const toastStore = useToastStore()

const route = useRoute()

const logout = async () => {
    const notifyId = toastStore.startToast('loading', 'Выполняется выход... 🚀', 'top-center')
    try {
        const { data } = await axios.post('/api/logout')
        if (data.successfully) {
            toastStore.stopToast(notifyId, "Выход выполнен успешно", "success")
            setTimeout(async () => {
                userStore.data = null
                await router.push({ name: 'auth' })
            }, 2000)
        }
        toastStore.stopToast(notifyId, data.message, "error")
        if (userStore?.data) {
            await router.push({ name: 'auth' })
        }
    }
    catch (err) {
        console.log(err)
    }
}

const isProfileMenuVisible = ref(false)
const currentLanguage = ref('en')
const router = useRouter()

const toggleLanguage = () => {
  currentLanguage.value = currentLanguage.value === 'en' ? 'ru' : 'en'
  locale.value = currentLanguage.value
}

const toggleProfileMenu = () => {
  isProfileMenuVisible.value = !isProfileMenuVisible.value
}

</script>

<style scoped>
.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 10px;
  background-color: var(--header-background-color);
  color: var(--header-text-color);
  position: relative; /* Ensure header is positioned relative to the profile menu */
  z-index: 10000; /* Ensure the header itself is above other content */
}

.header-right {
  display: flex;
  align-items: center;
}

.icon-button {
  background: none;
  border: none;
  cursor: pointer;
  margin-left: 10px;
}

.flag-icon {
  width: 24px;
  height: 24px;
}

.profile-icon {
  position: relative;
  cursor: pointer;
  margin-left: 10px;
}

.profile-icon.active img {
  border: 2px solid var(--secondary-color);
}

.profile-icon img {
  border-radius: 50%;
  width: 40px;
  height: 40px;
}

.profile-menu {
  position: absolute;
  top: 50px;
  right: 0;
  background-color: var(--profile-menu-background);
  border: 1px solid var(--profile-menu-border);
  padding: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  z-index: 10; /* Base z-index */
}

.profile-menu.menu-visible {
  z-index: 10000; /* Higher z-index when visible */
}

.profile-menu.light {
  --profile-menu-background: #ffffff;
  --profile-menu-border: #ccc;
}

.profile-menu.dark {
  --profile-menu-background: #333;
  --profile-menu-border: #555;
}

.profile-menu .menu-item {
  padding: 10px;
  cursor: pointer;
}

.profile-menu .menu-item:hover {
  background-color: var(--secondary-color);
  color: white;
}
</style>
