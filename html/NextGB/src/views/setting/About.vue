<script setup lang="ts">
import { Check, Link } from '@element-plus/icons-vue'

const aboutInfo = {
  name: 'NextGB',
  version: 'v1.0.0',
  description: '下一代国标平台，支持GB/T 28181-2016标准，同时支持SRS和ZLMediaKit',
  background: {
    title: '与 AI 共舞',
    sections: [
      {
        subtitle: '',
        content: '本项目是一次探索性实验，代码通过与Cursor AI助手对话生成。项目旨在验证当前AI技术在实际项目开发中的能力和潜力。从项目启动到完成基本功能，总计投入约40个小时，包括UI，架构，甚至目录结构。'
      },
      {
        subtitle: '',
        content: '对于我这样一个后端开发者而言，UI设计几乎是不可能完成的任务，而现在，竟成了对话框里的文字。"要一个左边设备树、右边视频墙的布局"，当这句话变成代码在浏览器上显示时，我竟有些兴奋。从"八方向的云台控制面板"到"能拖拽缩放的时间轴"，AI像一个魔术师，不断的给我惊喜。有时，我甚至不知如何描述，只会输入"帮我优化样式"，然后就是等待惊喜的出现。'
      },
      {
        subtitle: '',
        content: '看着这样一件由对话框诞生的作品，我有一种"作弊"的兴奋感，同时还有一些恍惚。坦白讲，我不知道这将会是替代，还是开发者能力的延伸。毫无疑问的是，它将推倒领域间的壁垒，而我们，也将成为这一切的参与者，抑或见证者。'
      },
      {
        subtitle: '',
        content: ''
      }
    ]
  }
}
</script>

<template>
  <div class="about-container">
    <div class="about-header">
      <el-image class="logo" src="/logo.svg" fit="contain" />
      <div class="title-section">
        <h1>{{ aboutInfo.name }}</h1>
        <div class="version-tag">{{ aboutInfo.version }}</div>
      </div>
    </div>

    <div class="about-content">
      <div class="description">
        {{ aboutInfo.description }}
      </div>
    </div>

    <div class="author-container">
      <div class="divider"></div>
      <div class="author-section">
        <div class="author-info">
          <div class="background-content">
            <h4 class="story-title">{{ aboutInfo.background.title }}</h4>
            <div v-for="(section, index) in aboutInfo.background.sections" 
                 :key="index" 
                 class="story-section">
              <h5 v-if="section.subtitle" class="section-subtitle">{{ section.subtitle }}</h5>
              <p class="section-content">{{ section.content }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.about-container {
  max-width: 900px;
  margin: 32px auto;
  padding: 0 24px;
}

.about-header {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 48px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.logo {
  width: 90px;
  height: 90px;
  filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.1));
}

.title-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.title-section h1 {
  margin: 0;
  font-size: 36px;
  font-weight: 600;
  background: linear-gradient(120deg, var(--el-color-primary), var(--el-color-primary-light-3));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.version-tag {
  background-color: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(var(--el-color-primary-rgb), 0.1);
}

.about-content {
  color: var(--el-text-color-regular);
}

.description {
  font-size: 18px;
  line-height: 1.7;
  margin-bottom: 40px;
  color: var(--el-text-color-secondary);
}

.author-container {
  margin-top: 60px;
}

.divider {
  height: 1px;
  background: linear-gradient(90deg, 
    transparent, 
    var(--el-border-color-lighter) 20%,
    var(--el-border-color-lighter) 80%,
    transparent
  );
  margin-bottom: 32px;
}

.author-section {
  padding: 32px;
  background: linear-gradient(145deg, var(--el-bg-color), var(--el-fill-color-light));
  border-radius: 16px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
}

.author-section h3 {
  font-size: 24px;
  margin-bottom: 20px;
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.author-info {
  margin-top: 24px;
}

.background-content {
  margin-top: 0;
}

.story-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--el-color-primary);
  margin-bottom: 24px;
}

.story-section {
  margin-bottom: 24px;
}

.section-subtitle {
  font-size: 18px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  margin-bottom: 12px;
}

.section-content {
  font-size: 16px;
  line-height: 1.8;
  color: var(--el-text-color-secondary);
  text-align: justify;
  margin-bottom: 16px;
}

.story-section:last-child .section-content {
  margin-bottom: 0;
}
</style> 