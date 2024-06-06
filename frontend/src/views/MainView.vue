<script setup>
import axios from '@/axios/index.js'
import { onMounted, ref } from 'vue'

const user = ref({
    id: null,
    role: null,
    lastActivity: null,
    dateCreated: null,
    name: null,
    surname: null,
    patronymic: null,
    gender: null,
    dateBorn: null,
    image: null,
})

const formatDateWithoutTime = (dateStr) => {
    const date = new Date(dateStr)
    const day = date.getDate()
    const month = date.getMonth() + 1 // Месяцы начинаются с 0, поэтому добавляем 1
    const year = date.getFullYear()

    return `${day}.${month}.${year}`
}

const formatDateWithTime = (dateStr) => {
    const date = new Date(dateStr)
    const day = date.getDate()
    const month = date.getMonth() + 1 // Месяцы начинаются с 0, поэтому добавляем 1
    const year = date.getFullYear()
    const hours = date.getHours()
    const minutes = date.getMinutes()
    const seconds = date.getSeconds()

    return `${day}.${month}.${year} ${hours}:${minutes}:${seconds}`
}

const getUser = async () => {
    try {
        const { data } = await axios.get('/api/getUser')
        console.log(data)
        if (data.successfully) {
            user.value = data.user
            console.log(user.value)
        }

    }
    catch (err) {
        console.log(err)
    }
}

onMounted(() => {
    //getUser()
})

</script>

<template>
    <div class="flex justify-center items-center">
        <div class="p-4 bg-blue-200 rounded">
            <h2 class="text-lg font-bold">Дата:</h2>
            <p>{{ formatDateWithoutTime(user.dateCreated) }}</p>
        </div>
        <div class="p-4 bg-blue-200 rounded">
            <h2 class="text-lg font-bold">Дата:</h2>
            <p>{{ formatDateWithTime(user.lastActivity) }}</p>
        </div>
    </div>
</template>

<style scoped>

</style>