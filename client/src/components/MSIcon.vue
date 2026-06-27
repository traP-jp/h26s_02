<script setup lang="ts">
defineProps<{
  name: string
}>()
</script>

<template>
  <span class="material-symbols-outlined">
    {{ name }}
  </span>
</template>

<style scoped>
/* 共通のアイコン用クラスを定義 */
.material-symbols-outlined {
  font-family: 'Material Symbols Outlined Variable', sans-serif;
  font-weight: normal;
  font-style: normal;
  font-size: 24px; /* アイコンの基本サイズ */
  display: inline-block;
  line-height: 1;
  text-transform: none;
  letter-spacing: normal;
  word-wrap: normal;
  white-space: nowrap;
  direction: ltr;
  user-select: none;
  -webkit-user-select: none; /* Safari用 */
  cursor: inherit;

  /* Material Symbols特有のバリエーション設定（お好みで調整可能） */
  font-variation-settings:
    'FILL' 0,
    'wght' 400,
    'GRAD' 0,
    'opsz' 24;
}
</style>
