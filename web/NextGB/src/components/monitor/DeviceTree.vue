<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { deviceApi } from '@/api'
import type { Device, ChannelInfo } from '@/types/api'
import { ElMessage } from 'element-plus'
import { Search, Refresh, List, Grid } from '@element-plus/icons-vue'

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
const searchQuery = ref('')
const expandedKeys = ref<string[]>([])

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
    devices.value = (response.data || []).map(formatDeviceData)

    const nodes: DeviceNode[] = []
    for (const device of devices.value) {
      try {
        const response = await deviceApi.getDeviceChannels(device.device_id)
        const deviceNode: DeviceNode = {
          device_id: device.device_id,
          label: device.device_id,
          children: (response.data || []).map((channel: ChannelInfo) => ({
            device_id: channel.device_id,
            label: channel.name || channel.device_id,
            isChannel: true,
            channelInfo: channel,
          })),
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

const emit = defineEmits<{
  (e: 'select', data: { device: Device | undefined; channel: ChannelInfo }): void
  (e: 'play', data: { device: Device | undefined; channel: ChannelInfo }): void
}>()

const handleSelect = (data: DeviceNode) => {
  if (data.isChannel && data.channelInfo) {
    emit('select', {
      device: devices.value.find((d) => d.device_id === data.channelInfo?.parent_id),
      channel: data.channelInfo,
    })
  }
}

const handleNodeDbClick = (data: DeviceNode) => {
  if (data.isChannel && data.channelInfo) {
    emit('play', {
      device: devices.value.find((d) => d.device_id === data.channelInfo?.parent_id),
      channel: data.channelInfo,
    })
  }
}

const viewMode = ref<'tree' | 'list'>('tree')

const filteredData = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()

  if (viewMode.value === 'list') {
    const allChannels = deviceNodes.value.flatMap(node => 
      (node.children || []).map(channel => ({
        ...channel,
        parentDeviceId: node.device_id
      }))
    )
    
    if (!query) return allChannels

    return allChannels.filter(channel => 
      channel.label.toLowerCase().includes(query) ||
      channel.device_id.toLowerCase().includes(query)
    )
  }

  if (!query) {
    expandedKeys.value = []
    return deviceNodes.value
  }

  expandedKeys.value = ['root']

  return deviceNodes.value.filter((node) => {
    const searchNode = (item: any): boolean => {
      const isMatch =
        item.label?.toLowerCase().includes(query) || item.device_id?.toLowerCase().includes(query)

      if (isMatch) {
        if (item.isChannel) {
          const parentDevice = deviceNodes.value.find((device) =>
            device.children?.some((channel) => channel.device_id === item.device_id),
          )
          if (parentDevice) {
            expandedKeys.value.push(parentDevice.device_id)
          }
        } else {
          expandedKeys.value.push(item.device_id)
        }
      }

      if (item.children) {
        const hasMatchingChild = item.children.some(searchNode)
        if (hasMatchingChild && !expandedKeys.value.includes(item.device_id)) {
          expandedKeys.value.push(item.device_id)
        }
        return isMatch || hasMatchingChild
      }
      return isMatch
    }
    return searchNode(node)
  })
})

const tooltipRef = ref()

const refreshDevices = () => {
  fetchDevices()
  tooltipRef.value?.hide()
}

onMounted(() => {
  fetchDevices()
})
</script>

<template>
  <div class="device-tree">
    <div class="search-box">
      <div class="search-wrapper">
        <el-input v-model="searchQuery" placeholder="搜索设备或通道..." clearable size="small">
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-tooltip :content="viewMode === 'tree' ? '切换到列表视图' : '切换到树形视图'" placement="top">
          <el-button
            class="view-mode-btn"
            :icon="viewMode === 'tree' ? List : Grid"
            size="small"
            @click="viewMode = viewMode === 'tree' ? 'list' : 'tree'"
          />
        </el-tooltip>
        <el-tooltip ref="tooltipRef" content="刷新设备列表" placement="top">
          <el-button
            type="primary"
            :icon="Refresh"
            size="small"
            :loading="loading"
            @click="refreshDevices"
          />
        </el-tooltip>
      </div>
    </div>

    <el-tree
      v-if="viewMode === 'tree'"
      v-loading="loading"
      :data="[
        {
          label: '所有设备',
          device_id: 'root',
          children: filteredData,
        },
      ]"
      :props="{ children: 'children', label: 'label' }"
      @node-click="handleSelect"
      node-key="device_id"
      highlight-current
      :expanded-keys="expandedKeys"
    >
      <template #default="{ node, data }">
        <span class="custom-tree-node" @dblclick.stop="handleNodeDbClick(data)">
          <span :class="data.isChannel ? 'channel-label' : 'device-label'">
            {{ data.label }}
          </span>
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

    <div v-else class="channel-list" v-loading="loading">
      <div
        v-for="channel in filteredData"
        :key="channel.device_id"
        class="channel-list-item"
        @click="handleSelect(channel)"
        @dblclick="handleNodeDbClick(channel)"
      >
        <span class="channel-label">{{ channel.label }}</span>
        <el-tag
          size="small"
          :type="channel.channelInfo?.status === 'ON' ? 'success' : 'danger'"
        >
          {{ channel.channelInfo?.status === 'ON' ? '在线' : '离线' }}
        </el-tag>
      </div>
    </div>
  </div>
