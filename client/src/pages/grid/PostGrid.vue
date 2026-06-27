<script setup lang="ts">
import PostImage from '@/components/PostImage.vue'
import ScrollView from '@/samples/ScrollView.vue'
import { ref } from 'vue'

const posts = ref<number[]>(Array.from({length: 100}, (_, i) => i + 1))

</script>

<template>
  <ScrollView>
    <div class="grid">
      <PostImage
        v-for="num in posts"
        :key="num"
        :num="num"
      />
    </div>
  </ScrollView>
</template>

<style scoped>
.grid {
  max-width: 800px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1px;
  padding: 1px;
}
</style>
