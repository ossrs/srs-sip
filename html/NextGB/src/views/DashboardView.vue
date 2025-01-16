<script setup lang="ts">
import { computed } from 'vue'
import { useMediaServers } from '@/stores/mediaServer'
import { useDevices } from '@/stores/devices'
import type { MediaServer } from '@/api/mediaserver/types'
import type { Device } from '@/api/types'
import { createMediaServer } from '@/api/mediaserver/factory'
import { ref, onMounted } from 'vue'

const mediaServers = useMediaServers()
const devices = useDevices()

const onlineServerCount = computed(
  () => mediaServers.value.filter((server: MediaServer) => server.status === 1).length,
)
const totalServerCount = computed(() => mediaServers.value.length)
const onlineDeviceCount = computed(
  () => devices.value.filter((device: Device) => device.status === 1).length,
)
const totalDeviceCount = computed(() => devices.value.length)

const totalStreams = ref(0)
const totalPlayers = ref(0)

const fetchStreamAndPlayerCount = async () => {
  let streamCount = 0
  let playerCount = 0
  for (const server of mediaServers.value) {
    if (server.status === 1) { // 只统计在线服务器
      try {
        const mediaServer = createMediaServer(server)
        const streams = await mediaServer.getStreamInfo()
        streamCount += streams.length
        // 统计所有流的客户端数量
        playerCount += streams.reduce((sum, stream) => sum + (stream.clients || 0), 0)
      } catch (error) {
        console.error(`获取服务器 ${server.name} 的流信息失败:`, error)
      }
    }
  }
  totalStreams.value = streamCount
  totalPlayers.value = playerCount
}

// 每30秒更新一次统计数据
onMounted(() => {
  fetchStreamAndPlayerCount()
  setInterval(fetchStreamAndPlayerCount, 30000)
})
</script>

<template>
  <div class="dashboard">
    <h1 class="dashboard-title">系统概览</h1>
    <div class="dashboard-grid">
      <div class="dashboard-card">
        <div class="card-header">
          <span class="card-title">流媒体服务器</span>
        </div>
        <div class="card-content">
          <div class="number">
            <span class="online">{{ onlineServerCount }}</span>
            <span class="separator">/</span>
            <span class="total">{{ totalServerCount }}</span>
          </div>
          <div class="label">在线/总数</div>
        </div>
      </div>

      <div class="dashboard-card">
        <div class="card-header">
          <span class="card-title">设备状态</span>
        </div>
        <div class="card-content">
          <div class="number">
            <span class="online">{{ onlineDeviceCount }}</span>
            <span class="separator">/</span>
            <span class="total">{{ totalDeviceCount }}</span>
          </div>
          <div class="label">在线/总数</div>
        </div>
      </div>

      <div class="dashboard-card">
        <div class="card-header">
          <span class="card-title">流数量</span>
        </div>
        <div class="card-content">
          <div class="number">
            <span class="total">{{ totalStreams }}</span>
          </div>
          <div class="label">总流数</div>
        </div>
      </div>

      <div class="dashboard-card">
        <div class="card-header">
          <span class="card-title">播放者数量</span>
        </div>
        <div class="card-content">
          <div class="number">
            <span class="total">{{ totalPlayers }}</span>
          </div>
          <div class="label">总播放数</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 24px;
  height: 100%;
  background-color: var(--el-bg-color-page);
}

.dashboard-title {
  margin: 0 0 24px 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.dashboard-card {
  background: var(--el-bg-color);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--el-box-shadow-light);
  transition: all 0.3s ease;
}

.dashboard-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--el-box-shadow);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.card-content {
  text-align: center;
}

.number {
  font-size: 48px;
  font-weight: 600;
  line-height: 1.2;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.number .online {
  color: var(--el-color-success);
}

.number .separator {
  color: var(--el-text-color-secondary);
  font-size: 36px;
  margin: 0 4px;
}

.number .total {
  color: var(--el-text-color-regular);
}

.label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .dashboard {
    padding: 16px;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .number {
    font-size: 36px;
  }

  .number .separator {
    font-size: 28px;
  }
}
</style>
