<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'
import BottomNav from '@/components/BottomNav.vue'
import TimelineItem from '@/pages/timeline/TimelineItem.vue'

import { api, type Post } from '@/schema'

const el = ref<HTMLElement | null>(null)
const posts = ref<Post[]>([])

const isLoadingTop = ref<boolean>(false)
const isLoadingBottom = ref<boolean>(false)

onMounted(async () => {
  await api.getPosts(undefined, 20)
})

// 上スクロール → 最新の投稿を取得してマージ
useInfiniteScroll(
  el,
  async () => {
    if (isLoadingTop.value) return
    isLoadingTop.value = true
    try {
      const previousScrollHeight = el.value?.scrollHeight ?? 0
      const latestPosts = await api.getPosts(undefined, 20)
      const existingIds = new Set(posts.value.map((item) => item.id))
      const newPosts = latestPosts.filter((post) => !existingIds.has(post.id))

      if (newPosts.length > 0) {
        posts.value.unshift(...newPosts) // 先頭に新着投稿を追加

        // スクロール位置を補正
        await nextTick()
        if (el.value) {
          const currentScrollHeight = el.value.scrollHeight
          el.value.scrollTop += currentScrollHeight - previousScrollHeight
        }
      }
    } catch (error) {
      console.error('[Timeline] 最新の投稿の取得・マージエラー:', error)
    } finally {
      // 10 秒まつ
      setTimeout(() => {
        isLoadingTop.value = false
      }, 10000)
    }
  },
  { direction: 'top', distance: 100 }
)

const shouldFetchMore = ref(true) // 追加の投稿を取得するかどうかのフラグ

// 下スクロール → 過去の投稿を取得
useInfiniteScroll(
  el,
  async () => {
    if (isLoadingBottom.value) return
    if (!shouldFetchMore.value) return // 追加の投稿がない場合は処理をスキップ
    isLoadingBottom.value = true
    try {
      const last = posts.value[posts.value.length - 1]
      if (last) {
        const prevPosts = await api.getPosts(last.id, 20)
        if (prevPosts.length > 0) {
          posts.value.push(...prevPosts)
          await nextTick()
        } else {
          shouldFetchMore.value = false // 追加の投稿がない場合はフラグを false にする
        }
      }
    } catch (error) {
      console.error('[Timeline] 過去の投稿の取得エラー:', error)
    } finally {
      isLoadingBottom.value = false
    }
  },
  { direction: 'bottom', distance: 100 }
)
</script>

<template>
  <div ref="el" class="timeline-container">
    <TimelineItem v-for="post in posts" :key="post.id" :post="post" />
  </div>
  <BottomNav />
</template>

<style scoped>
.timeline-container {
  height: 100dvh;
  overflow-y: auto;
}
</style>
