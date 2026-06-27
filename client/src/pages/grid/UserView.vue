<script setup lang="ts">
import PostGrid from '@/pages/grid/PostGrid.vue'
import GridHeader from '@/components/GridHeader.vue'
import BottomNav from '@/components/BottomNav.vue'
import UserIcon from '@/components/UserIcon.vue'

import { useRoute } from 'vue-router'
const route = useRoute()
const userId = route.params['userId'] as string
</script>

<template>
  <PostGrid />
  <GridHeader :posts="12" :reactions="30" :title="userId">
    <UserIcon :user-id="userId" />
  </GridHeader>
  <BottomNav />
</template>
