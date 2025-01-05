<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount, onMounted, onActivated, onDeactivated } from 'vue'
import { VideoCamera, Close, Camera, FullScreen, Refresh, Setting, Mute, Microphone, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Device, ChannelInfo, MediaServer, ApiResponse } from '@/types/api'
import { deviceApi } from '@/api'
import { SrsRtcPlayer } from '@/api/srs'
import { useDefaultMediaServer } from '@/stores/mediaServer'
import type { LayoutConfig } from '@/types/layout'

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

const props = defineProps<{
  modelValue: string
  defaultMuted?: boolean
  layouts: Record<string, LayoutConfig>
  showBorder?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'window-select': [data: { deviceId: string; channelId: string } | null]
}>()

const currentLayout = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const selectedDevices = ref<(DeviceWithChannel | null)[]>([])
const isFullscreen = ref(false)

// 使用共享的默认媒体服务器
const defaultMediaServer = useDefaultMediaServer()

// 计算属性
const gridStyle = computed(() => {
  const layout = props.layouts[currentLayout.value]
  return {
    gridTemplateColumns: `repeat(${layout.cols}, 1fr)`,
    gridTemplateRows: `repeat(${layout.rows}, 1fr)`,
  }
})

const maxDevices = computed(() => props.layouts[currentLayout.value].size)

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

    if (!defaultMediaServer?.value) {
      throw new Error('未找到可用的媒体服务器，请先在流媒体服务页面设置默认服务器')
    }

    if (defaultMediaServer.value.status === 0) {
      throw new Error('默认流媒体服务器不在线，请检查服务器状态')
    }

    const mediaServerAddr = `${defaultMediaServer.value.ip}:${defaultMediaServer.value.port}`
    const response = await deviceApi.inviteStream({
      media_server_addr: mediaServerAddr,
      device_id: device.channel!.parent_id,
      channel_id: device.channel!.device_id,
      sub_stream: 0,
      play_type: 0,
      start_time: 0,
      end_time: 0,
    })

    const streamData = response.data as unknown as StreamResponse
    if (!streamData?.url) {
      throw new Error('播放地址不存在')
    }

    await startWebRTCPlay(streamData.url, index, device)
  } catch (error) {
    console.error('启动播放失败:', error)
    device.error = true
    ElMessage.error(error instanceof Error ? error.message : '启动播放失败')
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

const play = async (device: Device & { channel: ChannelInfo }) => {
  let index = selectedDevices.value.findIndex(d => d === null);
  if (index === -1) {
    if (selectedDevices.value.length >= maxDevices.value) {
      ElMessage.warning('已达到最大分屏数量')
      return
    }
    index = selectedDevices.value.length
  }

  try {
    const deviceWithChannel: DeviceWithChannel = {
      ...device,
      channelInfo: device.channel,
      channel: device.channel,
      error: false,
      isMuted: props.defaultMuted
    }
    
    while (selectedDevices.value.length <= index) {
      selectedDevices.value.push(null)
    }
    
    selectedDevices.value[index] = deviceWithChannel
    await startStream(deviceWithChannel, index)
    
    const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
    if (videoElement) {
      videoElement.muted = props.defaultMuted ?? true
    }
  } catch (error) {
    console.error('添加设备失败:', error)
    ElMessage.error('添加设备失败')
  }
}

const removeDevice = async (index: number) => {
  const device = selectedDevices.value[index]
  if (!device) return

  await cleanupDevice(device)

  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (videoElement?.srcObject) {
    const stream = videoElement.srcObject as MediaStream
    stream.getTracks().forEach((track) => track.stop())
    videoElement.srcObject = null
  }

  // 将位置设为 null 而不是删除
  selectedDevices.value[index] = null
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
    if (!device) {
      throw new Error('设备不存在')
    }
    
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

// 清空所有设备
const clearAllDevices = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有设备吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    for (let i = 0; i < selectedDevices.value.length; i++) {
      if (selectedDevices.value[i]) {
        await removeDevice(i)
      }
    }
    
    // 用 null 填充数组而不是清空
    selectedDevices.value = new Array(props.layouts[currentLayout.value].size).fill(null)

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
  const maxSize = props.layouts[newLayout].size
  const activeDevices = selectedDevices.value.filter(d => d !== null).length
  
  if (activeDevices > maxSize) {
    try {
      await ElMessageBox.confirm(
        `切换布局将移除${activeDevices - maxSize}个设备，是否继续？`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        },
      )

      // 从后往前移除超出设备
      for (let i = selectedDevices.value.length - 1; i >= 0; i--) {
        if (selectedDevices.value[i] && i >= maxSize) {
          await removeDevice(i)
        }
      }
      
      // 调整数组大小
      selectedDevices.value.length = maxSize
      
      ElMessage.success('布局切换成功')
    } catch (err) {
      if (err !== 'cancel') {
        console.error('布局切换失败:', err)
        ElMessage.error('布局切换失败')
      }
      currentLayout.value = oldLayout
    }
  } else {
    // 如果设备数量不超过新布局，只需调整数组大小
    if (selectedDevices.value.length > maxSize) {
      selectedDevices.value.length = maxSize
    } else {
      while (selectedDevices.value.length < maxSize) {
        selectedDevices.value.push(null)
      }
    }
  }
})

