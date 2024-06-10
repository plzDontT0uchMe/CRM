<template>
  <div :class="['app', theme]">
    <Header />
    <div class="content">
      <Menu />
      <div class="workout-plans-content">
        <button class="btn btn-primary mb-6" @click="showCreateForm = true">
          {{ $t('createNewPlan') }}
        </button>
        <div v-if="showCreateForm" class="create-plan-form">
          <h2 class="text-2xl font-bold mb-4">{{ $t('newPlan') }}</h2>
          <form @submit.prevent="createPlan">
            <div class="form-group">
              <label for="planName">{{ $t('planName') }}</label>
              <input
                id="planName"
                v-model="newPlan.name"
                type="text"
                class="input input-custom w-full"
                required
              />
            </div>
            <div class="form-group">
              <label for="planDescription">{{ $t('planDescription') }}</label>
              <textarea
                id="planDescription"
                v-model="newPlan.description"
                class="textarea textarea-custom w-full"
                required
              ></textarea>
            </div>
            <div class="form-group">
              <label for="planExercises">{{ $t('planExercises') }}</label>
              <select
                id="planExercises"
                v-model="selectedExercise"
                class="select select-custom w-full"
              >
                <option v-for="exercise in exercises" :key="exercise.id" :value="exercise">
                  {{ exercise.name }}
                </option>
              </select>
              <button type="button" class="btn btn-secondary mt-2" @click="addExercise">
                {{ $t('addExercise') }}
              </button>
            </div>
            <div class="selected-exercises">
              <h3 class="text-xl font-bold mt-4">{{ $t('selectedExercises') }}</h3>
              <ul>
                <li v-for="exercise in newPlan.exercises" :key="exercise.id">
                  {{ exercise.name }} - {{ exercise.description }} ({{
                    exercise.muscles.join(', ')
                  }})
                </li>
              </ul>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-success mt-4">{{ $t('addPlan') }}</button>
              <button
                type="button"
                class="btn btn-secondary mt-4 ml-2"
                @click="showCreateForm = false"
              >
                {{ $t('cancel') }}
              </button>
            </div>
          </form>
        </div>
        <div v-if="loading" class="loading">{{ $t('loading') }}</div>
        <div v-else class="workout-plans-list">
          <WorkoutPlanCard
            v-for="plan in workoutPlans"
            :key="plan.id"
            :plan="plan"
            :selected="selectedPlanId === plan.id"
            @open="openPlan"
          />
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
import WorkoutPlanCard from '../components/WorkoutPlanCard.vue'

const store = useStore()
const theme = computed(() => store.state.theme)
const { t } = useI18n()

const exercises = ref([
  {
    id: 1,
    name: 'Push Up',
    description: 'A basic upper body strength exercise.',
    muscles: ['Chest', 'Shoulders', 'Triceps']
  },
  {
    id: 2,
    name: 'Squat',
    description: 'A lower body exercise that strengthens the legs and glutes.',
    muscles: ['Quadriceps', 'Glutes', 'Hamstrings']
  },
  {
    id: 3,
    name: 'Plank',
    description: 'An exercise for core strength and stability.',
    muscles: ['Abdominals', 'Lower Back', 'Shoulders']
  },
  {
    id: 4,
    name: 'Bench Press',
    description: 'A compound exercise for chest and triceps.',
    muscles: ['Chest', 'Triceps']
  },
  {
    id: 5,
    name: 'Pull Up',
    description: 'An exercise for back and biceps strength.',
    muscles: ['Back', 'Biceps']
  },
  {
    id: 6,
    name: 'Shoulder Press',
    description: 'An exercise targeting the shoulders and triceps.',
    muscles: ['Shoulders', 'Triceps']
  }
])

const workoutPlans = ref([
  {
    id: 1,
    name: 'Beginner Full Body Workout',
    description: 'A great workout plan for beginners to target all major muscle groups.',
    exercises: [
      {
        id: 1,
        name: 'Push Up',
        description: 'A basic upper body strength exercise.',
        muscles: ['Chest', 'Shoulders', 'Triceps']
      },
      {
        id: 2,
        name: 'Squat',
        description: 'A lower body exercise that strengthens the legs and glutes.',
        muscles: ['Quadriceps', 'Glutes', 'Hamstrings']
      },
      {
        id: 3,
        name: 'Plank',
        description: 'An exercise for core strength and stability.',
        muscles: ['Abdominals', 'Lower Back', 'Shoulders']
      }
    ]
  },
  {
    id: 2,
    name: 'Advanced Upper Body Workout',
    description: 'An advanced workout plan focusing on upper body strength.',
    exercises: [
      {
        id: 4,
        name: 'Bench Press',
        description: 'A compound exercise for chest and triceps.',
        muscles: ['Chest', 'Triceps']
      },
      {
        id: 5,
        name: 'Pull Up',
        description: 'An exercise for back and biceps strength.',
        muscles: ['Back', 'Biceps']
      },
      {
        id: 6,
        name: 'Shoulder Press',
        description: 'An exercise targeting the shoulders and triceps.',
        muscles: ['Shoulders', 'Triceps']
      }
    ]
  }
])

const selectedPlanId = ref(null)
const showCreateForm = ref(false)
const newPlan = ref({
  name: '',
  description: '',
  exercises: []
})
const selectedExercise = ref(null)
const loading = ref(false)

const openPlan = (planId) => {
  selectedPlanId.value = selectedPlanId.value === planId ? null : planId
}

const addExercise = () => {
  if (selectedExercise.value) {
    newPlan.value.exercises.push(selectedExercise.value)
    selectedExercise.value = null
  }
}

const createPlan = () => {
  const newId = workoutPlans.value.length + 1
  const plan = {
    id: newId,
    name: newPlan.value.name,
    description: newPlan.value.description,
    exercises: newPlan.value.exercises
  }
  workoutPlans.value.push(plan)
  newPlan.value = {
    name: '',
    description: '',
    exercises: []
  }
  showCreateForm.value = false
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

.workout-plans-content {
  flex-grow: 1;
  padding: 20px;
  text-align: center;
}

.loading {
  font-size: 1.5rem;
  font-weight: bold;
}

.workout-plans-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 20px;
}

.create-plan-form {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 10px;
}

.selected-exercises ul {
  list-style: none;
  padding: 0;
}

.selected-exercises li {
  margin-bottom: 5px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.input-custom,
.textarea-custom,
.select-custom {
  background-color: #fff;
  border: 1px solid #ccc;
  color: #333;
  padding: 0.5rem;
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
}

.input-custom:focus,
.textarea-custom:focus,
.select-custom:focus {
  border-color: #007bff;
  box-shadow: 0 0 5px rgba(0, 123, 255, 0.5);
  outline: none;
}

.card {
  max-width: 100%;
  flex: 1 1 calc(25% - 20px); /* 4 cards per row, 20px gap */
}

.card-expanded {
  transform: scale(1.05);
}

@media (max-width: 1200px) {
  .card {
    flex: 1 1 calc(33.33% - 20px); /* 3 cards per row on smaller screens */
  }
}

@media (max-width: 900px) {
  .card {
    flex: 1 1 calc(50% - 20px); /* 2 cards per row on smaller screens */
  }
}

@media (max-width: 600px) {
  .card {
    flex: 1 1 calc(100% - 20px); /* 1 card per row on smallest screens */
  }
}
</style>
