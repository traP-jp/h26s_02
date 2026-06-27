<script setup lang="ts">
import PostGrid from '@/pages/timeline/PostGrid.vue'
import GridHeader from '@/components/GridHeader.vue'
import BottomNav from '@/components/BottomNav.vue'
import UserIcon from '@/components/UserIcon.vue'
</script>

<template>
  <PostGrid />
  <GridHeader :posts="12" :reactions="30">
    <UserIcon user-id="kitsne" />
  </GridHeader>
  <BottomNav />
</template>
