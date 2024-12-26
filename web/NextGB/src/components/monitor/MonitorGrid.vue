<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount, onMounted } from 'vue'
import { VideoCamera, Close, Camera, FullScreen, Refresh, Setting, Mute, Microphone } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Device, ChannelInfo } from '@/types/api'
import { deviceApi } from '@/api'
import { SrsRtcPlayer } from '@/api/srs'

interface DeviceWithChannel extends Device {
  channelInfo?: ChannelInfo
  player?: any
  error?: boolean
  id?: string
  channel?: ChannelInfo
  isMuted?: boolean
}

interface StreamResponse {
  url: string
  [key: string]: any
}

// 布局配置
const layouts = {
  '1': { cols: 1, rows: 1, size: 1, label: '单屏' },
  '4': { cols: 2, rows: 2, size: 4, label: '四分屏' },
  '9': { cols: 3, rows: 3, size: 9, label: '九分屏' },
  '16': { cols: 4, rows: 4, size: 16, label: '十六分屏' },
}

// 状态管理
const selectedDevices = ref<DeviceWithChannel[]>([])
const currentLayout = ref<keyof typeof layouts>('9')
const isFullscreen = ref(false)
const showSettings = ref(false)
const autoReconnectEnabled = ref(true)
const reconnectInterval = ref(10)
let reconnectTimers: number[] = []

// 计算属性
const gridStyle = computed(() => {
  const layout = layouts[currentLayout.value]
  return {
    gridTemplateColumns: `repeat(${layout.cols}, 1fr)`,
    gridTemplateRows: `repeat(${layout.rows}, 1fr)`,
  }
})

const maxDevices = computed(() => layouts[currentLayout.value].size)

// 视频流控制
const startWebRTCPlay = async (url: string, index: number, device: DeviceWithChannel) => {
  const player = SrsRtcPlayer()
  device.player = player

  // @ts-ignore
  player.onaddstream = (event: { stream: MediaStream }) => {
    const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
    if (videoElement) {
      videoElement.srcObject = event.stream
    }
  }

  await player.play(url)
}

const startStream = async (device: DeviceWithChannel, index: number) => {
  try {
    device.error = false

    const response = await deviceApi.inviteStream({
      device_id: device.channel!.parent_id,
      channel_id: device.channel!.device_id,
      sub_stream: '0',
    })

    const streamData = response.data as unknown as StreamResponse
    if (!streamData?.url) {
      throw new Error('播放地址不存在')
    }

    await startWebRTCPlay(streamData.url, index, device)
  } catch (error) {
    console.error('启动播放失败:', error)
    device.error = true
    ElMessage.error('启动播放失败')
  }
}

const retryStream = async (index: number) => {
  const device = selectedDevices.value[index]
  if (device) {
    await startStream(device, index)
  }
}

// 设备管理
const cleanupDevice = async (device: DeviceWithChannel) => {
  if (device.player) {
    try {
      await device.player.close()
      device.player = null
    } catch (err) {
      console.error('关闭播放器失败:', err)
    }
  }
}

const addDevice = async (device: Device & { channel: ChannelInfo }) => {
  if (selectedDevices.value.length < maxDevices.value) {
    try {
      const deviceWithChannel: DeviceWithChannel = {
        ...device,
        channelInfo: device.channel,
        channel: device.channel,
        error: false,
      }
      selectedDevices.value.push(deviceWithChannel)
      const index = selectedDevices.value.length - 1

      await startStream(deviceWithChannel, index)
      saveLayoutState()
    } catch (error) {
      console.error('添加设备失败:', error)
      ElMessage.error('添加设备失败')
    }
  } else {
    ElMessage.warning('已达到最大分屏数量')
  }
}

const removeDevice = async (index: number) => {
  const device = selectedDevices.value[index]
  if (!device) return

  stopAutoReconnect(index)
  await cleanupDevice(device)

  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (videoElement?.srcObject) {
    const stream = videoElement.srcObject as MediaStream
    stream.getTracks().forEach((track) => track.stop())
    videoElement.srcObject = null
  }

  selectedDevices.value.splice(index, 1)
  saveLayoutState()
}

// 自动重连管理
const stopAutoReconnect = (index: number) => {
  if (reconnectTimers[index]) {
    clearInterval(reconnectTimers[index])
    delete reconnectTimers[index]
  }
}

const startAutoReconnect = (index: number) => {
  if (!autoReconnectEnabled.value) return
  stopAutoReconnect(index)

  const timer = window.setInterval(async () => {
    const device = selectedDevices.value[index]
    if (device?.error) {
      console.log(`尝试重连设备: ${device.name}`)
      await retryStream(index)
    } else {
      stopAutoReconnect(index)
    }
  }, reconnectInterval.value * 1000)

  reconnectTimers[index] = timer
}

