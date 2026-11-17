# Yukumo Script Web/WASM 适配计划

- 状态：Implemented
- 日期：2026-11-16
- 分支：`wasm-port`（当前重构分支）
- 前置：[1-layer-refactor.md](1-layer-refactor.md)（Implemented）
- 相关计划：[2-cli-layer-refactor.md](2-cli-layer-refactor.md)（CLI 内部分层；已合入，本计划的 CLI 落点按它的 action/schema/ui 结构实现）
- 实施记录：
  - `88e62cd` chore(plan): split wasm port plan from layer refactor
  - `51addf9` refactor(port): add core/port and inject stores into core services
  - `c2fbaed` refactor(synth): move artifact io behind artifact store
  - `311ea93` refactor(api): split pure api from cliopen composition root
  - `e5a085f` refactor(cli): adopt cliopen and bytes-only file inputs
  - `6d938e1` test(wasm): add node headless smoke harness
  - `8ff61ef` chore(ci): add wasm gates and port module checks
  - `db30838` refactor(api): expose cliopen defaults
  - `a2a9afc` test(core): cover fs store and pure api wiring
  - `afc7ed3` style(api): gofmt app struct
- 最终检查：`just ci`、`just wasm-smoke`、`just wasm-lint` 通过；`just check-replaces`、`just check-cli-imports` 通过。
- 范围：`core/port` 契约、依赖注入、`core/api` 去 IO、平台组合根、synth 的 IO port 化、bytes-only 文件入口、Node headless 冒烟、CI 门禁。
- 非范围：产品级 `web/` module、IndexedDB、Web Audio、浏览器 UI、WASI 目标、`runtime.NewWithWASM`。
- 约束：CLI 可以修改，但 CLI 的输出文本、交互流程、`-o` / `-p` / `-P` 行为、合成策略、YAML 文件格式与现状一致；有意差异集中在 §19。
- 结论：分层重构已经把 store / service / runtime 拆开，但 **port 不可见、组合根带 runtime、synth 直接 OS IO** 这三件事仍把 wasm 挡在外面。本计划补的是 seam，不是重写：
  1. 新增独立 public module `core/port`；
  2. `core/api` 只保留纯用例，平台组合根拆到 `core/api/cliopen`；
  3. synth 的 WAV / 缓存 / index 全部走 `ArtifactStore`，`-o` 走 `OutputWriter`；
  4. `-P` 与 phont 改成 bytes 入口，CLI 负责读文件；
  5. 用 `wasmtest/`（Node + `GOOS=js GOARCH=wasm`）验证 JS 实现 `port.Engine` 的路线。

## 目录

