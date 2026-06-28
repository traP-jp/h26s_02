<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useInfiniteScroll } from '@vueuse/core'
import PostImage from '@/components/PostImage.vue'
import { type Post } from '@/schema'

// 親から関数をPropsとして受け取ります
const props = defineProps<{
  loadNew: (existingIds: Set<string>) => Promise<Post[]>
  loadOld: (lastId?: string) => Promise<Post[]>
}>()

const el = ref<HTMLElement | null>(null)
const posts = ref<Post[]>([])

const isLoadingTop = ref<boolean>(false)
const isLoadingBottom = ref<boolean>(false)

// 初回マウント時にデータを取得
onMounted(async () => {
  isLoadingBottom.value = true
  try {
    const initialPosts = await props.loadOld()
    posts.value = initialPosts
  } catch (error) {
    console.error('[PostGrid] 初回データ読み込みエラー:', error)
  } finally {
    isLoadingBottom.value = false
  }
})

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
      isLoadingTop.value = false
    }
  },
  { direction: 'top', distance: 100 }
)

// 下スクロール → 過去の投稿を取得
useInfiniteScroll(
  el,
  async () => {
    // まだ1件も読み込まれていない場合はスキップ
    if (isLoadingBottom.value || posts.value.length === 0) return
    isLoadingBottom.value = true
    try {
      const last = posts.value[posts.value.length - 1]
      if (!last) return
      const prevPosts = await props.loadOld(last.id)
      if (prevPosts.length > 0) {
        posts.value.push(...prevPosts)
        await nextTick()
      }
    } catch (error) {
      console.error('[PostGrid] 過去の投稿の取得エラー:', error)
    } finally {
      isLoadingBottom.value = false
    }
  },
  { direction: 'bottom', distance: 100 }
)
</script>

<template>
  <div ref="el" class="grid">
    <PostImage v-for="post in posts" :key="post.id" :post="post" />
  </div>
</template>

<style scoped>
.grid {
  max-width: 800px;
  height: 100vh;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1px;
  padding: 1px;
}
</style>
