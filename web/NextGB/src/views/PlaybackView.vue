<script setup lang="ts">
import { ref, computed } from 'vue'
import DeviceTree from '@/components/monitor/DeviceTree.vue'
import VideoPlayer from '@/components/playback/VideoPlayer.vue'
import type { Device, ChannelInfo } from '@/types/api'
import { VideoPlay, VideoPause, VideoCamera, Download, Microphone } from '@element-plus/icons-vue'

const showBanner = ref(false)
const currentDevice = ref<Device>()
const currentChannel = ref<ChannelInfo>()
const selectedDate = ref<Date>()
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

const handleSearch = () => {
  if (!currentChannel.value || !selectedDate.value) return
  // TODO: 实现录像查询逻辑
}
</script>

<template>
  <div class="playback-container">
    <div class="left-panel">
      <DeviceTree @select="handleDeviceSelect" />
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
          <div class="date-picker">
            <el-date-picker
              v-model="selectedDate"
              type="date"
              placeholder="选择日期"
              :disabled="!currentChannel"
              size="small"
            />
            <el-button
              type="primary"
              size="small"
              :disabled="!currentChannel || !selectedDate"
              @click="handleSearch"
            >
              查询
            </el-button>
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
  height: 100px;
  background-color: #1a1a1a;
  display: flex;
  flex-direction: column;
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

.date-picker {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #242424;
  padding: 0 16px;
  gap: 16px;

  :deep(.el-date-editor) {
    --el-input-bg-color: transparent;
    --el-input-border-color: rgba(255, 255, 255, 0.2);
    --el-input-hover-border-color: rgba(255, 255, 255, 0.3);
    --el-input-text-color: #fff;
    --el-input-placeholder-color: rgba(255, 255, 255, 0.5);
  }

  :deep(.el-button) {
    --el-button-bg-color: var(--el-color-primary);
    --el-button-border-color: var(--el-color-primary);
    --el-button-hover-bg-color: var(--el-color-primary-light-3);
    --el-button-hover-border-color: var(--el-color-primary-light-3);
    --el-button-active-bg-color: var(--el-color-primary-dark-2);
    --el-button-active-border-color: var(--el-color-primary-dark-2);
    --el-button-disabled-bg-color: var(--el-color-primary-light-5);
    --el-button-disabled-border-color: var(--el-color-primary-light-5);
  }
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
