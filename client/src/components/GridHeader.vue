<script setup lang="ts">
defineProps<{
  posts?: number
  reactions?: number
  title: string
}>()
</script>

<template>
  <div class="header-container">
    <div class="left-icon">
      <slot />
    </div>
    <div>
      <div class="header-title">{{ title }}</div>
      <div class="header-nums">
        <div>
          <span class="info-num">{{ posts ?? 0 }}</span>
          <span class="info-label">投稿</span>
        </div>
        <div>
          <span class="info-num">{{ reactions ?? 0 }}</span>
          <span class="info-label">リアクション</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header-container {
  position: fixed;
  z-index: 10;
  top: 0;
  width: 100%;
  height: 120px;
  background-color: white;
  display: flex;
  align-items: center;
  padding-left: 32px;
}

.left-icon {
  width: 48px;
  height: 48px;
  margin-right: 12px;
}

.header-title {
  font-size: 1.5rem;
  font-weight: bold;
}

.header-nums {
  display: flex;
  gap: 0.7em;
  font-size: 0.9rem;
}

.info-num {
  font-weight: bold;
}

.info-label {
  color: #666666;
  font-weight: 500;
}
</style>
