<script setup lang="ts">
defineProps<{
  num: number
}>()
</script>

<template>
  <div class="image-wrapper">
    <img :src="`/images/image ${num}.png`" class="post-image" />
  </div>
</template>

<style module>
.image-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  background-color: #f0f0f0;
}

.post-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
</style>
