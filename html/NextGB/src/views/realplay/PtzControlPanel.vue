<script setup lang="ts">
import { ref, computed } from 'vue'
import { ArrowRight, VideoCamera } from '@element-plus/icons-vue'
import {
  ArrowUp,
  ArrowDown,
  ArrowLeft,
  ArrowRight as ArrowRightControl,
  TopLeft,
  TopRight,
  BottomLeft,
  BottomRight,
} from '@element-plus/icons-vue'
import { deviceApi } from '@/api'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  activeWindow?: {
    deviceId: string
    channelId: string
  } | null
  title?: string
}>()

const emit = defineEmits<{
  control: [direction: string]
}>()

const isCollapsed = ref(false)
const speed = ref(5)

const handlePtzStart = async (direction: string) => {
  if (!props.activeWindow) {
    ElMessage.warning('请先选择一个视频窗口')
    return
  }

  try {
    await deviceApi.controlPTZ({
      device_id: props.activeWindow.deviceId,
      channel_id: props.activeWindow.channelId,
      ptz: direction,
      speed: speed.value.toString(),
    })
    emit('control', direction)
  } catch (error) {
    console.error('PTZ control failed:', error)
  }
}

const handlePtzStop = async () => {
  if (!props.activeWindow) return

  try {
    await deviceApi.controlPTZ({
      device_id: props.activeWindow.deviceId,
      channel_id: props.activeWindow.channelId,
      ptz: 'stop',
      speed: speed.value.toString(),
    })
    emit('control', 'stop')
  } catch (error) {
    console.error('PTZ stop failed:', error)
  }
}

const isDisabled = computed(() => !props.activeWindow)
</script>

<template>
  <div class="ptz-control-panel" :class="{ collapsed: isCollapsed }">
    <div class="panel-header" @click="isCollapsed = !isCollapsed">
      <div class="header-title">
        <el-icon class="collapse-arrow" :class="{ collapsed: isCollapsed }">
          <ArrowRight />
        </el-icon>
        <el-icon class="title-icon"><VideoCamera /></el-icon>
        <span>{{ title || '云台控制' }}</span>
      </div>
    </div>
    <div class="panel-content">
      <div class="control-form">
        <div class="ptz-controls">
          <div class="direction-controls">
            <div class="direction-pad">
              <el-button
                class="direction-btn up"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('up')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><ArrowUp /></el-icon>
              </el-button>
              <el-button
                class="direction-btn right"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('right')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><ArrowRightControl /></el-icon>
              </el-button>
              <el-button
                class="direction-btn down"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('down')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <el-button
                class="direction-btn left"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('left')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><ArrowLeft /></el-icon>
              </el-button>
              <el-button
                class="direction-btn up-left"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('upleft')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><TopLeft /></el-icon>
              </el-button>
              <el-button
                class="direction-btn up-right"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('upright')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><TopRight /></el-icon>
              </el-button>
              <el-button
                class="direction-btn down-left"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('downleft')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><BottomLeft /></el-icon>
              </el-button>
              <el-button
                class="direction-btn down-right"
                :disabled="isDisabled"
                @mousedown="handlePtzStart('downright')"
                @mouseup="handlePtzStop"
                @mouseleave="handlePtzStop"
              >
                <el-icon><BottomRight /></el-icon>
              </el-button>
              <div class="direction-center"></div>
            </div>
          </div>
          <div class="speed-control">
            <div class="speed-value">{{ speed }}</div>
            <el-slider
              v-model="speed"
              :min="1"
              :max="10"
              :step="1"
              :show-tooltip="false"
              :disabled="isDisabled"
              vertical
              height="90px"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ptz-control-panel {
  background-color: var(--el-bg-color);
  border-radius: var(--el-border-radius-base);
  box-shadow: var(--el-box-shadow-lighter);
}

.panel-header {
  padding: 8px 12px;
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;

  &:hover {
    background-color: var(--el-fill-color-light);
  }
}

.header-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--el-text-color-primary);
}

.collapse-arrow {
  font-size: 12px;
  transition: transform 0.2s ease;
  color: var(--el-text-color-secondary);

  &.collapsed {
    transform: rotate(-90deg);
  }
}

.title-icon {
  font-size: 14px;
  color: var(--el-color-primary);
}

.panel-content {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  height: auto;
}

.ptz-control-panel.collapsed .panel-content {
  height: 0;
  padding: 0;
  opacity: 0;
  margin: 0;
  pointer-events: none;
}

.control-form {
  padding: 16px;
  transform-origin: top;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  justify-content: center;
}

.ptz-control-panel.collapsed .control-form {
  transform: scaleY(0);
}

.ptz-controls {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  width: fit-content;
}

.direction-controls {
  display: flex;
  justify-content: center;
  padding: 0;
  width: fit-content;
}

.direction-pad {
  position: relative;
  width: 120px;
  height: 120px;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  gap: 4px;
  margin: 0 auto;
}

.direction-btn {
  --el-button-bg-color: var(--el-color-primary-light-8);
  --el-button-border-color: var(--el-color-primary-light-5);
  --el-button-hover-bg-color: var(--el-color-primary-light-7);
  --el-button-hover-border-color: var(--el-color-primary-light-4);
  --el-button-active-bg-color: var(--el-color-primary-light-5);
  --el-button-active-border-color: var(--el-color-primary);

  border-radius: 4px;
  padding: 0;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;

  .el-icon {
    font-size: 16px;
  }

  &:hover {
    transform: scale(1.05);
  }

  &:active {
    transform: scale(0.95);
  }
}

.direction-btn.up {
  grid-column: 2;
  grid-row: 1;
}

.direction-btn.right {
  grid-column: 3;
  grid-row: 2;
}

.direction-btn.down {
  grid-column: 2;
  grid-row: 3;
}

.direction-btn.left {
  grid-column: 1;
  grid-row: 2;
}

.direction-btn.up-left {
  grid-column: 1;
  grid-row: 1;
}

.direction-btn.up-right {
  grid-column: 3;
  grid-row: 1;
}

.direction-btn.down-left {
  grid-column: 1;
  grid-row: 3;
}

.direction-btn.down-right {
  grid-column: 3;
  grid-row: 3;
}

.direction-center {
  grid-column: 2;
  grid-row: 2;
  background-color: var(--el-color-primary-light-8);
  border-radius: 4px;
}

.control-groups,
.control-group {
  display: none;
}

.speed-control {
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0;
  margin-right: 0;
  width: fit-content;
}

.speed-value {
  color: var(--el-color-primary);
  font-weight: 500;
  font-size: 13px;
  margin-bottom: 4px;
  width: 16px;
  text-align: center;
  background-color: var(--el-color-primary-light-9);
  border-radius: 2px;
  padding: 1px 0;
}

:deep(.el-slider) {
  --el-slider-button-size: 10px;
  --el-slider-height: 2px;
  height: 90px;
}

:deep(.el-slider.is-vertical) {
  margin: 0;
}
</style>
