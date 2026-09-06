# Archify 使用指南（pi agent 适配版）

> Archify 是一个 Agent Skill：把代码库或口述的系统描述变成**可交互、可验证**的系统图（自包含 HTML + SVG），支持 architecture / workflow / sequence / dataflow / lifecycle 五种图。
> 官方仓库：https://github.com/tt-a1i/archify · 本文档基于本机安装的 **v2.17.0-dev.1** 与官方 SKILL.md / references / schemas 整理，所有示例均在本机跑通。

**本机安装状态**

| 项 | 值 |
|---|---|
| 安装位置 | `~/.pi/agent/skills/archify/`（全局，零 runtime 依赖） |
| 运行环境 | Node v24.18.0（要求 ≥18） |
| 自检 | `node bin/archify.mjs doctor` 全绿 |
| Skill 目录 | SKILL.md · bin/ · schemas/ · references/ · examples/ · recipes/ · renderers/ |

---

## 1. 五种图怎么选

| 类型 | 适用场景 | 核心结构 | 本机 demo |
|---|---|---|---|
| `architecture` | 组件/服务/存储、信任边界、云拓扑 | components + boundaries + connections | `demos/blog.architecture.html` |
| `workflow` | 流程、审批门禁、CI/CD、工具调用 | lanes + nodes + edges | `cicd-release.html` |
| `sequence` | API 调用链、请求生命周期、异步时序 | participants + messages + activations | `demos/login.sequence.html` |
| `dataflow` | 管道、ETL、数据血缘、PII 边界 | stages + nodes + flows | `demos/order-pipeline.dataflow.html` |
| `lifecycle` | 状态机、重试、等待、终态 | lanes + states + transitions | `demos/order.lifecycle.html` |

拿不准时让 CLI 帮你判断：

```bash
node ~/.pi/agent/skills/archify/bin/archify.mjs guide "API 请求 Redis 缓存未命中" --json
node ~/.pi/agent/skills/archify/bin/archify.mjs guide "Kafka 重放与死信队列" --lang zh
```

---

## 2. 在 pi 里怎么用

### 2.1 三种触发方式

1. **自然语言自动触发**（新会话生效）：提到架构图/流程图/时序图/状态机/Mermaid 等关键词，pi 自动加载 skill。
2. **强制调用**：`/skill:archify 画一个 CI/CD 流程图`。
3. **我手动走流程**：任何会话里说"用 archify 画…"，agent 读取 SKILL.md 后按标准流程执行。

### 2.2 Prompt 模板库（官方推荐句式）

**从描述开始（不需要仓库）：**

```text
用 archify 画：浏览器 -> API -> Redis 缓存 -> PostgreSQL 回退。
```

**从真实代码生成（打开仓库后）：**

```text
分析这个仓库，然后用 archify 画一张高层运行时架构图。
展示 8–12 个核心组件、一条主路径、外部依赖和信任边界。
把支撑细节放进卡片，不要为了细节加边。
```

**按类型：**

```text
画 CI/CD 流程：lint → test → 审批 → 部署，含回滚分支。
画 API 时序：Web → API → Redis 未命中 → Postgres 回退 → 异步 trace。
画订单 ETL 数据流：标出 PII 边界和消费方。
画支付单状态机：含重试、等待、取消终态。
```

**对话式迭代（JSON 源保留，可精准修改）：**

```text
加一个 Redis / 把审批移到左边 / 高亮回滚路径 / 换成 signal-flow 风格 / 加 trace 动画
```

**Mermaid 转换**：直接贴 mermaid 代码（支持 `flowchart`/`graph` → workflow 或 architecture；`sequenceDiagram` → sequence；`stateDiagram` → lifecycle），说"用 archify 重画成可交互图"。

---

## 3. Agent 标准工作流（SKILL.md 快速创作路径）

