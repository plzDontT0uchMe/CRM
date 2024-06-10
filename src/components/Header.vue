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
      <div
        class="profile-icon"
        @click="toggleProfileMenu"
        :class="{ active: isProfileMenuVisible }"
      >
        <img :src="profilePhoto" alt="Profile" />
        <div
          v-if="isProfileMenuVisible"
          class="profile-menu"
          :class="[theme, { 'menu-visible': isProfileMenuVisible }]"
        >
          <p @click="viewProfile" class="menu-item">{{ $t('profile') }}</p>
          <p @click="viewSettings" class="menu-item">{{ $t('settings') }}</p>
          <p @click="logout" class="menu-item">{{ $t('logout') }}</p>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
import SunIcons from './icons/SunIcons.vue'
import MoonIcons from './icons/MoonIcons.vue'

const store = useStore()
const { t, locale } = useI18n()
const profilePhoto = computed(() => store.getters.profilePhoto)
const theme = computed(() => store.state.theme)

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

const viewProfile = () => {
  router.push('/profile')
}

const viewSettings = () => {
  router.push('/settings')
}

const logout = () => {
  router.push('/login')
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
