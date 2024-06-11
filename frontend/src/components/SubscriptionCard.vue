<template>
  <div class="subscription-card" :style="{ backgroundColor: color }">
    <h2>{{ $t(`plans.${name}`) }}</h2>
    <div class="price">{{ currentPrice }}</div>
    <ul>
      <li
        v-for="feature in allFeatures"
        :key="feature"
        :class="{ included: features.includes(feature) }"
      >
        <span class="icon">{{ features.includes(feature) ? '✔️' : '❌' }}</span>
        {{ $t(`features.${feature}`) }}
      </li>
    </ul>
    <div v-if="name === 'premium'" class="trainer-container">
      <div class="trainer-selection">
        <button @click="$emit('show-trainer-modal', true)" class="btn">
          {{ $t('chooseTrainer') }}
        </button>
        <p v-if="selectedTrainer" class="selected-trainer">{{ selectedTrainer.fullName }}</p>
      </div>
    </div>
    <div class="button-container mt-4">
      <button @click="buySubscription" class="btn">{{ $t('buy') }}</button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  name: String,
  price: Object,
  features: Array,
  color: String,
  selectedDuration: String,
  allFeatures: Array,
  selectedTrainer: Object
})

const { t } = useI18n()

const currentPrice = computed(() => props.price[props.selectedDuration])

const buySubscription = () => {
  alert(`Subscribed to ${t(`plans.${props.name}`)} Plan for ${props.selectedDuration}`)
}
</script>

<style scoped>
.subscription-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 20px;
  width: 400px;
  height: 450px;
  background-color: var(--card-background-color, white);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.subscription-card:hover {
  transform: scale(1.05);
}

.subscription-card h2 {
  margin: 10px 0;
  font-size: 1.5em;
}

.subscription-card .price {
  font-size: 1.25em;
  font-weight: bold;
  margin: 10px 0;
}

.subscription-card ul {
  list-style-type: none;
  padding: 0;
  text-align: left;
  margin: 0 0 10px 0;
  flex-grow: 1;
}

.subscription-card li {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
  color: #999;
  font-size: 0.95em;
}

.subscription-card li.included {
  color: #000;
  font-weight: bold;
}

.subscription-card .icon {
  margin-right: 10px;
}

.trainer-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin-top: 10px;
}

.trainer-container .trainer-selection {
  display: flex;
  align-items: center;
}

.trainer-container .selected-trainer {
  margin-left: 10px;
  font-size: 1em;
  color: var(--text-color);
}

.trainer-container button,
.subscription-card button {
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

.trainer-container button:hover,
.subscription-card button:hover {
  background-color: var(--button-hover-background-color);
}

.button-container.mt-4 {
  margin-top: 16px;
}
</style>
