import { createApp } from 'vue'
import main from './main.vue'
import { createVuestic } from 'vuestic-ui'
import 'vuestic-ui/css'
import Axios from 'axios'
Axios.get()

createApp(main).use(createVuestic()).mount('#app')
