<template>
  <div class="image-recognition-container">
    <!-- 左侧会话列表 -->
    <div class="session-list">
      <div class="session-list-header">
        <span>图像识别</span>
      </div>
      <ul class="session-list-ul">
        <li class="session-item active">
          图像识别助手
        </li>
      </ul>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-section">
      <div class="top-bar">
        <div class="left-group">
          <button class="back-btn" @click="$router.push('/menu')">← 返回</button>
          <h2>AI 图像识别助手</h2>
        </div>
        <ThemeToggle />
      </div>
      
      <div class="chat-messages" ref="chatContainerRef">
        <div
          v-for="(message, index) in messages"
          :key="index"
          :class="['message', message.role === 'user' ? 'user-message' : 'ai-message']"
        >
          <div class="message-header">
            <b>{{ message.role === 'user' ? '你' : 'AI' }}:</b>
          </div>
          <div class="message-content">
            <span>{{ message.content }}</span>
            <img v-if="message.imageUrl" :src="message.imageUrl" alt="上传的图片" />
          </div>
        </div>
      </div>

      <div class="chat-input">
        <form @submit.prevent="handleSubmit">
          <input
            ref="fileInputRef"
            type="file"
            accept="image/*"
            required
            @change="handleFileSelect"
          />
          <button type="submit" :disabled="!selectedFile">发送图片</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick } from 'vue'
import api from '../utils/api'
import ThemeToggle from '../components/ThemeToggle.vue'

export default {
  name: 'ImageRecognition',
  components: { ThemeToggle },
  setup() {
    const messages = ref([])
    const selectedFile = ref(null)
    const fileInputRef = ref()
    const chatContainerRef = ref()

    const handleFileSelect = (event) => {
      selectedFile.value = event.target.files[0]
    }

    const handleSubmit = async () => {
      if (!selectedFile.value) return

      const file = selectedFile.value
      const imageUrl = URL.createObjectURL(file)

      // Add user message to UI
      messages.value.push({
        role: 'user',
        content: `已上传图片: ${file.name}`,
        imageUrl: imageUrl,
      })

      await nextTick()
      scrollToBottom()

      // Create FormData
      const formData = new FormData()
      formData.append('image', file)

      try {
        const response = await api.post('/image/analyze', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })


        if (response.data && response.data.class_name) {
             let aiText = `识别结果: ${response.data.class_name}`
             if (response.data.analysis_text) {
                 aiText += `\n\n详细解析:\n${response.data.analysis_text}`
             }
            messages.value.push({
                role: 'assistant',
                content: aiText,
            })
        } else {
             messages.value.push({
                 role: 'assistant',
                 content: `[错误] ${response.data.status_msg || '识别失败'}`,
             })
        }
      } catch (error) {
        console.error('Upload error:', error)
        messages.value.push({
          role: 'assistant',
          content: `[错误] 无法连接到服务器或上传失败: ${error.message}`,
        })
      } finally {

        URL.revokeObjectURL(imageUrl)

            await nextTick()
        scrollToBottom()


        selectedFile.value = null
        if (fileInputRef.value) {
          fileInputRef.value.value = ''
        }
      }
    }

    const scrollToBottom = () => {
      if (chatContainerRef.value) {
        chatContainerRef.value.scrollTop = chatContainerRef.value.scrollHeight
      }
    }

    return {
      messages,
      selectedFile,
      fileInputRef,
      chatContainerRef,
      handleFileSelect,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.image-recognition-container {
  height: 100vh;
  display: flex;
  background-color: var(--bg-cyber);
  background-image: 
    linear-gradient(var(--cyber-grid) 1px, transparent 1px),
    linear-gradient(90deg, var(--cyber-grid) 1px, transparent 1px);
  background-size: 40px 40px;
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial;
  color: var(--text-primary);
  padding: 20px;
  gap: 20px;
  box-sizing: border-box;
  transition: background-color 0.3s;
}

/* 进场动画定义 */
@keyframes slideInLeft {
  from { opacity: 0; transform: translateX(-50px); }
  to { opacity: 1; transform: translateX(0); }
}

@keyframes slideInDown {
  from { opacity: 0; transform: translateY(-50px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideInUp {
  from { opacity: 0; transform: translateY(50px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 扫描波浪效果 */
.image-recognition-container::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, var(--cyber-wave), transparent);
  animation: scanWave 8s linear infinite;
  pointer-events: none;
  z-index: 0;
}

@keyframes scanWave {
  0% { left: -50%; }
  100% { left: 150%; }
}

/* 动态背景光晕 */
.image-recognition-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at 50% 50%, var(--cyber-pulse) 0%, transparent 60%);
  animation: bgPulse 20s ease-in-out infinite;
  pointer-events: none;
  z-index: 0;
}

@keyframes bgPulse {
  0%, 100% { transform: scale(1); opacity: 0.4; }
  50% { transform: scale(1.1); opacity: 0.6; }
}

.session-list {
  width: 280px;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  box-shadow: 0 10px 30px var(--shadow-color);
  position: relative;
  z-index: 2;
  animation: slideInLeft 1.0s ease-out;
}

.session-list-header {
  padding: 20px;
  text-align: center;
  font-weight: 600;
  background: transparent;
  border-bottom: 1px solid var(--border-color);
  color: var(--accent-color);
}

.session-list-ul {
  list-style: none;
  padding: 0;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.session-item {
  padding: 15px 20px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.2s ease;
  position: relative;
  color: var(--text-regular);
}

.session-item.active {
  background: rgba(var(--accent-rgb), 0.15);
  color: var(--accent-color);
  font-weight: 600;
  box-shadow: inset 4px 0 0 var(--accent-color);
}

.chat-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  box-shadow: 0 10px 30px var(--shadow-color);
  animation: slideInDown 1.0s ease-out; /* 聊天框从上方渐入 */
}

.top-bar {
  background: var(--header-bg);
  backdrop-filter: blur(10px);
  color: var(--text-primary);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  box-shadow: 0 2px 14px var(--shadow-color);
  border-bottom: 1px solid var(--border-color);
  gap: 12px;
}

.left-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  background: var(--bg-glass-light);
  border: 1px solid rgba(var(--accent-rgb), 0.3);
  color: var(--accent-color);
  padding: 8px 14px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px var(--glass-shadow);
}

.back-btn:hover {
  background: var(--bg-glass);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px var(--glass-shadow);
}

.top-bar h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--accent-color);
}

.chat-messages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 30px;
  padding-bottom: 120px; /* 为悬浮输入框留出空间 */
  display: flex;
  flex-direction: column;
  gap: 18px;
  position: relative;
  z-index: 1;
}

