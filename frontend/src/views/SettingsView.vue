<script setup>
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast.js'
import { ref } from 'vue'
import axios from '@/axios/index.js'
import moment from '@/moment/index.js'
import WomenIcon from '@/components/icons/WomenIcon.vue'
import UserIcon from '@/components/icons/UserIcon.vue'
import ManIcon from '@/components/icons/ManIcon.vue'
import { toast } from 'vue3-toastify'

const userStore = useUserStore()
const toastStore = useToastStore()
const user = ref({ ...userStore.data })

const image = ref({
    url: "",
    file: null
})

const updateUser = async () => {
    const notifyId = toastStore.startToast('loading', 'Сохранение изменений... 🚀', 'top-center')
    const data = new FormData()
    for (const key in user.value) {
        if (key === 'image') {
            if (image.value.file) {
                data.append('image', image.value.file)
                continue
            }
        }
        data.append(key, user.value[key])
    }
    try {
        const resp = await axios.post('/api/updateUser', data)
        toastStore.stopToast(notifyId, resp.data.message, resp.data.successfully ? 'success' : 'error')
        if (resp.data.successfully) {
            await userStore.fetchUser()
        }
    } catch (err) {
        console.log(err)
    }
}

const onFileChange = (e) => {
    image.value = {
        file: e.target.files[0],
        url: URL.createObjectURL(e.target.files[0])
    }
}

const hoverImage = ref(null)
</script>

<template>
    <div class="w-screen h-[80vh] flex flex-col justify-center items-center space-y-2">
        <div class="avatar">
            <div class="w-32 rounded">
                <div class="relative w-full h-full" v-if="user.image || image.url"
                     @mouseover="hoverImage === null ? hoverImage = image.url : hoverImage"
                     @mouseleave="hoverImage = null"
                >
                    <img v-if="image.url" alt="image" :src="image.url" />
                    <img v-else alt="image" :src="axios.defaults.baseURL + '/api/getImage/' + userStore?.data?.image" />
                    <transition name="show">
                        <div v-show="hoverImage === image.url"
                             class="absolute top-0 right-0">
                            <button class="btn btn-square btn-outline btn-xs btn-error" @click="image = {}">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none"
                                     viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                          d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
                        </div>
                    </transition>
                </div>
                <div v-else>
                    <ManIcon v-if="user.gender == 1" />
                    <WomenIcon v-else-if="user.gender == 2" />
                    <UserIcon v-else />
                </div>
            </div>
        </div>
        <div class="flex flex-col justify-center items-center space-y-2 max-w-96">
            <input type="file" class="file-input file-input-bordered w-full"
                   @change="onFileChange" />
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
            <p>{{ moment(user.lastActivity).fromNow() }}</p>
        </div>
        <button class="btn btn-primary btn-outline" @click="updateUser">Сохранить</button>
    </div>
</template>

<style scoped>
.show-enter-active,
.show-leave-active {
    transition: scale 0.3s ease-in-out;
    transition-delay: 0.2s;
}

.show-enter-from,
.show-leave-to {
    transition-delay: 0s;
    scale: 0;
}
</style>