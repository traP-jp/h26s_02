<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  name: string
  size?: number
}>()

const iconName = computed(() => {
  return `material-symbols:${props.name}`
})
</script>

<template>
  <Icon
    :icon="iconName"
    :width="size ? `${size}px` : '100%'"
    :height="size ? `${size}px` : '100%'"
    class="material-icon-svg"
  />
</template>

<style scoped>
.material-icon-svg {
  color: inherit;
}
</style>