/* scrollbar */
.chat-messages::-webkit-scrollbar {
  width: 8px;
}
.chat-messages::-webkit-scrollbar-thumb {
  background: var(--scrollbar-thumb);
  border-radius: 8px;
}
.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.message {
  max-width: 70%;
  padding: 14px 18px;
  border-radius: 18px;
  line-height: 1.6;
  word-wrap: break-word;
  position: relative;
  animation: messageSlideIn 0.28s ease-out;
  font-size: 15px;
  box-sizing: border-box;
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.user-message {
  align-self: flex-end;
  background: var(--accent-color);
  color: var(--text-inverse);
  border: 1px solid var(--accent-color);
  box-shadow: 0 6px 20px var(--glass-shadow);
}

.user-message::after {
  content: '';
  position: absolute;
  bottom: -6px;
  right: 18px;
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-top: 8px solid rgba(241, 196, 15, 0.2);
}

.ai-message {
  align-self: flex-start;
  background: var(--bg-secondary);
  backdrop-filter: blur(4px);
  color: var(--text-primary);
  box-shadow: 0 6px 20px var(--shadow-color);
  border: 1px solid var(--border-color);
}

.ai-message::after {
  content: '';
  position: absolute;
  bottom: -6px;
  left: 18px;
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-top: 8px solid var(--bg-secondary);
}

.message-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  color: var(--accent-color);
}

.message-header b {
  font-weight: 600;
}

/* message content */
.message-content {
  white-space: pre-wrap;
  word-break: break-word;
}

.message-content img {
  max-width: 250px;
  border-radius: 12px;
  display: block;
  margin-top: 12px;
  box-shadow: 0 4px 15px var(--shadow-color);
  border: 1px solid rgba(241, 196, 15, 0.2);
  transition: all 0.3s ease;
}

.message-content img:hover {
  transform: scale(1.05);
  border-color: #f1c40f;
}

/* input area */
.chat-input-container {
  position: absolute;
  bottom: 20px;
  left: 320px; /* 280px (session-list) + 20px (gap) + 20px (padding) */
  right: 20px;
  z-index: 10;
  animation: slideInUp 1.0s ease-out; /* 输入框从下方渐入 */
}

.chat-input {
  padding: 20px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  position: relative;
  box-shadow: 0 10px 30px var(--shadow-color);
}

.chat-input form {
  display: flex;
  gap: 20px;
}

.chat-input input[type="file"] {
  flex: 1;
  border: 1px dashed var(--border-color);
  border-radius: 12px;
  padding: 15px 20px;
  background: var(--bg-input);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 14px;
}

.chat-input input[type="file"]:hover {
  border-color: var(--accent-color);
  background: rgba(241, 196, 15, 0.05);
}

.chat-input input[type="file"]::file-selector-button {
  border: none;
  background: rgba(241, 196, 15, 0.2);
  border: 1px solid rgba(241, 196, 15, 0.3);
  padding: 8px 16px;
  border-radius: 8px;
  color: #f1c40f;
  cursor: pointer;
  font-weight: 600;
  margin-right: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.chat-input input[type="file"]::file-selector-button:hover {
  transform: translateY(-2px);
  background: rgba(241, 196, 15, 0.3);
  box-shadow: 0 4px 15px rgba(241, 196, 15, 0.2);
}

.chat-input button {
  padding: 15px 30px;
  border: none;
  border-radius: 12px;
  background: rgba(241, 196, 15, 0.8);
  color: #000;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(241, 196, 15, 0.2);
}

.chat-input button:hover {
  transform: translateY(-2px);
  background: #f1c40f;
  box-shadow: 0 8px 25px rgba(241, 196, 15, 0.4);
}

.chat-input button:disabled {
  background: var(--border-color);
  color: var(--text-regular);
  box-shadow: none;
  cursor: not-allowed;
  transform: none;
}
</style>
