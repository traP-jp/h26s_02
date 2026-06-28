<template>
  <div class="main-layout">
    <div class="background-text">躍動</div>

    <div class="radial-blur-container">
      <div v-show="imageLoaded" class="canvas-wrapper">
        <canvas ref="canvasRef"></canvas>
      </div>
    </div>

    <MotionView
      @update-blur-time="onBlurUpdate"
      @update-acceleration="onAccelerationUpdate"
      @update-rotation-rate="onUpdateRotationRate"
    />

    <div class="action-area">
      <button class="dummy-post-btn" @click="handlePost">投稿する</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MotionView from '@/pages/camera/MotionView.vue'
import { api } from '@/schema'

const router = useRouter()

const onBlurUpdate = (value: number) => {
  blurTime.value = value
  drawCanvas()
}

const resetBlur = () => {
  blurTime.value = 0
  drawCanvas()
}

// ここに「最大値」を保持します
const currentRotationRate = ref(0)

// ★変更：子コンポーネントから値が送られてくるたびにチェック
const onUpdateRotationRate = (rate: number) => {
  // 送られてきた値が、これまでの最大値よりも大きければ更新
  if (rate > currentRotationRate.value) {
    currentRotationRate.value = rate
    // 最大値が更新されたときだけキャンバスを再描画する
    drawCanvas()
  }
}

const MAX_BLUR_TIME = 1000
const canvasRef = ref<HTMLCanvasElement | null>(null)
const imageLoaded = ref<boolean>(false)
const centerX = ref<number>(0.5)
const centerY = ref<number>(0.5)
const blurTime = ref<number>(0)

const onAccelerationUpdate = (accelX: number, accelY: number) => {
  const biasX = -accelX / 10
  const biasY = accelY / 10
  centerX.value = Math.max(0, Math.min(1, 0.5 + biasX))
  centerY.value = Math.max(0, Math.min(1, 0.5 + biasY))
  drawCanvas()
}

let sourceImage: HTMLImageElement | null = null

const drawCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas || !sourceImage) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  canvas.width = sourceImage.width
  canvas.height = sourceImage.height

  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.drawImage(sourceImage, 0, 0)

  const strength = Math.min(blurTime.value / MAX_BLUR_TIME / 2, 3)

  // ★ 常にその時点での「最大値」が適用されます
  const maxRotation = Math.min(currentRotationRate.value * 0.0002, 0.4)

  const passes = 60
  const originX: number = canvas.width * centerX.value
  const originY: number = canvas.height * centerY.value

  ctx.globalAlpha = (1.0 / passes) * 4.0

  for (let i = 0; i < passes; i++) {
    const ratio = i / passes
    const scale = 1 + strength * ratio
    const angle = maxRotation * (ratio * ratio) * 0.2

    ctx.save()
    ctx.translate(originX, originY)
    ctx.rotate(angle)
    ctx.scale(scale, scale)
    ctx.translate(-originX, -originY)
    ctx.drawImage(sourceImage, 0, 0, canvas.width, canvas.height)
    ctx.restore()
  }
  ctx.globalAlpha = 1.0
}

onMounted(() => {
  // ここでの sessionStorage の読み込みは不要になったので削除しました
  const state = history.state as { capturedImage?: string }
  if (state && state.capturedImage) {
    loadImageFromDataUrl(state.capturedImage)
  } else {
    router.push('/camera')
  }
})

const loadImageFromDataUrl = (dataUrl: string) => {
  sourceImage = new Image()
  sourceImage.onload = () => {
    if (!sourceImage) return
    imageLoaded.value = true
    resetBlur()
    drawCanvas()
  }
  sourceImage.src = dataUrl
}

const handlePost = async () => {
  if (!canvasRef.value) return
  const blob = await new Promise<Blob | null>((resolve) => {
    canvasRef.value!.toBlob(resolve, 'image/png')
  })
  if (!blob) throw new Error('Canvas から Blob の生成に失敗しました。')

  const imageFile = new File([blob], 'image.png', { type: 'image/png' })
  const tags = ['sample', `tag`]

  await api.newPost({ image: imageFile, tags: tags })
  await router.push('/')
}
</script>

<style scoped>
/* 画面全体を上下中央揃えにするスタイル */
.main-layout {
  position: relative; /* 背景文字の基準点 */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
  box-sizing: border-box;
  padding: 24px;
  gap: 24px;
  overflow: hidden;
  background-color: #f7fafc;
}

/* ★ 画面下部に配置し、薄くした「躍動」の文字 */
.background-text {
  position: absolute;
  bottom: -2vh; /* 文字が大きくなった分、下からはみ出させるためにマイナスに調整 */
  left: 50%;
  transform: translateX(-50%) rotate(-8deg); /* 大きな文字が綺麗に見えるよう、傾きを少し緩やか（-8度）に */

  font-size: 45vw; /* ★ 30vw から 45vw に大幅アップ！ */
  font-weight: 900;
  font-family: 'Impact', 'Arial Black', 'Noto Sans JP', sans-serif;

  color: #000000;
  opacity: 0.08; /* 薄さはキープ（お好みで 0.1 くらいに上げてもカッコいいです） */

  white-space: nowrap;
  letter-spacing: -0.05em; /* 文字が離れすぎないように少し詰め気味に */
  pointer-events: none;
  z-index: 0;
  user-select: none;
}

/* 画像（Canvas）エリア */
.radial-blur-container {
  position: relative;
  z-index: 2; /* 背景文字（0）より「前」に出す */
  display: flex;
  justify-content: center;
  align-items: center;
  max-width: 100%;
  max-height: 65vh;
}

.canvas-wrapper {
  max-width: 100%;
  max-height: 100%;
  overflow: hidden;
  border-radius: 12px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.2); /* 影を少し強めに */
}

canvas {
  display: block;
  max-width: 100%;
  max-height: 65vh;
  object-fit: contain;
}

/* 2階目のMotionView（ボタン等）エリア */
:deep(.shake-container) {
  position: relative;
  z-index: 2; /* 背景文字（0）より「前」に出す */
}

.action-area {
  display: flex;
  justify-content: center;
  width: 100%;
}

.dummy-post-btn {
  padding: 14px 40px;
  font-size: 16px;
  font-weight: bold;
  color: #ffffff;
  background-color: #3182ce;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  box-shadow: 0 4px 6px rgba(49, 130, 206, 0.3);
  transition: all 0.2s ease;
}

.dummy-post-btn:hover {
  background-color: #2b6cb0;
  transform: translateY(-1px);
}

.dummy-post-btn:active {
  transform: translateY(1px);
}
</style>
