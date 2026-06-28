<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'
import BottomNav from '@/components/BottomNav.vue'
import TimelineItem from './TimelineItem.vue'

interface Post {
  id: number
}

const el = ref<HTMLElement | null>(null)

const items = ref<Post[]>([{ id: 101 }, { id: 102 }, { id: 103 }, { id: 104 }, { id: 105 }])

const isLoadingTop = ref<boolean>(false)
const isLoadingBottom = ref<boolean>(false)

useInfiniteScroll(
  el,
  () => {
    if (isLoadingTop.value) return
    isLoadingTop.value = true

    // 読み込みに 1 秒かかる想定
    setTimeout(async () => {
      const previousScrollHeight = el.value?.scrollHeight ?? 0

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

      // スクロール位置を修正
      await nextTick()
      if (el.value) {
        const currentScrollHeight = el.value.scrollHeight
        el.value.scrollTop += currentScrollHeight - previousScrollHeight
      }
      isLoadingTop.value = false
    }, 1000)
  },
  { direction: 'top', distance: 100 }
)

useInfiniteScroll(
  el,
  () => {
    if (isLoadingBottom.value) return
    isLoadingBottom.value = true

    // 読み込みに 1 秒かかる想定
    setTimeout(async () => {
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
      isLoadingBottom.value = false
    }, 1000)
  },
  { direction: 'bottom', distance: 100 }
)
</script>

<template>
  <div ref="el">
    <TimelineItem v-for="item in items" :key="item.id" :num="item.id" />
  </div>
  <BottomNav />
</template>
