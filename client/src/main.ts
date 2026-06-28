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
    console.log('[handlePostSubmit] 開始: 静的画像の取得を試みます。')

    // 1. public 内の画像を fetch で取得する（public の記述は不要で、/ から始めます）
    // スペースが含まれるファイル名の場合は、URLエンコード（%20）されるか、そのまま文字列で指定します
    const response = await fetch(`/images/image ${x}.png`)

    if (!response.ok) {
      throw new Error(`画像の取得に失敗しました: ${response.status} ${response.statusText}`)
    }

    // 2. レスポンスを Blob（バイナリデータ）に変換する
    const blob = await response.blob()

    // 3. Blob から api.newPost が受け取れる File オブジェクトを作成する
    // 第1引数: データの配列, 第2引数: ファイル名, 第3引数: オプション（MIMEタイプなど）
    const imageFile = new File([blob], 'image 1.png', { type: 'image/png' })

    // 4. 送信するタグの配列を定義する
    const tags = ['Vue', 'TypeScript']

    console.log('[handlePostSubmit] Fileオブジェクトの作成完了。newPost を呼び出します。', {
      fileName: imageFile.name,
      fileSize: imageFile.size,
      fileType: imageFile.type,
      tags,
    })

    // 5. 改修した newPost を呼び出す
    // ※ api オブジェクトの定義場所に合わせて適宜インポートや参照を行ってください
    const result = await api.newPost({
      image: imageFile,
      tags: tags,
    })

    console.log('[handlePostSubmit] 投稿処理が正常に完了しました。', result)
    // 必要に応じて、成功後の画面遷移や通知処理をここに記述します
  } catch (error) {
    console.error('[handlePostSubmit] 処理中にエラーが発生しました。', error)
    // 必要に応じて、エラーメッセージを画面に表示するなどの処理を記述します
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
