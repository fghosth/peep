# Archify 实战经验手册

> 来源：2026-09 在 pi agent 环境下从零安装并制作 6 张图（覆盖全部 5 种类型 + Delta 对比）的完整踩坑记录。
> 版本：archify v2.17.0-dev.1 · Node v24.18.0 · macOS + Chrome。
> 配套：使用方法见 `ARCHIFY-GUIDE.md`；本文只讲**官方文档没写、靠实际踩坑才知道的事**。

---

## 0. 一页速查：坑 → 对策

| # | 症状 | 一句话对策 |
|---|---|---|
| 1 | sequence 消息报 `keep y between 160 and 477` | 消息 y 超出可读时间线：减消息数或改 viewBox 高度 |
| 2 | `viewBox/1 must be >= 566/480` | viewBox 高度有**按内容动态计算的 schema 下限**，不是写死常数 |
| 3 | `does not honor fromSide "bottom"` | 手工 route/side 产生斜线段：**删掉手工控制改回 auto** |
| 4 | 终态/状态 `exceeds the vertical lifecycle area` | 按 `supportedFixes` 用 `yOffset: -60` 上移，别只加高画布 |
| 5 | `7px interior segment below the 8px micro-segment floor` | 一节点 3+ 出边触发端口展开：相邻同泳道的边加 `route: "straight"` |
| 6 | visual-check `viewport-overflow` 全视口失败 | 内容纵横比太高，压到 ~2:1（详见 §3 尺寸速查） |
| 7 | `desktop-readability` 失败、投影字号 4px | viewBox **太宽**触发 viewer 收窄阅读列到 ~930px，改窄画布/减内容 |
| 8 | 交付失败但旧 HTML 还在 | 这是**原子交付**特性：失败不覆盖旧文件，先修源再重跑 |
| 9 | 多写一个字段直接报错 | 全层级 `additionalProperties: false`，错误信息带最近 id/label 可定位 |
| 10 | 更新检查报 `silent/check-failed` | 网络受限时的正常表现，静默跳过即可，不影响任何功能 |
| 11 | lifecycle 中间 band 显示英文 "Interruptions + recovery" | 渲染器固定字符串，`zh-CN` 不覆盖它，属已知行为，非 bug |

---

## 1. 实测最顺的工作流

```text
写候选 JSON（全自动路由，零手工几何）
  → validate --json          ──失败──→ 读 diagnostics[].code/subject/supportedFixes
  →                            一次只修一处 → 重新 validate（循环）
  → validate 9/9 + 0 err + 0 warn
  → deliver --json           （原子提交，拿到 SHA-256 回执）
  → visual-check --json      （Chrome 实测 4 视口 + 深/浅截图）
  → fail? → 改内容纵横比/参数（见 §3）→ 重走 validate → deliver
  → 人工看截图（感知复查，自动化测不了审美）
```

**三条实测铁律：**

1. **自动优先**：`via / channelX / channelY / labelAt / fromSide / toSide / route` 一个都别先写。6 张图里所有手工几何控制都是被诊断"逼"出来的；先写的那个（lifecycle 的 `straight + fromSide/toSide`）反而是唯一产生斜线错误的。
2. **一次一处**：每轮只按一条诊断改一个点。连修两轮不见好就停，如实报告。
3. **validate 便宜，大胆循环**：单次 validate 秒级完成，可以用脚本批量试参数（见 §2.8）。

**实测耗时参考**（M 系列 Mac，无网络依赖）：validate < 1s，deliver 1–3s，visual-check 10–20s（起 Chrome 4 视口 + 4 截图）。

---

## 2. 踩坑完整档案（真实诊断 + 修复）

### 2.1 sequence：消息 y 超出可读时间线

**现象**：8 条消息、`viewBox: [820, 560]`，validate 失败。

```json
{ "code": "layout/constraint", "severity": "error",
  "message": "Message \"进入首页\" sits outside the readable timeline — keep y between 160 and 477." }
```

