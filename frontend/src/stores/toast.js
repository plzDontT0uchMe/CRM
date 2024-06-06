import { defineStore } from 'pinia'
import { toast } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'

const positionsToast = {
    'top-left': toast.POSITION.TOP_LEFT,
    'top-center': toast.POSITION.TOP_CENTER,
    'top-right': toast.POSITION.TOP_RIGHT,
    'bottom-left': toast.POSITION.BOTTOM_LEFT,
    'bottom-center': toast.POSITION.BOTTOM_CENTER,
    'bottom-right': toast.POSITION.BOTTOM_RIGHT
}

export const useToastStore = defineStore({
    id: 'toast',
    state: () => ({}),
    actions: {
        startToast(type, message, position) {
            switch (type) {
                case 'loading':
                    return toast.loading(
                        message,
                        {
                            position: positionsToast[position],
                            autoClose: false
                        }
                    )
                case 'success':
                    return toast.success(
                        message,
                        {
                            position: positionsToast[position],
                            autoClose: 1500,
                            closeOnClick: true,
                            closeButton: true,
                            isLoading: false,
                            pauseOnHover: false
                        }
                    )
                case 'error':
                    return toast.error(
                        message,
                        {
                            position: positionsToast[position],
                            autoClose: 1500,
                            closeOnClick: true,
                            closeButton: true,
                            isLoading: false,
                            pauseOnHover: false
                        }
                    )
                case 'info':
                    return toast.info(
                        message,
                        {
                            position: positionsToast[position],
                            autoClose: 1500,
                            closeOnClick: true,
                            closeButton: true,
                            isLoading: false,
                            pauseOnHover: false
                        }
                    )
                default:
                    return toast(
                        message,
                        {
                            position: positionsToast[position],
                            autoClose: 1500,
                            closeOnClick: true,
                            closeButton: true,
                            isLoading: false,
                            pauseOnHover: false
                        }
                    )
            }
        },
        stopToast(toastId, message, type) {
            setTimeout(() => {
                toast.update(toastId, {
                    render: message,
                    autoClose: 1500,
                    closeOnClick: true,
                    closeButton: true,
                    type: type,
                    isLoading: false,
                    pauseOnHover: false
                })
            }, 1000)
        }
    }
})