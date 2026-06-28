import { defineStore } from 'pinia'
import { ref, readonly } from 'vue'
import { api } from '@/schema'

export const useUserStore = defineStore('user', () => {
  const userId = ref<string>()
  const initUser = async () => {
    userId.value = (await api.getMe()).userName
  }

  return { userId: readonly(userId), initUser }
})