1. **选题**：从需求判断五种类型之一。
2. **读格式**：读对应 `schemas/<type>.schema.json` + `schemas/common.schema.json` + 一个同类 example（只学字段形状，不抄内容）。
3. **先写工件**：直接写候选 JSON（≤12 个主节点、一条主路径、少标签），`meta.quality_profile: "showcase"`；**先不加** `via/channelX/channelY/labelAt` 等几何控制，等诊断要求再加。
4. **每轮验证**：
   ```bash
   node bin/archify.mjs validate <type> <candidate.json> --quality showcase --json
   ```
   showcase 通过标准：**9/9 项检查 + 0 错误 + 0 警告**。
5. **原子交付**：
   ```bash
   node bin/archify.mjs deliver <type> <candidate.json> <output.html> --quality showcase --json
   ```
   失败不覆盖旧文件；成功返回规范/工件的 SHA-256 与字节数。
6. **浏览器实测**（交付成功后）：
   ```bash
   node bin/archify.mjs visual-check <output.html> --json
   ```

**验收三声明（不可混为一谈）**：`deliver` 证明确定性工件检查 ✓ → `visual-check` 证明真实浏览器行为 ✓ → 感知层面的视觉复查需要人或图像模型亲自看。任何一项通过都不代表其他两项。

**修复规则**：按诊断的 `subject` + `supportedFixes` 修，一次只改一处几何控制；连续两轮无改善就停止并如实报告。

---

## 4. CLI 命令大全

```bash
cd ~/.pi/agent/skills/archify   # 以下命令在此目录执行
```

| 命令 | 用途 |
|---|---|
| `render <type> <in.json> [out.html] [--quality]` | 渲染，不做交付门禁 |
| `validate <type> <in.json> [--quality] [--layout-json]` | schema + 布局全量校验，输出机器可读诊断 |
| `deliver <type> <in.json> <out.html> [--json] [--open]` | **最终验收**：原子渲染+检查+提交，输出 SHA-256 回执 |
| `compare architecture <base.json> <head.json> <out.html>` | 两份架构快照对比 → Before/Delta/After |
| `preview <type> <in.json> <out.html>` | 桌面实时预览（监听 JSON，验证通过才刷新，Ctrl-C 退出） |
| `visual-check <out.html> [--json]` | Chrome 实测 4 视口 + 深/浅截图，输出 PNG 侧车 |
| `migrate workflow <old.json> <new.json> --to-schema 2` | workflow v1 → v2 迁移 |
| `inspect <type> <in.json>` | 查看 IR 摘要 |
| `check <out.html>` | 校验已交付 HTML 完整性 |
| `guide [场景] [--json] [--lang en\|zh]` | 场景→类型推荐 |
| `brands [名称/域名] [--json]` | 查内置品牌图标目录 |
| `brands capture <url> [--json]` | 抓取并钉住一个新品牌图标（返回 digest） |
| `examples` | 列出内置示例 |
| `demo [目录]` | 生成全套示例 |
| `doctor` | 环境自检 |

`--repo-root <path>` 仅 architecture 的 render/validate/deliver/preview/compare 接受，用于仓库证据校验。

---

## 5. JSON IR 通用规范

所有类型共享：`schema_version` + `diagram_type` + `meta`（必填 `title`）+ 类型结构数组。**任何一层都 `additionalProperties: false`**——多写的字段直接报错。

### 5.1 meta 公共字段

| 字段 | 值 | 说明 |
|---|---|---|
| `title` | string（必填） | 图标题 |
| `locale` | `"en"` \| `"zh-CN"` | 只翻译**渲染器 UI**（图例默认文案、快捷键面板、`<html lang>`），不翻译你写的内容；中文图必须写 `zh-CN` |
| `quality_profile` | `"showcase"` \| `"standard"` | showcase=9 项严格构图检查（默认用它） |
| `animation` | `"trace"` \| `"none"` | 默认静态；演示/汇报场景才开 trace |
| `visual_preset` | `classic`(默认) \| `signal-flow` \| `blueprint` \| `editorial` | 只影响观感不改语义；用户点名才写 |
| `views` | ≤5 个 | 引导章节：`{id,label,focus:[节点ID],note}`，生成故事线播放 |
| `legend` | `{mode:"auto"\|"all"\|"hidden", entries}` | 图例改写；auto 只列图中实际出现的类型 |
| `viewBox` | `[w,h]` | 画布尺寸；**能省则省**（见 §10 的坑） |
| `output` | string | 建议输出路径 |
| `subtitle` | string | 默认省略，别用来复述标题 |

