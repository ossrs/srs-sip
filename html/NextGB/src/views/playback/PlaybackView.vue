<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, onActivated, onDeactivated, watch } from 'vue'
import DeviceTree from './DeviceTree.vue'
import MonitorGrid from '@/components/monitor/MonitorGrid.vue'
import DateTimeRangePanel from '@/components/common/DateTimeRangePanel.vue'
import type { Device, ChannelInfo, RecordInfoResponse } from '@/api/types'
import type { LayoutConfig } from '@/types/layout'
import { deviceApi } from '@/api'
import dayjs from 'dayjs'
import { VideoPlay, VideoPause, CloseBold, Microphone } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

type LayoutKey = '1'
type LayoutConfigs = Record<LayoutKey, LayoutConfig>

// 布局配置
const layouts: LayoutConfigs = {
  '1': { cols: 1, rows: 1, size: 1, label: '单屏' },
} as const

const monitorGridRef = ref()
const selectedChannels = ref<{ device: Device | undefined; channel: ChannelInfo }[]>([])
const volume = ref(100)
const currentLayout = ref<'1'>('1') // 固定为单屏模式
const defaultMuted = ref(true)
const activeWindow = ref<{ deviceId: string; channelId: string } | null>(null)
const isPlaying = ref(false) // 添加播放状态变量
const isFirstPlay = ref(true)

// 时间轴刻度显示控制
const timelineWidth = ref(0)
const showAllLabels = computed(() => timelineWidth.value >= 720) // 当宽度大于720px时显示所有标签
const showMediumLabels = computed(() => timelineWidth.value >= 480) // 当宽度大于480px时显示中等标签

// 时间轴光标位置
const cursorPosition = ref(0)
const cursorTime = ref('')
const isTimelineHovered = ref(false)

const getTimeFromEvent = (event: MouseEvent, element: HTMLElement) => {
  const rect = element.getBoundingClientRect()
  const position = ((event.clientX - rect.left) / rect.width) * 100
  const normalizedPosition = Math.max(0, Math.min(100, position))

  const totalMinutes = 24 * 60
  const minutes = Math.floor((normalizedPosition / 100) * totalMinutes)
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60

  return {
    position: normalizedPosition,
    time: dayjs().startOf('day').add(hours, 'hour').add(mins, 'minute'),
  }
}

const handleTimelineMouseMove = (event: MouseEvent) => {
  const timeline = event.currentTarget as HTMLElement
  const { position, time } = getTimeFromEvent(event, timeline)
  cursorPosition.value = position
  cursorTime.value = time.format('HH:mm:ss')
}

// 处理时间轴鼠标进入/离开
const handleTimelineMouseEnter = () => {
  isTimelineHovered.value = true
}

const handleTimelineMouseLeave = () => {
  isTimelineHovered.value = false
}

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

const handleWindowSelect = (data: { deviceId: string; channelId: string } | null) => {
  activeWindow.value = data
}

const recordSegments = ref<RecordInfoResponse[]>([])

const handleQueryRecord = async ({ start, end }: { start: string; end: string }) => {
  if (selectedChannels.value.length === 0) {
    ElMessage.warning('请先选择要查询的通道')
    return
  }

  try {
    const promises = selectedChannels.value.map(async ({ device, channel }) => {
      if (!device?.device_id || !channel.device_id) return []

      const response = await deviceApi.queryRecord({
        device_id: device.device_id,
        channel_id: channel.device_id,
        start_time: dayjs(start).unix(),
        end_time: dayjs(end).unix(),
      })
      return Array.isArray(response.data) ? response.data : []
    })

    const results = await Promise.all<RecordInfoResponse[]>(promises)
    recordSegments.value = results.flat()

    // 自动激活第一个选中的通道
    if (selectedChannels.value.length > 0) {
      const firstChannel = selectedChannels.value[0]
      if (firstChannel.device?.device_id && firstChannel.channel.device_id) {
        activeWindow.value = {
          deviceId: firstChannel.device.device_id,
          channelId: firstChannel.channel.device_id
        }
      }
    }
  } catch (error) {
    console.error('查询录像失败:', error)
    ElMessage.error('查询录像失败')
    recordSegments.value = []
  }
}

