import { createApp } from 'vue'
import main from './main.vue'
import { createVuestic } from 'vuestic-ui'
import 'vuestic-ui/css'

createApp(main).use(createVuestic()).mount('#app')