### 5.2 共享枚举

- **节点类型 `componentType`**：`frontend` `backend` `database` `cloud` `security` `messagebus` `external`
- **变体 `variant`**：`default` `emphasis`（主线强调）`security`（安全红）`dashed`（次要/异步）；sequence 消息额外有 `return`
- **ID 规则**：`^[a-zA-Z][a-zA-Z0-9_-]*$`，集合内唯一
- **品牌 `brand`**：节点可带内置品牌 ID（`archify brands "Claude" --json` 查询）或 `brands capture` 返回的 `{url, sha256}`；不猜测、不虚构

### 5.3 cards（底部结论卡）

```json
"cards": [
  { "dot": "cyan",   "title": "发布门禁", "items": ["要点一", "要点二"] },
  { "dot": "rose",   "title": "异常处置", "items": ["…"] },
  { "dot": "emerald"|"amber"|"orange", "title": "…", "items": ["…"] }
]
```

---

## 6. 各类型字段速查（含本机已验证 demo）

### 6.1 architecture — `demos/blog.architecture.json`

```jsonc
{
  "schema_version": 1, "diagram_type": "architecture",
  "meta": { "title": "博客平台架构", "locale": "zh-CN", "quality_profile": "showcase" },
  "components": [
    // pos/size 显式定位；label+sublabel；tag 可选角标
    { "id": "api", "type": "backend", "label": "API 服务", "sublabel": "Node.js :3000",
      "pos": [470, 300], "size": [130, 60] }
  ],
  "boundaries": [
    // kind: region | security-group …；真实所有权/信任边界才画
    { "kind": "region", "label": "云区域: cn-hangzhou", "wraps": ["cdn","api","redis","db"] }
  ],
  "connections": [
    { "id": "users-cdn", "from": "users", "to": "cdn", "label": "HTTPS", "variant": "emphasis" },
    // 几何控制仅在诊断要求时加：
    { "id": "api-redis", "from": "api", "to": "redis", "label": "read-through",
      "fromSide": "top", "toSide": "bottom", "labelDy": -68 }
  ]
}
```

- 布局建议：一条左→右主脊 + 短垂直分支；6–12 个主组件；外部角色放边界外
- 进阶：`meta.engineering_profile: "deployment-ownership"` 失闭式部署审查（要求 owner/region/跨边界机制齐全，缺事实就别开）
- 仓库证据：`meta.repository`（公开 GitHub URL + commit SHA）+ 组件 `sources`（文件+行范围），渲染时必须 `--repo-root` 且 Git 验证通过；节点会标 `SRC n` 可点开验证过的源码

### 6.2 workflow — `cicd-release.workflow.json`（新图一律 schema_version: 2）

```jsonc
{
  "schema_version": 2, "diagram_type": "workflow",
  "meta": { "title": "CI/CD 发布流水线", "locale": "zh-CN", "quality_profile": "showcase" },
  "lanes":  [ { "id": "pipeline", "label": "CI 流水线" },
              { "id": "ops", "label": "生产运维", "variant": "exception" } ],
  "phases": [ { "id": "ci", "label": "持续集成", "fromCol": 0, "toCol": 2 } ],   // 可选
  "mainPath": ["commit", "lint", "test", "approval", "deploy", "prod"],          // 主路径
  "nodes": [
    // col 固定 0..5；lane 决定行；type 决定配色
    { "id": "approval", "lane": "release", "col": 3, "type": "security",
      "label": "人工审批", "sublabel": "发布负责人确认", "tag": "发布门禁", "width": 132 }
  ],
  "edges": [
    { "id": "approval-deploy", "from": "approval", "to": "deploy", "label": "批准", "variant": "emphasis" },
    { "id": "approval-reject", "from": "approval", "to": "commit",
      "label": "拒绝，修复后重提", "variant": "dashed", "role": "error" }
    // role: main|branch|async|return|error；route: auto|straight|drop|outside-right|return-left|bottom-channel|up-channel
  ]
}
```

