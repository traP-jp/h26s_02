<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'
import BottomNav from '@/components/BottomNav.vue'
import TimelineItem from './TimelineItem.vue'

interface Post {
  id: number
}

const el = ref<HTMLElement | null>(null)

const items = ref<Post[]>([
  { id: 1 },
  { id: 2 },
  { id: 3 },
  { id: 4 },
  { id: 5 },
])

const isLoadingTop = ref<boolean>(false)
const isLoadingBottom = ref<boolean>(false)

useInfiniteScroll(
  el,
  () => {
    if (isLoadingTop.value) return
    console.log('[useInfiniteScroll: Top] スクロールが上部の終端位置に達しました。')
    loadMoreTop()
  },
  { direction: 'top', distance: 100 }
)

useInfiniteScroll(
  el,
  () => {
    if (isLoadingBottom.value) return
    console.log('[useInfiniteScroll: Bottom] スクロールが下部の終端位置に達しました。')
    loadMoreBottom()
  },
  { direction: 'bottom', distance: 100 }
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

      // 2. 配列の先頭にデータを追加
      const first = items.value[0]
      if (first) {
        const prev: Post[] = [
          { id: first.id - 5 },
          { id: first.id - 4 },
          { id: first.id - 3 },
          { id: first.id - 2 },
          { id: first.id - 1 },
        ]
        items.value.unshift(...prev)
      }
      // 3. DOMの描画更新を待機
      await nextTick()

      // 4. 追加後の高さから追加前の高さを引き、その差分だけスクロール位置を加算して視界のズレを補正する
      if (el.value) {
        const currentScrollHeight = el.value.scrollHeight
        el.value.scrollTop += currentScrollHeight - previousScrollHeight
      }

      console.log(`[loadMoreTop] データの追加が完了しました。`)
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

  setTimeout(async () => {
    try {
      // 配列の末尾にデータを追加
      const last = items.value[items.value.length - 1]
      if (last) {
        const next: Post[] = [
          { id: last.id + 1 },
          { id: last.id + 2 },
          { id: last.id + 3 },
          { id: last.id + 4 },
          { id: last.id + 5 },
        ]
        items.value.push(...next)
      }
      await nextTick()

      console.log(
        `[loadMoreBottom] データの追加が完了しました。`
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
    <div style="height: 20px; text-align: center; color: #666; width: 500px">
      <span v-if="isLoadingTop">↑ 過去のデータを読み込み中...</span>
    </div>

    <div ref="el" style="height: 800px; overflow-y: auto; border: 1px solid #ccc; width: 500px">
      <TimelineItem
        v-for="item in items"
        :key="item.id"
        :num="item.id"
      />
    </div>

    <div style="height: 20px; text-align: center; color: #666; width: 300px">
      <span v-if="isLoadingBottom">↓ 次のデータを読み込み中...</span>
    </div>
  </div>

  <BottomNav />
</template>
