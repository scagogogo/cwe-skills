# CWE Skills 文档站

基于 [VitePress](https://vitepress.dev/) 构建的 CWE Skills 官方文档站。

## 本地开发

```bash
cd website
npm install
npm run dev      # 启动开发服务器
npm run build    # 构建静态站点到 .vitepress/dist
npm run preview  # 本地预览构建产物
```

## 目录结构

- `index.md` — 首页
- `guide/` — 使用指南（25 篇）
- `sdk/` — Go SDK API 参考（120 篇）
- `cli/` — CLI 命令参考（67 篇）
- `enums/` — 枚举类型参考（10 篇）
- `wellknown/` — 知名列表（5 篇）
- `skills/` — AI Skills 接入（13 篇）
- `examples/` — 实战教程（8 篇）
- `.vitepress/config.ts` — 站点配置
- `.vitepress/theme/` — 自定义主题
- `public/` — 静态资源（架构图等）

## 部署

通过 GitHub Actions（`.github/workflows/deploy-website.yml`）自动构建并部署到 GitHub Pages。
