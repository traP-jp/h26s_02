<script setup lang="ts">
import MSIcon from '@/components/MSIcon.vue'
</script>

<template>
  <div class="navigation">
    <div class="buttons">
      <RouterLink to="/timeline" class="button">
        <MSIcon name="home-outline" :size="24" />
      </RouterLink>
      <div class="button">
        <MSIcon name="photo-camera-outline" :size="24" />
      </div>
      <div class="button">
        <MSIcon name="search" :size="24" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.navigation {
  position: fixed;
  z-index: 10;
  bottom: 0;
  width: 100%;
  height: 80px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.buttons {
  display: flex;
  justify-content: space-around;
  align-items: center;
  width: 180px;
  max-width: 400px;
  background-color: #ffffff88;
  backdrop-filter: blur(10px);
  padding: 2px 4px;
  border-radius: 24px;
}

.button {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 60px;
  height: 40px;
  border-radius: 20px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.button:hover {
  background-color: white;
}
</style>
