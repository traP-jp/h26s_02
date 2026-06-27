<script setup lang="ts">
defineProps<{
  num: number
}>()
</script>

<template>
  <div class="image-wrapper">
    <img :src="`/images/image ${(num % 10) + 1}.png`" class="post-image" loading="lazy" />
  </div>
</template>

<style module>
.image-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  background-color: #f0f0f0; /* 画像が存在しない場合の背景色 */
}

.post-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
</style>
