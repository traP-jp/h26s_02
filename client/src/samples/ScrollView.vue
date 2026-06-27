<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'

// スクロールを検知したい要素の参照
const el = ref<HTMLElement | null>(null)

// 表示するデータ（テスト用として1〜10の初期データを設定）
const data = ref<number[]>([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])

// 読み込み中フラグ（上下それぞれ用意）
const isLoadingTop = ref<boolean>(false)
const isLoadingBottom = ref<boolean>(false)

// 【上方向】の無限スクロール設定
useInfiniteScroll(
  el, // 明示的に要素の ref は握ってそう
  () => {
    console.log('[useInfiniteScroll: Top] スクロールが上部の終端位置に達しました。')
    loadMoreTop()
  },
  { direction: 'top', distance: 10 }
)

// 【下方向】の無限スクロール設定
useInfiniteScroll(
  el,
  () => {
    console.log('[useInfiniteScroll: Bottom] スクロールが下部の終端位置に達しました。')
    loadMoreBottom()
  },
  { direction: 'bottom', distance: 10 }
)

// 上方向のデータ追加処理
const loadMoreTop = (): void => {
  if (isLoadingTop.value) {
    console.log('[loadMoreTop] 現在読み込み中のためスキップします。')
    return
  }
  isLoadingTop.value = true
  console.log('[loadMoreTop] 過去データの追加読み込みを開始します。')

  setTimeout(async () => {
    try {
      // 1. データ追加前のスクロール領域の全体の高さを記憶
      const previousScrollHeight = el.value?.scrollHeight ?? 0

      const firstItem = data.value[0]!
      // 擬似的に古いデータ（マイナス方向）を生成
      const prevData: number[] = [
        firstItem - 5,
        firstItem - 4,
        firstItem - 3,
        firstItem - 2,
        firstItem - 1,
      ]

      // 2. 配列の先頭にデータを追加
      data.value.unshift(...prevData)

      // 3. DOMの描画更新を待機
      await nextTick()

      // 4. 追加後の高さから追加前の高さを引き、その差分だけスクロール位置を加算して視界のズレを補正する
      if (el.value) {
        const currentScrollHeight = el.value.scrollHeight
        el.value.scrollTop += currentScrollHeight - previousScrollHeight
      }

      console.log(`[loadMoreTop] データの追加が完了しました。現在の総件数: ${data.value.length}件`)
    } catch (error) {
      console.error('[loadMoreTop] エラーが発生しました:', error)
    } finally {
      isLoadingTop.value = false
    }
  }, 1000)
}

// 下方向のデータ追加処理
const loadMoreBottom = (): void => {
  if (isLoadingBottom.value) {
    console.log('[loadMoreBottom] 現在読み込み中のためスキップします。')
    return
  }
  isLoadingBottom.value = true
  console.log('[loadMoreBottom] 未来データの追加読み込みを開始します。')

  setTimeout(() => {
    try {
      const lastItem = data.value[data.value.length - 1]!
      // 擬似的に新しいデータを生成
      const nextData: number[] = [
        lastItem + 1,
        lastItem + 2,
        lastItem + 3,
        lastItem + 4,
        lastItem + 5,
      ]

      // 配列の末尾にデータを追加
      data.value.push(...nextData)

      console.log(
        `[loadMoreBottom] データの追加が完了しました。現在の総件数: ${data.value.length}件`
      )
    } catch (error) {
      console.error('[loadMoreBottom] エラーが発生しました:', error)
    } finally {
      isLoadingBottom.value = false
    }
  }, 1000)
}
</script>

<template>
  <div
    style="
      height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
    "
  >
    <div style="height: 20px; text-align: center; color: #666; width: 300px">
      <span v-if="isLoadingTop">↑ 過去のデータを読み込み中...</span>
    </div>

    <div ref="el" style="height: 300px; overflow-y: auto; border: 1px solid #ccc; width: 300px">
      <div v-for="item in data" :key="item" style="padding: 20px; border-bottom: 1px solid #eee">
        アイテム {{ item }}
      </div>
    </div>

    <div style="height: 20px; text-align: center; color: #666; width: 300px">
      <span v-if="isLoadingBottom">↓ 次のデータを読み込み中...</span>
    </div>
  </div>
</template>
