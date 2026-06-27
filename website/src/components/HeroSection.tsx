import { Typography, Button, Space } from 'antd'
import {
  RocketOutlined,
  SafetyCertificateOutlined,
  ApiOutlined,
} from '@ant-design/icons'
import { asset } from '../constants'

const { Title, Paragraph } = Typography

const HeroSection: React.FC = () => {
  return (
    <section
      id="hero"
      style={{
        padding: '100px 48px 64px',
        textAlign: 'center',
        background: '#0a0a0a',
      }}
    >
      {/* 标签 */}
      <Space size={6} style={{ marginBottom: 20 }}>
        <span
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 4,
            padding: '3px 12px',
            fontSize: 12,
            color: '#1677ff',
            background: 'rgba(22, 119, 255, 0.08)',
            border: '1px solid rgba(22, 119, 255, 0.15)',
            borderRadius: 2,
          }}
        >
          <RocketOutlined style={{ fontSize: 11 }} /> AI-Native
        </span>
        <span
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 4,
            padding: '3px 12px',
            fontSize: 12,
            color: '#52c41a',
            background: 'rgba(82, 196, 26, 0.08)',
            border: '1px solid rgba(82, 196, 26, 0.15)',
            borderRadius: 2,
          }}
        >
          <SafetyCertificateOutlined style={{ fontSize: 11 }} /> Zero Deps
        </span>
        <span
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 4,
            padding: '3px 12px',
            fontSize: 12,
            color: '#fa8c16',
            background: 'rgba(250, 140, 22, 0.08)',
            border: '1px solid rgba(250, 140, 22, 0.15)',
            borderRadius: 2,
          }}
        >
          <ApiOutlined style={{ fontSize: 11 }} /> 40+ Commands
        </span>
      </Space>

      {/* 主标题 */}
      <Title
        level={1}
        style={{
          color: '#fff',
          fontSize: 48,
          fontWeight: 700,
          letterSpacing: '-1.2px',
          marginBottom: 12,
          lineHeight: 1.15,
        }}
      >
        CWE Skills
      </Title>

      <Title
        level={2}
        style={{
          color: 'rgba(255,255,255,0.5)',
          fontSize: 20,
          fontWeight: 400,
          marginTop: 0,
          marginBottom: 10,
        }}
      >
        AI 原生的 CWE 集成层
      </Title>

      <Paragraph
        style={{
          color: 'rgba(255,255,255,0.4)',
          fontSize: 15,
          maxWidth: 640,
          margin: '0 auto 36px',
          lineHeight: 1.8,
        }}
      >
        统一 MITRE REST API、XML 目录与权威列表（Top 25 / OWASP / SANS），
        提供四种接入方式 — Skills、Go SDK、CLI、MCP
      </Paragraph>

      {/* CTA 按钮 */}
      <Space size={12}>
        <Button
          type="primary"
          size="large"
          href="#quickstart"
          style={{
            borderRadius: 3,
            height: 44,
            paddingInline: 28,
            fontSize: 14,
            fontWeight: 500,
          }}
        >
          快速开始
        </Button>
        <Button
          size="large"
          href="https://github.com/scagogogo/cwe-skills"
          target="_blank"
          style={{
            borderRadius: 3,
            height: 44,
            paddingInline: 28,
            fontSize: 14,
            borderColor: 'rgba(255,255,255,0.12)',
          }}
        >
          GitHub
        </Button>
      </Space>

      {/* 特性概览图 */}
      <div style={{ maxWidth: 860, margin: '56px auto 0' }}>
        <img
          src={asset('feature-tree.png')}
          alt="CWE Skills Feature Tree"
          style={{
            width: '100%',
            borderRadius: 3,
            border: '1px solid rgba(255,255,255,0.06)',
          }}
        />
      </div>
    </section>
  )
}

export default HeroSection
