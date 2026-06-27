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
    :width="size ? `${size}px` : undefined"
    :height="size ? `${size}px` : undefined"
  />
</template>

<style scoped>
.user-icon-container {
  border-radius: 50%;
  object-fit: cover;
  /* sizeが指定されていない（属性がない）場合、こちらが適用されます */
  width: 100%;
  height: 100%;
}
</style>
