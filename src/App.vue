<template>
  <div :class="['app', theme]">
    <router-view />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

const store = useStore()
const theme = computed(() => store.state.theme)

const toggleTheme = () => {
  store.commit('toggleTheme')
}
</script>

<style>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
}

:root {
  --primary-color: #16141c;
  --secondary-color: #6e52f9;
  --background-color: #24212b;
  --text-color: #ffffff;
  --label-color: #aaaaaa;
  --input-background-color: #333333;
  --input-text-color: #ffffff;
  --button-background-color: #6e52f9;
  --button-text-color: #ffffff;
  --button-hover-background-color: #5a43d2;
  --header-background-color: #16141c;
  --header-text-color: #ffffff;
  --menu-background-color: #1f1b29;
  --menu-text-color: #ffffff;
  --profile-menu-background: #444444;
  --profile-menu-border: #666666;
  --content-background-color: #282828;
  --content-text-color: #ffffff;
  --modal-background-color: white; /* Цвет фона модального окна */
  --modal-text-color: black;
}

.light {
  --background-color: #ffffff;
  --text-color: #000000;
  --label-color: #666666;
  --input-background-color: #f9f9f9;
  --input-text-color: #000000;
  --button-background-color: #6e52f9;
  --button-text-color: #ffffff;
  --button-hover-background-color: #5a43d2;
  --header-background-color: #1d232a; /* Прежний цвет для светлой темы */
  --header-text-color: #000000;
  --menu-background-color: #e0e0e0;
  --menu-text-color: #000000;
  --profile-menu-background: #ffffff;
  --profile-menu-border: #ccc;
  --content-background-color: #f0f0f0;
  --content-text-color: #000000;
  --modal-background-color: white; /* Цвет фона модального окна */
  --modal-text-color: black;
}

.dark {
  --primary-color: #16141c;
  --background-color: #24212b;
  --text-color: #ffffff;
  --label-color: #aaaaaa;
  --input-background-color: #333333;
  --input-text-color: #ffffff;
  --button-background-color: #6e52f9;
  --button-text-color: #ffffff;
  --button-hover-background-color: #5a43d2;
  --header-background-color: #16141c;
  --header-text-color: #ffffff;
  --menu-background-color: #1f1b29;
  --menu-text-color: #ffffff;
  --profile-menu-background: #444444;
  --profile-menu-border: #666666;
  --content-background-color: #282828;
  --content-text-color: #ffffff;
  --modal-background-color: #444444; /* Цвет фона модального окна */
  --modal-text-color: #fff;
}
</style>
