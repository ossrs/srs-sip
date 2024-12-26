<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import DeviceTree from '@/components/monitor/DeviceTree.vue'
import MonitorGrid from '@/components/monitor/MonitorGrid.vue'
import type { Device, ChannelInfo } from '@/types/api'

const monitorGridRef = ref()
const selectedChannel = ref<{ device: Device | undefined; channel: ChannelInfo } | null>(null)

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
</script>

<template>
  <div class="monitor-view">
    <div class="monitor-layout">
      <div class="device-tree-container">
        <DeviceTree 
          @select="handleDeviceSelect"
          @play="handleDevicePlay"
        />
      </div>
      <div class="monitor-grid-container">
        <MonitorGrid ref="monitorGridRef" />
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
  grid-template-columns: 240px 1fr;
  gap: 20px;
  height: calc(100vh - 180px);
}

.device-tree-container {
  height: 100%;
}

.monitor-grid-container {
  height: 100%;
}
</style>
