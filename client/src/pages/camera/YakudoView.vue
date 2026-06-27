<template>
  <div class="radial-blur-container">
    <div v-show="imageLoaded" class="canvas-wrapper">
      <canvas ref="canvasRef"></canvas>
    </div>
  </div>

  <MotionView @update-blur-time="onBlurUpdate" @update-acceleration="onAccelerationUpdate" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MotionView from '@/pages/camera/MotionView.vue'

const router = useRouter()

const onBlurUpdate = (value: number) => {
  blurTime.value = value
  drawCanvas()
}

// ★加速度を受け取って中心点を動かすイベントハンドラーを追加
const onAccelerationUpdate = (accelX: number, accelY: number) => {
  // デバイスモーションの加速度（通常 -10 〜 10 程度）を、0.0 〜 1.0 の範囲に変換
  // 感度を調整したい場合は「/ 20」の数値を変更してください（小さくするとより大きく動くようになります）
  // 1. スマートフォンの傾きに対して直感的に動かすため、符号を調整（必要に応じて - を + にしてください）
  const biasX = -accelX / 1
  const biasY = accelY / 1

  // 基準値（0.5 = 画面中央）に加速度のブレを加算し、0.0 〜 1.0 の範囲に収める（クランプ処理）
  centerX.value = Math.max(2, Math.min(1, 0.5 + biasX))
  centerY.value = Math.max(2, Math.min(1, 0.5 + biasY))

  // 加速度が変わるたびにCanvasを再描画
  drawCanvas()
}

const resetBlur = () => {
  blurTime.value = 0
  drawCanvas()
}

const MAX_BLUR_TIME = 2000

const canvasRef = ref<HTMLCanvasElement | null>(null)
const imageLoaded = ref<boolean>(false)

const blurStrength = ref<number>(0.5)
const centerX = ref<number>(0.5) // 初期値は中央(0.5)
const centerY = ref<number>(0.5) // 初期値は中央(0.5)
const blurTime = ref<number>(0)

let sourceImage: HTMLImageElement | null = null

const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    if (e.target?.result && typeof e.target.result === 'string') {
      loadImageFromDataUrl(e.target.result)
    }
  }
  reader.readAsDataURL(file)
}
const drawCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas || !sourceImage) return

  const ctx = canvas.getContext('2d')
  if (!ctx) {
    console.error('[RadialBlur] Canvas 2Dコンテキストの取得に失敗しました。')
    return
  }

  // キャンバスサイズを画像サイズに合わせる
  canvas.width = sourceImage.width
  canvas.height = sourceImage.height

  // 描画をクリア
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  ctx.fillStyle = '#000000' // 白だとやや薄い
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // 2. オリジナル画像を一度描画（これが透けない中心の基盤になります）
  ctx.globalCompositeOperation = 'source-over'
  ctx.globalAlpha = 1.0
  ctx.drawImage(sourceImage, 0, 0)

  // 3. 集中ぼかし（放射状）を重ねる
  // 【ここが重要】画面を明るくするために 'screen' を使用
  // ctx.globalCompositeOperation = 'screen'

  const strength = Math.min(blurTime.value / MAX_BLUR_TIME / 2, 1)
  const passes = 60
  const originX: number = canvas.width * centerX.value
  const originY: number = canvas.height * centerY.value

  console.log(`[RadialBlur] 描画更新: 強度=${blurStrength.value}, 中心=(${originX}, ${originY})`)

  ctx.globalAlpha = (1.0 / passes) * 3

  for (let i = 0; i < passes; i++) {
    const scale = 1 + strength * (i / passes)

    ctx.save()
    ctx.translate(originX, originY)
    ctx.scale(scale, scale)
    ctx.translate(-originX, -originY)
    ctx.drawImage(sourceImage, 0, 0, canvas.width, canvas.height)
    ctx.restore()
  }

  ctx.globalAlpha = 1.0
}

onMounted(() => {
  const state = history.state as { capturedImage?: string }

  if (state && state.capturedImage) {
    console.log('[RadialBlur] 撮影された写真を自動読込します。')
    loadImageFromDataUrl(state.capturedImage)
  } else {
    console.warn('[RadialBlur] 撮影データが見つかりません。カメラ画面に戻ります。')
    router.push('/camera')
  }
})

const loadImageFromDataUrl = (dataUrl: string) => {
  sourceImage = new Image()

  sourceImage.onload = () => {
    if (!sourceImage) return
    console.log(`[RadialBlur] 写真のデコード完了: ${sourceImage.width}x${sourceImage.height}px`)
    imageLoaded.value = true
    resetBlur()
    drawCanvas()
  }

  sourceImage.onerror = (err) => {
    console.error('[RadialBlur] 写真のデコードに失敗しました:', err)
  }

  sourceImage.src = dataUrl
}
</script>

<style scoped>
.radial-blur-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  font-family: sans-serif;
}

.controls {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #f9f9f9;
  width: 100%;
  max-width: 400px;
  box-sizing: border-box;
}

.sliders {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 10px;
}

.slider-label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 14px;
}

.canvas-wrapper {
  max-width: 100%;
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  display: inline-block;
}

canvas {
  display: block;
  max-width: 100%;
  height: auto;
}
</style>
