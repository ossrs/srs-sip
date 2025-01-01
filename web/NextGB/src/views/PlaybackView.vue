<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import DeviceTree from '@/components/monitor/DeviceTree.vue'
import MonitorGrid from '@/components/monitor/MonitorGrid.vue'
import DateTimeRangePanel from '@/components/common/DateTimeRangePanel.vue'
import type { Device, ChannelInfo } from '@/types/api'
import type { LayoutConfig } from '@/types/layout'
import {
  CaretRight,
  VideoPause,
  CircleClose,
  DArrowRight,
  DArrowLeft,
  Timer,
  Picture,
  Download,
  Microphone
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

type LayoutKey = '1'
type LayoutConfigs = Record<LayoutKey, LayoutConfig>

// 布局配置
const layouts: LayoutConfigs = {
  '1': { cols: 1, rows: 1, size: 1, label: '单屏' },
} as const

const monitorGridRef = ref()
const currentDevice = ref<Device>()
const currentChannel = ref<ChannelInfo>()
const volume = ref(100)
const playbackSpeed = ref(1.0)
const playbackSpeeds = [0.25, 0.5, 1.0, 2.0, 4.0, 8.0]
const currentLayout = ref<'1'>('1') // 固定为单屏模式
const defaultMuted = ref(true)
const activeWindow = ref<{ deviceId: string; channelId: string } | null>(null)

// 时间轴刻度显示控制
const timelineWidth = ref(0)
const showAllLabels = computed(() => timelineWidth.value >= 720) // 当宽度大于720px时显示所有标签
const showMediumLabels = computed(() => timelineWidth.value >= 480) // 当宽度大于480px时显示中等标签

// 屏幕尺寸类型
const screenType = computed(() => {
  if (timelineWidth.value >= 720) return '大屏'
  if (timelineWidth.value >= 480) return '中屏'
  return '小屏'
})

// 根据屏幕宽度调整时间轴高度
const timelineHeight = computed(() => {
  if (timelineWidth.value >= 720) return 60
  if (timelineWidth.value >= 480) return 48
  return 36
})

// 监听时间轴宽度变化
const updateTimelineWidth = () => {
  const timeline = document.querySelector('.timeline-scale')
  if (timeline) {
    timelineWidth.value = timeline.clientWidth
    console.log(`时间轴宽度: ${timelineWidth.value}px, 当前屏幕: ${screenType.value}`)
  }
}

// 监听窗口大小变化
onMounted(() => {
  updateTimelineWidth()
  window.addEventListener('resize', updateTimelineWidth)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateTimelineWidth)
})

const handleDeviceSelect = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  currentDevice.value = data.device
  currentChannel.value = data.channel
}

const handleDevicePlay = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  if (data.channel.device_id) {
    monitorGridRef.value?.play({
      ...data.device,
      channel: data.channel,
    })
  } else {
    ElMessage.warning('设备信息不完整')
  }
}

const handleWindowSelect = (data: { deviceId: string; channelId: string } | null) => {
  activeWindow.value = data
}

const handleSearch = ({ start, end }: { start: string; end: string }) => {
  console.log('查询时间范围：', {
    start,
    end,
    channel: currentChannel.value
  })
}

const clearAllDevices = () => {
  monitorGridRef.value?.clearAllDevices()
}
</script>

