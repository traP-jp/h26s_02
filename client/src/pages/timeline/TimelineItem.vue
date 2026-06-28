<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
import UserIcon from '@/components/UserIcon.vue'
import TimelineReaction from '@/pages/timeline/TimelineReaction.vue'
import { type Post } from '@/schema'
import { useRouter } from 'vue-router'
import { computed } from 'vue'

const router = useRouter()

const props = defineProps<{
  post: Post
}>()

const reactions = computed(() => {
  if (!props.post.reactions)
    return {
      1: { id: 1, count: 0, myReaction: false },
      2: { id: 2, count: 0, myReaction: false },
      3: { id: 3, count: 0, myReaction: false },
    }

  const r1 = props.post.reactions.find((r) => r.id === 1)
  const r2 = props.post.reactions.find((r) => r.id === 2)
  const r3 = props.post.reactions.find((r) => r.id === 3)

  const result = {
    1: { id: 1, count: r1 ? r1.count : 0, myReaction: false },
    2: { id: 2, count: r2 ? r2.count : 0, myReaction: false },
    3: { id: 3, count: r3 ? r3.count : 0, myReaction: false },
  }

  return result
})
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
        v-for="id in [1, 2, 3]"
        :id="id"
        :key="id"
        :post-id="post.id"
        :count="reactions[id].count"
        :my-reaction="reactions[id].myReaction"
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