- [1. 最终决策](#sec-1)
- [2. 目标依赖与 wasm 边界](#sec-2)
- [3. 目标目录](#sec-3)
- [4. `core/port` 契约](#sec-4)
  - [4.1 Store](#sec-4-1)
  - [4.2 ArtifactStore 与 OutputWriter](#sec-4-2)
  - [4.3 Engine 与 EngineParam](#sec-4-3)
  - [4.4 Flusher、Closer 与错误](#sec-4-4)
  - [4.5 key 语义](#sec-4-5)
- [5. Store 实现与分层](#sec-5)
  - [5.1 纯 store](#sec-5-1)
  - [5.2 fs store](#sec-5-2)
  - [5.3 ArtifactStore 的 key 布局](#sec-5-3)
- [6. API 层设计](#sec-6)
  - [6.1 `Deps` 与 `New`](#sec-6-1)
  - [6.2 `App` 生命周期](#sec-6-2)
  - [6.3 DTO 与 `SynthRequest`](#sec-6-3)
  - [6.4 文件入参：bytes-only](#sec-6-4)
- [7. 平台组合根 `core/api/cliopen`](#sec-7)
- [8. CLI 改造](#sec-8)
- [9. synth 的 IO 改造](#sec-9)
  - [9.1 合成流程](#sec-9-1)
  - [9.2 缓存与输出](#sec-9-2)
  - [9.3 错误与路径返回](#sec-9-3)
- [10. 同步 Store + 内存缓存 + Flush](#sec-10)
- [11. Engine 的 JS 适配契约](#sec-11)
- [12. wasm 打包边界](#sec-12)
- [13. 测试分层](#sec-13)
- [14. 静态检查与 depguard](#sec-14)
- [15. justfile 与 CI](#sec-15)
- [16. Node headless 冒烟](#sec-16)
- [17. 迁移步骤](#sec-17)
- [18. 验收清单](#sec-18)
- [19. 风险与取舍](#sec-19)
- [20. 后续：真正的 Web 实现](#sec-20)
- [附录 A：文件映射](#sec-a)
- [附录 B：命令速查](#sec-b)

---

<a id="sec-1"></a>
## 1. 最终决策

| 决策 | 结果 |
|---|---|
| 目标 | 让 Go 业务代码可以在 `GOOS=js GOARCH=wasm` 下运行；产品级浏览器 UI 不在本计划 |
| port 位置 | 新增独立 module `yukumo-script/core/port`（目录 `core/port/`） |
| port 可见性 | `Store` / `ArtifactStore` / `OutputWriter` / `Engine` / `EngineParam` 全部 public |
| Store 划分 | `ConfigStore`、`CharacterStore`、`ArtifactStore` 三个实例由组合根注入 |
| 持久化模型 | 同步 Store + 内存缓存 + `Flush`；IndexedDB 在 JS adapter 内异步写回 |
| Engine | 浏览器/Node 由 JS 实现 `port.Engine`；wasm build 不链接 `runtime` / wazero |
| API | `core/api` 只保留纯用例、DTO、`New(deps)`；`Open` 移到 `core/api/cliopen` |
| synth IO | WAV / 缓存 / `index.yaml` 走 `ArtifactStore`；`-o` 走 `OutputWriter` |
| 文件入参 | API 不再接收文件路径：`-P` 由 CLI 读 bytes，phont 改成 `SetPhontData` |
| OS 默认值 | `DefaultCharacterPath` / `DefaultOutputPath` / `EnsurePaths` 全部移到 `cliopen` / CLI |
| core 纯度 | `domain` / `service` / `core/api` 根包不 import `os`、`path/filepath`、`runtime`、wazero |
| 验证 | native 测试全绿 + `GOOS=js` build 闸门 + Node headless 冒烟（fake engine + 真 `aq_wasm.wasm`） |
| CI | `just ci` 保持 native；新增独立 `wasm` job 跑 `just wasm-smoke` / `just wasm-lint` |
| CLI | 可修改；行为不变，唯一有意差异见 §19.3 |
| 与旧 addendum | 原 [1-layer-refactor.md](1-layer-refactor.md) §16 收敛为指针，内容由本文件维护 |

目标依赖：

```text
cli
  ↓
core/api/cliopen ── os / core/internal/store/fs / runtime/aq → wazero
  ↓
core/api ── DTO + 用例 + New(deps)
  ↓
core/internal/service ── 依赖 core/port，不依赖任何具体实现
  ↓
core/internal/domain

wasmtest（GOOS=js）
  ↓
core/api → core/port
  （不 import core/api/cliopen、core/internal/store/fs、runtime、cli）
```

`core/port` 的 go.mod 只有 stdlib：

```text
yukumo-script/core/port     → （无第三方依赖）
```

`core/api` 根包在 js 下的依赖闭包不出现 wazero / runtime / os；只有 `core/api/cliopen` 会带 OS 与 runtime。

---

<a id="sec-2"></a>
## 2. 目标依赖与 wasm 边界

### 2.1 module

| module | 目录 | 职责 | 直接依赖 |
|---|---|---|---|
| `yukumo-script/core/port` | `core/port/` | Store / ArtifactStore / Engine 契约 | stdlib |
| `yukumo-script/core` | `core/` | domain + service + 纯 store | core/port、yaml |
| `yukumo-script/runtime` | `runtime/` | wazero、Yukw WASM、Pool、AQ 参数编码 | wazero |
| `yukumo-script/core/api` | `core/api/` | DTO、用例门面、纯 `New(deps)`；`cliopen` 子包是 CLI 平台组合根 | core、core/port、runtime（仅 cliopen）、yaml |
| `yukumo-script/cli` | `cli/` | cobra/huh/lipgloss、prompt、播放、展示 | core/api、core/api/cliopen、UI 库 |
| `yukumo-script/wasmtest` | `wasmtest/` | Node headless 冒烟；JS store / engine 参考 adapter | core/api、core/port、`syscall/js` |

### 2.2 依赖方向

```text
cli → core/api/cliopen → core/api → core/internal/service → core/port
                      ↘ core/internal/store/fs → runtime/aq → wazero

wasmtest → core/api → core/port
```

禁止方向：

```text
core/port     → 任何项目包或第三方库
core          → core/api / cli / runtime / wazero
runtime       → 任何 yukumo-script 包（包括 core/port）
core/api 根包 → runtime / wazero / cli / os / path/filepath
core/api/cliopen → cli
cli           → core/internal / runtime / core/port / wazero
wasmtest      → runtime / wazero / cli / core/internal / core/api/cliopen
```

### 2.3 为什么 port 必须是独立 module

- `core/internal/service` 与 `core/api` 都要 import port；如果 port 留在 `core` module 内，web/ 只 import `core/port` 也会把 `core` 的 `go.mod` 依赖（yaml）带进 module 图。
- 独立 module 后，`wasmtest` 与未来 `web/` 的 `go list -deps` 只出现 port 自身，不会出现 yaml / core。
- 代价是 `go.work`、replace、`check-replaces` 多一项，见 §15、§17。

---

<a id="sec-3"></a>
## 3. 目标目录

```text
.
├── go.work
├── go.work.sum
├── justfile
├── README.md
├── README.zh.md
├── docs/
│   └── refactorings/
│       ├── 1-layer-refactor.md          # §16 改成指针
│       ├── 2-cli-layer-refactor.md      # 已存在；本计划不改
│       └── 3-wasm-port.md               # 本文件
├── scripts/
│   ├── check-replaces.go
│   ├── prepare-wasm.go                  # 复制 GOROOT/lib/wasm/wasm_exec.js
│   └── check-wasm-deps.go               # GOOS=js 依赖闭包断言
├── core/
│   ├── port/                            # module yukumo-script/core/port
│   │   ├── go.mod
│   │   ├── .golangci.yml
│   │   ├── store.go                     # Store / Reader / Writer / Lister / Deleter / ErrNotFound
│   │   ├── artifact.go                  # ArtifactStore / OutputWriter / Flusher
│   │   └── engine.go                    # Engine / EngineParam / PresetCatalog / Closer
│   ├── go.mod
│   ├── go.sum
│   ├── .golangci.yml
│   ├── internal/
│   │   ├── domain/                      # 不再 import os
│   │   ├── store/                       # 纯实现：memory.go、yaml.go
│   │   │   └── fs/                      # OS 实现：store.go、artifact.go、output.go
│   │   └── service/
│   │       ├── config/
│   │       ├── character/
│   │       │   └── voice/
│   │       ├── nameutil/
│   │       └── synth/                   # 不再 import os / filepath
│   ├── api/
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── .golangci.yml
│   │   ├── app.go                       # Deps / New / App / Close / CleanIndex
│   │   ├── model.go                     # DTO + DefaultConfig + SetPhontData
│   │   ├── convert.go
│   │   ├── config.go
│   │   ├── character.go
│   │   ├── voice.go
│   │   ├── synth.go
│   │   ├── errors.go
│   │   ├── cliopen/                     # 平台组合根
│   │   │   ├── open.go                  # Open / EnsurePaths / resolveConfigDir
│   │   │   ├── defaults.go              # OS 默认路径
│   │   │   └── aqengine.go              # port.Engine ↔ aq.Pool
│   │   └── test/
│   │       ├── runtime/aq/
│   │       └── synth/
│   └── test/
│       ├── service/...
│       └── store/fs/
├── runtime/
│   ├── go.mod
│   └── aq/
├── cli/                                 # 见 §8
│   ├── go.mod
│   ├── cmd/yukumo-script/main.go
│   ├── root.go
│   ├── ...
│   └── synth/voice_file.go              # -P 解析 + 读文件
├── wasmtest/                            # module yukumo-script/wasmtest
│   ├── go.mod
│   ├── .golangci.yml
│   ├── main.go                          # 所有 .go 文件 //go:build js
│   ├── smoke.go
│   ├── jsstore.go
│   ├── jsartifact.go
│   ├── jsengine.go
│   ├── fakeengine.go
│   ├── store.cjs
│   ├── aq_engine.cjs
│   ├── harness.cjs
│   └── .build/                          # 生成物，gitignore
│       ├── wasm_exec.js
│       └── smoke.wasm
└── third-party/
    ├── aq_wasm.h
    └── aq_wasm.wasm
```

说明：

- `core/port` 是独立 module，`core`、`core/api`、`cli`、`wasmtest` 都要在各自 `go.mod` 里 require + replace 它。
- `core/internal/store` 保持纯包（Memory + YAML helper）；OS 实现拆到 `core/internal/store/fs`，只有 `core/api/cliopen` import 它。
- `core/api` 根包不再有 `Open`、`aqengine`、`os`；CLI 通过 `core/api/cliopen` 拿 FS + runtime 组合根。
- `wasmtest` 只在 `GOOS=js` 下构建，不进入 native build/test。
- `third-party/aq_wasm.wasm` 位置不变；Node harness 直接读它。

---

<a id="sec-4"></a>
## 4. `core/port` 契约

<a id="sec-4-1"></a>
### 4.1 Store

```go
package port

import "errors"

var ErrNotFound = errors.New("store: not found")

type Reader interface {
	Read(key string) ([]byte, error)
}

type Writer interface {
	Write(key string, data []byte) error
}

type ExistenceChecker interface {
	Exists(key string) (bool, error)
}

type Lister interface {
	List(prefix string) ([]string, error)
}

type Deleter interface {
	Delete(key string) error
}

type Store interface {
	Reader
	Writer
	ExistenceChecker
	Lister
	Deleter
}
```

与当前 `core/internal/store.Store` 的差异：

- 接口从 internal 移到 public `core/port`；`core/internal/store` 只保留实现与 YAML helper。
- `List` 的参数和返回都按逻辑 key 命名，不再叫 `dir` / `name`。
- 缺 prefix 的 `List` 返回空 slice（不是 `ErrNotFound`），保持现有 FS 行为。
- `Memory.List` 必须排序，保证跨 adapter 的确定顺序。

<a id="sec-4-2"></a>
### 4.2 ArtifactStore 与 OutputWriter

```go
type Flusher interface {
	Flush(ctx context.Context) error
}

// ArtifactStore 是 synth 的输出与缓存存储。key 是逻辑 key，不是绝对路径。
type ArtifactStore interface {
	Store
	// Path 只用于展示 / 日志：FS 返回绝对路径，内存/JS 返回 key 或空串。
	// 调用方不得把 Path 的结果再当 key 传回。
	Path(key string) string
}

// OutputWriter 只负责显式的 -o 路径。CLI 注入 FS 实现；web 不注入。
type OutputWriter interface {
	Write(path string, data []byte) error
}
```

为什么 `Path` 在 port 上：

- CLI 的 `synth` 输出行要打印真实文件路径；`SynthResult.Path` 必须能来自 `ArtifactStore`。
- 把 `Path` 与 `Read/Write` 放在同一个 adapter 里，避免 service 再 import `filepath`。
- web adapter 的 `Path` 返回 key 或空串；web 不打印路径。
- `-o` 的绝对路径不经过 ArtifactStore，走 `OutputWriter`，保持 key 语义干净。

<a id="sec-4-3"></a>
### 4.3 Engine 与 EngineParam

```go
type EngineParam struct {
	Version int
	Preset  int
	Speed   int
	Phont   []byte
	Base    int
	Volume  int
	Pitch   int
	Accent  int
	Lmd     int
	Fsc     int
}

type Synthesizer interface {
	Synthesize(ctx context.Context, text string, p EngineParam) ([]byte, error)
}

type PresetCatalog interface {
	HasPreset(ctx context.Context, version, preset int32) (bool, error)
	PresetNames(ctx context.Context, version int32) ([]string, error)
}

type Engine interface {
	Synthesizer
	PresetCatalog
}

type Closer interface {
	Close(ctx context.Context) error
}
```

说明：

- `Engine` 同时覆盖合成与 preset 查询；`voice.Service` 只需要 `PresetCatalog`，`synth.Service` 只需要 `Synthesizer`。
- `runtime/aq.Param` 不实现 `port.EngineParam`；两者保持独立，由 `core/api/cliopen/aqengine.go` 做字段映射（约 20 行），避免 runtime 反向依赖 port。
- `engine` 在 `Deps` 里可为 nil：此时配置 / 角色 / 音色 CRUD 仍可用，`Synthesize` / `PresetNames` 返回明确错误。

<a id="sec-4-4"></a>
### 4.4 Flusher、Closer 与错误

- `Flusher` 可选；`App.Close` 对实现它的 store 逐个 `Flush`。
- `Closer` 可选；`App.Close` 对实现它的 engine 调 `Close`。
- `port.ErrNotFound` 是 store 层错误；`domain.ErrNotFound` 是业务层错误。`api` 继续暴露 `domain` 的 `ErrInvalid` / `ErrNotFound` / `ErrExists`，不要把 `port.ErrNotFound` 直接泄漏给 CLI。
- `core/internal/store` 保留：

  ```go
  var ErrNotFound = port.ErrNotFound
  ```

  让内部测试和 YAML helper 的 `errors.Is` 不用改两套常量。

<a id="sec-4-5"></a>
### 4.5 key 语义

| 操作 | 契约 |
|---|---|
| `Read` | 返回完整 bytes；不存在返回 `ErrNotFound` |
| `Write` | 原子 replace；成功返回后读到完整新值 |
| `Exists` | 不存在返回 `false, nil`，不是错误 |
| `List(prefix)` | 返回 prefix 下的 **immediate** 逻辑 key，按确定顺序排序 |
| `Delete` | 不存在返回 `ErrNotFound` |
| key | 逻辑 key，使用 `/`，不以 `/` 开头，不含 `..`，不假设平台分隔符 |
| 并发 | 由 adapter 保证；core 不假设全局锁 |
| `Path` | 仅展示；FS 返回绝对路径，内存/JS 返回 key 或空串 |
| `Flush` | 可选；web adapter 在关闭/检查点写回 IndexedDB |

FS adapter 必须拒绝会逃出 root 的 key（绝对路径、`..`、Windows drive 前缀）；`OutputWriter` 的 path 不受此限制（它明确是 OS 路径）。

---

<a id="sec-5"></a>
## 5. Store 实现与分层

<a id="sec-5-1"></a>
### 5.1 纯 store

`core/internal/store`（不 import `os`）：

```text
store.go      # type aliases: Reader/Writer/... = port.*; ErrNotFound = port.ErrNotFound
memory.go     # Memory：map + RWMutex；List 排序；Path(key)=key；Flush no-op
yaml.go       # LoadYAML / SaveYAML / LoadYAMLOrDefault，参数改成 port.Reader / port.Writer
```

`Memory` 同时实现 `port.Store`、`port.ArtifactStore`、`port.Flusher`，用于：

- `core` 的 service 单元测试；
- `core/api` 的 `app_test.go`；
- `synth` 的 filename/cache 单元测试。

<a id="sec-5-2"></a>
### 5.2 fs store

`core/internal/store/fs`（只有这个包 import `os` / `path/filepath`）：

```go
type Store struct{ root string }

func New(root string) *Store
func (s *Store) Read(key string) ([]byte, error)
func (s *Store) Write(key string, data []byte) error
func (s *Store) Exists(key string) (bool, error)
func (s *Store) List(prefix string) ([]string, error)
func (s *Store) Delete(key string) error
func (s *Store) Path(key string) string
func (s *Store) Root() string

type Output struct{}

func (Output) Write(path string, data []byte) error // MkdirAll + WriteFile
```

行为与当前 `store.FS` 保持一致：

- `Write` 先 `MkdirAll`，写 `.tmp`，再 `Rename`；
- `Read` / `Delete` 把 `fs.ErrNotExist` 映射为 `port.ErrNotFound`；
- `List` 缺目录返回空 slice；
- `Path` = `filepath.Join(root, filepath.FromSlash(key))`，并对绝对 key / `..` 做拒绝。

`fs.Store` 同时实现 `port.ArtifactStore`；`fs.Output` 实现 `port.OutputWriter`。

<a id="sec-5-3"></a>
### 5.3 ArtifactStore 的 key 布局

CLI 的 ArtifactStore root 仍然是 `config.paths.output_path`，因此磁盘布局保持：

```text
<output_path>/
├── index.yaml                         # hash → key
├── remilia_default_2026-11-16_....wav # 正常生成输出
└── cache-<hash>.wav                   # 仅 -o 场景的独立缓存副本
```

规则：

- 生成输出的 key 是 `<sanitized-name>.wav`，直接放在 ArtifactStore root；`Path` 返回 `<output_path>/<key>`。
- `index.yaml` 是 hash → key 的映射；key 对 FS 是相对文件名，对 JS 是逻辑 key。
- `-o` 写文件走 `OutputWriter`，不进入 ArtifactStore；为了让 `-o` 也能命中缓存，缓存写入 `cache-<hash>.wav`，而不是把绝对路径塞进 index。
- 旧 `index.yaml` 条目仍然有效：它们的 key 就是当前输出文件 basename；`CleanIndex` 继续按 `Exists` 清理。
- 没有 ArtifactStore 时（`Deps.Artifacts == nil`）：不做缓存，也不能生成默认输出；只允许 `-o` + `OutputWriter` 或纯 bytes 返回。

`-o` 的缓存是本次有意修正：当前实现把 `filepath.Base(-o path)` 写进 index，但文件在 output root 之外，缓存实际不可命中。新实现用独立 cache key，行为可预期。

---

<a id="sec-6"></a>
## 6. API 层设计

<a id="sec-6-1"></a>
### 6.1 `Deps` 与 `New`

```go
package api

type Deps struct {
	ConfigStore    port.Store
	CharacterStore port.Store
	Artifacts      port.ArtifactStore // 可选；nil 表示不缓存、不写默认输出
	Outputs        port.OutputWriter  // 可选；nil 表示不支持 -o
	Engine         port.Engine        // 可选；nil 表示合成/预设查询不可用
	Defaults       Config             // 缺省配置；Paths 由平台层填
}

func New(deps Deps) (*App, error)
```

`New` 只做构造，不做 IO：

1. 校验 `ConfigStore` / `CharacterStore` 非 nil；
2. `config.New(deps.ConfigStore, fromConfig(deps.Defaults))`；
3. `character.New(cfgSvc, deps.CharacterStore)`；
4. `voice.New(charSvc, deps.Engine)`；
5. `synth.New(cfgSvc, voiceSvc, deps.Engine, deps.Artifacts, deps.Outputs)`；
6. 保存 engine 与三个 store 的引用，供 `Close` 使用。

> `Deps.Defaults` 应该传 `api.DefaultConfig()`（可再改 `Paths`），不要用零值；零值会把 `default_interactive` / `always_write` / `reuse_cache` 的默认值变成 false。

`CleanIndex` 不在 `New` 里自动执行；由平台组合根显式调用，避免构造函数做 IO：

```go
func (a *App) CleanIndex() error
```

<a id="sec-6-2"></a>
### 6.2 `App` 生命周期

```go
type App struct {
	config  *config.Service
	chars   *character.Service
	voices  *voice.Service
	synth   *synth.Service
	engine  port.Engine
	stores  []port.Store
}

func (a *App) Close(ctx context.Context) error
```

`Close` 顺序：

1. 若 `engine` 实现 `port.Closer`，关闭 engine；
2. 对 `stores` 中实现 `port.Flusher` 的逐个 `Flush(ctx)`；
3. 用 `errors.Join` 收集错误。

CLI 的 `defer app.Close(context.Background())` 不变；`cliopen.Open` 创建的 `aq.Pool` 由 `aqEngine.Close` 关闭。

<a id="sec-6-3"></a>
### 6.3 DTO 与 `SynthRequest`

DTO 保持现有字段（`yaml` tag 不变），只做两处调整：

```go
type SynthRequest struct {
	VoiceRef    string // character:voice（与 VoiceFile 二选一）
	VoiceFile   []byte // 已读取的 Character YAML bytes
	VoiceName   string // VoiceFile 模式下必填
	VoiceBase   string // VoiceFile 模式下 character.name 为空时的 fallback
	Text        string
	Version     int
	SpeedFactor float64
	Output      string
	Preview     bool
	ReuseCache  *bool
}

type SynthResult struct {
	WAV       []byte
	Path      string
	Wrote     bool
	FromCache bool
}
```

- `VoicePathSpec` 删除；CLI 在调用前解析并读文件。
- `SynthResult` 保持不变，CLI 的输出行逻辑不用改。
- `api.Config` / `Character` / `Voice` / `SynthParam` 的 YAML 字段与现有 DTO 完全一致。
- `api.DefaultConfig()` 返回 **不含 OS 路径** 的默认值：

  ```go
  func DefaultConfig() Config {
      return toConfig(domain.DefaultConfig(domain.PathConfig{}))
  }
  ```

  `cliopen` 再把 OS 路径填进 `Config.Paths`。

<a id="sec-6-4"></a>
### 6.4 文件入参：bytes-only

删除：

```go
func SetPhontPath(p *SynthParam, path string) error
```

新增：

```go
func SetPhontData(p *SynthParam, data []byte) error
```

规则：

- `p == nil` → `ErrInvalid`；
- `len(data) == 0` → `p.PhontData = nil`（等价于清除）；
- 非空 → 深拷贝到 `p.PhontData`；
- `.phont` 扩展名 / 文件存在性 / 文件为空的校验移到 CLI（`cli/character/voice/phont.go`）。

`-P` 的 `path:voice` 解析也移到 CLI：

```go
// cli/synth/voice_file.go
func parseVoicePathSpec(spec string) (yamlPath, voiceName string, err error)
func loadVoiceFile(yamlPath, voiceName string) (data []byte, name, base string, err error)
func baseNameNoExt(path string) string
```

`core/internal/service/synth` 只保留纯函数：

```go
func LoadVoiceBytes(data []byte, voiceName, base string) (domain.Character, domain.Voice, []byte, error)
```

`base` 为空时，`character.name` 为空则回退到 `"character"`；CLI 传入文件名 stem，保持现有行为。

---

<a id="sec-7"></a>
## 7. 平台组合根 `core/api/cliopen`

`core/api/cliopen` 是唯一的 FS + runtime 组合根：

```go
package cliopen

type Options struct {
	ConfigDir     string
	CharacterPath string // 可选；测试/嵌入覆盖
	OutputPath    string // 可选；测试/嵌入覆盖
}

func Open(opts Options) (*api.App, error)
func EnsurePaths(app *api.App) error
func Defaults() (api.Config, error) // 含 OS 路径
```

`Open` 流程：

1. `resolveConfigDir(opts.ConfigDir)`：`Options.ConfigDir` > `YUKUMO_CONFIG_DIR` > `os.UserConfigDir()/yukumo-script`。
2. `cfgStore := fs.New(dir)`。
3. `defaults, err := domainDefaults()`（OS 默认路径）。
4. `cfgSvc := config.New(cfgStore, defaults)`；`cfg, err := cfgSvc.Get()`；失败返回错误。
5. 若 `opts.CharacterPath` / `opts.OutputPath` 非空，覆盖 `cfg.Paths`（测试/嵌入用）。
6. `charStore := fs.New(cfg.Paths.CharacterPath)`；`artifacts := fs.New(cfg.Paths.OutputPath)`。
7. `outputs := fs.Output{}`。
8. `pool := aq.NewPool()`；`engine := aqEngine{pool: pool}`。
9. `app, err := api.New(api.Deps{ ... Defaults: apiDefaults(defaults) })`；失败时 `pool.Close` 并返回。
10. `ensurePaths(cfg.Paths.CharacterPath, cfg.Paths.OutputPath)`；失败时 `app.Close` 并返回。
11. `app.CleanIndex()`；失败时 `app.Close` 并返回。
12. 返回 `app`。

`EnsurePaths`：

```go
func EnsurePaths(app *api.App) error {
	cfg, err := app.Config()
	if err != nil {
		return err
	}
	for _, p := range []string{cfg.Paths.CharacterPath, cfg.Paths.OutputPath} {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if err := os.MkdirAll(p, 0o755); err != nil {
			return fmt.Errorf("create config path %q: %w", p, err)
		}
	}
	return nil
}
```

`aqEngine`：

```go
type aqEngine struct{ pool *aq.Pool }

func (e aqEngine) Synthesize(ctx context.Context, text string, p port.EngineParam) ([]byte, error)
func (e aqEngine) HasPreset(ctx context.Context, version, preset int32) (bool, error)
func (e aqEngine) PresetNames(ctx context.Context, version int32) ([]string, error)
func (e aqEngine) Close(ctx context.Context) error
```

`Open` 不再 panic：`os.UserConfigDir` 失败、`EnsurePaths` 失败都以 error 返回（当前 `domain.DefaultCharacterPath` 会 panic）。

---

<a id="sec-8"></a>
## 8. CLI 改造

CLI 行为保持不变，但以下文件需要跟着 API 调整：

| 文件 | 改动 |
|---|---|
| `cli/root.go` | `api.Open` → `cliopen.Open`；其余不变 |
| `cli/action/config.go` | `app.SaveConfig` 成功后调用 `cliopen.EnsurePaths(app)`（CLI 分层后实际落点） |
| `cli/parse/parse.go` | 新增 `ReadPhont(path) ([]byte, error)`，校验 `.phont` 与空文件 |
| `cli/schema/voice.go` | phont 字段 `SetPhontPath` → `parse.ReadPhont` + `api.SetPhontData` |
| `cli/action/voice.go` | 同上；空值清除 phont |
| `cli/ui/voice_param.go` | `pickPhont` 用 `parse.ReadPhont` + `api.SetPhontData` |
| `cli/synth/root.go` | `-P` 改为 `parseVoicePathSpec` + 读 bytes；写 `VoiceFile` / `VoiceName` / `VoiceBase` |
| `cli/synth/voice_file.go` | 新增 `-P` 解析、读文件、base name |
| `cli/test/action/action_test.go` | `cliopen.Open` + `Options` 路径覆盖，避免写用户目录 |
| `cli/go.mod` | 增加 `yukumo-script/core/port` 的 require + replace（传递依赖） |

`-P` 新流程：

```go
yamlPath, voiceName, err := parseVoicePathSpec(spec)
data, name, base, err := loadVoiceFile(yamlPath, voiceName)
req.VoiceFile = data
req.VoiceName = name
req.VoiceBase = base
```

`readPhont` 落在 `cli/parse`（CLI 层保留路径校验）：

```go
func ReadPhont(path string) ([]byte, error) {
	path = strings.TrimSpace(path)
	if path == "" { return nil, fmt.Errorf("phont path is required") }
	if !strings.EqualFold(filepath.Ext(path), ".phont") {
		return nil, fmt.Errorf("phont path must end with .phont")
	}
	data, err := os.ReadFile(path)
	if err != nil { return nil, err }
	if len(data) == 0 { return nil, fmt.Errorf("phont file is empty") }
	return data, nil
}
```

CLI 仍然不允许 import：

```text
yukumo-script/core/internal/...
yukumo-script/runtime/...
yukumo-script/core/port
```

CLI 允许 import：

```text
yukumo-script/core/api
yukumo-script/core/api/cliopen
github.com/spf13/cobra
github.com/charmbracelet/huh
github.com/charmbracelet/lipgloss
gopkg.in/yaml.v3
```

---

<a id="sec-9"></a>
## 9. synth 的 IO 改造

<a id="sec-9-1"></a>
### 9.1 合成流程

`core/internal/service/synth` 不再 import `os` / `path/filepath`：

```go
type Service struct {
	config    Config
	voices    Voices
	engine    port.Synthesizer
	artifacts port.ArtifactStore
	outputs   port.OutputWriter
}

func New(cfg Config, voices Voices, engine port.Synthesizer, artifacts port.ArtifactStore, outputs port.OutputWriter) *Service
```

`Synthesize` 流程：

1. 校验 text；
2. 解析 voice：
   - `VoiceRef` → `ParseVoiceRef` + `voices.Get`；
   - `VoiceFile` → `LoadVoiceBytes(data, name, base)`；
   - 两者都空或都非空 → `ErrInvalid`；
3. `SelectParam`；计算 speed；`TaskHash`；
4. 读 config；确定 reuse；
5. 解析输出目标（见 §9.2）；
6. reuse 且 `artifacts != nil`：`lookupCache(hash)`；
   - 命中且不写：返回 `{WAV, FromCache: true}`；
   - 命中且写：读缓存 bytes → 写到目标 → 返回 `{WAV, Path, Wrote: true, FromCache: true}`；
7. miss 且 `engine == nil` → `ErrInvalid`；否则 `engine.Synthesize`；
8. 不写：返回 `{WAV}`；
9. 写：写目标；`putCache(hash, target, wav)`；返回 `{WAV, Path, Wrote: true}`。

`-o` 与默认输出的分支：

```go
type target struct {
	external string // -o 的绝对路径；空表示 store key
	key      string // ArtifactStore key；空表示 external
	write    bool
}

func (s *Service) resolveTarget(cfg domain.Config, req Request, charName, voiceName string) (target, error) {
	if strings.TrimSpace(req.Output) != "" {
		return target{external: req.Output, write: true}, nil
	}
	if !cfg.Synth.WriteAlways() && req.Preview {
		return target{}, nil
	}
	if s.artifacts == nil {
		return target{}, fmt.Errorf("%w: artifact store unavailable", domain.ErrInvalid)
	}
	base := FormatOutputName(cfg.Synth.OutputFormat, charName, voiceName)
	key, err := s.uniqueKey(base, ".wav")
	return target{key: key, write: true}, err
}
```

写：

```go
func (s *Service) writeTarget(t target, data []byte) error {
	if t.external != "" {
		if s.outputs == nil {
			return fmt.Errorf("%w: output writer unavailable", domain.ErrInvalid)
		}
		return s.outputs.Write(t.external, data)
	}
	return s.artifacts.Write(t.key, data)
}

func (s *Service) targetPath(t target) string {
	if t.external != "" {
		return t.external
	}
	return s.artifacts.Path(t.key)
}
```

<a id="sec-9-2"></a>
### 9.2 缓存与输出

`index.yaml` / 缓存的 key 全部走 ArtifactStore：

```go
const indexKey = "index.yaml"

func (s *Service) loadIndex() (synthIndex, error)     // s.artifacts.Read(indexKey)
func (s *Service) saveIndex(idx synthIndex) error     // s.artifacts.Write(indexKey, yaml)
func (s *Service) lookupCache(hash string) (key string, ok bool, err error)
func (s *Service) putCache(hash string, t target, wav []byte) error
func (s *Service) CleanIndex() error
func (s *Service) uniqueKey(base, ext string) (string, error)
```

语义：

- 生成输出：`index[hash] = key`，key 是 `<base>.wav`；
- `-o`：`index[hash] = cacheKey(hash)`，`cacheKey` 形如 `cache-<hash>.wav`，同时把 WAV 副本写进 ArtifactStore；
- `lookupCache`：index 里找不到 / key 为空 → miss；key 在 store 里不存在 → 删除该 entry 并 miss；
- `CleanIndex`：遍历 index，`artifacts.Exists(key)` 为 false 的删除；
- `uniqueKey`：用 `artifacts.Exists` 循环 `base` / `base-1` / …，不再 `os.Stat`；
- 缓存命中且需要写时，读缓存 bytes 再写到目标，不再 `copyFile`。

`CleanIndex` 在 `api.App.CleanIndex` 暴露；`cliopen.Open` 启动时调用。

<a id="sec-9-3"></a>
### 9.3 错误与路径返回

- `SynthResult.Path`：
  - 默认输出 → `artifacts.Path(key)`；FS 是绝对路径，JS 是逻辑 key；
  - `-o` → `req.Output`；
  - 不写 → 空串。
- `SynthResult.Wrote` / `FromCache` 语义不变。
- `putCache` 失败仍然返回部分结果（`out, err`），与现状一致。
- `Outputs == nil` 且 `req.Output != ""` → `ErrInvalid`；`Artifacts == nil` 且需要默认输出 → `ErrInvalid`。

---

<a id="sec-10"></a>
## 10. 同步 Store + 内存缓存 + Flush

本计划不实现 IndexedDB，但 port 必须让下面的 web adapter 可实现：

**启动**

1. JS 侧用 IndexedDB 把 config / character / artifact 的 key-value 读入内存
   （或读入一个 `Map` 后交给 Go）。
2. 内存 store 实现 `port.Store` / `port.ArtifactStore`，读操作同步。
3. 读完之后再跑 `api.New` + `CleanIndex` + 业务调用。

**运行**

- Go 侧所有 `Read/Write/List/Exists/Delete` 都只碰内存。
- 每次 `Write` 后，JS adapter 可以：
  - fire-and-forget 写回 IndexedDB（最终一致）；或
  - 只标 dirty，等待 `Flush`。
- 选择 fire-and-forget 时，`Flush` 只负责把未完成的写队列收敛。

**关闭**

- `App.Close` → 每个 store 若实现 `port.Flusher` → `Flush(ctx)`。
- JS adapter 在 `Flush` 里等待所有 pending 的 IDB 事务。

**契约**

| 场景 | 结果 |
|---|---|
| 读 | 永远同步；数据已预载 |
| 写 | 立即对 Go 可见；对 IDB 最终一致 |
| 崩溃 | 最多丢最后一次未 flush 的写；可接受 |
| localStorage | 只用于小配置/元数据；WAV 走 IndexedDB 或 Blob |
| 多标签页 | 本计划不定义；后续可用 BroadcastChannel / 版本号 |

不要同时保留“同步 store + 异步 store”两套 service 语义；只保留同步 port，把异步藏在 adapter 里。

---

<a id="sec-11"></a>
## 11. Engine 的 JS 适配契约

浏览器/Node 不运行 `runtime/aq`，而是由 JS 直接实例化 `third-party/aq_wasm.wasm`。该模块 0 imports，导出一个 linear memory，JS 可以直接：

```js
const { instance } = await WebAssembly.instantiate(wasmBytes);
```

Go 侧实现 `port.Engine` 的 adapter（`wasmtest/jsengine.go`，未来可平移到 `web/`）：

```go
type JSEngine struct{ ref js.Value }

func (e JSEngine) Synthesize(ctx context.Context, text string, p port.EngineParam) ([]byte, error)
func (e JSEngine) HasPreset(ctx context.Context, version, preset int32) (bool, error)
func (e JSEngine) PresetNames(ctx context.Context, version int32) ([]string, error)
```

JS 侧暴露同步对象：

```js
globalThis.__yukumo = {
  engine: {
    synthesize(text, param, phont /* Uint8Array | null */) { /* → Uint8Array */ },
    presetNames(version) { /* → string[] */ },
    hasPreset(version, preset) { /* → boolean */ },
  },
  stores: {
    config: /* Store */,
    characters: /* Store */,
    artifacts: /* ArtifactStore */,
  },
};
```

`param` 对象字段与 `port.EngineParam` 一一对应：

```js
{ version, preset, speed, base, volume, pitch, accent, lmd, fsc }
```

### 11.1 ABI 映射

JS adapter 必须复刻 `runtime/aq` 的调用序列：

1. `text` 用 `TextEncoder` 编 UTF-8；
2. `malloc(len)` + `mem.set(text, ptr)`；
3. `phont` 非空时 `malloc` + copy；`phontPtr = 0` 表示无；
4. `malloc(40)` 写 40 字节 little-endian `AqSynthParam`：
   - `0..4` version
   - `4..8` preset
   - `8..12` speed
   - `12..16` phontPtr
   - `16..20` base
   - `20..24` volume
   - `24..28` pitch
   - `28..32` accent
   - `32..36` lmd
   - `36..40` fsc
5. `malloc(4)` 两个 slot：out_ptr / out_len，清零；
6. `aq_synthesize(textPtr, textLen, paramPtr, phontLen, outPtrSlot, outLenSlot)`；
7. 返回码非 0：
   - `-1/-2/-3` → `AQ_ERR_ARG` / `AQ_ERR_VERSION` / `AQ_ERR_OPEN`；
   - 其余 → `aq_error_message(version, code)` 读 C 字符串；
   - 抛 `Error`，Go adapter 转成 error；
8. 成功 → 读 out_ptr / out_len → `mem.slice(wavPtr, wavPtr+wavLen)` → `aq_free(wavPtr)`；
9. `free` 所有临时指针。

`aq_has_preset(version, preset)` 直接返回 `int32`；`aq_preset_names(version, outCount)` 返回指针表，先读 count，再按 `tablePtr + i*4` 读 C 字符串。

注意：

- wasm `memory.buffer` 会在 `malloc` 后增长；每次调用后重新创建 `Uint8Array` / `DataView`，不要缓存旧 buffer。
- Go `syscall/js` 调用是同步的；`aq_synthesize` 会阻塞 JS 事件循环。浏览器端后续应放到 Web Worker（§20）。
- `js.Value.Call` 在 JS 抛异常时会 panic；Go adapter 用 `defer recover` 转 error。
- `phont` 为空时传 `null`，不要传零长度 `Uint8Array`。
- 大 WAV 复制会产生两份内存；后续优化（Transferable ArrayBuffer / Worker）不在本计划。

---

<a id="sec-12"></a>
## 12. wasm 打包边界

### 12.1 允许被 wasm 引入的包

```text
yukumo-script/core/api          # 纯用例 + DTO + New
yukumo-script/core/port         # 契约
yukumo-script/core/internal/... # 通过 core/api 传递依赖（纯包）
gopkg.in/yaml.v3                # 通过 core/api 传递依赖
```

不允许：

```text
yukumo-script/core/api/cliopen
yukumo-script/core/internal/store/fs
yukumo-script/runtime
github.com/tetratelabs/wazero
yukumo-script/cli
项目包直接 import os / path/filepath
```

> `os` 仍会通过 `fmt`、`crypto/sha256`、`yaml` 等标准库/第三方包**传递**进入 wasm 构建，这是 Go 标准库自身的依赖，不是项目代码在调用文件 IO。门禁区分两件事：
>
> 1. 编译闭包里不能出现 `runtime` / wazero / `cli` / `cliopen` / `store/fs`；
> 2. **项目自己的包**不能直接 import `os` / `path/filepath`。
>
> `fmt` 在 js 下会 import `os`（标准库实现），这是允许的；我们的业务代码不调用 `os.Open/ReadFile/MkdirAll` 即可。

### 12.2 为什么拆 `cliopen` 就够

Go 只编译被 import 的包。`core/api` 根包不 import `cliopen`，所以 `import "yukumo-script/core/api"` 不会把 runtime 带进来；`cliopen` 自己 import runtime，只被 CLI import。

`core/api/go.mod` 仍然 require runtime，因为同一个 module 里有 `cliopen`；这只会影响 module 图，不影响 web 的编译闭包。真正的门禁是 `go list -deps`（§14.4）。

### 12.3 build tag 与平台文件

不依赖 `//go:build !js` 作为长期 seam。平台差异只通过包边界表达：

```text
core/api              纯
core/api/cliopen      OS + runtime
core/internal/store/fs OS
```

`wasmtest` 的所有 Go 文件都加 `//go:build js`，避免 native build / lint 尝试编译 `syscall/js`。

### 12.4 验证矩阵

| 目标 | 命令 | 期望 |
|---|---|---|
| native 全模块 | `just ci` | 全绿 |
| wasm 可编译 | `GOWORK=off GOOS=js GOARCH=wasm go -C core/api build ./...` | 成功 |
| wasm 依赖干净 | `go run scripts/check-wasm-deps.go` | 无 runtime / wazero / cli / fs 实现；项目包无 os/filepath 直接 import |
| wasm 运行 | `just wasm-smoke` | Node 退出码 0 |

`runtime` 本身在 `GOOS=js` 下也能编译（已验证），但它不进 `wasmtest` 的依赖闭包；不需要为 runtime 增加 wasm 闸门。

---

<a id="sec-13"></a>
## 13. 测试分层

### 13.1 迁移表

| 现在 | 之后 | 类型 |
|---|---|---|
| `core/api/test/runtime/aq/*` | 不变（仍 import runtime） | runtime 集成 |
| `core/api/test/synth/cache_test.go` | 用 `store.Memory` 作为 ArtifactStore；真实 WASM 合成仍由 `core/api/test/runtime/aq` 与 `wasmtest` 覆盖 | synth 集成 |
| `core/api/test/synth/service_test.go` | 用 `store.Memory` 作为 ArtifactStore + 真实 `aq.Pool` | synth 集成 |
| `core/api/test/synth/adapter_test.go` | adapter 用 `port.EngineParam` | synth 集成 |
| `core/test/service/synth/filename_test.go` | 删除 `UniquePath` 用例；`uniqueKey` 由 synth 集成测试覆盖 | 纯单元 |
| `core/test/service/synth/voice_test.go` | `LoadVoiceFile` → `LoadVoiceBytes` | 纯单元 |
| `core/test/service/config/service_test.go` | 去掉 EnsurePaths 用例；默认值用注入的 defaults | 纯单元 |
| `core/test/service/character/**` | `character.New(nil, mem)` 注入 Memory store | 纯单元 |
| `core/test/store/fs/fs_test.go` | 新增 FS Store / Output / key escape 测试 | OS 单元 |
| `cli/test/action/action_test.go` | 用 `cliopen.Open` + `Options` 路径覆盖 | CLI 单元 |
| `cli/test/prompt/**` | 不变 | CLI 单元 |
| `runtime/aq/errors_test.go` | 不变 | runtime 包内单元 |

### 13.2 新增测试

1. `core/api/app_test.go`
   - `api.New` 用 Memory store 构造，Config 默认值/往返、Character CRUD 与 `Close`；
   - 缺少 ConfigStore / CharacterStore 时返回错误。
2. `core/api/convert_test.go`
   - 现有用例保持。
3. `core/test/store/fs/fs_test.go`
   - FS Store 读写/List/Delete/key escape/Path + OutputWriter 写 `-o` 路径。
4. `wasmtest`
   - fake engine：CRUD + 配置持久化 + `Close` flush；
   - 真 `aq_wasm.wasm`：`HasPreset` / `Synthesize` / 缓存命中 / WAV header；
   - `go list -deps` 门禁。

### 13.3 不减少测试

- `core/api/test/runtime/aq` 继续真实跑 WASM（native）。
- `core/api/test/synth` 继续用真实 `aq.Pool` 跑至少一条合成。
- Node smoke 用真实 `aq_wasm.wasm` 跑一条合成。
- 删除的只有“把 OS 路径逻辑放在 core”的测试（EnsurePaths / UniquePath 的 FS 版），它们迁移到 `cliopen` / `fs` / CLI 测试。

---

<a id="sec-14"></a>
## 14. 静态检查与 depguard

### 14.1 module 与配置

```text
core/port/.golangci.yml
core/.golangci.yml
runtime/.golangci.yml
core/api/.golangci.yml
cli/.golangci.yml
wasmtest/.golangci.yml
```

每个配置显式启用：

```yaml
linters:
  enable:
    - depguard
    - staticcheck
```

`just lint` 覆盖前 5 个 native module；`just wasm-lint` 用 `GOOS=js GOARCH=wasm` 跑 `wasmtest`。

### 14.2 depguard 规则表

| module | 文件 | 禁止 |
|---|---|---|
| core/port | `**/*.go` | 所有 `yukumo-script/`、wazero、yaml、UI 库 |
| core | `internal/domain/**` | 所有 `yukumo-script/` |
| core | `internal/store/**` | `core/internal/service`、`core/api`、`cli`、runtime、wazero、UI 库 |
| core | `internal/store/fs/**` | 除 `core/port` 外的所有 `yukumo-script/` |
| core | `internal/service/**` | `os`、`path/filepath`、runtime、`core/api`、`cli`、wazero、UI 库 |
| core | `test/**` | runtime、`core/api`、`cli`、UI 库 |
| core | `**/*.go` | wazero、UI 库（兜底） |
| runtime | `**/*.go` | 所有 `yukumo-script/`、yaml、UI 库 |
| core/api | `*.go`（根包） | `yukumo-script/runtime`、wazero、`os`、`path/filepath`、`cli`、UI 库 |
| core/api | `test/**` | `cli`、UI 库、wazero（runtime 允许） |
| core/api | `cliopen/**` | `cli`、UI 库 |
| cli | `**/*.go` | `core/internal`、runtime、`core/port`、wazero |
| cli | `cmd/**` | `core`、`core/api`、`core/port`、runtime、UI 库、yaml |
| wasmtest | `**/*.go` | runtime、wazero、`cli`、`core/internal`、`core/api/cliopen`、UI 库、yaml |

外部库归属：

| 库 | 只允许 |
|---|---|
| `github.com/spf13/cobra` | `cli` |
| `github.com/charmbracelet/huh` | `cli` |
| `github.com/charmbracelet/lipgloss` | `cli` |
| `github.com/tetratelabs/wazero` | `runtime` |
| `gopkg.in/yaml.v3` | `core`、`core/api`、`cli` |

### 14.3 负向验证

1. 在 `core/internal/service/synth` 临时加：

   ```go
   import "os"
   ```

   期望 `just lint` 报 depguard violation；删除。

2. 在 `core/api` 根包临时加：

   ```go
   import _ "yukumo-script/runtime/aq"
   ```

   期望 `just lint` 报 depguard violation；删除。

3. 在 `wasmtest` 临时加：

   ```go
   import _ "yukumo-script/runtime/aq"
   ```

   期望 `just wasm-lint` 报 depguard violation；删除。

4. 在 `core/api/cliopen` 临时 import `yukumo-script/cli`，期望 lint 报错；删除。

### 14.4 `GOOS=js` 依赖门禁

`scripts/check-wasm-deps.go` 在 `wasmtest` 目录下执行：

```bash
GOWORK=off GOOS=js GOARCH=wasm go list -deps -f '{{.ImportPath}}' ./...
```

断言全量依赖不出现：

```text
yukumo-script/runtime
github.com/tetratelabs/wazero
yukumo-script/cli
yukumo-script/core/api/cliopen
yukumo-script/core/internal/store/fs
```

再对 `ImportPath` 以 `yukumo-script/` 开头的包检查 **直接** `Imports`，断言不出现：

```text
os
path/filepath
```

`os` 会通过 `fmt` / `yaml` 等传递进入闭包，这是允许的；门禁只禁止项目包直接 import。这条比 depguard 更硬：它看的是实际编译闭包，而不是单个源文件的 import。

---

<a id="sec-15"></a>
## 15. justfile 与 CI

### 15.1 justfile

新增/调整：

```just
# 原生 module 构建
build-modules:
    go -C core/port build ./...
    go -C core build ./...
    go -C runtime build ./...
    go -C core/api build ./...
    go -C cli build ./...

test:
    go -C core/port test ./...
    go -C core test ./...
    go -C runtime test ./...
    go -C core/api test ./...
    go -C cli test ./...

vet:
    go -C core/port vet ./...
    go -C core vet ./...
    go -C runtime vet ./...
    go -C core/api vet ./...
    go -C cli vet ./...

lint:
    cd core/port && golangci-lint run
    cd core && golangci-lint run
    cd runtime && golangci-lint run
    cd core/api && golangci-lint run
    cd cli && golangci-lint run

# 生成 wasm_exec.js，构建 js/wasm，跑依赖门禁与 Node 冒烟
wasm-prepare:
    go run scripts/prepare-wasm.go

wasm-smoke: wasm-prepare
    GOOS=js GOARCH=wasm go -C wasmtest build -o .build/smoke.wasm ./...
    go run scripts/check-wasm-deps.go
    node wasmtest/harness.cjs

wasm-lint:
    cd wasmtest && GOOS=js GOARCH=wasm golangci-lint run
```

`ci` 保持 native：

```just
ci: check-replaces build build-modules test vet
```

`build` 仍然依赖 `lint`，所以 native CI 的 lint 顺序不变。

### 15.2 `.gitignore`

新增：

```gitignore
/wasmtest/.build/
```

不要 ignore `third-party/aq_wasm.wasm`。

### 15.3 CI

`.github/workflows/ci.yml` 增加独立 job：

```yaml
  wasm:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.27"

      - uses: actions/setup-node@v4
        with:
          node-version: "22"

      - uses: extractions/setup-just@v2

      - name: Install golangci-lint
        run: |
          go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
          echo "$(go env GOPATH)/bin" >> "$GITHUB_PATH"

      - name: wasm smoke
        run: just wasm-smoke

      - name: wasm lint
        run: just wasm-lint
```

`wasm-smoke` 依赖 `third-party/aq_wasm.wasm`（checkout 自带）和 `third-party/aq_wasm.h`（不参与运行）。

### 15.4 `go.work` / replace / `check-replaces`

`go.work`：

```go
go 1.27

use (
	./core
	./core/port
	./core/api
	./runtime
	./cli
	./wasmtest
)
```

各 `go.mod` 的本地 replace：

```text
core/go.mod          → yukumo-script/core/port => ./port
core/api/go.mod      → yukumo-script/core => .., yukumo-script/core/port => ../port, yukumo-script/runtime => ../../runtime
cli/go.mod           → yukumo-script/core => ../core, yukumo-script/core/port => ../core/port,
                       yukumo-script/core/api => ../core/api, yukumo-script/runtime => ../runtime
wasmtest/go.mod      → yukumo-script/core => ../core, yukumo-script/core/port => ../core/port,
                       yukumo-script/core/api => ../core/api, yukumo-script/runtime => ../runtime
```

`just check-replaces` 继续跑 `scripts/check-replaces.go`，期望输出包含：

```text
yukumo-script/core       -> core
yukumo-script/core/port  -> core/port
yukumo-script/runtime    -> runtime
yukumo-script/core/api   -> core/api
```

---

<a id="sec-16"></a>
## 16. Node headless 冒烟

### 16.1 结构

```text
wasmtest/
├── go.mod
├── .golangci.yml
├── main.go          # //go:build js
├── smoke.go         # 场景
├── jsstore.go       # port.Store ↔ JS
├── jsartifact.go    # port.ArtifactStore ↔ JS
├── jsengine.go      # port.Engine ↔ JS
├── fakeengine.go    # Go fake
├── store.cjs        # JS MemoryStore / MemoryArtifactStore
├── aq_engine.cjs    # JS AqEngine（真 wasm）
├── harness.cjs      # 入口
└── .build/          # wasm_exec.js + smoke.wasm（gitignore）
```

### 16.2 Go 侧

`main.go`（以及所有 `wasmtest` Go 文件）加 `//go:build js`，避免 native build / lint 编译 `syscall/js`：

```go
//go:build js

func main() {
	if err := runSmoke(); err != nil {
		js.Global().Get("console").Call("error", "wasm smoke failed: "+err.Error())
		js.Global().Get("process").Set("exitCode", 1)
		return
	}
	js.Global().Get("process").Set("exitCode", 0)
}
```

`jsstore.go`：

```go
type JSStore struct{ ref js.Value }

func (s JSStore) Read(key string) ([]byte, error) {
	v := s.ref.Call("read", key)
	if v.IsNull() || v.IsUndefined() {
		return nil, port.ErrNotFound
	}
	out := make([]byte, v.Get("length").Int())
	js.CopyBytesToGo(out, v)
	return out, nil
}

func (s JSStore) Write(key string, data []byte) error
func (s JSStore) Exists(key string) (bool, error)
func (s JSStore) List(prefix string) ([]string, error)
func (s JSStore) Delete(key string) error
func (s JSStore) Flush(ctx context.Context) error // s.ref.Call("flush")
```

`jsartifact.go` 在 `JSStore` 上加 `Path(key) string`（JS 返回 `mem://<key>` 或空串）。

`jsengine.go`：

```go
type JSEngine struct{ ref js.Value }

func (e JSEngine) Synthesize(ctx context.Context, text string, p port.EngineParam) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("js engine: %v", r)
		}
	}()
	param := js.Global().Get("Object").New()
	param.Set("version", p.Version)
	// ... 其余字段
	var phont js.Value
	if len(p.Phont) > 0 {
		phont = js.Global().Get("Uint8Array").New(len(p.Phont))
		js.CopyBytesToJS(phont, p.Phont)
	} else {
		phont = js.Null()
	}
	res := e.ref.Call("synthesize", text, param, phont)
	out = make([]byte, res.Get("length").Int())
	js.CopyBytesToGo(out, res)
	return out, nil
}
```

`HasPreset` / `PresetNames` 同理，`ref.Call("hasPreset", version, preset)` / `ref.Call("presetNames", version)`。

`fakeengine.go`：

```go
type FakeEngine struct{}

func (FakeEngine) Synthesize(ctx context.Context, text string, p port.EngineParam) ([]byte, error)
func (FakeEngine) HasPreset(ctx context.Context, version, preset int32) (bool, error) // true
func (FakeEngine) PresetNames(ctx context.Context, version int32) ([]string, error)  // []string{"fake"}
```

### 16.3 JS 侧

`store.cjs`：

```js
class MemoryStore {
  constructor(prefix) { this.prefix = prefix; this.data = new Map(); }
  read(key) { return this.data.has(key) ? this.data.get(key) : null; }
  write(key, bytes) { this.data.set(key, new Uint8Array(bytes)); }
  exists(key) { return this.data.has(key); }
  list(prefix) { /* immediate keys, sorted */ }
  delete(key) { if (!this.data.has(key)) throw new Error('store: not found'); this.data.delete(key); }
  path(key) { return `mem://${key}`; }
  flush() {}
}
```

`aq_engine.cjs`：按 §11.1 实现 `synthesize` / `presetNames` / `hasPreset`。

`harness.cjs`：

```js
const fs = require('node:fs');
const path = require('node:path');

// 对齐 GOROOT/lib/wasm/wasm_exec_node.js 的宿主环境
globalThis.require = require;
globalThis.fs = fs;
globalThis.path = path;
globalThis.TextEncoder ??= require('node:util').TextEncoder;
globalThis.TextDecoder ??= require('node:util').TextDecoder;
globalThis.performance ??= require('node:perf_hooks').performance;
globalThis.crypto ??= require('node:crypto').webcrypto;

require('./.build/wasm_exec.js'); // 设置 globalThis.Go

async function main() {
  const wasmBytes = fs.readFileSync(path.join(__dirname, '..', 'third-party', 'aq_wasm.wasm'));
  const { instance } = await WebAssembly.instantiate(wasmBytes);
  const engine = makeAqEngine(instance);

  const config = new MemoryStore('config');
  const characters = new MemoryStore('characters');
  const artifacts = new MemoryStore('artifacts');

  globalThis.__yukumo = {
    engine,
    stores: { config, characters, artifacts },
  };

  const go = new globalThis.Go();
  go.argv = [];
  go.env = { ...process.env };
  const program = fs.readFileSync(path.join(__dirname, '.build', 'smoke.wasm'));
  const { instance: goInstance } = await WebAssembly.instantiate(program, go.importObject);
  await go.run(goInstance);
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});
```

`wasm_exec.js` 从 GOROOT 复制，不提交：

```go
// scripts/prepare-wasm.go
// 1. go env GOROOT
// 2. 复制 <GOROOT>/lib/wasm/wasm_exec.js → wasmtest/.build/wasm_exec.js
// 3. 创建 wasmtest/.build
```

### 16.4 冒烟场景

`smoke.go` 按顺序：

1. **fake engine + 内存 store**
   - `api.New(Deps{ConfigStore: config, CharacterStore: characters, Artifacts: artifacts, Engine: FakeEngine{}, Defaults: api.DefaultConfig()})`；
   - `app.SaveConfig(cfg)` / `app.Config()` 往返；
   - `app.AddCharacter` / `ListCharacters` / `GetCharacter`；
   - `app.AddVoice` / `ListVoices` / `DuplicateVoice` / `DeleteVoice`；
   - `app.Close`，确认 `Flush` 被调用（JS 侧计数）。
2. **真 aq_wasm engine**
   - 用 `JSEngine` 替换 `FakeEngine`；
   - `AddCharacter` + `AddVoice`，version 1 / preset 0 / speed 100；
   - `app.Synthesize(ctx, SynthRequest{ VoiceRef: "smoke:default", Text: "こんにちは" })`；
   - 断言：`WAV` 长度 > 44，前 4 字节是 `RIFF`，`Wrote == true`，`Path != ""`；
   - 再次同参数合成：`FromCache == true`，bytes 与第一次相等；
   - `app.CleanIndex()` 不报错；
   - `app.Close`。
3. **失败路径**
   - `Deps.Engine == nil` 时 `Synthesize` 返回错误而不是 panic；
   - `Deps.Artifacts == nil` 且需要默认输出时返回错误；
   - `Deps.Outputs == nil` 且 `Output != ""` 时返回错误。

### 16.5 成功标准

- `node wasmtest/harness.cjs` 退出码 0；
- 输出包含 `wasm smoke: ok`；
- `just wasm-smoke` 同时跑 `check-wasm-deps`，失败即退出非 0。

---

<a id="sec-17"></a>
## 17. 迁移步骤

> 结构性搬迁建议在独立 commit 完成；中间 commit 不要求可构建，但每个 Phase 结束时 `just ci` 应绿。CLI 改动集中在 Phase 4。

### Phase 0：基线

1. 记录当前 `wasm-port` commit（与 `layer-refactor` 同点）。
2. 跑 `just ci`，确认全绿。
3. 手动跑一次 CLI smoke，记录输出：

   ```bash
   YUKUMO_CONFIG_DIR=$(mktemp -d) just run config list
   YUKUMO_CONFIG_DIR=$(mktemp -d) just run character list
   YUKUMO_CONFIG_DIR=$(mktemp -d) just run synth -v remilia:default "こんにちは"
   ```

4. 记录当前 `output_path` 下的 `index.yaml` / WAV 布局。

### Phase 1：`core/port` + store 拆分

1. 创建 `core/port/go.mod`、`store.go`、`artifact.go`、`engine.go`。
2. `core/internal/store` 改成纯包：`Memory`、YAML helper；删除 `fs.go`。
3. 新建 `core/internal/store/fs`：`Store` / `Output`；实现 `port.ArtifactStore` / `port.OutputWriter`。
4. `core/go.mod` 加 `core/port` require + replace。
5. `core/internal/service/{config,character,synth}` 改成依赖 `port`：
   - `config.New(store port.Store, defaults domain.Config)`；
   - `character.New(cfg, store port.Store)`；
   - `synth.New(cfg, voices, engine port.Synthesizer, artifacts port.ArtifactStore, outputs port.OutputWriter)`。
6. 此时 `core/api` 还会红；继续 Phase 2。

### Phase 2：`core/api` 纯化 + `cliopen`

1. `core/api/app.go` 改成 `Deps` / `New` / `App` / `Close` / `CleanIndex`。
2. 删除 `core/api` 根包的 `Options` / `Open` / `resolveConfigDir` / `aqengine.go`。
3. 新建 `core/api/cliopen/`：
   - `open.go` / `defaults.go` / `aqengine.go`；
   - `Open` 按 §7 流程；
   - `EnsurePaths`。
4. `core/api/go.mod` 加 `core/port` require + replace；runtime 继续保留（cliopen）。
5. `core/api` 根包检查：`grep -R '"os"\|"path/filepath"\|yukumo-script/runtime' core/api/*.go` 应为空。

### Phase 3：synth IO / bytes-only

1. `core/internal/service/synth` 去掉 `os` / `filepath`：
   - `index.go` 改用 ArtifactStore；
   - `service.go` 用 `target` / `writeTarget` / `targetPath`；
   - `filename.go` 的 `UniquePath` → `uniqueKey`；
   - `voice.go` 的 `LoadVoiceFile` → `LoadVoiceBytes`。
2. `domain/config.go` 去掉 `os` / `filepath`；`DefaultConfig(paths PathConfig)`。
3. `config.Service` 去掉 `EnsurePaths` 与 `ensurePaths`；默认值从构造参数来。
4. `voice.LoadPhont` / `SetPhontPath` 删除；phont 校验移到 CLI。
5. `api`：
   - `SynthRequest` 改 bytes；
   - `SetPhontData` 替代 `SetPhontPath`；
   - `DefaultConfig` 加；
   - `SynthResult` 保持不变。
6. 迁移测试（§13.1）。

### Phase 4：CLI

1. `cli/root.go` 改用 `cliopen.Open`。
2. `cli/action/config.go` 保存后调用 `cliopen.EnsurePaths`。
3. `cli/parse/parse.go` 新增 `ReadPhont`；`cli/schema/voice.go`、`cli/action/voice.go`、`cli/ui/voice_param.go` 改用 `SetPhontData`。
4. `cli/synth/root.go` + 新 `voice_file.go` 改用 bytes。
5. `cli/test/action/action_test.go` 改用 `cliopen.Open` + 路径覆盖。
6. `cli/go.mod` 加 `core/port` require + replace；`go mod tidy`。
7. `just ci` 必须绿。

### Phase 5：wasmtest + Node

1. 创建 `wasmtest` module 与 replace。
2. 写 Go adapter（`jsstore.go` 含 JSArtifactStore、`jsengine.go`、`fakeengine.go`）与 `smoke.go`。
3. 写 JS（`store.cjs` / `aq_engine.cjs` / `harness.cjs`）。
4. 写 `scripts/prepare-wasm.go` / `scripts/check-wasm-deps.go`。
5. `just wasm-smoke` 手动绿。

### Phase 6：lint / justfile / CI / docs

1. 6 份 `.golangci.yml`；`just lint` / `just wasm-lint`。
2. `justfile` 增加 `wasm-prepare` / `wasm-smoke` / `wasm-lint`；`build-modules` / `test` / `vet` / `lint` 加 `core/port`。
3. `.github/workflows/ci.yml` 增加 `wasm` job。
4. `.gitignore` 增加 `/wasmtest/.build/`。
5. README / README.zh 增加 wasm 段落（可选，但建议）。
6. `1-layer-refactor.md` §16 改成指针。
7. 跑 §18 全量验收。

### 17.1 实际落地差异

- 实际提交顺序是 Phase 1 → synth IO → API 拆分 → CLI → wasmtest → CI；与计划 Phase 2/3 的顺序对调，提交记录见头部。
- CLI 已合入 `2-cli-layer-refactor` 的结构，phont 读取落在 `cli/parse.ReadPhont`，`EnsurePaths` 落在 `cli/action/config.go`，不再有 `cli/character/voice/phont.go`。
- `cliopen.Options` 增加 `CharacterPath` / `OutputPath` 覆盖，供 CLI 测试避免写用户目录。
- `core/internal/store/fs` 的 Store / ArtifactStore / Output 合并在 `fs.go`，没有单独的 `artifact.go` / `output.go`。
- `core/api/test/synth` 用 `store.Memory` 作为 ArtifactStore；真实 WASM 合成仍由 `core/api/test/runtime/aq` 和 `wasmtest` 覆盖。
- `UniquePath` 的直接单测删除；`uniqueKey` 由 synth 集成测试覆盖。
- README / README.zh 的 wasm 段落未加（可选），不影响 gate。

### 17.2 完成状态

| Phase | 状态 | 提交 |
|---|---|---|
| Phase 0 基线 | 完成 | `88e62cd` |
| Phase 1 `core/port` + store | 完成 | `51addf9` |
| Phase 2 `core/api` + `cliopen` | 完成 | `311ea93` |
| Phase 3 synth IO | 完成 | `c2fbaed` |
| Phase 4 CLI | 完成 | `e5a085f` |
| Phase 5 `wasmtest` | 完成 | `6d938e1` |
| Phase 6 lint / CI / docs | 完成 | `8ff61ef` |

---

<a id="sec-18"></a>
## 18. 验收清单

### 18.1 native

```bash
just ci
just build
just run
```

- `just ci` 全绿；
- 4 个原有 native module + `core/port` 的 build/test/vet/lint 全绿；
- `just build` 仍生成根目录 `yukumo-script`；
- CLI help / 交互不变。

### 18.2 core 纯度

```bash
grep -R '"os"\|"path/filepath"' core/internal/service core/internal/domain core/api/*.go   # 应为空
grep -R 'yukumo-script/runtime\|tetratelabs/wazero' core/internal core/api/*.go             # 应为空
```

### 18.3 wasm 构建 / 依赖

```bash
GOOS=js GOARCH=wasm GOWORK=off go -C core/api build ./...
GOOS=js GOARCH=wasm GOWORK=off go -C wasmtest build ./...
go run scripts/check-wasm-deps.go
```

- 两条 build 成功；
- 依赖门禁无 `runtime` / `wazero` / `cli` / `core/internal/store/fs`；项目包无 `os` / `path/filepath` 直接 import。

### 18.4 Node 冒烟

```bash
just wasm-smoke
just wasm-lint
```

- fake engine 的 CRUD / Config / Flush 全通过；
- 真 `aq_wasm.wasm` 合成出 RIFF WAV；
- 第二次合成 `FromCache == true`；
- 退出码 0。

### 18.5 CLI 行为

```bash
YUKUMO_CONFIG_DIR=$(mktemp -d) just run config list
YUKUMO_CONFIG_DIR=$(mktemp -d) just run config edit paths.output_path /tmp/yukumo-out
YUKUMO_CONFIG_DIR=$(mktemp -d) just run character add remilia "example"
YUKUMO_CONFIG_DIR=$(mktemp -d) just run character list
YUKUMO_CONFIG_DIR=$(mktemp -d) just run synth -v remilia:default "こんにちは"
YUKUMO_CONFIG_DIR=$(mktemp -d) just run synth -P /path/character.yaml:default "こんにちは"
```

- 输出文本与 Phase 0 记录一致；
- `-o` 写指定路径；
- `-p` 仍能播放；
- phont 选择仍工作；
- `config edit` 后输出目录 / 角色目录被创建；
- `output_path` 下仍是 `index.yaml` + WAV 布局。

### 18.6 缓存

- 同参数第二次合成 `(cached)`；
- `-o` 同参数第二次合成命中缓存（本次修正）；
- 删除缓存 WAV 后 `CleanIndex` 清掉 entry；
- 旧 `index.yaml`（basename 形式）仍可命中。

### 18.7 边界

```bash
grep -R 'yukumo-script/core/internal' cli --include='*.go'   # 应为空
grep -R 'yukumo-script/runtime' cli --include='*.go'         # 应为空
grep -R 'yukumo-script/core/port' cli --include='*.go'       # 应为空
```

### 18.8 replace

```bash
just check-replaces
```

输出包含 core / core/port / runtime / core/api 四个映射，无冲突。

---

<a id="sec-19"></a>
## 19. 风险与取舍

### 19.1 独立 port module

- 多一个 module 与 replace 项；`check-replaces` 必须覆盖。
- 收益是 web/ 的依赖图真正干净；`core/port` 只有 stdlib。
- 如果后续觉得 6 个 module 太重，可以合并回 `core`，但那时 web 的 module 图会重新带上 yaml。

### 19.2 `core/api` 仍 require runtime

- `core/api/go.mod` 为了 `cliopen` 保留 runtime；`core/api` 根包的编译闭包不含它。
- 误判风险：以后有人在 `core/api` 根包 import `cliopen`，web 就会把 runtime 带进来；depguard + `check-wasm-deps` 双重拦截。

### 19.3 行为差异

1. **`-o` 缓存**：从“把 `-o` basename 写进 index（通常命中不了）”改成独立的 `cache-<hash>.wav`，缓存可命中。
2. **`EnsurePaths` 错误**：从 panic 改成返回 error；CLI 会打印错误而不是崩溃。
3. **`DefaultCharacterPath` 错误**：`os.UserConfigDir` 失败时返回 error，不再 panic。
4. **`-P` 文件读取位置**：从 service 移到 CLI；错误文本可能略有差异，但输出/行为不变。
5. **`SynthRequest.VoicePathSpec` 删除**：CLI 内部 API 变化；没有外部用户。
6. **`core/api.Open` 删除**：CLI 改用 `cliopen.Open`；没有外部用户。

除以上，CLI 输出文本、交互流程、合成策略、YAML 字段不变。

### 19.4 ArtifactStore 与 web

- 本计划只在 `wasmtest` 里实现 JS 内存版；产品级 IndexedDB 版留给 §20。
- `ArtifactStore.Path` 在 web 返回逻辑 key；任何“打印路径”的 UI 在 web 上要自己决定怎么展示。
- 缓存 key 布局对 web 是逻辑字符串，JS adapter 可以自由映射到 IDB object store / prefix。

### 19.5 JS engine

- `aq_synthesize` 同步执行，会阻塞 JS 事件循环；长文本/多请求需要 Web Worker（§20）。
- wasm `memory` 增长后旧 `DataView` 失效；adapter 必须每次重新取 buffer。
- `syscall/js` 的异常是 panic；adapter 必须 recover。
- 大 WAV 在 JS ↔ Go 之间复制；后续可以用 Worker + Transferable 优化。

### 19.6 同步 store 的崩溃窗口

- 采用内存缓存 + `Flush`，崩溃可能丢最后一次未 flush 的写。
- 如果以后要求“每次写立即持久”，需要改异步 Store（service 全链路 async），成本高；本计划明确不这样做。

### 19.7 wasm_exec.js 与 Node

- `wasm_exec.js` 跟随 Go 版本，不提交；`wasm-prepare` 从 `GOROOT/lib/wasm` 复制。
- CI 需要 Node；本计划把 Node 放到独立 job，`just ci` 不依赖 Node。
- Node 版本在 CI 固定；本地用任意 LTS 即可。

### 19.8 测试与真实 WASM

- native 的 `core/api/test/synth` 仍用真实 `aq.Pool`；
- Node smoke 用真实 `aq_wasm.wasm`；
- 只有 `synth` 的纯单元测试用 fake；不会把真实合成测试变成 fake。

---

<a id="sec-20"></a>
## 20. 后续：真正的 Web 实现

本计划完成后，`web/` 可以按下面的顺序增量做，不需要再改 core：

1. **`web/` module**（`yukumo-script/web`）
   - Go 侧：`web/jsstore.go` / `web/jsartifact.go` / `web/jsengine.go`（从 `wasmtest` 平移）；
   - JS 侧：IndexedDB store、Web Audio、UI。
   - `go.work` 加 `./web`；depguard 允许 `core/api` + `core/port`，禁止 `core/internal` / `runtime` / `cli` / `cliopen`。
2. **IndexedDB ArtifactStore**
   - 实现 `port.ArtifactStore` + `port.Flusher`；
   - 预载 / 写回策略按 §10。
3. **Web Audio**
   - `SynthResult.WAV` 直接 `AudioContext.decodeAudioData`；
   - `-o` / `-p` 在 web 上不适用，UI 自己决定下载/播放。
4. **Web Worker**
   - 把 Go wasm + `aq_wasm.wasm` 放进 Worker，避免阻塞主线程；
   - `port.Store` / `port.Engine` 的 JS 侧改成 Worker 消息（此时 Go 侧 adapter 通过 `postMessage` 同步等待需要 SharedArrayBuffer；或者保持主线程同步，后续再评估）。
5. **WASI（可选）**
   - 如果以后要 wasmtime/Node CLI，需要另一套 `port.Store`（WASI FS）与 `port.OutputWriter`；
   - `core/api` / `core/port` 不需要改。
6. **`runtime.NewWithWASM`（可选）**
   - 只有决定在 wasm 内用 Go runtime + wazero 时才需要；当前选择 JS engine，不做。

---

<a id="sec-a"></a>
## 附录 A：文件映射

| 现在 | 之后 |
|---|---|
| `core/internal/store/store.go` | `core/port/store.go` + `core/internal/store/store.go`（type alias） |
| `core/internal/store/fs.go` | `core/internal/store/fs/fs.go`（Store + ArtifactStore + Output） |
| `core/internal/store/memory.go` | 改实现 port + ArtifactStore + Flusher |
| `core/internal/store/yaml.go` | 参数改 port |
| `core/internal/service/synth/index.go` | 改用 ArtifactStore |
| `core/internal/service/synth/filename.go` | 删除 `UniquePath`；`uniqueKey` 移到 `index.go` |
| `core/internal/service/synth/voice.go` | `LoadVoiceFile` → `LoadVoiceBytes` |
| `core/internal/domain/config.go` | 删除 `DefaultCharacterPath` / `DefaultOutputPath` / `os` |
| `core/api/app.go` | `Deps` / `New` / `App` / `Close` / `CleanIndex` |
| `core/api/aqengine.go` | `core/api/cliopen/aqengine.go` |
| （无） | `core/api/cliopen/open.go`（`Options` + `EnsurePaths` + `Open`） |
| （无） | `core/api/cliopen/defaults.go`（OS 默认路径 + `apiDefaults`） |
| `core/api/model.go` | `SetPhontPath` → `SetPhontData`；加 `DefaultConfig` |
| `core/api/synth.go` | `SynthRequest` 改 bytes |
| `core/internal/service/character/voice/phont.go` | 删除；CLI 用 `cli/parse.ReadPhont` + `api.SetPhontData` |
| （无） | `cli/synth/voice_file.go` |
| （无） | `core/test/store/fs/fs_test.go` |
| （无） | `core/api/app_test.go` |
| （无） | `wasmtest/**` |
| （无） | `scripts/prepare-wasm.go` |
| （无） | `scripts/check-wasm-deps.go` |
| `docs/refactorings/1-layer-refactor.md` §16 | 指向本文件 |
| （无） | `docs/refactorings/3-wasm-port.md` |

---

<a id="sec-b"></a>
## 附录 B：命令速查

```bash
# native
just ci
just build
just run
just lint
just test
just vet
just check-replaces

# wasm 构建 / 依赖
GOOS=js GOARCH=wasm GOWORK=off go -C core/api build ./...
GOOS=js GOARCH=wasm GOWORK=off go -C wasmtest build ./...
go run scripts/check-wasm-deps.go

# Node 冒烟
just wasm-prepare
just wasm-smoke
just wasm-lint

# 手动模块检查
(cd core/port && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd core && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd runtime && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd core/api && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd cli && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)

# wasm 依赖闭包
(cd wasmtest && GOWORK=off GOOS=js GOARCH=wasm go list -deps -f '{{.ImportPath}}' ./...)

# 负向检查（临时加 import 后跑，再删除）
(cd core && GOWORK=off golangci-lint run)
(cd core/api && GOWORK=off golangci-lint run)
(cd wasmtest && GOWORK=off GOOS=js GOARCH=wasm golangci-lint run)
```