<template>
  <div class="playback-container">
    <div class="left-panel">
      <DeviceTree 
        @select="handleDeviceSelect"
        @play="handleDevicePlay"
      />
      <DateTimeRangePanel
        title="录像查询"
        @search="handleSearch"
      />
    </div>
    <div class="right-panel">
      <div class="playback-panel">
        <MonitorGrid 
          ref="monitorGridRef" 
          v-model="currentLayout"
          :layouts="layouts"
          :default-muted="defaultMuted"
          @window-select="handleWindowSelect"
        />
        <div class="timeline-panel" :style="{ height: `${timelineHeight}px` }">
          <div class="timeline-ruler">
            <div class="timeline-scale">
              <div v-for="hour in 24" :key="hour" 
                class="hour-mark"
                :class="{
                  'major-mark': (hour - 1) % 6 === 0,
                  'medium-mark': (hour - 1) % 3 === 0 && (hour - 1) % 6 !== 0,
                  'minor-mark': (hour - 1) % 3 !== 0
                }"
              >
                <div 
                  v-if="(hour - 1) % 6 === 0 || (showMediumLabels && (hour - 1) % 3 === 0) || showAllLabels"
                  class="hour-label"
                >
                  {{ (hour - 1).toString().padStart(2, '0') }}:00
                </div>
                <div class="hour-line"></div>
                <div class="half-hour-mark"></div>
              </div>
              <div class="hour-mark major-mark" style="flex: 0 0 auto;">
                <div class="hour-label">24:00</div>
                <div class="hour-line"></div>
              </div>
            </div>
            <div class="timeline-pointer" :style="{ left: '0%' }">
              <div class="pointer-head"></div>
            </div>
          </div>
        </div>
        <div class="control-panel">
          <div class="control-group">
            <el-button-group>
              <el-button
                :icon="CaretRight"
                :disabled="!activeWindow"
                size="small"
                title="播放"
              />
              <el-button
                :icon="VideoPause"
                :disabled="!activeWindow"
                size="small"
                title="暂停"
              />
              <el-button
                :icon="CircleClose"
                :disabled="!activeWindow"
                @click="clearAllDevices"
                size="small"
                title="停止"
              />
            </el-button-group>

            <el-button-group>
              <el-button
                :icon="DArrowLeft"
                :disabled="!activeWindow || playbackSpeed <= playbackSpeeds[0]"
                @click="playbackSpeed = Math.max(playbackSpeed / 2, playbackSpeeds[0])"
                size="small"
                title="减速"
              />
              <el-button
                :icon="Timer"
                :disabled="!activeWindow"
                @click="playbackSpeed = 1.0"
                size="small"
                :class="{ 'speed-active': playbackSpeed !== 1.0 }"
                title="当前速度"
              >
                <span class="speed-text">{{ playbackSpeed }}x</span>
              </el-button>
              <el-button
                :icon="DArrowRight"
                :disabled="!activeWindow || playbackSpeed >= playbackSpeeds[playbackSpeeds.length - 1]"
                @click="playbackSpeed = Math.min(playbackSpeed * 2, playbackSpeeds[playbackSpeeds.length - 1])"
                size="small"
                title="加速"
              />
            </el-button-group>

            <div class="time-info">
              <span class="current-time">00:00:00</span>
              <span class="time-separator">/</span>
              <span class="total-time">00:00:00</span>
            </div>

            <el-button-group>
              <el-button 
                :icon="Picture" 
                :disabled="!activeWindow" 
                size="small"
                title="截图"
              />
              <el-button 
                :icon="Download" 
                :disabled="!activeWindow" 
                size="small"
                title="下载"
              />
            </el-button-group>
          </div>

          <div class="volume-control">
            <el-icon><Microphone /></el-icon>
            <el-slider
              v-model="volume"
              :max="100"
              :min="0"
              :disabled="!activeWindow"
              size="small"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.playback-container {
  height: 100%;
  display: flex;
  gap: 16px;
}

