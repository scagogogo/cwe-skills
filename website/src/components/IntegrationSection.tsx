import { Typography, Row, Col, Space } from 'antd'
import {
  RobotOutlined,
  CodeOutlined,
  CodeSandboxOutlined,
  ApiOutlined,
  CopyOutlined,
  DownloadOutlined,
  ToolOutlined,
} from '@ant-design/icons'

const { Title, Paragraph, Text } = Typography

interface IntegrationMethod {
  num: number
  icon: React.ReactNode
  title: string
  subtitle: string
  description: string
  setup: string
  tags: string[]
  color: string
}

const methods: IntegrationMethod[] = [
  {
    num: 1,
    icon: <RobotOutlined />,
    title: 'Skills',
    subtitle: 'AI 代理接入',
    description:
      '复制粘贴一段提示词到 AI 代理系统提示词中，即可获得完整的 CWE 查询能力。',
    setup: '复制下方提示词',
    tags: ['Claude', 'GPT', 'AI Agent'],
    color: '#1677ff',
  },
  {
    num: 2,
    icon: <CodeOutlined />,
    title: 'Go SDK',
    subtitle: 'Go 应用和库',
    description:
      '类型安全的 Go SDK，内建速率控制和重试。核心 SDK 零依赖，仅使用 Go 标准库。',
    setup: 'go get github.com/scagogogo/cwe-skills',
    tags: ['Typed', 'Rate Limit', 'Zero Deps'],
    color: '#52c41a',
  },
  {
    num: 3,
    icon: <CodeSandboxOutlined />,
    title: 'CLI',
    subtitle: 'Shell 脚本和开发工作流',
    description:
      '40+ 子命令覆盖 CWE ID 操作、API 查询、离线搜索过滤、关系导航、树构建等。',
    setup: '从 Releases 下载',
    tags: ['40+ Commands', 'JSON Output', '30+ Platforms'],
    color: '#fa8c16',
  },
  {
    num: 4,
    icon: <ApiOutlined />,
    title: 'MCP',
    subtitle: 'MCP 兼容 AI 工具',
    description:
      'MCP Server 模式，支持 MCP 协议的 AI 工具可直接调用 CWE 查询能力。',
    setup: '即将推出',
    tags: ['MCP Protocol', 'Coming Soon'],
    color: '#13c2c2',
  },
]

const IntegrationSection: React.FC = () => {
  return (
    <section
      id="integration"
      style={{
        padding: '64px 48px',
        background: '#0c0c0c',
      }}
    >
      <div style={{ maxWidth: 1120, margin: '0 auto' }}>
        <div style={{ textAlign: 'center', marginBottom: 40 }}>
          <Title
            level={2}
            style={{
              color: '#fff',
              fontSize: 32,
              fontWeight: 600,
              marginBottom: 8,
            }}
          >
            四种接入方式
          </Title>
          <Paragraph
            style={{
              color: 'rgba(255,255,255,0.4)',
              fontSize: 15,
              maxWidth: 520,
              margin: '0 auto',
            }}
          >
            从 AI 代理到 Shell 脚本，选择最适合你的方式
          </Paragraph>
        </div>

        <Row gutter={[16, 16]}>
          {methods.map((method) => (
            <Col xs={24} sm={12} key={method.num}>
              <div
                style={{
                  padding: 24,
                  background: '#111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  borderRadius: 3,
                  height: '100%',
                }}
              >
                {/* 标题行 */}
                <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 14 }}>
                  <div
                    style={{
                      width: 40,
                      height: 40,
                      borderRadius: 3,
                      background: `${method.color}12`,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontSize: 20,
                      color: method.color,
                      flexShrink: 0,
                    }}
                  >
                    {method.icon}
                  </div>
                  <div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span
                        style={{
                          fontSize: 11,
                          color: method.color,
                          background: `${method.color}12`,
                          padding: '1px 8px',
                          borderRadius: 2,
                          fontWeight: 600,
                        }}
                      >
                        #{method.num}
                      </span>
                      <Text strong style={{ color: '#fff', fontSize: 17, fontWeight: 600 }}>
                        {method.title}
                      </Text>
                    </div>
                    <Text style={{ color: 'rgba(255,255,255,0.4)', fontSize: 12 }}>
                      {method.subtitle}
                    </Text>
                  </div>
                </div>

                {/* 描述 */}
                <Paragraph
                  style={{
                    color: 'rgba(255,255,255,0.55)',
                    fontSize: 13,
                    lineHeight: 1.8,
                    marginBottom: 14,
                  }}
                >
                  {method.description}
                </Paragraph>

                {/* 安装命令 */}
                <div
                  style={{
                    background: '#0a0a0a',
                    borderRadius: 2,
                    padding: '8px 14px',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 8,
                    marginBottom: 12,
                    border: '1px solid rgba(255,255,255,0.04)',
                  }}
                >
                  {method.num === 1 ? (
                    <CopyOutlined style={{ color: method.color, fontSize: 12 }} />
                  ) : method.num === 3 ? (
                    <DownloadOutlined style={{ color: method.color, fontSize: 12 }} />
                  ) : method.num === 4 ? (
                    <ToolOutlined style={{ color: method.color, fontSize: 12 }} />
                  ) : (
                    <CodeOutlined style={{ color: method.color, fontSize: 12 }} />
                  )}
                  <Text
                    code
                    style={{
                      color: 'rgba(255,255,255,0.6)',
                      fontSize: 12,
                      background: 'transparent',
                      border: 'none',
                      padding: 0,
                    }}
                  >
                    {method.setup}
                  </Text>
                </div>

                {/* 标签 */}
                <Space size={4} wrap>
                  {method.tags.map((tag) => (
                    <span
                      key={tag}
                      style={{
                        padding: '2px 8px',
                        fontSize: 11,
                        color: 'rgba(255,255,255,0.4)',
                        background: 'rgba(255,255,255,0.03)',
                        border: '1px solid rgba(255,255,255,0.06)',
                        borderRadius: 2,
                      }}
                    >
                      {tag}
                    </span>
                  ))}
                </Space>
              </div>
            </Col>
          ))}
        </Row>

        {/* 架构图 */}
        <div style={{ maxWidth: 860, margin: '40px auto 0' }}>
          <img
            src="/cwe-skills/architecture.png"
            alt="Architecture"
            style={{
              width: '100%',
              borderRadius: 3,
              border: '1px solid rgba(255,255,255,0.06)',
            }}
          />
        </div>
      </div>
    </section>
  )
}

export default IntegrationSection
