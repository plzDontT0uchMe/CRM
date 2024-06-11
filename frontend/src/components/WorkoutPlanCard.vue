<template>
  <div :class="['card', { 'card-expanded': selected }]" @click="handleClick">
    <div class="card-body">
      <h2 class="card-title">{{ plan.name }}</h2>
      <p>{{ plan.description }}</p>
    </div>
    <div v-if="selected" class="exercises-list">
      <div
        v-for="exercise in plan.exercises"
        :key="exercise.id"
        class="exercise-card bg-gray-100 p-4 mb-4"
      >
        <h4 class="font-bold">{{ exercise.name }}</h4>
        <p>{{ exercise.description }}</p>
        <div class="mt-2 flex flex-wrap gap-2">
          <div v-for="muscle in exercise.muscles" :key="muscle" class="badge badge-outline">
            {{ muscle }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  plan: Object,
  selected: Boolean
})

const emits = defineEmits(['open'])

const handleClick = () => {
  emits('open', props.plan.id)
}
</script>

<style scoped>
.card {
  transition:
    transform 0.2s,
    box-shadow 0.2s;
  width: 100%;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  margin: 10px;
}

.card:hover {
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.card-expanded {
  transform: scale(1.05);
}

.card-body {
  padding: 1rem;
}

.exercises-list {
  padding: 10px;
}

.exercise-card {
  margin-top: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.badge-outline {
  border: 1px solid var(--button-background-color);
  color: var(--button-background-color);
}
</style>
