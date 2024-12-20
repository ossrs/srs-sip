<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { VideoCamera, Close, Camera, FullScreen } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { Device, ChannelInfo } from '@/types/api'
import { deviceApi } from '@/api'
import { SrsRtcPlayer } from '@/api/srs'

interface DeviceWithChannel extends Device {
  channelInfo?: ChannelInfo;
}

interface StreamResponse {
  url: string;
  [key: string]: any;
}

// 布局配置
const layouts = {
  '1': { cols: 1, rows: 1, size: 1 },
  '4': { cols: 2, rows: 2, size: 4 },
  '9': { cols: 3, rows: 3, size: 9 },
}

const selectedDevices = ref<DeviceWithChannel[]>([])
const currentLayout = ref<keyof typeof layouts>('9')

const gridStyle = computed(() => {
  const layout = layouts[currentLayout.value]
  return {
    gridTemplateColumns: `repeat(${layout.cols}, 1fr)`,
    gridTemplateRows: `repeat(${layout.rows}, 1fr)`,
  }
})

const maxDevices = computed(() => layouts[currentLayout.value].size)

const addDevice = async (device: Device & { channel: ChannelInfo }) => {
  if (selectedDevices.value.length < maxDevices.value) {
    try {
      const response = await deviceApi.inviteStream({
        device_id: device.channel.parent_id,
        channel_id: device.channel.device_id,
        sub_stream: '0',
      })

      console.log('invite请求成功:', response)

      // 使用响应中的url进行WebRTC播放
      const streamData = response.data as unknown as StreamResponse
      if (!streamData?.url) {
        throw new Error('播放地址不存在')
      }
      startWebRTCPlay(streamData.url, selectedDevices.value.length)

      const deviceWithChannel: DeviceWithChannel = {
        ...device,
        channelInfo: device.channel
      }
      selectedDevices.value.push(deviceWithChannel)
    } catch (error) {
      console.error('invite请求失败:', error)
      ElMessage.error('启动播放失败')
    }
  } else {
    ElMessage.warning('已达到最大分屏数量')
  }
}

// 新增函数，用于启动WebRTC播放
function startWebRTCPlay(url: string, index: number) {
  const player = SrsRtcPlayer();
  // @ts-ignore - 忽略类型检查，因为SrsRtcPlayer的类型定义可能不完整
  player.onaddstream = (event: { stream: MediaStream }) => {
    console.log('Start play, event: ', event);
    const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement;
    if (videoElement) {
      videoElement.srcObject = event.stream;
    } else {
      console.error('未找到对应的 video 元素');
    }
  };

  player.play(url).catch(error => {
    console.error('播放失败:', error);
    ElMessage.error('视频播放失败');
  });
}

const removeDevice = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement;
  if (videoElement && videoElement.srcObject) {
    const stream = videoElement.srcObject as MediaStream;
    stream.getTracks().forEach(track => track.stop());
    videoElement.srcObject = null;
  }
  selectedDevices.value.splice(index, 1);
}

// 添加布切换时的设备处理
watch(currentLayout, (newLayout) => {
  const maxSize = layouts[newLayout].size
  if (selectedDevices.value.length > maxSize) {
    selectedDevices.value = selectedDevices.value.slice(0, maxSize)
    ElMessage.warning('已自动移除超出布局限制的设备')
  }
})

// 抓图功能
const captureImage = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement;
  if (!videoElement) {
    console.error('未找到视频元素');
    return;
  }

  try {
    const canvas = document.createElement('canvas');
    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      throw new Error('无法创建canvas上下文');
    }
    
    ctx.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
    
    // 将图片保存为文件
    const link = document.createElement('a');
    link.download = `capture-${Date.now()}.png`;
    link.href = canvas.toDataURL('image/png');
    link.click();
    
    ElMessage.success('抓图成功');
  } catch (err) {
    console.error('抓图失败:', err);
    ElMessage.error('抓图失败');
  }
}

