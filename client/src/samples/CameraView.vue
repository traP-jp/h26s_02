<template>
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
        class="switch-button"
        @click="switchCamera"
      >
        {{
          isSwitching
            ? '切り替え中...'
            : facingMode === 'user'
              ? 'アウトカメラへ切り替え'
              : 'インカメラへ切り替え'
        }}
      </button>

      <button :disabled="!isCameraActive" class="capture-button" @click="captureAndDownload">
        キャプチャしてダウンロード
      </button>
    </div>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'

const videoRef = ref<HTMLVideoElement | null>(null)
const errorMessage = ref<string>('')
const isCameraActive = ref<boolean>(false)
const isSwitching = ref<boolean>(false) // 切り替え中の状態を管理

let currentStream: MediaStream | null = null

// 現在のカメラの向きを管理する状態 (初期値はインカメラ)
const facingMode = ref<'user' | 'environment'>('user')

const startCamera = async () => {
  console.log(`カメラの起動を試みます。向き: ${facingMode.value}`)
  errorMessage.value = ''
  isCameraActive.value = false

  try {
    // 既存のストリームがあれば確実に停止させる（切り替え時）
    if (currentStream) {
      console.log('既存のストリームを一度停止させます。')
      stopCamera()
    }

    const constraints = {
      video: {
        facingMode: facingMode.value, // リアクティブな状態を使用
      },
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

// カメラを切り替える関数
const switchCamera = async () => {
  if (isSwitching.value) return // 切り替え中は処理をスキップ

  console.log('カメラの切り替え処理を開始します。')
  isSwitching.value = true // 切り替え中フラグを立てる

  // 1. 状態をトグルする (user <-> environment)
  facingMode.value = facingMode.value === 'user' ? 'environment' : 'user'
  console.log(`向きを ${facingMode.value} に設定しました。`)

  // 2. 新しい向きでカメラを起動 (startCamera内で既存ストリームの停止も行う)
  await startCamera()

  isSwitching.value = false // 切り替え中フラグを下ろす
  console.log('カメラの切り替え処理が完了しました。')
}

// キャプチャしてダウンロードする関数 (前回と同じ)
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
.camera-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.camera-video {
  background-color: #000;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.button-area {
  display: flex;
  gap: 1rem;
}

.switch-button,
.capture-button {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: bold;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition:
    background-color 0.2s,
    opacity 0.2s;
}

.switch-button {
  background-color: #f5222d; /* 切り替えボタンは赤系 */
}

.switch-button:hover:not(:disabled) {
  background-color: #ff4d4f;
}

.capture-button {
  background-color: #1890ff; /* キャプチャボタンは青系 */
}

.capture-button:hover:not(:disabled) {
  background-color: #40a9ff;
}

.switch-button:disabled,
.capture-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
  opacity: 0.8;
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
