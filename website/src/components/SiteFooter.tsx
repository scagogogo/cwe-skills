import { Typography, Space } from 'antd'
import { GithubOutlined } from '@ant-design/icons'

const { Text, Link } = Typography

const SiteFooter: React.FC = () => {
  return (
    <footer
      style={{
        background: '#0c0c0c',
        borderTop: '1px solid rgba(255,255,255,0.06)',
        padding: '32px 48px',
        textAlign: 'center',
      }}
    >
      <Space direction="vertical" size={10} style={{ width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'center', gap: 20 }}>
          <Link
            href="https://github.com/scagogogo/cwe-skills"
            target="_blank"
            style={{ color: 'rgba(255,255,255,0.4)', fontSize: 18 }}
          >
            <GithubOutlined />
          </Link>
        </div>

        <Text style={{ color: 'rgba(255,255,255,0.25)', fontSize: 12 }}>
          CWE Skills — AI 原生的 CWE 集成层
        </Text>

        <Space size={14}>
          <Link
            href="https://cwe.mitre.org/"
            target="_blank"
            style={{ color: 'rgba(255,255,255,0.3)', fontSize: 11 }}
          >
            CWE (MITRE)
          </Link>
          <Link
            href="https://pkg.go.dev/github.com/scagogogo/cwe-skills"
            target="_blank"
            style={{ color: 'rgba(255,255,255,0.3)', fontSize: 11 }}
          >
            Go Doc
          </Link>
          <Link
            href="https://github.com/scagogogo/cwe-skills/tree/main/docs/skills"
            target="_blank"
            style={{ color: 'rgba(255,255,255,0.3)', fontSize: 11 }}
          >
            Skills 文档
          </Link>
        </Space>

        <Text style={{ color: 'rgba(255,255,255,0.15)', fontSize: 11 }}>
          MIT License · Built with React + Ant Design
        </Text>
      </Space>
    </footer>
  )
}

export default SiteFooter
