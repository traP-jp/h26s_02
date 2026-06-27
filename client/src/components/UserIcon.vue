<script setup lang="ts">
defineProps<{
  userId: string
  size?: number
}>()
</script>

<template>
  <img
    class="user-icon-container"
    :src="`https://q.trap.jp/api/v3/public/icon/${userId}`"
    loading="lazy"
    :width="size"
    :height="size"
    :style="{
      width: size ? `${size}px` : '100%',
      height: size ? `${size}px` : '100%',
    }"
  />
</template>

<style scoped>
.user-icon-container {
  border-radius: 50%;
  object-fit: cover;
}
</style>