- 语义不变式：主路径单调向右；语义标签**永不删除**（不是间距修复手段）；异常回路走主泳道外侧
- 可选 `semanticChecks`：`allowedRoots` / `allowedTerminals` / `requiredEdges` / `requiredPaths`，编译前校验拓扑事实
- 旧 v1（固定坐标）用 `migrate workflow old.json new.json --to-schema 2` 升级

### 6.3 sequence — `demos/login.sequence.json`

```jsonc
{
  "schema_version": 1, "diagram_type": "sequence",
  "meta": { "title": "用户登录时序", "locale": "zh-CN", "viewBox": [1000, 480], "quality_profile": "showcase" },
  "participants": [
    { "id": "auth", "type": "security", "label": "认证中心", "sublabel": "密码校验" }
  ],
  "segments": [ { "from": 150, "to": 415, "label": "认证与查询" } ],   // 可选：纵向分段
  "messages": [
    // y 是纵向时间轴（有可读范围限制）；variant 支持 return
    { "id": "post-login", "from": "app", "to": "api", "y": 185, "label": "POST /login", "variant": "emphasis" },
    { "id": "verify-password", "from": "api", "to": "auth", "y": 228, "label": "校验密码", "variant": "security" }
  ],
  "activations": [
    { "participant": "api", "from": 185, "to": 370, "type": "backend" }  // 可选：激活条
  ]
}
```

- `meta.column_fit`: 默认 `fixed`（108px 列距）；viewBox 很宽时用 `"spread"` 把宽度变成列距和标签空间
- sequence **不用**自动端口展开；消息的纵向顺序即语义顺序

### 6.4 dataflow — `demos/order-pipeline.dataflow.json`

```jsonc
{
  "schema_version": 1, "diagram_type": "dataflow",
  "meta": { "title": "订单数据管道", "locale": "zh-CN", "viewBox": [1080, 560], "quality_profile": "showcase" },
  "stages": [ { "label": "采集" }, { "label": "接入" }, { "label": "处理" }, { "label": "存储" }, { "label": "消费" } ],
  "nodes": [
    // stage=列（0..4），row=行（并行流分离）
    { "id": "kafka", "type": "messagebus", "label": "Kafka", "sublabel": "订单主题",
      "stage": 2, "row": 2, "tag": "有序" }
  ],
  "flows": [
    // classification 是数据分级标签（如 "PII 接触"），dataflow 特有
    { "id": "risk-check", "from": "gateway", "to": "risk",
      "label": "敏感数据校验", "classification": "PII 接触", "variant": "security" }
  ]
}
```

- 只给数据契约/分级/跨边界移动写标签；不写泛化动词

### 6.5 lifecycle — `demos/order.lifecycle.json`

```jsonc
{
  "schema_version": 1, "diagram_type": "lifecycle",
  "meta": { "title": "订单状态机", "locale": "zh-CN", "viewBox": [1050, 580], "quality_profile": "showcase" },
  "lanes": [ { "id": "main", "label": "订单推进" }, { "id": "terminal", "label": "终态出口" } ],
  "states": [
    // type: start|active|waiting|decision|success|failure|neutral|external
    // 主轨 col 0..4；事件/终态带 col 0..2，且终态列 N 对齐主列 N+2
    { "id": "created",  "type": "start",   "label": "已创建", "lane": "main", "col": 0, "step": "01" },
    { "id": "cancelled","type": "failure", "label": "已取消", "lane": "terminal", "col": 0,
      "yOffset": -60 }   // 终态纵向微调（诊断建议时才用）
  ],
  "transitions": [
    { "from": "created", "to": "paid" },
    { "from": "paid", "to": "cancelled", "variant": "security" },   // 可恢复状态要有真实回路
    { "from": "shipped", "to": "delivered", "variant": "emphasis" }
  ]
}
```

