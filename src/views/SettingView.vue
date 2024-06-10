<template>
  <div :class="['app', theme]">
    <Header />
    <div class="content">
      <Menu />
      <div class="main-content">
        <div class="tabs">
          <button :class="{ active: selectedTab === 'account' }" @click="selectTab('account')">
            {{ $t('account') }}
          </button>
          <button :class="{ active: selectedTab === 'security' }" @click="selectTab('security')">
            {{ $t('security') }}
          </button>
          <button :class="{ active: selectedTab === 'info' }" @click="selectTab('info')">
            {{ $t('info') }}
          </button>
          <button
            :class="{ active: selectedTab === 'notification' }"
            @click="selectTab('notifications')"
          >
            {{ $t('notifications') }}
          </button>
        </div>
        <div v-if="selectedTab === 'account'" class="tab-content">
          <SettingsForm :settings="trainerSettings" @close="selectedTrainerProfile = null" />
        </div>
        <div v-if="selectedTab === 'security'" class="tab-content">
          <h2>{{ $t('securitySettings') }}</h2>
          <!-- Добавьте контент для настроек безопасности -->
        </div>
        <div v-if="selectedTab === 'info'" class="tab-content">
          <ProfileForm :trainer="selectedTrainerProfile" @close="selectedTrainerProfile = null" />
        </div>
        <div v-if="selectedTab === 'notification'" class="tab-content">
          <h2>{{ $t('notificationSettings') }}</h2>
          <!-- Добавьте контент для настроек уведомлений -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
import Header from '../components/Header.vue'
import Menu from '../components/Menu.vue'
import ProfileForm from '../components/ProfileForm.vue'
import SettingsForm from '../components/SettingsForm.vue'

const store = useStore()
const theme = computed(() => store.state.theme)
const { t } = useI18n()

const selectedTab = ref('account')

const selectTab = (tab) => {
  selectedTab.value = tab
}

const selectedTrainerProfile = ref({
  firstName: 'John',
  lastName: 'Doe',
  middleName: 'A.',
  dateOfBirth: '1980-01-01',
  gender: 'Male',
  sportType: 'Weightlifting',
  achievements: 'World Champion 2010, 2012',
  avatar: 'https://via.placeholder.com/150',
  rating: 4,
  experience: 10,
  position: 'trainer',
  lastActive: new Date('2023-06-01T12:00:00Z')
})

const trainerSettings = ref({
  ...selectedTrainerProfile.value
})
</script>

<style scoped>
.app {
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

.main-content {
  flex-grow: 1;
  padding: 20px;
}

.tabs {
  display: flex;
  justify-content: flex-start; /* Выровнять вкладки по левому краю */
  margin-bottom: 20px;
}

.tabs button {
  background-color: transparent;
  color: inherit;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s;
  border-bottom: 2px solid transparent;
}

.tabs button.active {
  background-color: transparent;
  border-bottom: 2px solid #007bff;
  color: #007bff;
}

.tab-content {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
}

.tab-content > * {
  width: 100%;
}
</style>
