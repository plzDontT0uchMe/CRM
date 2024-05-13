<script setup>
import SettingsIcon from "@/components/icons/SettingsIcon.vue";
import { ref } from "vue";

const time = ref(144);

const isStarted = ref(false);
const timer = ref();

const startTimer = () => {
    isStarted.value = true;
    timer.value = setInterval(() => {
        time.value--;
    }, 1000);
}
const stopTimer = () => {
    isStarted.value = false;
    clearInterval(timer.value)
    time.value = 60;
}
</script>

<template>
<div class="page">
    <div class="container">
        <div class="timerContainer">
            <div class="timer" @click="!isStarted ? startTimer() : stopTimer()"> {{ time }} </div>
            <div class="timerSetting">
                <SettingsIcon />
            </div>
        </div>
    </div>
</div>
</template>

<style scoped>
.page{
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100vw;
    height: 100vh;
    user-select: none;
}
.container{
    display: flex;
    justify-content: center;
    align-items: center;
    width: 414px;
    height: 896px;
    background: #242424;
}
.timerContainer{
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
}
.timer{
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    width: 320px;
    height: 320px;
    background: white;
    border-radius: 180px;
    font-size: 32px;
    margin-bottom: 8px;
}
.timerSetting{
    display: flex;
    justify-content: center;
    align-items: center;
    width: 64px;
    height: 64px;
    background: white;
    border-radius: 180px;
}
</style>