- "可恢复"必须落在拓扑上：failure 类型 + 真实 transition 回到活动状态；卡片/章节里写"重试"不算数

### 6.6 Architecture Delta（变更评审）

```bash
node bin/archify.mjs compare architecture \
  examples/checkout-platform.base.architecture.json \
  examples/checkout-platform.head.architecture.json \
  ~/output/checkout-delta.html --json
```

输出 Before/Delta/After 三视图，精确列出**新增/删除/变更/移动/改道**的事实（本机已生成 `demos/checkout-delta.html`）。只对比已验证快照，不推断风险或合并安全性。

---

## 7. Viewer 交互（打开 HTML 后）

| 按键 | 功能 |
|---|---|
| `?` | 图例指南 |
| `/` | 搜索节点（标签 + 稳定 ID） |
| 点节点 | Semantic Passport：上游/下游事实 + 可复制深链 |
| `R` | 路径探测（两节点间**已授权**的有向路由） |
| `L` | 语义透镜（对比两类角色的流量） |
| `M` | 语义雷达（视口镜像总览） |
| `P` / `[` `]` | 播放引导故事 / 切章节（来自 `meta.views`，≤5 章） |
| `F` | 演示舞台（Presentation Stage） |
| `S` / `T` / `E` | 切换视觉预设 / 深浅主题 / 导出菜单 |

**导出**：全图 PNG（复制/下载）、JPEG/WebP、双主题 SVG、trace 动画 WebM、1200×630 分享卡、路由分享卡（Route Probe 之后）、可达性分享卡（Reach 查询之后）。导出自动剥离一切临时查看状态。

**深链**：`#focus=<id>&reach=upstream|downstream`、`#route=<a>~<b>`、`#lens=<a>~<b>`、`#view=<id>`、`#relation=<id>`（需在边/消息/flows 上显式写 `id`）。

**真实边界**：交互只复用已授权的节点与关系，不发明拓扑、不声称运行时影响；Reach 只是"授权可达性"，不是"爆炸半径"。

---

## 8. 质量体系速查

**showcase 的 9 项工件检查**：single_svg · finite_svg · orthogonal_arrows · label_route_clearance · relationship_crossings · relationship_corridors · container_border_runs · route_rhythm · legend_clearance

**标签净距公式**（authoring-contract）：

```text
净距 > 标签掩码宽 + 8px
标签掩码宽 ≈ 6.5px × ASCII 单元数 + 13px   （CJK 字符按 2 个单元计）
```

**修复顺序**：meta/schema 错误 → 节点重叠/越界 → 边穿节点/端点方向 → 交叉/走廊/边界贴边/路由节奏 → 标签净距。每次修完重新 `validate`。

**诊断格式**（机器可读，稳定 code）：

```json
{ "code": "layout/constraint", "severity": "error",
  "message": "…keep y between 160 and 477",
  "subject": { …精确到集合和下标… },
  "evidence": { …测量数据… },
  "supportedFixes": [ "…建议动作…" ] }
```

**交付回执模板**（handoff receipt）：

```text
diagram_type: …
output: /absolute/path.html
specification_sha256 / artifact_sha256: …
validation: 9/9 showcase, 0 errors, 0 warnings
browser_evidence: passed | failed | skipped
visual_review: passed | skipped (image reader unavailable) | failed
correction_rounds: 0 | 1 | 2
```

---

## 9. 本机 demo 清单（全部 validate 9/9 + visual-check 通过）

