<template>
  <div class="dashboard">
    <h1 class="dashboard-title">系统概览</h1>
    <div class="dashboard-grid">
      <div class="dashboard-card">
        <div class="card-header">
          <el-icon><VideoCamera /></el-icon>
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
          <el-icon><Monitor /></el-icon>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Monitor, VideoCamera } from '@element-plus/icons-vue'
import { useMediaServers, fetchMediaServers } from '@/stores/mediaServer'
import { useDevices, fetchDevicesAndChannels } from '@/stores/devices'

const mediaServers = useMediaServers()
const devices = useDevices()

const onlineServerCount = ref(0)
const totalServerCount = ref(0)
const onlineDeviceCount = ref(0)
const totalDeviceCount = ref(0)

// 更新数据的函数
const updateData = () => {
  // 服务器统计
  const servers = mediaServers.value
  onlineServerCount.value = servers.filter(server => server.status === 1).length
  totalServerCount.value = servers.length

  // 设备统计
  const devicesList = devices.value
  onlineDeviceCount.value = devicesList.filter(device => device.status === 'online').length
  totalDeviceCount.value = devicesList.length
}

// 定时器引用
let timer: number

onMounted(async () => {
  // 初始化数据
  await Promise.all([
    fetchMediaServers(),
    fetchDevicesAndChannels()
  ])
  
  // 初始更新
  updateData()
  // 设置定时更新
  timer = setInterval(updateData, 5000)
})

onUnmounted(() => {
  // 清理定时器
  if (timer) {
    clearInterval(timer)
  }
})
</script>

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
  gap: 12px;
  margin-bottom: 24px;
}

.card-header .el-icon {
  font-size: 24px;
  color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
  padding: 12px;
  border-radius: 8px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-regular);
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