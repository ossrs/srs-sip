<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { VideoCamera, Close, FullScreen } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { Device, ChannelInfo } from '@/types/api'
import { deviceApi } from '@/api'

// 布局配置
const layouts = {
  '1': { cols: 1, rows: 1, size: 1 },
  '4': { cols: 2, rows: 2, size: 4 },
  '9': { cols: 3, rows: 3, size: 9 },
}

const selectedDevices = ref<Device[]>([])
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
        sub_stream: "0"
      })
      
      selectedDevices.value.push(device)
      
    } catch (error) {
      console.error('invite请求失败:', error)
    }
  } else {
    ElMessage.warning('已达到最大分屏数量')
  }
}

const removeDevice = (index: number) => {
  selectedDevices.value.splice(index, 1)
}

// 添加布局切换时的设备处理
watch(currentLayout, (newLayout) => {
  const maxSize = layouts[newLayout].size
  if (selectedDevices.value.length > maxSize) {
    selectedDevices.value = selectedDevices.value.slice(0, maxSize)
    ElMessage.warning('已自动移除超出布局限制的设备')
  }
})

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
        </el-button-group>
      </div>
    </div>
    <div class="grid-container" :style="gridStyle">
      <div v-for="i in maxDevices" :key="i" class="grid-item">
        <template v-if="selectedDevices[i - 1]">
          <div class="video-container">
            <div class="video-placeholder">
              <span>{{ selectedDevices[i - 1].name }}</span>
              <div class="channel-info" v-if="selectedDevices[i - 1].channelInfo">
                <span>通道: {{ selectedDevices[i - 1].channelInfo.name }}</span>
                <span>状态: {{ selectedDevices[i - 1].channelInfo.status }}</span>
              </div>
            </div>
            <div class="video-controls">
              <el-button-group>
                <el-button type="primary" circle size="small">
                  <el-icon><FullScreen /></el-icon>
                </el-button>
                <el-button type="danger" circle size="small" @click="removeDevice(i - 1)">
                  <el-icon><Close /></el-icon>
                </el-button>
              </el-button-group>
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
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: #1a1a1a;
}

.video-controls {
  position: absolute;
  top: 10px;
  right: 10px;
  opacity: 0;
  transition: opacity 0.3s;
}

.video-container:hover .video-controls {
  opacity: 1;
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

.channel-info {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
</style>
