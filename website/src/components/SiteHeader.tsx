import { Button, Space } from 'antd'
import { GithubOutlined, CodeOutlined } from '@ant-design/icons'

interface NavItem {
  key: string
  label: string
}

const navItems: NavItem[] = [
  { key: 'problem', label: '问题' },
  { key: 'integration', label: '接入方式' },
  { key: 'features', label: '功能特性' },
  { key: 'skills', label: 'Skills' },
  { key: 'quickstart', label: '快速开始' },
]

const SiteHeader: React.FC = () => {
  const handleClick = (key: string) => {
    const el = document.getElementById(key)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }

  return (
    <header
      style={{
        position: 'sticky',
        top: 0,
        zIndex: 100,
        background: 'rgba(10, 10, 10, 0.92)',
        backdropFilter: 'blur(12px)',
        borderBottom: '1px solid rgba(255,255,255,0.06)',
        display: 'flex',
        alignItems: 'center',
        padding: '0 48px',
        height: 56,
      }}
    >
      {/* Logo */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          marginRight: 40,
          cursor: 'pointer',
          flexShrink: 0,
        }}
        onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
      >
        <CodeOutlined style={{ fontSize: 20, color: '#1677ff' }} />
        <span
          style={{
            fontSize: 16,
            fontWeight: 600,
            color: '#fff',
            letterSpacing: '-0.3px',
          }}
        >
          CWE Skills
        </span>
      </div>

      {/* 导航链接 - 始终水平展示，不折叠 */}
      <nav
        style={{
          flex: 1,
          display: 'flex',
          alignItems: 'center',
          gap: 4,
          overflow: 'visible',
          whiteSpace: 'nowrap',
        }}
      >
        {navItems.map((item) => (
          <a
            key={item.key}
            onClick={(e) => {
              e.preventDefault()
              handleClick(item.key)
            }}
            style={{
              color: 'rgba(255,255,255,0.55)',
              fontSize: 14,
              padding: '6px 14px',
              borderRadius: 4,
              transition: 'all 0.2s',
              cursor: 'pointer',
              textDecoration: 'none',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.color = '#fff'
              e.currentTarget.style.background = 'rgba(255,255,255,0.06)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.color = 'rgba(255,255,255,0.55)'
              e.currentTarget.style.background = 'transparent'
            }}
          >
            {item.label}
          </a>
        ))}
      </nav>

      {/* 右侧按钮 */}
      <Space size={8} style={{ flexShrink: 0 }}>
        <Button
          size="small"
          href="https://github.com/scagogogo/cwe-skills"
          target="_blank"
          icon={<GithubOutlined />}
          style={{
            color: 'rgba(255,255,255,0.55)',
            border: '1px solid rgba(255,255,255,0.1)',
            borderRadius: 4,
            background: 'transparent',
          }}
        >
          GitHub
        </Button>
      </Space>
    </header>
  )
}

export default SiteHeader
