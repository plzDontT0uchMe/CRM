<script setup>
import axios from '@/axios/index.js'
import { ref } from 'vue'

const items = ref([])

const getItems = async () => {
    try {
        const resp = await axios.get('http://127.0.0.1:3001/api/users')
        console.log(resp)
        items.value = resp.data
    } catch (error) {
        console.error(error)
    }
}

</script>

<template>
    <div class="cursor-pointer h-10 w-10 bg-lime-500 m-16" @click="getItems">

    </div>
</template>

<style scoped>
</style>