<script setup lang="ts">
import { ElMessage } from 'element-plus'
import {
  Delete,
  View,
  VideoCamera,
  Microphone,
  Upload,
  Download,
  Connection,
  VideoPlay,
  User,
  Refresh,
} from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import type { MediaServer } from '@/api/mediaserver/types'
import type { StreamInfo } from '@/api/mediaserver/types'
import { createMediaServer } from '@/api/mediaserver/factory'
import zlmLogo from '@/assets/zlm-logo.png'
import srsLogo from '@/assets/srs-logo.ico'

const props = defineProps<{
  server: MediaServer
}>()

const version = ref<string>('获取中...')

const emit = defineEmits<{
  (e: 'delete', server: MediaServer): void
  (e: 'set-default', server: MediaServer): void
}>()

const isDefault = computed(() => props.server.isDefault === 1)

const handleDelete = () => {
  emit('delete', props.server)
}

const handleSetDefault = () => {
  emit('set-default', props.server)
}

const dialogVisible = ref(false)
const streams = ref<StreamInfo[]>([])
const loading = ref(false)

const expandedRowKeys = ref<string[]>([])

const handleRowExpand = (row: StreamInfo) => {
  const index = expandedRowKeys.value.indexOf(row.id)
  if (index > -1) {
    expandedRowKeys.value.splice(index, 1)
  } else {
    expandedRowKeys.value.push(row.id)
  }
}

const withLoading = async (operation: () => Promise<void>) => {
  loading.value = true
  try {
    await operation()
  } catch (error) {
    console.error('获取流信息失败:', error)
    ElMessage.error('获取流信息失败')
  } finally {
    loading.value = false
  }
}

const handleView = async () => {
  dialogVisible.value = true
  await withLoading(fetchStreamInfo)
}

const fetchStreamInfo = async () => {
  const mediaServer = createMediaServer(props.server)
  const streamInfos = await mediaServer.getStreamInfo()

  let clientsMap: Record<string, any[]> = {}
  
  // 根据服务器类型使用不同的获取客户端信息逻辑
  if (props.server.type === 'ZLM') {
    // ZLM需要为每个流单独获取客户端信息
    await Promise.all(
      streamInfos.map(async (stream) => {
        const clients = await mediaServer.getClientInfo({ stream_id: stream.id })
        clientsMap[stream.id] = clients
      })
    )
  } else {
    // SRS可以一次性获取所有客户端信息
    const clients = await mediaServer.getClientInfo()
    clientsMap = clients.reduce(
      (acc: Record<string, any[]>, client) => {
        if (!acc[client.stream]) {
          acc[client.stream] = []
        }
        acc[client.stream].push(client)
        return acc
      },
      {}
    )
  }

  // 将客户端信息关联到对应的流
  streams.value = streamInfos.map((stream: StreamInfo) => ({
    ...stream,
    clients_info: clientsMap[stream.id] || [],
  }))
}

const refreshStreams = () => withLoading(fetchStreamInfo)

onMounted(async () => {
  try {
    const mediaServer = createMediaServer(props.server)
    const versionInfo = await mediaServer.getVersion()
    version.value = versionInfo.version
  } catch (error) {
    version.value = '获取失败'
    console.error('获取版本信息失败:', error)
  }
})

const formatBytes = (bytes?: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

const copyUrl = (url: string) => {
  navigator.clipboard.writeText(url)
  ElMessage.success('URL已复制到剪贴板')
}

const formatDuration = (ms: number) => {
  console.log('Duration input:', ms, typeof ms)

  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    return `${hours}小时${minutes % 60}分`
  } else if (minutes > 0) {
    return `${minutes}分${seconds % 60}秒`
  } else {
    return `${seconds}秒`
  }
}

const handleDisconnect = async (stream: StreamInfo) => {
  // TODO: 实现断开流的API调用
  ElMessage.warning('断开流功能即将实现')
}
</script>