// 布局状态管理
const saveLayoutState = () => {
  try {
    localStorage.setItem(
      'monitorGridLayout',
      JSON.stringify({
        layout: currentLayout.value,
        devices: selectedDevices.value.map((d) => ({
          name: d.name,
          channelInfo: d.channelInfo,
        })),
      }),
    )
  } catch (err) {
    console.error('保存布局状态失败:', err)
  }
}

const restoreLayoutState = async () => {
  try {
    const savedState = localStorage.getItem('monitorGridLayout')
    if (savedState) {
      const { layout, devices } = JSON.parse(savedState)
      currentLayout.value = layout
      for (const device of devices) {
        if (device.channelInfo) {
          await addDevice({ ...device, channel: device.channelInfo })
        }
      }
    }
  } catch (err) {
    console.error('恢复布局状态失败:', err)
  }
}

// 视频控制
const handleVideoDoubleClick = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (!videoElement) return

  try {
    if (!document.fullscreenElement) {
      videoElement.requestFullscreen()
    } else {
      document.exitFullscreen()
    }
  } catch (err) {
    console.error('视频全屏切换失败:', err)
    ElMessage.error('全屏切换失败')
  }
}

const handleVideoError = (index: number, event: Event) => {
  console.error('视频播放错误:', event)
  const device = selectedDevices.value[index]
  if (device) {
    device.error = true
    if (autoReconnectEnabled.value) {
      startAutoReconnect(index)
    }
  }
}

const toggleGridFullscreen = async () => {
  const gridContainer = document.querySelector('.grid-container') as HTMLElement
  if (!gridContainer) {
    console.error('未找到视频网格容器')
    return
  }

  try {
    if (!document.fullscreenElement) {
      await gridContainer.requestFullscreen()
      isFullscreen.value = true
    } else {
      await document.exitFullscreen()
      isFullscreen.value = false
    }
  } catch (err) {
    console.error('全屏切换失败:', err)
    ElMessage.error('全屏切换失败')
  }
}

const captureImage = async (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (!videoElement) {
    ElMessage.error('未找到视频元素')
    return
  }

  try {
    const canvas = document.createElement('canvas')
    canvas.width = videoElement.videoWidth
    canvas.height = videoElement.videoHeight
    const ctx = canvas.getContext('2d')
    if (!ctx) {
      throw new Error('无法创建canvas上下文')
    }

    ctx.drawImage(videoElement, 0, 0, canvas.width, canvas.height)

    const device = selectedDevices.value[index]
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-')
    const filename = `${device.name || 'capture'}-${timestamp}.png`

    const link = document.createElement('a')
    link.download = filename
    link.href = canvas.toDataURL('image/png')
    link.click()

    ElMessage.success('抓图成功')
  } catch (err) {
    console.error('抓图失败:', err)
    ElMessage.error('抓图失败')
  }
}

const getControlSize = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (!videoElement) return { btnSize: 24, iconSize: 12 }

  const width = videoElement.clientWidth
  if (width < 300) {
    return { btnSize: 20, iconSize: 10 }
  } else if (width < 500) {
    return { btnSize: 24, iconSize: 12 }
  } else if (width < 800) {
    return { btnSize: 28, iconSize: 14 }
  } else {
    return { btnSize: 32, iconSize: 16 }
  }
}

// 清空所有设备
const clearAllDevices = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有设备吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    while (selectedDevices.value.length > 0) {
      await removeDevice(0)
    }

    ElMessage.success('已清空所有设备')
  } catch (err) {
    if (err !== 'cancel') {
      console.error('清空设备失败:', err)
      ElMessage.error('清空设备失败')
    }
  }
}

// 布局切换处理
watch(currentLayout, async (newLayout, oldLayout) => {
  const maxSize = layouts[newLayout].size
  if (selectedDevices.value.length > maxSize) {
    try {
      await ElMessageBox.confirm(
        `切换布局将移除${selectedDevices.value.length - maxSize}个设备，是否继续？`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        },
      )

      const devicesToRemove = selectedDevices.value.slice(maxSize)
      for (const device of devicesToRemove) {
        const index = selectedDevices.value.indexOf(device)
        if (index !== -1) {
          await removeDevice(index)
        }
      }
      ElMessage.success('布局切换成功')
    } catch (err) {
      if (err !== 'cancel') {
        console.error('布局切换失败:', err)
        ElMessage.error('布局切换失败')
      }
      currentLayout.value = oldLayout
    }
  }
  saveLayoutState()
})

