/* 主题入口样式 —— 引入自定义样式覆盖默认 VitePress 主题 */
import DefaultTheme from 'vitepress/theme'
import CopyPrompt from './components/CopyPrompt.vue'
import './index.css'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    // 全局注册 AI 提示词一键复制组件
    app.component('CopyPrompt', CopyPrompt)
  },
}
