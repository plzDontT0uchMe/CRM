import axios from 'axios'

const customAxios = axios.create({
    url: '',
    baseURL: '',
    withCredentials: false,
})
export default customAxios