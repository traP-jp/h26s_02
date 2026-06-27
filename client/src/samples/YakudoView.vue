<template>
  <div class="radial-blur-container">
    <div class="controls">
      <input type="file" accept="image/*" @change="handleImageUpload" />

      <div v-if="imageLoaded" class="sliders">
        <label class="slider-label">
          ぼかしの強さ (ブレ幅):
          <input
            v-model.number="blurStrength"
            type="range"
            min="0"
            max="1"
            step="0.01"
            @input="drawCanvas"
          />
          {{ blurStrength }}
        </label>
        <label class="slider-label">
          中心位置 X:
          <input
            v-model.number="centerX"
            type="range"
            min="0"
            max="1"
            step="0.01"
            @input="drawCanvas"
          />
          {{ Math.round(centerX * 100) }}%
        </label>
        <label class="slider-label">
          中心位置 Y:
          <input
            v-model.number="centerY"
            type="range"
            min="0"
            max="1"
            step="0.01"
            @input="drawCanvas"
          />
          {{ Math.round(centerY * 100) }}%
        </label>
      </div>
    </div>

    <div v-show="imageLoaded" class="canvas-wrapper">
      <canvas ref="canvasRef"></canvas>
    </div>
  </div>

  <MotionView @update-blur-time="onBlurUpdate" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MotionView from './MotionView.vue'

const router = useRouter()

// ==========================================
// 【修正】不足していたカメラからの引き継ぎデータ用 ref 変数を追加
// ==========================================
const initialBlurTime = ref<number>(0)
const initialAcceleration = ref({ x: 0, y: 0, z: 0 })

const onBlurUpdate = (value: number) => {
  blurTime.value = value
  drawCanvas()
}

const resetBlur = () => {
  blurTime.value = 0
  drawCanvas()
}

const MAX_BLUR_TIME = 2000

// テンプレート参照の型定義
const canvasRef = ref<HTMLCanvasElement | null>(null)
const imageLoaded = ref<boolean>(false)

// パラメータの型定義
const blurStrength = ref<number>(0.5)
const centerX = ref<number>(0.5)
const centerY = ref<number>(0.5)
const blurTime = ref<number>(0)

// 画像オブジェクトを保持
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

  ctx.fillStyle = '#000000'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // 2. オリジナル画像を一度描画
  ctx.globalCompositeOperation = 'source-over'
  ctx.globalAlpha = 1.0
  ctx.drawImage(sourceImage, 0, 0)

  // ==========================================
  // 【修正】撮影時の躍動感 ＋ スライダーの強さ ＋ フリフリ躍動感 を全て合算
  // ==========================================
  // スライダーの blurStrength（0〜1）を MAX_BLUR_TIME に掛け合わせて ms 相当に換算します
  const sliderBonusTime = blurStrength.value * MAX_BLUR_TIME
  const totalBlurTime = initialBlurTime.value + blurTime.value + sliderBonusTime

  const strength = Math.min(totalBlurTime / MAX_BLUR_TIME / 2, 1)
  const passes = 60
  const originX: number = canvas.width * centerX.value
  const originY: number = canvas.height * centerY.value

  console.log(`[RadialBlur] 描画更新: 総ぼかし時間=${totalBlurTime}ms, 強度=${strength}`)

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
  const state = history.state as {
    capturedImage?: string
    blurTime?: number
    acceleration?: { x: number; y: number; z: number }
  }

  if (state && state.capturedImage) {
    console.log('[RadialBlur] 撮影された写真とセンサーデータを読み込みます。')
    initialBlurTime.value = state.blurTime || 0
    initialAcceleration.value = state.acceleration || { x: 0, y: 0, z: 0 }
    loadImageFromDataUrl(state.capturedImage)
  } else {
    console.warn('[RadialBlur] 撮影データがないため、テスト用画像を表示します。')
    initialBlurTime.value = 0
    // 直接遷移したとき用のダミー画像（Picsum）
    loadImageFromDataUrl('https://picsum.photos/800/800')
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