// 计算控制按钮大小
const getControlSize = (index: number) => {
  const videoElement = document.getElementById(`video-player-${index}`) as HTMLVideoElement;
  if (!videoElement) return { btnSize: 24, iconSize: 12 };
  
  const width = videoElement.clientWidth;
  // 根据视频宽度计算按钮和图标大小
  if (width < 300) {
    return { btnSize: 20, iconSize: 10 };
  } else if (width < 500) {
    return { btnSize: 24, iconSize: 12 };
  } else if (width < 800) {
    return { btnSize: 28, iconSize: 14 };
  } else {
    return { btnSize: 32, iconSize: 16 };
  }
}

// 全屏切换函数
const toggleGridFullscreen = () => {
  const gridContainer = document.querySelector('.grid-container') as HTMLElement;
  if (!gridContainer) {
    console.error('未找到视频网格容器');
    return;
  }

  try {
    if (!document.fullscreenElement) {
      // 进入全屏
      gridContainer.requestFullscreen();
    } else {
      // 退出全屏
      document.exitFullscreen();
    }
  } catch (err) {
    console.error('全屏切换失败:', err);
    ElMessage.error('全屏切换失败');
  }
}

defineExpose({
  addDevice,
})
</script>

<template>
  <div class="monitor-grid">
    <div class="grid-toolbar">
      <div class="layout-buttons">
        <el-button-group>
          <el-button :type="currentLayout === '1' ? 'primary' : ''" @click="currentLayout = '1'">
            单屏
          </el-button>
          <el-button :type="currentLayout === '4' ? 'primary' : ''" @click="currentLayout = '4'">
            四分屏
          </el-button>
          <el-button :type="currentLayout === '9' ? 'primary' : ''" @click="currentLayout = '9'">
            九分屏
          </el-button>
          <el-button type="primary" @click="toggleGridFullscreen">
            <el-icon><FullScreen /></el-icon>
          </el-button>
        </el-button-group>
      </div>
    </div>
    <div class="grid-container" :style="gridStyle">
      <div v-for="i in maxDevices" :key="i" class="grid-item">
        <template v-if="selectedDevices[i - 1]">
          <div class="video-container">
            <div class="video-placeholder">
              <video :id="'video-player-' + (i-1)" width="100%" height="100%" controls autoplay></video>
              <div class="device-info">
                <span>{{ selectedDevices[i - 1]?.name }}</span>
              </div>
            </div>
            <div class="video-controls">
              <div class="control-bar">
                <el-button 
                  type="primary" 
                  class="control-btn" 
                  @click="captureImage(i - 1)"
                  :style="{ 
                    width: getControlSize(i - 1).btnSize + 'px', 
                    height: getControlSize(i - 1).btnSize + 'px',
                    fontSize: getControlSize(i - 1).iconSize + 'px'
                  }"
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
                    fontSize: getControlSize(i - 1).iconSize + 'px'
                  }"
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
  justify-content: flex-end;
  align-items: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.grid-container {
  flex: 1;
  display: grid;
  gap: 2px;
  height: calc(100% - 60px);
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
}

.video-placeholder video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.device-info {
  position: absolute;
  top: 10px;
  left: 10px;
  z-index: 1;
  background: rgba(0, 0, 0, 0.5);
  padding: 5px 10px;
  border-radius: 4px;
}

.video-controls {
  position: absolute;
  top: 10px;
  right: 10px;
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
  border-radius: 3px;
  gap: 1px;
}

.control-btn {
  padding: 0 !important;
  border: none !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  background: transparent !important;
  transition: all 0.2s ease !important;
  border-radius: 2px !important;
}

.control-btn:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  transform: scale(1.1);
}

.control-btn {
  color: #fff !important;
}

.control-btn:hover {
  color: #409EFF !important;
}

.control-btn.el-button--danger:hover {
  color: #F56C6C !important;
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
  color: #666;
  gap: 10px;

  .el-icon {
    font-size: 32px;
  }
}

:deep(.el-button-group) {
  .el-button {
    padding: 8px 15px;
  }
}

/* 添加全屏样式 */
.grid-container:fullscreen {
  background: #000;
  padding: 20px;
}

.grid-container:-webkit-full-screen {
  background: #000;
  padding: 20px;
}

.grid-container:-moz-full-screen {
  background: #000;
  padding: 20px;
}

.grid-container:-ms-fullscreen {
  background: #000;
  padding: 20px;
}
</style>
