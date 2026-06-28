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

const handlePostSubmit = async (x: number) => {
  try {
    const response = await fetch(`/images/image ${x}.png`)
    const blob = await response.blob()
    const imageFile = new File([blob], `image ${x}.png`, { type: 'image/png' })
    const tags = ['image', `${x}`]
    const result = await api.newPost({
      image: imageFile,
      tags: tags,
    })

    console.log('[handlePostSubmit] 投稿処理が正常に完了しました。', result)
  } catch (error) {
    console.error('[handlePostSubmit] 処理中にエラーが発生しました。', error)
  }
}

// 仮の画像を 30 個くらい入れる
import { api } from '@/schema'
if ((await api.getPosts()).length === 0) {
  for (let i = 0; i < 10; i++) {
    await handlePostSubmit(i + 1)
  }
}

console.log('Mode:', import.meta.env.MODE)
app.mount('#app')
