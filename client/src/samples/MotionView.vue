<template>
  <div class="shake-container">
    <div v-if="!isPermissionGranted" class="init-screen">
      <p>この機能はデバイスのセンサーを利用します。</p>
      <button class="start-btn" @click="requestAccess">開始する</button>
    </div>

    <div v-else class="play-screen">
      <p>スマホを振ってみてください！</p>
      <div class="device-icon" :class="{ shaking: isShaking }">［スマートフォン］</div>
    </div>

    <div class="debug-panel">
      <p>【デバッグログ】</p>
      <ul class="log-list">
        <li v-for="(log, index) in debugLogs" :key="index">{{ log }}</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'

const emit = defineEmits<{
  'update-blur-time': [blurTime: number]
}>()

const isPermissionGranted = ref(false)
const isShaking = ref(false)
const debugLogs = ref<string[]>([])

// 画面にログを出力するためのヘルパー関数
const addLog = (message: string) => {
  debugLogs.value.push(message)
}
const soundUrl = '/sound.mp3'
let audio: HTMLAudioElement | null = null

const SHAKE_THRESHOLD = 1500
const BLUR_THRESHOLD = 3000
let blurTime = 0
let lastX = 0,
  lastY = 0,
  lastZ = 0
let lastUpdate = 0

const handleMotion = (event: DeviceMotionEvent) => {
  if (isShaking.value) return
  const current = event.accelerationIncludingGravity
  // eslint-disable-next-line @typescript-eslint/prefer-optional-chain
  if (!current || current.x === null || current.y === null || current.z === null) return

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
      addLog(`[Motion] 強いシェイクの検知: Blur Time ${blurTime} ms`)

      blurTime += diffTime
    }

    blurTime = Math.min(blurTime, 3000)

    emit('update-blur-time', blurTime)

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

const requestAccess = async () => {
  addLog('[Init] ボタンが押されました')

  try {
    audio = new Audio(soundUrl)
    audio.play().catch((e) => {
      addLog(`[Audio Init Error] ${String(e)}`)
    })
    audio.pause()
    audio.currentTime = 0
    addLog('[Init] 音声の初期化完了')

    // DeviceMotionEventが存在するか確認
    if (typeof DeviceMotionEvent === 'undefined') {
      addLog(
        '[Error] DeviceMotionEventが未定義です。HTTPS環境またはlocalhostでアクセスしていますか？'
      )
      return
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-member-access
    if (typeof (DeviceMotionEvent as any).requestPermission === 'function') {
      addLog('[Init] iOS向け許可プロンプトを要求します')
      // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
      const permissionState = await (DeviceMotionEvent as any).requestPermission()
      addLog(`[Init] 許可ステータス: ${permissionState}`)

      if (permissionState === 'granted') {
        window.addEventListener('devicemotion', handleMotion, false)
        isPermissionGranted.value = true
        addLog('[Init] センサーを登録しました')
      } else {
        addLog('[Error] センサーへのアクセスが拒否されました')
      }
    } else {
      addLog('[Init] プロンプト不要環境。センサーを登録します')
      window.addEventListener('devicemotion', handleMotion, false)
      isPermissionGranted.value = true
    }
  } catch (error: unknown) {
    if (error instanceof Error) {
      addLog(`[Exception] 例外エラー発生: ${error.message}`)
    } else {
      addLog(`[Exception] 例外エラー発生: ${String(error)}`)
    }
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
  min-height: 50vh;
  text-align: center;
  font-family: sans-serif;
  padding: 20px;
}

.start-btn {
  padding: 12px 24px;
  font-size: 16px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 20px;
}

.device-icon {
  font-size: 24px;
  margin-top: 20px;
  font-weight: bold;
  transition: transform 0.1s ease-in-out;
}

.shaking {
  animation: shake-animation 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
}

@keyframes shake-animation {
  10%,
  90% {
    transform: translate3d(-5px, 0, 0) rotate(-5deg);
  }
  20%,
  80% {
    transform: translate3d(5px, 0, 0) rotate(5deg);
  }
  30%,
  50%,
  70% {
    transform: translate3d(-10px, 0, 0) rotate(-10deg);
  }
  40%,
  60% {
    transform: translate3d(10px, 0, 0) rotate(10deg);
  }
}

.debug-panel {
  margin-top: 40px;
  padding: 10px;
  background-color: #f5f5f5;
  border: 1px solid #ccc;
  width: 100%;
  max-width: 400px;
  text-align: left;
}

.log-list {
  font-size: 12px;
  color: #333;
  padding-left: 20px;
  word-break: break-all;
}
</style>