<script setup lang="ts">
import BottomNav from '@/components/BottomNav.vue'
import TimelineItem from './TimelineItem.vue'
</script>

<template>
  <TimelineItem v-for="num in 100" :key="num" :num="num" />
  <BottomNav />
</template>
