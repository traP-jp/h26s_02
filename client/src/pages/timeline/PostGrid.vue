<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
const posts = Array.from({ length: 100 }, (_, index) => index)
</script>

<template>
  <div class="grid">
    <PostImage v-for="num in posts" :key="num" :num="num" />
  </div>
</template>

<style scoped>
.grid {
  max-width: 800px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1px;
  padding: 1px;
}
</style>
