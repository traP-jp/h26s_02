<template>
  <div class="shake-container">
    <div v-if="!isPermissionGranted" class="init-screen">
      <p>センサーを準備しています...</p>
    </div>

    <div class="action-area">
      <button class="dummy-post-btn" @click="handleDummyPost">投稿する</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

// 変更：update-acceleration も親へ送れるように定義を追加
const emit = defineEmits<{
  'update-blur-time': [blurTime: number]
  'update-acceleration': [accelX: number, accelY: number]
}>()

const isPermissionGranted = ref(false)
const isShaking = ref(false)
const debugLogs = ref<string[]>([])

const addLog = (message: string) => {
  debugLogs.value.push(message)
}
const soundUrl = '/sound.mp3'
let audio: HTMLAudioElement | null = null

const SHAKE_THRESHOLD = 1500
const BLUR_THRESHOLD = 3000
// 変更：script setup内での参照用にリアクティブもしくは変数管理
const blurTime = ref(0)
let lastX = 0,
  lastY = 0,
  lastZ = 0
let lastUpdate = 0

const handleMotion = (event: DeviceMotionEvent) => {
  if (isShaking.value) return
  const current = event.accelerationIncludingGravity
  if (!current || current.x === null || current.y === null || current.z === null) return

  // ★追加：親コンポーネント（1階目）にデバイスの傾き（x, y）をリアルタイム伝達
  emit('update-acceleration', current.x, current.y)

  const currentTime = Date.now()
  if (currentTime - lastUpdate > 100) {
    const diffTime = currentTime - lastUpdate
    lastUpdate = currentTime

    const speed =
      (Math.abs(current.x + current.y + current.z - lastX - lastY - lastZ) / diffTime) * 10000

    if (speed > SHAKE_THRESHOLD) {
      addLog(`[Motion] シェイク検知: Speed ${Math.floor(speed)}`)
      playSound()

      isShaking.value = true
      setTimeout(() => {
        isShaking.value = false
      }, 500)
    }

    const isStrongShake = speed > BLUR_THRESHOLD

    if (isStrongShake) {
      addLog(`[Motion] 強いシェイクの検知: Blur Time ${blurTime.value} ms`)
      blurTime.value += diffTime
    }

    blurTime.value = Math.min(blurTime.value, 3000)
    emit('update-blur-time', blurTime.value)

    lastX = current.x
    lastY = current.y
    lastZ = current.z
  }
}

const playSound = () => {
  if (audio) {
    audio.currentTime = 0
    audio.play().catch((e) => {
      if (e instanceof Error) {
        addLog(`[Audio Error] 再生失敗: ${e.message}`)
      } else {
        addLog(`[Audio Error] 再生失敗: ${String(e)}`)
      }
    })
  }
}

onMounted(() => {
  addLog('[Auto] ページが読み込まれたため、自動でセンサーを要求します')
  void requestAccess()
})

const requestAccess = async () => {
  try {
    audio = new Audio(soundUrl)
    audio.play().catch((e) => {
      addLog(`[Audio Init Error] ${String(e)}`)
    })
    audio.pause()
    audio.currentTime = 0

    if (typeof DeviceMotionEvent === 'undefined') {
      addLog('[Error] DeviceMotionEventが未定義です。')
      return
    }

    if (typeof (DeviceMotionEvent as any).requestPermission === 'function') {
      const permissionState = await (DeviceMotionEvent as any).requestPermission()
      if (permissionState === 'granted') {
        window.addEventListener('devicemotion', handleMotion, false)
        isPermissionGranted.value = true
      }
    } else {
      window.addEventListener('devicemotion', handleMotion, false)
      isPermissionGranted.value = true
    }
  } catch (error: unknown) {
    console.error(error)
  }
}

onUnmounted(() => {
  if (isPermissionGranted.value) {
    window.removeEventListener('devicemotion', handleMotion, false)
  }
})

const handleDummyPost = () => {
  console.log('[Dummy Post] ボタンがクリックされました。現在のぼかし時間:', blurTime.value)
  alert(`画像をPOSTしました！（ダミー処理 / blurTime: ${blurTime.value}ms）`)
}
</script>

<style scoped>
.shake-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  /* 変更：親要素のFlexboxの中で自然に縮むように、固定高さを削除 */
  width: 100%;
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
