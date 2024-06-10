<template>
  <nav class="menu">
    <ul>
      <li>
        <router-link to="/" active-class="active-link">{{ $t('calendar') }}</router-link>
      </li>
      <li>
        <router-link to="/membership" active-class="active-link">{{
          $t('memberShip')
        }}</router-link>
      </li>
      <li>
        <router-link to="/exercise" active-class="active-link">{{ $t('exercises') }}</router-link>
      </li>
      <li>
        <router-link to="/workout-plans">{{ $t('workoutPlans') }}</router-link>
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<style scoped>
.menu {
  background-color: var(--menu-background-color);
  color: var(--menu-text-color);
  padding: 20px;
  width: 220px;
  flex-shrink: 0;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
}

ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

li {
  margin-bottom: 15px;
}

li a {
  color: var(--menu-text-color);
  text-decoration: none;
  display: block;
  padding: 12px;
  border-radius: 8px;
  transition:
    background-color 0.3s,
    color 0.3s;
}

li a:hover,
li a.active-link {
  background-color: var(--secondary-color);
  color: white;
}
</style>
