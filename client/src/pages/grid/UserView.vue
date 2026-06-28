<script setup lang="ts">
import PostGrid from '@/pages/grid/PostGrid.vue'
import GridHeader from '@/components/GridHeader.vue'
import BottomNav from '@/components/BottomNav.vue'
import UserIcon from '@/components/UserIcon.vue'
import { api } from '@/schema'

import { useRoute } from 'vue-router'
const route = useRoute()
const userId = route.params['userId'] as string

const loadNew = async (existingIds: Set<string>) => {
  const latestPosts = await api.getUserPosts(userId)
  const newPosts = latestPosts.filter((post) => !existingIds.has(post.id))
  return newPosts
}

const loadOld = async (lastId?: string) => {
  return await api.getUserPosts(userId, lastId)
}
</script>

<template>
  <PostGrid :load-new="loadNew" :load-old="loadOld" />
  <GridHeader :posts="12" :reactions="30" :title="userId">
    <UserIcon :user-id="userId" />
  </GridHeader>
  <BottomNav />
</template>
