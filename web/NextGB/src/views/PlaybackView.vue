<script setup lang="ts">
import { ref, computed } from 'vue'
import DeviceTree from '@/components/monitor/DeviceTree.vue'
import VideoPlayer from '@/components/playback/VideoPlayer.vue'
import DateTimeRangePanel from '@/components/common/DateTimeRangePanel.vue'
import type { Device, ChannelInfo } from '@/types/api'
import { VideoPlay, VideoPause, VideoCamera, Download, Microphone } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const showBanner = ref(false)
const currentDevice = ref<Device>()
const currentChannel = ref<ChannelInfo>()
const videoPlayerRef = ref()
const volume = ref(100)

const isPlaying = computed(() => videoPlayerRef.value?.isPlaying || false)

const handlePlay = () => {
  videoPlayerRef.value?.play()
}

const handlePause = () => {
  videoPlayerRef.value?.pause()
}

const handleDeviceSelect = (data: { device: Device | undefined; channel: ChannelInfo }) => {
  currentDevice.value = data.device
  currentChannel.value = data.channel
  videoPlayerRef.value?.pause()
}

const handleSearch = ({ start, end }: { start: string; end: string }) => {
  console.log('查询时间范围：', {
    start,
    end,
    channel: currentChannel.value
  })
}

const option = {
  grid: {
    left: '3%',
    right: '3%',
    bottom: '15%',
    containLabel: true
  },
  xAxis: {
    type: 'time',
    min: dayjs().startOf('day').valueOf(),
    max: dayjs().endOf('day').valueOf(),
    axisLabel: {
      formatter: (value: number) => {
        return dayjs(value).format('HH:mm')
      }
    },
    splitLine: {
      show: true
    },
    axisTick: {
      alignWithLabel: true
    }
  }
}
</script>

<template>
  <div class="playback-container">
    <div class="left-panel">
      <DeviceTree @select="handleDeviceSelect" />
      <DateTimeRangePanel
        title="录像查询"
        @search="handleSearch"
      />
    </div>
    <div class="right-panel">
      <div class="playback-panel">
        <div class="video-container">
          <div
            class="video-banner"
            :class="{ show: showBanner }"
            @mouseenter="showBanner = true"
            @mouseleave="showBanner = false"
          >
            <div class="banner-content">
              <div class="info-item">
                <span class="label">码率：</span>
                <span class="value">2048 Kbps</span>
              </div>
              <div class="info-item">
                <span class="label">编码：</span>
                <span class="value">H.264</span>
              </div>
              <div class="info-item">
                <span class="label">分辨率：</span>
                <span class="value">1920×1080</span>
              </div>
              <div class="info-item">
                <span class="label">帧率：</span>
                <span class="value">25fps</span>
              </div>
            </div>
          </div>
          <VideoPlayer
            v-if="currentChannel"
            ref="videoPlayerRef"
            :device="currentDevice"
            :channel="currentChannel"
          />
          <div v-if="!currentChannel" class="placeholder">
            <el-empty description="请选择通道" />
          </div>
        </div>
        <div class="timeline-panel">
          <div class="timeline-ruler">
            <div class="timeline-scale">
              <div v-for="hour in 24" :key="hour" 
                class="hour-mark"
                :class="{
                  'major-mark': (hour - 1) % 6 === 0,
                  'medium-mark': (hour - 1) % 3 === 0 && (hour - 1) % 6 !== 0
                }"
              >
                <div class="hour-label">{{ (hour - 1).toString().padStart(2, '0') }}:00</div>
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
          <div class="control-buttons">
            <el-button-group>
              <el-button
                :icon="VideoPlay"
                :disabled="!currentChannel || isPlaying"
                @click="handlePlay"
                size="small"
              />
              <el-button
                :icon="VideoPause"
                :disabled="!currentChannel || !isPlaying"
                @click="handlePause"
                size="small"
              />
              <el-button :icon="VideoCamera" :disabled="!currentChannel" size="small" />
              <el-button :icon="Download" :disabled="!currentChannel" size="small" />
            </el-button-group>
            <el-slider
              v-model="volume"
              :max="100"
              :min="0"
              :disabled="!currentChannel"
              size="small"
              style="width: 100px"
            >
              <template #prepend>
                <el-icon><Microphone /></el-icon>
              </template>
            </el-slider>
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

.video-container {
  flex: 1;
  background-color: #000;
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 16/9;
  position: relative;
  overflow: hidden;
}

.placeholder {
  width: 100%;
  height: 100%;
  color: #909399;
  display: flex;
  align-items: center;
  justify-content: center;

  :deep(.el-empty) {
    color: #909399;
  }
}