**根因**：消息 y 的可读上限 ≈ `viewBox[1] - 83`。y=480 > 477。
**修复**：把 viewBox 加高到 `[820, 600]`（上限变为 517）→ validate ✓。
**但注意**：这只过了 validate，visual-check 阶段又因首屏溢出被打回（见 2.6），最终方案是**减消息到 6 条 + `[1000, 480]`**。

### 2.2 lifecycle：viewBox 高度下限是动态的

```json
{ "code": "schema/minimum",
  "message": "/meta/viewBox/1 must be >= 566 {\"comparison\":\">=\",\"limit\":566}",
  "supportedFixes": ["use a value >= 566"] }
```

**根因**：viewBox 最小高度按状态数/泳道内容动态推导，不是 schema 里写死的常数。同样的图，宽度从 900 改到 1500 后下限也变。
**经验**：见到 `schema/minimum` 直接采纳 `supportedFixes` 给的数值（可略加余量）。

### 2.3 lifecycle：手工 route/side 产生斜线

**现象**：`paid→cancelled` 写了 `"route": "straight", "fromSide": "bottom", "toSide": "top"`，但终点列与起点列**没有垂直对齐**（终态列 N 对齐主列 N+2，而 paid 在主列 1），straight 画成对角线：

```json
{ "code": "clean-flow/endpoint-side-direction",
  "message": "... first segment [248,188] -> [402,450] does not honor fromSide \"bottom\"
              — it must run vertical downward; keep automatic routing, or choose
              fromSide/toSide and via points whose first and final segments cross
              state borders perpendicularly." }
```

**修复**：删掉这三个手工字段 → `{ "from": "paid", "to": "cancelled", "variant": "security" }`，auto 路由自己搞定。
**教训**：`fromSide: "bottom"` 是**方向契约**（第一段必须垂直向下出边），端点不对齐时 straight 必然斜线。官方示例里能用 straight 是因为它们的端点恰好对齐。

### 2.4 lifecycle：终态纵向越界 + yOffset

```json
{ "code": "layout/constraint",
  "message": "State \"cancelled\" exceeds the vertical lifecycle area — keep y between 64 and 448
              (adjust yOffset or increase meta.viewBox[1])." }
```

**关键**：诊断明确给了**两个**修法。加高 viewBox 到 640 能过布局，但会撞上首屏溢出（2.6）；正确解是 `yOffset`：

```json
{ "id": "cancelled", "lane": "terminal", "col": 0, "type": "failure", "yOffset": -60 }
```

终态上移 60px 后既满足纵向区域，又保住了扁画布。

### 2.5 lifecycle：7px 微线段

```json
{ "code": "composition/micro-segment",
  "message": "transitions[1] \"paid\" -> \"shipped\" has a 7px interior segment
              [325,150] -> [325,157] that is below the 8px micro-segment floor
              — move route/via or channel coordinates so each lifecycle turn has a readable run-up." }
```

**根因**：`paid` 有三条出边（shipped / cancelled / auto 展开），自动端口展开给相邻同泳道的水平边留了个 7px 的小拐点。
**修复**：这条边加 `"route": "straight"`（相邻同泳道，直线即正交，无副作用）。
**记忆点**：micro-segment 下限 8px；showcase 还要求内部段 ≥16px、端点段 ≥8px。

### 2.6 全类型：首屏溢出（最隐蔽的坑）

**现象**：validate 9/9、deliver 成功，但 visual-check 全视口失败：

```text
viewer/viewport-overflow → The rendered artifact overflows the 1440x900 light viewport.
containment: scroll 1440x1254 / inner 1440x900 (overflowY: true)
```

**重要事实**：官方自带的 sequence/lifecycle/dataflow 示例 HTML **同样溢出**（实测 scrollHeight 1586/1245/1302 vs 900），只有 architecture（自动 viewBox）通过。也就是说这是"高画布内容的通病"，不是你画错了。

