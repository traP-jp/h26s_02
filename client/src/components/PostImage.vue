<script setup lang="ts">
import { ref } from 'vue'
import { type Post } from '@/schema'

defineProps<{
  post: Post
}>()

const isZoomed = ref<boolean>(false)

const toggleZoom = (): void => {
  isZoomed.value = !isZoomed.value
}
</script>

<template>
  <div class="image-wrapper" @click="toggleZoom">
    <img :src="post.imageUrl" class="post-image" loading="lazy" />

    <Teleport to="body">
      <Transition name="fade">
        <div v-if="isZoomed" class="modal-overlay" @click="toggleZoom">
          <img :src="post.imageUrl" class="modal-image" />
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.post-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  cursor: zoom-in;
}

.post-image:hover {
  opacity: 0.9;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  cursor: zoom-out;
}

.modal-image {
  max-width: 90%;
  max-height: 90%;
  object-fit: contain;
  border-radius: 4px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
