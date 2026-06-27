<script setup lang="ts">
import MSIcon from '@/components/MSIcon.vue'
</script>

<template>
  <div class="navigation">
    <MSIcon name="home" />
    Hello!
  </div>
</template>

<style scoped>
.navigation {
  position: fixed;
  z-index: 10;
  bottom: 0;
  width: 100%;
  height: 80px;
  background-color: white;
}
</style>
