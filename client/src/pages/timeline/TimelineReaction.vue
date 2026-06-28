<script setup lang="ts">
import { computed, ref } from 'vue'
import MSIcon from '@/components/MSIcon.vue'
import { api } from '@/schema'

const props = defineProps<{
  postId: string
  id: number
  count: number
  myReaction: boolean
}>()

type ReactionInfo = {
  id: number
  iconName: string
  color: string
}

// 私がそのリアクションを押しているかどうか
const isMyReaction = ref(props.myReaction)
const reactionCount = computed(() => {
  const reactionCountByOthers = props.count - (props.myReaction ? 1 : 0)
  return reactionCountByOthers + (isMyReaction.value ? 1 : 0)
})

const reactionInfoList: Record<number, ReactionInfo> = {
  1: { id: 1, iconName: 'thumb-up', color: '#27C200' },
  2: { id: 2, iconName: 'favorite', color: '#FF65C7' },
  3: { id: 3, iconName: 'local-fire-department', color: '#ff8400' },
}

const reactionInfo = computed(() => {
  const info = reactionInfoList[props.id]
  if (!info) throw new Error(`Invalid reaction id: ${props.id}`)
  return info
})

const iconNameOutline = computed(() =>
  isMyReaction.value ? reactionInfo.value.iconName : `${reactionInfo.value.iconName}-outline`
)

const toggleReaction = async () => {
  if (isMyReaction.value) {
    await api.deleteReaction(props.postId, props.id)
    isMyReaction.value = false
  } else {
    await api.postReaction(props.postId, props.id)
    isMyReaction.value = true
  }
}
</script>

<template>
  <div
    class="tl-item-reaction"
    :style="{
      color: isMyReaction ? reactionInfo.color : '#aaaaaa',
    }"
    @click="toggleReaction"
  >
    <MSIcon :name="iconNameOutline" class="reaction-icon" :size="24" />
    <span class="tl-item-reaction-count">{{ reactionCount }}</span>
  </div>
</template>

<style scoped>
.tl-item-reaction {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  padding-right: 0.1rem;
  cursor: pointer;
}

.tl-item-reaction:hover {
  opacity: 0.8;
}

.tl-item-reaction-count {
  font-weight: bold;
}
</style>