| 文件 | 类型 | 说明 |
|---|---|---|
| `cicd-release.workflow.json` → `cicd-release.html` | workflow v2 | CI/CD：lint→test→审批→部署 + 拒绝回路 + 回滚分支（1 轮通过） |
| `demos/blog.architecture.json` → `.html` | architecture | 博客平台：读者→CDN→API→Redis/PostgreSQL + region 边界（1 轮通过） |
| `demos/login.sequence.json` → `.html` | sequence | 登录时序：POST /login→校验密码→查库→签发 JWT（1 轮修复后通过） |
| `demos/order-pipeline.dataflow.json` → `.html` | dataflow | 订单五段管道 + 风控/PII 边界（1 轮通过） |
| `demos/order.lifecycle.json` → `.html` | lifecycle | 订单状态机 + 取消/退款终态（2 轮修复后通过） |
| `demos/checkout-delta.html` | compare | 官方 checkout base/head 对比图（含 added/removed/changed 事实清单） |

重新渲染任意 demo：

```bash
cd ~/.pi/agent/skills/archify
node bin/archify.mjs deliver architecture /Volumes/data/project/archify/demos/blog.architecture.json \
  /Volumes/data/project/archify/demos/blog.architecture.html --quality showcase --json --open
```

---

## 10. 实战经验与坑（本机踩过，官方文档未明说的）

1. **首屏溢出**：sequence/lifecycle 用高 viewBox（如 820×600）时页面会纵向滚动，`visual-check` 判 fail（官方自带示例同样如此）。对策：**压扁内容纵横比到 ~2:1**——sequence 精简消息数 + `[1000,480]`；lifecycle `[1050,580]`。
2. **viewer 反向收窄**：viewBox 宽度过大（如 1600）时 viewer 会把阅读列收窄到 ~930px 保高度，导致投影字号 < 6px 下限（`composition/desktop-readability` 失败）。宽度适中比一味加宽更安全。
3. **lifecycle 终态越界**：宽画布下终态 y 会超出可读区域，按诊断用 `yOffset: -60` 上移即可。
4. **微线段**：一个节点有 3+ 条出边时自动端口展开可能产生 7px 拐点（< 8px 下限），对相邻同泳道的边加 `route: "straight"` 解决。
5. **endpoint 方向契约**：`fromSide: "bottom"` 要求第一段必须垂直向下；两端点未对齐时 straight 路由会产生斜线报错，此时**去掉手工 route/side 改回 auto**。
6. **schema 一票否决**：`additionalProperties: false` 全层级生效，拼写错误（如 `colour`）直接被拒；错误信息会带上最近的 id/label 定位。
7. **更新检查网络失败**：`scripts/check-update.mjs` 在本机网络下 `silent/check-failed`，按契约静默跳过，不影响任何功能；设 `ARCHIFY_UPDATE_CHECK_DISABLED=1` 可彻底关闭。
8. **workflow v1 vs v2**：新图一律 `schema_version: 2`（可读布局契约）；只有保留旧固定几何时才留 v1，且不要只改版本号——用 `migrate` 迁移。
9. **中文**：`meta.locale: "zh-CN"` 只管渲染器 UI；你写的内容（节点/标签/卡片）本来就要用中文写。渲染器个别固定 band 文案（如 lifecycle 的 "Interruptions + recovery"）仍是英文，属已知行为。
10. **Mermaid 不是渲染目标**：贴 Mermaid 是让它读拓扑语义后**重新创作** Archify JSON，不是套皮换肤。

---

## 11. 信息边界（何时不要用）

- Archify 不是通用绘图编辑器、不是 Mermaid 主题、不做自动 Mermaid 解析、不提供托管分享和 WYSIWYG 编辑
- `deployment-ownership` 工程档不会探测真实基础设施；缺事实就关掉或补事实，不要为过验证而删除档案
- Viewer 导出是传播资产，不能替代已验证 HTML、交付回执和真实视觉复查