.control-panel {
  height: 60px;
  background-color: #1a1a1a;
  display: flex;
  align-items: center;
}

.control-buttons {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 0;
  color: #fff;

  .el-button-group {
    margin: 0 16px;
    .el-button {
      width: 32px;
      height: 32px;
      padding: 0;
      --el-button-bg-color: transparent;
      --el-button-border-color: transparent;
      --el-button-hover-bg-color: rgba(255, 255, 255, 0.1);
      --el-button-hover-border-color: transparent;
      --el-button-active-bg-color: rgba(255, 255, 255, 0.2);
      --el-button-text-color: #fff;
      --el-button-disabled-text-color: rgba(255, 255, 255, 0.3);
      --el-button-disabled-bg-color: transparent;
      --el-button-disabled-border-color: transparent;

      :deep(.el-icon) {
        font-size: 16px;
      }
    }
  }

  :deep(.el-date-editor) {
    --el-input-bg-color: transparent;
    --el-input-border-color: rgba(255, 255, 255, 0.2);
    --el-input-hover-border-color: rgba(255, 255, 255, 0.3);
    --el-input-text-color: #fff;
    --el-input-placeholder-color: rgba(255, 255, 255, 0.5);
  }

  :deep(.el-slider) {
    --el-slider-main-bg-color: var(--el-color-primary);
    --el-slider-runway-bg-color: rgba(255, 255, 255, 0.2);
    --el-slider-stop-bg-color: rgba(255, 255, 255, 0.3);
    --el-slider-disabled-color: rgba(255, 255, 255, 0.2);
    .el-slider__runway {
      height: 4px;
    }
    .el-slider__button {
      border-color: var(--el-color-primary);
      background-color: var(--el-color-primary);
    }
  }
}

.search-panel {
  background-color: var(--el-bg-color);
  border-radius: var(--el-border-radius-base);
  box-shadow: var(--el-box-shadow-lighter);
}

.search-panel-header {
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

.search-panel-content {
  transition: all 0.2s ease;
  overflow: hidden;
}

.search-panel.collapsed .search-panel-content {
  height: 0;
  padding: 0;
  opacity: 0;
}

.datetime-range {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.datetime-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.datetime-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  padding-left: 4px;
}

.shortcuts {
  display: flex;
  gap: 8px;
  padding: 0 4px;
  
  :deep(.el-button) {
    height: 24px;
    padding: 0 8px;
    
    &.is-disabled {
      color: var(--el-text-color-disabled);
    }
  }
}

.calendar-wrapper {
  :deep(.el-input__wrapper) {
    padding: 0 8px;
    height: 32px;
  }
  
  :deep(.el-input__inner) {
    font-size: 13px;
  }
  
  :deep(.el-date-editor) {
    --el-date-editor-width: 100%;
  }
}

.search-form {
  padding: 12px;
}

:deep(.el-picker-panel) {
  --el-datepicker-border-color: var(--el-border-color-lighter);
}

.video-banner {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 40px;
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0.8), rgba(0, 0, 0, 0));
  transform: translateY(-100%);
  transition: transform 0.3s ease;
  z-index: 10;
  padding: 8px 16px;
  color: #fff;
}

.video-banner.show {
  transform: translateY(0);
}

.banner-content {
  display: flex;
  align-items: center;
  gap: 24px;
  height: 100%;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
}

.info-item .label {
  color: rgba(255, 255, 255, 0.7);
}

.info-item .value {
  color: #fff;
  font-family: monospace;
}

.video-banner::before {
  content: '';
  position: absolute;
  top: -20px;
  left: 0;
  right: 0;
  height: 20px;
}

.timeline-panel {
  height: 60px;
  background-color: #242424;
  position: relative;
  overflow: hidden;
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
    bottom: 20px;
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
    bottom: 20px;
  }
}

.hour-label {
  position: absolute;
  bottom: 20px;
  left: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  transform: translateX(-50%);
  transition: color 0.2s ease;
}

.hour-line {
  position: relative;
  width: 1px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.15);
  transition: all 0.2s ease;
}

.half-hour-mark {
  position: absolute;
  left: 50%;
  bottom: 0;
  width: 1px;
  height: 6px;
  background-color: rgba(255, 255, 255, 0.1);
  transition: all 0.2s ease;
}

.timeline-pointer {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: var(--el-color-primary);
  transform: translateX(-50%);
  transition: all 0.2s ease;
  
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
  transition: all 0.2s ease;
  
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

.hour-mark:hover {
  .hour-line {
    height: 20px;
    background-color: rgba(255, 255, 255, 0.5);
  }
  
  .hour-label {
    color: #fff;
  }
  
  .half-hour-mark {
    height: 10px;
    background-color: rgba(255, 255, 255, 0.3);
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
