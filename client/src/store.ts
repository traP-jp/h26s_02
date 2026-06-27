import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useOrientationStore = defineStore('orientation', () => {
  const isPortrait = ref(true)
  let cleanup: (() => void) | null = null

  // リスナーを初期化・登録するアクション
  const initListener = () => {
    // 既に登録済みの場合は重複を避ける
    if (cleanup) return

    const mediaQuery = window.matchMedia('(orientation: portrait)')

    // 現在の状態を取得してステートを更新
    isPortrait.value = mediaQuery.matches
    console.log(
      `[OrientationStore] 初期状態をセットしました: ${isPortrait.value ? '縦 (Portrait)' : '横 (Landscape)'}`
    )

    // 画面の向きが変わった際のハンドラ
    const handleOrientationChange = (e: MediaQueryListEvent) => {
      isPortrait.value = e.matches
      console.log(
        `[OrientationStore] 画面の向きが変更されました: ${isPortrait.value ? '縦 (Portrait)' : '横 (Landscape)'}`
      )
    }

    // リスナーの登録
    mediaQuery.addEventListener('change', handleOrientationChange)
    console.log('[OrientationStore] リスナーを登録しました。')

    // 解除用の関数をセット
    cleanup = () => {
      mediaQuery.removeEventListener('change', handleOrientationChange)
      console.log('[OrientationStore] リスナーを解除しました。')
    }
  }

  // リスナーを解除するアクション
  const removeListener = () => {
    if (cleanup) {
      cleanup()
      cleanup = null
    }
  }

  return {
    isPortrait,
    initListener,
    removeListener,
  }
})
