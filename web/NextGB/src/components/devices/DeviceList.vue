<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Search, Plus, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { deviceApi } from '@/api'
import type { Device, ChannelInfo } from '@/types/api'

const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const deviceList = ref<Device[]>([])

const dialogVisible = ref(false)
const currentDevice = ref<Device | null>(null)
const channels = ref<ChannelInfo[]>([])
const channelsLoading = ref(false)

const formatDeviceData = (device: Device): Device => {
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
    deviceList.value = (response.data as Device[]).map(formatDeviceData)
  } catch (error) {
    console.error('获取设备列表失败:', error)
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

const showDeviceDetails = async (device: Device) => {
  currentDevice.value = device
  dialogVisible.value = true
  channelsLoading.value = true

  try {
    const response = await deviceApi.getDeviceChannels(device.device_id)
    channels.value = response.data as ChannelInfo[]
  } catch (error) {
    console.error('获取设备通道失败:', error)
    ElMessage.error('获取设备通道失败')
  } finally {
    channelsLoading.value = false
  }
}

const paginatedDevices = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredDevices.value.slice(start, start + pageSize.value)
})

onMounted(() => {
  fetchDevices()
})
</script>

<template>
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
      <el-button type="primary" @click="handleSearch"> 搜索 </el-button>
      <el-button type="success" :loading="loading" @click="fetchDevices">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <el-table v-loading="loading" :data="paginatedDevices" border>
      <template #empty>
        <el-empty :description="searchQuery ? '未找到匹配的设备' : '暂无设备数据'" />
      </template>
      <el-table-column prop="device_id" label="设备ID" />
      <el-table-column prop="source_addr" label="地址" />
      <el-table-column prop="network_type" label="网络类型" />
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 'online' ? 'success' : 'danger'">
            {{ row.status === 'online' ? '在线' : '离线' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="showDeviceDetails(row)"> 详情 </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="filteredDevices.length"
        @current-change="handleCurrentChange"
        layout="total, prev, pager, next"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="`设备详情 - ${currentDevice?.device_id}`"
      width="60%"
    >
      <div class="device-details">
        <div class="device-info">
          <h3>设备信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="设备ID">{{
              currentDevice?.device_id
            }}</el-descriptions-item>
            <el-descriptions-item label="地址">{{
              currentDevice?.source_addr
            }}</el-descriptions-item>
            <el-descriptions-item label="网络类型">{{
              currentDevice?.network_type
            }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentDevice?.status === 'online' ? 'success' : 'danger'">
                {{ currentDevice?.status === 'online' ? '在线' : '离线' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="channel-list">
          <h3>��道列表</h3>
          <el-table v-loading="channelsLoading" :data="channels" border>
            <el-table-column prop="name" label="通道名称" />
            <el-table-column prop="device_id" label="通道ID" />
            <el-table-column prop="manufacturer" label="厂商" />
            <el-table-column prop="status" label="状态">
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
</template>

<style scoped>
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
</style>
