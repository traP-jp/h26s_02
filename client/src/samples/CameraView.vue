<template>
  <div class="camera-layout-wrapper">
    <button
      :disabled="!isCameraActive || isSwitching"
      class="switch-button"
      @click="switchCamera"
    >
      {{
        isSwitching
          ? '⌛️'
          : facingMode === 'user'
          ? '🔄'
          : '🔄'
      }}
    </button>

    <div class="camera-container">
      <h2>カメラ映像のテスト</h2>

      <video
        ref="videoRef"
        width="640"
        height="480"
        autoplay
        playsinline
        class="camera-video"
      ></video>

      <div class="button-area">
        <button 
          :disabled="!isCameraActive || isSwitching"
          :class="['iphone-shutter-button', { 'is-pressing': isPressing }]" 
          @mousedown="startPress"
          @mouseup="endPress"
          @mouseleave="endPress"
          @touchstart.passive="startPress"
          @touchend.passive="endPress"
          @click="captureAndDownload"
        ></button>
      </div>

      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'

const videoRef = ref<HTMLVideoElement | null>(null)
const errorMessage = ref<string>('')
const isCameraActive = ref<boolean>(false)
const isSwitching = ref<boolean>(false)

// 押下状態を管理するフラグを新設
const isPressing = ref<boolean>(false)

let currentStream: MediaStream | null = null
const facingMode = ref<'user' | 'environment'>('user')

// 押し始めたとき
const startPress = () => {
  if (!isCameraActive.value || isSwitching.value) return
  isPressing.value = true
}

// 指・マウスが離れたとき
const endPress = () => {
  isPressing.value = false
}

const startCamera = async () => {
  console.log(`カメラの起動を試みます。向き: ${facingMode.value}`)
  errorMessage.value = ''
  isCameraActive.value = false

  try {
    if (currentStream) {
      console.log('既存のストリームを一度停止させます。')
      stopCamera()
    }

    const constraints = {
      video: { facingMode: facingMode.value },
      audio: false,
    }

    currentStream = await navigator.mediaDevices.getUserMedia(constraints)
    console.log('カメラのストリームを正常に取得しました。')

    if (videoRef.value) {
      videoRef.value.srcObject = currentStream
      videoRef.value.onloadedmetadata = () => {
        isCameraActive.value = true
        console.log('映像のメタデータを読み込み、キャプチャと切り替えを有効化しました。')
      }
    } else {
      console.error('video要素の参照(ref)が取得できず、ストリームを設定できませんでした。')
    }
  } catch (error: unknown) {
    console.error('カメラの取得に失敗しました:', error)
    if (error instanceof Error) {
      if (error.name === 'NotAllowedError') {
        errorMessage.value = 'ユーザーによってカメラへのアクセスが拒否されました。'
      } else if (error.name === 'NotFoundError') {
        errorMessage.value = 'カメラデバイスが見つかりません。'
      } else {
        errorMessage.value = `予期せぬエラーが発生しました: ${error.message || error.name}`
      }
    } else {
      errorMessage.value = '予期せぬエラーが発生しました。'
    }
  }
}

const stopCamera = () => {
  if (currentStream) {
    console.log('カメラのストリームを停止し、リソースを解放します。')
    currentStream.getTracks().forEach((track) => {
      track.stop()
      console.log(`トラック (${track.kind}) を停止しました。`)
    })
    currentStream = null
    isCameraActive.value = false
  }
}

const switchCamera = async () => {
  if (isSwitching.value) return
  console.log('カメラの切り替え処理を開始します。')
  isSwitching.value = true
  facingMode.value = facingMode.value === 'user' ? 'environment' : 'user'
  console.log(`向きを ${facingMode.value} に設定しました。`)
  await startCamera()
  isSwitching.value = false
  console.log('カメラの切り替え処理が完了しました。')
}

const captureAndDownload = () => {
  console.log('キャプチャ処理を開始します。')
  const video = videoRef.value

  if (!video) {
    console.error('video要素が存在しないため、キャプチャを中断しました。')
    return
  }

  try {
    const canvas = document.createElement('canvas')
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight

    const ctx = canvas.getContext('2d')
    if (!ctx) {
      console.error('Canvasの2Dコンテキスト取得に失敗しました。')
      return
    }

    ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
    console.log(`Canvasへの描画が完了しました (サイズ: ${canvas.width}x${canvas.height})。`)

    const dataUrl = canvas.toDataURL('image/png')
    const link = document.createElement('a')
    link.href = dataUrl
    const filename = `capture_${new Date().getTime()}.png`
    link.download = filename

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    console.log(`画像のダウンロード処理を実行しました: ${filename}`)
  } catch (error) {
    console.error('キャプチャまたはダウンロード中にエラーが発生しました:', error)
  }
}

onMounted(startCamera)
onBeforeUnmount(stopCamera)
</script>

<style scoped>
.camera-layout-wrapper {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  width: 100%;
  box-sizing: border-box;
  padding: 2rem;
}

.camera-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  width: 100%;
  max-width: 640px;
}

.camera-video {
  background-color: #000;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  height: auto;
}

.button-area {
  display: flex;
  justify-content: center;
  width: 100%;
  margin-top: 1rem;
}

.switch-button {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: bold;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s, opacity 0.2s;
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  background-color: #555;
  z-index: 10;
}

.switch-button:hover:not(:disabled) {
  background-color: #ffffff;
}

.switch-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
  opacity: 0.8;
}


.iphone-shutter-button {
  position: relative;
  width: 76px;
  height: 76px;
  border-radius: 50%;
  background-color: #ffffff;
  border: 4px solid #000000;
  padding: 0;
  margin: 0;
  cursor: pointer;
  outline: none;
  box-shadow: 0 0 0 2px #ffffff inset;
  

  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  
  transition: transform 0.08s ease-out;
  appearance: none;
  -webkit-appearance: none;
}

/* ボタンの内側の白い円 */
.iphone-shutter-button::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: calc(100% - 8px);
  height: calc(100% - 8px);
  background-color: #ffffff;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  transition: width 0.08s ease-out, height 0.08s ease-out, background-color 0.08s ease-out;
}


.iphone-shutter-button.is-pressing {
  transform: scale(0.9); 
}

.iphone-shutter-button.is-pressing::after {
  width: calc(100% - 24px);
  background-color: #8e8e93; 
}


.iphone-shutter-button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.error-message {
  color: #ff4d4f;
  font-weight: bold;
  padding: 0.5rem;
  border: 1px solid #ff4d4f;
  border-radius: 4px;
  background-color: #fff1f0;
}
</style>