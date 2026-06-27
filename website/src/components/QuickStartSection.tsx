import { Typography, Row, Col, Tabs, Space } from 'antd'
import {
  CodeSandboxOutlined,
  CodeOutlined,
  CopyOutlined,
  DownloadOutlined,
} from '@ant-design/icons'

const { Title, Paragraph, Text } = Typography

const cliCode = `# 安装 CLI (Linux amd64)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/\\
cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# 基础操作
cwe parse CWE-79 89 cwe-352
cwe validate CWE-79
cwe wellknown check CWE-79

# MITRE API 查询
cwe show CWE-79
cwe relations parents CWE-79

# 离线搜索和过滤
cwe search --xml cwec_latest.xml --keyword Injection
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable

# JSON 输出
cwe parse CWE-79 -o json`

const sdkCode = `import (
    "context"
    "github.com/scagogogo/cwe-skills"
)

// 解析和验证 CWE ID
id, _ := cweskills.ParseCWEID("CWE-79")
if cweskills.IsCWEID("CWE-89") { /* valid */ }

// 查询 MITRE REST API
client := cweskills.NewAPIClient()
defer client.Close()
weakness, _ := client.GetWeakness(context.Background(), 79)

// 本地注册表
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()

// 关系导航
nav := cweskills.NewNavigator(registry)
ancestors := nav.Ancestors(79)
path := nav.ShortestPath(79, 1)

// 知名列表
cweskills.IsInTop25(79)      // true
cweskills.IsInOWASPTop10(79) // true`

const skillsCode = `## CWE Skills

你可以使用 \`cwe\` CLI 工具进行 CWE 操作。

### 安装
\`\`\`bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/\\
download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
\`\`\`

### 核心命令
| 命令 | 功能 |
|------|------|
| cwe parse CWE-79 | 解析 CWE ID |
| cwe validate CWE-79 | 验证格式 |
| cwe show CWE-79 | MITRE API 查询 |
| cwe wellknown check CWE-79 | 检查知名列表 |
| cwe search --xml <file> --keyword ... | 离线搜索 |

### 输出
所有命令支持 \`-o json\` 输出结构化 JSON`

const installOptions = [
  {
    icon: <DownloadOutlined />,
    label: 'Release 下载',
    desc: '预编译二进制',
    color: '#1677ff',
  },
  {
    icon: <CodeOutlined />,
    label: '从源码编译',
    desc: 'go build',
    color: '#52c41a',
  },
  {
    icon: <CopyOutlined />,
    label: '包管理器',
    desc: 'brew / scoop / go install',
    color: '#13c2c2',
  },
]

const CodeBlock: React.FC<{ title: string; code: string }> = ({ title, code }) => (
  <div
    style={{
      background: '#111',
      borderRadius: 3,
      border: '1px solid rgba(255,255,255,0.06)',
      overflow: 'hidden',
    }}
  >
    {/* 窗口标题栏 */}
    <div
      style={{
        padding: '8px 16px',
        borderBottom: '1px solid rgba(255,255,255,0.06)',
        display: 'flex',
        alignItems: 'center',
        gap: 6,
      }}
    >
      <div style={{ width: 8, height: 8, borderRadius: 1, background: '#ff5f57' }} />
      <div style={{ width: 8, height: 8, borderRadius: 1, background: '#febc2e' }} />
      <div style={{ width: 8, height: 8, borderRadius: 1, background: '#28c840' }} />
      <Text style={{ marginLeft: 10, color: 'rgba(255,255,255,0.25)', fontSize: 11 }}>
        {title}
      </Text>
    </div>
    <pre
      style={{
        padding: 16,
        color: 'rgba(255,255,255,0.7)',
        fontSize: 12.5,
        lineHeight: 1.7,
        margin: 0,
        overflowX: 'auto',
      }}
    >
      {code}
    </pre>
  </div>
)

const QuickStartSection: React.FC = () => {
  return (
    <section
      id="quickstart"
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
            快速开始
          </Title>
          <Paragraph
            style={{
              color: 'rgba(255,255,255,0.4)',
              fontSize: 15,
              maxWidth: 520,
              margin: '0 auto',
            }}
          >
            三种安装方式，几分钟即可上手
          </Paragraph>
        </div>

        {/* 安装方式 */}
        <Row gutter={[16, 16]} style={{ marginBottom: 36 }}>
          {installOptions.map((opt) => (
            <Col xs={24} sm={8} key={opt.label}>
              <div
                style={{
                  padding: 20,
                  background: '#111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  borderRadius: 3,
                  textAlign: 'center',
                }}
              >
                <div
                  style={{
                    width: 36,
                    height: 36,
                    borderRadius: 3,
                    background: `${opt.color}12`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 17,
                    color: opt.color,
                    margin: '0 auto 10px',
                  }}
                >
                  {opt.icon}
                </div>
                <Text strong style={{ display: 'block', color: '#fff', fontSize: 14, marginBottom: 3 }}>
                  {opt.label}
                </Text>
                <Text style={{ color: 'rgba(255,255,255,0.4)', fontSize: 12 }}>
                  {opt.desc}
                </Text>
              </div>
            </Col>
          ))}
        </Row>

        {/* 代码示例 */}
        <div style={{ maxWidth: 860, margin: '0 auto' }}>
          <Tabs
            defaultActiveKey="cli"
            centered
            items={[
              {
                key: 'cli',
                label: (
                  <Space size={4}>
                    <CodeSandboxOutlined style={{ fontSize: 13 }} />
                    <span style={{ fontSize: 13 }}>CLI</span>
                  </Space>
                ),
                children: <CodeBlock title="terminal" code={cliCode} />,
              },
              {
                key: 'sdk',
                label: (
                  <Space size={4}>
                    <CodeOutlined style={{ fontSize: 13 }} />
                    <span style={{ fontSize: 13 }}>Go SDK</span>
                  </Space>
                ),
                children: <CodeBlock title="main.go" code={sdkCode} />,
              },
              {
                key: 'skills',
                label: (
                  <Space size={4}>
                    <CopyOutlined style={{ fontSize: 13 }} />
                    <span style={{ fontSize: 13 }}>Skills</span>
                  </Space>
                ),
                children: <CodeBlock title="system-prompt.md" code={skillsCode} />,
              },
            ]}
          />
        </div>

        {/* 支持平台 */}
        <div style={{ textAlign: 'center', marginTop: 40 }}>
          <Text style={{ color: 'rgba(255,255,255,0.3)', fontSize: 12 }}>
            预编译二进制支持 30+ 平台
          </Text>
          <div style={{ display: 'flex', justifyContent: 'center', gap: 6, marginTop: 10, flexWrap: 'wrap' }}>
            {['Linux', 'macOS', 'Windows', 'FreeBSD', 'NetBSD', 'OpenBSD', 'AIX', 'Illumos', 'Solaris'].map(
              (platform) => (
                <span
                  key={platform}
                  style={{
                    padding: '2px 10px',
                    fontSize: 11,
                    color: 'rgba(255,255,255,0.35)',
                    background: 'rgba(255,255,255,0.03)',
                    border: '1px solid rgba(255,255,255,0.06)',
                    borderRadius: 2,
                  }}
                >
                  {platform}
                </span>
              ),
            )}
          </div>
        </div>
      </div>
    </section>
  )
}

export default QuickStartSection
