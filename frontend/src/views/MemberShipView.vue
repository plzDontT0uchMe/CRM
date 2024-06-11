<template>
  <div v-if="isLoaded" :class="['app', theme]">
    <Header />
    <div class="content">
      <Menu />
      <div class="subscription-content">
        <h1>{{ $t('memberShip') }}</h1>
        <div class="duration-selector-container">
          <div class="duration-selector">
            <button
              v-for="(duration, index) in durations"
              :key="index"
              :class="{ active: selectedDuration === duration }"
              @click="selectDuration(duration)"
            >
              {{ $t(`durations.${duration}`) }}
            </button>
          </div>
        </div>
        <div class="subscriptions">
          <SubscriptionCard
            v-for="subscription in subscriptions"
            :key="subscription.id"
            :id="subscription.id"
            :name="subscription.name"
            :price="subscription.price"
            :duration="selectedDuration"
            :features="subscription.possibilities"
            :color="subscription.color"
            :selectedDuration="selectedDuration"
            :allFeatures="allFeatures"
            :selectedTrainer="selectedTrainer"
            @show-trainer-modal="openTrainerModal"
          />
        </div>
      </div>
    </div>

    <!-- DaisyUI Modal -->
    <input type="checkbox" id="trainer-modal" class="modal-toggle" v-model="showTrainerModal" />
    <div class="modal">
      <div class="modal-box relative">
        <button class="btn-close" @click="showTrainerModal = false">✕</button>
        <h2>{{ $t('chooseTrainer') }}</h2>
        <ul>
          <li
            v-for="trainer in trainers"
            :key="trainer.id"
            class="trainer-item border border-gray-300 p-4 rounded-lg mb-4"
          >
            <img :src="axios.defaults.baseURL + '/api/getImage/' + trainer.image" alt="trainer-avatar" class="trainer-avatar" />
            <div class="trainer-info">
              <h3>{{ trainer.name }} {{trainer.surname}} {{trainer.patronymic}}</h3>
              <p>{{ $t('age') }}: {{ trainer.dateBorn }}</p>
              <p>{{ $t('experience') }}: {{ trainer.trainerInfo[0].exp }} {{ $t('years') }}</p>
              <p>{{ $t('gender') }}: {{ trainer.gender }}</p>
              <p>{{ $t('sportType') }}: {{ trainer.trainerInfo[0].sport }}</p>
              <p>{{ $t('achievements') }}: {{ trainer.trainerInfo[0].achievements }}</p>
              <div class="rating">
                <span
                  v-for="star in 5"
                  :key="star"
                  class="mask mask-star-2"
                  :class="{ 'bg-warning': star <= Math.round(trainer.rating) }"
                ></span>
              </div>
              <button @click="selectTrainer(trainer)" class="btn btn-primary mt-2">
                {{ $t('choose') }}
              </button>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
    <span v-else class="loading"></span>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeMount } from 'vue'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
import Header from '../components/Header.vue'
import Menu from '../components/Menu.vue'
import SubscriptionCard from '../components/SubscriptionCard.vue'
import axios from '@/axios/index.js'
import { useUserStore } from '@/stores/user.js'

const userStore = useUserStore()
const subscriptions = ref('')

const getSubs = async () => {
    try {
        const resp = await axios.get('/api/getSubs')
        if (resp.data.successfully) {
            subscriptions.value = resp.data.subscriptions
            allFeatures.value = resp.data.subscriptions[2].possibilities
        }
    }
    catch (err) {
        console.log(err)
    }
}

const getTrainers = async () => {
    try {
        const resp = await axios.get('/api/getTrainers')
        if (resp.data.successfully) {
            trainers.value = resp.data.trainers
        }
    }
    catch (err) {
        console.log(err)
    }
}

const isLoaded = ref(false)

onMounted(async () => {
    isLoaded.value = false
    await getSubs()
    await getTrainers()
    selectedTrainer.value = userStore.data.subscription.trainer
    isLoaded.value = true
})

const store = useStore()
const theme = computed(() => store.state.theme)
const { t } = useI18n()

const durations = ['oneMonth', 'threeMonths', 'sixMonths', 'twelveMonths']
const selectedDuration = ref('oneMonth')

