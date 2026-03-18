<template>
  <div class="image-recognition-container">
    <!-- Glass Panel Wrapper -->
    <div class="glass-wrapper">
      <!-- 左侧会话列表 -->
      <div class="session-list">
        <div class="session-list-header">
          <span>
            <LogoEye :size="24" class="header-logo" />
            图像识别
          </span>
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
        </div>
        
        <div class="chat-messages" ref="chatContainerRef">
          <!-- Empty State -->
          <div v-if="messages.length === 0" class="empty-state">
             <div class="empty-icon">
                <LogoEye :size="60" />
             </div>
             <h3>Upload an image to analyze</h3>
             <p>I can identify objects, read text, and describe scenes.</p>
          </div>

          <div
            v-for="(message, index) in messages"
            :key="index"
            :class="['message', message.role === 'user' ? 'user-message' : 'ai-message']"
          >
            <div class="message-header">
              <span v-if="message.role === 'user'" class="avatar-u">U</span>
              <LogoEye v-else :size="20" />
              <b>{{ message.role === 'user' ? '你' : 'AI' }}</b>
            </div>
            <div class="message-content">
              <span>{{ message.content }}</span>
              <img v-if="message.imageUrl" :src="message.imageUrl" alt="上传的图片" />
            </div>
          </div>
        </div>

        <div class="chat-input">
          <form @submit.prevent="handleSubmit">
            <div class="input-wrapper">
                <input
                    ref="fileInputRef"
                    type="file"
                    accept="image/*"
                    required
                    @change="handleFileSelect"
                    id="file-upload"
                    class="file-input"
                />
                <label for="file-upload" class="file-label">
                    <span v-if="!selectedFile">Choose Image...</span>
                    <span v-else>{{ selectedFile.name }}</span>
                </label>
                <button type="submit" :disabled="!selectedFile" class="send-btn">
                    发送图片
                </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick } from 'vue'
import api from '../utils/api'
import LogoEye from '../components/LogoEye.vue'

export default {
  name: 'ImageRecognition',
  components: { LogoEye },
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

        // Note: keeping objectURL valid for preview in chat history
        // In a real app you might want to upload to S3 and use that URL
        // URL.revokeObjectURL(imageUrl)

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
  justify-content: center;
  align-items: center;
  font-family: 'Inter', sans-serif;
  color: #111827;
  padding: 40px;
}

.glass-wrapper {
  display: flex;
  width: 100%;
  max-width: 1200px;
  height: 90vh;
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.session-list {
  width: 280px;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: rgba(243, 244, 246, 0.6);
  border-right: 1px solid rgba(229, 231, 235, 0.6);
  z-index: 2;
}

.session-list-header {
  padding: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid rgba(229, 231, 235, 0.6);
}

.session-list-header span {
  font-size: 1.25rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  gap: 10px;
  background: linear-gradient(135deg, #EC4899 0%, #4F46E5 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  /* Reset for LogoEye */
}
.header-logo {
    -webkit-text-fill-color: initial;
}


.session-list-ul {
  list-style: none;
  padding: 12px;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.session-item {
  padding: 12px 16px;
  margin-bottom: 4px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
  color: #4B5563;
  border: 1px solid transparent;
}

.session-item.active {
  background: rgba(253, 242, 248, 0.8); /* Pink-50 */
  color: #DB2777; /* Pink-600 */
  font-weight: 600;
  border-color: #FCE7F3;
}

.chat-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.4);
  position: relative;
}

.top-bar {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  color: #1F2937;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  border-bottom: 1px solid rgba(229, 231, 235, 0.6);
}

.left-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid #E5E7EB;
  color: #6B7280;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
  font-size: 14px;
}

.back-btn:hover {
  border-color: #4F46E5;
  color: #4F46E5;
  background: #fff;
}

.top-bar h2 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1F2937;
}

.chat-messages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 30px 15%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9CA3AF;
  gap: 16px;
}

.empty-icon {
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.empty-state h3 {
    font-size: 1.5rem;
    color: #374151;
    margin: 0;
}

.message {
  max-width: 80%;
  padding: 16px 24px;
  border-radius: 16px;
  line-height: 1.6;
  word-wrap: break-word;
  font-size: 15px;
  box-sizing: border-box;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.user-message {
  align-self: flex-end;
  background: #4F46E5;
  color: white;
  border-bottom-right-radius: 4px;
}

.ai-message {
  align-self: flex-start;
  background: rgba(255, 255, 255, 0.9);
  color: #1F2937;
  border: 1px solid rgba(229, 231, 235, 0.6);
  border-bottom-left-radius: 4px;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  font-size: 12px;
  opacity: 0.9;
  font-weight: 500;
}

.avatar-u {
    width: 20px;
    height: 20px;
    background: rgba(255,255,255,0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
}

.message-content img {
  max-width: 300px;
  border-radius: 12px;
  display: block;
  margin-top: 12px;
  border: 1px solid rgba(255,255,255,0.2);
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.chat-input {
  padding: 24px 15%;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(229, 231, 235, 0.6);
}

.input-wrapper {
    display: flex;
    gap: 12px;
    align-items: center;
}

.file-input {
    display: none;
}

.file-label {
    flex: 1;
    border: 1px dashed #9CA3AF;
    border-radius: 12px;
    padding: 12px 16px;
    background: rgba(255,255,255,0.8);
    color: #6B7280;
    cursor: pointer;
    text-align: center;
    transition: all 0.2s;
}

.file-label:hover {
    border-color: #4F46E5;
    color: #4F46E5;
    background: #fff;
}

.send-btn {
  padding: 12px 28px;
  border: none;
  border-radius: 12px;
  background: #4F46E5;
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.2);
}

.send-btn:hover:not(:disabled) {
  background: #4338CA;
  transform: translateY(-1px);
  box-shadow: 0 6px 8px -1px rgba(79, 70, 229, 0.3);
}

.send-btn:disabled {
  background: #E5E7EB;
  cursor: not-allowed;
  box-shadow: none;
  transform: none;
  color: #9CA3AF;
}
</style>