<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
import UserIcon from '@/components/UserIcon.vue'
import TimelineReaction from '@/pages/timeline/TimelineReaction.vue'
import { type Post, type Reaction } from '@/schema'
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'

const router = useRouter()

const props = defineProps<{
  post: Post
}>()

const reactions = ref<Record<number, Reaction> | null>(null)

onMounted(() => {
  reactions.value = props.post.reactions
  for (let i = 1; i <= 3; i++) {
    reactions.value[i] ??= { id: i, count: 0, myReaction: false }
  }
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
      <RouterLink v-for="tag in post.tags" :key="tag" :to="`/search?tag=${tag}`">
        #{{ tag }}
      </RouterLink>
    </div>
    <div v-if="reactions" class="tl-item-reactions">
      <TimelineReaction
        v-for="id in [1, 2, 3]"
        :id="id"
        :key="id"
        :post-id="post.id"
        :count="reactions[id]!.count"
        :my-reaction="reactions[id]!.myReaction"
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
