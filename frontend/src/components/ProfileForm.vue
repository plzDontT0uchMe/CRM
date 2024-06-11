<script setup xmlns="http://www.w3.org/1999/html">

import { onMounted, ref, watch } from 'vue'
import axios from '@/axios/index.js'
import { useRoute } from 'vue-router'
import router from '@/router/index.js'
import moment from '@/moment/index.js'
import WomenIcon from '@/components/icons/WomenIcon.vue'
import UserIcon from '@/components/icons/UserIcon.vue'
import ManIcon from '@/components/icons/ManIcon.vue'

const route = useRoute()

const user = ref(null)
const isLoaded = ref(false)

const getUser = async () => {
    try {
        const { data } = await axios.get(`/api/getUser/${route.params.id}`)
        if (data.successfully) {
            user.value = data.user
            isLoaded.value = true
            return
        }
        router.push('/notfound')
    } catch (err) {
        console.log(err)
    }
}

onMounted(() => {
    getUser()
})

watch(() => route.params.id, () => {
    getUser()
})

import { computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

</script>

<template>
    <div class="profile-form" v-if="isLoaded">
        <div class="profile-header">
            <div class="profile-photo-container">
                <div class="profile-photo">
                    <img v-if="user.image" :src="axios.defaults.baseURL + '/api/getImage/' + user?.image" />
                    <div v-else>
                        <ManIcon v-if="user?.gender == 1" />
                        <WomenIcon v-else-if="user?.gender == 2" />
                        <UserIcon v-else />
                    </div>
                </div>
                <span
                    v-for="star in 5"
                    :key="star"
                    class="mask mask-star-2"
                    :class="{ 'bg-warning': star <= rating }"
                ></span>
                <p class="last-active">{{ $t('lastActive') }}: {{moment(user.lastActivity).fromNow()}}</p>
            </div>
            <div class="profile-info"></div>
        </div>
        <div class="profile-details" v-if="isLoaded">
            <div class="form-columns">
                <div class="form-column">
                    <div class="form-row">
                        <label>{{ $t('firstName') }}:</label>
                        <p>{{ user.name }}</p>
                    </div>
                    <div class="form-row">
                        <label>{{ $t('lastName') }}:</label>
                        <p>{{ user.surname }}</p>
                    </div>
                    <div class="form-row">
                        <label>{{ $t('gender') }}:</label>
                        <p>{{ user.gender }}</p>
                    </div>
                    <div class="form-row">
                        <label>{{ $t('accountCreated') }}:</label>
                        <p>{{ user.dateCreated }}</p>
                    </div>
                </div>
                <div class="form-column">
                    <div class="form-row">
                        <label>{{ $t('middleName') }}:</label>
                        <p>{{ user.patronymic }}</p>
                    </div>
                    <div class="form-row">
                        <label>{{ $t('dob') }}:</label>
                        <p>{{ user.dateBorn }}</p>
                    </div>
                    <div class="form-row">
                        <label>{{ $t('role') }}:</label>
                        <p>{{ user.position }}</p>
                    </div>
                    <template v-if="user.position === 'trainer'">
                        <div class="form-row">
                            <label>{{ $t('experience') }}:</label>
                            <p>{{ user.trainerInfo[0].exp }} {{ $t('years') }}</p>
                        </div>
                        <div class="form-row">
                            <label>{{ $t('sportType') }}:</label>
                            <p>{{ user.trainerInfo[0].sport }}</p>
                        </div>
                        <div class="form-row">
                            <label>{{ $t('achievements') }}:</label>
                            <p>{{ user.trainerInfo[0].achievements }}</p>
                        </div>
                    </template>
                </div>
            </div>
        </div>
    </div>
    <span class="loading loading-lg" v-else></span>


    <!--    <div class="w-screen h-[80vh] flex flex-col justify-center items-center space-y-2" v-if="isLoaded">
            <div class="avatar">
                <div class="w-32 rounded">
                    <img v-if="user.image" :src="axios.defaults.baseURL + '/api/getImage/' + user?.image" />
                    <div v-else>
                        <ManIcon v-if="user?.gender == 1"/>
                        <WomenIcon v-else-if="user?.gender == 2"/>
                        <UserIcon v-else />
                    </div>
                </div>
            </div>
            <div class="flex flex-col justify-center items-center space-y-2 max-w-96">
                <label class="input input-bordered flex items-center gap-2 w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
                         class="w-4 h-4 opacity-70">
                        <path
                            d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
                    </svg>
                    <input type="text" class="grow" placeholder="Name" v-model="user.name" />
                </label>
                <label class="input input-bordered flex items-center gap-2 w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
                         class="w-4 h-4 opacity-70">
                        <path
                            d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
                    </svg>
                    <input type="text" class="grow" placeholder="Surname" v-model="user.surname" />
                </label>
                <label class="input input-bordered flex items-center gap-2 w-full">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
                         class="w-4 h-4 opacity-70">
                        <path
                            d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
                    </svg>
                    <input type="text" class="grow" placeholder="Patronymic" v-model="user.patronymic" />
                </label>
                <select class="select select-bordered w-full" v-model.number="user.gender">
                    <option value="0" disabled selected>...</option>
                    <option value="1">Мужчина</option>
                    <option value="2">Женщина</option>
                </select>
                <input type="date" class="input input-bordered w-full" v-model="user.dateBorn" />
                <p>Дата создания: {{ moment(user.dateCreated).format('DD-MM-YYYY') }}</p>
                <p>Последняя активность: {{ moment(user.lastActivity).format('DD-MM-YYYY HH:mm:ss') }}</p>
                <p>{{moment(user.lastActivity).fromNow()}}</p>
            </div>
        </div>
        <span class="loading loading-lg" v-else></span>-->
</template>

<style scoped>
.profile-form {
    display: flex;
    flex-direction: column;
    width: 100%;
    background-color: var(--background-color);
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    color: var(--text-color);
}

.profile-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
}

.profile-photo-container {
    text-align: center;
    margin-right: 20px;
    flex-shrink: 0;
    position: relative;
}

.profile-photo {
    width: 150px;
    height: 150px;
    background-color: #ddd;
    display: block;
    margin-bottom: 10px;
    border-radius: 50%;
}

.last-active {
    margin-top: 10px;
    font-size: 0.875em;
    color: var(--label-color);
}

.profile-info {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
}

.rating {
    display: flex;
    align-items: center;
    margin-top: 10px;
}

.profile-details {
    display: flex;
    flex-direction: column;
}

.form-columns {
    display: flex;
    justify-content: space-between;
    margin-bottom: 15px;
}

.form-column {
    flex: 1;
    margin-right: 10px;
}

.form-column:last-child {
    margin-right: 0;
}

.form-row {
    display: flex;
    flex-direction: column;
    margin-bottom: 15px;
}

label {
    margin-bottom: 5px;
    color: var(--label-color);
    font-weight: bold;
}

p {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
    background-color: var(--input-background-color);
    color: var(--input-text-color);
    width: 100%;
    box-sizing: border-box;
    margin: 0;
}

.mask-star-2 {
    width: 24px;
    height: 24px;
    margin-right: 4px;
}
</style>