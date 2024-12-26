<script setup lang="ts">
import { ref } from 'vue'
import { ArrowRight, VideoCamera } from '@element-plus/icons-vue'
import { ArrowUp, ArrowDown, ArrowLeft, ArrowRight as ArrowRightControl } from '@element-plus/icons-vue'

const props = defineProps<{
  title?: string
}>()

const emit = defineEmits<{
  control: [direction: string]
}>()

const isCollapsed = ref(false)

const handlePtzControl = (direction: string) => {
  emit('control', direction)
}
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
              <el-button class="direction-btn up" @click="handlePtzControl('up')">
                <el-icon><ArrowUp /></el-icon>
              </el-button>
              <el-button class="direction-btn right" @click="handlePtzControl('right')">
                <el-icon><ArrowRightControl /></el-icon>
              </el-button>
              <el-button class="direction-btn down" @click="handlePtzControl('down')">
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <el-button class="direction-btn left" @click="handlePtzControl('left')">
                <el-icon><ArrowLeft /></el-icon>
              </el-button>
              <div class="direction-center"></div>
            </div>
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
}

.ptz-control-panel.collapsed .control-form {
  transform: scaleY(0);
}

.ptz-controls {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.direction-controls {
  display: flex;
  justify-content: center;
  padding: 0;
}

.direction-pad {
  position: relative;
  width: 120px;
  height: 120px;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  gap: 4px;
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
</style> 