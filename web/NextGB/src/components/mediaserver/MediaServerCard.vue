<script setup lang="ts">
import { Delete, View } from '@element-plus/icons-vue'
import { computed } from 'vue'
import type { MediaServer } from '@/types/api'

const props = defineProps<{
  server: MediaServer
}>()

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
</script>

<template>
  <el-card class="server-card">
    <div class="server-header">
      <img src="@/assets/srs-logo.ico" class="server-icon" />
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
      <p>最新交互时间: {{ new Date().toLocaleString() }}</p>
      <p>端口: {{ server.port }}</p>
    </div>
    <div class="server-footer">
      <el-button-group>
        <el-button type="primary" size="small" :icon="View">查看</el-button>
        <el-button type="danger" size="small" :icon="Delete" @click="handleDelete">删除</el-button>
      </el-button-group>
    </div>
  </el-card>
</template>

<style scoped>
.server-card {
  transition: all 0.3s;
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
</style> 