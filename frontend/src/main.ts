import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import naive from 'naive-ui'
import './style.css';

const meta = document.createElement('meta')
meta.name = 'naive-ui-style'
document.head.appendChild(meta)

const app = createApp(App)
app.use(naive)
app.use(router)
app.mount('#app')
