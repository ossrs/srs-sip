<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Device, ChannelInfo } from '@/api/types'
import { ElMessage } from 'element-plus'
import {
  useDevices,
  useChannels,
  useDevicesLoading,
  fetchDevicesAndChannels,
} from '@/stores/devices'
import SearchBox from '@/components/common/SearchBox.vue'

interface DeviceNode {
  device_id: string
  label: string
  children?: DeviceNode[]
  isChannel?: boolean
  channelInfo?: ChannelInfo
}

const devices = useDevices()
const channels = useChannels()
const loading = useDevicesLoading()
const searchQuery = ref('')
const expandedKeys = ref<string[]>([])
const selectedChannels = ref<Set<string>>(new Set())
const viewMode = ref<'tree' | 'list'>('list')

const deviceNodes = computed(() => {
  const nodes: DeviceNode[] = []
  for (const device of devices.value) {
    const deviceChannels = channels.value.filter(
      (channel: ChannelInfo) => channel.parent_id === device.device_id,
    )
    const deviceNode: DeviceNode = {
      device_id: device.device_id,
      label: device.name || '未命名',
      children: deviceChannels.map((channel: ChannelInfo) => ({
        device_id: channel.device_id,
        label: `${channel.name}`,
        isChannel: true,
        channelInfo: channel,
      })),
    }
    nodes.push(deviceNode)
  }
  return nodes
})

const refreshDevices = async () => {
  try {
    await fetchDevicesAndChannels()
  } catch (error) {
    ElMessage.error('刷新设备列表失败')
  }
  tooltipRef.value?.hide()
}

const emit = defineEmits<{
  (
    e: 'update:selectedChannels',
    channels: { device: Device | undefined; channel: ChannelInfo }[],
  ): void
}>()

const handleChannelSelect = (channelId: string, checked: boolean) => {
  if (checked) {
    // 清除之前选中的所有通道
    selectedChannels.value.clear()
    // 只添加当前选中的通道
    selectedChannels.value.add(channelId)
  } else {
    selectedChannels.value.delete(channelId)
  }

  // 发送选中的通道信息
  const selectedChannelInfos = Array.from(selectedChannels.value)
    .map((id) => {
      const channel = channels.value.find((c: ChannelInfo) => c.device_id === id)
      if (channel) {
        return {
          device: devices.value.find((d: Device) => d.device_id === channel.parent_id),
          channel,
        }
      }
      return null
    })
    .filter((info) => info !== null) as { device: Device | undefined; channel: ChannelInfo }[]

  emit('update:selectedChannels', selectedChannelInfos)
}

const filteredData = computed(() => {
  const nodes = deviceNodes.value
  const query = searchQuery.value.trim().toLowerCase()

  if (viewMode.value === 'list') {
    const allChannels = nodes.flatMap((node) =>
      (node.children || []).map((channel) => ({
        ...channel,
        parentDeviceId: node.device_id,
      })),
    )

    if (!query) {
      return [
        {
          label: '所有通道',
          device_id: 'root',
          children: allChannels,
        },
      ]
    }

    const filteredChannels = allChannels.filter(
      (channel) =>
        channel.label.toLowerCase().includes(query) ||
        channel.device_id.toLowerCase().includes(query),
    )

    return [
      {
        label: '所有通道',
        device_id: 'root',
        children: filteredChannels,
      },
    ]
  }

  if (!query) {
    expandedKeys.value = []
    return nodes
  }

  expandedKeys.value = ['root']

  return nodes.filter((node) => {
    const searchNode = (item: any): boolean => {
      const isMatch =
        item.label?.toLowerCase().includes(query) || item.device_id?.toLowerCase().includes(query)

      if (isMatch) {
        if (item.isChannel) {
          const parentDevice = nodes.find((device) =>
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
</script>

<template>
  <div class="device-tree">
    <SearchBox
      v-model:searchQuery="searchQuery"
      v-model:viewMode="viewMode"
      :loading="loading"
      :show-view-mode-switch="true"
      @refresh="refreshDevices"
    />

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
      node-key="device_id"
      :expanded-keys="expandedKeys"
      default-expand-all
    >
      <template #default="{ node, data }">
        <span class="custom-tree-node">
          <el-checkbox
            v-if="data.isChannel"
            :model-value="selectedChannels.has(data.device_id)"
            @update:model-value="
              (checked) => handleChannelSelect(data.device_id, checked as boolean)
            "
          />
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
      <el-tree
        :data="filteredData"
        :props="{ children: 'children', label: 'label' }"
        node-key="device_id"
        :default-expanded-keys="['root']"
      >
        <template #default="{ node, data }">
          <span class="custom-tree-node">
            <el-checkbox
              v-if="data.isChannel"
              :model-value="selectedChannels.has(data.device_id)"
              @update:model-value="
                (checked) => handleChannelSelect(data.device_id, checked as boolean)
              "
            />
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
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
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

  .action-buttons {
    display: flex;
    gap: 0;
    margin-left: auto;

    .action-btn {
      color: var(--el-text-color-regular);
      border-color: transparent;
      background-color: transparent;
      padding: 5px 8px;
      height: 32px;
      width: 32px;
      border-radius: 4px;

      &:hover {
        color: var(--el-color-primary);
        background-color: var(--el-fill-color-light);
      }

      &:focus {
        outline: none;
      }

      :deep(.el-icon) {
        font-size: 16px;
      }

      &.refresh-btn {
        color: var(--el-color-primary);

        &:hover {
          background-color: var(--el-color-primary-light-9);
        }
      }
    }
  }
}

.el-tree {
  flex: 1;
  padding: 0;
  overflow-y: auto;
  border-radius: 4px;

  :deep(.el-tree-node__content) {
    height: 36px;
    padding-left: 8px !important;
    border-radius: 4px;
    transition: all 0.2s ease;

    &:hover {
      background-color: var(--el-fill-color-light);
    }
  }

  :deep(.el-tree-node__expand-icon) {
    font-size: 16px;
    color: var(--el-text-color-secondary);
    transition: transform 0.2s ease;

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
    color: var(--el-text-color-primary);
    font-size: 14px;
  }

  .channel-label {
    color: var(--el-text-color-regular);
    font-size: 13px;
  }

  .el-tag {
    margin-left: auto;
    transition: all 0.2s ease;
  }

  :deep(.el-checkbox) {
    margin-right: 4px;
  }
}

.channel-list {
  flex: 1;
  overflow-y: auto;
  padding: 2px;
  border-radius: 4px;

  :deep(.el-tree) {
    background: transparent;

    .el-tree-node__content {
      height: 36px;
      padding-left: 8px !important;
      border-radius: 4px;
      transition: all 0.2s ease;

      &:hover {
        background-color: var(--el-fill-color-light);
      }

      &.is-current {
        background-color: var(--el-color-primary-light-9);
        color: var(--el-color-primary);
      }
    }

    .el-tree-node__expand-icon {
      font-size: 16px;
      color: var(--el-text-color-secondary);
      transition: transform 0.2s ease;

      &.expanded {
        transform: rotate(90deg);
      }

      &.is-leaf {
        color: transparent;
      }
    }
  }
}
</style>
