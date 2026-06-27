<script setup lang="ts">
import { computed } from 'vue'
import MSIcon from '@/components/MSIcon.vue'

const props = defineProps<{
  id: number
  count: number
  isActive: boolean
}>()

type ReactionInfo = {
  id: number
  iconName: string
  color: string
}

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
  props.isActive ? reactionInfo.value.iconName : `${reactionInfo.value.iconName}-outline`
)
</script>

<template>
  <div
    class="tl-item-reaction"
    :style="{
      color: isActive ? reactionInfo.color : '#aaaaaa',
    }"
  >
    <MSIcon :name="iconNameOutline" class="reaction-icon" :size="24" />
    <span class="tl-item-reaction-count">{{ count }}</span>
  </div>
</template>

<style scoped>
.tl-item-reaction {
  display: flex;
  align-items: center;
  gap: 0.3rem;
}

.tl-item-reaction-count {
  font-weight: bold;
}
</style>
