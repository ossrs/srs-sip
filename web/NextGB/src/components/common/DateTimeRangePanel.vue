<script setup lang="ts">
import { ref } from 'vue'
import { ArrowRight, VideoCamera, Search } from '@element-plus/icons-vue'

const props = defineProps<{
  title?: string
}>()

const emit = defineEmits<{
  search: [{ start: string; end: string }]
}>()

const isCollapsed = ref(false)
const startDateTime = ref('')
const endDateTime = ref('')

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
  emit('search', {
    start: startDateTime.value,
    end: endDateTime.value,
  })
}
</script>

<template>
  <div class="datetime-range-panel" :class="{ collapsed: isCollapsed }">
    <div class="panel-header" @click="isCollapsed = !isCollapsed">
      <div class="header-title">
        <el-icon class="collapse-arrow" :class="{ collapsed: isCollapsed }">
          <ArrowRight />
        </el-icon>
        <el-icon class="title-icon"><VideoCamera /></el-icon>
        <span>{{ title || '时间范围' }}</span>
      </div>
    </div>
    <div class="panel-content">
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
            <el-button text size="small" @click="handleShortcut('today')"> 今天 </el-button>
            <el-button text size="small" @click="handleShortcut('yesterday')"> 昨天 </el-button>
            <el-button text size="small" @click="handleShortcut('lastWeek')"> 最近一周 </el-button>
          </div>
        </div>
        <template v-if="!isCollapsed">
          <el-button
            type="primary"
            :disabled="!startDateTime || !endDateTime"
            @click="handleSearch"
            style="width: 100%"
          >
            <el-icon><Search /></el-icon>
            查询
          </el-button>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.datetime-range-panel {
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
  transition: all 0.2s ease;
  overflow: hidden;
}

.datetime-range-panel.collapsed .panel-content {
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
</style>
