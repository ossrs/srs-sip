<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { Device, ChannelInfo } from '@/types/api'

const props = defineProps<{
  device?: Device
  channel?: ChannelInfo
}>()

const videoRef = ref<HTMLVideoElement>()

// 播放控制
const isPlaying = ref(false)

const play = () => {
  if (videoRef.value) {
    videoRef.value.play()
    isPlaying.value = true
  }
}

const pause = () => {
  if (videoRef.value) {
    videoRef.value.pause()
    isPlaying.value = false
  }
}

// 暴露方法给父组件
defineExpose({
  play,
  pause,
  isPlaying
})

// 组件卸载时停止播放
onUnmounted(() => {
  if (videoRef.value) {
    videoRef.value.pause()
  }
})
</script>

<template>
  <div class="video-player">
    <video
      ref="videoRef"
      class="video-element"
      controls
      controlsList="nodownload"
      :poster="channel?.info?.snapshot_url"
    >
      <source :src="channel?.info?.stream_url" type="video/mp4">
      您的浏览器不支持 video 标签
    </video>
  </div>
</template>

<style scoped>
.video-player {
  width: 100%;
  height: 100%;
  background-color: #000;
  position: relative;
  overflow: hidden;
}

.video-element {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

/* 自定义播放器控件样式 */
.video-element::-webkit-media-controls {
  background-color: rgba(0, 0, 0, 0.5);
}

.video-element::-webkit-media-controls-panel {
  display: flex;
  align-items: center;
  padding: 0 10px;
}

.video-element::-webkit-media-controls-play-button {
  display: none;
}

.video-element::-webkit-media-controls-timeline {
  display: none;
}

.video-element::-webkit-media-controls-current-time-display,
.video-element::-webkit-media-controls-time-remaining-display {
  color: #fff;
}

.video-element::-webkit-media-controls-volume-slider {
  display: none;
}

.video-element::-webkit-media-controls-mute-button {
  display: none;
}
</style> 