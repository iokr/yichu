import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './styles/tokens.css'
import './styles/element-override.scss'
import './styles/main.css'
import App from './App.vue'
import router from './router'
import { useMetaStore } from './stores/meta'
import { useDictStore } from './stores/dict'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus, { locale: zhCn })

app.mount('#app')

const meta = useMetaStore(pinia)
const dict = useDictStore(pinia)

void meta.load()
void dict.ensure()

window.addEventListener('yichu:unauthorized', () => {
  const current = router.currentRoute.value
  if (current.path === '/login') return
  void router.replace({ path: '/login', query: current.fullPath === '/' ? {} : { redirect: current.fullPath } })
})
