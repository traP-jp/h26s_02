<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useOrientationStore } from '@/store' // パスは環境に合わせて変更してください

const orientationStore = useOrientationStore()

onMounted(() => {
  orientationStore.initListener()
})

onUnmounted(() => {
  orientationStore.removeListener()
})
</script>

<template>
  <div>
    <p v-if="orientationStore.isPortrait">現在は縦画面です</p>
    <p v-else>現在は横画面です</p>
  </div>
</template>
