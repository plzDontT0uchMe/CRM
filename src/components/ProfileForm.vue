<template>
  <div class="profile-form">
    <div class="profile-header">
      <div class="profile-photo-container">
        <img :src="photoUrl" alt="Profile Photo" class="profile-photo" />
        <div class="rating">
          <span
            v-for="star in 5"
            :key="star"
            class="mask mask-star-2"
            :class="{ 'bg-warning': star <= rating }"
          ></span>
        </div>
        <p class="last-active">{{ $t('lastActive') }}: {{ formattedLastActive }}</p>
      </div>
      <div class="profile-info"></div>
    </div>
    <div class="profile-details">
      <div class="form-columns">
        <div class="form-column">
          <div class="form-row">
            <label>{{ $t('firstName') }}:</label>
            <p>{{ firstName }}</p>
          </div>
          <div class="form-row">
            <label>{{ $t('middleName') }}:</label>
            <p>{{ middleName }}</p>
          </div>
          <div class="form-row">
            <label>{{ $t('gender') }}:</label>
            <p>{{ gender }}</p>
          </div>
          <div class="form-row">
            <label>{{ $t('accountCreated') }}:</label>
            <p>{{ createdAt }}</p>
          </div>
        </div>
        <div class="form-column">
          <div class="form-row">
            <label>{{ $t('lastName') }}:</label>
            <p>{{ lastName }}</p>
          </div>
          <div class="form-row">
            <label>{{ $t('dob') }}:</label>
            <p>{{ dob }}</p>
          </div>
          <div class="form-row">
            <label>{{ $t('role') }}:</label>
            <p>{{ position }}</p>
          </div>
          <template v-if="position === 'trainer'">
            <div class="form-row">
              <label>{{ $t('experience') }}:</label>
              <p>{{ experience }} {{ $t('years') }}</p>
            </div>
            <div class="form-row">
              <label>{{ $t('sportType') }}:</label>
              <p>{{ sportType }}</p>
            </div>
            <div class="form-row">
              <label>{{ $t('achievements') }}:</label>
              <p>{{ achievements }}</p>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import { formatDistanceToNow } from 'date-fns'

const store = useStore()

const firstName = ref('John')
const lastName = ref('Doe')
const middleName = ref('Smith')
const dob = ref('1990-01-01')
const gender = ref('Male')
const lastActive = ref(new Date('2023-06-01T12:00:00Z'))
const createdAt = ref('2020-01-01')
const position = ref('trainer') // Проверка позиции пользователя
const experience = ref(10) // Добавлено поле стаж работы
const sportType = ref('Weightlifting') // Добавлено поле вид спорта
const achievements = ref('World Champion 2010, 2012') // Добавлено поле достижения
const rating = ref(5) // Добавлен рейтинг пользователя
const photoUrl = computed(() => store.getters.profilePhoto)

const formattedLastActive = computed(() => {
  return formatDistanceToNow(new Date(lastActive.value), { addSuffix: true })
})
</script>

<style scoped>
.profile-form {
  display: flex;
  flex-direction: column;
  width: 100%;
  background-color: var(--background-color);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  color: var(--text-color);
}

.profile-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.profile-photo-container {
  text-align: center;
  margin-right: 20px;
  flex-shrink: 0;
  position: relative;
}

.profile-photo {
  width: 150px;
  height: 150px;
  background-color: #ddd;
  display: block;
  margin-bottom: 10px;
  border-radius: 50%;
}

.last-active {
  margin-top: 10px;
  font-size: 0.875em;
  color: var(--label-color);
}

.profile-info {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.rating {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.profile-details {
  display: flex;
  flex-direction: column;
}

.form-columns {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
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

p {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: var(--input-background-color);
  color: var(--input-text-color);
  width: 100%;
  box-sizing: border-box;
  margin: 0;
}

.mask-star-2 {
  width: 24px;
  height: 24px;
  margin-right: 4px;
}
</style>