// 生命周期钩子
onMounted(() => {
  // 确保有初始网格
  if (selectedDevices.value.length === 0) {
    selectedDevices.value = new Array(props.layouts[currentLayout.value].size).fill(null)
  }

  document.addEventListener('fullscreenchange', () => {
    isFullscreen.value = !!document.fullscreenElement
  })
})

onBeforeUnmount(() => {
  // 只在组件真正销毁时清理资源
  if (!(document as any)._vue_app_is_switching_route) {
    selectedDevices.value.forEach((device) => {
      if (device?.player) {
        try {
          device.player.close()
        } catch (err) {
          console.error('关闭播放器失败:', err)
        }
      }
    })
  }
})

defineExpose({
  play,
  clearAllDevices
})

const toggleMute = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
  if (!videoElement) return

  const device = selectedDevices.value[index]
  if (!device) return

  device.isMuted = !device.isMuted
  videoElement.muted = device.isMuted
}

const activeIndex = ref<number | null>(null)

const handleVideoClick = (index: number) => {
  const device = selectedDevices.value[index]
  activeIndex.value = index
  
  if (device && device.channel) {
    emit('window-select', {
      deviceId: device.device_id,
      channelId: device.channel.device_id
    })
  } else {
    emit('window-select', null)
  }
}

// 添加激活/停用处理
onActivated(() => {
  console.log('MonitorGrid activated')
  // 检查并恢复所有视频播放
  selectedDevices.value.forEach((device, index) => {
    if (device) {
      const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
      if (videoElement && videoElement.paused) {
        videoElement.play().catch(err => {
          console.error('恢复视频播放失败:', err)
        })
      }
    }
  })
})

onDeactivated(() => {
  console.log('MonitorGrid deactivated')
  // 可以选择暂停视频播放，但不销毁资源
  selectedDevices.value.forEach((device, index) => {
    if (device) {
      const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
      if (videoElement && !videoElement.paused) {
        // 可以选择是否暂停视频
        // videoElement.pause()
      }
    }
  })
})
</script>

<template>
  <div class="monitor-grid">
    <div class="grid-container" :style="gridStyle">
      <div
        v-for="(device, index) in selectedDevices"
        :key="index"
        class="grid-item"
        :class="{ 
          active: index === activeIndex,
          'with-border': props.showBorder 
        }"
        @click="handleVideoClick(index)"
      >
        <template v-if="device !== null">
          <video
            :id="`video-player-${index}`"
            class="video-player"
            autoplay
            :muted="device.isMuted ?? true"
            @dblclick="handleVideoDoubleClick(index)"
            @error="handleVideoError(index, $event)"
          />
          <div class="video-overlay">
            <div class="device-info">
              {{ device.name }} - {{ device.channel?.name ?? '' }}
            </div>
            <div class="video-controls">
              <div class="control-bar">
                <el-button
                  class="control-btn"
                  @click.stop="toggleMute(index)"
                  :title="device.isMuted ? '取消静音' : '静音'"
                >
                  <el-icon>
                    <component :is="device.isMuted ? 'Mute' : 'Microphone'" />
                  </el-icon>
                </el-button>
                <el-button
                  class="control-btn"
                  @click.stop="captureImage(index)"
                  :title="'抓图'"
                  :disabled="device.error"
                >
                  <el-icon><Camera /></el-icon>
                </el-button>
                <el-button
                  class="control-btn"
                  @click.stop="handleVideoDoubleClick(index)"
                  :title="'全屏'"
                >
                  <el-icon><FullScreen /></el-icon>
                </el-button>
                <el-button
                  class="control-btn is-danger"
                  @click.stop="removeDevice(index)"
                  :title="'关闭'"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </template>
        <div v-else class="empty-slot">
          <el-icon><VideoCamera /></el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.monitor-grid {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color);
  border-radius: 4px;
  box-shadow: var(--el-box-shadow-lighter);
}

