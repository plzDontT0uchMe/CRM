<template>
  <div class="settings-form">
    <div class="settings-header">
      <div class="settings-photo-container">
        <img :src="settings.avatar" alt="Profile Photo" class="settings-photo" />
        <input type="file" @change="onPhotoChange" />
      </div>
    </div>
    <div class="settings-info">
      <div class="form-columns">
        <div class="form-column">
          <div class="form-row">
            <label>{{ $t('firstName') }}:</label>
            <input v-model="settings.firstName" type="text" />
          </div>
          <div class="form-row">
            <label>{{ $t('middleName') }}:</label>
            <input v-model="settings.middleName" type="text" />
          </div>
          <div class="form-row">
            <label>{{ $t('lastName') }}:</label>
            <input v-model="settings.lastName" type="text" />
          </div>
          <div class="form-row">
            <label>{{ $t('dob') }}:</label>
            <input v-model="settings.dateOfBirth" type="date" />
          </div>
        </div>
        <div class="form-column">
          <div class="form-row">
            <label>{{ $t('gender') }}:</label>
            <select v-model="settings.gender">
              <option value="Male">{{ $t('male') }}</option>
              <option value="Female">{{ $t('female') }}</option>
              <option value="Other">{{ $t('other') }}</option>
            </select>
          </div>
          <div class="form-row">
            <label>{{ $t('position') }}:</label>
            <select v-model="settings.position">
              <option value="trainer">{{ $t('trainer') }}</option>
              <option value="user">{{ $t('user') }}</option>
            </select>
          </div>
          <template v-if="settings.position === 'trainer'">
            <div class="form-row">
              <label>{{ $t('experience') }}:</label>
              <input v-model="settings.experience" type="number" />
            </div>
            <div class="form-row">
              <label>{{ $t('sportType') }}:</label>
              <input v-model="settings.sportType" type="text" />
            </div>
            <div class="form-row">
              <label>{{ $t('achievements') }}:</label>
              <textarea v-model="settings.achievements"></textarea>
            </div>
          </template>
        </div>
      </div>
    </div>
    <div class="settings-actions">
      <button @click="saveSettings" class="btn-primary">{{ $t('save') }}</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useStore } from 'vuex'

const props = defineProps({
  settings: Object
})

const { t } = useI18n()
const store = useStore()

const onPhotoChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      props.settings.avatar = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const saveSettings = () => {
  store.commit('updateProfilePhoto', props.settings.avatar)
  alert('Settings saved successfully!')
}
</script>

<style scoped>
.settings-form {
  display: flex;
  flex-direction: column;
  width: 100%;
  background-color: var(--background-color);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  color: var(--text-color);
}

.settings-header {
  display: flex;
  justify-content: left;
  margin-bottom: 20px;
}

.settings-photo-container {
  text-align: left;
  position: relative;
}

.settings-photo {
  width: 150px;
  height: auto; /* Сохранить пропорции изображения */
  background-color: #ddd;
  display: block;
  margin-bottom: 10px;
  border-radius: 8px; /* Прямоугольное изображение с закругленными углами */
}

.settings-info {
  display: flex;
  justify-content: center;
}

.form-columns {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.form-column {
  flex: 1;
  margin-right: 10px;
}

.form-column:last-child {
  margin-right: 0;
}

.form-row {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
}

label {
  margin-bottom: 5px;
  color: var(--label-color);
  font-weight: bold;
}

input,
select,
textarea {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: var(--input-background-color);
  color: var(--input-text-color);
  width: 100%;
  box-sizing: border-box;
  margin: 0;
}

.settings-actions {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.btn-primary {
  background-color: var(--button-background-color);
  color: var(--button-text-color);
  border: none;
  padding: 10px 20px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 1em;
  cursor: pointer;
  border-radius: 25px;
  transition: background-color 0.3s;
}

.btn-primary:hover {
  background-color: var(--button-hover-background-color);
}
</style>
