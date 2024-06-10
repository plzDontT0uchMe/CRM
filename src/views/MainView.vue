<template>
  <div class="main-container">
    <Header />
    <div class="content">
      <Menu />
      <div class="main-content">
        <DxScheduler
          class="calendar"
          :data-source="dataSource"
          :views="views"
          :current-view="currentView"
          :current-date="currentDate"
          :start-day-hour="6"
          :end-day-hour="22"
          height="700px"
          @appointment-added="onAppointmentAdded"
          @appointment-updated="onAppointmentUpdated"
          @appointment-deleted="onAppointmentDeleted"
          @appointment-dbl-click="onAppointmentDblClick"
        >
          <DxView type="day" name="Day" />
          <DxView type="week" name="Week" />
          <DxView type="workWeek" name="Work Week" />
          <DxView type="month" name="Month" />
          <DxEditing :allow-adding="true" :allow-updating="true" :allow-deleting="true" />
        </DxScheduler>
        <div v-if="showExerciseForm" class="exercise-form">
          <h3>Exercises Completed</h3>
          <ul>
            <li v-for="exercise in exercises" :key="exercise.id">
              {{ exercise.text }} ({{ exercise.startDate }} - {{ exercise.endDate }})
            </li>
          </ul>
          <button @click="closeExerciseForm">Close</button>
        </div>
        <div v-if="showAddForm" class="add-form">
          <h3>Add Exercise</h3>
          <form @submit.prevent="addExercise">
            <label>
              Exercise:
              <input v-model="newExercise.text" type="text" required />
            </label>
            <label>
              Start Time:
              <input v-model="newExercise.startDate" type="datetime-local" required />
            </label>
            <label>
              End Time:
              <input v-model="newExercise.endDate" type="datetime-local" required />
            </label>
            <button type="submit">Add</button>
          </form>
          <button @click="closeAddForm">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from '../components/Header.vue'
import Menu from '../components/Menu.vue'
import { DxScheduler, DxView, DxEditing } from 'devextreme-vue/scheduler'

const dataSource = ref([
  {
    text: 'Morning Workout',
    startDate: new Date('2024-06-10T08:00:00'),
    endDate: new Date('2024-06-10T10:00:00')
  },
  {
    text: 'Evening Yoga',
    startDate: new Date('2024-06-10T18:00:00'),
    endDate: new Date('2024-06-10T19:30:00')
  }
])

const views = ['day', 'week', 'workWeek', 'month']
const currentView = ref('month')
const currentDate = ref(new Date())
const showExerciseForm = ref(false)
const showAddForm = ref(false)
const exercises = ref([])
const newExercise = ref({
  text: '',
  startDate: '',
  endDate: ''
})

const onAppointmentAdded = (e) => {
  dataSource.value.push(e.appointmentData)
  console.log('Appointment added:', e.appointmentData)
}

const onAppointmentUpdated = (e) => {
  const index = dataSource.value.findIndex(
    (item) => item.startDate === e.oldData.startDate && item.endDate === e.oldData.endDate
  )
  if (index > -1) {
    dataSource.value.splice(index, 1, e.appointmentData)
  }
  console.log('Appointment updated:', e.appointmentData)
}

const onAppointmentDeleted = (e) => {
  const index = dataSource.value.findIndex(
    (item) =>
      item.startDate === e.appointmentData.startDate && item.endDate === e.appointmentData.endDate
  )
  if (index > -1) {
    dataSource.value.splice(index, 1)
  }
  console.log('Appointment deleted:', e.appointmentData)
}

const onAppointmentDblClick = (e) => {
  const clickedDate = new Date(e.appointmentData.startDate)
  if (clickedDate < new Date()) {
    exercises.value = dataSource.value.filter(
      (item) => item.startDate <= clickedDate && item.endDate >= clickedDate
    )
    showExerciseForm.value = true
  } else {
    newExercise.value.startDate = clickedDate.toISOString().slice(0, 16)
    newExercise.value.endDate = new Date(clickedDate.getTime() + 60 * 60 * 1000)
      .toISOString()
      .slice(0, 16) // default end time +1 hour
    showAddForm.value = true
  }
}

const addExercise = () => {
  dataSource.value.push({
    text: newExercise.value.text,
    startDate: new Date(newExercise.value.startDate),
    endDate: new Date(newExercise.value.endDate)
  })
  closeAddForm()
}

const closeExerciseForm = () => {
  showExerciseForm.value = false
}

const closeAddForm = () => {
  showAddForm.value = false
}
</script>

<style scoped>
.calendar {
  z-index: 1; /* Lower z-index to ensure forms appear above */
}
.main-container {
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

.main-content {
  flex-grow: 1;
  padding: 20px;
}

.exercise-form,
.add-form {
  position: absolute;
  top: 20%;
  left: 50%;
  transform: translate(-50%, -20%);
  background-color: white;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
  z-index: 10000; /* Ensure forms appear on top */
}
</style>
