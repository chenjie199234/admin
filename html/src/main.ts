import { createApp } from 'vue'
import app from './app.vue'
import { createVuestic } from 'vuestic-ui'
import 'vuestic-ui/css'

createApp(app).use(createVuestic()).mount('#app')