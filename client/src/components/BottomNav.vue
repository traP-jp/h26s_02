<script setup lang="ts">
import { ref } from 'vue'
import MSIcon from '@/components/MSIcon.vue'
import SearchInput from '@/components/SearchInput.vue'
const isSearching = ref(false)

const nav = navigator as any
const isMobile: boolean = nav.userAgentData.mobile
</script>

<template>
  <div class="navigation">
    <div class="buttons">
      <RouterLink to="/" class="button">
        <MSIcon name="home-outline" :size="24" />
      </RouterLink>
      <RouterLink v-if="isMobile" to="/camera" class="button">
        <MSIcon name="photo-camera-outline" :size="24" />
      </RouterLink>
      <div class="button" @click="isSearching = true">
        <MSIcon name="search" :size="24" />
      </div>
    </div>
  </div>
  <SearchInput v-if="isSearching" @close="isSearching = false" />
</template>

<style scoped>
.navigation {
  position: fixed;
  z-index: 10;
  bottom: 0;
  left: 0;
  width: 100vw;
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
