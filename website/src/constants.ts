// Vite 的 base 路径，用于 public 目录下的资源引用
// 开发环境: '/cwe-skills/'
// 生产环境: '/cwe-skills/'
export const BASE = import.meta.env.BASE_URL

// 图片资源路径辅助函数
export const asset = (path: string) => `${BASE}${path}`.replace(/\/+/g, '/')
