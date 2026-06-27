<script setup lang="ts">
import MSIcon from '@/components/MSIcon.vue'
import PostImage from '@/components/PostImage.vue'
import UserIcon from '@/components/UserIcon.vue'

defineProps<{
  num: number
}>()
</script>

<template>
  <div class="tl-item">
    <div class="tl-item-header">
      <UserIcon user-id="kitsne" :size="24" />
      <div class="tl-item-userid">kitsne</div>
    </div>
    <PostImage :num="num" />
    <div class="tl-item-hash">
      <span>#四川</span>
      <span>#油淋鶏</span>
    </div>
    <div class="tl-item-reactions">
      <div class="tl-item-reaction">
        <MSIcon name="thumb-up-outline" class="reaction-icon" :size="24" />
        <span class="tl-item-reaction-count">11</span>
      </div>
      <div class="tl-item-reaction">
        <MSIcon name="favorite-outline" class="reaction-icon" :size="24" />
        <span class="tl-item-reaction-count">6</span>
      </div>
      <div class="tl-item-reaction">
        <MSIcon name="local-fire-department-outline" class="reaction-icon" :size="24" />
        <span class="tl-item-reaction-count">7</span>
      </div>
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
  font-size: 0.8rem;
  font-weight: bold;
}

.tl-item-hash {
  margin-top: 4px;
  font-size: 0.8rem;
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
  color: #cccccc;
}

.tl-item-reaction {
  display: flex;
  align-items: center;
  gap: 0.3rem;
}

.tl-item-reaction-count {
  font-weight: bold;
}
</style>