**机理**：
- 页面高度 ≈ 头部(60) + 图显示高 + 图例(70) + 结论卡(170) + 边距
- 图显示高 = `viewBox[1] × 阅读列宽 / viewBox[0]`，**SVG 不会缩到 100% 以下**
- 所以 viewBox 纵横比（宽:高）接近或超过 2:1 才能首屏放下

**修复案例**：sequence 从 `[820,600]`+8 消息 → 6 消息+`[1000,480]`（纵横比 2.08）；lifecycle `[900,640]`（1.41，溢出 252px）→ `[1050,580]`+yOffset（1.81，通过）。

### 2.7 全类型：viewBox 太宽 → viewer 反向收窄 → 字号不达标

**现象**：把 sequence 拉宽到 `[1600,600]`（想靠缩放压高度），validate 报：

```json
{ "code": "composition/desktop-readability",
  "evidence": { "viewportWidth": 1440, "availableDiagramWidth": 930,
                "viewBoxWidth": 1600, "scale": 0.58125,
                "text": "浏览器", "sourceFontPx": 7,
                "projectedFontPx": 4.07, "minimumProjectedFontPx": 6 } }
```

**根因**：viewer 为保住首屏高度，会把阅读列**收窄**（实测收到 ~930px），缩放反而从 1.0 掉到 0.58，投影字号跌破 6px 下限。**加宽是反向操作。**
**经验**：宽度上限经验值 ~1300；想压高度靠"减内容/压扁语义布局"，不靠加宽画布。

### 2.8 兜底大招：脚本化网格搜索参数

viewBox/yOffset 这类参数组合空间小、validate 又快，直接写循环穷举，机器替你找可行解（本次 lifecycle 就是靠它找到 `1050×580 + yOffset -60`）：

```js
// node 循环调用 validate，找到第一个 pass 的组合
for (const w of [1050,1100,1150,1200])
  for (const h of [540,560,580,600])
    for (const dy of [-40,-60,-80]) {
      // 写临时 JSON → execSync validate --json → 解析 ok+composition.status
      if (pass) { 记录并写回正式文件; break outer; }
    }
```

注意：这只适合"数值微调"类问题；语义问题（缺状态、缺回路）穷举救不了。

### 2.9 schema：全层级严格模式

```text
workflow schema validation failed:
  /nodes/3 (id/label: "router") must NOT have additional properties {"additionalProperty":"colour"}
```

（官方 README 示例）多写、拼错字段直接被拒，**不会静默忽略**。好处是错误信息带实例路径 + 最近元素的 id/label，定位很快。改字段名前先翻 `schemas/<type>.schema.json` 和 `common.schema.json`。

### 2.10 原子交付的保护行为

`deliver` 失败（exit ≠ 0）时：**旧输出文件原样保留**、清理私有状态、绝不调用打开器。此时不要对着旧文件跑 visual-check——它测的是"上一版可信产物"，不是你刚被拒的候选。正确动作：读交付 JSON 里的 `error`/`diagnostics`，修源文件，重跑。

### 2.11 环境类

- **更新检查**：`scripts/check-update.mjs` 在无外网环境返回 `{"status":"silent","reason":"check-failed"}`，按契约静默继续即可；彻底关闭用 `ARCHIFY_UPDATE_CHECK_DISABLED=1`。
- **安装**（github.com 直连不通的网络）：全量 clone 会超时；用 `gh api -H "Accept: application/vnd.github.raw" repos/tt-a1i/archify/contents/archify.zip` 下载官方打包（1.3MB，77 文件，零依赖）解压到 skill 目录即可。
- **中文**：`meta.locale: "zh-CN"` 只管渲染器 UI；lifecycle 中间 band 的 "Interruptions + recovery" 等个别固定串仍为英文，官方未本地化，不算缺陷。

---

## 3. 首屏尺寸速查表（viewBox 经验值）

