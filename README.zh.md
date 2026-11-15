# Yukumo Script

AquesTalk WASM（`third-party/aq_wasm`）语音合成与角色/音色管理 CLI。

[English](README.md)

## 架构

分层（依赖方向只向下）：

```
cli → service → domain
              ↘ store
              ↘ runtime(aq) → aq_wasm.wasm（Yukw）→ AquesTalk 引擎
```

| 层 | 职责 |
|---|---|
| `domain` | 纯数据模型，数据默认值和规范化 |
| `store` | 通用数据读写，YAML 数据编解码 |
| `runtime/aq` | 底层引擎，wazero 实例、`aq_*` 导出、多线程实例池 |
| `service/*` | 各种业务逻辑实现，通过接口依赖 store/runtime/其它 service |
| `cli/*` | 交互层，cobra 命令、huh TUI、音频预览 |

**Yukw**：用 Rust 编写的兼容层，把多个 AquesTalk 产品的 WASM 串成单一模块 `aq_wasm.wasm`，并导出统一的 `aq_*` 合成/查询 API，以及 `malloc` / `free` 等宿主侧内存分配接口。本仓库的 `runtime/aq` 经 wazero 加载该模块。

### 多线程

引擎侧通过 `aq.Pool` 支持多 goroutine 并行使用；播放预览也可在后台 goroutine 进行。

- 由于 WASM 模块本身不线程安全，一个 `Runtime`（一个实例）同一时刻只能被一个 goroutine 使用，禁止跨 goroutine 共用同一实例。
- `aq.Pool` 按需创建实例：`Acquire` 借出独占 `Runtime`，用完 `Release` 归还；`Do` / `Synthesize` / `HasPreset` / `PresetNames` 使用同一套租借，因此多个合成/查询可以并行，彼此各用各的实例。
- 交互模式 `-p` / preview：播放在后台 goroutine，不阻塞 huh 界面；再次合成会先停掉上一次播放进程。

### 领域模型

- `Character`：`name`（id）, `description`, `voices[]`；文件 `{Sanitize(name)}.yaml`。
- `Voice`：`name`（id） + `params[]`；同一 voice 最多 3 条 param，版本唯一 ∈ {1,2,10}。
- `SynthParam`：对齐 `AqSynthParam`。公共：`preset`, `speed`(50–300)。v2：`phont_data`。v10：`base`(F1E/F2E/M1E), `volume`, `pitch`, `accent`, `lmd`, `fsc`。`Normalize` 按版本清无关字段并 clamp。`speed` 仍挂在音色上，表示该预设的急促程度；合成时再乘 CLI `-s` 因子。
- 合成有效语速：`clamp(round(voice.speed * factor), 50, 300)`；`factor≤0` 视为 `1`。

### 解耦

- 实现了 `store.Store` 接口及 `Reader`、`Writer`、`ExistenceChecker`、`Lister`、`Deleter` 子接口，使得增删查改可以换不同存储实现。单元测试用 `Memory`，不必碰真实文件系统；生产用 `FS` 输出。
- `service/character/voice` 只依赖 `Characters`（读改角色）和 `PresetCatalog`（问引擎某 preset 是否合法），不直接依赖 character 的磁盘路径，也不绑定某一个全局 WASM 实例。因此同一套 voice 逻辑既能挂真实角色 service + `aq.Pool`，也能在测试里换成内存角色与假 catalog。
- `service/synth` 只依赖 `Config`、`Voices`、`Engine` 三个接口：向谁要配置、向谁取音色、用谁合成，都可以替换。输出策略（写不写、写哪）留在 service；终端打印路径与播放预览留在 CLI（`runRequest`），互不拖拽。
- `cli/setting` 描述「有哪些字段、如何读写模型」；`cli/prompt` 负责「怎么问用户」。config / character / voice param 共用同一套 Ask、Browse、Normalize，改交互控件不必改业务字段表，反之亦然。
- `cli/browse` 在 `cli.default_interactive=true` 时给「有子命令」的节点挂菜单导航；真正干活的仍是各叶子命令自己的 `RunE`，命令树与交互壳分离。

## 构建与运行

需可用的 `third-party/aq_wasm.wasm`（或设 `YUKUMO_AQ_WASM`）。

```bash
go build -o yukumo-script .
go test ./test/...

# 交互根菜单（cli.default_interactive 默认为 true）
go run .

# 创建角色 + 音色（二选一，勿对同一角色重复加同名 default）
# A：character add 时自动带 default（v1）
go run . config edit characters.init_default_voice true
go run . character add remilia "example"
# B：手动添加音色
# go run . character add remilia "example"
# go run . character remilia voice add default 1
# 注意：默认 interactive 下，voice add 会进入 param 编辑；确认默认值退出即可。
# 纯脚本可先：go run . config edit cli.default_interactive false

go run . synth -v remilia:default "こんにちは"
go run . synth -v remilia:default -p "こんにちは"          # 预览播放
go run . synth -v remilia:default -o /tmp/out.wav "テスト" # 指定输出
go run . synth -P ./remilia.yaml:default "こんにちは"      # 覆盖角色文件路径
```

## CLI 命令

全局：无显式子命令且 interactive 开启时进入根菜单。`get` 在 interactive 下对 config/character **隐藏**（改走 edit/browse）。

### `config`

| 命令 | 行为 |
|---|---|
| `config list` | 打印全部设置键与当前值 |
| `config get <key>` | 打印单键（non-interactive） |
| `config edit <key> <value>` | 一次性改键 |
| `config edit` / `config edit <key>` | interactive：树形浏览 `paths.*` / `characters.*` / `synth.*` / `cli.*`；bool 选中即翻转 |

