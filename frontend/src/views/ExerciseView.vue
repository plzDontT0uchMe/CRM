<template>
  <div :class="['app', theme]">
    <Header />
    <div class="content">
      <Menu />
      <div class="exercise-content">
        <h1 class="text-3xl font-bold mb-6">{{ $t('exercises') }}</h1>
        <div class="exercise-list">
          <ExerciseCard
            v-for="exercise in exercises"
            :key="exercise.id"
            :exercise="exercise"
            @open-modal="openModal"
          />
        </div>
      </div>
    </div>

    <!-- DaisyUI Modal -->
    <div v-if="showModal" class="modal modal-open">
      <div class="modal-box relative bg-white text-black modal-content">
        <button
          class="btn btn-sm btn-circle absolute right-2 top-2 text-black bg-transparent border-none"
          @click="showModal = false"
        >
          ✕
        </button>
        <h2 class="text-lg font-bold">{{ selectedExercise.name }}</h2>
        <img :src="axios.defaults.baseURL + '/api/getImage/' + selectedExercise.image" alt="exercise-image" class="w-full h-auto mt-4" />
        <h3 class="text-lg font-bold mt-4">Описание</h3>
        <p class="mt-2 description">{{ selectedExercise.description }}</p>
        <h3 class="text-lg font-bold mt-4">Мускулы</h3>
        <div class="mt-2 flex flex-wrap gap-2">
          <div v-for="muscle in selectedExercise.muscles" :key="muscle" class="badge badge-outline">
            {{ muscle }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
import Header from '../components/Header.vue'
import Menu from '../components/Menu.vue'
import ExerciseCard from '../components/ExerciseCard.vue'
import axios from '@/axios/index.js'

const store = useStore()
const theme = computed(() => store.state.theme)
const { t } = useI18n()

const showModal = ref(false)
const selectedExercise = ref({})

const openModal = (exercise) => {
  selectedExercise.value = exercise
  showModal.value = true
}

const getExercises = async () => {
    try{
        const resp = await axios.get("/api/getExercises")
        if (resp.data.successfully) {
            exercises.value = resp.data.exercises
        }
    } catch (error) {
        console.error(error)
    }
}

onMounted(async () => {
    await getExercises()
})

const exercises = ref([
  {
    id: 1,
    name: 'Push Up',
    image: 'https://via.placeholder.com/300',
    description: 'A basic upper body strength exercise.',
    muscles: ['Chest', 'Shoulders', 'Triceps']
  },
  {
    id: 2,
    name: 'Squat',
    image: 'https://via.placeholder.com/300',
    description: 'A lower body exercise that strengthens the legs and glutes.',
    muscles: ['Quadriceps', 'Glutes', 'Hamstrings']
  },
  {
    id: 3,
    name: 'Plank',
    image: 'https://via.placeholder.com/300',
    description: 'An exercise for core strength and stability.',
    muscles: ['Abdominals', 'Lower Back', 'Shoulders']
  }
])
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

.exercise-content {
  flex-grow: 1;
  padding: 20px;
  text-align: center;
}

.exercise-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
}

.modal-box {
  background-color: #f9f9f9; /* светлый фон */
  color: #000; /* черный текст */
  overflow: hidden; /* убираем прокрутку */
  display: flex;
  flex-direction: column;
  align-items: center;
}

.modal-content {
  max-width: 80%;
  max-height: 90%;
}

.modal-box img {
  max-height: 300px; /* ограничиваем высоту изображения */
  width: auto;
  position: left;
}

.modal-box p.description {
  max-height: 100px; /* ограничиваем высоту описания */
  overflow: hidden;
}

.modal-box .badge-outline {
  margin-bottom: 4px; /* добавляем немного пространства между бейджами */
}

.btn-circle {
  background-color: transparent;
  border: none;
  color: var(--button-background-color);
}

.btn-circle:hover {
  background-color: #e0e0e0;
}

/* Медиазапросы для адаптивности */
@media (min-width: 640px) {
  .modal-content {
    max-width: 60%;
    max-height: 70%;
  }
}

@media (min-width: 1024px) {
  .modal-content {
    max-width: 40%;
    max-height: 60%;
  }
}
</style>
