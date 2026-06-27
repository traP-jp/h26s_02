<script setup lang="ts">
import { ref } from 'vue'
import BottomNav from '@/components/BottomNav.vue'
import TimelineItem from './TimelineItem.vue'
import ScrollView from '@/samples/ScrollView.vue';

interface Post {
  id: number
  imageUrl: string
}

const items = ref<Post[]>([
  { id: 1, imageUrl: 'https://picsum.photos/300/400?1' },
  { id: 2, imageUrl: 'https://picsum.photos/300/400?2' },
  { id: 3, imageUrl: 'https://picsum.photos/300/400?3' },
  { id: 4, imageUrl: 'https://picsum.photos/300/400?4' },
  { id: 5, imageUrl: 'https://picsum.photos/300/400?5' },
])

const onReachTop = () => {
  const first = items.value[0]!

  const prev: Post[] = [
    { id: first.id - 5, imageUrl: `https://picsum.photos/300/400?${first.id + 5}` },
    { id: first.id - 4, imageUrl: `https://picsum.photos/300/400?${first.id + 4}` },
    { id: first.id - 3, imageUrl: `https://picsum.photos/300/400?${first.id + 3}` },
    { id: first.id - 2, imageUrl: `https://picsum.photos/300/400?${first.id + 2}` },
    { id: first.id - 1, imageUrl: `https://picsum.photos/300/400?${first.id + 1}` },
  ]

  items.value.unshift(...prev)
}

const onReachBottom = () => {
  const last = items.value[items.value.length - 1]
  if (!last) return

  const next: Post[] = [
    { id: last.id + 1, imageUrl: `https://picsum.photos/300/400?${last.id + 1}` },
    { id: last.id + 2, imageUrl: `https://picsum.photos/300/400?${last.id + 2}` },
    { id: last.id + 3, imageUrl: `https://picsum.photos/300/400?${last.id + 3}` },
    { id: last.id + 4, imageUrl: `https://picsum.photos/300/400?${last.id + 4}` },
    { id: last.id + 5, imageUrl: `https://picsum.photos/300/400?${last.id + 5}` },
  ]

  items.value.push(...next)
}
</script>

<template>
  <div style="width: 500px; margin: 0 auto;">
    <ScrollView
      @reach-top="onReachTop"
      @reach-bottom="onReachBottom"
    >
      <TimelineItem
        v-for="item in items"
        :key="item.id"
        :num="item.id"
      />
    </ScrollView>
  </div>

  <BottomNav />
</template>