### `character`

| 命令 | 行为 |
|---|---|
| `character list` | JSON 列表 |
| `character get [name]` | 打印角色；缺名 interactive 可选 |
| `character add [name] [description]` | 创建；可受 `init_default_voice` 影响 |
| `character edit [name] [field] [value]` | 改 `name`/`description`；interactive 为字段菜单，含进入 `voices` |
| `character delete [name]` | 删除 |
| `character <name> voice <verb> …` | 音色子命令（见下） |

### `character <name> voice …`

| 动词 | 行为 |
|---|---|
| `list` | 列出音色 |
| `get [voice]` | 打印音色 |
| `add [voice]` | 创建；interactive 选版本并编辑 param |
| `edit [voice] …` | 改名或 `edit <voice> <1\|2\|10> <field> <value>`；interactive 编辑各版本 param |
| `delete [voice]` | 删除 |
| `duplicate [voice] [newName]` | 深拷贝；默认名 `nameutil` 生成 |
| （仅 verb 省略且 interactive） | 打开 voice 菜单 |

Param 字段：`preset`, `speed`；v2 +`phont_data`（选文件嵌入）；v10 +`base`,`volume`,`pitch`,`accent`,`lmd`,`fsc`。preset 经 `PresetCatalog`（WASM）校验。

### `synth`

| 命令 / 标志 | 行为 |
|---|---|
| `synth -v character:voice … TEXT` | 从配置的 `character_path` 取角色；`-s` 默认 `1` |
| `synth -P file.yaml:voice … TEXT` | 用指定 Character YAML（**最后一个** `:` 分割路径与音色）；与 `-v` **互斥** |
| 共用：`[-V 1\|2\|10] [-s factor] [-o path] [-p] [--reuse-cache[=false]]` | 版本 / 语速因子 / 输出 / 预览 / 是否读缓存（flag 覆盖配置，省略则用 `synth.reuse_cache`） |
| `synth voice` | **仅 interactive**：表单循环；voice 树键为 `character:voice`；preview **异步** |
| `synth custom` | stub，未实现 |

**版本**：`-V` 省略 → voice 内 **1→2→10** 首个 param；皆无则 `params empty` / `definition invalid`。  
**引用**：`角色:音色`（第一个 `:`）。`-P` 为 `--voice-path` 短选项（cobra 仅支持单字符短旗标；语义上即「voice-path」）。

**缓存（`output_path/index.yaml`）**

- 任务哈希（SHA-256 hex）：文本、有效 speed、完整 `SynthParam`（含 phont）、以及 `-P` 时整份 YAML 字节。
- index 为简单 map：`hash: basename.wav`。启动时清理指向缺失文件的条目。
- `reuse_cache` / `--reuse-cache` 启用且命中：不重新合成；**复制**到新路径（`-o` 或模板生成的新唯一名），避免多线程挪动同一文件。源文件保留。
- 命中但文件已丢：从 index 删掉该条目并重合成，不报错。
- 不落盘（如仅 preview 且 `always_write=false`）时不写 index；若有缓存可只读入 WAV 供预览。

## 配置（`config.yaml`）

配置**目录**为 `$YUKUMO_CONFIG_DIR`（若设置）或 `{UserConfigDir}/yukumo-script`；文件为该目录下的 `config.yaml`。  
启动时对 `character_path` / `output_path` 做 `MkdirAll`，失败（含无权限）**panic**；空路径字段回填默认值。

| Key | 默认 | 说明 |
|---|---|---|
| `paths.character_path` | `{UserConfigDir}/yukumo-script/characters` | 角色 YAML 目录 |
| `paths.output_path` | `{TempDir}/yukumo-script` | 合成输出目录 |
| `characters.init_default_voice` | `false` | `character add` 时是否自动加 `default` 音色（v1）；与手动 `voice add` **二选一**，勿重复同名 |
| `synth.output_format` | `{character}:{voice}_{yyyy}-{MM}-{dd}:{hh}-{mm}-{ss}_{ms}` | 输出文件名模板（仅花括号） |
| `synth.always_write` | `true` | 是否总是写文件（见合成输出） |
| `synth.reuse_cache` | `true` | 按任务哈希复用 `index.yaml` 中的已有 WAV |
| `cli.default_interactive` | `true` | 缺参时进入 huh；根命令挂载 browse |
| `cli.with_accessibility` | `false` | huh 无障碍模式 |

布尔缺省用指针区分「未设置」与显式 `false`。

**输出文件名模板**（仅 `{…}`，无紧凑串）

- 命名：`{character}`/`{c}`，`{voice}`/`{v}`
- 时间：同字母重复长度 = 零填充宽度。`y` 年、`M` 月、`d` 日、`h`/`H` 时（24h）、`m` 分、`s` 秒；`{ms}` 毫秒（三位）
- 默认 → 如 `remilia:default_2026-11-15:13-04-05_123.wav`（`{MM}`=月，`{mm}`=分）
- 冲突时自动加序号：`name.wav` → `name-1.wav` …

**合成输出策略**

| `always_write` | `-o` | `-p` preview | 行为 |
|---|---|---|---|
| true | 无 | 任意 | 按模板写入 `output_path` |
| true/false | 有 | 任意 | 写入 `-o` |
| false | 无 | 有 | 不输出；预览用临时文件（播完删） |
| false | 无 | 无 | 仍按模板写入 |
