# musicbrainz-cli

<p align="center">
  <img src="doc/mbz.png" alt="mbz" height="200" />
  &nbsp;&nbsp;
  <img src="doc/cobra.png" alt="cobra" height="200" />
</p>

[English](./README.md) | [简体中文](./README.zh-CN.md)

通过 [MusicBrainz Web Service v2](https://musicbrainz.org/doc/MusicBrainz_API) 检索 artist、release 与 release group 元数据的命令行工具。输出为标准 JSON，便于脚本与管道集成。

## 功能特性

- **Search** — 使用 Lucene 语法搜索 artist / release / release group
- **按 artist MBID 搜 release / release group** — `search release` 与 `search releasegroup` 支持 `--artist-mbid` 过滤
- **仅 album** — `search release` 与 `search releasegroup` 自动在 API 查询中追加 `primarytype:album`
- **Lookup** — 按 MusicBrainz ID（MBID）精确查询 artist / release / release group
- **分页** — 通过 `--limit` / `--pageno` 控制（仅 search）
- **分数过滤** — 自动丢弃 score &lt; 50 的搜索结果
- **JSON 输出** — 成功结果输出至 stdout，错误输出至 stderr
- **输出模式** — 默认 `simple`（精简字段），可选 `full`（完整 API 结构）

## 环境要求

- Go 1.24 或更高版本
- 可访问 `https://musicbrainz.org`

## 安装

从源码安装：

```bash
git clone https://github.com/liuguancheng/musicbrainz-cli.git
cd musicbrainz-cli
go install .
```

或本地构建：

```bash
go build -o mbz .
```

可执行文件名为 `mbz`。

## 快速开始

```bash
# 搜索艺术家（默认精简 JSON）
mbz search artist "Beatles"

# 带分页搜索发行版
mbz search release 'release:"Abbey Road" AND artist:"Beatles"' --limit 5

# 按 MBID 查询
mbz lookup artist b10bbbfc-cf9e-42e6-888b-88b6b374d5d4

# 管道处理
mbz search artist "Beatles" | jq '.results[].artist'
```

## 命令参考

```
mbz
├── search
│   ├── artist <query>        搜索艺术家
│   ├── release [query]       搜索发行版（query 可选；可配合 --artist-mbid）
│   └── releasegroup [query]  搜索 release group（query 可选；可配合 --artist-mbid）
└── lookup
    ├── artist <mbid>         按 MBID 查询艺术家
    ├── release <mbid>        按 MBID 查询发行版
    └── releasegroup <mbid>   按 MBID 查询 release group
```

### 全局参数

| 参数 | 短选项 | 默认值 | 说明 |
|------|--------|--------|------|
| `--limit` | `-l` | `25` | 每页条数（1–100），仅 search 生效 |
| `--pageno` | `-p` | `1` | 页码（≥ 1），仅 search 生效 |
| `--output` | | `simple` | 输出模式：`simple` 或 `full` |
| `--user-agent` | | 自动生成 | HTTP User-Agent 覆盖 |
| `--contact` | | 项目仓库 URL | 联系方式，嵌入 User-Agent |
| `--api-url` | | `https://musicbrainz.org/ws/2/` | WS2 API 地址 |

### lookup 专属参数

| 参数 | 说明 |
|------|------|
| `--inc` | 可多次指定，请求附加关联数据（如 `releases`、`artist-credits`、`media`） |

### search release / releasegroup 专属参数

| 参数 | 说明 |
|------|------|
| `--artist-mbid` | 按 artist MBID 过滤（对应 Lucene `arid`）；可与可选的位置参数 query 组合 |

`search release` 与 `search releasegroup` 的 `query` 与 `--artist-mbid` 至少提供一个。

## 使用示例

### 搜索艺术家

```bash
mbz search artist "Beatles"
mbz search artist 'artist:"The Beatles"' --limit 10 --pageno 1
mbz search artist 'artist:"The Beatles"' --limit 10 --pageno 2   # 翻页
mbz search artist "Beatles" --output full                          # 完整 API JSON
```

### 搜索发行版

```bash
mbz search release 'release:"Abbey Road" AND artist:"Beatles"'
mbz search release "Abbey Road" --limit 5

# 按 artist MBID 列出该 artist 的全部 release
mbz search release --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4

# 关键词 + artist MBID 组合
mbz search release "Abbey Road" --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
```

### 搜索 release group

```bash
mbz search releasegroup "Abbey Road"
mbz search releasegroup --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
mbz search releasegroup "Abbey Road" --artist-mbid b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
```

### 按 MBID 查询

```bash
mbz lookup artist b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
mbz lookup release 464a321e-97a0-4654-8a7a-d1d88e8496e0 --inc artist-credits --inc media
mbz lookup releasegroup abbc4905-c25f-4c67-8e2d-19329ec48b1f
```

## 输出格式

### 输出模式

| 模式 | 参数 | 说明 |
|------|------|------|
| `simple` | `--output simple`（默认） | 抽取关键字段，缺失字段不出现在 JSON 中 |
| `full` | `--output full` | 完整 MusicBrainz API 实体结构 |

**simple 模式字段**（仅在 API 返回中存在时输出）：

`mbid`、`score`、`artist`、`release`、`releasegroup`、`type`、`country`、`date`、`format`、`barcode`、`alias`、`primary_alias`、`tag`

### Search 成功响应（simple）

Artist search 示例：

```json
{
  "type": "artist_search",
  "output": "simple",
  "query": "Beatles",
  "pageno": 1,
  "limit": 25,
  "min_score": 50,
  "count": 42,
  "current_count": 1,
  "has_data": true,
  "created": "2026-07-01T12:00:00Z",
  "results": [
    {
      "mbid": "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4",
      "score": 100,
      "artist": "The Beatles",
      "type": "Group",
      "country": "GB"
    }
  ]
}
```

Release search 会额外包含 `"primary_type": "album"`（通过 Lucene `primarytype:album` 在 API 侧过滤）：

```json
{
  "type": "release_search",
  "output": "simple",
  "query": "(Abbey Road) AND primarytype:album",
  "pageno": 1,
  "limit": 25,
  "min_score": 50,
  "primary_type": "album",
  "count": 42,
  "current_count": 1,
  "has_data": true,
  "created": "2026-07-01T12:00:00Z",
  "results": [
    {
      "mbid": "464a321e-97a0-4654-8a7a-d1d88e8496e0",
      "score": 100,
      "release": "Abbey Road"
    }
  ]
}
```

Release group search 使用相同 album 过滤；simple 结果以 `releasegroup` 表示标题：

```json
{
  "type": "releasegroup_search",
  "output": "simple",
  "query": "(Abbey Road) AND primarytype:album",
  "pageno": 1,
  "limit": 25,
  "min_score": 50,
  "primary_type": "album",
  "count": 42,
  "current_count": 1,
  "has_data": true,
  "created": "2026-07-01T12:00:00Z",
  "results": [
    {
      "mbid": "abbc4905-c25f-4c67-8e2d-19329ec48b1f",
      "score": 100,
      "releasegroup": "Abbey Road",
      "artist": "The Beatles",
      "type": "Album"
    }
  ]
}
```

### Search 成功响应（full）

包含原始 `results` 及顶层 `scores` 映射（MBID → score）。

### Lookup 成功响应（simple）

```json
{
  "type": "artist_lookup",
  "output": "simple",
  "id": "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4",
  "result": {
    "mbid": "b10bbbfc-cf9e-42e6-888b-88b6b374d5d4",
    "artist": "The Beatles"
  }
}
```

### 错误响应（stderr）

```json
{
  "error": "limit must be between 1 and 100",
  "code": "INVALID_ARGUMENT"
}
```

| 退出码 | 含义 |
|--------|------|
| `0` | 成功 |
| `1` | API 或运行时错误 |
| `2` | 参数校验失败 |

## 查询语法

搜索命令使用 [Apache Lucene 语法](https://lucene.apache.org/core/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description)，字段说明见 [MusicBrainz Search 文档](https://musicbrainz.org/doc/MusicBrainz_API/Search)。

| 场景 | 查询示例 |
|------|----------|
| 按艺术家名 | `artist:"The Beatles"` |
| 简单关键词 | `Beatles` |
| 按发行版名 + 艺术家 | `release:"Abbey Road" AND artist:"Beatles"` |
| 按条码 | `barcode:602527306377` |
| 按 artist MBID 搜 release / release group | `--artist-mbid b10bbbfc-...` 或 Lucene `arid:b10bbbfc-...` |

## 注意事项

1. **速率限制** — MusicBrainz 要求每个客户端最多 **1 次请求/秒**，请勿并发调用。
2. **User-Agent** — 必须提供有意义的 User-Agent（应用名、版本、联系方式）。
3. **分页** — 默认 `--limit=25`，`--pageno=1`；`limit` 有效范围 1–100。当 `(pageno - 1) * limit >= count` 时 `has_data` 为 `false`，表示已越页。
4. **计数字段** — `count` 为 API 返回的总命中数（跨分页；release search 含 `primarytype:album` 条件，不含 CLI score 过滤）。`current_count` 为本页经 score 过滤后实际输出的条数。`has_data` 在 `(pageno - 1) * limit < count` 时为 `true`。
5. **分数过滤** — search 自动丢弃 score &lt; 50 的结果；仅影响 `current_count` 与 `results`，不影响 `count` / `has_data`。
6. **Release search album 过滤** — `search release` 始终在发往 API 的 Lucene 查询中追加 `primarytype:album`，仅返回 album 类型发行版；JSON 中的 `primary_type` 字段标明该过滤条件。
7. **Release group search** — `search releasegroup` 使用相同 album 与 score 过滤；Lucene 文本默认搜索 `releasegroup` 字段；simple 输出以 `releasegroup`（非 `release`）表示标题。
8. **lookup** — lookup 命令不使用分页参数。

## 开发

运行测试：

```bash
go test ./...
```

构建：

```bash
go build -o mbz .
```

## 依赖

- [go.uploadedlobster.com/musicbrainzws2](https://pkg.go.dev/go.uploadedlobster.com/musicbrainzws2) — MusicBrainz WS2 Go 客户端
- [github.com/spf13/cobra](https://github.com/spf13/cobra) — CLI 框架

## License

MIT