const selectDuration = (duration) => {
  selectedDuration.value = duration
}

const allFeatures = ref('')

const showTrainerModal = ref(false)
const selectedTrainer = ref(null)

const trainers = ref([
  {
    id: 1,
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
    fullName: 'John A. Doe',
    description: 'Expert in weightlifting and strength training.'
  },
  {
    id: 2,
    firstName: 'Jane',
    lastName: 'Smith',
    middleName: 'B.',
    dateOfBirth: '1985-02-15',
    gender: 'Female',
    sportType: 'Cardio',
    achievements: 'Marathon Winner 2015, 2017',
    avatar: 'https://via.placeholder.com/150',
    rating: 3,
    experience: 8,
    fullName: 'Jane B. Smith',
    description: 'Specializes in cardio and endurance training.'
  },
  {
    id: 3,
    firstName: 'Alice',
    lastName: 'Johnson',
    middleName: 'C.',
    dateOfBirth: '1990-03-30',
    gender: 'Female',
    sportType: 'Yoga',
    achievements: 'Yoga Instructor of the Year 2018',
    avatar: 'https://via.placeholder.com/150',
    rating: 5,
    experience: 5,
    fullName: 'Alice C. Johnson',
    description: 'Yoga and flexibility expert.'
  }
])

const openTrainerModal = () => {
  showTrainerModal.value = true
}

const selectTrainer = (trainer) => {
  selectedTrainer.value = trainer
  showTrainerModal.value = false
}

const getAge = (dateOfBirth) => {
  const today = new Date()
  const birthDate = new Date(dateOfBirth)
  let age = today.getFullYear() - birthDate.getFullYear()
  const monthDifference = today.getMonth() - birthDate.getMonth()
  if (monthDifference < 0 || (monthDifference === 0 && today.getDate() < birthDate.getDate())) {
    age--
  }
  return age
}
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

.subscription-content {
  flex-grow: 1;
  padding: 20px;
  text-align: center;
}

.duration-selector-container {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.duration-selector {
  display: flex;
  justify-content: center;
  border: 1px solid #ccc;
  border-radius: 20px;
  overflow: hidden;
  width: auto;
}

.duration-selector-container {
  padding: 10px;
  width: fit-content;
  margin: 10px auto;
}

.duration-selector button {
  background-color: transparent;
  color: inherit;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
  font-size: 0.875em;
  transition:
    background-color 0.3s,
    border-bottom 0.3s;
  border-bottom: 2px solid transparent;
}

.duration-selector button.active {
  background-color: transparent;
  border-bottom: 2px solid #007bff;
  color: #007bff;
}

.subscriptions {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
}

.subscription-card {
  width: 400px;
  height: 450px; /* Уменьшена высота карточек */
}

button {
  margin: 20px;
  padding: 10px 20px;
  cursor: pointer;
}

.trainer-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  border: 1px solid #ccc; /* Рамка вокруг каждого тренера */
  border-radius: 8px; /* Радиус для закругленных углов */
  padding: 20px; /* Увеличен padding для тренеров */
}

.trainer-avatar {
  width: 100px; /* Увеличен размер фотографий тренеров */
  height: 100px; /* Увеличен размер фотографий тренеров */
  border-radius: 50%;
  margin-right: 20px; /* Увеличен отступ справа от фотографии */
  object-fit: cover; /* Пропорциональное отображение фотографий */
}

.trainer-info {
  display: flex;
  flex-direction: column;
}

.modal-box {
  max-width: 90%;
  overflow-y: auto;
  background-color: var(--modal-background-color, white); /* Цвет фона модального окна */
  color: var(--modal-text-color, black); /* Цвет текста модального окна */
  position: relative;
}

.modal-box .btn-close {
  position: absolute;
  top: -16px; /* Вверх и правее */
  right: -16px; /* Вверх и правее */
  background: none;
  border: none;
  color: var(--button-background-color);
  font-size: 24px;
  cursor: pointer;
}

.modal-action .btn {
  margin-top: 20px;
  background-color: var(--button-background-color);
  color: var(--button-text-color);
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

.rating {
  display: flex;
  margin-top: 10px;
}

.mask-star-2 {
  width: 24px;
  height: 24px;
  margin-right: 4px;
}
</style>
