<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'
import PostImage from '@/components/PostImage.vue'
import { type Post } from '@/schema'

// 親から関数をPropsとして受け取ります
const props = defineProps<{
  loadNew: (existingIds: Set<string>) => Promise<Post[]>
}>()

const el = ref<HTMLElement | null>(null)
const posts = ref<Post[]>([])

const isLoadingTop = ref<boolean>(false)

// 上スクロール → 最新の投稿を取得してマージ
useInfiniteScroll(
  el,
  async () => {
    if (isLoadingTop.value) return
    isLoadingTop.value = true

    try {
      const previousScrollHeight = el.value?.scrollHeight ?? 0
      const existingIds = new Set(posts.value.map((item) => item.id))

      const newPosts = await props.loadNew(existingIds)
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
      console.error('[PostGrid] 最新の投稿の取得・マージエラー:', error)
    } finally {
      // 10 秒まつ
      setTimeout(() => {
        isLoadingTop.value = false
      }, 10000)
    }
  },
  { direction: 'top', distance: 100 }
)
</script>

<template>
  <div ref="el" class="grid">
    <PostImage v-for="post in posts" :key="post.id" :post="post" />
  </div>
</template>

<style scoped>
.grid {
  margin-top: 120px;
  max-width: 800px;
  height: 100dvh;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1px;
  padding: 1px;
}
</style>
