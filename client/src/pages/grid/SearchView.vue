<script setup lang="ts">
import PostGrid from '@/pages/grid/PostGrid.vue'
import GridHeader from '@/components/GridHeader.vue'
import BottomNav from '@/components/BottomNav.vue'
import MSIcon from '@/components/MSIcon.vue'
import { api } from '@/schema'

import { useRoute } from 'vue-router'
const route = useRoute()
const tag = route.query['tag'] as string

const loadNew = async (existingIds: Set<string>) => {
  const latestPosts = await api.getTagPosts(tag)
  const newPosts = latestPosts.filter((post) => !existingIds.has(post.id))
  return newPosts
}
</script>

<template>
  <PostGrid :load-new="loadNew" />
  <GridHeader :title="tag">
    <MSIcon name="tag" class="hash-tag" />
  </GridHeader>
  <BottomNav />
</template>

<style scoped>
.hash-tag {
  color: #aaaaaa;
}
</style>
