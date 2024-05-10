<script setup>
import TheWelcome from '../components/TheWelcome.vue'
</script>

<template>
  <main class="bg-red-600">
    <TheWelcome />
      <div class="avatar">
          <div class="w-24 rounded">
              <img src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
          </div>
      </div>
      <div class="flex flex-col gap-4 w-52">
          <div class="flex gap-4 items-center">
              <div class="skeleton w-16 h-16 rounded-full shrink-0"></div>
              <div class="flex flex-col gap-4">
                  <div class="skeleton h-4 w-20"></div>
                  <div class="skeleton h-4 w-28"></div>
              </div>
          </div>
          <div class="skeleton h-32 w-full"></div>
      </div>
  </main>
</template>
