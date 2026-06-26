import { createApp } from 'vue'
import '@/base.css'
import App from '@/App.vue'

import '@fontsource-variable/noto-sans-jp/index.css' // 'Noto Sans JP Variable'

// eslint-disable-next-line @typescript-eslint/no-unsafe-argument
const app = createApp(App)
app.mount('#app')
