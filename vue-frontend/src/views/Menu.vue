<template>
  <div class="menu-container">
    <el-header class="header">
      <h1>CyrusAI</h1>
      <el-button type="danger" @click="handleLogout">退出登录</el-button>
    </el-header>
    <el-main class="main">
      <div class="menu-grid">
        <el-card class="menu-item" @click="$router.push('/ai-chat')">
          <div class="card-content">
            <el-icon size="48" color="#409eff"><ChatDotRound /></el-icon>
            <h3>AI聊天</h3>
            <p>与AI进行智能对话</p>
          </div>
        </el-card>
        <el-card class="menu-item" @click="$router.push('/image-recognition')">
          <div class="card-content">
            <el-icon size="48" color="#67c23a"><Camera /></el-icon>
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

export default {
  name: 'MenuView',
  components: {
    ChatDotRound,
    Camera
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
  background-color: #050505;
  background-image: 
    linear-gradient(rgba(241, 196, 15, 0.25) 1px, transparent 1px),
    linear-gradient(90deg, rgba(241, 196, 15, 0.25) 1px, transparent 1px);
  background-size: 30px 30px;
  position: relative;
  overflow: hidden;
}

.menu-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, rgba(241, 196, 15, 0.05) 0%, transparent 70%);
  animation: pulse 15s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

.header {
  background: rgba(26, 26, 26, 0.9);
  backdrop-filter: blur(10px);
  color: #f1c40f;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30px;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(241, 196, 15, 0.2);
  position: relative;
  z-index: 2;
}

.header h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 2px;
  text-transform: uppercase;
  color: #f1c40f;
  text-shadow: 0 0 10px rgba(241, 196, 15, 0.3);
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
  background: rgba(26, 26, 26, 0.8) !important;
  backdrop-filter: blur(15px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(241, 196, 15, 0.2) !important;
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
  background: linear-gradient(135deg, rgba(241, 196, 15, 0.1) 0%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.menu-item:hover {
  transform: translateY(-10px);
  border-color: #f1c40f !important;
  box-shadow: 0 15px 40px rgba(241, 196, 15, 0.15);
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
  color: #f1c40f !important;
}

.menu-item:hover .el-icon {
  transform: scale(1.2) rotate(5deg);
  color: #fff !important;
  text-shadow: 0 0 15px #f1c40f;
}

.card-content h3 {
  margin: 0 0 15px 0;
  color: #fff;
  font-size: 24px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.menu-item:hover h3 {
  color: #f1c40f;
}

.card-content p {
  margin: 0;
  color: #999;
  font-size: 16px;
  line-height: 1.6;
  transition: all 0.3s ease;
}

.menu-item:hover p {
  color: #ccc;
}
</style>