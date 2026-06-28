<template>
  <div class="hashtag-input-wrapper" @click="focusInput">
    <span v-for="(tag, index) in modelValue" :key="index" class="hashtag-tag">
      #{{ tag }}
      <button type="button" class="remove-btn" @click.stop="removeTag(index)">&times;</button>
    </span>

    <div class="input-prefix-wrapper">
      <span class="prefix-hash">#</span>
      <input
        ref="inputRef"
        v-model="currentInput"
        type="text"
        class="hashtag-inner-input"
        :placeholder="modelValue.length === 0 ? 'タグを入力...' : ''"
        @keydown="handleKeydown"
        @blur="addTag"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'

const props = defineProps<{
  modelValue: string[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void
}>()

const currentInput = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

const focusInput = () => {
  inputRef.value?.focus()
}

const addTag = () => {
  // 全角・半角スペース、および先頭の#をトリミング
  let cleaned = currentInput.value.replace(/^#/, '').trim()

  if (cleaned) {
    // 重複チェック
    if (!props.modelValue.includes(cleaned)) {
      const updated = [...props.modelValue, cleaned]
      console.log('[HashtagInputSimple] Added tag:', cleaned, updated)
      emit('update:modelValue', updated)
    }
    currentInput.value = ''
  }
}

const removeTag = (index: number) => {
  const updated = props.modelValue.filter((_, i) => i !== index)
  console.log('[HashtagInputSimple] Removed tag at index:', index, updated)
  emit('update:modelValue', updated)
}

const handleKeydown = (event: KeyboardEvent) => {
  // スペースキーが押されたらタグを確定
  if (event.key === ' ' || event.key === 'Spacebar') {
    event.preventDefault()
    addTag()
  }

  // 入力が空の状態でバックスペースが押されたら、最後のタグを削除
  if (event.key === 'Backspace' && currentInput.value === '' && props.modelValue.length > 0) {
    event.preventDefault()
    removeTag(props.modelValue.length - 1)
  }
}
</script>

<style scoped>
.hashtag-input-wrapper {
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 6px 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
  background: #fff;
  cursor: text;
}

.hashtag-input-wrapper:focus-within {
  border-color: #000000;
  box-shadow: 0 0 0 1px #000000;
}

.hashtag-tag {
  background: #f3f4f6;
  color: #1f2937;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.remove-btn {
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 0;
}

.remove-btn:hover {
  color: #ef4444;
}

.input-prefix-wrapper {
  display: flex;
  align-items: center;
  position: relative;
  flex-grow: 1;
}

.prefix-hash {
  color: #9ca3af;
  font-size: 14px;
  user-select: none;
  margin-right: 2px;
}

.hashtag-inner-input {
  border: none;
  outline: none;
  padding: 0;
  margin: 0;
  font-size: 14px;
  width: 100%;
  min-width: 60px;
  flex-grow: 1;
}
</style>