| 类型 | 推荐做法 | 实测通过组合 |
|---|---|---|
| architecture | **不写 viewBox**（渲染器按内容+图例自动算，天然首屏） | 5 组件图自动通过 |
| workflow | **不写 viewBox** | 3 泳道×6 列自动通过（viewBoxWidth 1032） |
| sequence | 内容压到 ~6-8 条消息；`[1000, 480]`；超宽画布配 `column_fit: "spread"`（但见 2.7） | `[1000,480]` + 6 消息 ✓ |
| dataflow | 五段以下、rows ≤4 时 `[1080, 560]` 安全 | ✓ |
| lifecycle | 状态 ≤7 个：`[1050, 580]` + 终态 `yOffset: -60` | ✓ |

**判定口诀**：validate 过 → deliver 过 → visual-check 还挂 → 先看纵横比（目标 ≥1.8:1），减内容或压高度；别加宽。

---

## 4. 诊断码速查表

| 诊断 code | 含义 | 典型修法 |
|---|---|---|
| `schema/additionalProperties` | 多写/拼错字段 | 对照 schema 改名或删除 |
| `schema/minimum` | 数值低于动态下限（viewBox 等） | 采纳 supportedFixes 数值 + 余量 |
| `layout/constraint` | 元素越出可读区域 | yOffset 微调 / 改 viewBox / 减内容 |
| `clean-flow/endpoint-side-direction` | 手工 side/route 产生非正交段 | 删手工控制回 auto |
| `composition/micro-segment` | 线段 < 8px（内部 <16px） | 相邻同向边改 `route:"straight"` |
| `composition/desktop-readability` | 投影字号 < 6px（viewer 收窄所致） | **减窄** viewBox / 减内容 |
| `composition/desktop-readability`(溢出) | 首屏 scrollHeight > innerHeight | 压纵横比 ≥1.8:1 |
| `viewer/viewport-overflow` | visual-check 实测溢出 | 同上，改完重走 validate→deliver |
| `label_route_clearance` 类 | 标签距路线太近 | 先挪位置/加间距，最后才缩短文案（语义标签不许删） |

---

## 5. 本次 6 张图的制作账本

| 图 | 类型 | 初始问题 | 修复轮数 | 最终参数 |
|---|---|---|---|---|
| CI/CD 发布流水线 | workflow v2 | 无（一次通过） | 0 | 无 viewBox |
| 博客平台架构 | architecture | 无（一次通过） | 0 | 无 viewBox |
| 用户登录时序 | sequence | 消息 y 越界 → 首屏溢出 → 收窄字号 | 3 | 6 消息 + `[1000,480]` |
| 订单数据管道 | dataflow | 无（一次通过） | 0 | `[1080,560]` |
| 订单状态机 | lifecycle | viewBox 下限 → 斜线 → 微线段 → 溢出 → 收窄 | 5（含网格搜索） | `[1050,580]` + 终态 yOffset:-60 + shipped 边 straight |
| checkout 对比 | compare | 无（官方 base/head 对） | 0 | —（对比产物本身有溢出特性，属官方已知） |

**规律**：类型越"纵向生长"（sequence 消息堆叠、lifecycle 多泳道带）越容易撞首屏约束；architecture/workflow 有自动尺寸，几乎不用操心。

---

## 6. 下次做图 Checklist

- [ ] 类型选对了吗？拿不准跑 `guide "<场景>" --json`
- [ ] 候选 JSON 里**没写**任何手工几何控制（via/channel*/labelAt/route/fromSide/toSide）
- [ ] 中文图写了 `meta.locale: "zh-CN"`，内容本身也是中文
- [ ] `quality_profile: "showcase"` 别漏
- [ ] validate 9/9 + 0 warn 才算过，"4 项通过"只是 basic
- [ ] deliver 成功才允许 visual-check；失败先看诊断，别测旧文件
- [ ] visual-check 挂了先查纵横比，再查 viewer 收窄（availableDiagramWidth 是否被压到 ~930）
- [ ] 最后**亲眼看截图**（深/浅两张）：自动化证明不了审美
- [ ] 汇报时三个声明分开写：deliver ✓ / visual-check ✓ / 人工视觉复查 ✓，轮数如实的 correction_rounds
