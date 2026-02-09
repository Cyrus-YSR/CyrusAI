<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>登录</h2>
        </div>
      </template>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            placeholder="请输入密码"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            type="text"
            @click="$router.push('/register')"
            style="width: 100%"
          >
            还没有账号？去注册
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '../utils/api'

export default {
  name: 'LoginView',
  setup() {
    const router = useRouter()
    const loginFormRef = ref()
    const loading = ref(false)
    const loginForm = ref({
      username: '',
      password: ''
    })

    const loginRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ]
    }

    const handleLogin = async () => {
      try {
        await loginFormRef.value.validate()
        loading.value = true
        const response = await api.post('/user/login', {
          username: loginForm.value.username,
          password: loginForm.value.password
        })
        if (response.data.status_code === 1000) {
          localStorage.setItem('token', response.data.token)
          ElMessage.success('登录成功')
          router.push('/menu')
        } else {
          ElMessage.error(response.data.status_msg || '登录失败')
        }
      } catch (error) {
        console.error('Login error:', error)
        ElMessage.error('登录失败，请重试')
      } finally {
        loading.value = false
      }
    }

    return {
      loginFormRef,
      loading,
      loginForm,
      loginRules,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  /* Nano Banana 动态背景 */
  background-color: #050505;
  background-image: 
    linear-gradient(rgba(241, 196, 15, 0.25) 1px, transparent 1px),
    linear-gradient(90deg, rgba(241, 196, 15, 0.25) 1px, transparent 1px);
  background-size: 40px 40px;
  position: relative;
  overflow: hidden;
}

/* 扫描波浪效果 */
.login-container::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(241, 196, 15, 0.4), transparent);
  animation: scanWave 8s linear infinite;
  pointer-events: none;
  z-index: 0;
}

@keyframes scanWave {
  0% { left: -50%; }
  100% { left: 150%; }
}

.login-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(241, 196, 15, 0.08) 0%, transparent 60%);
  animation: pulse 10s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.5; }
  50% { transform: scale(1.2); opacity: 0.8; }
}

.login-card {
  width: 420px;
  background: rgba(26, 26, 26, 0.8) !important;
  backdrop-filter: blur(15px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(241, 196, 15, 0.2) !important;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  animation: slideIn 0.8s ease-out;
  position: relative;
  z-index: 1;
  overflow: hidden;
}

.login-card::before {
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

.login-card:hover {
  transform: translateY(-10px);
  border-color: #f1c40f !important;
  box-shadow: 0 15px 40px rgba(241, 196, 15, 0.15);
}

.login-card:hover::before {
  opacity: 1;
}

:deep(.el-card__header) {
  border-bottom: none !important;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.card-header {
  text-align: center;
  padding: 30px 0 20px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.card-header h2 {
  margin: 0;
  color: #f1c40f;
  font-size: 28px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2px;
  text-shadow: 0 0 10px rgba(241, 196, 15, 0.3);
}

.el-form-item {
  margin-bottom: 24px;
}

:deep(.el-form-item__label) {
  color: #f1c40f !important;
  font-weight: 600;
}

.el-input {
  transition: all 0.3s ease;
}

.el-input:focus-within {
  transform: scale(1.02);
}

.el-button {
  height: 48px;
  border-radius: 12px;
  font-weight: 700;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.el-button--primary {
  background: #f1c40f;
  border: none;
  color: #000;
  box-shadow: 0 4px 15px rgba(241, 196, 15, 0.3);
}

.el-button--primary:hover {
  background: #f39c12;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(241, 196, 15, 0.5);
}

.el-button--text {
  color: #aaa;
}

.el-button--text:hover {
  color: #f1c40f;
}

</style>