<template>
  <el-card class="server-card" :class="{ 'is-default': isDefault }">
    <div v-if="isDefault" class="default-ribbon"></div>
    <div class="server-header">
      <img :src="server.type === 'ZLM' ? zlmLogo : srsLogo" class="server-icon" />
      <div class="server-info">
        <h3>
          {{ server.name }}
          <el-tag size="small" type="info" class="type-tag">{{ server.type }}</el-tag>
        </h3>
        <div class="server-ip">{{ server.ip }}</div>
      </div>
      <div class="status-tags">
        <el-tag :type="server.status === 1 ? 'success' : 'danger'" class="status-tag">
          {{ server.status === 1 ? '在线' : '离线' }}
        </el-tag>
        <el-tag
          :type="isDefault ? 'warning' : 'info'"
          class="default-tag"
          @click="handleSetDefault"
          style="cursor: pointer"
        >
          {{ isDefault ? '默认节点' : '设为默认' }}
        </el-tag>
      </div>
    </div>
    <div class="server-body">
      <p>版本: {{ version }}</p>
    </div>
    <div class="server-footer">
      <el-button-group>
        <el-button type="primary" size="small" :icon="View" @click="handleView">查看</el-button>
        <el-button type="danger" size="small" :icon="Delete" @click="handleDelete">删除</el-button>
      </el-button-group>
    </div>
  </el-card>

  <el-dialog
    v-model="dialogVisible"
    :title="`${server.name} - 流信息`"
    width="90%"
    class="stream-dialog"
    destroy-on-close
  >
    <div class="stream-dashboard">
      <div class="dashboard-item">
        <div class="dashboard-icon">
          <el-icon><Connection /></el-icon>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-value">{{ streams.length }}</div>
          <div class="dashboard-label">总流数</div>
        </div>
      </div>
      <div class="dashboard-divider"></div>
      <div class="dashboard-item">
        <div class="dashboard-icon active">
          <el-icon><VideoPlay /></el-icon>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-value success">{{ streams.filter((s) => s.active).length }}</div>
          <div class="dashboard-label">活跃流数</div>
        </div>
      </div>
      <div class="dashboard-divider"></div>
      <div class="dashboard-item">
        <div class="dashboard-icon primary">
          <el-icon><User /></el-icon>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-value primary">
            {{ streams.reduce((sum, s) => sum + s.clients, 0) }}
          </div>
          <div class="dashboard-label">总客户端数</div>
        </div>
      </div>
    </div>

    <el-table
      v-loading="loading"
      :data="streams"
      style="width: 100%"
      border
      stripe
      class="stream-table"
      :empty-text="loading ? '加载中...' : '暂无流数据'"
      :expand-row-keys="expandedRowKeys"
      row-key="id"
      @row-click="(row) => handleRowExpand(row)"
    >
      <el-table-column prop="name" label="流名称" min-width="120" />
      <el-table-column label="URL" min-width="200" show-overflow-tooltip>
        <template #default="{ row }">
          <el-link type="primary" :underline="false" @click="copyUrl(row.url)">
            {{ row.url }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="clients" label="客户端数" width="100" align="center" />
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.active ? 'success' : 'danger'" size="small">
            {{ row.active ? '活跃' : '断开' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="编码信息" min-width="320">
        <template #default="{ row }">
          <div class="codec-info">
            <template v-if="row.video">
              <el-tooltip content="视频编码信息" placement="top">
                <div class="info-chip video">
                  <el-icon><VideoCamera /></el-icon>
                  <span class="codec">{{ row.video.codec }}</span>
                  <span class="detail">{{ row.video.width }}x{{ row.video.height }}</span>
                </div>
              </el-tooltip>
            </template>
            <template v-if="row.audio">
              <el-tooltip content="音频编码信息" placement="top">
                <div class="info-chip audio">
                  <el-icon><Microphone /></el-icon>
                  <span class="codec">{{ row.audio.codec }}</span>
                  <span class="detail">{{ row.audio.sampleRate }}Hz</span>
                </div>
              </el-tooltip>
            </template>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="传输数据" width="300">
        <template #default="{ row }">
          <div class="transfer-info">
            <el-tooltip :content="server.type === 'ZLM' ? '下行速率' : '累计下行流量'" placement="top">
              <div class="info-chip download">
                <el-icon><Download /></el-icon>
                <span class="value">{{ formatBytes(row.send_bytes) }}</span>
                <span v-if="server.type === 'ZLM'" class="rate-unit">/s</span>
              </div>
            </el-tooltip>
            <el-tooltip :content="server.type === 'ZLM' ? '上行速率' : '累计上行流量'" placement="top">
              <div class="info-chip upload">
                <el-icon><Upload /></el-icon>
                <span class="value">{{ formatBytes(row.recv_bytes) }}</span>
                <span v-if="server.type === 'ZLM'" class="rate-unit">/s</span>
              </div>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" align="center">
        <template #default="{ row }">
          <el-button
            type="danger"
            size="small"
            :disabled="!row.active"
            @click.stop="handleDisconnect(row)"
          >
            断开
          </el-button>
        </template>
      </el-table-column>
      <el-table-column type="expand" :width="0">
        <template #default="{ row }">
          <div class="expanded-details">
            <el-table
              :data="row.clients_info || []"
              border
              stripe
              size="small"
              class="client-table"
              :show-header="false"
              style="width: 600px"
            >
              <el-table-column prop="id" width="100">
                <template #default="{ row: client }">
                  <el-tooltip content="客户端ID" placement="top">
                    <span>{{ client.id }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column prop="ip" width="110">
                <template #default="{ row: client }">
                  <el-tooltip content="IP地址" placement="top">
                    <span>{{ client.ip }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column prop="type" width="60">
                <template #default="{ row: client }">
                  <el-tooltip content="类型" placement="top">
                    <span>{{ client.type }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column prop="url" min-width="150">
                <template #default="{ row: client }">
                  <el-tooltip content="URL" placement="top">
                    <el-link type="primary" :underline="false" @click.stop="copyUrl(client.url)">
                      {{ client.url }}
                    </el-link>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column prop="alive" width="80" align="right">
                <template #default="{ row: client }">
                  <el-tooltip content="存活时间" placement="top">
                    <span>{{ formatDuration(client.alive) }}</span>
                  </el-tooltip>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="refreshStreams">
          <el-icon><Refresh /></el-icon>刷新
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped>
.server-card {
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
}

.is-default {
  border: 1px solid #e6a23c;
}

.default-ribbon {
  position: absolute;
  top: 0;
  left: 0;
  width: 100px;
  height: 100px;
  overflow: hidden;
  pointer-events: none;
}

.default-ribbon::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 150%;
  height: 24px;
  background: #e6a23c;
  transform: rotate(-45deg) translateX(-50%);
  transform-origin: 0 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.server-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.server-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.server-icon {
  width: 40px;
  height: 40px;
  object-fit: cover;
  object-position: left;
}

.server-info {
  flex: 1;
}

.server-info h3 {
  margin: 0 0 5px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.server-ip {
  color: #909399;
  font-size: 14px;
}

.server-body {
  color: #666;
  font-size: 14px;
  line-height: 1.8;
}

.server-body p {
  margin: 5px 0;
}

.server-footer {
  margin-top: 15px;
  display: flex;
  justify-content: flex-end;
}

.status-tags {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.status-tag,
.default-tag {
  white-space: nowrap;
}

.default-tag {
  transition: all 0.3s;
}

.default-tag:hover {
  opacity: 0.8;
}

.type-tag {
  margin-left: 8px;
  font-weight: normal;
  vertical-align: middle;
}

:deep(.el-radio-group) {
  width: 100%;
}

.stream-dialog :deep(.el-dialog__header) {
  padding: 20px 24px;
  margin-right: 0;
  border-bottom: 1px solid #dcdfe6;
}

.stream-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.stream-dialog :deep(.el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #dcdfe6;
}

.stream-dashboard {
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin-bottom: 16px;
  padding: 16px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.dashboard-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 24px;
}

.dashboard-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #f5f7fa;
  color: #909399;
  transition: all 0.3s ease;
}

.dashboard-icon.active {
  background: #f0f9eb;
  color: #67c23a;
}

.dashboard-icon.primary {
  background: #ecf5ff;
  color: #409eff;
}

.dashboard-icon :deep(.el-icon) {
  font-size: 20px;
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dashboard-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  line-height: 1.2;
}

.dashboard-value.success {
  color: #67c23a;
}

.dashboard-value.primary {
  color: #409eff;
}

.dashboard-label {
  font-size: 14px;
  color: #909399;
}

.dashboard-divider {
  width: 1px;
  height: 60px;
  background: linear-gradient(180deg, transparent, #dcdfe6 50%, transparent);
}

.stream-table {
  margin-top: 16px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.stream-table :deep(.el-table__header-wrapper) {
  border-radius: 8px 8px 0 0;
}

.stream-table :deep(.el-table__header) th {
  background-color: #f5f7fa;
  font-weight: 600;
  height: 40px;
  padding: 4px 0;
}

.stream-table :deep(.el-table__row) {
  transition: all 0.3s ease;
  height: 48px;
}

.stream-table :deep(.el-table__cell) {
  padding: 4px 8px;
}

.stream-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.codec-info,
.transfer-info {
  display: flex;
  align-items: center;
  gap: 6px;
}

.info-chip {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: 12px;
  font-size: 12px;
  white-space: nowrap;
  cursor: default;
  transition: all 0.3s ease;
}

.info-chip .el-icon {
  margin-right: 4px;
  font-size: 12px;
}

.info-chip.video {
  background: linear-gradient(45deg, #409eff22, #409eff11);
  color: #409eff;
}

.info-chip.audio {
  background: linear-gradient(45deg, #67c23a22, #67c23a11);
  color: #67c23a;
}

.info-chip.download {
  background: linear-gradient(45deg, #409eff22, #409eff11);
  color: #409eff;
}

.info-chip.upload {
  background: linear-gradient(45deg, #e6a23d22, #e6a23d11);
  color: #e6a23d;
}

.info-chip .unit {
  font-size: 10px;
  opacity: 0.8;
  margin-left: 2px;
}

.info-chip .rate-unit {
  font-size: 10px;
  opacity: 0.8;
  margin-left: 1px;
}

.info-chip:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.info-chip .codec {
  font-weight: 500;
  margin-right: 4px;
}

.info-chip .detail,
.info-chip .value {
  opacity: 0.9;
}

/* 优化表格样式 */
.stream-table :deep(.el-table__row) {
  height: 60px;
}

.stream-table :deep(.el-table__cell) {
  padding: 8px 12px;
}

.expanded-details {
  padding: 4px 8px 4px 48px;
  background: #f8f9fa;
}

.client-table {
  --el-table-border-color: #e4e7ed;
  --el-table-row-hover-bg-color: #ecf5ff;
}

.client-table :deep(.el-table__row) {
  height: 28px;
}

.client-table :deep(.el-table__cell) {
  padding: 2px 4px;
}

.client-table :deep(.cell) {
  color: #606266;
  font-size: 12px;
  line-height: 1.4;
}

.client-table :deep(.el-link) {
  font-size: 12px;
}

.stream-table :deep(.el-table__expand-column) {
  display: none;
}

.stream-table :deep(.el-table__row) {
  cursor: pointer;
}
</style>
