import { Typography, Row, Col, Collapse } from 'antd'
import { CheckCircleOutlined } from '@ant-design/icons'

const { Title, Paragraph, Text } = Typography

interface SkillItem {
  num: number
  title: string
  description: string
  category: 'basic' | 'api' | 'local' | 'sdk'
}

const skills: SkillItem[] = [
  { num: 1, title: 'CWE ID 解析与验证', description: '解析、验证、格式化 CWE ID', category: 'basic' },
  { num: 2, title: 'CWE ID 提取与比较', description: '从文本提取、比较 ID', category: 'basic' },
  { num: 3, title: '知名列表', description: 'CWE Top 25、OWASP Top 10、SANS Top 25', category: 'basic' },
  { num: 4, title: '枚举类型', description: '抽象层级、状态、关系类型', category: 'basic' },
  { num: 5, title: 'API: 获取弱点详情', description: '从 MITRE API 获取', category: 'api' },
  { num: 6, title: 'API: 关系查询', description: '通过 API 查询关系', category: 'api' },
  { num: 7, title: 'API: 版本检查', description: '检查 MITRE API 版本', category: 'api' },
  { num: 8, title: '本地: 搜索与过滤', description: '搜索和多条件过滤', category: 'local' },
  { num: 9, title: '本地: 注册表操作', description: '加载、查询、导出本地数据', category: 'local' },
  { num: 10, title: '本地: 关系导航', description: '离线导航关系', category: 'local' },
  { num: 11, title: '本地: 树构建', description: '构建和遍历层次树', category: 'local' },
  { num: 12, title: 'SDK: 序列化', description: 'JSON、XML、CSV 导入导出', category: 'sdk' },
]

const categoryColors: Record<string, string> = {
  basic: '#1677ff',
  api: '#fa8c16',
  local: '#52c41a',
  sdk: '#13c2c2',
}

const categoryLabels: Record<string, string> = {
  basic: '基础',
  api: 'API',
  local: '离线',
  sdk: 'SDK',
}

const SkillsSection: React.FC = () => {
  return (
    <section
      id="skills"
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
            12 个 Skills 技能
          </Title>
          <Paragraph
            style={{
              color: 'rgba(255,255,255,0.4)',
              fontSize: 15,
              maxWidth: 520,
              margin: '0 auto',
            }}
          >
            渐进式技能文档，面向 AI 代理和开发者 — 从简到深
          </Paragraph>
        </div>

        {/* 分类图例 */}
        <div style={{ display: 'flex', justifyContent: 'center', gap: 20, marginBottom: 32 }}>
          {Object.entries(categoryLabels).map(([key, label]) => (
            <div key={key} style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <div
                style={{
                  width: 6,
                  height: 6,
                  borderRadius: 1,
                  background: categoryColors[key],
                }}
              />
              <Text style={{ color: 'rgba(255,255,255,0.4)', fontSize: 12 }}>
                {label}
              </Text>
            </div>
          ))}
        </div>

        <Row gutter={[16, 16]}>
          {skills.map((skill) => (
            <Col xs={24} sm={12} lg={8} key={skill.num}>
              <div
                style={{
                  padding: 16,
                  background: '#111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  borderRadius: 3,
                  height: '100%',
                  display: 'flex',
                  alignItems: 'flex-start',
                  gap: 10,
                }}
              >
                <div
                  style={{
                    width: 28,
                    height: 28,
                    borderRadius: 3,
                    background: `${categoryColors[skill.category]}12`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 12,
                    color: categoryColors[skill.category],
                    flexShrink: 0,
                    fontWeight: 700,
                  }}
                >
                  {skill.num}
                </div>
                <div style={{ flex: 1 }}>
                  <Text strong style={{ color: '#fff', fontSize: 13, display: 'block', marginBottom: 3 }}>
                    {skill.title}
                  </Text>
                  <Text style={{ color: 'rgba(255,255,255,0.4)', fontSize: 12, display: 'block', marginBottom: 6 }}>
                    {skill.description}
                  </Text>
                  <span
                    style={{
                      fontSize: 10,
                      color: categoryColors[skill.category],
                      background: `${categoryColors[skill.category]}10`,
                      padding: '1px 6px',
                      borderRadius: 2,
                    }}
                  >
                    {categoryLabels[skill.category]}
                  </span>
                </div>
              </div>
            </Col>
          ))}
        </Row>

        {/* CLI 命令树 */}
        <div style={{ maxWidth: 860, margin: '40px auto 0' }}>
          <Collapse
            ghost
            items={[
              {
                key: 'cli-tree',
                label: (
                  <Text style={{ color: 'rgba(255,255,255,0.5)', fontSize: 13 }}>
                    <CheckCircleOutlined style={{ marginRight: 6 }} />
                    查看 CLI 命令树
                  </Text>
                ),
                children: (
                  <img
                    src="/cwe-skills/cli-command-tree.png"
                    alt="CLI Command Tree"
                    style={{
                      width: '100%',
                      borderRadius: 3,
                      border: '1px solid rgba(255,255,255,0.06)',
                    }}
                  />
                ),
              },
            ]}
          />
        </div>
      </div>
    </section>
  )
}

export default SkillsSection
