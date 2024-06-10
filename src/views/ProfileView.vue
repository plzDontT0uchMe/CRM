<template>
  <div class="profile-container">
    <Header />
    <div class="content">
      <Menu />
      <div class="profile-content">
        <ProfileForm />
      </div>
    </div>
  </div>
</template>

<script setup>
import Header from '../components/Header.vue'
import Menu from '../components/Menu.vue'
import ProfileForm from '../components/ProfileForm.vue'
</script>

<style scoped>
.profile-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
}

.content {
  display: flex;
  flex: 1;
  background-color: var(--content-background-color);
  color: var(--content-text-color);
}

.menu {
  height: 100%;
  flex-shrink: 0;
}

.profile-content {
  flex-grow: 1;
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
