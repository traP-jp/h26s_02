<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
import UserIcon from '@/components/UserIcon.vue'
import TimelineReaction from '@/pages/timeline/TimelineReaction.vue'
import { type Post } from '@/schema'
import { useRouter } from 'vue-router'

const router = useRouter()

defineProps<{
  post: Post
}>()
</script>

<template>
  <div class="tl-item">
    <div class="tl-item-header" @click="router.push(`/user/${post.userName}`)">
      <UserIcon :user-id="post.userName" :size="32" />
      <div class="tl-item-userid">{{ post.userName }}</div>
    </div>
    <PostImage :post="post" />
    <div class="tl-item-hash">
      <span v-for="tag in post.tags" :key="tag">#{{ tag }}</span>
    </div>
    <div class="tl-item-reactions">
      <TimelineReaction
        v-for="reaction in post.reactions"
        :id="reaction.id"
        :key="reaction.id"
        :count="reaction.count"
        :is-active="true"
      />
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
  cursor: pointer;
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
