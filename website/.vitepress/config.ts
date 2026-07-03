import { defineConfig } from 'vitepress'

// CWE Skills 官方文档站 VitePress 配置
// 仓库地址
const repoURL = 'https://github.com/scagogogo/cwe-skills'
const docsBase = '/cwe-skills/' // GitHub Pages 子路径，仓库名

export default defineConfig({
  lang: 'zh-CN',
  title: 'CWE Skills',
  titleTemplate: 'CWE Skills · AI原生CWE集成',
  description: 'AI原生的 CWE（通用缺陷枚举）集成层 — Skills、Go SDK、CLI 与 MCP 四种接入方式',
  base: docsBase,
  cleanUrls: true,
  lastUpdated: true,

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: `${docsBase}favicon.svg` }],
    ['meta', { name: 'theme-color', content: '#3c6c8f' }],
    ['meta', { property: 'og:title', content: 'CWE Skills · AI原生CWE集成' }],
    ['meta', { property: 'og:description', content: '统一 MITRE API、XML 目录与权威列表的 CWE 集成层' }],
  ],

  themeConfig: {
    logo: '/favicon.svg',

    // 顶部导航
    nav: [
      { text: '指南 📖', link: '/guide/what-is-cwe-skills' },
      { text: 'SDK API 🔧', link: '/sdk/overview' },
      { text: 'CLI 命令 💻', link: '/cli/overview' },
      { text: '枚举参考 📚', link: '/enums/abstraction' },
      {
        text: '更多 ⋯',
        items: [
          { text: '🦾 Skills（AI代理）', link: '/skills/' },
          { text: '🏆 知名列表', link: '/wellknown/cwe-top-25' },
          { text: '🗺️ 示例与教程', link: '/examples/' },
          { text: '📝 更新日志', link: '/changelog' },
        ],
      },
    ],

    // 社交链接
    socialLinks: [
      { icon: 'github', link: repoURL },
    ],

    // 搜索
    search: {
      provider: 'local',
      options: {
        translations: {
          button: { buttonText: '搜索文档', buttonAriaLabel: '搜索文档' },
          modal: {
            displayDetails: '显示详情',
            resetButtonTitle: '清除查询',
            backButtonTitle: '返回',
            noResultsText: '没有结果',
            footer: {
              selectText: '选择',
              navigateText: '切换',
              closeText: '关闭',
            },
          },
        },
      },
    },

    // 大纲（右侧目录）
    outline: {
      level: [2, 3],
      label: '本页目录',
    },

    docFooter: {
      prev: '上一页',
      next: '下一页',
    },

    lastUpdatedText: '最后更新于',
    returnToTopLabel: '回到顶部',
    sidebarMenuLabel: '菜单',

    // 页脚
    footer: {
      message: '基于 MIT 许可证发布',
      copyright: '© 2024-present scagogogo · CWE Skills',
    },

    // 侧边栏 —— 覆盖整个文档站的导航结构
    sidebar: {
      // ===== 指南 =====
      '/guide/': [
        {
          text: '开始 🚀',
          collapsed: false,
          items: [
            { text: 'CWE Skills 是什么？', link: '/guide/what-is-cwe-skills' },
            { text: '为什么需要 CWE Skills', link: '/guide/why' },
            { text: '解决了什么问题', link: '/guide/problem-solved' },
            { text: '工作原理 🧠', link: '/guide/how-it-works' },
            { text: '快速开始', link: '/guide/quick-start' },
            { text: '安装 ⬇️', link: '/guide/installation' },
          ],
        },
        {
          text: '核心概念 💡',
          collapsed: false,
          items: [
            { text: 'CWE 是什么', link: '/guide/concept-cwe' },
            { text: '抽象层级 (Pillar/Class/Base/Variant)', link: '/guide/concept-abstraction' },
            { text: '结构类型 (Simple/Chain/Composite)', link: '/guide/concept-structure' },
            { text: '关系类型 (ChildOf/Requires…)', link: '/guide/concept-relationship' },
            { text: '后果范围与影响', link: '/guide/concept-consequence' },
            { text: '视图 (View)', link: '/guide/concept-view' },
            { text: '类别 (Category)', link: '/guide/concept-category' },
            { text: '复合元素 (CompoundElement)', link: '/guide/concept-compound' },
          ],
        },
        {
          text: '四种接入方式 🔌',
          collapsed: false,
          items: [
            { text: '总览', link: '/guide/integrations' },
            { text: 'Skills — AI 代理', link: '/guide/integration-skills' },
            { text: 'Go SDK', link: '/guide/integration-sdk' },
            { text: 'CLI 命令行', link: '/guide/integration-cli' },
            { text: 'MCP（规划中）', link: '/guide/integration-mcp' },
          ],
        },
        {
          text: '进阶 ⚙️',
          collapsed: true,
          items: [
            { text: '在线 vs 离线模式', link: '/guide/online-offline' },
            { text: '速率限制与重试', link: '/guide/rate-limit-retry' },
            { text: '错误处理', link: '/guide/error-handling' },
            { text: '输出格式 (text/JSON)', link: '/guide/output-format' },
            { text: '性能与零依赖', link: '/guide/performance' },
            { text: '常见问题 FAQ', link: '/guide/faq' },
          ],
        },
      ],

      // ===== SDK API 参考 =====
      '/sdk/': [
        {
          text: 'SDK 总览 🔧',
          collapsed: false,
          items: [
            { text: 'SDK 概览', link: '/sdk/overview' },
            { text: '包与版本 cweskills', link: '/sdk/package' },
            { text: '模块地图 🗺️', link: '/sdk/module-map' },
            { text: '安装与导入', link: '/sdk/install' },
          ],
        },
        {
          text: 'CWE ID 工具 🆔',
          collapsed: false,
          items: [
            { text: 'cwe_utils 概览', link: '/sdk/cwe-utils' },
            { text: 'ParseCWEID 解析', link: '/sdk/parse-cwe-id' },
            { text: 'FormatCWEID 格式化', link: '/sdk/format-cwe-id' },
            { text: 'IsCWEID 判断', link: '/sdk/is-cwe-id' },
            { text: 'ValidateCWEID 验证', link: '/sdk/validate-cwe-id' },
            { text: 'ExtractCWEIDs 提取', link: '/sdk/extract-cwe-ids' },
            { text: 'ExtractFirstCWEID', link: '/sdk/extract-first-cwe-id' },
            { text: 'CompareCWEIDs 比较', link: '/sdk/compare-cwe-ids' },
            { text: 'FormatCWEIDFromInt', link: '/sdk/format-cwe-id-from-int' },
          ],
        },
        {
          text: '数据模型 🧱',
          collapsed: false,
          items: [
            { text: 'model.go 概览', link: '/sdk/model' },
            { text: 'CWE 弱点', link: '/sdk/cwe-struct' },
            { text: 'CWE 类型判断方法', link: '/sdk/cwe-type-methods' },
            { text: 'CWE 关系获取方法', link: '/sdk/cwe-relationship-methods' },
            { text: 'Mitigation 缓解措施', link: '/sdk/mitigation' },
            { text: 'Consequence 后果', link: '/sdk/consequence' },
            { text: 'DemonstrativeExample', link: '/sdk/demonstrative-example' },
            { text: 'ObservedExample', link: '/sdk/observed-example' },
            { text: 'Reference 参考文献', link: '/sdk/reference' },
            { text: 'ApplicablePlatforms', link: '/sdk/applicable-platforms' },
            { text: 'Introduction 引入方式', link: '/sdk/introduction' },
            { text: 'AlternateTerm 备用术语', link: '/sdk/alternate-term' },
            { text: 'ContentHistory 内容历史', link: '/sdk/content-history' },
            { text: 'Category 类别', link: '/sdk/category' },
            { text: 'View 视图', link: '/sdk/view' },
            { text: 'CompoundElement 复合元素', link: '/sdk/compound-element' },
          ],
        },
        {
          text: '枚举类型 📚',
          collapsed: false,
          items: [
            { text: 'enums.go 概览', link: '/sdk/enums' },
            { text: 'Abstraction 抽象层级', link: '/sdk/enum-abstraction' },
            { text: 'Structure 结构类型', link: '/sdk/enum-structure' },
            { text: 'Status 状态', link: '/sdk/enum-status' },
            { text: 'LikelihoodOfExploit', link: '/sdk/enum-likelihood' },
            { text: 'RelationshipNature', link: '/sdk/enum-relationship-nature' },
            { text: 'ConsequenceScope', link: '/sdk/enum-consequence-scope' },
            { text: 'ConsequenceImpact', link: '/sdk/enum-consequence-impact' },
            { text: 'ViewType 视图类型', link: '/sdk/enum-view-type' },
            { text: 'PlatformType 平台类型', link: '/sdk/enum-platform-type' },
          ],
        },
        {
          text: '知名列表 🏆',
          collapsed: false,
          items: [
            { text: 'wellknown_ids 概览', link: '/sdk/wellknown-ids' },
            { text: 'CWETop25 列表', link: '/sdk/cwe-top-25' },
            { text: 'OWASPTop10 映射', link: '/sdk/owasp-top-10' },
            { text: 'SANSTop25 列表', link: '/sdk/sans-top-25' },
            { text: 'IsInTop25', link: '/sdk/is-in-top-25' },
            { text: 'IsInOWASPTop10', link: '/sdk/is-in-owasp-top-10' },
            { text: 'IsInSANSTop25', link: '/sdk/is-in-sans-top-25' },
            { text: 'GetOWASPCategory', link: '/sdk/get-owasp-category' },
            { text: 'GetOWASPCategories', link: '/sdk/get-owasp-categories' },
            { text: 'IsInWellKnownView', link: '/sdk/is-in-wellknown-view' },
          ],
        },
        {
          text: '注册表与索引 🗃️',
          collapsed: false,
          items: [
            { text: 'registry.go 概览', link: '/sdk/registry' },
            { text: 'Registry 注册/查询', link: '/sdk/registry-operations' },
            { text: 'BuildIndexes 构建索引', link: '/sdk/build-indexes' },
            { text: '关系索引查询', link: '/sdk/relationship-indexes' },
            { text: 'ExportJSON / ImportJSON', link: '/sdk/registry-json' },
          ],
        },
        {
          text: '关系导航 🧭',
          collapsed: false,
          items: [
            { text: 'navigator.go 概览', link: '/sdk/navigator' },
            { text: 'Parents / Children', link: '/sdk/nav-parents-children' },
            { text: 'Ancestors / Descendants', link: '/sdk/nav-ancestors-descendants' },
            { text: 'Siblings / Peers', link: '/sdk/nav-siblings-peers' },
            { text: 'CanPrecede / CanFollow', link: '/sdk/nav-precede-follow' },
            { text: 'Requires / RequiredBy', link: '/sdk/nav-requires' },
            { text: 'CanAlsoBe', link: '/sdk/nav-can-also-be' },
            { text: 'ChainMembers / CompositeMembers', link: '/sdk/nav-chain-composite' },
            { text: 'ShortestPath 最短路径', link: '/sdk/nav-shortest-path' },
            { text: 'IsAncestorOf / IsRelated', link: '/sdk/nav-ancestor-related' },
            { text: 'RelationshipDepth', link: '/sdk/nav-relationship-depth' },
          ],
        },
        {
          text: '树构建 🌳',
          collapsed: false,
          items: [
            { text: 'tree.go 概览', link: '/sdk/tree' },
            { text: 'TreeNode 节点', link: '/sdk/tree-node' },
            { text: 'BuildTree 构建树', link: '/sdk/build-tree' },
            { text: 'BuildForest 构建森林', link: '/sdk/build-forest' },
            { text: 'BuildViewTree 视图树', link: '/sdk/build-view-tree' },
            { text: 'Walk / WalkBFS 遍历', link: '/sdk/tree-walk' },
            { text: 'Path 路径查找', link: '/sdk/tree-path' },
            { text: 'LeafNodes 叶子节点', link: '/sdk/tree-leaf-nodes' },
            { text: 'MaxDepth / Count', link: '/sdk/tree-depth-count' },
          ],
        },
        {
          text: '搜索与过滤 🔍',
          collapsed: false,
          items: [
            { text: 'search.go 概览', link: '/sdk/search' },
            { text: 'FindByID', link: '/sdk/find-by-id' },
            { text: 'FindByKeyword', link: '/sdk/find-by-keyword' },
            { text: 'FindByAbstraction', link: '/sdk/find-by-abstraction' },
            { text: 'FindByStatus', link: '/sdk/find-by-status' },
            { text: 'FindByLikelihood', link: '/sdk/find-by-likelihood' },
            { text: 'FindByConsequenceScope', link: '/sdk/find-by-consequence-scope' },
            { text: 'FindByStructure', link: '/sdk/find-by-structure' },
            { text: 'FindTopLevel / FindBaseWeaknesses', link: '/sdk/find-top-level' },
            { text: 'FindChains / FindComposites', link: '/sdk/find-chains-composites' },
            { text: 'filter.go 过滤', link: '/sdk/filter' },
            { text: 'FilterOption 选项', link: '/sdk/filter-option' },
            { text: 'SortByID / SortByName', link: '/sdk/sort' },
            { text: 'GroupByAbstraction 等', link: '/sdk/group-by' },
            { text: 'Deduplicate 去重', link: '/sdk/deduplicate' },
          ],
        },
        {
          text: '统计 📊',
          collapsed: false,
          items: [
            { text: 'stats.go 概览', link: '/sdk/stats' },
            { text: 'ComputeStatistics', link: '/sdk/compute-statistics' },
            { text: 'CountByAbstraction', link: '/sdk/count-by-abstraction' },
            { text: 'CountByStatus', link: '/sdk/count-by-status' },
            { text: 'CountByLikelihood', link: '/sdk/count-by-likelihood' },
            { text: 'CountByScope', link: '/sdk/count-by-scope' },
          ],
        },
        {
          text: '序列化 📦',
          collapsed: false,
          items: [
            { text: 'serializer.go 概览', link: '/sdk/serializer' },
            { text: 'MarshalJSON / UnmarshalJSON', link: '/sdk/marshal-json' },
            { text: 'MarshalJSONList', link: '/sdk/marshal-json-list' },
            { text: 'MarshalXML / UnmarshalXML', link: '/sdk/marshal-xml' },
            { text: 'MarshalCSV / UnmarshalCSV', link: '/sdk/marshal-csv' },
            { text: 'ExportCSV (Registry)', link: '/sdk/export-csv' },
          ],
        },
        {
          text: 'MITRE API 客户端 🌐',
          collapsed: false,
          items: [
            { text: 'api_client.go 概览', link: '/sdk/api-client' },
            { text: 'NewAPIClient 与选项', link: '/sdk/new-api-client' },
            { text: 'GetWeakness', link: '/sdk/api-get-weakness' },
            { text: 'GetCategory', link: '/sdk/api-get-category' },
            { text: 'GetView', link: '/sdk/api-get-view' },
            { text: 'GetCWEs 批量', link: '/sdk/api-get-cwes' },
            { text: 'GetVersion 版本', link: '/sdk/api-get-version' },
            { text: 'GetParents / GetChildren', link: '/sdk/api-parents-children' },
            { text: 'GetAncestors / GetDescendants', link: '/sdk/api-ancestors-descendants' },
            { text: 'APIResponse 响应类型', link: '/sdk/api-response' },
          ],
        },
        {
          text: 'HTTP 客户端与限流 🚦',
          collapsed: false,
          items: [
            { text: 'http_client.go 概览', link: '/sdk/http-client' },
            { text: 'HTTPClientOption', link: '/sdk/http-client-option' },
            { text: 'Get / Post 请求方法', link: '/sdk/http-methods' },
            { text: '重试机制', link: '/sdk/http-retry' },
            { text: 'http_rate_limiter 限流', link: '/sdk/rate-limiter' },
            { text: 'RateLimiter API', link: '/sdk/rate-limiter-api' },
          ],
        },
        {
          text: 'XML 解析与数据获取 📥',
          collapsed: false,
          items: [
            { text: 'xml_parser.go 概览', link: '/sdk/xml-parser' },
            { text: 'NewXMLParser', link: '/sdk/new-xml-parser' },
            { text: 'ParseFile / Parse', link: '/sdk/xml-parse' },
            { text: 'data_fetcher 概览', link: '/sdk/data-fetcher' },
            { text: 'BasicFetcher', link: '/sdk/basic-fetcher' },
            { text: 'MultipleFetcher', link: '/sdk/multiple-fetcher' },
            { text: 'TreeFetcher', link: '/sdk/tree-fetcher' },
          ],
        },
        {
          text: '错误处理 ⚠️',
          collapsed: false,
          items: [
            { text: 'errors.go 概览', link: '/sdk/errors' },
            { text: 'CWEError 基础错误', link: '/sdk/cwe-error' },
            { text: 'InvalidCWEIDError', link: '/sdk/invalid-cwe-id-error' },
            { text: 'CWENotFoundError', link: '/sdk/cwe-not-found-error' },
            { text: 'APIError', link: '/sdk/api-error' },
            { text: 'RateLimitError', link: '/sdk/rate-limit-error' },
            { text: 'ValidationError', link: '/sdk/validation-error' },
            { text: 'ParseError', link: '/sdk/parse-error' },
            { text: 'RelationshipError', link: '/sdk/relationship-error' },
          ],
        },
      ],

      // ===== CLI 命令 =====
      '/cli/': [
        { text: 'CLI 总览 💻', link: '/cli/overview' },
        { text: '安装 CLI', link: '/cli/install' },
        { text: '全局参数 -o/--output', link: '/cli/global-flags' },
        {
          text: 'ID 工具 🆔',
          collapsed: false,
          items: [
            { text: 'parse 解析', link: '/cli/parse' },
            { text: 'validate 验证', link: '/cli/validate' },
            { text: 'format 格式化', link: '/cli/format' },
            { text: 'extract 提取', link: '/cli/extract' },
            { text: 'compare 比较', link: '/cli/compare' },
            { text: 'compare-int', link: '/cli/compare-int' },
          ],
        },
        {
          text: '枚举 📚',
          collapsed: false,
          items: [
            { text: 'enum 总览', link: '/cli/enum' },
            { text: 'enum abstraction', link: '/cli/enum-abstraction' },
            { text: 'enum status', link: '/cli/enum-status' },
            { text: 'enum relationship', link: '/cli/enum-relationship' },
            { text: 'enum structure', link: '/cli/enum-structure' },
            { text: 'enum likelihood', link: '/cli/enum-likelihood' },
            { text: 'enum consequence-scope', link: '/cli/enum-consequence-scope' },
            { text: 'enum consequence-impact', link: '/cli/enum-consequence-impact' },
            { text: 'enum view-type', link: '/cli/enum-view-type' },
            { text: 'enum platform', link: '/cli/enum-platform' },
          ],
        },
        {
          text: '知名列表 🏆',
          collapsed: false,
          items: [
            { text: 'wellknown 总览', link: '/cli/wellknown' },
            { text: 'wellknown top25', link: '/cli/wellknown-top25' },
            { text: 'wellknown owasp', link: '/cli/wellknown-owasp' },
            { text: 'wellknown sans', link: '/cli/wellknown-sans' },
            { text: 'wellknown check', link: '/cli/wellknown-check' },
          ],
        },
        {
          text: 'API 命令 🌐',
          collapsed: false,
          items: [
            { text: 'show 弱点详情', link: '/cli/show' },
            { text: 'show category', link: '/cli/show-category' },
            { text: 'show view', link: '/cli/show-view' },
            { text: 'relations 总览', link: '/cli/relations' },
            { text: 'relations parents', link: '/cli/relations-parents' },
            { text: 'relations children', link: '/cli/relations-children' },
            { text: 'relations ancestors', link: '/cli/relations-ancestors' },
            { text: 'relations descendants', link: '/cli/relations-descendants' },
            { text: 'api-version', link: '/cli/api-version' },
          ],
        },
        {
          text: '本地搜索过滤 🔍',
          collapsed: false,
          items: [
            { text: 'search 搜索', link: '/cli/search' },
            { text: 'filter 过滤', link: '/cli/filter' },
            { text: 'stats 统计', link: '/cli/stats' },
          ],
        },
        {
          text: '注册表 🗃️',
          collapsed: false,
          items: [
            { text: 'registry 总览', link: '/cli/registry' },
            { text: 'registry load', link: '/cli/registry-load' },
            { text: 'registry get', link: '/cli/registry-get' },
            { text: 'registry contains', link: '/cli/registry-contains' },
            { text: 'registry list-views', link: '/cli/registry-list-views' },
            { text: 'registry list-categories', link: '/cli/registry-list-categories' },
            { text: 'registry parents/children', link: '/cli/registry-relations' },
            { text: 'registry ancestors/descendants', link: '/cli/registry-anc-desc' },
            { text: 'registry peers', link: '/cli/registry-peers' },
            { text: 'registry view-members', link: '/cli/registry-view-members' },
            { text: 'registry category-members', link: '/cli/registry-category-members' },
            { text: 'registry member-of', link: '/cli/registry-member-of' },
            { text: 'registry export', link: '/cli/registry-export' },
          ],
        },
        {
          text: '导航 🧭',
          collapsed: false,
          items: [
            { text: 'nav 总览', link: '/cli/nav' },
            { text: 'nav parents', link: '/cli/nav-parents' },
            { text: 'nav children', link: '/cli/nav-children' },
            { text: 'nav ancestors', link: '/cli/nav-ancestors' },
            { text: 'nav descendants', link: '/cli/nav-descendants' },
            { text: 'nav siblings', link: '/cli/nav-siblings' },
            { text: 'nav peers', link: '/cli/nav-peers' },
            { text: 'nav precede', link: '/cli/nav-precede' },
            { text: 'nav follow', link: '/cli/nav-follow' },
            { text: 'nav requires', link: '/cli/nav-requires' },
            { text: 'nav required-by', link: '/cli/nav-required-by' },
            { text: 'nav can-also-be', link: '/cli/nav-can-also-be' },
            { text: 'nav chain-members', link: '/cli/nav-chain-members' },
            { text: 'nav composite-members', link: '/cli/nav-composite-members' },
            { text: 'nav shortest-path', link: '/cli/nav-shortest-path' },
            { text: 'nav is-ancestor', link: '/cli/nav-is-ancestor' },
            { text: 'nav is-related', link: '/cli/nav-is-related' },
            { text: 'nav depth', link: '/cli/nav-depth' },
          ],
        },
        {
          text: '树 🌳',
          collapsed: false,
          items: [
            { text: 'tree 总览', link: '/cli/tree' },
            { text: 'tree build', link: '/cli/tree-build' },
            { text: 'tree forest', link: '/cli/tree-forest' },
            { text: 'tree view', link: '/cli/tree-view' },
            { text: 'tree path', link: '/cli/tree-path' },
            { text: 'tree leaves', link: '/cli/tree-leaves' },
          ],
        },
        {
          text: '其他 🛠️',
          collapsed: false,
          items: [
            { text: 'version 版本', link: '/cli/version' },
          ],
        },
      ],

      // ===== 枚举参考 =====
      '/enums/': [
        {
          text: '枚举参考 📚',
          collapsed: false,
          items: [
            { text: '总览', link: '/enums/overview' },
            { text: 'Abstraction 抽象层级', link: '/enums/abstraction' },
            { text: 'Structure 结构类型', link: '/enums/structure' },
            { text: 'Status 状态', link: '/enums/status' },
            { text: 'LikelihoodOfExploit 利用可能性', link: '/enums/likelihood' },
            { text: 'RelationshipNature 关系类型', link: '/enums/relationship-nature' },
            { text: 'ConsequenceScope 后果范围', link: '/enums/consequence-scope' },
            { text: 'ConsequenceImpact 后果影响', link: '/enums/consequence-impact' },
            { text: 'ViewType 视图类型', link: '/enums/view-type' },
            { text: 'PlatformType 平台类型', link: '/enums/platform-type' },
          ],
        },
      ],

      // ===== 知名列表 =====
      '/wellknown/': [
        {
          text: '知名列表 🏆',
          collapsed: false,
          items: [
            { text: '总览', link: '/wellknown/overview' },
            { text: 'CWE Top 25', link: '/wellknown/cwe-top-25' },
            { text: 'OWASP Top 10 (2021)', link: '/wellknown/owasp-top-10' },
            { text: 'SANS Top 25', link: '/wellknown/sans-top-25' },
            { text: '知名视图', link: '/wellknown/well-known-views' },
          ],
        },
      ],
    },

    editLink: {
      pattern: `${repoURL}/edit/main/website/:path`,
      text: '在 GitHub 上编辑此页',
    },

    externalLinkIcon: true,
  },
})
