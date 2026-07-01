# musicbrainz-cli

<p align="center">
  <img src="doc/mbz.png" alt="mbz" height="200" />
  &nbsp;&nbsp;
  <img src="doc/cobra.png" alt="cobra" height="200" />
</p>

[English](./README.md) | [简体中文](./README.zh-CN.md)

通过 [MusicBrainz Web Service v2](https://musicbrainz.org/doc/MusicBrainz_API) 检索 artist 与 release 元数据的命令行工具。输出为标准 JSON，便于脚本与管道集成。

## 功能特性

- **Search** — 使用 Lucene 语法搜索 artist / release
- **Lookup** — 按 MusicBrainz ID（MBID）精确查询
- **分页** — 通过 `--limit` / `--offset` 控制（仅 search）
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
│   ├── artist <query>     搜索艺术家
│   └── release <query>    搜索发行版
└── lookup
    ├── artist <mbid>      按 MBID 查询艺术家
    └── release <mbid>     按 MBID 查询发行版
```

### 全局参数

| 参数 | 短选项 | 默认值 | 说明 |
|------|--------|--------|------|
| `--limit` | `-l` | `25` | 每页条数（1–100），仅 search 生效 |
| `--offset` | `-o` | `0` | 起始偏移（≥ 0），仅 search 生效 |
| `--output` | | `simple` | 输出模式：`simple` 或 `full` |
| `--user-agent` | | 自动生成 | HTTP User-Agent 覆盖 |
| `--contact` | | 项目仓库 URL | 联系方式，嵌入 User-Agent |
| `--api-url` | | `https://musicbrainz.org/ws/2/` | WS2 API 地址 |

### lookup 专属参数

| 参数 | 说明 |
|------|------|
| `--inc` | 可多次指定，请求附加关联数据（如 `releases`、`artist-credits`、`media`） |

## 使用示例

### 搜索艺术家

```bash
mbz search artist "Beatles"
mbz search artist 'artist:"The Beatles"' --limit 10 --offset 0
mbz search artist 'artist:"The Beatles"' --limit 10 --offset 10   # 翻页
mbz search artist "Beatles" --output full                          # 完整 API JSON
```

### 搜索发行版

```bash
mbz search release 'release:"Abbey Road" AND artist:"Beatles"'
mbz search release "Abbey Road" --limit 5
```

### 按 MBID 查询

```bash
mbz lookup artist b10bbbfc-cf9e-42e6-888b-88b6b374d5d4
mbz lookup release 464a321e-97a0-4654-8a7a-d1d88e8496e0 --inc artist-credits --inc media
```

## 输出格式

### 输出模式

| 模式 | 参数 | 说明 |
|------|------|------|
| `simple` | `--output simple`（默认） | 抽取关键字段，缺失字段不出现在 JSON 中 |
| `full` | `--output full` | 完整 MusicBrainz API 实体结构 |

**simple 模式字段**（仅在 API 返回中存在时输出）：

`mbid`、`score`、`artist`、`release`、`type`、`country`、`date`、`format`、`barcode`、`alias`、`primary_alias`、`tag`

### Search 成功响应（simple）

```json
{
  "type": "artist_search",
  "output": "simple",
  "query": "Beatles",
  "offset": 0,
  "limit": 25,
  "min_score": 50,
  "count": 1,
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

## 注意事项

1. **速率限制** — MusicBrainz 要求每个客户端最多 **1 次请求/秒**，请勿并发调用。
2. **User-Agent** — 必须提供有意义的 User-Agent（应用名、版本、联系方式）。
3. **分页** — 默认 `--limit=25`，`--offset=0`；`limit` 有效范围 1–100。
4. **分数过滤** — search 自动丢弃 score &lt; 50 的结果；`count` 为过滤后的条数。
5. **lookup** — lookup 命令不使用分页参数。

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
