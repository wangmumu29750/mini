import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import './style.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

const authStore = useAuthStore()
void authStore.refreshCurrentUser().catch(() => {
  authStore.clearSession()
})

app.mount('#app')