.left-panel {
  width: 280px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.right-panel {
  flex: 1;
  background-color: #fff;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
}

.playback-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.control-panel {
  height: 60px;
  background-color: #1a1a1a;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  gap: 24px;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.time-info {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 8px;
}

.time-separator {
  color: rgba(255, 255, 255, 0.3);
  margin: 0 2px;
}

.current-time {
  color: var(--el-color-primary);
}

.total-time {
  color: rgba(255, 255, 255, 0.5);
}

.volume-control {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 140px;
  color: rgba(255, 255, 255, 0.8);
  
  :deep(.el-icon) {
    font-size: 18px;
  }
}

.el-button-group {
  .el-button {
    width: 36px;
    height: 36px;
    padding: 0;
    --el-button-bg-color: transparent;
    --el-button-border-color: transparent;
    --el-button-hover-bg-color: rgba(255, 255, 255, 0.1);
    --el-button-hover-border-color: transparent;
    --el-button-active-bg-color: rgba(255, 255, 255, 0.15);
    --el-button-text-color: rgba(255, 255, 255, 0.85);
    --el-button-disabled-text-color: rgba(255, 255, 255, 0.3);
    --el-button-disabled-bg-color: transparent;
    --el-button-disabled-border-color: transparent;

    :deep(.el-icon) {
      font-size: 20px;
    }

    &:hover:not(:disabled) {
      --el-button-text-color: var(--el-color-primary);
    }

    &:has(.speed-text) {
      width: auto;
      padding: 0 12px;
    }
  }
}

:deep(.el-slider) {
  --el-slider-main-bg-color: var(--el-color-primary);
  --el-slider-runway-bg-color: rgba(255, 255, 255, 0.15);
  --el-slider-stop-bg-color: rgba(255, 255, 255, 0.2);
  --el-slider-disabled-color: rgba(255, 255, 255, 0.1);
  
  .el-slider__runway {
    height: 3px;
  }
  
  .el-slider__button {
    border: none;
    width: 10px;
    height: 10px;
    background-color: #fff;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
    transition: transform 0.2s ease;
    
    &:hover {
      transform: scale(1.3);
    }
  }

  .el-slider__bar {
    height: 3px;
  }
}

.speed-text {
  font-size: 13px;
  margin-left: 6px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.speed-active {
  --el-button-text-color: var(--el-color-primary) !important;
  --el-button-bg-color: rgba(var(--el-color-primary-rgb), 0.1) !important;
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

.timeline-panel {
  background-color: #242424;
  position: relative;
  overflow: visible;
  user-select: none;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 8px 0;
}

.timeline-ruler {
  height: 100%;
  position: relative;
  padding: 0 24px;
  display: flex;
  align-items: flex-end;
}

.timeline-scale {
  height: 100%;
  display: flex;
  position: relative;
  width: 100%;
}

.hour-mark {
  flex: 1;
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding-top: 16px;
}

.hour-label {
  position: absolute;
  left: 0;
  top: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  transform: translateX(-50%);
  white-space: nowrap;
  line-height: 1;
}

.hour-line {
  position: relative;
  width: 1px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.15);
}

.half-hour-mark {
  position: absolute;
  left: 50%;
  bottom: 0;
  width: 1px;
  height: 6px;
  background-color: rgba(255, 255, 255, 0.1);
}

.timeline-pointer {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: var(--el-color-primary);
  transform: translateX(-50%);
  
  &::after {
    content: '';
    position: absolute;
    top: 0;
    bottom: 0;
    left: 50%;
    width: 4px;
    transform: translateX(-50%);
    background: linear-gradient(
      90deg,
      transparent,
      rgba(var(--el-color-primary-rgb), 0.2),
      transparent
    );
  }
}

.pointer-head {
  position: absolute;
  top: -1px;
  left: 50%;
  transform: translateX(-50%);
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: var(--el-color-primary);
  box-shadow: 0 0 6px rgba(var(--el-color-primary-rgb), 0.6);
  
  &::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background-color: #fff;
  }
}

.major-mark {
  .hour-line {
    height: 16px;
    background-color: rgba(255, 255, 255, 0.4);
    width: 2px;
  }
  
  .hour-label {
    color: rgba(255, 255, 255, 0.95);
    font-weight: 500;
    font-size: 12px;
  }
}

.medium-mark {
  .hour-line {
    height: 12px;
    background-color: rgba(255, 255, 255, 0.25);
    width: 1.5px;
  }
  
  .hour-label {
    color: rgba(255, 255, 255, 0.7);
  }
}

.minor-mark {
  .hour-line {
    height: 8px;
    background-color: rgba(255, 255, 255, 0.15);
    width: 1px;
  }
  
  .hour-label {
    color: rgba(255, 255, 255, 0.4);
    font-size: 10px;
    opacity: 0.8;
  }
}

.timeline-ruler:hover .timeline-pointer {
  box-shadow: 0 0 8px rgba(var(--el-color-primary-rgb), 0.2);
}

.timeline-panel::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.1) 20%,
    rgba(255, 255, 255, 0.1) 80%,
    transparent
  );
}
</style>
