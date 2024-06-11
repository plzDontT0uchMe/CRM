<template>
  <div
    class="card w-96 bg-white shadow-xl cursor-pointer transform transition-transform hover:scale-105"
    @click="$emit('open-modal', exercise)"
  >
    <figure>
      <img :src="axios.defaults.baseURL + '/api/getImage/' + exercise.image" alt="exercise-image" class="exercise-image" />
    </figure>
    <div class="card-body">
      <h2 class="card-title">{{ exercise.name }}</h2>
      <p>{{ exercise.description }}</p>
      <div class="card-actions justify-end">
        <div v-for="muscle in exercise.muscles" :key="muscle" class="badge badge-outline">
          {{ muscle }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '@/axios/index.js'

const props = defineProps({
  exercise: Object
})

const { t } = useI18n()
</script>

<style scoped>
.exercise-image {
  width: 100%;
  height: auto;
  object-fit: cover;
}

.card {
  transition:
    transform 0.2s ease-in-out,
    box-shadow 0.2s ease-in-out;
}

.card:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.btn {
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

.btn:hover {
  background-color: var(--button-hover-background-color);
}

.badge-outline {
  border: 1px solid var(--button-background-color);
  color: var(--button-background-color);
}
</style>
