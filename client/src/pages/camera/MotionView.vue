<template>
  <div class="shake-container">
    <div v-if="!isPermissionGranted" class="init-screen">
      <p>センサーを準備しています...</p>
    </div>
    <div v-else class="ready-screen">
      <p>端末を躍動させろ！</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

// ★変更：update-rotation-rate を emit できるように定義を追加
const emit = defineEmits<{
  'update-blur-time': [blurTime: number]
  'update-acceleration': [accelX: number, accelY: number]
  'update-rotation-rate': [rate: number]
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
const blurTime = ref(0)
let lastX = 0,
  lastY = 0,
  lastZ = 0
let lastUpdate = 0

const handleMotion = (event: DeviceMotionEvent) => {
  if (isShaking.value) return

  // 1. 加速度の処理
  const current = event.accelerationIncludingGravity
  if (!current || current.x === null || current.y === null || current.z === null) return

  emit('update-acceleration', current.x, current.y)

  // ★追加：回転速度（ジャイロ）の処理
  // rotationRate.alpha(Z軸), beta(X軸), gamma(Y軸) から、スマホのひねり速度を計算します
  const rotation = event.rotationRate
  if (rotation && rotation.alpha !== null && rotation.beta !== null && rotation.gamma !== null) {
    // 3軸の回転速度の絶対値を合成して「回転全体の激しさ」を算出（単位: 度/秒）
    const rotationSpeed =
      Math.abs(rotation.alpha) + Math.abs(rotation.beta) + Math.abs(rotation.gamma)

    // 親コンポーネントへ通知
    emit('update-rotation-rate', rotationSpeed)
  }

  // 2. 従来のシェイク・ブラー時間の計算
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
</script>

<style scoped>
.shake-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}
</style>
