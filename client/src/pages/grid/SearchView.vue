<script setup lang="ts">
import PostGrid from '@/pages/grid/PostGrid.vue'
import GridHeader from '@/components/GridHeader.vue'
import BottomNav from '@/components/BottomNav.vue'
import MSIcon from '@/components/MSIcon.vue'
</script>

<template>
  <PostGrid />
  <GridHeader :posts="12" :reactions="30">
    <MSIcon name="tag" class="hash-tag" />
  </GridHeader>
  <BottomNav />
</template>

<style scoped>
.hash-tag {
  color: #aaaaaa;
}
</style>
