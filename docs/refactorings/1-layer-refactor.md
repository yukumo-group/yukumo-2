# Yukumo Script 分层重构计划

- 状态：Implemented
- 日期：2026-11-16
- 分支：`layer-refactor`
- 实施记录：
  - `3e95c63` docs(refactor): add layered architecture plan
  - `abb0d21` refactor(core): move domain, store, and service under core
  - `58bc5a5` refactor(runtime): extract runtime and core modules
  - `2bfd99a` refactor(api): add use-case facade and engine port
  - `680ea1a` fix(api): preserve YAML field names in DTOs
  - `e4b82c6` refactor(cli): move CLI into its own module behind api
  - `d35d07e` chore(lint): add justfile, depguard/staticcheck, and CI gates
  - `f59b997` docs: update architecture guides and refactor status
  - `69a3b4f` chore(just): build every module in ci
  - `0cedf9a` fix(ci): track the CLI entrypoint
- 最终检查：`just ci` 通过（replace 检查、4 module build/test/vet、depguard + staticcheck 0 issues）。
- 范围：模块边界、依赖方向、API 门面、测试分层、CI 护栏；不改 CLI 输出、交互流程、合成策略、文件格式。
- Web/WASM：可移植性评估与后续增量设计见 [16. Web/WASM 可移植性 addendum](#sec-16)，不属于当前 Phase 0–6 的实现范围。

## 目录

- [1. 最终决策](#sec-1)
- [2. 目标目录](#sec-2)
- [3. 模块与依赖](#sec-3)
- [4. 模块文件](#sec-4)
- [5. 编译期边界与 linter](#sec-5)
- [6. API 层设计](#sec-6)
- [7. Engine port 解耦](#sec-7)
- [8. CLI 改造](#sec-8)
- [9. 测试分层](#sec-9)
- [10. 静态检查：depguard + staticcheck](#sec-10)
- [11. replace 一致性脚本](#sec-11)
- [12. CI](#sec-12)
- [13. 迁移步骤](#sec-13)
- [14. 验收清单](#sec-14)
- [15. 风险与取舍](#sec-15)
- [16. Web/WASM 可移植性 addendum](#sec-16)
  - [16.1 结论](#sec-16-1)
  - [16.2 已经可复用的部分](#sec-16-2)
  - [16.3 缺口](#sec-16-3)
  - [16.4 必需的架构增量](#sec-16-4)
  - [16.5 Store 契约](#sec-16-5)
  - [16.6 同步还是异步](#sec-16-6)
  - [16.7 WASM 打包边界](#sec-16-7)
  - [16.8 优先级](#sec-16-8)
  - [16.9 待定问题](#sec-16-9)
  - [16.10 对现有 plan 的改动点](#sec-16-10)
- [附录 A：文件映射](#sec-a)
- [附录 B：命令速查](#sec-b)

---

<a id="sec-1"></a>
## 1. 最终决策

| 决策 | 结果 |
|---|---|
| 根目录 | **没有根 `go.mod`、没有根 main 包**；只保留 `go.work`、`justfile`、`scripts/`、`docs/`、`.github/`、`third-party/` 等聚合文件 |
| 入口 | `main` 移入 CLI module：`cli/cmd/yukumo-script/main.go` |
| 模块 | `core`、`runtime`、`core/api`、`cli`，共 4 个 module |
| 连接方式 | 根 `go.work` 连接 4 个 module |
| runtime 位置 | 从 core 移出，顶层 `runtime/`，module path `yukumo-script/runtime` |
| core 依赖 | 生产代码和测试都不依赖 runtime/wazero；`core/go.mod` 只剩 yaml |
| 集成测试 | runtime 集成测试 + synth 集成测试迁到 API module（API 同时依赖 core 和 runtime） |
| API 位置 | `core/api`（必须在 core 路径下，才能 import `core/internal`） |
| 单向依赖 | Go internal 规则负责 compiler 硬边界；depguard 负责 layer 方向与外部依赖隔离 |
| 静态检查 | depguard（依赖方向）+ staticcheck（通用缺陷），统一由 `just lint` 执行 |
| 命令入口 | 根 `justfile`：`build` 依赖 `lint`；`ci` = replace 检查 + build + test + vet |
| CI | 调用 `just ci`；每模块 `GOWORK=off` 单独 build/vet/test；另一个脚本核对所有 replace 目标一致 |

目标依赖：

```text
cli/cmd/yukumo-script (main)
  ↓
cli  ── cobra / huh / lipgloss
  ↓
core/api  ── 用例门面 + DTO + 组合根 + aq adapter
  ↓                    ↘
core/internal/service   runtime/aq ── wazero + Yukw aq_wasm.wasm
  ↓
core/internal/{domain, store}
```

`core` 的 go.mod 只有：

```text
yukumo-script/core        → gopkg.in/yaml.v3
```

`runtime` 的 go.mod 只有：

```text
yukumo-script/runtime     → github.com/tetratelabs/wazero
```

“core 不依赖 wazero” 在 module 级别成立，不再只是“代码没直接 import”。

---

<a id="sec-2"></a>
## 2. 目标目录

```text
.
├── go.work
├── go.work.sum
├── README.md
├── README.zh.md
├── docs/
│   └── refactor-plan.md
├── justfile                            # 跨平台命令集线器
├── scripts/
│   └── check-replaces.go
├── .github/
│   └── workflows/
│       └── ci.yml
├── core/                               # module yukumo-script/core
│   ├── go.mod
│   ├── go.sum
│   ├── .golangci.yml
│   ├── internal/
│   │   ├── domain/
│   │   ├── store/
│   │   └── service/
│   │       ├── config/
│   │       ├── character/
│   │       │   └── voice/
│   │       ├── nameutil/
│   │       └── synth/
│   ├── api/                            # module yukumo-script/core/api
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── .golangci.yml
│   │   ├── app.go
│   │   ├── model.go
│   │   ├── convert.go
│   │   ├── aqengine.go
│   │   ├── config.go
│   │   ├── character.go
│   │   ├── voice.go
│   │   ├── synth.go
│   │   ├── errors.go
│   │   └── test/
│   │       ├── runtime/aq/             # 原 test/runtime/aq
│   │       └── synth/                  # 原 test/service/synth 的集成测试
│   └── test/
│       └── service/                    # core 的纯单元测试（无 runtime）
│           ├── character/
│           ├── config/
│           ├── nameutil/
│           └── synth/                  # 只留 filename/voice 等纯函数测试
├── runtime/                            # module yukumo-script/runtime
│   ├── go.mod
│   ├── go.sum
│   ├── .golangci.yml
│   └── aq/
│       ├── load.go
│       ├── param.go
│       ├── pool.go
│       ├── runtime.go
│       ├── synthesize.go
│       └── errors_test.go
├── cli/                                # module yukumo-script/cli
│   ├── go.mod
│   ├── go.sum
│   ├── .golangci.yml
│   ├── cmd/
│   │   └── yukumo-script/
│   │       └── main.go
│   ├── root.go
│   ├── browse.go
│   ├── prompt/
│   ├── setting/
│   ├── config/
│   ├── character/
│   ├── synth/
│   └── test/prompt/                    # 原 test/cli/prompt
└── third-party/                        # 不动
    ├── aq_wasm.h
    └── aq_wasm.wasm
```

说明：

- 根目录不再有 `main.go`，所以不需要根 `go.mod`。`go.work` 已经足够聚合 4 个 module。
- `main` 放在 CLI module 内，避免额外创建 `cmd` module；CLI module 本身就是交互层，入口放这里合理。
- `core/api` 必须留在 `core` 路径下。Go 的 `internal` 规则按 import path 判断：`yukumo-script/core/api` 位于 `yukumo-script/core` 下，才能 import `yukumo-script/core/internal/...`；`yukumo-script/cli` 不行。
- runtime 顶层化后是 public package，CLI 在编译器层面可以 import，所以由 depguard 明确禁止。
- `third-party/aq_wasm.wasm` 保持根目录。runtime 的 `loadWASM` 继续向上查找，CI 在 `core/api` 下跑集成测试也能找到。

### 2.1 `.gitignore` 必须同步调整

当前 `.gitignore` 有：

```gitignore
*.yaml
*.yml
```

这会忽略：

- `.github/workflows/ci.yml`
- 各 module 的 `.golangci.yml`

迁移时必须加白名单，例如：

```gitignore
!.github/
!.github/**
!.golangci.yml
!**/.golangci.yml
```

否则 CI / linter（depguard、staticcheck）配置不会进 Git，检查会“静默缺失”。

另外，二进制名 `yukumo-script` 必须锚定根目录：

```gitignore
/yukumo-script
```

如果写成 `yukumo-script`，Git 会同时忽略 `cli/cmd/yukumo-script/` 目录，导致 CI checkout 缺入口文件。

---

<a id="sec-3"></a>
## 3. 模块与依赖

| module | 目录 | 职责 | 直接依赖 |
|---|---|---|---|
| `yukumo-script/core` | `core/` | domain + store + service | yaml、stdlib |
| `yukumo-script/runtime` | `runtime/` | wazero、Yukw WASM、Pool、AQ 参数编码 | wazero、stdlib |
| `yukumo-script/core/api` | `core/api/` | DTO、用例门面、组合根、aq adapter | core、runtime、stdlib |
| `yukumo-script/cli` | `cli/` | cobra/huh/lipgloss、prompt、播放、展示、main | core/api、UI 库、stdlib |

依赖方向：

```text
cli (含 main) -> core/api -> core/internal/service -> core/internal/{domain,store}
                        \-> runtime/aq -> wazero
```

禁止方向：

```text
core        -> runtime / core/api / cli
runtime     -> 任何 yukumo-script 包
core/api    -> cli
cli         -> core/internal / runtime
cli/cmd     -> core、runtime（只能 import cli）
```

---

<a id="sec-4"></a>
## 4. 模块文件

> 版本号和 indirect 项以 `go mod tidy` 的最终结果为准；下面是结构示意。

### 4.1 `go.work`

```go
go 1.27

use (
	./core
	./runtime
	./core/api
	./cli
)
```

根目录不 `use .`，因为没有根 module。

### 4.2 `core/go.mod`

```go
module yukumo-script/core

go 1.27

require gopkg.in/yaml.v3 v3.0.1
```

不再有 runtime、wazero 或本地 replace。core 的测试也只 import core internal + yaml，不 import runtime。

### 4.3 `runtime/go.mod`

```go
module yukumo-script/runtime

go 1.27

require github.com/tetratelabs/wazero v1.12.0
```

package 路径是 `yukumo-script/runtime/aq`。

### 4.4 `core/api/go.mod`

```go
module yukumo-script/core/api

go 1.27

require (
	yukumo-script/core v0.0.0
	yukumo-script/runtime v0.0.0
)

replace (
	yukumo-script/core => ..
	yukumo-script/runtime => ../../runtime
)
```

`core/api` 在 `core` 下，所以可以 import `core/internal`；runtime 是 public module，API 组合根直接依赖它。

### 4.5 `cli/go.mod`

```go
module yukumo-script/cli

go 1.27

require (
	github.com/charmbracelet/huh v1.0.0
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/spf13/cobra v1.10.2
	yukumo-script/core/api v0.0.0
	yukumo-script/core v0.0.0 // indirect
	yukumo-script/runtime v0.0.0 // indirect
)

replace (
	yukumo-script/core/api => ../core/api
	yukumo-script/core => ../core
	yukumo-script/runtime => ../runtime
)
```

说明：

- CLI 源码不 import core/internal 或 runtime；core/runtime 只是 API 的传递依赖，列出来是为了 `GOWORK=off` / `go mod tidy` 能解析。
- `replace` 是本地 monorepo 的解析兜底；发布模块时移除并打 tag。
- 根目录没有 `go.mod`，也没有根 `go.sum`；旧的根 go.sum 删除。

### 4.6 入口

`cli/cmd/yukumo-script/main.go`：

```go
package main

import (
	"fmt"
	"os"

	"yukumo-script/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

构建与运行：

```bash
# 手动方式（workspace 模式，适合 IDE 之外临时验证）
go build -o yukumo-script ./cli/cmd/yukumo-script
go run ./cli/cmd/yukumo-script

# 推荐方式：统一走 justfile
just build
just run
```

不再使用 `go run .` / `go build .`。

### 4.7 Justfile：跨平台命令集线器

根目录新增 `justfile`。所有构建、测试、vet、depguard/staticcheck、replace 检查都从这里进，CI 也调用同一组 recipe，避免 CI 和本地命令漂移。

核心设计：

- `export GOWORK := "off"`：所有 recipe 默认在 module 模式下跑，强制每个 module 用自己的 `go.mod` + `replace`。
- Go 命令优先用 `go -C <dir>`（Go 1.20+），避免依赖 `cd` 的 shell 差异。
- `build` 依赖 `lint`，所以 depguard 和 staticcheck 一定在 build 前执行。
- `ci` 依赖 `check-replaces`、`build`、`test`、`vet`，顺序固定。
- `just` 本身跨平台；recipe 只使用 `go`、`cd`、`&&` 这类 POSIX / cmd 都支持的命令，不依赖 bash 专有语法。

示意：

```just
# justfile
export GOWORK := "off"

default: check

# 本地完整检查
check: check-replaces lint vet test

# CI 入口：先查 replace，再 build（build 会先跑 depguard + staticcheck），最后 test/vet
ci: check-replaces build test vet

# build 前强制 depguard + staticcheck
build: lint
    go -C cli build -o ../yukumo-script ./cmd/yukumo-script

run *ARGS:
    go -C cli run ./cmd/yukumo-script {{ARGS}}

install-tools:
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

test:
    go -C core test ./...
    go -C runtime test ./...
    go -C core/api test ./...
    go -C cli test ./...

vet:
    go -C core vet ./...
    go -C runtime vet ./...
    go -C core/api vet ./...
    go -C cli vet ./...

lint:
    cd core && golangci-lint run
    cd runtime && golangci-lint run
    cd core/api && golangci-lint run
    cd cli && golangci-lint run

check-replaces:
    go run scripts/check-replaces.go
```

说明：

- `go -C cli build -o ../yukumo-script ./cmd/yukumo-script` 的输出落在仓库根目录；因为 `GOWORK=off`，它实际验证的是 CLI module 的 `replace` 链是否完整。
- `go -C cli run ./cmd/yukumo-script` 的进程 cwd 是 `cli/`；runtime 的 WASM loader 会从 `cli/` 向上找到根 `third-party/aq_wasm.wasm`，不需要移动 wasm。
- `golangci-lint` 需要提前安装（`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`）；CI 会装，本地可用 `just install-tools`（可选）。
- 每个 recipe 的每一行由 just 独立执行；`cd core && golangci-lint run` 写成一行，避免 cwd 丢失。
- `go.work` 仍保留给 IDE 和手动 workspace 开发；`just` 的检查和构建全部走 `GOWORK=off`，CI 额外做一次 workspace build 验证 `go.work`。

---

<a id="sec-5"></a>
## 5. 编译期边界与 linter

### 5.1 编译期硬边界

```text
yukumo-script/core/api          路径在 yukumo-script/core 下 → 可以 import core/internal
        ↓
yukumo-script/core/internal/... internal 只允许 yukumo-script/core 前缀

yukumo-script/cli               路径不在 core 下 → 禁止 import core/internal
yukumo-script                   根目录没有 module，也没有 main
```

验证：临时在 CLI 加：

```go
import _ "yukumo-script/core/internal/domain"
```

期望：

```text
use of internal package yukumo-script/core/internal/domain not allowed
```

### 5.2 depguard 负责方向，staticcheck 负责通用检查

runtime 顶层化后是 public package，编译器不会阻止 CLI import `yukumo-script/runtime/aq`。因此用 depguard 检查：

- core service 不能 import runtime；
- runtime 不能 import 任何项目包；
- cli 不能 import core/internal 或 runtime；
- cmd 不能 import core 或 runtime；
- API 不能 import cli；
- UI 库只在 cli，wazero 只在 runtime，yaml 只在 core。

staticcheck 是通用静态检查，不参与分层规则；它和 depguard 一样在每个 module 的 `.golangci.yml` 中显式启用，由 `just lint` 一起执行。

每个 module 一份 `.golangci.yml`：

```text
core/.golangci.yml
runtime/.golangci.yml
core/api/.golangci.yml
cli/.golangci.yml
```

CI / 本地统一用 `just lint`；它内部在每个 module 目录里跑 `GOWORK=off golangci-lint run`，这样 depguard 的 `files` 路径相对各自 module 根，规则清晰，staticcheck 也按 module 独立跑。

示例（`core/.golangci.yml` 的结构示意）：

```yaml
version: "2"

linters:
  enable:
    - depguard
    - staticcheck
  settings:
    depguard:
      rules:
        domain:
          files:
            - "internal/domain/**"
          deny:
            - pkg: "yukumo-script/"
              desc: "domain is the bottom layer"
        store:
          files:
            - "internal/store/**"
          deny:
            - pkg: "yukumo-script/"
              desc: "store is infrastructure; must not import application layers"
        service:
          files:
            - "internal/service/**"
          deny:
            - pkg: "yukumo-script/runtime"
              desc: "service owns the engine port; runtime is an adapter"
            - pkg: "yukumo-script/core/api"
              desc: "api depends on service, not the other way around"
            - pkg: "yukumo-script/cli"
              desc: "service must not depend on the interaction layer"
            - pkg: "github.com/tetratelabs/wazero"
            - pkg: "github.com/spf13/cobra"
            - pkg: "github.com/charmbracelet/huh"
            - pkg: "github.com/charmbracelet/lipgloss"
        core-test:
          files:
            - "test/**"
          deny:
            - pkg: "yukumo-script/runtime"
            - pkg: "yukumo-script/core/api"
            - pkg: "yukumo-script/cli"
        core-external:
          files:
            - "**/*.go"
          deny:
            - pkg: "github.com/tetratelabs/wazero"
            - pkg: "github.com/spf13/cobra"
            - pkg: "github.com/charmbracelet/huh"
            - pkg: "github.com/charmbracelet/lipgloss"
```

> depguard 的 `files` 是相对 module 根的文件 glob；具体字段名以 CI 使用的 golangci-lint v2 版本为准。若版本支持 `allow` / `list-mode: strict`，可以把 service 改成 allowlist，进一步收紧。

### 5.3 为什么还要 `GOWORK=off`

`go.work` 会让 workspace 内所有 module 都成为 main module，本地可以跨 module import。CI 对每个 module 跑：

```bash
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off go test ./...
```

这样每个 module 只能靠自己的 `go.mod` + `replace` 解析依赖，模块边界在构建层面也成立。

---

<a id="sec-6"></a>
## 6. API 层设计

### 6.1 对外 DTO

API 不暴露 `domain` / `store` / `aq` 类型，只暴露自己的数据结构。

```go
package api

type Character struct {
	Name        string
	Description string
	Voices      []Voice
}

type Voice struct {
	Name   string
	Params []SynthParam
}

type SynthParam struct {
	Version   int
	Preset    int
	Speed     int
	PhontData []byte
	Base      int
	Volume    int
	Pitch     int
	Accent    int
	Lmd       int
	Fsc       int
}
```

配置 DTO 使用 effective value：

```go
type Config struct {
	Characters CharacterConfig
	Paths      PathConfig
	CLI        CLIConfig
	Synth      SynthConfig
}

type CharacterConfig struct {
	InitDefaultVoice bool
}

type PathConfig struct {
	CharacterPath string
	OutputPath    string
}

type CLIConfig struct {
	DefaultInteractive bool
	WithAccessibility  bool
}

type SynthConfig struct {
	OutputFormat string
	AlwaysWrite  bool
	ReuseCache   bool
}
```

> 行为差异：当前 domain 的 `*bool` 在 nil 时保存可能写出 `null`；API `SaveConfig` 会把 effective bool 写成显式 `true/false`。这是有意的清理。若要求 YAML 逐字节不变，把 `CLIConfig` / `SynthConfig` 改成 `*bool`。

### 6.2 模型辅助

```go
const (
	Version1  = 1
	Version2  = 2
	Version10 = 10
)

func DefaultSynthParam(version int) SynthParam
func (p *SynthParam) Normalize()
func BaseName(base int) string
func BaseNames() []string
func SetPhontPath(p *SynthParam, path string) error
```

实现转成 `domain` 后调用，避免在 API 复制 clamp / 默认值规则。

错误保持 `errors.Is` 可用：

```go
var (
	ErrInvalid  = domain.ErrInvalid
	ErrNotFound = domain.ErrNotFound
	ErrExists   = domain.ErrExists
)
```

### 6.3 `App` 与用例方法

```go
type Options struct {
	ConfigDir string // 空则读 YUKUMO_CONFIG_DIR / UserConfigDir
}

func Open(opts Options) (*App, error)
func (a *App) Close(ctx context.Context) error

func (a *App) Config() (Config, error)
func (a *App) SaveConfig(c Config) error

func (a *App) ListCharacters() ([]Character, error)
func (a *App) GetCharacter(name string) (Character, error)
func (a *App) AddCharacter(c Character) error
func (a *App) UpdateCharacter(oldName string, c Character) error
func (a *App) DeleteCharacter(name string) error

func (a *App) ListVoices(charName string) ([]Voice, error)
func (a *App) GetVoice(charName, voiceName string) (Voice, error)
func (a *App) AddVoice(ctx context.Context, charName string, v Voice) error
func (a *App) UpdateVoice(ctx context.Context, charName, oldName string, v Voice) error
func (a *App) DeleteVoice(charName, voiceName string) error
func (a *App) DuplicateVoice(ctx context.Context, charName, srcName, newName string) (string, error)

func (a *App) PresetNames(ctx context.Context, version int) ([]string, error)
func (a *App) HasPreset(ctx context.Context, version, preset int) (bool, error)

func (a *App) Synthesize(ctx context.Context, req SynthRequest) (SynthResult, error)
```

```go
type SynthRequest struct {
	VoiceRef      string
	VoicePathSpec string
	Text          string
	Version       int
	SpeedFactor   float64
	Output        string
	Preview       bool
	ReuseCache    *bool
}

type SynthResult struct {
	WAV       []byte
	Path      string
	Wrote     bool
	FromCache bool
}
```

`App` 持有：

```go
config *config.Service
chars  *character.Service
voices *voice.Service
synth  *synth.Service
engine *aq.Pool
```

`Open` 负责组合根：

1. 解析 config dir（`Options.ConfigDir` > `YUKUMO_CONFIG_DIR` > `UserConfigDir/yukumo-script`）。
2. `store.NewFS` → `config.New` → `Get`。
3. `EnsurePaths`。
4. `aq.NewPool`（runtime/aq）。
5. 构造 character / voice / synth service；synth 使用 `aqEngine` adapter。
6. `CleanIndex`；失败时关闭 pool 并返回错误。
7. CLI 不再知道 store / aq / service 的存在。

`DuplicateVoice` 需要把内部 `voice.Service.Duplicate` 改成返回最终名字：

```go
func (s *Service) Duplicate(ctx context.Context, charName, srcName, newName string) (string, error)
```

这样 CLI 不再依赖 `service/nameutil`。

### 6.4 DTO 转换

`core/api/convert.go`：

```go
toCharacter(domain.Character) Character
fromCharacter(Character) domain.Character
toVoice / fromVoice
toParam / fromParam
toConfig / fromConfig
```

要求：

- 切片和 `PhontData` 深拷贝，避免 API 边界共享底层数组。
- `fromConfig` 把 API bool 转成 domain `*bool`。
- 默认值、clamp、preset 校验仍由 core 负责，API 不复制业务规则。
- API 导出签名只出现 API 类型和标准库类型。

### 6.5 API 作为集成测试位置

API 同时依赖 core 和 runtime，是天然的集成测试层：

```text
core/api/test/runtime/aq/   → 真实 WASM runtime 测试
core/api/test/synth/        → 真实 aq.Pool 的 synth service 集成测试
```

这样：

- `core/go.mod` 不需要 runtime / wazero；
- `runtime/go.mod` 只保留自己的单元测试；
- API 负责“core + runtime 一起工作”的验证；
- CLI 的测试只测 prompt / 输出格式，不碰 WASM。

---

<a id="sec-7"></a>
## 7. Engine port 解耦

为了让 `core` 在**生产代码和 go.mod** 上都不依赖 runtime，`service/synth` 不直接使用 `aq.Param`。

### 7.1 service 侧

`core/internal/service/synth` 定义：

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

type Engine interface {
	Synthesize(ctx context.Context, text string, p EngineParam) ([]byte, error)
}
```

`service/synth` 内部：

```go
wav, err := s.engine.Synthesize(ctx, text, toEngineParam(p, speed))
```

删除：

- `import "yukumo-script/runtime/aq"`
- `toAQParam`
- `Engine` 接口里的 `aq.Param`

### 7.2 API adapter

`core/api/aqengine.go`：

```go
type aqEngine struct {
	pool *aq.Pool
}

func (e aqEngine) Synthesize(ctx context.Context, text string, p synth.EngineParam) ([]byte, error) {
	return e.pool.Synthesize(ctx, text, aq.Param{
		Version: int32(p.Version),
		Preset:  int32(p.Preset),
		Speed:   int32(p.Speed),
		Phont:   p.Phont,
		Base:    int32(p.Base),
		Volume:  int32(p.Volume),
		Pitch:   int32(p.Pitch),
		Accent:  int32(p.Accent),
		Lmd:     int32(p.Lmd),
		Fsc:     int32(p.Fsc),
	})
}
```

`Open`：

```go
pool := aq.NewPool()
engine := aqEngine{pool: pool}
synthSvc := synth.New(cfgSvc, voiceSvc, engine)
```

`voice.New(charSvc, pool)` 不变，因为 `PresetCatalog` 接口只用原始类型，不暴露 `aq` 类型。

### 7.3 API 集成测试里的 adapter

`core/api/test/synth` 里增加一个等价的测试 adapter，仍然用真实 `aq.Pool` 跑 WASM：

```go
type aqEngine struct{ pool *aq.Pool }

func (e aqEngine) Synthesize(ctx context.Context, text string, p synth.EngineParam) ([]byte, error) {
	return e.pool.Synthesize(ctx, text, aq.Param{ /* 同上映射 */ })
}
```

约 15 行重复，换 core 的干净依赖图，值得。

---

<a id="sec-8"></a>
## 8. CLI 改造

### 8.1 允许的 import

`cli` 的 Go 文件只允许：

```text
标准库
yukumo-script/core/api
yukumo-script/cli/...（同 module 内部）
github.com/spf13/cobra
github.com/charmbracelet/huh
github.com/charmbracelet/lipgloss
```

不允许：

```text
yukumo-script/core/internal/...
yukumo-script/runtime/...
yukumo-script/core 的其他包
```

`cli/cmd/yukumo-script/main.go` 只允许 import `yukumo-script/cli` + 标准库。

### 8.2 启动流程

`cli.Execute`：

```go
func Execute() error {
	app, err := api.Open(api.Options{})
	if err != nil {
		return err
	}
	defer func() { _ = app.Close(context.Background()) }()

	cfg, err := app.Config()
	if err != nil {
		return err
	}
	prompt.SetAccessible(cfg.CLI.WithAccessibility)

	root := newRoot(app, cfg.CLI.DefaultInteractive)
	if cfg.CLI.DefaultInteractive {
		enableInteractiveBrowse(root)
	}
	return root.ExecuteContext(context.Background())
}
```

删除：

- `configDir()`
- `store.NewFS`
- `svcconfig.New`
- `aq.NewPool`
- `synthSvc.CleanIndex`
- `root.PersistentPostRun` 里关闭 pool（改为 defer）
- 所有 `svccharacter` / `svcvoice` / `svcconfig` / `svcsynth` / `store` / `aq` import

`newRoot`：

```go
func newRoot(app *api.App, interactive bool) *cobra.Command
```

### 8.3 命令构造器

从“多个 service”改成“一个 `*api.App`”：

```go
cliconfig.NewCmd(app, cfgOpts)
clicharacter.NewCmd(app, charOpts)
clisynth.NewCmd(app, synthOpts)
```

对应内部：

- `character.NewCmd(app, opts)`
- `voice.Dispatch(cmd, app, opts, charName, args)`
- `voice.RunMenu(cmd, app, charName)`
- `synth.NewCmd(app, opts)`
- `newVoiceCmd(app, opts)`

### 8.4 类型替换

| 现在 | 之后 |
|---|---|
| `domain.Character` | `api.Character` |
| `domain.Voice` | `api.Voice` |
| `domain.SynthParam` | `api.SynthParam` |
| `domain.Config` | `api.Config` |
| `domain.Version1/2/10` | `api.Version1/2/10` |
| `domain.DefaultSynthParam` | `api.DefaultSynthParam` |
| `domain.BaseName` / `BaseNames` | `api.BaseName` / `api.BaseNames()` |
| `domain.ErrInvalid` | `api.ErrInvalid` |
| `p.Normalize()` | `p.Normalize()`（API 方法） |
| `svcvoice.SetPhontPath` | `api.SetPhontPath` |
| `nameutil.ResolveCopy` | 删除，使用 `app.DuplicateVoice` 返回值 |

### 8.5 重点文件

- `cli/cmd/yukumo-script/main.go`：新入口。
- `cli/root.go`：只保留 `Execute` / `newRoot` / 命令注册。
- `cli/browse.go`：不变。
- `cli/config/fields.go`：`setting.Entry[api.Config]`，直接访问 `c.Paths` / `c.Synth` / `c.CLI`。
- `cli/config/edit.go`：`app.Config()` / `app.SaveConfig(cfg)`。
- `cli/config/get.go`、`list.go`：`app.Config()`。
- `cli/character/add.go`：构造 `api.Character`，调用 `app.AddCharacter`。
- `cli/character/edit.go`：`api.Character`，`app.GetCharacter` / `app.UpdateCharacter`；`voices` 菜单传 `app`。
- `cli/character/fields.go`：`setting.Entry[api.Character]`。
- `cli/character/voice/fields.go`：`setting.Entry[api.SynthParam]`；`api.BaseName`、`api.BaseNames()`、`api.SetPhontPath`。
- `cli/character/voice/param.go`：`app.PresetNames` / `app.HasPreset`；`api.BaseNames()`。
- `cli/character/voice/duplicate.go`：`finalName, err := app.DuplicateVoice(...)`，不再 import `nameutil`。
- `cli/synth/root.go`：`api.SynthRequest`，`app.Synthesize`。
- `cli/synth/run.go`：`api.SynthResult`。
- `cli/synth/voice_form.go`：`app.ListCharacters` / `app.ListVoices` / `api.SynthRequest`。
- `cli/synth/pick.go`：`app.ListCharacters` / `app.ListVoices`。
- `cli/synth/preview.go`：不变。

CLI 的输出文本、交互流程、预览逻辑保持现状。

---

<a id="sec-9"></a>
## 9. 测试分层

### 9.1 迁移表

| 现在 | 之后 | 类型 |
|---|---|---|
| `test/runtime/aq/runtime_test.go` | `core/api/test/runtime/aq/runtime_test.go` | runtime 集成 |
| `test/runtime/aq/synthesize_test.go` | `core/api/test/runtime/aq/synthesize_test.go` | runtime 集成 |
| `test/service/synth/cache_test.go` | `core/api/test/synth/cache_test.go` | synth 集成 |
| `test/service/synth/service_test.go` | `core/api/test/synth/service_test.go` | synth 集成 |
| `test/service/synth/filename_test.go` | `core/test/service/synth/filename_test.go` | 纯单元 |
| `test/service/synth/voice_test.go` | `core/test/service/synth/voice_test.go` | 纯单元 |
| `test/service/character/*` | `core/test/service/character/*` | 纯单元 |
| `test/service/character/voice/*` | `core/test/service/character/voice/*` | 纯单元 |
| `test/service/config/*` | `core/test/service/config/*` | 纯单元 |
| `test/service/nameutil/*` | `core/test/service/nameutil/*` | 纯单元 |
| `test/cli/prompt/*` | `cli/test/prompt/*` | CLI 单元 |
| `internal/runtime/aq/errors_test.go` | `runtime/aq/errors_test.go` | runtime 包内单元 |
| `internal/cli/synth/run_test.go` | `cli/synth/run_test.go` | CLI 包内单元 |

### 9.2 结果

- `core/test/service/synth` 只剩 `filename_test.go`、`voice_test.go`，不 import runtime。
- `core/go.mod` 只有 yaml；`GOWORK=off go list -m all` 在 core 下不会出现 wazero。
- `runtime` module 只保留 `errors_test.go`，不跑 WASM 集成。
- API module 的测试既 import core internal，也 import runtime，覆盖真实 WASM 合成、缓存、Pool 并发。
- CLI module 的测试只覆盖 prompt 和输出格式。

### 9.3 新增测试

1. `core/api/convert_test.go`
   - Character / Voice / SynthParam / Config 往返转换。
   - 深拷贝：修改 DTO 不影响原 domain 数据。
   - `Normalize` / `DefaultSynthParam` / `BaseName` 与 domain 一致。
   - Config bool ↔ `*bool` 映射。

2. `core/api/app_test.go`（建议）
   - `api.Open(api.Options{ConfigDir: t.TempDir()})`。
   - 角色 CRUD、配置读写。
   - 添加一个无参数 voice（不触发 WASM preset 校验），验证 voice CRUD。
   - `Close`。

3. API 集成测试保留真实 `aq.Pool` adapter，确保迁移没有把合成测试变成 fake。

---

<a id="sec-10"></a>
## 10. 静态检查：depguard + staticcheck

### 10.1 原则

- 编译器负责“CLI 不能 import core/internal”。
- depguard 负责“所有 module 的单向依赖 + 外部库归属”。
- staticcheck 负责“通用静态缺陷”，补齐之前缺失的检查。
- `GOWORK=off` 负责“每个 module 的 go.mod 真的能独立解析”。

### 10.2 每个 module 一份配置

每个 `.golangci.yml` 至少显式启用（避免默认 linter 集随版本变化）：

```yaml
linters:
  enable:
    - depguard
    - staticcheck
```

```text
core/.golangci.yml
runtime/.golangci.yml
core/api/.golangci.yml
cli/.golangci.yml
```

CI / 本地统一入口：

```bash
just lint
```

`just lint` 内部等价于：

```bash
(cd core && GOWORK=off golangci-lint run)
(cd runtime && GOWORK=off golangci-lint run)
(cd core/api && GOWORK=off golangci-lint run)
(cd cli && GOWORK=off golangci-lint run)
```

### 10.3 depguard 规则表

| module | 文件 | 禁止 |
|---|---|---|
| core | `internal/domain/**` | 所有 `yukumo-script/` |
| core | `internal/store/**` | 所有 `yukumo-script/` |
| core | `internal/service/**` | `yukumo-script/runtime`、`yukumo-script/core/api`、`yukumo-script/cli`、UI 库、wazero |
| core | `test/**` | `yukumo-script/runtime`、`yukumo-script/core/api`、`yukumo-script/cli` |
| core | `**/*.go` | UI 库、wazero（全局兜底） |
| runtime | `**/*.go` | 所有 `yukumo-script/`、UI 库、yaml |
| core/api | `**/*.go` | `yukumo-script/cli`、UI 库、wazero、yaml |
| cli | `**/*.go` | `yukumo-script/core/internal`、`yukumo-script/runtime`、wazero |
| cli | `cmd/**` | `yukumo-script/core`、`yukumo-script/runtime`、UI 库、wazero、yaml |

外部库归属：

| 库 | 只允许 |
|---|---|
| `github.com/spf13/cobra` | `cli` |
| `github.com/charmbracelet/huh` | `cli` |
| `github.com/charmbracelet/lipgloss` | `cli` |
| `github.com/tetratelabs/wazero` | `runtime` |
| `gopkg.in/yaml.v3` | `core`、`cli`（CLI YAML 输出） |

### 10.4 staticcheck

- 4 个 module 的 `.golangci.yml` 都显式启用 `staticcheck`。
- `just lint` 会同时跑 depguard 和 staticcheck；`build` 依赖 `lint`，所以 CI 的 build 前一定经过 staticcheck。
- 首次启用可能暴露存量问题。处理原则：
  1. 优先直接修。
  2. 确实需要保留的，按行/按规则精确 `exclude`，并在旁边写原因或 issue 链接。
  3. 禁止全局关闭 staticcheck，也不要用宽泛的 `//nolint:staticcheck`。
- 测试文件也在检查范围；本仓库没有 generated Go 代码需要排除。
- staticcheck 与 `go vet` 互补：`go vet` 继续保留在 `just vet`，`just ci` 两者都跑。

### 10.5 负向验证

在 CLI 临时加：

```go
import _ "yukumo-script/runtime/aq"
```

期望 `just lint` 报 depguard violation；删除后恢复。这个验证说明 runtime 顶层化后仍不能绕过 API。

---

<a id="sec-11"></a>
## 11. replace 一致性脚本

位置：`scripts/check-replaces.go`，只用标准库；根目录没有 module，但 `go run scripts/check-replaces.go` 可以直接运行（Go 支持单文件运行）。

脚本做什么：

1. 在仓库根执行 `go work edit -json`，拿到 `use` 的 module 目录。
   - 注意：justfile 全局 `export GOWORK := "off"` 后，子进程默认也会继承 `GOWORK=off`，而 `GOWORK=off go work edit` 会报 “no go.work file found”。脚本在调用 `go work edit -json` 时必须显式把子进程环境改成 `GOWORK=auto`。
2. 对每个 module 目录执行 `GOWORK=off go mod edit -json`，解析 `Replace`。
3. 只处理本地 replace（`New.Version == ""`）：
   - 把 `New.Path` 相对 module 目录解析成绝对路径，`EvalSymlinks` 规范化。
4. 按 `Old.Path` 分组，检查：
   - 所有 go.mod 中同一个 `Old.Path` 的本地 replace 都指向**同一个目录**；
   - 目标目录存在，且其中的 `go.mod` 的 `module` 等于 `Old.Path`；
   - `go.work` 里的每个 module 都能在需要的地方被 replace 到。
5. 打印实际映射，例如：

```text
yukumo-script/core       -> core
yukumo-script/runtime    -> runtime
yukumo-script/core/api   -> core/api
```

6. 任何冲突、缺失、路径不存在、module path 不匹配都 exit 1。

用途：

- runtime 改名/移动时，API 和 CLI 的 replace 必须同步；脚本会直接报错。
- root 没有 go.mod 后，所有本地解析都依赖各 module 的 replace + go.work；这个脚本是防漏改的关键。

justfile 的 `check-replaces` recipe 就是跑它；`just ci` 和 CI 的第一步都会调用：

```bash
just check-replaces
# 等价于
go run scripts/check-replaces.go
```

---

<a id="sec-12"></a>
## 12. CI

新增 `.github/workflows/ci.yml`，CI 只负责装工具和调用 `just ci`：

```yaml
name: ci

on:
  push:
  pull_request:

jobs:
  checks:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.27"

      - uses: extractions/setup-just@v2

      - name: Install golangci-lint
        run: |
          go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
          echo "$(go env GOPATH)/bin" >> "$GITHUB_PATH"

      - name: CI
        run: just ci

      # 额外用 workspace 模式验证 go.work 本身
      - name: Workspace build
        run: go build -o yukumo-script ./cli/cmd/yukumo-script
```

`just ci` 的依赖顺序：

```text
ci
├── check-replaces      # scripts/check-replaces.go
├── build
│   └── lint            # depguard + staticcheck，先于 build
├── test                # 4 个 module，GOWORK=off
└── vet                 # 4 个 module，GOWORK=off
```

说明：

- justfile 全局 `export GOWORK := "off"`，所以 `just ci` 里的 build/test/vet/lint 都是 module 模式，验证各自的 `replace` 够用。
- `check-replaces` 用 Go 脚本，跨平台，不依赖 bash。
- `golangci-lint` 用 `go install` 安装到 `$(go env GOPATH)/bin`；不要用 action 默认的“在仓库根跑 lint”模式，因为根目录没有 module。
- 最后的 workspace build 不经过 just，保留默认 `GOWORK=auto`，专门验证 `go.work` 能解析 4 个 module 和入口。
- API 的集成测试需要 `third-party/aq_wasm.wasm`；checkout 自带，loader 会从 `core/api/test/...` 向上找到根 `third-party/`。
- 本地等价命令就是 `just ci`；只想跑测试时用 `just test`。

---

<a id="sec-13"></a>
## 13. 迁移步骤

> 结构性搬迁建议在独立分支完成。中间 commit 可能不可构建；不要为了“每一步都绿”引入兼容 shim，成本更高。

### Phase 0：基线

1. 记录当前 commit。
2. 跑 `go test ./...`、`go vet ./...`，确认全绿。
3. 确认 `third-party/aq_wasm.wasm` 存在。

### Phase 1：建立新目录并搬文件

0. 先建目标目录，避免 `git mv` 因父目录不存在失败：
   ```bash
   mkdir -p core/internal runtime/aq cli/cmd/yukumo-script
   ```
1. `git mv internal/domain core/internal/domain`
2. `git mv internal/store core/internal/store`
3. `git mv internal/service core/internal/service`
4. `git mv internal/runtime/aq runtime/aq`
5. `git mv internal/cli cli`
6. `git mv main.go cli/cmd/yukumo-script/main.go`
7. 更新 import：
   - `yukumo-script/internal/domain` → `yukumo-script/core/internal/domain`
   - `yukumo-script/internal/store` → `yukumo-script/core/internal/store`
   - `yukumo-script/internal/service/...` → `yukumo-script/core/internal/service/...`
   - `yukumo-script/internal/runtime/aq` → `yukumo-script/runtime/aq`
   - `yukumo-script/internal/cli/...` → `yukumo-script/cli/...`
8. 此时 CLI 还不能 import core/internal，会红；继续完成 API 和 CLI 切换。

### Phase 2：module 文件与 go.work

1. 创建 `core/go.mod`（只 yaml）。
2. 创建 `runtime/go.mod`（wazero）。
3. 创建 `core/api/go.mod`（core + runtime + replace）。
4. 创建 `cli/go.mod`（api + UI + 传递 local replace）。
5. 创建根 `go.work`，`use` 4 个 module。
6. 删除根 `go.mod`、根 `go.sum`。
7. `go work sync` 或逐模块 `go mod tidy`。

### Phase 3：API 与 engine port

1. `service/synth` 增加 `EngineParam` / `Engine`，删除 `aq` import 和 `toAQParam`。
2. 创建 `core/api` 的 DTO / convert / App / aqengine。
3. 把 `voice.Service.Duplicate` 改成返回最终名字。
4. 把 runtime 集成测试、synth 集成测试移到 `core/api/test/...`。
5. `core/test/service/synth` 只留纯函数测试。
6. 跑：
   ```bash
   (cd core && GOWORK=off go test ./...)
   (cd runtime && GOWORK=off go test ./...)
   (cd core/api && GOWORK=off go test ./...)
   ```

### Phase 4：CLI 切换

1. `cli/root.go` 改用 `api.Open` / `app.Config()` / defer `Close`。
2. 所有命令构造器接收 `*api.App`。
3. 替换 domain 类型为 API DTO。
4. 把 `cli/cmd/yukumo-script/main.go` 的 import 改成 `yukumo-script/cli`。
5. 跑：
   ```bash
   (cd cli && GOWORK=off go test ./...)
   go build -o yukumo-script ./cli/cmd/yukumo-script
   ```

### Phase 5：linter、justfile 与 CI

1. 调整 `.gitignore`，白名单放行 `.github/**` 和 `**/.golangci.yml`（否则 `*.yml` 会忽略它们）。
2. 添加 4 个 `.golangci.yml`，显式启用 `depguard` + `staticcheck`。
3. 添加 `scripts/check-replaces.go`。
4. 添加根 `justfile`（见 4.7 节）：
   - `export GOWORK := "off"`
   - `build: lint`
   - `test` / `vet` / `lint` 覆盖 4 个 module
   - `check-replaces`
   - `ci: check-replaces build test vet`
5. 添加 `.github/workflows/ci.yml`，CI 调用 `just ci`，并额外做一次 workspace build。
6. 本地验证：
   ```bash
   just check-replaces
   just lint
   just build
   just test
   just vet
   just ci
   ```

### Phase 6：文档

更新 README.md / README.zh.md：

- 新架构图。
- 新构建命令：`just build`（底层是 `go -C cli build ...`）。
- 新运行命令：`just run`。
- 不再有根 go.mod / `go run .`。
- `just test` / `just lint` / `just ci` 说明。
- 4 个 module 的独立检查说明。
- depguard / staticcheck / CI / replace 脚本说明。

---

<a id="sec-14"></a>
## 14. 验收清单

### 14.1 根目录结构

```bash
test ! -f go.mod
test ! -f go.sum
test ! -f main.go
test -f go.work
test -f justfile
test ! -f scripts/test.sh
```

### 14.2 构建

```bash
just build
just run
```

手动验证 help 可以直接走 module 命令：

```bash
go -C cli run ./cmd/yukumo-script --help
```

### 14.3 模块独立检查

```bash
just ci
```

`just ci` 会依次执行 replace 检查、depguard + staticcheck、build、test、vet；统一走 `GOWORK=off`。需要单独跑时：

```bash
just check-replaces
just lint
just build
just test
just vet
```

`just lint` 必须零 staticcheck findings；如果确有无法立即修的历史问题，只允许按 file/rule 精确 exclude，并附原因。

### 14.4 replace 一致性

```bash
just check-replaces
```

期望输出 core/runtime/core/api 三个映射且无冲突。

### 14.5 core 不再依赖 wazero

```bash
(cd core && GOWORK=off go list -m all | grep -i wazero)   # 应为空
grep -R "yukumo-script/runtime" core --include='*.go'     # 应为空
grep -R "tetratelabs/wazero" core --include='*.go'        # 应为空
```

### 14.6 CLI 边界

```bash
grep -R "yukumo-script/core/internal" cli --include='*.go'  # 应为空
grep -R "yukumo-script/runtime" cli --include='*.go'        # 应为空
```

### 14.7 编译期负向验证

临时在 CLI 加 `import _ "yukumo-script/core/internal/domain"`，确认编译报 `use of internal package ... not allowed`，删除。

### 14.8 depguard 负向验证

临时在 CLI 加 `import _ "yukumo-script/runtime/aq"`，确认 `just lint` 报 depguard violation，删除。

### 14.9 行为回归

- 现有测试全绿，测试数量不减少。
- API 集成测试真实跑 WASM。
- 手动 smoke：
  ```bash
  YUKUMO_CONFIG_DIR=$(mktemp -d) go -C cli run ./cmd/yukumo-script config list
  YUKUMO_CONFIG_DIR=$(mktemp -d) go -C cli run ./cmd/yukumo-script character add remilia "example"
  YUKUMO_CONFIG_DIR=$(mktemp -d) go -C cli run ./cmd/yukumo-script character list
  ```
  或使用 `just run` 的交互流程。

---

<a id="sec-15"></a>
## 15. 风险与取舍

1. **根目录没有 go.mod**
   - `go test ./...` / `go run .` 不再可用；日常用 `just ...`，底层是 `go -C <module>`。
   - 优点是没有“根 module 偷偷依赖所有层”的入口。

2. **main 放在 CLI module**
   - 不需要额外 `cmd` module；入口天然属于交互层。
   - 如果以后想拆独立 `cmd` module，可以再把 `cli/cmd` 移出，成本很低。

3. **runtime 顶层化后是 public**
   - 编译器不阻止 CLI import runtime；由 depguard 阻止。
   - 如果要求编译器硬阻止，需要把 runtime 放回 `core/internal`，但那样 runtime 又和 core 绑在一起，与“移出 core”冲突。

4. **depguard / staticcheck 是 linter，不是编译器**
   - 本地只跑 `go build` 不会发现反向依赖或 staticcheck 问题；`just lint` / CI 会。
   - `GOWORK=off` 的模块检查补上 go.mod 层面的约束。

5. **staticcheck 首次启用**
   - 可能暴露存量问题；优先修，其次精确 exclude 并写原因，禁止全局关闭。
   - 如果某条规则暂时无法修，在对应 `.golangci.yml` 里按 file/rule 精确排除，并登记 issue。

6. **core 测试全部无 runtime**
   - runtime/synth 集成测试迁到 API，API 必须保留 `third-party/aq_wasm.wasm` 的路径查找能力。
   - CI checkout 自带 wasm；本地从 API 目录跑测试也能向上找到。

7. **replace 矩阵**
   - API 和 CLI 需要本地 replace；runtime/core 不需要。
   - `scripts/check-replaces.go` 防止改名漏改。
   - 发布模块时移除 replace 并打 tag。

8. **DTO 映射**
   - 约 200–300 行转换代码；默认值 / Normalize 委托 domain，加 round-trip 测试。
   - Config bool 从指针语义改成 effective value；若要求逐字节兼容则改回 `*bool`。

9. **行为不变**
   - 本次只做分层与依赖倒置，不改 CLI 输出、交互流程、合成策略、文件格式。
   - 唯一有意差异是 Config 的显式 bool 写回。

10. **未来 GUI**
   - 新增 `gui/` module，路径 `yukumo-script/gui`。
   - 只能 import `yukumo-script/core/api`；depguard 禁止 core/internal 和 runtime。
   - `go.work` 增加 `./gui`，replace 脚本自动核对。

11. **`.gitignore` 的 YAML 通配**
    - 当前 `*.yml` / `*.yaml` 会忽略 `.github/workflows/ci.yml` 和 `.golangci.yml`。
    - 必须按 2.1 节加白名单，否则 CI、depguard、staticcheck 配置不会进 Git。

12. **justfile 与 GOWORK=off**
    - justfile 全局 `export GOWORK := "off"`；`go work edit` 必须由脚本显式切回 `GOWORK=auto`。
    - recipe 只用 `go -C`、`cd ... && ...`，不依赖 bash 专有语法，macOS/Linux/Windows（cmd 或 Git Bash）都能跑。
    - `golangci-lint` 和 `just` 本身需要安装；CI 负责，本地可用 `just install-tools` 装 linter。

---

<a id="sec-16"></a>
## 16. Web/WASM 可移植性 addendum

> 状态：评估与后续增量设计，不属于当前 Phase 0–6 的实现范围。未启用 web 前，前 15 节按现状执行；启用 web 前需要落地本节列出的 port / 组合根变更。
>
> 结论：当前分层方向可复用，但 store 解耦**不足以直接支撑 JS + IndexedDB/localStorage**。需要改架构，主要改 port、依赖注入、平台组合根和 synth 的 IO 覆盖；不是重写。

<a id="sec-16-1"></a>
### 16.1 结论

- `store.Store` 的**字节级 KV 形状是正确的**：JS 只需要存取 bytes，不需要理解 Character / Voice / YAML / 合成 hash。
- 但当前设计把“选择哪个 store”和“直接 OS IO”留在了 core/API 内部，外部 JS adapter 无法替换。
- 需要补的不是 web 实现，而是 **port 可见性、依赖注入、平台组合根、synth artifact store** 这几条 seam。
- 这些改动可以分阶段做；如果 web 是确定目标，建议在本次重构内完成前三条，避免后续回头改 service/API。

<a id="sec-16-2"></a>
### 16.2 已经可复用的部分

- `store.Store` 是 byte-level 的，JS adapter 只需要实现 KV/Blob 存储。
- `store.Memory` 证明了后端可替换；它也可以作为 web 的内存缓存基础。
- API DTO / 用例方法已经与 domain 解耦，web 不需要 import `core/internal/domain`。
- CLI 的 cobra/huh/lipgloss 已隔离，不会进入 wasm。
- `runtime` 已独立成 module，`synth.Engine` port 已存在；web 可以用 JS aq-wasm engine 替换 runtime。
- `go.work`、`justfile`、depguard、CI 的结构不需要因为 web 改变。

<a id="sec-16-3"></a>
### 16.3 当前缺口

1. **Port 不可见、不可注入**
   - `store` 在 `core/internal/store`，外部 `web/` module 不能 import 并实现它。
   - `character.Service.backend()` 会在内部 `store.NewFS(root)`。
   - 计划中的 `api.Open()` 也直接构造 FS store 和 `aq.Pool`。
   - 结果是：即使 JS 提供了 IndexedDB store，也没有入口把它交给 service。

2. **Store 覆盖范围不足**
   - `synth` 直接使用 `os.ReadFile` / `os.WriteFile` / `os.MkdirAll` / `os.Stat` 处理 WAV、缓存、`index.yaml`、输出路径。
   - `config.Service` 用 `os.MkdirAll` 做 `EnsurePaths`。
   - `domain.DefaultConfig` 用 `os.UserConfigDir` / `os.TempDir`。
   - `voice.LoadPhont` 用 `os.ReadFile`。
   - `runtime/aq` 用 `os.ReadFile` 加载 wasm。
   - JS 实现了 store，也只能覆盖 config/character 的 YAML，覆盖不到最重要的合成输出与缓存。

3. **同步 Store 与 IndexedDB 的异步模型冲突**
   - `Store` 是同步接口；IndexedDB 是异步的；`localStorage` 同步但容量小，不适合 WAV。
   - 如果 Go 业务代码同步调用 IDB，普通 wasm_exec 无法直接等待 Promise。
   - 必须选择：同步 port + 内存缓存 + flush，或异步 port + 全链路 async。

4. **路径语义过强**
   - `List(dir)`、`Config.Paths`、`SynthRequest.Output`、`SynthResult.Path` 都带文件系统语义。
   - IDB/localStorage 更适合**逻辑 key + namespace/prefix**。
   - 字节接口本身没问题，但需要把 key、prefix list、原子写、not-found 写进契约。

5. **API 组合根把 runtime 带进来**
   - 计划中的 `core/api` 同时包含 DTO/用例、`Open`（FS + `aq.NewPool`）和 `aqEngine` adapter（import runtime）。
   - web 即使想用 JS aq-wasm engine，只要 import 根 `core/api`，就可能把 runtime/wazero 一起编译进来。
   - 需要把“纯 API”和“CLI 平台组合根”拆开。

6. **Engine port 在 internal**
   - `synth.Engine` / `EngineParam` 在 `core/internal/service/synth`，外部 web engine adapter 无法 import。
   - 需要把 engine 契约移到 public port。

<a id="sec-16-4"></a>
### 16.4 必需的架构增量

#### 16.4.1 Public `core/port`

新增 public package：

```text
core/port/
  store.go      # Store, Reader, Writer, Lister, Deleter, ErrNotFound, Flusher
  engine.go     # Engine, EngineParam
```

- `core/internal/service` 依赖 `port`，而不是 `store` 的具体实现。
- `core/internal/store` 的 FS/Memory 实现 `port.Store`。
- `core/api` 使用 `port`。
- 未来 `web/` module import `core/port`，实现 JS store 和 JS engine。
- `port` 不暴露 domain 类型，只使用 `[]byte`、`string`、基本类型。

#### 16.4.2 `api.New(deps)` + 平台组合根

```go
type Deps struct {
    ConfigStore    port.Store
    CharacterStore port.Store
    ArtifactStore  port.Store // synth 输出/缓存；可选但 web synth 需要
    Engine         port.Engine
    // Defaults / platform info...
}

func New(deps Deps) (*App, error)
```

- `core/api` 只做纯用例、DTO、`New(deps)`；不 import `os`、`store.FS`、`runtime`。
- CLI 组合根移到 `core/api/cliopen`：
  - 创建 FS store；
  - 创建 `aq.Pool`；
  - 读取 config；
  - 调用 `api.New`；
  - 负责 `EnsurePaths`、`CleanIndex` 等平台副作用。
- CLI 仍然只 import `core/api`（允许子包 `core/api/cliopen`），不 import runtime/store。
- 未来 `core/api/webopen` 或独立 `web/` module：
  - 创建 JS store；
  - 创建 JS engine（或 runtime engine）；
  - 调用同一个 `api.New`。

#### 16.4.3 synth 的 IO port 化

- `synth.Service` 不再直接 `os.*`。
- WAV 输出、缓存、`index.yaml` 走 `ArtifactStore`（可以复用 `port.Store`）。
- `-o` 明确文件路径、`-p` 播放等 CLI 行为移到 CLI。
- API 始终返回 WAV bytes；web 直接交给 Web Audio。
- `-P` 的 YAML 文件读取移到 CLI；service 接受 bytes 或已解析的 voice。
- `UniquePath` 改成基于 `Store.Exists` 的逻辑 key 分配，不再依赖 `os.Stat`。

#### 16.4.4 OS 相关默认值移出 domain/core

- `DefaultCharacterPath` / `DefaultOutputPath` / `EnsurePaths` 移到 `cliopen` 或平台层。
- domain 保持纯模型；web 的默认值是 host 提供的逻辑命名空间。
- `Config.Paths` 可以保留给 CLI 编辑，但 service 不再用它构造 FS。

#### 16.4.5 phont 与 runtime loader

- `SetPhontPath` 拆成 `SetPhontData([]byte)`；CLI 读文件，web 直接传 File bytes。
- `runtime/aq` 增加 `NewWithWASM([]byte)` 或 `Loader` port；如果 web 用 JS aq-wasm engine，则不需要 runtime。
- `Engine` port 移到 `core/port`，web engine 直接实现。

<a id="sec-16-5"></a>
### 16.5 Store 契约

建议的 public port：

```go
type Store interface {
    Read(key string) ([]byte, error)
    Write(key string, data []byte) error
    Exists(key string) (bool, error)
    List(prefix string) ([]string, error)
    Delete(key string) error
}

var ErrNotFound = errors.New("store: not found")

// 可选：让 web adapter 在 Close 时把内存缓存刷回 IDB。
type Flusher interface {
    Flush(ctx context.Context) error
}
```

语义要求：

| 操作 | 契约 |
|---|---|
| `Read` | 返回完整 bytes；不存在返回 `ErrNotFound` |
| `Write` | 原子 replace；成功返回后读到的是完整新值 |
| `Exists` | 不存在返回 `false, nil`，不是错误 |
| `List(prefix)` | 返回该 prefix 下的 immediate 逻辑 key，按确定顺序 |
| `Delete` | 不存在返回 `ErrNotFound` |
| key | 逻辑 key，不是绝对路径；不含平台路径分隔符假设 |
| 并发 | 由 adapter 保证安全；core 不假设全局锁 |
| `Flush` | 可选；web adapter 用它在关闭/检查点持久化 |

JS 映射：

| Store | IndexedDB | localStorage |
|---|---|---|
| `Read` | object store get | `getItem` |
| `Write` | put（事务内） | `setItem`（二进制需编码） |
| `Exists` | count/get | key 存在性 |
| `List` | key range / prefix scan | key scan |
| `Delete` | delete | `removeItem` |
| 缺失 | `ErrNotFound` | `ErrNotFound` |

YAML、Character、Voice、hash 仍留在 Go/core；JS 不需要实现这些。

<a id="sec-16-6"></a>
### 16.6 同步还是异步

**方案 A：同步 Store + 内存缓存 + Flush（推荐）**

- 启动：web adapter 把 IDB 数据读入内存缓存。
- Go 业务代码同步读写缓存。
- 在 `App.Close` 或检查点调用 `Flush(ctx)` 写回 IDB。
- 优点：service 代码不变；与 CLI/FS 一致。
- 缺点：不是每次写都持久；崩溃可能丢最后一次未 flush 的变更。

**方案 B：异步 Store**

- `Read/Write/List` 全部带 `context.Context`，service 全链路 async。
- 优点：每次操作直接持久化。
- 缺点：改动大，所有 service/API 方法都要 async 化；CLI 也要适配。

**方案 C：JS 在调用 Go 前后自己存取**

- Go 侧只拿内存快照，glue 负责加载/保存。
- 优点：Go 不碰 IDB。
- 缺点：glue 最重，状态边界容易出错。

建议：如果没有“每次写必须立即持久”的硬要求，选 A；如果要直接 per-op IDB，选 B。不要为了 web 同时保留两套 store 语义。

`localStorage` 只适合小量配置/元数据；WAV 缓存应使用 IndexedDB 或交给 JS 的 Blob 存储。

<a id="sec-16-7"></a>
### 16.7 WASM 打包边界

- web build 只应 import：
  - `core/api`（纯用例/DTO）；
  - `core/port`（Store/Engine 契约）；
  - 可选：`runtime`（如果决定在 wasm 内用 Go runtime）；
  - 未来 `web/` module 的 JS adapter。
- 不应 import：
  - `cli`（cobra/huh/lipgloss）；
  - `core/internal/store` 的 FS；
  - `core/api/cliopen`（FS + runtime 组合根）。
- 如果 web 用 JS aq-wasm：
  - 实现 `port.Engine`；
  - 不需要 import runtime；
  - `synth` 通过 `port.Engine` 调用 JS。
- 如果 web 在 wasm 内继续用 runtime：
  - runtime 需要 `NewWithWASM` / loader port；
  - `api.New` 接受 runtime engine adapter。
- Go 只链接实际 import 的包；但若 `core/api` 根包自己 import runtime，web import 它就会带上 runtime。因此必须把 `cliopen`/`aqEngine` 从纯 API 包拆出去。
- `go.work` module 结构不需要变；未来增加 `web/` module（module path `yukumo-script/web`），depguard 规则如下：
  - 允许 `yukumo-script/core/api`、`yukumo-script/core/port`；
  - 禁止 `yukumo-script/core/internal`、`yukumo-script/cli`、`yukumo-script/runtime`（除非 web 明确使用 runtime）。

<a id="sec-16-8"></a>
### 16.8 优先级

如果 web 是确定目标，建议分三档：

**P0：现在就该做（否则以后要返工）**

1. public `core/port`（Store / Engine / EngineParam / ErrNotFound）。
2. service 依赖 port，不再内部 `NewFS`。
3. `api.New(deps)` 与 `api/cliopen` 拆分。
4. `synth.Engine` 移到 public port。

**P1：web synth 前必须做**

1. `synth` 的 WAV/缓存/index 走 `ArtifactStore`。
2. `config` / `domain` 的 OS 默认值与 `EnsurePaths` 移到平台层。
3. 明确 sync + Flush 还是 async Store。
4. `-o` / `-p` / `-P` 等 CLI 文件行为从 service 移出。

**P2：web 真正实现时再做**

1. JS store adapter（IndexedDB/localStorage）。
2. JS engine adapter（aq-wasm）。
3. Web Audio / UI。
4. Web 端 config/default 命名空间。

如果 web 还只是可能性：至少做 P0；P1 可以在 web 立项前补，但 `synth` 的 store 覆盖不足是最大的返工点。

<a id="sec-16-9"></a>
### 16.9 待定问题

1. web 是继续在 wasm 内用 Go `runtime`/wazero，还是直接用 JS 的 aq-wasm engine？
   - 决定 runtime 是否进入 web build，以及是否需要 `NewWithWASM`。
2. web 是否需要“每次写立即持久”，还是可以接受内存缓存 + `Flush`？
   - 决定 Store 同步/异步。
3. web 的 config `paths.*` 是有意义的逻辑 namespace，还是完全忽略？
   - 决定 Config DTO 是否保留 path 字段。
4. web 是否需要 `-P` / phont 文件路径，还是全部改成 bytes/file picker？
   - 决定 API 是否保留 path 版本。
5. config/character/artifact 用**一个** namespaced store，还是三个独立 store？
   - 决定 `Deps` 结构和 CLI 的 FS root 划分。

<a id="sec-16-10"></a>
### 16.10 对现有 plan 的改动点

| 现有章节 | web addendum 要求 |
|---|---|
| 1. 最终决策 | 增加 public `core/port`；API 拆 `New(deps)` / `cliopen` |
| 4.4 `core/api/go.mod` | 纯 API 不 import runtime；`cliopen` 子包 import runtime |
| 6.3 `App` / `Open` | `Open` 移到 `cliopen`；`App` 由 `New(deps)` 构造 |
| 6.4 DTO 转换 | 不变，DTO 仍不暴露 domain |
| 7. Engine port | `Engine` / `EngineParam` 移到 `core/port` |
| 9. 测试分层 | API 集成测试可继续用 runtime；core 测试用 fake/port |
| 10. depguard | 增加 `core/port` 规则；CLI 不得 import port；web 允许 port |
| 13. 迁移步骤 | P0 可并入 Phase 3/4；P1 作为独立 phase |
| 14. 验收清单 | 增加“无 os 依赖的纯 core/API”检查（web 目标时） |

---

<a id="sec-a"></a>
## 附录 A：文件映射

| 现在 | 之后 |
|---|---|
| `internal/domain` | `core/internal/domain` |
| `internal/store` | `core/internal/store` |
| `internal/service` | `core/internal/service` |
| `internal/runtime/aq` | `runtime/aq`（module `yukumo-script/runtime`） |
| `internal/cli` | `cli`（module `yukumo-script/cli`） |
| `main.go` | `cli/cmd/yukumo-script/main.go` |
| `test/runtime/aq` | `core/api/test/runtime/aq` |
| `test/service/synth/cache_test.go` | `core/api/test/synth/cache_test.go` |
| `test/service/synth/service_test.go` | `core/api/test/synth/service_test.go` |
| `test/service/synth/filename_test.go` | `core/test/service/synth/filename_test.go` |
| `test/service/synth/voice_test.go` | `core/test/service/synth/voice_test.go` |
| `test/service/character/**` | `core/test/service/character/**` |
| `test/service/config/**` | `core/test/service/config/**` |
| `test/service/nameutil/**` | `core/test/service/nameutil/**` |
| `test/cli/prompt` | `cli/test/prompt` |
| `internal/runtime/aq/errors_test.go` | `runtime/aq/errors_test.go` |
| `internal/cli/synth/run_test.go` | `cli/synth/run_test.go` |
| `go.mod` / `go.sum`（根） | 删除 |
| `third-party/` | 位置不变 |

<a id="sec-b"></a>
## 附录 B：命令速查

```bash
# 日常入口（推荐）
just build
just run
just test
just vet
just lint
just check-replaces
just ci

# 构建 / 运行（直接命令，workspace 模式）
go build -o yukumo-script ./cli/cmd/yukumo-script
go run ./cli/cmd/yukumo-script

# 模块独立检查（justfile 内部也是这些命令）
(cd core && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd runtime && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd core/api && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)
(cd cli && GOWORK=off go build ./... && GOWORK=off go vet ./... && GOWORK=off go test ./...)

# depguard + staticcheck
(cd core && GOWORK=off golangci-lint run)
(cd runtime && GOWORK=off golangci-lint run)
(cd core/api && GOWORK=off golangci-lint run)
(cd cli && GOWORK=off golangci-lint run)

# replace 一致性
go run scripts/check-replaces.go
```
