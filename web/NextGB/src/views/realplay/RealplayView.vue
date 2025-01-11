<script setup lang="ts">
import { ref, onActivated, onDeactivated } from 'vue'
import { ElMessage } from 'element-plus'
import { FullScreen, Setting, Delete } from '@element-plus/icons-vue'
import DeviceTree from './DeviceTree.vue'
import MonitorGrid from '@/components/monitor/MonitorGrid.vue'
import PtzControlPanel from '@/views/realplay/PtzControlPanel.vue'
import type { Device, ChannelInfo } from '@/api/types'
import type { LayoutConfig } from '@/types/layout'

// 所有可用的布局键
type LayoutKey = '1' | '4' | '9' | '16'
type LayoutConfigs = Record<LayoutKey, LayoutConfig>

// 布局配置
const layouts: LayoutConfigs = {
  '1': { cols: 1, rows: 1, size: 1, label: '单屏' },
  '4': { cols: 2, rows: 2, size: 4, label: '四分屏' },
  '9': { cols: 3, rows: 3, size: 9, label: '九分屏' },
  '16': { cols: 4, rows: 4, size: 16, label: '十六分屏' },
} as const

const monitorGridRef = ref()
const selectedChannel = ref<{ device: Device | undefined; channel: ChannelInfo } | null>(null)
const activeWindow = ref<{ deviceId: string; channelId: string } | null>(null)
const currentLayout = ref<LayoutKey>('9')
const showSettings = ref(false)
const defaultMuted = ref(true)

const handleDeviceSelect = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  selectedChannel.value = data
}

const handleDevicePlay = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  if (data.channel.device_id) {
    monitorGridRef.value?.play({
      ...data.device,
      channel: data.channel,
      play_type: 0,
      start_time: 0,
      end_time: 0,
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

const clearAll = () => {
  monitorGridRef.value?.clear()
}

const toggleGridFullscreen = async () => {
  const gridContainer = document.querySelector('.monitor-grid') as HTMLElement
  if (!gridContainer) {
    console.error('未找到视频网格容器')
    return
  }

  try {
    if (!document.fullscreenElement) {
      await gridContainer.requestFullscreen()
    } else {
      await document.exitFullscreen()
    }
  } catch (err) {
    console.error('全屏切换失败:', err)
    ElMessage.error('全屏切换失败')
  }
}

// 添加激活/停用处理
onActivated(() => {
  console.log('MonitorView activated')
  // 如果需要在重新激活时执行某些操作，可以在这里添加
})

onDeactivated(() => {
  console.log('MonitorView deactivated')
  // 组件被缓存，不需要清理视频资源
})

// 组件名称（用于 keep-alive）
defineOptions({
  name: 'RealplayView',
})
</script>

<template>
  <div class="monitor-view">
    <div class="monitor-layout">
      <div class="left-panel">
        <DeviceTree @select="handleDeviceSelect" @play="handleDevicePlay" />
        <PtzControlPanel
          title="云台控制"
          :active-window="activeWindow"
          @control="handlePtzControl"
        />
      </div>
      <div class="monitor-grid-container">
        <div class="grid-toolbar">
          <div class="layout-controls">
            <el-radio-group v-model="currentLayout" size="small">
              <el-radio-button v-for="(layout, key) in layouts" :key="key" :value="key">
                {{ layout.label }}
              </el-radio-button>
            </el-radio-group>
          </div>
          <div class="toolbar-actions">
            <el-button-group>
              <el-button size="small" @click="showSettings = true" :title="'设置'">
                <el-icon><Setting /></el-icon>
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="clearAll"
                :title="'清空所有设备'"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
              <el-button size="small" @click="toggleGridFullscreen" :title="'全屏'">
                <el-icon><FullScreen /></el-icon>
              </el-button>
            </el-button-group>
          </div>
        </div>
        <MonitorGrid
          ref="monitorGridRef"
          v-model="currentLayout"
          :layouts="layouts"
          :default-muted="defaultMuted"
          :show-border="true"
          @window-select="handleWindowSelect"
        />
      </div>
    </div>
  </div>

  <!-- 设置对话框 -->
  <el-dialog v-model="showSettings" title="设置" width="400px" destroy-on-close>
    <el-form label-width="120px">
      <el-form-item label="默认静音">
        <el-switch v-model="defaultMuted" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="showSettings = false">取消</el-button>
        <el-button type="primary" @click="showSettings = false">确定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped>
.monitor-view {
  height: 100%;
}

.monitor-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
  height: 100%;
}

.left-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.monitor-grid-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.grid-toolbar {
  padding: 12px 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--el-bg-color);
  border-radius: 4px 4px 0 0;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  gap: 10px;
}

:deep(.el-button-group .el-button--small) {
  padding: 5px 11px;
}

:deep(.el-radio-group .el-radio-button__inner) {
  padding: 5px 15px;
}

:deep(.el-radio-group) {
  --el-button-bg-color: var(--el-fill-color-blank);
  --el-button-hover-bg-color: var(--el-fill-color);
  --el-button-active-bg-color: var(--el-color-primary);
  --el-button-text-color: var(--el-text-color-regular);
  --el-button-hover-text-color: var(--el-text-color-primary);
  --el-button-active-text-color: #fff;
  --el-button-border-color: var(--el-border-color);
  --el-button-hover-border-color: var(--el-border-color-hover);
}
</style>
