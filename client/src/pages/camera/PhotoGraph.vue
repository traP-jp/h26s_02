<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'

const videoRef = ref<HTMLVideoElement | null>(null)
const errorMessage = ref<string>('')
const isCameraActive = ref<boolean>(false)
const isSwitching = ref<boolean>(false)

const isPressing = ref<boolean>(false)

let currentStream: MediaStream | null = null
const facingMode = ref<'user' | 'environment'>('user')

const startPress = () => {
  if (!isCameraActive.value || isSwitching.value) return
  isPressing.value = true
}

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

<template>
  <div class="camera-layout-wrapper">
    <!-- <div>
        <RouterLink to="/" class="back-button">
          ←
        </RouterLink>
    </div> -->

    <button
      :disabled="!isCameraActive || isSwitching"
      :class="['switch-button', { 'is-switching': isSwitching }]"
      @click="switchCamera"
      aria-label="カメラ切り替え"
    >
      <span class="switch-icon"></span>
    </button>

    <div class="camera-container">
      <!-- <h2>カメラ映像のテスト</h2> -->

      <div class="video-cropper">
        <video ref="videoRef" autoplay playsinline class="camera-video"></video>
      </div>

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

<style scoped>
.camera-layout-wrapper {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  width: 100%;
  box-sizing: border-box;
  padding: 2rem 0;
  background-color: #000000;
}

.camera-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  width: 100%;
}

.camera-container h2 {
  color: #ffffff;
  margin: 0;
}

.camera-video {
  background-color: #000;
  border-radius: 0;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  height: auto;
  /*アスペクト比*/
  aspect-ratio: 1 / 1;
  object-fit: cover;
}

.button-area {
  display: flex;
  justify-content: center;
  width: 100%;
  margin-top: 1rem;
}

.switch-button {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  width: 48px;
  height: 48px;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  z-index: 10;

  background-color: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);

  display: flex;
  align-items: center;
  justify-content: center;

  transition:
    background-color 0.25s ease,
    transform 0.2s ease,
    opacity 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.switch-button:hover:not(:disabled) {
  background-color: rgba(255, 255, 255, 0.25);
  transform: scale(1.05);
}

.switch-button:disabled {
  cursor: not-allowed;
  opacity: 0.4;
  background-color: rgba(255, 255, 255, 0.05);
}

.switch-icon {
  position: relative;
  width: 20px;
  height: 20px;
  border: 2px solid #ffffff;
  border-radius: 50%;
  border-right-color: transparent;
  border-left-color: transparent;
  transform: rotate(-45deg);
  transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.switch-icon::before,
.switch-icon::after {
  content: '';
  position: absolute;
  width: 0;
  height: 0;
  border-style: solid;
}

.switch-icon::before {
  top: -4px;
  right: -2px;
  border-width: 4px 0 4px 5px;
  border-color: transparent transparent transparent #ffffff;
}

.switch-icon::after {
  bottom: -4px;
  left: -2px;
  border-width: 4px 5px 4px 0;
  border-color: transparent #ffffff transparent transparent;
}

.iphone-shutter-button {
  position: relative;
  width: 76px;
  height: 76px;
  border-radius: 50%;
  background-color: #000000;
  border: 4px solid #ffffff;
  padding: 0;
  margin: 0;
  cursor: pointer;
  outline: none;
  box-shadow: 0 0 0 2px #000000 inset;

  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;

  transition: transform 0.08s ease-out;
  appearance: none;
  -webkit-appearance: none;
}

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
  transition:
    width 0.08s ease-out,
    height 0.08s ease-out,
    background-color 0.08s ease-out;
}

.iphone-shutter-button.is-pressing {
  transform: scale(0.9);
}

.iphone-shutter-button.is-pressing::after {
  width: calc(100% - 8px);
  height: calc(100% - 8px);
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

.back-button {
  position: absolute;
  top: 20px;
  left: 20px;

  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
}
</style>
