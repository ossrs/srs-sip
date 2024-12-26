<script setup lang="ts">
import { ref, computed } from 'vue'
import DeviceTree from '@/components/monitor/DeviceTree.vue'
import VideoPlayer from '@/components/playback/VideoPlayer.vue'
import type { Device, ChannelInfo } from '@/types/api'
import { VideoPlay, VideoPause, VideoCamera, Download, Microphone, ArrowRight, Search } from '@element-plus/icons-vue'

const showBanner = ref(false)
const currentDevice = ref<Device>()
const currentChannel = ref<ChannelInfo>()
const startDateTime = ref('')
const endDateTime = ref('')
const videoPlayerRef = ref()
const volume = ref(100)
const isSearchPanelCollapsed = ref(false)

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

const formatDateTime = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const handleShortcut = (type: string) => {
  const now = new Date()
  const start = new Date()
  
  switch (type) {
    case 'today': {
      start.setHours(0, 0, 0, 0)
      break
    }
    case 'yesterday': {
      start.setDate(start.getDate() - 1)
      start.setHours(0, 0, 0, 0)
      now.setDate(now.getDate() - 1)
      break
    }
    case 'lastWeek': {
      start.setDate(start.getDate() - 7)
      start.setHours(0, 0, 0, 0)
      break
    }
  }
  
  startDateTime.value = formatDateTime(start)
  now.setHours(23, 59, 59)
  endDateTime.value = formatDateTime(now)
}

const handleSearch = () => {
  if (!startDateTime.value || !endDateTime.value) return
  console.log('查询时间范围：', {
    start: startDateTime.value,
    end: endDateTime.value,
    channel: currentChannel.value
  })
}
</script>

<template>
  <div class="playback-container">
    <div class="left-panel">
      <DeviceTree @select="handleDeviceSelect" />
      <div class="search-panel" :class="{ collapsed: isSearchPanelCollapsed }">
        <div class="search-panel-header" @click="isSearchPanelCollapsed = !isSearchPanelCollapsed">
          <div class="header-title">
            <el-icon class="collapse-arrow" :class="{ collapsed: isSearchPanelCollapsed }">
              <ArrowRight />
            </el-icon>
            <el-icon class="title-icon"><VideoCamera /></el-icon>
            <span>录像查询</span>
          </div>
        </div>
        <div class="search-panel-content">
          <div class="search-form">
            <div class="form-item calendar-wrapper">
              <div class="datetime-range">
                <div class="datetime-item">
                  <div class="datetime-label">开始时间：</div>
                  <el-date-picker
                    v-model="startDateTime"
                    type="datetime"
                    :editable="false"
                    placeholder="开始时间"
                    value-format="YYYY-MM-DD HH:mm:ss"
                    style="width: 100%"
                  />
                </div>
                <div class="datetime-item">
                  <div class="datetime-label">结束时间：</div>
                  <el-date-picker
                    v-model="endDateTime"
                    type="datetime"
                    :editable="false"
                    placeholder="结束时间"
                    value-format="YYYY-MM-DD HH:mm:ss"
                    style="width: 100%"
                  />
                </div>
              </div>
              <div class="shortcuts">
                <el-button
                  text
                  size="small"
                  @click="handleShortcut('today')"
                >
                  今天
                </el-button>
                <el-button
                  text
                  size="small"
                  @click="handleShortcut('yesterday')"
                >
                  昨天
                </el-button>
                <el-button
                  text
                  size="small"
                  @click="handleShortcut('lastWeek')"
                >
                  最近一周
                </el-button>
              </div>
            </div>
            <template v-if="!isSearchPanelCollapsed">
              <el-button
                type="primary"
                :disabled="!startDateTime || !endDateTime"
                @click="handleSearch"
                style="width: 100%"
              >
                <el-icon><Search /></el-icon>
                查询录像
              </el-button>
            </template>
          </div>
        </div>
      </div>
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
</style>
