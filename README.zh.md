# Yukumo Script

AquesTalk WASM（`third-party/aq_wasm`）语音合成与角色/音色管理 CLI。

[English](README.md)

## 构建与运行

需可用的 `third-party/aq_wasm.wasm`（或设 `YUKUMO_AQ_WASM`）。开发命令使用 `just`；`just lint` / `just build` 需要 golangci-lint v2（可用 `just install-tools`）。

```bash
just build          # depguard + staticcheck 通过后再 build
just test           # 全部 module
just ci             # replace 检查 + lint + build + test + vet

# 运行 CLI（无参数时打印帮助）
just run
# 或：go run ./cli/cmd/yukumo-script

# 创建角色 + 音色（二选一，勿对同一角色重复加同名 default）
# A：character add 时自动带 default（v1）
go run ./cli/cmd/yukumo-script config edit characters.init_default_voice true
go run ./cli/cmd/yukumo-script character add remilia "example"
# B：手动添加音色
# go run ./cli/cmd/yukumo-script character add remilia "example"
# go run ./cli/cmd/yukumo-script character remilia voice add default 1
# voice add 必须给出音色名；版本 1|2|10 可选（默认 1）

go run ./cli/cmd/yukumo-script synth -v remilia:default "こんにちは"
go run ./cli/cmd/yukumo-script synth -v remilia:default -p "こんにちは"          # 预览播放（阻塞至播完）
go run ./cli/cmd/yukumo-script synth -v remilia:default -o /tmp/out.wav "テスト" # 指定输出
go run ./cli/cmd/yukumo-script synth -P ./remilia.yaml:default "こんにちは"      # 覆盖角色文件路径
```

## CLI 命令

CLI 为脚本化（仅 cobra）：无交互提示或菜单；缺少必需参数会报错退出。无参数的 `yukumo-script`、`config`、`character` 会打印帮助。

### `config`

| 命令 | 行为 |
|---|---|
| `config list` | 打印全部设置键与当前值 |
| `config get <key>` | 打印单键 |
| `config edit <key> <value>` | 一次性改键（两个参数均必填） |

### `character`

| 命令 | 行为 |
|---|---|
| `character list` | JSON 列表 |
| `character get <name>` | 打印角色 |
| `character add <name> [description]` | 创建；可受 `init_default_voice` 影响 |
| `character edit <name> <field> <value>` | 改 `name` / `description` |
| `character delete <name>` | 删除 |
| `character <name> voice <verb> …` | 音色子命令（见下） |

### `character <name> voice …`

| 动词 | 行为 |
|---|---|
| `list` | 列出音色 |
| `get <voice>` | 打印音色 |
| `add <voice> [1\|2\|10]` | 按版本默认 param 创建（版本默认 `1`） |
| `edit <voice> name <newName>` 或 `edit <voice> <1\|2\|10> <field> <value>` | 改名或设置某版本的 param 字段 |
| `delete <voice>` | 删除 |
| `duplicate <voice> [newName]` | 深拷贝；默认名 `nameutil` 生成 |

动词必填；仅写 `character <name> voice` 会报错。

Param 字段：`preset`, `speed`；v2 +`phont_data`（选文件嵌入）；v10 +`base`,`volume`,`pitch`,`accent`,`lmd`,`fsc`。preset 经 `PresetCatalog`（WASM）校验。

### `synth`

| 命令 / 标志 | 行为 |
|---|---|
| `synth -v character:voice … TEXT` | 从配置的 `character_path` 取角色；`-s` 默认 `1` |
| `synth -P file.yaml:voice … TEXT` | 用指定 Character YAML（**最后一个** `:` 分割路径与音色）；与 `-v` **互斥** |
| 共用：`[-V 1\|2\|10] [-s factor] [-o path] [-p] [--reuse-cache[=false]]` | 版本 / 语速因子 / 输出 / 阻塞式预览播放 / 是否读缓存（flag 覆盖配置，省略则用 `synth.reuse_cache`） |
| `synth custom` | stub，未实现 |

**版本**：`-V` 省略 → voice 内 **1→2→10** 首个 param；皆无则 `params empty` / `definition invalid`。  
**引用**：`角色:音色`（第一个 `:`）。`-P` 为 `--voice-path` 短选项（cobra 仅支持单字符短旗标；语义上即「voice-path」）。

**缓存（`output_path/index.yaml`）**

- 任务哈希（SHA-256 hex）：文本、有效 speed、完整 `SynthParam`（含 phont）、以及 `-P` 时整份 YAML 字节。
- index 为简单 map：`hash: basename.wav`。启动时清理指向缺失文件的条目。
- `reuse_cache` / `--reuse-cache` 启用且命中：不重新合成；**复制**到新路径（`-o` 或模板生成的新唯一名），避免多线程挪动同一文件。源文件保留。
- 命中但文件已丢：从 index 删掉该条目并重合成，不报错。
- 不落盘（如仅 preview 且 `always_write=false`）时不写 index；若有缓存可只读入 WAV 供预览。
- 每次合成结束后打印耗时，写在输出路径（及 `(cached)`）后面，例如 `/tmp/a.wav (cached) 12.34ms`。

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
