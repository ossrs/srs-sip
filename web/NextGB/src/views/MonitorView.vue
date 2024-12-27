<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import DeviceTree from '@/components/monitor/DeviceTree.vue'
import MonitorGrid from '@/components/monitor/MonitorGrid.vue'
import PtzControlPanel from '@/components/monitor/PtzControlPanel.vue'
import type { Device, ChannelInfo } from '@/types/api'

const monitorGridRef = ref()
const selectedChannel = ref<{ device: Device | undefined; channel: ChannelInfo } | null>(null)
const activeWindow = ref<{ deviceId: string; channelId: string } | null>(null)

const handleDeviceSelect = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  selectedChannel.value = data
}

const handleDevicePlay = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  if (data.channel.device_id) {
    monitorGridRef.value?.addDevice({
      ...data.device,
      channel: data.channel,
    })
  } else {
    ElMessage.warning('设备信息不完整')
  }
}

const handleWindowSelect = (data: { deviceId: string; channelId: string } | null) => {
  activeWindow.value = data
}

const handlePtzControl = (direction: string) => {
  if (!activeWindow.value) {
    ElMessage.warning('请先选择视频窗口')
    return
  }
  console.log('云台控制:', direction, activeWindow.value)
}
</script>

<template>
  <div class="monitor-view">
    <div class="monitor-layout">
      <div class="left-panel">
        <DeviceTree 
          @select="handleDeviceSelect"
          @play="handleDevicePlay"
        />
        <PtzControlPanel
          title="云台控制"
          :active-window="activeWindow"
          @control="handlePtzControl"
        />
      </div>
      <div class="monitor-grid-container">
        <MonitorGrid 
          ref="monitorGridRef" 
          @window-select="handleWindowSelect"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.monitor-view {
  height: 100%;
}

.monitor-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
  height: calc(100vh - 180px);
}

.left-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.monitor-grid-container {
  height: 100%;
}
</style>
