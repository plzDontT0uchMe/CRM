import axios from 'axios'

const customAxios = axios.create({
    url: '',
    baseURL: import.meta.env.MODE == 'development' ? "http://127.0.0.1:3000" : "http://45.156.23.236:3000",
    withCredentials: false,
    validateStatus: () => true,
})
export default customAxios