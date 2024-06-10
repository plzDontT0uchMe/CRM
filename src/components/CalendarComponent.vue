<template>
  <div>
    <FullCalendar :options="calendarOptions" @dateClick="handleDateClick" />
    <div v-if="showModal" class="modal">
      <div class="modal-content">
        <span class="close" @click="closeModal">&times;</span>
        <h2>Add Note</h2>
        <form @submit.prevent="addEvent">
          <label for="title">Title:</label>
          <input type="text" id="title" v-model="noteTitle" required />
          <label for="description">Description:</label>
          <textarea id="description" v-model="noteDescription"></textarea>
          <button type="submit">Add Note</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { FullCalendar } from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'

export default defineComponent({
  name: 'CalendarComponent',
  components: {
    FullCalendar
  },
  setup() {
    const calendarOptions = ref({
      plugins: [dayGridPlugin, interactionPlugin],
      initialView: 'dayGridMonth',
      events: [],
      headerToolbar: {
        left: 'prev,next today',
        center: 'title',
        right: 'dayGridMonth,dayGridWeek,dayGridDay'
      },
      locale: 'en'
    })

    const showModal = ref(false)
    const noteTitle = ref('')
    const noteDescription = ref('')
    const selectedDate = ref(null)

    const handleDateClick = (info) => {
      selectedDate.value = info.dateStr
      showModal.value = true
    }

    const closeModal = () => {
      showModal.value = false
      noteTitle.value = ''
      noteDescription.value = ''
    }

    const addEvent = () => {
      if (noteTitle.value) {
        calendarOptions.value.events.push({
          title: noteTitle.value,
          date: selectedDate.value,
          description: noteDescription.value
        })
        closeModal()
      }
    }

    return {
      calendarOptions,
      showModal,
      noteTitle,
      noteDescription,
      handleDateClick,
      closeModal,
      addEvent
    }
  }
})
</script>

<style>
@import '@fullcalendar/common/main.css';
@import '@fullcalendar/daygrid/main.css';

.fullcalendar {
  max-width: 100%;
  margin: 0 auto;
}

.modal {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  z-index: 1;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.4);
}

.modal-content {
  background-color: #fefefe;
  padding: 20px;
  border: 1px solid #888;
  width: 400px;
  position: relative;
}

.close {
  position: absolute;
  top: 10px;
  right: 10px;
  color: #aaa;
  font-size: 28px;
  font-weight: bold;
  cursor: pointer;
}

.close:hover,
.close:focus {
  color: black;
  text-decoration: none;
  cursor: pointer;
}
</style>
