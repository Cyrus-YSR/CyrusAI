<template>
  <div class="menu-container">
    <el-header class="header">
      <div class="header-brand">
        <h1>CyrusAI</h1>
      </div>
      <div class="header-actions">
        <ThemeToggle />
        <el-button type="danger" @click="handleLogout">退出登录</el-button>
      </div>
    </el-header>
    <el-main class="main">
      <div class="menu-grid">
        <el-card class="menu-item" @click="$router.push('/ai-chat')">
          <div class="card-content">
            <el-icon size="48"><ChatDotRound /></el-icon>
            <h3>AI聊天</h3>
            <p>与AI进行智能对话</p>
          </div>
        </el-card>
        <el-card class="menu-item" @click="$router.push('/image-recognition')">
          <div class="card-content">
            <el-icon size="48"><Camera /></el-icon>
            <h3>图像识别</h3>
            <p>上传图片进行AI识别</p>
          </div>
        </el-card>
      </div>
    </el-main>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatDotRound, Camera } from '@element-plus/icons-vue'
import ThemeToggle from '../components/ThemeToggle.vue'

export default {
  name: 'MenuView',
  components: {
    ChatDotRound,
    Camera,
    ThemeToggle
  },
  setup() {
    const router = useRouter()

    const handleLogout = async () => {
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        localStorage.removeItem('token')
        ElMessage.success('退出登录成功')
        router.push('/login')
      } catch {
        // 用户取消操作
      }
    }

    return {
      handleLogout
    }
  }
}
</script>

<style scoped>
.menu-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  /* Nano Banana 动态背景 */
  background-color: var(--bg-cyber);
  background-image: 
    linear-gradient(var(--cyber-grid) 1px, transparent 1px),
    linear-gradient(90deg, var(--cyber-grid) 1px, transparent 1px);
  background-size: 30px 30px;
  position: relative;
  overflow: hidden;
  transition: background-color 0.3s;
}

.menu-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, var(--cyber-pulse) 0%, transparent 70%);
  animation: pulse 15s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

.header {
  background: var(--header-bg);
  backdrop-filter: blur(10px);
  color: var(--accent-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30px;
  box-shadow: 0 2px 20px var(--shadow-color);
  border-bottom: 1px solid var(--cyber-wave);
  position: relative;
  z-index: 2;
  transition: background-color 0.3s, border-color 0.3s;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 2px;
  text-transform: uppercase;
  color: var(--accent-color);
  text-shadow: 0 0 10px var(--cyber-pulse);
  transition: color 0.3s;
}

.main {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  z-index: 1;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 40px;
  max-width: 900px;
  width: 100%;
  padding: 40px;
  animation: gridFadeIn 1s ease-out;
}

@keyframes gridFadeIn {
  from {
    opacity: 0;
    transform: translateY(50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.menu-item {
  cursor: pointer;
  background: var(--bg-glass) !important;
  backdrop-filter: var(--glass-backdrop);
  border-radius: 20px;
  box-shadow: var(--glass-shadow);
  border: var(--glass-border) !important;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  overflow: hidden;
}

.menu-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(var(--accent-rgb), 0.1) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.menu-item:hover {
  transform: translateY(-10px);
  border-color: var(--accent-color) !important;
  box-shadow: 0 15px 40px rgba(var(--accent-rgb), 0.15);
}

.menu-item:hover::before {
  opacity: 1;
}

.card-content {
  text-align: center;
  padding: 50px 30px;
  position: relative;
  z-index: 1;
}

.el-icon {
  display: block;
  margin: 0 auto 20px;
  transition: all 0.3s ease;
  color: var(--accent-color) !important;
}

.menu-item:hover .el-icon {
  transform: scale(1.2) rotate(5deg);
  color: var(--icon-hover-color) !important;
  text-shadow: 0 0 15px var(--accent-color);
}

.card-content h3 {
  margin: 0 0 15px 0;
  color: var(--text-primary);
  font-size: 24px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.menu-item:hover h3 {
  color: var(--accent-color);
}

.card-content p {
  margin: 0;
  color: var(--text-regular);
  font-size: 16px;
  line-height: 1.6;
  transition: all 0.3s ease;
}

.menu-item:hover p {
  color: var(--text-secondary);
}
</style>