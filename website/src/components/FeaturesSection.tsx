import { Typography, Row, Col } from 'antd'
import {
  DatabaseOutlined,
  AppstoreOutlined,
  IdcardOutlined,
  SafetyCertificateOutlined,
  CloudServerOutlined,
  FileSearchOutlined,
  SearchOutlined,
  FilterOutlined,
  ApartmentOutlined,
  PartitionOutlined,
  ExportOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons'
import { asset } from '../constants'

const { Title, Paragraph, Text } = Typography

interface FeatureItem {
  icon: React.ReactNode
  title: string
  description: string
  color: string
}

const features: FeatureItem[] = [
  {
    icon: <DatabaseOutlined />,
    title: '完整 CWE 数据模型',
    description: '弱点、类别、视图、复合元素',
    color: '#1677ff',
  },
  {
    icon: <AppstoreOutlined />,
    title: '类型化枚举',
    description: '抽象层级、状态、关系类型',
    color: '#13c2c2',
  },
  {
    icon: <IdcardOutlined />,
    title: 'CWE ID 工具',
    description: '解析、格式化、验证、提取、比较',
    color: '#52c41a',
  },
  {
    icon: <SafetyCertificateOutlined />,
    title: '知名列表',
    description: 'CWE Top 25、OWASP Top 10、SANS Top 25',
    color: '#fa8c16',
  },
  {
    icon: <CloudServerOutlined />,
    title: 'MITRE REST API 客户端',
    description: '速率限制、自动重试、结构化错误',
    color: '#1677ff',
  },
  {
    icon: <FileSearchOutlined />,
    title: 'XML 目录解析器',
    description: '离线解析 MITRE XML，无需网络',
    color: '#52c41a',
  },
  {
    icon: <SearchOutlined />,
    title: '搜索与过滤',
    description: '关键字、抽象层级、状态、排序',
    color: '#fa8c16',
  },
  {
    icon: <FilterOutlined />,
    title: '内存注册表',
    description: '存储、索引、查询及关系索引',
    color: '#13c2c2',
  },
  {
    icon: <ApartmentOutlined />,
    title: '关系导航',
    description: '父/子/祖先/后代/最短路径',
    color: '#1677ff',
  },
  {
    icon: <PartitionOutlined />,
    title: '树构建',
    description: '构建层次树、遍历、查找路径',
    color: '#52c41a',
  },
  {
    icon: <ExportOutlined />,
    title: '序列化',
    description: 'JSON、XML、CSV 导入导出',
    color: '#13c2c2',
  },
  {
    icon: <ThunderboltOutlined />,
    title: '零依赖',
    description: '核心 SDK 仅使用 Go 标准库',
    color: '#fa8c16',
  },
]

const FeaturesSection: React.FC = () => {
  return (
    <section
      id="features"
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
            功能特性
          </Title>
          <Paragraph
            style={{
              color: 'rgba(255,255,255,0.4)',
              fontSize: 15,
              maxWidth: 520,
              margin: '0 auto',
            }}
          >
            从 ID 解析到关系导航，覆盖 CWE 数据操作的完整链路
          </Paragraph>
        </div>

        <Row gutter={[16, 16]}>
          {features.map((feature, index) => (
            <Col xs={24} sm={12} lg={8} key={index}>
              <div
                style={{
                  padding: 18,
                  background: '#111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  borderRadius: 3,
                  height: '100%',
                  display: 'flex',
                  alignItems: 'flex-start',
                  gap: 12,
                }}
              >
                <div
                  style={{
                    width: 32,
                    height: 32,
                    borderRadius: 3,
                    background: `${feature.color}12`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 15,
                    color: feature.color,
                    flexShrink: 0,
                  }}
                >
                  {feature.icon}
                </div>
                <div>
                  <Text
                    strong
                    style={{
                      display: 'block',
                      color: '#fff',
                      fontSize: 14,
                      marginBottom: 3,
                    }}
                  >
                    {feature.title}
                  </Text>
                  <Text
                    style={{
                      color: 'rgba(255,255,255,0.4)',
                      fontSize: 12,
                      lineHeight: 1.5,
                    }}
                  >
                    {feature.description}
                  </Text>
                </div>
              </div>
            </Col>
          ))}
        </Row>

        {/* SDK API 地图 */}
        <div style={{ maxWidth: 860, margin: '40px auto 0' }}>
          <img
            src={asset('sdk-api-map.png')}
            alt="SDK API Map"
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

export default FeaturesSection
