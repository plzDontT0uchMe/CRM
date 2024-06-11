<script setup>

import { onMounted, ref, watch } from 'vue'
import axios from '@/axios/index.js'
import { useRoute } from 'vue-router'
import router from '@/router/index.js'
import moment from '@/moment/index.js'
import WomenIcon from '@/components/icons/WomenIcon.vue'
import UserIcon from '@/components/icons/UserIcon.vue'
import ManIcon from '@/components/icons/ManIcon.vue'

const route = useRoute();

const user = ref(null)
const isLoaded = ref(false)

/*const getUser = async () => {
    try {
        const { data } = await axios.get(`/api/getUser/${route.params.id}`)
        if (data.successfully) {
            user.value = data.user
            isLoaded.value = true
            return
        }
        router.push('/notfound')
    }
    catch (err) {
        console.log(err)
    }
}

onMounted(() => {
    getUser()
})

watch(() => route.params.id, () => {
    getUser()
})*/

import Header from '../components/Header.vue'
import Menu from '../components/Menu.vue'
import ProfileForm from '../components/ProfileForm.vue'

</script>

<template>
    <div class="profile-container">
        <Header />
        <div class="content">
            <Menu />
            <div class="profile-content">
                <ProfileForm />
            </div>
        </div>
    </div>




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
.profile-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
}

.content {
    display: flex;
    flex: 1;
    background-color: var(--content-background-color);
    color: var(--content-text-color);
}

.menu {
    height: 100%;
    flex-shrink: 0;
}

.profile-content {
    flex-grow: 1;
    padding: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
}
</style>