const handleStop = () => {
  monitorGridRef.value?.stop(0)
  isPlaying.value = false // 设置播放状态为 false
  
  // 确保 activeWindow 始终指向第一个屏幕
  if (selectedChannels.value.length > 0) {
    const firstChannel = selectedChannels.value[0]
    if (firstChannel.device?.device_id && firstChannel.channel.device_id) {
      activeWindow.value = {
        deviceId: firstChannel.device.device_id,
        channelId: firstChannel.channel.device_id
      }
    }
  }
}

// 监听 selectedChannels 变化，确保 activeWindow 始终指向第一个屏幕
watch(selectedChannels, (newChannels) => {
  if (newChannels.length > 0) {
    const firstChannel = newChannels[0]
    if (firstChannel.device?.device_id && firstChannel.channel.device_id) {
      activeWindow.value = {
        deviceId: firstChannel.device.device_id,
        channelId: firstChannel.channel.device_id
      }
    }
  }
}, { immediate: true })

const calculatePosition = (time: string) => {
  const hour = dayjs(time).hour()
  const minute = dayjs(time).minute()
  return ((hour * 60 + minute) / (24 * 60)) * 100
}

const calculateWidth = (start: string, end: string) => {
  const startMinutes = dayjs(start).hour() * 60 + dayjs(start).minute()
  const endMinutes = dayjs(end).hour() * 60 + dayjs(end).minute()
  return ((endMinutes - startMinutes) / (24 * 60)) * 100
}

// 添加激活/停用处理
onActivated(() => {
  console.log('PlaybackView activated')
})

onDeactivated(() => {
  console.log('PlaybackView deactivated')
})

// 组件名称（用于 keep-alive）
defineOptions({
  name: 'PlaybackView',
})

const handleTimelineDoubleClick = async (event: MouseEvent) => {
  if (selectedChannels.value.length === 0) {
    ElMessage.warning('请先选择要播放的通道')
    return
  }

  const timeline = event.currentTarget as HTMLElement
  const { time } = getTimeFromEvent(event, timeline)
  const endTime = dayjs().endOf('day')

  // 只播放第一个选中的通道
  const { device, channel } = selectedChannels.value[0]
  if (!device?.device_id || !channel.device_id) {
    ElMessage.warning('设备信息不完整')
    return
  }

  try {
    monitorGridRef.value?.play({
      ...device,
      channel: channel,
      play_type: 1, // 1 表示回放
      start_time: time.unix(),
      end_time: endTime.unix(), // 使用当天 23:59:59 的时间戳
    })
    isPlaying.value = true // 设置播放状态为 true
  } catch (error) {
    console.error('播放录像失败:', error)
    ElMessage.error('播放录像失败')
  }
}

// 处理播放/暂停切换
const handlePlayPause = async () => {
  if (!activeWindow.value || selectedChannels.value.length === 0) return

  if (!isPlaying.value) {
    // 开始播放
    const { device, channel } = selectedChannels.value[0]
    if (!device?.device_id || !channel.device_id) {
      ElMessage.warning('设备信息不完整')
      return
    }

    try {
      // 如果有录像段，则根据是否是第一次播放来决定是调用 play 还是 resume
      if (recordSegments.value.length === 0) {
        ElMessage.warning('没有可播放的录像')
        return
      }

      if (isFirstPlay.value) {
        monitorGridRef.value?.play({
          ...device,
          channel: channel,
          play_type: 1, // 1 表示回放
          start_time: recordSegments.value[0].start_time,
          end_time: recordSegments.value[0].end_time,
        })
        isFirstPlay.value = false
      } else {
        monitorGridRef.value?.resume(0)
      }
      isPlaying.value = true
    } catch (error) {
      console.error('播放录像失败:', error)
      ElMessage.error('播放录像失败')
      return
    }
  } else {
    // 暂停播放
    try {
      const { device, channel } = selectedChannels.value[0]
      if (!device?.device_id || !channel.device_id) {
        ElMessage.warning('设备信息不完整')
        return
      }
      console.log('暂停录像')
      monitorGridRef.value?.pause(0)
      isPlaying.value = false
    } catch (error) {
      console.error('暂停录像失败:', error)
      ElMessage.error('暂停录像失败')
    }
  }
}
</script>

