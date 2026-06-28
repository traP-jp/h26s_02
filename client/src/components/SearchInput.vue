<script setup lang="ts">
import { ref } from 'vue'
import MSIcon from './MSIcon.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const searchQuery = ref('')
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'search', query: string): void
}>()

const executeSearch = () => {
  const query = searchQuery.value.trim()
  if (!query) return

  emit('search', query)
  void router.push('/search?tag=' + encodeURIComponent(query))
}
</script>

<template>
  <div class="overlay" @click.self="emit('close')">
    <div class="search">
      <MSIcon name="search" :size="24" />
      <input
        v-model.lazy="searchQuery"
        class="search-input"
        type="search"
        placeholder="躍動を検索"
        @keyup.enter="executeSearch"
      />
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  z-index: 20;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100dvh;
  background-color: #00000088;
  /* backdrop-filter: blur(10px); */
}

.search {
  position: fixed;
  z-index: 30;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  border-radius: 24px;
  height: 48px;
  width: min(80%, 600px);
  display: flex;
  align-items: center;
  padding: 0 16px;
  gap: 8px;
  background-color: #ffffff88;
  backdrop-filter: blur(10px);
}

.search-input {
  flex-grow: 1;
  height: 100%;
  border: none;
  outline: none;
}

input[type='search']::-webkit-search-cancel-button {
  -webkit-appearance: none;
  appearance: none;
}

/* Edge / Internet Explorer のテキスト入力欄のバツボタンを消す */
input::-ms-clear {
  display: none;
}
</style>
