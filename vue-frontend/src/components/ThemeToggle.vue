<template>
  <el-button 
    circle 
    class="theme-toggle" 
    @click="toggleTheme"
    :title="isDark ? '切换到亮色模式' : '切换到暗色模式'"
  >
    <el-icon :size="20">
      <component :is="isDark ? Moon : Sunny" />
    </el-icon>
  </el-button>
</template>

<script setup>
import { ref, onMounted } from 'vue'
// Icons are used in template by component :is="...", so they are technically used
// but ESLint might not detect it if not directly referenced in script.
// However, since we are using <component :is="..."> with string names 'Moon' and 'Sunny',
// we need to register them or make them available.
// In <script setup>, imported components are automatically available.
// If ESLint complains unused, it might be because we pass string 'Moon'/'Sunny' instead of the component object itself.
import { Moon, Sunny } from '@element-plus/icons-vue'

const isDark = ref(true)


const updateTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'dark'
  isDark.value = savedTheme === 'dark'
  document.documentElement.setAttribute('data-theme', savedTheme)
}

onMounted(() => {
  updateTheme()
})

const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  localStorage.setItem('theme', newTheme)
  updateTheme()
}
</script>

<style scoped>
.theme-toggle {
  background: transparent !important;
  border: none !important;
  color: var(--accent-color) !important;
  transition: all 0.3s;
  z-index: 1000;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle:hover {
  background: var(--bg-secondary) !important;
  color: var(--accent-hover) !important;
  transform: rotate(15deg);
}
</style>
