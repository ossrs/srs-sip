<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Search, Plus, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useDevices, useChannels, fetchDevicesAndChannels } from '@/stores/devices'
import type { Device, ChannelInfo } from '@/types/api'

const devices = useDevices()
const allChannels = useChannels()
const deviceList = ref<ExtendedDevice[]>([])

interface ExtendedDevice extends Device {
  channelCount?: number
}

const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

const dialogVisible = ref(false)
const currentDevice = ref<ExtendedDevice | null>(null)
const channels = ref<ChannelInfo[]>([])

const formatDeviceData = (device: Device): ExtendedDevice => {
  return {
    ...device,
    status: device.status || 'offline',
    name: device.name || device.device_id,
    channelCount: allChannels.value.filter(channel => channel.device_id === device.device_id).length,
  }
}

const fetchDevices = async (showError = true) => {
  try {
    loading.value = true
    await fetchDevicesAndChannels()
    deviceList.value = devices.value.map(formatDeviceData)
  } catch (error) {
    console.error('获取设备列表失败:', error)
    if (showError) {
      ElMessage.error('获取设备列表失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}

const filteredDevices = computed(() => {
  if (!searchQuery.value.trim()) return deviceList.value

  const query = searchQuery.value.trim().toLowerCase()
  return deviceList.value.filter((device) => {
    return (
      device.name?.toLowerCase().includes(query) ||
      device.device_id?.toLowerCase().includes(query) ||
      device.source_addr?.toLowerCase().includes(query) ||
      device.network_type?.toLowerCase().includes(query)
    )
  })
})

const handleSearch = () => {
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

const handleRefresh = () => {
  fetchDevices(true)
}

const showDeviceDetails = async (device: ExtendedDevice) => {
  currentDevice.value = device
  dialogVisible.value = true
  channels.value = allChannels.value.filter(channel => channel.device_id === device.device_id)
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'online':
      return 'success'
    case 'offline':
      return 'danger'
    default:
      return 'warning'
  }
}

const getStatusText = (status: string) => {
  switch (status.toLowerCase()) {
    case 'online':
      return '在线'
    case 'offline':
      return '离线'
    default:
      return '未知'
  }
}

const paginatedDevices = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredDevices.value.slice(start, start + pageSize.value)
})

</script>
<template>
  <div class="device-list-view">
    <h1>设备管理</h1>
    <div class="device-list">
      <div class="toolbar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索设备ID、名称、地址或网络类型..."
          class="search-input"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
        <el-button type="success" :loading="loading" @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>

      <el-table v-loading="loading" :data="paginatedDevices" border stripe>
        <template #empty>
          <el-empty :description="searchQuery ? '未找到匹配的设备' : '暂无设备数据'" />
        </template>
        <el-table-column prop="device_id" label="设备ID" min-width="120" show-overflow-tooltip />
        <el-table-column prop="name" label="设备名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="source_addr" label="地址" min-width="120" show-overflow-tooltip />
        <el-table-column prop="network_type" label="网络类型" min-width="100" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="channelCount" label="通道数量" width="100" align="center">
          <template #default="{ row }">
            {{ row.channelCount || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click.stop="showDeviceDetails(row)"> 查看详情 </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="filteredDevices.length"
          @current-change="handleCurrentChange"
          layout="total, prev, pager, next, jumper"
        />
      </div>

      <el-dialog
        v-model="dialogVisible"
        :title="`设备详情 - ${currentDevice?.name || currentDevice?.device_id}`"
        width="70%"
        destroy-on-close
      >
        <div class="device-details">
          <div class="device-info">
            <h3>设备信息</h3>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="设备ID">
                {{ currentDevice?.device_id }}
              </el-descriptions-item>
              <el-descriptions-item label="设备名称">
                {{ currentDevice?.name }}
              </el-descriptions-item>
              <el-descriptions-item label="地址">
                {{ currentDevice?.source_addr }}
              </el-descriptions-item>
              <el-descriptions-item label="网络类型">
                {{ currentDevice?.network_type }}
              </el-descriptions-item>
              <el-descriptions-item label="状态">
                <el-tag :type="currentDevice?.status === 'online' ? 'success' : 'danger'">
                  {{ currentDevice?.status === 'online' ? '在线' : '离线' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="通道数量">
                {{ currentDevice?.channelCount || 0 }}
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <div class="channel-list">
            <h3>通道列表</h3>
            <el-table :data="channels" border stripe>
              <el-table-column type="index" label="序号" width="60" align="center" />
              <el-table-column prop="name" label="通道名称" min-width="120" show-overflow-tooltip />
              <el-table-column
                prop="device_id"
                label="通道ID"
                min-width="120"
                show-overflow-tooltip
              />
              <el-table-column prop="manufacturer" label="厂商" min-width="120" />
              <el-table-column prop="status" label="状态" width="100" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'ON' ? 'success' : 'danger'">
                    {{ row.status === 'ON' ? '在线' : '离线' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<style scoped>
.device-list-view {
  height: 100%;
}

h1 {
  margin-bottom: 20px;
}

.device-list {
  width: 100%;
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  height: calc(100vh - 180px);
  display: flex;
  flex-direction: column;
}

.el-table {
  flex: 1;
  margin: 20px 0;
}

.toolbar {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.search-input {
  width: 300px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.device-details {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.device-info h3,
.channel-list h3 {
  margin-bottom: 15px;
  font-weight: 500;
  color: #303133;
}

.channel-list {
  margin-top: 20px;
}

.channel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}
</style>

