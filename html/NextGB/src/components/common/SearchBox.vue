<script setup lang="ts">
import { ref } from 'vue'
import { Search, Refresh, Expand, List } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

interface Props {
  loading?: boolean
  showViewModeSwitch?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  showViewModeSwitch: false,
})

const searchQuery = ref('')
const viewMode = ref<'tree' | 'list'>('tree')
const tooltipRef = ref()

const emit = defineEmits<{
  (e: 'update:searchQuery', value: string): void
  (e: 'update:viewMode', value: 'tree' | 'list'): void
  (e: 'refresh'): void
}>()

const handleSearchInput = (value: string) => {
  searchQuery.value = value
  emit('update:searchQuery', value)
}

const handleViewModeChange = () => {
  viewMode.value = viewMode.value === 'tree' ? 'list' : 'tree'
  emit('update:viewMode', viewMode.value)
}

const handleRefresh = () => {
  emit('refresh')
  tooltipRef.value?.hide()
}
</script>

<template>
  <div class="search-box">
    <div class="search-wrapper">
      <el-input
        v-model="searchQuery"
        placeholder="搜索设备或通道..."
        clearable
        size="small"
        @input="handleSearchInput"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <div class="action-buttons">
        <el-tooltip
          v-if="showViewModeSwitch"
          :content="viewMode === 'tree' ? '切换到列表视图' : '切换到树形视图'"
          placement="top"
        >
          <el-button
            class="action-btn"
            :icon="viewMode === 'list' ? Expand : List"
            size="small"
            @click="handleViewModeChange"
          />
        </el-tooltip>
        <el-tooltip ref="tooltipRef" content="刷新设备列表" placement="top">
          <el-button
            class="action-btn refresh-btn"
            :icon="Refresh"
            size="small"
            :loading="loading"
            @click="handleRefresh"
          />
        </el-tooltip>
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-box {
  margin-bottom: 16px;
}

.search-wrapper {
  display: flex;
  gap: 8px;

  .el-input {
    flex: 1;
  }

  .action-buttons {
    display: flex;
    gap: 0;
    margin-left: auto;

    .action-btn {
      color: var(--el-text-color-regular);
      border-color: transparent;
      background-color: transparent;
      padding: 5px 8px;
      height: 32px;
      width: 32px;
      border-radius: 4px;

      &:hover {
        color: var(--el-color-primary);
        background-color: var(--el-fill-color-light);
      }

      &:focus {
        outline: none;
      }

      :deep(.el-icon) {
        font-size: 16px;
      }

      &.refresh-btn {
        color: var(--el-color-primary);

        &:hover {
          background-color: var(--el-color-primary-light-9);
        }
      }
    }
  }
}
</style>
