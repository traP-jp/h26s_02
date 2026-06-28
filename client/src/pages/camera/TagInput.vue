<template>
  <div class="hashtag-input-wrapper" @click="focusInput">
    <span v-for="(tag, index) in modelValue" :key="index" class="hashtag-tag">
      #{{ tag }}
      <button type="button" @click.stop="removeTag(index)" class="remove-btn" aria-label="削除">
        &times;
      </button>
    </span>

    <div class="input-container">
      <span class="prefix">#</span>
      <input
        ref="inputRef"
        v-model="inputValue"
        @keydown="handleKeydown"
        type="text"
        placeholder="タグを入力..."
        class="hashtag-input"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  modelValue: string[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void
}>()

const inputValue = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

// コンポーネント全体をクリックした際に入力欄にフォーカスを当てる
const focusInput = () => {
  inputRef.value?.focus()
}

const handleKeydown = (event: KeyboardEvent) => {
  // 未変換の日本語入力中（IME変換中）は処理をスキップ
  if (event.isComposing) return

  // スペースキーまたはEnterキーでタグを追加
  if (event.key === ' ' || event.key === 'Enter') {
    event.preventDefault() // スペース文字そのものが入力されるのを防ぐ
    addTag()
  }

  // Backspaceキーで直前のタグを削除（入力欄が空の時のみ）
  if (event.key === 'Backspace' && inputValue.value === '') {
    if (props.modelValue.length > 0) {
      const newTags = [...props.modelValue]
      newTags.pop()
      emit('update:modelValue', newTags)
    }
  }
}

const addTag = () => {
  // 先頭に誤って#を入力した場合は除去し、前後の空白をトリム
  const trimmed = inputValue.value.trim().replace(/^#/, '')

  // 空文字でなく、かつ重複していない場合のみ追加
  if (trimmed && !props.modelValue.includes(trimmed)) {
    emit('update:modelValue', [...props.modelValue, trimmed])
  }

  // 入力欄をリセット
  inputValue.value = ''
}

const removeTag = (index: number) => {
  const newTags = [...props.modelValue]
  newTags.splice(index, 1)
  emit('update:modelValue', newTags)
}
</script>

<style scoped>
.hashtag-input-wrapper {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  border: 1px solid #d1d5db;
  padding: 8px;
  border-radius: 6px;
  background-color: #ffffff;
  cursor: text;
  min-height: 42px;
}

.hashtag-input-wrapper:focus-within {
  border-color: #3b82f6;
  outline: 2px solid #93c5fd;
}

.hashtag-tag {
  background-color: #e0f2fe;
  color: #0284c7;
  padding: 4px 10px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 500;
}

.remove-btn {
  background: none;
  border: none;
  color: #0284c7;
  margin-left: 6px;
  padding: 0;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
}

.remove-btn:hover {
  color: #0369a1;
}

.input-container {
  display: flex;
  align-items: center;
  flex-grow: 1;
  min-width: 120px;
}

.prefix {
  color: #6b7280;
  margin-right: 2px;
  font-weight: bold;
}

.hashtag-input {
  border: none;
  outline: none;
  flex-grow: 1;
  font-size: 14px;
  background: transparent;
}
</style>
