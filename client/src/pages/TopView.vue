<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/schema'

const userName = ref<string | null>(null)
onMounted(async () => {
  const me = await api.getMe()
  userName.value = me.userName
})
</script>

<template>
  <div style="display: flex; flex-direction: column">
    <div v-if="userName">Hello, {{ userName }}!</div>
    <RouterLink to="/samples">Samples</RouterLink>
    <RouterLink to="/timeline">Timeline</RouterLink>
  </div>
</template>
