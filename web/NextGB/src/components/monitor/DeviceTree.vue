<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { deviceApi } from '@/api'
import type { Device, ChannelInfo } from '@/types/api'
import { ElMessage } from 'element-plus'

interface DeviceNode {
  device_id: string
  label: string
  children?: DeviceNode[]
  isChannel?: boolean
  channelInfo?: ChannelInfo
}

const devices = ref<Device[]>([])
const deviceNodes = ref<DeviceNode[]>([])
const loading = ref(false)

const formatDeviceData = (device: any): Device => {
  return {
    device_id: device.device_id,
    source_addr: device.source_addr,
    network_type: device.network_type,
    status: 'online',
    name: device.device_id,
  }
}

const fetchDevices = async () => {
  try {
    loading.value = true
    const response = await deviceApi.getDevices()
    devices.value = response.data.map(formatDeviceData)

    const nodes: DeviceNode[] = []
    for (const device of devices.value) {
      try {
        const response = await deviceApi.getDeviceChannels(device.device_id)
        const deviceNode: DeviceNode = {
          device_id: device.device_id,
          label: device.device_id,
          children: response.data.map(channel => ({
            device_id: channel.device_id,
            label: channel.name || channel.device_id,
            isChannel: true,
            channelInfo: channel
          }))
        }
        nodes.push(deviceNode)
      } catch (error) {
        console.error(`获取设备 ${device.device_id} 的通道失败:`, error)
      }
    }
    deviceNodes.value = nodes
  } catch (error) {
    console.error('获取设备列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSelect = (data: DeviceNode) => {
  if (data.isChannel && data.channelInfo) {
    emit('select', {
      device: devices.value.find((d) => d.device_id === data.channelInfo?.parent_id),
      channel: data.channelInfo,
    })
  }
}

const emit = defineEmits<{
  (e: 'select', data: { device: Device | undefined; channel: ChannelInfo }): void
}>()

onMounted(() => {
  fetchDevices()
})
</script>

<template>
  <div class="device-tree">
    <div class="tree-header">
      <h3>设备列表</h3>
    </div>
    <el-tree
      v-loading="loading"
      :data="[
        {
          label: '所有设备',
          children: deviceNodes,
        },
      ]"
      :props="{ children: 'children', label: 'label' }"
      @node-click="handleSelect"
      node-key="device_id"
      highlight-current
    >
      <template #default="{ node, data }">
        <span class="custom-tree-node">
          {{ data.label }}
          <el-tag
            v-if="data.isChannel"
            size="small"
            :type="data.channelInfo?.status === 'ON' ? 'success' : 'danger'"
          >
            {{ data.channelInfo?.status === 'ON' ? '在线' : '离线' }}
          </el-tag>
        </span>
      </template>
    </el-tree>
  </div>
</template>

<style scoped>
.device-tree {
  width: 100%;
  height: 100%;
  background: #fff;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
}

.tree-header {
  padding: 15px;
  border-bottom: 1px solid #eee;
}

.el-tree {
  flex: 1;
  padding: 15px;
  overflow-y: auto;
}

.custom-tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-tree-node__content) {
  height: 32px;
}
</style>
