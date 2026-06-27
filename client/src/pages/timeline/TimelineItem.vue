<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
import UserIcon from '@/components/UserIcon.vue'
import TimelineReaction from '@/pages/timeline/TimelineReaction.vue'

defineProps<{
  num: number
}>()
</script>

<template>
  <div class="tl-item">
    <div class="tl-item-header">
      <UserIcon user-id="kitsne" :size="32" />
      <div class="tl-item-userid">kitsne</div>
    </div>
    <PostImage :num="num" />
    <div class="tl-item-hash">
      <span>#四川</span>
      <span>#油淋鶏</span>
    </div>
    <div class="tl-item-reactions">
      <TimelineReaction v-for="id in 3" :id="id" :key="id" :count="id * 10" :is-active="true" />
    </div>
  </div>
</template>

<style scoped>
.tl-item {
  width: 100%;
  padding: 8px;
}

.tl-item-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 8px;
}

.tl-item-userid {
  font-weight: bold;
}

.tl-item-hash {
  margin-top: 4px;
  color: #ff8400;
}

.tl-item-hash span {
  font-weight: bold;
  margin-right: 0.5rem;
}

.tl-item-reactions {
  display: flex;
  gap: 0.8rem;
  margin-top: 4px;
  color: #808080;
}
</style>