<template>
  <div class="playback-container">
    <div class="left-panel">
      <DeviceTree v-model:selectedChannels="selectedChannels" />
      <DateTimeRangePanel title="录像查询" @search="handleQueryRecord" />
    </div>
    <div class="right-panel">
      <div class="playback-panel">
        <MonitorGrid
          ref="monitorGridRef"
          v-model="currentLayout"
          :layouts="layouts"
          :default-muted="defaultMuted"
          :show-border="false"
          @window-select="handleWindowSelect"
        />
        <div class="timeline-panel" :style="{ height: `${timelineHeight}px` }">
          <div class="timeline-ruler">
            <div
              class="timeline-scale"
              @mousemove="handleTimelineMouseMove"
              @mouseenter="handleTimelineMouseEnter"
              @mouseleave="handleTimelineMouseLeave"
              @dblclick="handleTimelineDoubleClick"
            >
              <div class="timeline-marks">
                <div
                  v-for="hour in 24"
                  :key="hour"
                  class="hour-mark"
                  :class="{
                    'major-mark': (hour - 1) % 6 === 0,
                    'medium-mark': (hour - 1) % 3 === 0 && (hour - 1) % 6 !== 0,
                    'minor-mark': (hour - 1) % 3 !== 0,
                  }"
                >
                  <div
                    v-if="
                      (hour - 1) % 6 === 0 ||
                      (showMediumLabels && (hour - 1) % 3 === 0) ||
                      showAllLabels
                    "
                    class="hour-label"
                  >
                    {{ (hour - 1).toString().padStart(2, '0') }}
                  </div>
                  <div class="hour-line"></div>
                  <div class="half-hour-mark"></div>
                </div>
                <div class="hour-mark major-mark" style="flex: 0 0 auto">
                  <div class="hour-label">24</div>
                  <div class="hour-line"></div>
                </div>
              </div>

              <div
                class="timeline-cursor"
                :class="{ visible: isTimelineHovered }"
                :style="{ left: `${cursorPosition}%` }"
              >
                <div class="cursor-time" :class="{ visible: isTimelineHovered }">
                  {{ cursorTime }}
                </div>
              </div>

              <div class="record-segments">
                <div
                  v-for="(segment, index) in recordSegments"
                  :key="index"
                  class="record-segment"
                  :style="{
                    left: `${calculatePosition(segment.start_time)}%`,
                    width: `${calculateWidth(segment.start_time, segment.end_time)}%`,
                  }"
                  :title="`${dayjs(segment.start_time).format('HH:mm:ss')} - ${dayjs(segment.end_time).format('HH:mm:ss')}`"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="control-panel">
          <div class="control-group">
            <el-button-group>
              <el-button
                :icon="isPlaying ? VideoPause : VideoPlay"
                :disabled="!activeWindow"
                size="large"
                :title="isPlaying ? '暂停' : '播放'"
                @click="handlePlayPause"
              />
              <el-button
                :icon="CloseBold"
                :disabled="!activeWindow"
                @click="handleStop"
                size="large"
                title="停止"
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

.timeline-marks {
  height: 100%;
  display: flex;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  width: 100%;
  pointer-events: none;
}

.timeline-scale {
  height: 100%;
  position: relative;
  width: 100%;
  cursor: pointer;
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

.timeline-cursor {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 1px;
  background-color: var(--el-color-warning);
  transform: translateX(-50%);
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 2;

  &.visible {
    opacity: 1;
  }

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
      rgba(var(--el-color-warning-rgb), 0.2),
      transparent
    );
  }
}

.cursor-time {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--el-color-warning);
  color: #000;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  white-space: nowrap;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s ease;

  &.visible {
    opacity: 1;
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

.record-segments {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 24px;
  pointer-events: none;
}

.record-segment {
  position: absolute;
  height: 8px;
  bottom: 12px;
  background-color: var(--el-color-success);
  opacity: 0.85;
  pointer-events: auto;
  cursor: pointer;

  & + & {
    margin-left: 2px;
  }
}
</style>