// 生命周期钩子
onMounted(async () => {
  await restoreLayoutState()

  document.addEventListener('fullscreenchange', () => {
    isFullscreen.value = !!document.fullscreenElement
  })
})

onBeforeUnmount(() => {
  reconnectTimers.forEach((timer) => {
    if (timer) clearInterval(timer)
  })

  selectedDevices.value.forEach((device) => {
    if (device.player) {
      try {
        device.player.close()
      } catch (err) {
        console.error('关闭播放器失败:', err)
      }
    }
  })
})

defineExpose({
  addDevice,
  clearAllDevices,
})

const toggleMute = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (!videoElement) return

  const device = selectedDevices.value[index]
  if (!device) return

  device.isMuted = !device.isMuted
  videoElement.muted = device.isMuted
}
</script>

<template>
  <div class="monitor-grid">
    <div class="grid-toolbar">
      <div class="toolbar-left">
        <el-button-group>
          <el-button
            v-for="(layout, key) in layouts"
            :key="key"
            :type="currentLayout === key ? 'primary' : ''"
            @click="currentLayout = key"
          >
            {{ layout.label }}
          </el-button>
        </el-button-group>
      </div>
      <div class="toolbar-right">
        <el-button-group>
          <el-button type="primary" @click="showSettings = true" :title="'设置'">
            <el-icon><Setting /></el-icon>
          </el-button>
          <el-button type="danger" @click="clearAllDevices" :title="'清空所有设备'">
            清空
          </el-button>
          <el-button
            type="primary"
            @click="toggleGridFullscreen"
            :title="isFullscreen ? '退出全屏' : '全屏显示'"
          >
            <el-icon><FullScreen /></el-icon>
          </el-button>
        </el-button-group>
      </div>
    </div>

    <div class="grid-container" :class="{ 'is-fullscreen': isFullscreen }" :style="gridStyle">
      <div v-for="i in maxDevices" :key="i" class="grid-item">
        <template v-if="selectedDevices[i - 1]">
          <div class="video-container" :class="{ 'has-error': selectedDevices[i - 1].error }">
            <div class="video-placeholder">
              <video
                :id="'video-player-' + (i - 1)"
                width="100%"
                height="100%"
                autoplay
                @error="handleVideoError(i - 1, $event)"
                @dblclick="handleVideoDoubleClick(i - 1)"
              ></video>

              <!-- 错误状态 -->
              <div v-if="selectedDevices[i - 1].error" class="error-overlay">
                <span>播放失败</span>
                <el-button type="primary" size="small" @click="retryStream(i - 1)">
                  <el-icon><Refresh /></el-icon>
                  重试
                </el-button>
              </div>
            </div>

            <div class="video-controls">
              <div class="control-bar">
                <el-button
                  type="primary"
                  class="control-btn"
                  @click="toggleMute(i - 1)"
                  :style="{
                    width: getControlSize(i - 1).btnSize + 'px',
                    height: getControlSize(i - 1).btnSize + 'px',
                    fontSize: getControlSize(i - 1).iconSize + 'px',
                  }"
                  :title="selectedDevices[i - 1].isMuted ? '取消静音' : '静音'"
                >
                  <el-icon>
                    <component :is="selectedDevices[i - 1].isMuted ? 'Mute' : 'Microphone'" />
                  </el-icon>
                </el-button>
                <el-button
                  type="primary"
                  class="control-btn"
                  @click="captureImage(i - 1)"
                  :style="{
                    width: getControlSize(i - 1).btnSize + 'px',
                    height: getControlSize(i - 1).btnSize + 'px',
                    fontSize: getControlSize(i - 1).iconSize + 'px',
                  }"
                  :disabled="selectedDevices[i - 1].error"
                  :title="'抓图'"
                >
                  <el-icon><Camera /></el-icon>
                </el-button>
                <el-button
                  type="danger"
                  class="control-btn"
                  @click="removeDevice(i - 1)"
                  :style="{
                    width: getControlSize(i - 1).btnSize + 'px',
                    height: getControlSize(i - 1).btnSize + 'px',
                    fontSize: getControlSize(i - 1).iconSize + 'px',
                  }"
                  :title="'移除'"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </template>
        <div v-else class="empty-grid">
          <el-icon><VideoCamera /></el-icon>
          <span>未选择设备</span>
        </div>
      </div>
    </div>

    <!-- 设置对话框 -->
    <el-dialog v-model="showSettings" title="设置" width="400px" destroy-on-close>
      <el-form label-width="120px">
        <el-form-item label="自动重连">
          <el-switch v-model="autoReconnectEnabled" />
        </el-form-item>
        <el-form-item label="重连间隔(秒)">
          <el-input-number
            v-model="reconnectInterval"
            :min="5"
            :max="60"
            :step="5"
            :disabled="!autoReconnectEnabled"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showSettings = false">取消</el-button>
          <el-button type="primary" @click="showSettings = false"> 确定 </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.monitor-grid {
  height: 100%;
  background: #fff;
  border-radius: 4px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.grid-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  gap: 10px;
}