</template>

<style scoped>
.device-tree {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  border-radius: 4px;
  padding: 16px;
}

.search-box {
  margin-bottom: 16px;
}

.search-wrapper {
  display: flex;
  gap: 8px;

  .el-input {
    flex: 1;
  }

  .view-mode-btn {
    color: var(--el-text-color-regular);
    border-color: transparent;
    background-color: transparent;
    padding: 5px 8px;

    &:hover {
      color: var(--el-color-primary);
      background-color: var(--el-fill-color-light);
    }

    &:focus {
      outline: none;
    }
  }

  :deep(.el-input__wrapper) {
    box-shadow: 0 0 0 1px #dcdfe6 inset;

    &:hover {
      box-shadow: 0 0 0 1px #c0c4cc inset;
    }

    &.is-focus {
      box-shadow: 0 0 0 1px #409eff inset;
    }
  }
}

.tree-header {
  padding: 15px;
  border-bottom: 1px solid #eee;
}

.el-tree {
  flex: 1;
  padding: 0;
  overflow-y: auto;

  :deep(.el-tree-node) {
    &.is-expanded > .el-tree-node__children {
      padding-left: 20px;
    }

    &::before {
      content: '';
      height: 100%;
      width: 1px;
      position: absolute;
      left: -12px;
      top: -4px;
      border-left: 1px dotted #c0c4cc;
    }

    &:last-child::before {
      height: 20px;
    }
  }

  :deep(.el-tree-node__content) {
    height: 32px;
    padding-left: 8px !important;

    &:hover {
      background-color: #f5f7fa;
    }

    &.is-current {
      background-color: #ecf5ff;
      color: #409eff;
    }
  }

  :deep(.el-tree-node__expand-icon) {
    font-size: 16px;
    color: #909399;

    &.expanded {
      transform: rotate(90deg);
    }

    &.is-leaf {
      color: transparent;
    }
  }
}

.custom-tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 0 4px;
  user-select: none;

  .device-label {
    font-weight: 500;
    color: #303133;
  }

  .channel-label {
    color: #606266;
    font-size: 13px;
  }

  .el-tag {
    margin-left: auto;
  }
}

:deep(.el-tree-node__content) {
  user-select: none;
}

.channel-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.channel-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  user-select: none;

  &:hover {
    background-color: #f5f7fa;
  }

  .channel-label {
    color: #606266;
    font-size: 13px;
  }
}

.search-wrapper {
  .el-button-group {
    margin-right: 8px;
  }
}
</style>
