import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import vue3GoogleLogin from 'vue3-google-login'

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(vue3GoogleLogin, {
  clientId: process.env.VUE_APP_GOOGLE_CLIENT_ID
})
app.mount('#app')