<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { deviceApi } from '@/api'
import type { Device, ChannelInfo } from '@/types/api'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'

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

const filteredData = computed(() => {
  if (!searchQuery.value.trim()) {
    expandedKeys.value = [] // 清空展开的节点
    return deviceNodes.value
  }

  const query = searchQuery.value.trim().toLowerCase()
  expandedKeys.value = ['root'] // 添加根节点，确保"所有设���"始终展开

  const filteredNodes = deviceNodes.value.filter((node) => {
    // 递归搜索设备和通道
    const searchNode = (item: any): boolean => {
      const isMatch =
        item.label?.toLowerCase().includes(query) || item.device_id?.toLowerCase().includes(query)

      // 如果当前节点匹配
      if (isMatch) {
        if (item.isChannel) {
          // 如果是通道节点匹配，将其父设备节点添加到展开列表中
          const parentDevice = deviceNodes.value.find((device) =>
            device.children?.some((channel) => channel.device_id === item.device_id),
          )
          if (parentDevice) {
            expandedKeys.value.push(parentDevice.device_id)
          }
        } else {
          // 如果是设备节点匹配，将其自身添加到展开列表中
          expandedKeys.value.push(item.device_id)
        }
      }

      // 递归搜索子节点
      if (item.children) {
        const hasMatchingChild = item.children.some(searchNode)
        // 如果子节点匹配，将当前节点ID添加到展开列表中
        if (hasMatchingChild && !expandedKeys.value.includes(item.device_id)) {
          expandedKeys.value.push(item.device_id)
        }
        return isMatch || hasMatchingChild
      }
      return isMatch
    }
    return searchNode(node)
  })

  return filteredNodes
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
        <span class="custom-tree-node">
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
</style>
