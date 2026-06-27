<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
const posts = Array.from({ length: 10 }, (_, index) => index + 1)
</script>

<template>
  <div class="timeline">
    <PostImage v-for="num in posts" :key="num" :num="num" />
  </div>
</template>

<style scoped>
.timeline {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1px;
  padding: 1px;
}
</style>
