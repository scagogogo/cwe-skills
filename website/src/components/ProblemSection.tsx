import { Typography, Row, Col } from 'antd'
import {
  CloudServerOutlined,
  ThunderboltOutlined,
  RobotOutlined,
  CodeOutlined,
  WifiOutlined,
} from '@ant-design/icons'
import { asset } from '../constants'

const { Title, Paragraph, Text } = Typography

interface ProblemItem {
  icon: React.ReactNode
  challenge: string
  solution: string
  color: string
}

const problems: ProblemItem[] = [
  {
    icon: <CloudServerOutlined />,
    challenge: 'CWE 数据分散在多个来源',
    solution: '统一集成层，一个入口访问所有 CWE 数据',
    color: '#1677ff',
  },
  {
    icon: <ThunderboltOutlined />,
    challenge: 'MITRE API 速率限制、响应慢',
    solution: 'Go SDK 内建速率控制、重试、类型化',
    color: '#52c41a',
  },
  {
    icon: <RobotOutlined />,
    challenge: 'AI 代理无法原生查询 CWE',
    solution: 'Skills — 复制粘贴即可获得查询能力',
    color: '#fa8c16',
  },
  {
    icon: <CodeOutlined />,
    challenge: 'Shell 脚本需要 CWE 验证',
    solution: 'CLI — 40+ 子命令，文本/JSON 双格式',
    color: '#13c2c2',
  },
  {
    icon: <WifiOutlined />,
    challenge: '离线场景无法访问 CWE 数据',
    solution: 'XML 解析器 — 完全离线的搜索和导航',
    color: '#eb2f96',
  },
]

const ProblemSection: React.FC = () => {
  return (
    <section
      id="problem"
      style={{
        padding: '64px 48px',
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
            解决什么问题
          </Title>
          <Paragraph
            style={{
              color: 'rgba(255,255,255,0.4)',
              fontSize: 15,
              maxWidth: 520,
              margin: '0 auto',
            }}
          >
            CWE 数据分散在各处，CWE Skills 将这一切统一到一个集成层
          </Paragraph>
        </div>

        {/* 数据流架构图 */}
        <div style={{ maxWidth: 860, margin: '0 auto 40px' }}>
          <img
            src={asset('data-flow.png')}
            alt="Data Flow"
            style={{
              width: '100%',
              borderRadius: 3,
              border: '1px solid rgba(255,255,255,0.06)',
            }}
          />
        </div>

        {/* 问题-解决方案卡片 */}
        <Row gutter={[16, 16]}>
          {problems.map((item, index) => (
            <Col xs={24} sm={12} lg={8} key={index}>
              <div
                style={{
                  padding: 20,
                  background: '#111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  borderRadius: 3,
                  height: '100%',
                }}
              >
                <div
                  style={{
                    width: 36,
                    height: 36,
                    borderRadius: 3,
                    background: `${item.color}12`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    marginBottom: 14,
                    fontSize: 17,
                    color: item.color,
                  }}
                >
                  {item.icon}
                </div>
                <Text
                  style={{
                    display: 'block',
                    color: 'rgba(255,255,255,0.5)',
                    fontSize: 13,
                    marginBottom: 6,
                  }}
                >
                  {item.challenge}
                </Text>
                <Text
                  strong
                  style={{
                    display: 'block',
                    color: '#fff',
                    fontSize: 14,
                  }}
                >
                  → {item.solution}
                </Text>
              </div>
            </Col>
          ))}
        </Row>
      </div>
    </section>
  )
}

export default ProblemSection
