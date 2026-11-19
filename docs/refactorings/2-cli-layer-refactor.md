# CLI 层内部分层重构计划（可执行修订版）

- 状态：Implemented
- 日期：2026-11-16（2026-11-16 修订：依赖边界可执行版）
- 前置：[`1-layer-refactor.md`](./1-layer-refactor.md) 已完成（CLI 只依赖 `core/api`）
- 分支：`layer-refactor-cli`
- 范围：`cli` module **内部**把 cobra shell、domain action、interactive UI 拆开；不改 `core` / `runtime` / `core/api` 边界，不改 CLI 对外输出文本、交互流程、合成策略、文件格式（明确修正见 §15.1）。
- 不在范围：Web/WASM、新 GUI module、把 voice 迁出 CLI、改 `cli.default_interactive` 的配置语义。
- Serve 模式：stdin/stdout 帧协议评估见 [16. Serve（stdin/stdout）addendum](#sec-16)，不属于当前 Phase 0–7 的实现范围。
- 本次修订解决的问题：
  - `Streams` 不再放 `cli` 根包，避免「根包 → 子包 → 根包」循环。
  - domain schema 从 shell 包移到 `cli/schema`，避免 `ui` / `action` 反向依赖 shell 包。
  - `action` 只做 use case，不打印、不播放；新增 `cli/render` 与 `cli/play`。
  - Browse 的 MenuSpec 模型、cobra 适配器、UI 循环三者分离。
  - depguard 增加 exact-root 规则，并用 import-graph 脚本兜住传递依赖。

## 实施记录

- Phase 0：基线（本文件状态改为 In Progress，分支 `layer-refactor-cli`）
- `cc07202` docs(refactor): start cli layer refactor
- `c708dac` refactor(cli): introduce streams io contract
- `c558a29` refactor(cli): extract schema and parse packages
- `ac6e195` refactor(cli): extract action render and play layers
- `aec1a87` refactor(cli): extract interactive ui package
- `1bc49f6` refactor(cli): replace browse mutation with menu spec
- `f4e3708` refactor(cli): unify voice routing without options
- 收尾：`refactor(cli): add import guard and tests`（本文件状态更新为 Implemented）
- 最终检查：`just ci` 通过（4 module build/test/vet，depguard + staticcheck 0 issues，`scripts/check-cli-imports.go` 通过）。

## 目录

- [1. 最终决策](#sec-1)
- [2. 问题与动机](#sec-2)
- [3. 目标目录](#sec-3)
- [4. 层职责与依赖](#sec-4)
- [5. Streams 与 IO 契约](#sec-5)
- [6. Cobra shell 设计](#sec-6)
- [7. Action 层设计](#sec-7)
- [8. UI 层设计](#sec-8)
- [9. Browse 与 MenuSpec](#sec-9)
- [10. Voice 入口统一](#sec-10)
- [11. 编译期 / linter 边界](#sec-11)
- [12. 测试分层](#sec-12)
- [13. 迁移步骤](#sec-13)
- [14. 验收清单](#sec-14)
- [15. 风险与取舍](#sec-15)
- [16. Serve（stdin/stdout）addendum](#sec-16)
  - [16.1 结论](#sec-16-1)
  - [16.2 三种运输层](#sec-16-2)
  - [16.3 落点](#sec-16-3)
  - [16.4 对 action / Streams 的约束](#sec-16-4)
  - [16.5 协议草图](#sec-16-5)
  - [16.6 与 browse / interactive 的隔离](#sec-16-6)
  - [16.7 待定问题](#sec-16-7)
  - [16.8 对本文 Phase 的影响](#sec-16-8)
- [附录 A：文件映射](#sec-a)
- [附录 B：典型 RunE 前后对比](#sec-b)

---

<a id="sec-1"></a>
## 1. 最终决策

| 决策 | 结果 |
|---|---|
| 模块边界 | **不新增 go module**；仍是 `yukumo-script/cli`，只做 package 内部分层 |
| 主依赖方向 | `shell`（cobra）→ `ui` → `action` → `core/api`；shell 也可直接调 action |
| IO 契约 | `cli/streams.Streams`（`In` / `Out` / `Err` / `Ctx`）；**action 不接收 Streams**，只接收 `context.Context` 与领域参数 |
| 输出 | `cli/render` 负责文本/YAML 渲染；`cli/play` 负责系统播放器副作用；action 不打印、不播放 |
| Domain schema | `cli/schema`（纯数据 + 校验）；`cli/parse`（纯解析）；`setting` 只保留 huh-facing 的 `Ask` / `BrowseKeys` / `Metas` |
| Cobra 适配 | `cli/internal/cobraio`（`FromCmd`）；`cli/internal/cobraui`（cobra 树 → MenuSpec） |
| Browse | **停止** `attachBrowse` 改写父命令 `RunE`；root 与 `config` 显式进入 `ui.Browse` |
| Voice | schema 归 `cli/schema`；`voice.Dispatch` 只做 argv 路由；菜单/表单在 `ui`；持久化在 `action` |
| Interactive 开关 | 仍来自 `cfg.CLI.DefaultInteractive`；只在 shell 里分支，不把 flag 渗进 action |
| 行为 | Phase 内行为与输出不变；允许内部函数签名与包路径变化。§15.1 列出的「非交互误入 huh」作为明确修正 |
| 护栏 | depguard：exact root 禁止、action/ui/schema/parse/streams/menu/render/play 分层规则；另加 import-graph 脚本进 `just ci` |

目标依赖（CLI 内部）：

```text
cli/cmd/yukumo-script (main)
  ↓
cli (root: Execute / newRoot / browse bootstrap)
  ↓
cli/{config,character,synth}            ← cobra shell，不 import cli 根包
  ├──► cli/ui
  ├──► cli/action
  ├──► cli/schema
  ├──► cli/render
  ├──► cli/play                         ← 仅 synth / ui 使用
  ├──► cli/streams
  └──► cli/internal/{cobraio,cobraui}   ← cobra 适配器

cli/ui ──► cli/action ──► core/api
   │           │
   │           └──► cli/schema ──► cli/parse
   ├──► cli/menu ──► cli/streams
   ├──► cli/schema
   ├──► cli/render ──► core/api
   ├──► cli/play ──► core/api
   ├──► cli/prompt ──► huh / lipgloss / cli/parse
   ├──► cli/setting ──► cli/prompt + cli/schema
   └──► cli/streams
```

禁止方向：

```text
任意 cli 子包        → yukumo-script/cli            （exact 根包；depguard + import graph）
ui                   → cobra / config / character / synth
action               → cobra / huh / lipgloss / prompt / setting / ui / render / play / streams
                       / menu / config / character / synth
schema / parse / streams → prompt / setting / ui / action / render / play / menu / shell / cobra / huh
menu                 → cobra / ui / action / schema / setting / prompt / shell
render / play        → cobra / huh / prompt / setting / ui / action / shell
prompt               → cobra / api / schema / setting / ui / action / shell
setting              → cobra / api / action / ui / shell
cmd                  → 除 cli 根包外，不直接碰 CLI 子包或 core/api
```

> `ui` 允许 import `prompt` / `setting` / `schema` / `action` / `render` / `play` / `menu` / `streams` / `core/api`。  
> `action` 允许 import `core/api` / `schema` / 标准库（`context` / `fmt` / `time` 等）。  
> shell（各 `*_cmd.go` / `NewCmd`）允许 import cobra、ui、action、schema、render、play、streams、`internal/cobraio`、`internal/cobraui`。  
> **shell 不得 import `yukumo-script/cli` 根包。**

---

<a id="sec-2"></a>
## 2. 问题与动机

`1-layer-refactor` 解决了「CLI 相对 core/runtime 的边界」。当前 `cli` **内部**仍把两种 UI 叠在同一套函数上：

1. **Cobra shell**：argv、flags、help、子命令树
2. **huh / prompt UI**：Select、输入、面包屑、嵌套菜单

具体耦合：

| 现象 | 位置 | 后果 |
|---|---|---|
| 全局 `Options.Interactive` 在建树时改 `Args` / `Hidden` / 是否注入 browse | `root.go`、`character/root.go`、`config/root.go` | 测两种模式要建两棵树；browse 与 argv 语义纠缠 |
| `attachBrowse` 给无 `RunE` 的父命令挂菜单 | `browse.go` | 父命令行为被 interactive 改写 |
| browse 以 `RunE(cmd, nil)` 调叶子 | `browse.go` `runCommand` | 叶子必须能处理空 args；one-shot 与菜单共享同一函数 |
| helpers 传 `*cobra.Command` 只为 IO | `applyEdit`、`runRequest`、`editParam*`、menus | action/ui 无法脱离 cobra 测试 |
| `fields.go` / `pick.go` 留在 shell 包 | `config/fields.go`、`character/fields.go`、`character/pick.go`、`synth/pick.go`、`character/voice/fields.go` | 一旦 ui 需要这些 helper，就与 shell → ui 形成 import cycle |
| voice 不在 cobra 子树 | `character/root.go` + `voice.Dispatch` / `RunMenu` | browse 发现不了 voice；两套导航并行 |
| 同一 `run*` 双模式分支 | `character/edit.go`、`config/edit.go`、`voice/*.go` | validation / prompt / persist 缠在一起 |
| action 同时做 use case / 打印 / 播放 | `synth/run.go`、`character/applyEdit` | serve 想要「只返回 bytes」时被迫加 quiet 开关；SRP 模糊 |

`prompt` 与 `setting` 已经示范了正确边界（`setting.Ask(..., io.Writer)`）。本次把同一模式推到 config / character / synth / voice，并把「共享 schema」「输出」「播放」从 action/shell 中拆出去。

---

<a id="sec-3"></a>
## 3. 目标目录

```text
cli/
├── go.mod
├── go.sum
├── .golangci.yml
├── cmd/
│   └── yukumo-script/
│       └── main.go                 # 不变
├── root.go                         # Execute / newRoot
├── browse.go                       # enableInteractiveRoot；不再 walk/attachBrowse
├── streams/
│   └── streams.go                  # Streams + In/Out/Err/Ctx；无 cobra
├── parse/
│   └── parse.go                    # ExpandPath / ParseBool / ParseRangedInt / ClampInt
├── schema/
│   ├── field.go                    # Kind / Field / Entry[T] / Keys / Lookup / Normalize
│   ├── config.go                   # ConfigFields()
│   ├── character.go                # CharacterFields()
│   └── voice.go                    # VoiceParamFields(version)
├── setting/
│   ├── ask.go                      # Ask / BrowseKeys（huh-facing）
│   └── meta.go                     # Metas（Entry → prompt.LeafMeta）
├── prompt/                         # huh 原语 + crumbs；纯解析委托 cli/parse
│   ├── bool.go
│   ├── key.go
│   ├── nav.go
│   ├── number.go
│   ├── path.go
│   ├── prompt.go
│   ├── run.go
│   ├── select.go
│   └── tree.go
├── internal/
│   ├── cobraio/
│   │   └── fromcmd.go              # FromCmd(*cobra.Command) streams.Streams
│   └── cobraui/
│       └── menu.go                 # FromCobra(*cobra.Command) menu.Node
├── menu/
│   └── menu.go                     # Node{Name,Short,Children,Run(streams.Streams)}
├── action/                         # 无 cobra / huh / IO / play
│   ├── config.go
│   ├── character.go
│   ├── voice.go
│   └── synth.go
├── render/                         # 文本 / YAML / 合成结果行
│   ├── yaml.go
│   ├── line.go
│   └── synth.go
├── play/
│   └── play.go                     # 系统播放器副作用（原 synth/preview.go）
├── ui/                             # interactive 循环 / 表单 / pick
│   ├── browse.go
│   ├── config_edit.go
│   ├── character_edit.go
│   ├── character_pick.go
│   ├── character_add.go
│   ├── character_delete.go
│   ├── voice_menu.go
│   ├── voice_pick.go
│   ├── voice_add.go
│   ├── voice_edit.go
│   ├── voice_delete.go
│   ├── voice_duplicate.go
│   ├── voice_param.go
│   ├── synth_pick.go
│   └── synth_form.go
├── config/                         # 只留 cobra shell
│   ├── root.go
│   ├── list.go
│   ├── get.go
│   └── edit.go
├── character/
│   ├── root.go
│   ├── list.go
│   ├── get.go
│   ├── add.go
│   ├── edit.go
│   ├── delete.go
│   └── voice/
│       └── dispatch.go             # 薄 argv 路由；无 cobra / 无 Options
├── synth/
│   ├── root.go                     # newSynthCmd + runSynth
│   └── voice.go                    # newVoiceCmd（交互表单入口）
└── test/
    ├── prompt/                     # 既有
    ├── schema/                     # 新增：Normalize / Lookup
    ├── action/                     # 新增：temp api.App + 返回值断言
    ├── render/                     # 新增：formatResultLine 等
    ├── shell/                      # 新增：Args / Hidden / 路由元数据
    └── arch/                       # 可选：import graph（与 scripts 二选一）
```

说明：

- **不强制**把 `config/` 等拆成 `config/cmd` 子目录；同包内 `*_cmd` 职责清晰 + depguard 按目录禁 import 即可。
- `fields.go` **不再留在** `config/` / `character/` / `character/voice/`；统一进 `schema`。
- `pick.go` / `voice_form.go` / `preview.go` 不再留在 shell 包；分别进 `ui` / `play`。
- `synth/run.go` 与 `synth/run_test.go` 删除；格式化进 `render/synth.go`，use case 进 `action/synth.go`。
- `browse.go` 只保留 `enableInteractiveRoot`；菜单循环在 `ui/browse.go`，MenuSpec 模型在 `menu`，cobra 适配在 `internal/cobraui`。

---

<a id="sec-4"></a>
## 4. 层职责与依赖

| 层 / 包 | 职责 | 允许依赖 | 禁止 |
|---|---|---|---|
| `cmd/` | `main` | `cli` + stdlib | cobra、huh、api、runtime |
| `cli` root | `Execute` / `newRoot` / root browse 入口 | 子包 shell、ui、menu、streams、internal 适配器、prompt、api、cobra | 被任意子包 import |
| shell（`config` / `character` / `synth`） | `Use` / `Args` / `Flags`；解析 argv；选 interactive 或 one-shot；构造 Streams | cobra、ui、action、schema、render、play、streams、internal 适配器、api | 长循环、schema 定义、pick 组件、业务细节复制、import 根包 |
| `ui/` | 收集输入、菜单循环、面包屑、pick、表单；完成后调 action | streams、menu、schema、setting、prompt、action、render、play、api、lipgloss | cobra、config/character/synth shell、import 根包 |
| `action/` | 调用 `api.App` 的 use case；返回领域结果 + 可选警告；不打印、不播放 | api、schema、stdlib（context/fmt/time） | cobra、huh、lipgloss、prompt、setting、ui、render、play、streams、menu、shell、根包 |
| `schema/` | `Field` / `Entry` / `Keys` / `Lookup` / `Normalize` 与 config/character/voice 字段定义 | api、parse、stdlib | cobra、huh、lipgloss、prompt、setting、ui、action、shell、根包 |
| `parse/` | 纯解析：路径展开、bool、range int | 仅 stdlib | 任何项目包 |
| `streams/` | IO 契约 | 仅 stdlib | 任何项目包 |
| `setting/` | huh-facing 的 `Ask` / `BrowseKeys` / `Metas` | schema、prompt、stdlib | cobra、api、action、ui、shell、根包 |
| `prompt/` | huh 封装、crumbs、tree；纯解析委托 `parse` | huh、lipgloss、parse、stdlib | cobra、api、schema、setting、action、ui、shell、根包 |
| `menu/` | MenuSpec 模型 | streams、stdlib | cobra、ui、action、schema、setting、prompt、shell、根包 |
| `render/` | 文本 / YAML / 合成结果行 | api、yaml、stdlib | cobra、huh、prompt、setting、ui、action、play、shell、根包 |
| `play/` | 系统播放器副作用 | api、stdlib | cobra、huh、prompt、setting、ui、action、render、shell、根包 |
| `internal/cobraio` | `FromCmd` | cobra、streams | ui、action、schema、shell |
| `internal/cobraui` | cobra 树 → `menu.Node` | cobra、menu、streams | ui、action、schema、shell |

`Options.Interactive` **只存在于 shell**。`action` 不知道 interactive；`ui` 默认就是 interactive。

依赖图：

```text
shell ──► ui ──► action ──► core/api
  │         │
  │         ├──► render ──► core/api
  │         ├──► play ──► core/api
  │         ├──► menu ──► streams
  │         ├──► schema ──► parse
  │         └──► prompt / setting
  └──► action / schema / render / play
```

---

<a id="sec-5"></a>
## 5. Streams 与 IO 契约

### 5.1 类型：叶子包，无 cobra

```go
// cli/streams/streams.go
package streams

import (
	"context"
	"io"
)

type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
	Ctx context.Context
}
```

`Streams` 放在 `cli/streams`，**不是** `cli` 根包。否则 shell 子包为了 `cli.Streams` 必须 import 根包，而根包又 import 子包，形成循环；`ui` / `action` 也会被卷入。

### 5.2 `FromCmd` 放 internal 适配器

```go
// cli/internal/cobraio/fromcmd.go
package cobraio

import (
	"context"
	"io"
	"strings"

	"yukumo-script/cli/streams"

	"github.com/spf13/cobra"
)

func FromCmd(cmd *cobra.Command) streams.Streams {
	if cmd == nil {
		return streams.Streams{
			In:  strings.NewReader(""),
			Out: io.Discard,
			Err: io.Discard,
			Ctx: context.Background(),
		}
	}
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	return streams.Streams{
		In:  cmd.InOrStdin(),
		Out: cmd.OutOrStdout(),
		Err: cmd.ErrOrStderr(),
		Ctx: ctx,
	}
}
```

- shell 每个 `RunE` 第一行：`s := cobraio.FromCmd(cmd)`。
- `ui` / `render` / `play` 不 import cobra。
- **不要**在 browse leaf 里用 `FromCmd(nil)`；`nil` 只用于测试/缺省，`Out` / `Err` 是 `io.Discard`，会把 browse 输出吞掉。

### 5.3 替换规则

| 现在 | 之后 |
|---|---|
| `cmd.OutOrStdout()` | `s.Out` |
| `cmd.ErrOrStderr()` | `s.Err` |
| `cmd.Context()` / `cmdContext(cmd)` | `s.Ctx` |
| `func applyEdit(cmd *cobra.Command, …)` | `action.ApplyConfigEdit(ctx, app, …)` + `render` |
| `func runRequest(cmd *cobra.Command, …)` | `action.Synthesize(ctx, app, req)` + `render` + `play` |
| menu 里 `fmt.Fprintln(cmd.ErrOrStderr(), err)` | `fmt.Fprintln(s.Err, err)` 或 `render.Line(s.Err, …)` |

### 5.4 与 `setting` / `render` / `play` 对齐

- `setting.Ask(..., s.Err)`：不变，只把 `cmd` 换成 `s.Err`。
- `render.*` 接收 `io.Writer`，不接收 `Streams`，便于单测。
- `play.Preview(ctx, res, async)` 接收 `context.Context`，不接收 `Streams`。
- `action` 不接收 `Streams`；需要取消的 API 调 `ctx`（目前 `Synthesize` / `AddVoice` / `UpdateVoice` / `DuplicateVoice` / `PresetNames` / `HasPreset`）。

---

<a id="sec-6"></a>
## 6. Cobra shell 设计

### 6.1 原则

每个命令文件只做四件事：

1. 声明 `Use` / `Short` / `Args` / `Flags`
2. `RunE` 第一行 `s := cobraio.FromCmd(cmd)`
3. 选择 `ui.*` 或 `action.*`
4. 需要输出时调用 `render.*` / `play.*`

`Args` 策略（按命令，不再用一条通用规则）：

| 命令 | 非 interactive | interactive |
|---|---|---|
| `config edit [key] [value]` | `ExactArgs(2)` | `MaximumNArgs(2)`；`len<2` → `ui.EditConfig`；`len=2` → one-shot |
| `character edit [name] [field] [value]` | `ExactArgs(3)` | `MaximumNArgs(3)`；`len<3` → `ui.EditCharacter`；`len=3` → one-shot |
| `character add [name] [description]` | `RangeArgs(1,2)` | `MaximumNArgs(2)`；始终走 `ui.AddCharacter`（已有值不重复问） |
| `character delete [name]` | `ExactArgs(1)` | `MaximumNArgs(1)`；始终走 `ui.DeleteCharacter`（即使给了 name 也确认） |
| `character get [name]` | `ExactArgs(1)` | `MaximumNArgs(1)`；`Hidden=true`；缺 name → `ui.PickCharacter` |
| `config get <key>` | `ExactArgs(1)` | `Hidden=true`；显式调用仍走 one-shot |
| `synth` | `ArbitraryArgs`；缺 text/voice 报错 | `ArbitraryArgs`；缺 text → `ui.EnsureSynthText`；缺 voice → `ui.PickVoiceRef` |
| `synth voice` | 报错（interactive only） | `ui.SynthVoiceForm` |
| `character <name> voice …` | argv dispatch；缺 verb 报错 | `ui.VoiceMenu`；有 verb 走对应 one-shot/ui 分支 |

**关键例外**：不能把「参数足够就走 one-shot」当成所有命令的规则。`character delete` 在 interactive 下总是确认；`voice edit` 在 interactive 下总是进菜单；这些要按上表/§14.4 保留。

### 6.2 `Execute` / `newRoot`

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
		enableInteractiveRoot(root)
	}
	return root.ExecuteContext(context.Background())
}

func newRoot(app *api.App, interactive bool) *cobra.Command {
	root := &cobra.Command{
		Use:           "yukumo-script",
		Short:         "Yukumo script and voice tooling",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(cliconfig.NewCmd(app, cliconfig.Options{Interactive: interactive}))
	root.AddCommand(clicharacter.NewCmd(app, clicharacter.Options{Interactive: interactive}))
	root.AddCommand(clisynth.NewCmd(app, clisynth.Options{Interactive: interactive}))
	return root
}
```

`enableInteractiveRoot` 只给 root 设 `RunE`，不再 walk：

```go
// cli/browse.go
func enableInteractiveRoot(root *cobra.Command) {
	spec := cobraui.FromCobra(root)
	root.RunE = func(cmd *cobra.Command, _ []string) error {
		return ui.Browse(cobraio.FromCmd(cmd), spec)
	}
}
```

### 6.3 Group root 的显式入口

- `config`：interactive 时 `RunE` → `ui.Browse(s, cobraui.FromCobra(cmd))`；非 interactive 保持无 `RunE`（cobra help）。
- `character`：无 args → `cmd.Help()`（与现状一致，不改成 browse）；有 `<name> voice …` → `voice.Dispatch`。
- `synth`：无 args → `runSynth`（interactive 时按字段补 prompt；非 interactive 缺 text/voice 报错）。
- root：interactive 无 args → browse；非 interactive 保持 help。

### 6.4 禁止事项

- shell 里写超过 ~20 行的菜单循环（应在 `ui`）。
- shell 里直接 `app.Update*` + 打印（应在 `action` + `render`）。
- 向 ui/action 传递 `*cobra.Command`。
- 子包 import `yukumo-script/cli` 根包。

---

<a id="sec-7"></a>
## 7. Action 层设计

### 7.1 职责

`action` 是 use case 层：把已经凑齐的参数交给 `api.App`，返回领域结果；**不打印、不播放、不做 huh、不接收 Streams**。需要取消的 API 调用接收 `context.Context`。

示例：

```go
package action

type EditResult struct {
	Value   string // 用于渲染 "key: value"
	Warning string // 例如 clamp 提示；无警告时为空
}

func ListConfig(ctx context.Context, app *api.App) (api.Config, error)
func GetConfigField(ctx context.Context, app *api.App, key string) (string, error)
func ApplyConfigEdit(ctx context.Context, app *api.App, key, raw string, normalize bool) (EditResult, error)

func ListCharacters(ctx context.Context, app *api.App) ([]api.Character, error)
func GetCharacter(ctx context.Context, app *api.App, name string) (api.Character, error)
func AddCharacter(ctx context.Context, app *api.App, name, description string) (api.Character, error)
func ApplyCharacterEdit(ctx context.Context, app *api.App, oldName, field, raw string, normalize bool) (api.Character, EditResult, error)
func DeleteCharacter(ctx context.Context, app *api.App, name string) error

func ListVoices(ctx context.Context, app *api.App, charName string) ([]api.Voice, error)
func GetVoice(ctx context.Context, app *api.App, charName, name string) (api.Voice, error)
func AddVoice(ctx context.Context, app *api.App, charName string, v api.Voice) (api.Voice, error)
func SaveVoice(ctx context.Context, app *api.App, charName, oldName string, v api.Voice) (api.Voice, error)
func ApplyVoiceName(ctx context.Context, app *api.App, charName, oldName, newName string) (api.Voice, EditResult, error)
func ApplyVoiceParam(ctx context.Context, app *api.App, charName, voiceName string, version int, field, raw string, normalize bool) (api.Voice, EditResult, error)
func DeleteVoice(ctx context.Context, app *api.App, charName, name string) error
func DuplicateVoice(ctx context.Context, app *api.App, charName, src, newName string) (string, error)

func Synthesize(ctx context.Context, app *api.App, req api.SynthRequest) (api.SynthResult, time.Duration, error)
```

- `action` 使用 `schema.ConfigFields()` / `schema.CharacterFields()` / `schema.VoiceParamFields(version)` 做 `Lookup` / `Normalize`。
- `EditResult.Warning` 由 shell/ui 写 `s.Err`；action 不持有 writer。
- `Synthesize` **只**调用 `app.Synthesize` 并计时，不打印、不 preview。
- `SaveVoice` / `ApplyVoice*` 的 `ctx` 传给 `app.UpdateVoice` / `app.AddVoice` / `app.DuplicateVoice`。
- `ApplyVoiceParam` 保持现状：`phont_data` 走 `api.SetPhontPath`，其余字段 `entry.Set` 后调用 `p.Normalize()`。

### 7.2 从现状迁移

| action 函数 | 来源现状 |
|---|---|
| `ListConfig` / `GetConfigField` / `ApplyConfigEdit` | `config/list.go` / `get.go` / `edit.go` 的持久化部分 |
| `ListCharacters` / `GetCharacter` / `AddCharacter` / `ApplyCharacterEdit` / `DeleteCharacter` | `character/*` |
| `ListVoices` / `GetVoice` / `AddVoice` / `SaveVoice` / `ApplyVoiceName` / `ApplyVoiceParam` / `DeleteVoice` / `DuplicateVoice` | `character/voice/*` |
| `Synthesize` | `synth/run.go` `runRequest` 的 api + 计时部分 |

### 7.3 Preview / 播放的归属

`preview.go`（系统播放器）是 OS 副作用，不是 huh，也不是 cobra。

决策：**移入 `cli/play`**，由 shell / ui 在 action 返回后调用。

- one-shot：`action.Synthesize` → `render.SynthResult` → 若 `req.Preview` 则 `play.Preview(ctx, res, false)`。
- interactive：`action.Synthesize` → 更新状态（`render.SynthStatusText`）→ 若 preview 则 `play.Preview(ctx, res, true)`。
- serve：**不调用** render/play；只拿 `res.WAV` / `res.Path`。

这样 action 不需要 `Quiet` 开关，serve 也不会因为 action 默认打印而踩坑。

### 7.4 输出 / YAML / 警告

`cli/render`：

```go
package render

func Line(w io.Writer, format string, args ...any) error
func YAML(w io.Writer, v any) error
func Warning(w io.Writer, s string) error
func SynthResultText(res api.SynthResult, elapsed time.Duration) string
func SynthStatusText(res api.SynthResult, elapsed time.Duration) string
func SynthResult(w io.Writer, res api.SynthResult, elapsed time.Duration) error
```

- `SynthResultText` 保持 one-shot 的 `formatResultLine`；`SynthStatusText` 保持 interactive 的 `doSynthesize` 状态文案（`wrote` / `cached` / `not written`）。

- YAML 编码（`enc.SetIndent(2)`）从 `config/list.go`、`character/list.go`、`character/voice/list.go`、`character/voice/get.go` 移到 `render.YAML`。
- `formatResultLine` / `formatElapsed` 从 `synth/run.go` 移到 `render/synth.go`；`synth/run_test.go` 改为 `test/render/synth_test.go`。
- 警告文本保持现状：`clamped %s → %d (range [%d, %d])`。

### 7.5 不放进 action 的东西

- `prompt.*` / `setting.Ask` / `setting.BrowseKeys`
- 字段菜单循环、voice 菜单、synth voice form
- cobra `Args` 校验文案（可返回 `fmt.Errorf`，但校验时机在 shell）
- `io.Writer` / `Streams` / `lipgloss`

### 7.6 测试

- action 测试使用 `api.Open(api.Options{ConfigDir: t.TempDir()})` 得到真实 App；`api.App` 是具体 struct，不能 fake。
- 纯格式化在 `render` 测试；action 的 synth 若需要 wasm，放到集成测试，不做假 App。
- 如果需要无 wasm 的 action 单测，可在 action 定义 consumer interface；本计划不强制，优先 temp App。

---

<a id="sec-8"></a>
## 8. UI 层设计

### 8.1 迁入内容（完整清单）

| 现函数 / 文件 | 目标 |
|---|---|
| `config.runInteractiveEdit` | `ui.EditConfig` |
| `character.runInteractiveEdit` | `ui.EditCharacter` |
| `character/add.go` 的 prompt 分支 | `ui.AddCharacter` |
| `character/delete.go` 的 confirm / pick | `ui.DeleteCharacter` |
| `character/get.go` / `pick.go` 的 pick | `ui.PickCharacter` |
| `voice.RunMenu` | `ui.VoiceMenu` |
| `voice/add.go` 的 interactive 分支 | `ui.AddVoice` |
| `voice/edit.go` 的 interactive 分支 | `ui.EditVoice` |
| `voice/delete.go` 的 interactive 分支 | `ui.DeleteVoice` |
| `voice/duplicate.go` 的 interactive 分支 | `ui.DuplicateVoice` |
| `voice/get.go` / `pick.go` 的 pick | `ui.PickVoice` |
| `voice/param.go` 的交互部分 | `ui.EditVoiceParam` |
| `synth/runInteractiveVoice` | `ui.SynthVoiceForm` |
| `synth/runSynth` 的 `EnsureString` / `ensureVoiceRef` | `ui.EnsureSynthText` / `ui.PickVoiceRef` |
| browse 循环 | `ui.Browse`（消费 `menu.Node`） |

`ui` **不 import** `config` / `character` / `synth`；它只依赖 `schema` / `setting` / `prompt` / `action` / `render` / `play` / `menu` / `streams` / `api`。

### 8.2 签名约定

```go
func Browse(s streams.Streams, spec menu.Node) error

func EditConfig(s streams.Streams, app *api.App, partialKey, partialValue string) error
func EditCharacter(s streams.Streams, app *api.App, partialName, partialField, partialValue string) error
func AddCharacter(s streams.Streams, app *api.App, name, description string) error
func DeleteCharacter(s streams.Streams, app *api.App, name string) error
func PickCharacter(s streams.Streams, app *api.App, name string) (string, error)

func VoiceMenu(s streams.Streams, app *api.App, charName string) error
func PickVoice(s streams.Streams, app *api.App, charName, name string) (string, error)
func AddVoice(s streams.Streams, app *api.App, charName, name string, version int) error
func EditVoice(s streams.Streams, app *api.App, charName, name string) error
func DeleteVoice(s streams.Streams, app *api.App, charName, name string) error
func DuplicateVoice(s streams.Streams, app *api.App, charName, src, newName string) error
func EditVoiceParam(s streams.Streams, app *api.App, version int, p *api.SynthParam) error

func EnsureSynthText(s streams.Streams, text *string) error
func PickVoiceRef(s streams.Streams, app *api.App, ref *string) error
func SynthVoiceForm(s streams.Streams, app *api.App) error
```

- 部分填好的 argv 作为 `partial*` 传入，避免 ui 解析 cobra args。
- `AddVoice` 的 `version == 0` 表示「未提供，让用户在 ui 里选」；`1` / `2` / `10` 表示已由 argv 指定。
- `DuplicateVoice` 的 `newName` 为空时保持现状：interactive 提示 `name + " - copy"` 作为默认值；one-shot 交给 `api.DuplicateVoice` 处理默认名。
- `EditVoiceParam` 的 `p == nil` 表示新建 `api.DefaultSynthParam(version)`；非 nil 表示编辑现有 draft。
- `ui` 需要字段 schema 时 import `schema`，不 import domain shell 包。
- pick 的展示样式（dim description）在 `ui` 内实现；原 `character/pick.go`、`synth/pick.go`、`voice/pick.go` 删除。
- 从 `synth/voice_form.go` 一并迁入的 helper：`voiceTree` / `formatFactor` / `formatPreset` / `paramMetas` / `paramByVersion` / `textDisplay` / `editOutput` 等；从 `voice/param.go` 迁入 `pickPreset` / `pickBase` / `pickPhont`；从 `voice/version.go` 迁入 `SelectVersion`；从 `voice/duplicate.go` 迁入 `defaultCopyName`；从 `character/edit.go` 迁入 `leafLabel`。
- `EditVoiceParam` 在内存中改 `*api.SynthParam`，保存由调用方 `action.SaveVoice` / `action.AddVoice` 完成；与现行为一致。
- 错误：可恢复的 abort（`prompt.ErrAborted` / huh 的 user abort）在 ui 内吞掉或转 `nil`；真正的失败上抛。
- 成功路径一律 `action.*`；`ui` 允许直接 `app.Config()` / `app.GetCharacter()` / `app.GetVoice()` 做只读展示，但写操作必须走 action（避免打印/持久化分叉）。
- 所有 `app.*` 的 `context.Context` 参数统一传 `s.Ctx`；不再出现 `context.Background()`（`Execute` / `api.Open` 的根 context 除外）。

### 8.3 Crumb 所有权

- browse 进入叶子前 `prompt.Reset(path...)`
- ui 内部用 `prompt.With` / `Replace`
- 叶子返回后 browse 再 `Reset` 回父路径（保持现行为）
- `prompt` 的 crumbs 仍是包级状态；本计划不引入并发 UI，单线程假设写进 `prompt/nav.go` 注释。若要并发，后续把 crumbs 改成 Navigator。

---

<a id="sec-9"></a>
## 9. Browse 与 MenuSpec

### 9.1 现状问题

```go
// browse.go（现状）
func attachBrowse(cmd *cobra.Command) {
	start := cmd
	cmd.RunE = func(*cobra.Command, []string) error {
		return browseCommands(start)
	}
}
```

这会：

- 改写父命令在 interactive 配置下的 argv 行为；
- 让「菜单是否可进」依赖「有没有 RunE」；
- 与 voice 的手写菜单无法统一；
- 让 browse 逻辑与 cobra 强绑定。

### 9.2 目标：menu 模型 / cobra 适配 / ui 循环分离

```go
// cli/menu/menu.go
package menu

type Node struct {
	Name     string
	Short    string
	Children []Node
	Run      func(s streams.Streams) error // nil = 目录
}
```

```go
// cli/internal/cobraui/menu.go
func FromCobra(cmd *cobra.Command) menu.Node
```

`FromCobra` 只读生成 MenuSpec：

- 过滤 `!c.IsAvailableCommand()`、`c.Hidden`、`help`、`completion`、`serve`；
- `Short` 原样保留；
- 叶子 `Run` 用 browse 传入的 `s` 执行，**不能用 `FromCmd(nil)`**：

```go
func runLeaf(cmd *cobra.Command, s streams.Streams) error {
	leaf := &cobra.Command{Use: cmd.Name()}
	leaf.SetIn(s.In)
	leaf.SetOut(s.Out)
	leaf.SetErr(s.Err)
	leaf.SetContext(s.Ctx)
	switch {
	case cmd.RunE != nil:
		return cmd.RunE(leaf, nil)
	case cmd.Run != nil:
		cmd.Run(leaf, nil)
		return nil
	default:
		return fmt.Errorf("command %q is not runnable", cmd.Name())
	}
}
```

> 约定：browse 叶子必须是无 flag 的 interactive/one-shot 入口。`synth` 是目录（有 `voice` / `custom` 子命令），不会被当作叶子；`synth voice` 没有 flag。若未来某个叶子需要 flag，必须在 shell 里显式注册 `menu.Node`，不要依赖适配器。

### 9.3 启动方式（R1）

- **root**：interactive 且无 args → `ui.Browse(rootSpec)`；非 interactive 保持 help。
- **config**：interactive 时显式 `RunE` → `ui.Browse(configSpec)`；非 interactive 保持 help。
- **character**：interactive 无 args → help（与现状一致，不改成 browse）。
- **synth**：interactive 无 args → `runSynth`（与现状一致）。
- **不再** `walk` + 给每个 group `attachBrowse`。

这样 root / config 的菜单与现状一致，character / synth 的无参行为也与现状一致；不需要“统一成全部 browse”的 UX 变化。

### 9.4 菜单路径与面包屑

- `ui.Browse` 用 `menu.Node` 栈维护路径；顶层显示 `Exit`，子层显示 `Back`。
- 叶子调用前 `prompt.Reset(path + leaf.Name())`；返回后 `Reset` 回父路径。
- `cobraui` 不处理 crumbs；`ui` 独占。

### 9.5 Voice 的衔接

- `menu.Node` 的 character 编辑入口仍由 `ui.EditCharacter` 内部提供 `voices/` 选项（与现状一致）。
- 若后续要让 browse 直接发现 `voice`，在 `character` 节点注册一个显式的 `voices/` 叶子：先 `ui.PickCharacter` 选角色，再进 `ui.VoiceMenu`（或让 `ui.VoiceMenu` 接受空 `charName` 并自己 pick）；这属于 §10 的可选增强，不阻塞主路径。

---

<a id="sec-10"></a>
## 10. Voice 入口统一

### 10.1 现状

```text
argv:  character <name> voice <verb> …   → voice.Dispatch
TUI:   character edit → voices/          → voice.RunMenu
browse: （看不到 voice）
```

`Dispatch` / `RunMenu` / 各 `run*` 共用，且都带 `*cobra.Command` + `Options.Interactive`。

### 10.2 目标

```text
shell: character … voice …     → 解析 verb → ui 或 action
ui:    VoiceMenu               → 选 verb → ui 收集 → action
browse: （可选）character 节点显式 voices/ → VoiceMenu
```

- `character/voice/dispatch.go` 只做 argv → verb 路由：

```go
func Dispatch(s streams.Streams, app *api.App, interactive bool, charName string, args []string) error
```

- one-shot 分支保留现有语义：
  - `voice edit <voice> <version> <field> <value>` → `action.ApplyVoiceParam`，输出保持 `v%d.%s updated`；
  - `voice edit <voice> name <value>` 与 `voice duplicate` / `voice add` / `voice delete` 的输出文案不变；
  - `voice list|get` 走 `action` + `render.YAML`。
- 不再有 `voice.Options`；`interactive` 从 `character` shell 传入。
- crumb 归属保持现状：`Dispatch` 无 verb 时先 `prompt.With(charName)`，`ui.VoiceMenu` 内部再 `prompt.With("voices")`；从 `ui.EditCharacter` 进入时外层已有 charName。
- `voice` 包不 import cobra；它 import `streams` / `ui` / `action` / `render` / `parse` / `schema` / `api`。
- `ui` **不 import `voice`**；voice 的字段 schema 在 `schema.VoiceParamFields`，pick/菜单在 ui。
- 可选增强（独立 commit，不阻塞主路径）：把 `character <name> voice` 注册为真正的 cobra 子命令树，让 help / completion / browse 同源。注意 `character voice` 与 `character <name> voice` 的语义冲突，做之前先定 UX；本计划不强制。

### 10.3 `Options.Interactive` 从 voice 包消失

- `voice.Options` 删除。
- interactive 分支只在 `character` shell / `voice.Dispatch` 外壳。
- `ui.VoiceMenu` 恒为 interactive；`action.AddVoice` 等不接收 interactive 标志。

---

<a id="sec-11"></a>
## 11. 编译期 / linter 边界

### 11.1 depguard

在现有 `cli/.golangci.yml` 上扩展。以下为示意；字段以 golangci-lint v2 为准，`files` glob 需按 depguard 文档用 `**/` 前缀并做负向验证。

```yaml
version: "2"

linters:
  enable:
    - depguard
    - staticcheck
  settings:
    depguard:
      rules:
        # 保留现有：CLI 只能走 core/api
        cli:
          files:
            - "**/*.go"
          deny:
            - pkg: "yukumo-script/core/internal"
              desc: "CLI must go through core/api"
            - pkg: "yukumo-script/runtime"
              desc: "CLI must go through core/api"
            - pkg: "github.com/tetratelabs/wazero"

        # 子包不得 import 根包；$ = exact match
        no-root:
          files:
            - "**/action/**/*.go"
            - "**/character/**/*.go"
            - "**/config/**/*.go"
            - "**/menu/**/*.go"
            - "**/parse/**/*.go"
            - "**/play/**/*.go"
            - "**/prompt/**/*.go"
            - "**/render/**/*.go"
            - "**/schema/**/*.go"
            - "**/setting/**/*.go"
            - "**/streams/**/*.go"
            - "**/synth/**/*.go"
            - "**/ui/**/*.go"
          deny:
            - pkg: "yukumo-script/cli$"
              desc: "subpackages must not import root; use leaf packages"

        action:
          files:
            - "**/action/**/*.go"
          deny:
            - pkg: "github.com/spf13/cobra"
            - pkg: "github.com/charmbracelet/huh"
            - pkg: "github.com/charmbracelet/lipgloss"
            - pkg: "yukumo-script/cli/prompt"
            - pkg: "yukumo-script/cli/setting"
            - pkg: "yukumo-script/cli/ui"
            - pkg: "yukumo-script/cli/render"
            - pkg: "yukumo-script/cli/play"
            - pkg: "yukumo-script/cli/streams"
            - pkg: "yukumo-script/cli/menu"
            - pkg: "yukumo-script/cli/internal"
            - pkg: "yukumo-script/cli/config"
            - pkg: "yukumo-script/cli/character"
            - pkg: "yukumo-script/cli/synth"

        ui:
          files:
            - "**/ui/**/*.go"
          deny:
            - pkg: "github.com/spf13/cobra"
            - pkg: "yukumo-script/cli/config"
            - pkg: "yukumo-script/cli/character"
            - pkg: "yukumo-script/cli/synth"
            - pkg: "yukumo-script/cli/internal"
            - pkg: "yukumo-script/cli$"

        schema:
          files:
            - "**/schema/**/*.go"
          deny:
            - pkg: "github.com/spf13/cobra"
            - pkg: "github.com/charmbracelet/huh"
            - pkg: "github.com/charmbracelet/lipgloss"
            - pkg: "yukumo-script/cli/prompt"
            - pkg: "yukumo-script/cli/setting"
            - pkg: "yukumo-script/cli/ui"
            - pkg: "yukumo-script/cli/action"
            - pkg: "yukumo-script/cli/internal"
            - pkg: "yukumo-script/cli$"

        prompt:
          files:
            - "**/prompt/**/*.go"
          deny:
            - pkg: "github.com/spf13/cobra"
            - pkg: "yukumo-script/core/api"
            - pkg: "yukumo-script/cli/schema"
            - pkg: "yukumo-script/cli/setting"
            - pkg: "yukumo-script/cli/action"
            - pkg: "yukumo-script/cli/ui"
            - pkg: "yukumo-script/cli$"

        setting:
          files:
            - "**/setting/**/*.go"
          deny:
            - pkg: "github.com/spf13/cobra"
            - pkg: "yukumo-script/core/api"
            - pkg: "yukumo-script/cli/action"
            - pkg: "yukumo-script/cli/ui"
            - pkg: "yukumo-script/cli$"

        streams-parse:
          files:
            - "**/streams/**/*.go"
            - "**/parse/**/*.go"
          deny:
            - pkg: "yukumo-script"
            - pkg: "github.com/spf13/cobra"
            - pkg: "github.com/charmbracelet/huh"
            - pkg: "github.com/charmbracelet/lipgloss"

        menu-render-play:
          files:
            - "**/menu/**/*.go"
            - "**/render/**/*.go"
            - "**/play/**/*.go"
          deny:
            - pkg: "github.com/spf13/cobra"
            - pkg: "github.com/charmbracelet/huh"
            - pkg: "github.com/charmbracelet/lipgloss"
            - pkg: "yukumo-script/cli/prompt"
            - pkg: "yukumo-script/cli/setting"
            - pkg: "yukumo-script/cli/ui"
            - pkg: "yukumo-script/cli/action"
            - pkg: "yukumo-script/cli/internal"
            - pkg: "yukumo-script/cli$"
```

补充：

- `cmd/**` 规则保持不变。
- `yukumo-script/cli/config` 这类前缀 deny 会同时覆盖其子包；这是期望行为。
- `yukumo-script/cli$` 才是 exact 根；不要写成 `yukumo-script/cli`（会误伤所有子包）。
- 负向验证：在 `action/`、`ui/`、`schema/`、`menu/` 临时加违禁 import，确认 `just lint` 红；删除后恢复。

### 11.2 import-graph 脚本（兜住传递依赖）

depguard 只检查直接 import。新增 `scripts/check-cli-imports.go`：

- 在 `cli/` 下执行 `go list -json ./...`，构建包图；
- 对每个包计算传递闭包，断言禁止边（与 §1 的禁止方向一致）：
  - 任意 `yukumo-script/cli/...` 不得依赖 `yukumo-script/cli`；
  - `ui` 不得传递依赖 `cobra` / `config` / `character` / `synth`；
  - `action` 不得传递依赖 `cobra` / `huh` / `lipgloss` / `prompt` / `setting` / `ui` / `render` / `play` / `streams` / shell；
  - `schema` 只允许 `core/api` + `parse`；`parse` / `streams` 只允许 stdlib；
  - `menu` 只允许 `streams` + stdlib；
  - `render` 只允许 `core/api` + yaml + stdlib；`play` 只允许 `core/api` + stdlib；
  - `setting` 只允许 `schema` + `prompt` + stdlib；`prompt` 只允许 huh + lipgloss + `parse` + stdlib。
- 在 `justfile` 增加：

```make
check-cli-imports:
    go run scripts/check-cli-imports.go

ci: check-replaces check-cli-imports build ...
check: check-replaces check-cli-imports lint ...
```

脚本失败即 CI 失败；这样即便 depguard glob 写错，传递依赖仍有兜底。

---

<a id="sec-12"></a>
## 12. 测试分层

| 层 | 测什么 | 怎么测 |
|---|---|---|
| `prompt` | 既有 path/tree/bool/number/nav | 保持 `cli/test/prompt` |
| `parse` | path / bool / ranged int | 纯表驱动 |
| `schema` | `Lookup` / `Normalize` / clamp warning / keys | 纯表驱动 |
| `action` | use case 返回值和持久化 | `api.Open(api.Options{ConfigDir: t.TempDir()})`；不启 huh；`api.App` 是具体 struct，**不能 fake** |
| `render` | YAML / 结果行 / warning | `bytes.Buffer`；原 `synth/run_test.go` 迁入 |
| `menu` / `cobraui` | 过滤 hidden/help/completion/serve；leaf 用传入 Streams | 小 cobra 树 + `bytes.Buffer` |
| `ui` | 默认不加重测试（huh 交互难单测）；关键纯逻辑抽函数测 | 可选 |
| `shell` | Args 松紧、Hidden、interactive/one-shot 路由 | 轻量；不跑 huh |
| `arch` | import graph | `scripts/check-cli-imports.go`；可选 `cli/test/arch` |
| 回归 | `just ci` + 手动 smoke（interactive / one-shot 各一条路径） | 见 §14 |

迁移期：`synth/run_test.go` 随 `formatResultLine` 移到 `test/render/synth_test.go`，测试数量不减少。

---

<a id="sec-13"></a>
## 13. 迁移步骤

> 与 `1-layer` 不同：本次 **每步必须可构建**。每个 Phase 结束跑 `just ci`（至少 `go -C cli test ./...` + build）。允许内部签名变化，不引入长期兼容 shim。

### Phase 0：基线

1. 确认 `1-layer-refactor` 已在主干 / 当前分支。
2. `just ci` 全绿。
3. 记下当前 interactive / one-shot 的 smoke 命令（§14.4）。
4. 开分支 `cli-layer-refactor`。

### Phase 1：Streams 叶子包 + exact-root 护栏（无行为变化）

1. 新增 `cli/streams/streams.go`（`In` / `Out` / `Err` / `Ctx`，无 cobra）。
2. 新增 `cli/internal/cobraio/fromcmd.go`。
3. 把 shell 内 `*cobra.Command` 仅用于 IO 的函数改成 `streams.Streams`：
   - `config.applyEdit` / `character.applyEdit`
   - `synth.runRequest` / `cmdContext`
   - `character/voice` 的 `editParam*` / menu 错误输出
4. **不要**新建 `action` / `ui`；本阶段只去掉伪 cobra 依赖。
5. `.golangci.yml` 加 `no-root` 规则（`yukumo-script/cli$`）。
6. `just ci`。

### Phase 2：schema / parse 拆分（无行为变化）

1. 新增 `cli/parse`，把 `prompt.ExpandPath` / `ParseBool` / `ParseRangedInt` / `ClampInt` 移到 `parse`；`prompt` 保留薄 wrapper。
2. 新增 `cli/schema`：
   - `field.go`：`Kind` / `Field` / `Entry[T]` / `Keys` / `Lookup` / `Normalize`（返回 `value, warning, err`）。
   - `config.go` / `character.go` / `voice.go`：从 `config/fields.go`、`character/fields.go`、`character/voice/fields.go` 迁移。
3. `setting` 删掉 `entry.go`；只留 `ask.go` / `meta.go`，依赖 `schema` + `prompt`。
4. shell / voice 改用 `schema.*Fields`；删除 shell 包里的 `fields.go`。
5. depguard 加 `schema` / `streams-parse` 规则。
6. `just ci`。

### Phase 3：action / render / play（无行为变化）

1. 新增 `cli/action`：按 §7.1 迁移 use case；不打印、不播放、不接收 Streams。
2. 新增 `cli/render`：YAML / 行 / 合成结果 / 警告。
3. 新增 `cli/play`：从 `synth/preview.go` 迁移播放器逻辑。
4. shell / 现有 interactive 函数改为 `action.*` + `render.*` + `play.*`。
5. 迁移 `formatResultLine` 测试到 `test/render`。
6. depguard 加 `action` / `menu-render-play` 规则。
7. `just ci`。

### Phase 4：抽出 ui 循环 / 表单 / pick（无行为变化）

1. 新增 `cli/ui`；迁 §8.1 的完整清单，包括 add/delete/duplicate/get 的 prompt 分支。
2. 把 `character/pick.go`、`synth/pick.go`、`voice/pick.go`、`voice/menu.go`、`voice/param.go`、`voice/version.go` 的交互部分迁入 `ui`；`voice/duplicate.go` 的 `defaultCopyName` 也迁 `ui`。删除 shell 包里的 pick/menu/param/version 文件。
3. shell `RunE` 变成：interactive 且需要 → `ui.*`，否则 `action.*` + `render.*`。
4. `voice.Dispatch` 改成 `Dispatch(s, app, interactive, charName, args)`；删除 `voice.Options`。
5. depguard 加 `ui` 规则。
6. `just ci`。

### Phase 5：Browse 去 mutation（行为对齐 R1）

1. 新增 `cli/menu` 与 `cli/internal/cobraui`。
2. `ui.Browse` 消费 `menu.Node`；leaf 使用传入的 `Streams`。
3. 删除 `enableInteractiveBrowse` 的 `walk` + `attachBrowse`；只保留 root 的 `enableInteractiveRoot`。
4. `config.NewCmd` 在 interactive 时显式设置 `RunE` → `ui.Browse`。
5. 手动确认：interactive 下 root / config 菜单与改前一致；character / synth 无参行为不变。
6. `just ci`。

### Phase 6：Voice 入口统一

1. `ui.VoiceMenu` 只调 action；`Dispatch` 变薄。
2. `character edit` 的 voices 入口与 `VoiceMenu` 共用同一实现。
3. 确认 `ui` 不 import `voice`。
4. （可选）cobra 子命令化 voice，使 `help` 与 browse 同源；不阻塞。
5. `just ci`。

### Phase 7：收尾与文档

1. 清理死代码：`attachBrowse` / `browseCommands` / `runRequest` / `preview.go` / 双模式 `run*` 残留。
2. 新增 `scripts/check-cli-imports.go`，接入 `justfile`.
3. 更新本文件状态为 Implemented，附 commit 列表。
4. 更新 `README.zh.md` 一句「CLI 内部分为 shell / action / ui / render」。
5. 最终 `just ci`。

---

<a id="sec-14"></a>
## 14. 验收清单

### 14.1 结构

```bash
test -d cli/streams
test -d cli/parse
test -d cli/schema
test -d cli/menu
test -d cli/action
test -d cli/ui
test -d cli/render
test -d cli/play
test -d cli/internal/cobraio
test -d cli/internal/cobraui
```

### 14.2 边界

```bash
# action / ui / schema / parse / streams / menu / render / play / prompt / setting 不得 import cobra
grep -R 'github.com/spf13/cobra' \
  cli/action cli/ui cli/schema cli/parse cli/streams cli/menu cli/render cli/play cli/prompt cli/setting \
  --include='*.go'
# 应为空

# action 不得 import huh / lipgloss / prompt / setting / ui / render / play / streams
grep -R 'charmbracelet/huh\|charmbracelet/lipgloss\|yukumo-script/cli/prompt\|yukumo-script/cli/setting\|yukumo-script/cli/ui\|yukumo-script/cli/render\|yukumo-script/cli/play\|yukumo-script/cli/streams' \
  cli/action --include='*.go'
# 应为空

# 子包不得 import 根包（exact）；根包自己除外
grep -R '"yukumo-script/cli"' \
  cli/action cli/ui cli/schema cli/parse cli/streams cli/menu cli/render cli/play \
  cli/config cli/character cli/synth cli/prompt cli/setting \
  --include='*.go'
# 应为空

# ui 不得 import shell 包
grep -R 'yukumo-script/cli/config\|yukumo-script/cli/character\|yukumo-script/cli/synth' \
  cli/ui --include='*.go'
# 应为空
```

### 14.3 Lint / CI

```bash
just lint                 # 含新 depguard 规则，0 issues
just check-cli-imports    # 传递依赖图断言
just ci
```

负向：在 `action/` 加 cobra import，确认 lint 红；删掉。在 `ui/` 加 `yukumo-script/cli/config` import，确认 lint 红；删掉。

### 14.4 行为回归（输出与流程不变，§15.1 例外）

One-shot（示例）：

```bash
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- config list
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- config get paths.output_path
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- config edit paths.output_path /tmp/out
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- character add remilia "example"
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- character list
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- character get remilia
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- character edit remilia description "updated"
YUKUMO_CONFIG_DIR=$(mktemp -d) just run -- character delete remilia
```

Interactive：

- 启动后进入 browse / 根菜单；
- `config` → edit 字段；
- `character` → get / add / edit / delete；
- `character` → edit → `voices/` 菜单；
- `synth` → `voice` 表单能合成（若环境有 wasm）；
- `synth` 无参在 interactive 下先问 text、再选 voice（与现状一致）。

对比重构前：文案、面包屑层级、Back/Exit 标签一致；§15.1 的「非交互误入 huh」例外改为报错。

### 14.5 Browse 不再 patch

```bash
grep -n 'attachBrowse\|enableInteractiveBrowse\|browseCommands\|menuCommands' cli/*.go
```

理想状态：无定义；`browse.go` 只保留 `enableInteractiveRoot`。

### 14.6 测试数量

`go -C cli test ./...` 通过；既有 prompt / synth 格式化测试不丢（后者迁到 `test/render`）。

### 14.7 Serve 预留

- `action.Synthesize` 不打印、不播放；
- `render` / `play` 不被 serve 调用；
- `Streams` 有 `In`；
- `MenuSpec` 过滤 `serve`（hidden + 显式）。

---

<a id="sec-15"></a>
## 15. 风险与取舍

### 15.1 明确的行为修正（非交互不再误入 huh）

现状部分路径在非 interactive 下仍会启动 huh，这是 bug；本次修正为报错：

- `synth` 无 text 且非 interactive：现状 `prompt.EnsureString` 会弹输入；目标返回 `text is required`。
- `character <name> voice get|edit|delete|duplicate` 缺 voice name 且非 interactive：现状 `ensureVoiceName` 会弹选择；目标返回 `voice name is required`。
- `voice list` 多余 args：现状 `cmd.Usage()`；目标返回带 usage 的 `fmt.Errorf`，文案尽量保持一致。

除此之外，§14.4 的文案、层级、Back/Exit 标签不变。

### 15.2 不新开 module

- 同 module 内分层靠约定 + depguard + import-graph 脚本，没有 `internal` 跨包硬限制。
- 换取：无 replace 矩阵变化、无 go.work 改动、与 `1-layer` 兼容。

### 15.3 Browse UX

- 去掉 `attachBrowse` 后，root / config 显式 `RunE`；character / synth 无参行为与现状一致。
- 风险：漏给某个 group 显式 `RunE` 会变成 Help；Phase 5 按 §14.4 手动点一遍。
- 缓解：`cobraui` 单测断言节点集合与 hidden 过滤；root/config parity 用例。

### 15.4 Voice 双入口

- Phase 6 之前，browse 仍可能看不到 voice；分层仍成立，但导航债留下。
- Phase 6 为必做；cobra 子命令化标为可选，避免 `character voice` 与 `character <name> voice` 语义冲突。

### 15.5 Action 返回结果、render/play 分离

- action 不再直接打印，shell/ui 多一步 `render` / `play`；换来 serve 不需要 quiet 开关，SRP 更清晰。
- 代价：每个命令多 1–3 行编排；可接受。

### 15.6 Interactive 仍是进程级配置

- 不引入 `--interactive` flag；与现状一致。
- 分层后若要加 per-invocation flag，只需改 shell。

### 15.7 大步搬迁 vs 渐进

- 强制每 Phase 可构建；CLI 文件多但依赖面小，适合渐进。
- Phase 1/2 是纯搬运，风险最低；Phase 3/4 是行为回归重点。

### 15.8 行为不变优先于目录完美

- 允许 `config` / `character` 暂不拆 `cmd` 子目录；但 schema、pick、param 必须离开 shell 包，否则循环。
- 目录完美度让位于 depguard 边界与签名清晰。

### 15.9 与 `1-layer` §8 的关系

- `1-layer` §8 要求 CLI 只依赖 api、并保持 browse 不变。
- 本次在 **已满足 §8 的前提下** 重构 CLI 内部；完成后更新本文件状态，不必回改 `1-layer` 正文（可在其 §8 加一行「内部再分层见 2-cli-layer-refactor」）。

### 15.10 未来 Serve 模式

- 评估与落点见 [16. Serve（stdin/stdout）addendum](#sec-16)；不纳入 Phase 0–7。
- 分层本身已允许第三种运输层；实现时勿把二进制帧与文本/TUI 输出混在同一 stdout。

---

<a id="sec-16"></a>
## 16. Serve（stdin/stdout）addendum

> **实现以 [`4-flutter-gui.md`](./4-flutter-gui.md) 为准。** 本节 length-prefixed 帧草图作废；产品协议为 NDJSON JSON-RPC 2（`yukumo-script serve`）。  
> 状态：评估与后续增量设计，不属于当时 Phase 0–7 的实现范围（历史记录保留）。  
> 目的：确认「其他应用通过 stdin/stdout 持续交互」是否被本计划的 shell / action / ui 分层允许，并记下落点与约束，便于以后立项引用。  
> 相关：[`1-layer-refactor.md`](./1-layer-refactor.md) 的 `api.App.Synthesize`（返回 `WAV []byte`）；Web/WASM 是另一条 transport；GUI 桥梁见 [`4-flutter-gui.md`](./4-flutter-gui.md)。

<a id="sec-16-1"></a>
### 16.1 结论

- **架构允许。** Serve 是第三种运输层：不是 huh interactive，也不是 one-shot 文本 CLI，而是长驻协议循环。
- 只依赖已有 `core/api`：`App.Synthesize` → `SynthResult.WAV`；**不必**为 serve 先改 core/runtime。
- 落在 **cobra shell（+ 可选薄协议包）→ action（纯合成）→ api**；**不进** `ui` / `prompt` / `browse` / `render` / `play`。
- 关键风险是 **stdout 污染**（路径文案、usage、huh），不是模块边界。

<a id="sec-16-2"></a>
### 16.2 三种运输层

| 运输层 | 入口 | 依赖 | stdout |
|---|---|---|---|
| One-shot CLI | `synth …` 等 argv | shell → action → render → api | 文本行（路径 / 耗时） |
| Interactive TUI | browse / forms | shell → ui → action | 终端 UI（stderr 提示可并存） |
| Serve（本 addendum） | `serve` 子命令 | shell →（可选 action）→ api | **仅** length-prefixed 帧 |

`cfg.CLI.DefaultInteractive` **不得**改变 `serve` 的 Args / RunE / 输出语义。

<a id="sec-16-3"></a>
### 16.3 落点

```text
cli/cmd/yukumo-script
  ↓
cli root 注册 serve（Hidden from browse / MenuSpec）
  ↓
cli/serve（或 synth 下的 shell）  读帧、写帧、组 SynthRequest
  ↓
api.App.Synthesize（或 action.Synthesize：无打印、无 preview）
```

| 层 | 放什么 | 不放什么 |
|---|---|---|
| shell `serve` | cobra、`SilenceUsage`/`SilenceErrors`、帧编解码循环、ctx | huh、browse、render、play、路径文案 |
| action（可选） | `Synthesize(ctx, app, req)` 纯合成 | `fmt.Fprintln`、`render`、`play` |
| ui / prompt / render / play | — | 全部不参与 |
| core/api | 现有 `Synthesize` 即可 | 除非以后要正式 “bytes-only、禁止落盘” 标志 |

更远：若多宿主共用同一帧格式，可将编解码抽到 `cli/serve/frame` 或独立小包；仍只依赖 `core/api`，与 `1-layer` 一致。

<a id="sec-16-4"></a>
### 16.4 对 action / Streams 的约束

本修订已经把约束提前满足：

1. **action.Synthesize 不打印、不播放**，返回 `(SynthResult, elapsed, error)`；one-shot / interactive 自己 `render` + `play`。
2. Serve 只要 `res.WAV` / `res.Path`；不调用 `render` / `play`。
3. Serve 请求固定：`Preview: false`；`Output` 空或明确「只返回 bytes、不强制写文件」，避免循环里播系统播放器或抢 stdout。
4. `Streams`：`In` 读帧，`Out` 写帧，`Err` 诊断；协议错误与合成错误的呈现策略在立项时定（见 §16.7）。
5. `action` 不接收 `Streams`；serve 也不需要 action 的 IO。

<a id="sec-16-5"></a>
### 16.5 协议草图

示意（非最终规范；立项时再定版本与 schema）：

```text
请求/响应帧：4-byte little-endian length + payload
```

```go
// 示意：挂在 cobra 上的 serve 循环骨架
serve := &cobra.Command{
	Use:          "serve",
	Hidden:       true,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := injectedApp // 从 Execute 注入的 *api.App
		s := cobraio.FromCmd(cmd)
		in := bufio.NewReader(s.In)
		w := bufio.NewWriter(s.Out)
		for {
			var n uint32
			if err := binary.Read(in, binary.LittleEndian, &n); err != nil {
				if err == io.EOF {
					return nil // parent 关闭 stdin：干净退出
				}
				return err
			}
			payload := make([]byte, n)
			if _, err := io.ReadFull(in, payload); err != nil {
				return err
			}
			req, err := decodeSynthRequest(payload) // 纯 text 或 JSON 等，立项时定
			if err != nil {
				return err
			}
			req.Preview = false
			res, _, err := action.Synthesize(s.Ctx, app, req)
			if err != nil {
				return err // 或写错误帧到协议里，见 §16.7
			}
			if err := binary.Write(w, binary.LittleEndian, uint32(len(res.WAV))); err != nil {
				return err
			}
			if _, err := w.Write(res.WAV); err != nil {
				return err
			}
			if err := w.Flush(); err != nil {
				return err
			}
		}
	},
}
```

最小 payload 可以是 UTF-8 文本；产品形态更可能是 JSON/msgpack（`voice` / `version` / `speed` / `text`）再映射到 `api.SynthRequest`。

生命周期：沿用 `Execute` 的 `api.Open` + defer `Close`；用 `s.Ctx` 支持取消。Parent 关 stdin → `EOF` → 退出。

<a id="sec-16-6"></a>
### 16.6 与 browse / interactive 的隔离

- `serve` **排除**出 MenuSpec / browse（`Hidden: true`，`cobraui` 显式过滤）。
- 不要对 `serve` 调用 `attachBrowse` 或任何 huh。
- Root 在 interactive 模式下仍可注册 `serve`；用户显式 `… serve` 进入协议模式，与默认 browse 无关。

<a id="sec-16-7"></a>
### 16.7 待定问题

1. 帧 payload：纯 text vs JSON/msgpack？是否需要协议 version 字节？
2. 错误：合成失败是退出进程、还是写错误帧（仍 length-prefix，payload 为错误码/消息）？
3. 是否允许多种 opcode（synth / list voices / health），还是只做 synth？
4. 缓存 / 落盘：serve 默认 `ReuseCache` 与是否写 WAV 文件，是否与 CLI one-shot 共用配置？
5. 并发：单连接串行是否足够，是否要在帧层支持 request id？
6. 与 Web/WASM：serve 是本机子进程 IPC；web 是另一套 host。二者都只应依赖 `core/api`，不共享帧格式除非有意为之。

<a id="sec-16-8"></a>
### 16.8 对本文 Phase 的影响

| 章节 / Phase | Serve addendum 要求 |
|---|---|
| §4 层职责 | 运输层枚举含 serve；ui/render/play 不参与 |
| §5 Streams | `In`/`Out`/`Err` 分离；serve 占满 `Out` 做二进制 |
| §7 Action | `Synthesize` 返回结果；无打印、无 play |
| §9 Browse | MenuSpec 排除 `serve` |
| Phase 0–7 | **不实现** serve；只避免封死上述约束 |
| 立项后 | 新 Phase：协议定稿 + `cli/serve` + 测试（temp config / fake host） |

不要求为了 serve 回头改 `1-layer` 的 module 图；若以后要 GUI/宿主复用同一协议，再考虑把帧编解码移出 `cli`（仍依赖 `core/api`）。

---

<a id="sec-a"></a>
## 附录 A：文件映射

| 现在 | 之后 |
|---|---|
| （无） | `cli/streams/streams.go` |
| （无） | `cli/parse/parse.go` |
| （无） | `cli/schema/{field,config,character,voice}.go` |
| `config/fields.go` | `cli/schema/config.go` |
| `character/fields.go` | `cli/schema/character.go` |
| `character/voice/fields.go` | `cli/schema/voice.go`（含 `paramEntries`） |
| `prompt` 的纯解析 | `cli/parse`；`prompt` 保留 wrapper |
| `setting/entry.go` | `cli/schema/field.go` |
| `setting/ask.go` / `setting` 的 Metas | `setting/ask.go` / `setting/meta.go` |
| `config/edit.go` 的 `runInteractiveEdit` | `ui/config_edit.go` |
| `config/edit.go` 的 `applyEdit` | `action.ApplyConfigEdit` + `render` |
| `character/edit.go` 的 `runInteractiveEdit` | `ui/character_edit.go` |
| `character/edit.go` 的 `applyEdit` | `action.ApplyCharacterEdit` + `render` |
| `character/add.go` interactive | `ui/character_add.go` |
| `character/delete.go` interactive | `ui/character_delete.go` |
| `character/pick.go` | `ui/character_pick.go` |
| `character/get.go` pick | `ui/character_pick.go` |
| `character/voice/menu.go` | `ui/voice_menu.go` |
| `character/voice/param.go` 交互部分 | `ui/voice_param.go` |
| `character/voice/pick.go` | `ui/voice_pick.go` |
| `character/voice/add.go` interactive | `ui/voice_add.go` |
| `character/voice/edit.go` interactive | `ui/voice_edit.go` |
| `character/voice/delete.go` interactive | `ui/voice_delete.go` |
| `character/voice/duplicate.go` interactive | `ui/voice_duplicate.go` |
| `character/voice/{add,edit,delete,duplicate,get,list}.go` 持久化 | `action/voice.go` + `render` |
| `character/voice/dispatch.go` | 薄路由，留在 voice；无 cobra / 无 Options |
| `synth/run.go` | `action/synth.go` + `render/synth.go` |
| `synth/preview.go` | `play/play.go` |
| `synth/pick.go` | `ui/synth_pick.go` |
| `synth/voice_form.go` | `ui/synth_form.go` + `synth/voice.go`（shell 入口） |
| `synth/run_test.go` | `test/render/synth_test.go` |
| `browse.go` `attachBrowse` / `enableInteractiveBrowse` | 删除；`browse.go` 只留 `enableInteractiveRoot` |
| `browse.go` 循环 | `ui/browse.go` |
| `browse.go` cobra → 菜单 | `internal/cobraui/menu.go` + `menu/menu.go` |
| （无） | `scripts/check-cli-imports.go` |
| `cli/.golangci.yml` | 增加 action/ui/schema/no-root 等规则 |

---

<a id="sec-b"></a>
## 附录 B：典型 RunE 前后对比

### 现在（耦合）

```go
func newEditCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use: "edit [name] [field] [value]",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(3)(cmd, args)
			}
			return cobra.ExactArgs(3)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive && len(args) < 3 {
				return runInteractiveEdit(cmd, app, args) // 循环 + api + 打印
			}
			return runOneShotEdit(cmd, app, args) // api + 打印，仍持有 *cobra.Command
		},
	}
}
```

### 之后（分层）

```go
func newEditCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use: "edit [name] [field] [value]",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(3)(cmd, args)
			}
			return cobra.ExactArgs(3)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if opts.Interactive && len(args) < 3 {
				return ui.EditCharacter(s, app, argAt(args, 0), argAt(args, 1), argAt(args, 2))
			}
			if len(args) < 3 {
				return fmt.Errorf("name, field, and value are required")
			}
			c, result, err := action.ApplyCharacterEdit(s.Ctx, app, args[0], args[1], args[2], true)
			if err != nil {
				return err
			}
			if result.Warning != "" {
				_ = render.Warning(s.Err, result.Warning)
			}
			return render.Line(s.Out, "%s: %s\n", args[1], result.Value)
		},
	}
}
```

**Bottom line：** 在 `1-layer` 已把 CLI 隔在 `core/api` 之外的前提下，本次在 `cli` module 内复制 `prompt`/`setting` 的成功经验——**cobra 做壳，action 做 use case，render 做输出，play 做 OS 副作用，ui 做 huh**——并用 `streams` / `schema` / `menu` / `internal` 适配器消除反向依赖；browse 从「改写 RunE」改为「显式菜单入口」，voice 从双轨导航收敛到同一 `ui.VoiceMenu`。未来 stdin/stdout **serve** 作为第三种运输层已评估可接入，见 [§16](#sec-16)；不纳入本次 Phase。
