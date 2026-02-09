<template>
  <div id="app">
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #000000;
  color: #ffffff;
}

#app {
  height: 100%;
}

/* 页面切换动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.4s cubic-bezier(0.55, 0, 0.1, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #1a1a1a;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #f1c40f;
  border-radius: 4px;
  transition: background 0.3s ease;
}

::-webkit-scrollbar-thumb:hover {
  background: #f39c12;
}

/* 输入框自动填充深色模式修复 */
input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 1000px #2a2a2a inset !important;
  -webkit-text-fill-color: #ffffff !important;
  transition: background-color 5000s ease-in-out 0s;
}

/* 全局动态背景动画类 */
.cyber-bg {
  background-color: #050505;
  background-image: 
    linear-gradient(rgba(241, 196, 15, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(241, 196, 15, 0.1) 1px, transparent 1px);
  background-size: 40px 40px;
  position: relative;
  overflow: hidden;
}

/* 扫描波浪效果 */
.cyber-bg::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(241, 196, 15, 0.15), transparent);
  animation: scanWave 8s linear infinite;
  pointer-events: none;
  z-index: 0;
}

@keyframes scanWave {
  0% { left: -50%; }
  100% { left: 150%; }
}

.cyber-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, rgba(241, 196, 15, 0.08) 0%, transparent 60%);
  animation: bgPulse 10s ease-in-out infinite;
  pointer-events: none;
}

@keyframes bgPulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.2); }
}

/* Element Plus 组件样式覆盖 */
:root {
  --el-color-primary: #f1c40f;
  --el-color-primary-light-3: #f39c12;
  --el-color-primary-dark-2: #d4ac0d;
}

.el-button {
  font-weight: 600;
  border-radius: 8px;
}

.el-button--primary {
  background-color: #f1c40f !important;
  border-color: #f1c40f !important;
  color: #1a1a1a !important;
}

.el-button--primary:hover {
  background-color: #f39c12 !important;
  border-color: #f39c12 !important;
  color: #000 !important;
}

.el-button--danger {
  background-color: #1a1a1a !important;
  border-color: #f1c40f !important;
  color: #f1c40f !important;
}

.el-button--danger:hover {
  background-color: #333 !important;
  border-color: #f39c12 !important;
  color: #f39c12 !important;
}

.el-input__wrapper {
  background-color: #2a2a2a !important;
  box-shadow: 0 0 0 1px #444 inset !important;
}

.el-input__wrapper.is-focus {
  box-shadow: 0 0 0 1px #f1c40f inset !important;
}

.el-input__inner {
  color: #fff !important;
}

.el-card {
  background-color: #1a1a1a !important;
  border-color: #333 !important;
  color: #fff !important;
}

.el-message {
  background: #1a1a1a !important;
  border-color: #333 !important;
}

.el-message--success .el-message__content {
  color: #f1c40f !important;
}

.el-message--error .el-message__content {
  color: #e74c3c !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-enter-from,
  .page-leave-to {
    transform: translateX(0);
    opacity: 0;
  }

  .page-enter-active,
  .page-leave-active {
    transition: opacity 0.3s ease;
  }
}
</style>