.grid-container {
  flex: 1;
  display: grid;
  gap: 0px;
  padding: 0px;
  height: 100%;
  background: var(--el-bg-color-page);
  border-radius: 4px;
  
  &.is-fullscreen {
    padding: 16px;
    background: #000;
    gap: 16px;
  }
}

.grid-item {
  position: relative;
  background-color: var(--el-fill-color-darker);
  transition: border-color 0.2s ease;
  cursor: pointer;
  overflow: hidden;
  
  &.with-border {
    border: 2px solid transparent;
    
    &.active {
      border-color: var(--el-color-primary);
    }
  }
}

.video-player {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background: #000;
}

.video-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  padding: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.5) 0%, transparent 30%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.grid-item:hover .video-overlay {
  opacity: 1;
}

.device-info {
  display: none;
}

.video-controls {
  display: flex;
  justify-content: flex-end;
}

.control-bar {
  display: flex;
  gap: 2px;
  background: rgba(0, 0, 0, 0.5);
  padding: 2px;
  border-radius: 0 0 0 4px;
  backdrop-filter: blur(4px);
}

.control-btn {
  padding: 0 !important;
  border: none !important;
  background: transparent !important;
  color: #fff !important;
  
  &:hover {
    background: rgba(255, 255, 255, 0.1) !important;
    color: var(--el-color-primary) !important;
  }
  
  &.is-danger:hover {
    color: var(--el-color-danger) !important;
  }
}

/* 根据布局调整按钮大小 */
.grid-container[style*="repeat(1, 1fr)"] .control-btn {
  width: 32px !important;
  height: 32px !important;
  font-size: 16px !important;
}

.grid-container[style*="repeat(2, 1fr)"] .control-btn {
  width: 28px !important;
  height: 28px !important;
  font-size: 14px !important;
}

.grid-container[style*="repeat(3, 1fr)"] .control-btn {
  width: 20px !important;
  height: 20px !important;
  font-size: 10px !important;
}

.grid-container[style*="repeat(4, 1fr)"] .control-btn {
  width: 16px !important;
  height: 16px !important;
  font-size: 8px !important;
}

.control-bar {
  .el-icon {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.empty-slot {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1a1a1a;
  
  .el-icon {
    font-size: 24px;
    opacity: 0.5;
    color: #fff;
  }

  &:hover {
    background: #242424;
    
    .el-icon {
      opacity: 0.8;
      color: var(--el-color-primary);
    }
  }
}

.monitor-grid {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color);
  border-radius: 4px;
  box-shadow: var(--el-box-shadow-lighter);
}

.grid-toolbar {
  padding: 12px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--el-bg-color);
  border-radius: 4px 4px 0 0;
}

:deep(.el-button-group .el-button--small) {
  padding: 5px 11px;
}

:deep(.el-radio-group .el-radio-button__inner) {
  padding: 5px 15px;
}

.video-container {
  width: 100%;
  height: 100%;
  position: relative;
  transition: all 0.3s ease;
}

.video-container:active {
  transform: none;
}

:deep(.el-radio-group) {
  --el-button-bg-color: var(--el-fill-color-blank);
  --el-button-hover-bg-color: var(--el-fill-color);
  --el-button-active-bg-color: var(--el-color-primary);
  --el-button-text-color: var(--el-text-color-regular);
  --el-button-hover-text-color: var(--el-text-color-primary);
  --el-button-active-text-color: #fff;
  --el-button-border-color: var(--el-border-color);
  --el-button-hover-border-color: var(--el-border-color-hover);
}
</style>
