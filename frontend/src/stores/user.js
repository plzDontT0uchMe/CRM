import { defineStore } from 'pinia'
import axios from '@/axios/index.js'

export const useUserStore = defineStore({
    id: 'user',
    state: () => ({
        data: null,
        loading: true,
    }),
    actions: {
        async fetchUser() {
            this.loading = true
            try {
                const { data } = await axios.get('/api/getUser')
                if (data.successfully) {
                    this.data = data.user;
                }
            } catch (err) {
                console.log(err)
            } finally {
                this.loading = false
            }
        }
    }
})