.grid-container {
  flex: 1;
  display: grid;
  gap: 4px;
  height: calc(100% - 60px);
  transition: all 0.3s ease;
}

.grid-container.is-fullscreen {
  padding: 20px;
  background: #000;
}

.grid-item {
  background: #000;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
  aspect-ratio: 16/9;
}

.video-container {
  width: 100%;
  height: 100%;
  position: relative;
  transition: all 0.3s ease;
}

.video-container.has-error .video-placeholder {
  border: 2px solid #f56c6c;
}

.video-placeholder {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: #1a1a1a;
  border: 2px solid transparent;
  transition: border-color 0.3s ease;
}

.video-placeholder video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.7);
  gap: 10px;
  z-index: 3;
  color: #f56c6c;
}

.video-controls {
  position: absolute;
  top: 0;
  right: 0;
  opacity: 0;
  transition: all 0.3s ease;
  z-index: 2;
}

.video-container:hover .video-controls {
  opacity: 1;
}

.control-bar {
  display: flex;
  background: rgba(0, 0, 0, 0.5);
  padding: 2px;
  gap: 2px;
  backdrop-filter: blur(4px);
}

.control-btn {
  padding: 0 !important;
  border: none !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  background: transparent !important;
  transition: all 0.2s ease !important;
  border-radius: 0 !important;
}

.control-btn:first-child {
  border-radius: 0 0 0 4px !important;
}

.control-btn:last-child {
  border-radius: 0 4px 0 0 !important;
}

.control-btn:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  transform: scale(1.1);
}

.control-btn {
  color: #fff !important;
}

.control-btn:hover {
  color: #409eff !important;
}

.control-btn.el-button--danger:hover {
  color: #f56c6c !important;
}

.control-btn .el-icon {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-grid {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.4);
  gap: 10px;
  background: #1a1a1a;
  border-radius: 4px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.empty-grid:hover {
  color: #409eff;
  background: #1d1d1d;
}

.empty-grid .el-icon {
  font-size: 32px;
  opacity: 0.6;
  transition: all 0.3s ease;
}

.empty-grid:hover .el-icon {
  opacity: 1;
  transform: scale(1.1);
}

:deep(.el-button-group) {
  .el-button {
    padding: 8px 15px;
    transition: all 0.3s ease;
  }

  .el-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
}

/* 全屏样式 */
.grid-container:fullscreen,
.grid-container:-webkit-full-screen,
.grid-container:-moz-full-screen,
.grid-container:-ms-fullscreen {
  background: #000;
  padding: 20px;
}

.grid-container.is-fullscreen {
  padding: 20px;
  background: #000;
  gap: 8px;
}

.grid-container.is-fullscreen .video-container {
  border-radius: 8px;
  overflow: hidden;
}

.video-container video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
  transition: all 0.3s ease;
}

.video-container:hover video {
  object-fit: contain;
}

.device-info {
  display: none;
}

.video-container:hover .device-info {
  display: none;
}

.control-bar {
  display: flex;
  background: rgba(0, 0, 0, 0.5);
  padding: 4px;
  border-radius: 4px;
  gap: 4px;
  backdrop-filter: blur(4px);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.empty-grid {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.4);
  gap: 10px;
  background: #1a1a1a;
  border-radius: 4px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.empty-grid:hover {
  color: #409eff;
  background: #1d1d1d;
}

.empty-grid .el-icon {
  font-size: 32px;
  opacity: 0.6;
  transition: all 0.3s ease;
}

.empty-grid:hover .el-icon {
  opacity: 1;
  transform: scale(1.1);
}

/* 优化全屏模式 */
.grid-container.is-fullscreen {
  padding: 20px;
  background: #000;
  gap: 8px;
}

.grid-container.is-fullscreen .video-container {
  border-radius: 8px;
  overflow: hidden;
}

/* 优化按钮样式 */
:deep(.el-button-group) {
  .el-button {
    padding: 8px 15px;
    transition: all 0.3s ease;
  }

  .el-button:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
}

/* 优化错误状态 */
.error-overlay {
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(4px);
}

.error-overlay .el-button {
  transition: all 0.3s ease;
}

.error-overlay .el-button:hover {
  transform: scale(1.1);
}

/* 添加点击效果样式 */
.video-container {
  width: 100%;
  height: 100%;
  position: relative;
  transition: all 0.3s ease;
}

.video-container:active {
  transform: none;
}
</style>
