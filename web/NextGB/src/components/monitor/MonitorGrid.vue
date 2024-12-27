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
} as const

type LayoutKey = keyof typeof layouts

const currentLayout = ref<LayoutKey>('9')
const selectedDevices = ref<(DeviceWithChannel | null)[]>([])
const isFullscreen = ref(false)
const showSettings = ref(false)
const defaultMuted = ref(true)

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
      isMuted: defaultMuted.value
    }
    
    while (selectedDevices.value.length <= index) {
      selectedDevices.value.push(null)
    }
    
    selectedDevices.value[index] = deviceWithChannel
    await startStream(deviceWithChannel, index)
    
    const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement
    if (videoElement) {
      videoElement.muted = defaultMuted.value
    }
    
    saveLayoutState()
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
  saveLayoutState()
}

// 布局状态管理
const saveLayoutState = () => {
  try {
    localStorage.setItem(
      'monitorGridLayout',
      JSON.stringify({
        layout: currentLayout.value,
        devices: selectedDevices.value
          .map((d, index) => d ? {
            name: d.name,
            channelInfo: d.channelInfo,
            index: index // 保存位置信息
          } : null)
          .filter(d => d !== null) // 只保存非空设备
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
      
      // 初始化数组大小
      selectedDevices.value = new Array(layouts[layout].size).fill(null)
      
      // 恢复设备到原来的位置
      for (const device of devices) {
        if (device?.channelInfo) {
          const index = device.index
          if (index < layouts[layout].size) {
            await addDevice({ ...device, channel: device.channelInfo })
          }
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

    for (let i = 0; i < selectedDevices.value.length; i++) {
      if (selectedDevices.value[i]) {
        await removeDevice(i)
      }
    }
    
    // 用 null 填充数组而不是清空
    selectedDevices.value = new Array(layouts[currentLayout.value].size).fill(null)

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
  saveLayoutState()
})

// 生命周期钩子
onMounted(async () => {
  // 确保有初始网格
  if (selectedDevices.value.length === 0) {
    selectedDevices.value = new Array(layouts[currentLayout.value].size).fill(null)
  }
  
  await restoreLayoutState()

  document.addEventListener('fullscreenchange', () => {
    isFullscreen.value = !!document.fullscreenElement
  })
})

onBeforeUnmount(() => {
  selectedDevices.value.forEach((device) => {
    if (device?.player) {
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

const emit = defineEmits<{
  'window-select': [data: { deviceId: string; channelId: string } | null]
}>()

const activeIndex = ref<number | null>(null)

const handleVideoClick = (index: number) => {
  const device = selectedDevices.value[index]
  activeIndex.value = index
  
  if (device && device.channel) {
    emit('window-select', {
      deviceId: device.channel.parent_id,
      channelId: device.channel.device_id
    })
  } else {
    emit('window-select', null)
  }
}
</script>

<template>
  <div class="monitor-grid">
    <div class="grid-toolbar">
      <div class="layout-controls">
        <el-radio-group v-model="currentLayout" size="small">
          <el-radio-button v-for="(layout, key) in layouts" :key="key" :label="key as LayoutKey">
            {{ layout.label }}
          </el-radio-button>
        </el-radio-group>
      </div>
      <div class="toolbar-actions">
        <el-button-group>
          <el-button size="small" @click="toggleGridFullscreen">
            <el-icon><FullScreen /></el-icon>
          </el-button>
        </el-button-group>
      </div>
    </div>
    <div class="grid-container" :style="gridStyle">
      <div
        v-for="(device, index) in selectedDevices"
        :key="index"
        class="grid-item"
        :class="{ active: index === activeIndex }"
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
              <el-button-group>
                <el-button size="small" @click.stop="retryStream(index)" v-if="device.error">
                  <el-icon><Refresh /></el-icon>
                </el-button>
                <el-button size="small" @click.stop="removeDevice(index)">
                  <el-icon><Close /></el-icon>
                </el-button>
              </el-button-group>
            </div>
          </div>
        </template>
        <div v-else class="empty-slot">
          <el-icon><VideoCamera /></el-icon>
          <span>点击添加视频</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.monitor-grid {
  height: 100%;
  background: #fff;
  border-radius: 4px;
  padding: 10px;
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
  position: relative;
  background-color: var(--el-fill-color-darker);
  border: 2px solid transparent;
  transition: border-color 0.2s ease;
  cursor: pointer;
  aspect-ratio: 16/9;
  overflow: hidden;
  
  &.active {
    border-color: var(--el-color-primary);
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
  justify-content: space-between;
  padding: 8px;
  background: linear-gradient(to bottom, rgba(0,0,0,0.5) 0%, transparent 30%, transparent 70%, rgba(0,0,0,0.5) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.grid-item:hover .video-overlay {
  opacity: 1;
}

.device-info {
  color: #fff;
  font-size: 12px;
  text-shadow: 0 1px 2px rgba(0,0,0,0.5);
}

.video-controls {
  display: flex;
  justify-content: flex-end;
  gap: 4px;
}

.empty-slot {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  background: var(--el-fill-color-light);
  
  .el-icon {
    font-size: 24px;
    opacity: 0.7;
  }

  &:hover {
    background: var(--el-fill-color);
    
    .el-icon {
      opacity: 1;
      color: var(--el-color-primary);
    }
  }
}

.grid-container {
  flex: 1;
  display: grid;
  gap: 8px;
  padding: 8px;
  height: calc(100% - 60px);
  background: var(--el-bg-color-page);
  border-radius: 4px;
  
  &.is-fullscreen {
    padding: 16px;
    background: #000;
    gap: 16px;
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
</style>
