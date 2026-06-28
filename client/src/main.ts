import { createApp } from 'vue'
import '@/base.css'
import App from '@/App.vue'

import '@fontsource-variable/noto-sans-jp/index.css' // 'Noto Sans JP Variable'
import '@fontsource/m-plus-1p/index.css' // 'M PLUS 1p'
import '@fontsource-variable/reddit-sans/index.css' // 'Reddit Sans Variable'
import '@fontsource-variable/material-symbols-outlined/index.css' // 'Material Symbols Outlined Variable'

// eslint-disable-next-line @typescript-eslint/no-unsafe-argument
const app = createApp(App)

import router from '@/router'
app.use(router)

import { createPinia } from 'pinia'
const pinia = createPinia()
app.use(pinia)

import { useUserStore } from '@/store'
const userStore = useUserStore(pinia)
await userStore.initUser()

console.log('Mode:', import.meta.env.MODE)
app.mount('#app')
