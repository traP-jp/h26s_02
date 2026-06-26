import { createApp } from 'vue'
import '@/base.css'
import router from '@/router'
import App from '@/App.vue'

import '@fontsource-variable/noto-sans-jp/index.css' // 'Noto Sans JP Variable'

// eslint-disable-next-line @typescript-eslint/no-unsafe-argument
const app = createApp(App)

app.use(router)
console.log('Mode:', import.meta.env.MODE)
app.mount('#app')
