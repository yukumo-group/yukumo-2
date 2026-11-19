import 'dart:async';
import 'dart:convert';

import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../bridge/client.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import '../synth/wav_preview.dart';
import 'voice_param_editor.dart';

const _kPreviewText = 'ゆっくりしていってね';

class VoiceEditPage extends StatefulWidget {
  const VoiceEditPage({
    super.key,
    required this.characterName,
    required this.voiceName,
    this.embedded = false,
    this.onBack,
    this.onUpdated,
    this.onVoiceRenamed,
  });

  final String characterName;
  final String voiceName;

  /// When true, render only the editor body (no Scaffold / AppBar).
  final bool embedded;

  /// Shown as a back control in the embedded title (e.g. dual-column overlay).
  final VoidCallback? onBack;

  final VoidCallback? onUpdated;
  final ValueChanged<String>? onVoiceRenamed;

  @override
  State<VoiceEditPage> createState() => _VoiceEditPageState();
}

class _VoiceEditPageState extends State<VoiceEditPage> {
  Timer? _saveTimer;
  YukumoClient? _client;
  String _characterName = '';
  String _voiceName = '';
  String _description = '';
  List<Voice> _voices = const [];
  int _selectedVersion = 1;
  bool _loading = true;
  bool _ready = false;
  bool _saving = false;
  bool _saveQueued = false;
  bool _dirty = false;
  bool _previewing = false;
  String? _error;

  Voice? get _voice {
    for (final v in _voices) {
      if (v.name == _voiceName) return v;
    }
    return null;
  }

  @override
  void initState() {
    super.initState();
    _characterName = widget.characterName;
    _voiceName = widget.voiceName;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_load());
    });
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _client = AppScope.of(context);
  }

  @override
  void didUpdateWidget(covariant VoiceEditPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.characterName != widget.characterName ||
        oldWidget.voiceName != widget.voiceName) {
      unawaited(_switchTarget(widget.characterName, widget.voiceName));
    }
  }

  @override
  void dispose() {
    _saveTimer?.cancel();
    if (_dirty && _ready) {
      unawaited(_persist(notify: false));
    }
    super.dispose();
  }

  Future<void> _switchTarget(String characterName, String voiceName) async {
    await _flushSave();
    if (!mounted) return;
    _characterName = characterName;
    _voiceName = voiceName;
    await _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _ready = false;
      _dirty = false;
      _error = null;
    });
    try {
      final c = await (_client ?? AppScope.of(context)).characters.get(
        CharacterGetReq(name: widget.characterName),
      );
      if (!mounted) return;
      final voices = List<Voice>.of(c.voices ?? const <Voice>[]);
      final exists = voices.any((v) => v.name == widget.voiceName);
      setState(() {
        _characterName = c.name;
        _description = c.description;
        _voices = voices;
        _voiceName = widget.voiceName;
        _loading = false;
        _ready = exists;
        _error = exists ? null : '找不到音色「${widget.voiceName}」';
      });
      if (!exists) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (!mounted) return;
          if (widget.embedded) {
            widget.onBack?.call();
          } else {
            Navigator.of(context).maybePop();
          }
        });
      }
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = '$e';
        _loading = false;
      });
      if (!widget.embedded) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (mounted) Navigator.of(context).maybePop();
        });
      }
    }
  }

  void _scheduleSave({Duration delay = const Duration(milliseconds: 450)}) {
    if (!_ready || _loading) return;
    _dirty = true;
    _saveTimer?.cancel();
    _saveTimer = Timer(delay, () {
      unawaited(_persist());
    });
  }

  Future<void> _flushSave() async {
    _saveTimer?.cancel();
    if (!_dirty || !_ready) return;
    await _persist();
  }

  Future<void> _persist({bool notify = true}) async {
    if (!_ready) return;
    if (_saving) {
      _saveQueued = true;
      return;
    }
    final client = _client;
    if (client == null) return;

    _saving = true;
    final voices = List<Voice>.of(_voices);
    final characterName = _characterName;
    final description = _description;
    final snapshotDirty = _dirty;
    try {
      await client.characters.update(
        CharacterUpdateReq(
          oldName: characterName,
          character: Character(
            name: characterName,
            description: description,
            voices: voices,
          ),
        ),
      );
      _dirty = false;
      if (notify && mounted) widget.onUpdated?.call();
    } catch (e) {
      _dirty = snapshotDirty || _dirty;
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('保存失败: $e')),
        );
      }
    } finally {
      _saving = false;
      if (_saveQueued) {
        _saveQueued = false;
        await _persist(notify: notify);
      }
    }
  }

  Future<void> _handleBack() async {
    await _flushSave();
    widget.onBack?.call();
  }

  void _onVoiceChanged(Voice next) {
    if (next.name != _voiceName && _voices.any((v) => v.name == next.name)) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('音色名已存在')),
      );
      return;
    }
    final renamed = next.name != _voiceName;
    setState(() {
      _voices = [
        for (final v in _voices)
          if (v.name == _voiceName) next else v,
      ];
      _voiceName = next.name;
    });
    if (renamed) widget.onVoiceRenamed?.call(next.name);
    _scheduleSave();
  }

  Future<void> _preview() async {
    final voice = _voice;
    if (voice == null || _previewing) return;
    final params = voice.params ?? const <SynthParam>[];
    final hasVersion = params.any((p) => p.version == _selectedVersion);
    if (!hasVersion) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('请先启用 AQ$_selectedVersion')),
      );
      return;
    }

    setState(() => _previewing = true);
    try {
      final client = _client ?? AppScope.of(context);
      final res = await client.synth.render(
        SynthRequest(
          voice: voice,
          voiceBase: _characterName,
          text: _kPreviewText,
          version: _selectedVersion,
          speedFactor: 1,
          preview: false,
          noWrite: true,
        ),
      );
      if (!mounted) return;
      final b64 = res.wAV;
      if (b64 == null || b64.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('预览无音频数据')),
        );
        return;
      }
      unawaited(() async {
        try {
          await playWavBytes(base64Decode(b64));
        } catch (e) {
          if (!mounted) return;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('播放失败: $e')),
          );
        }
      }());
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('预览失败: $e')),
      );
    } finally {
      if (mounted) setState(() => _previewing = false);
    }
  }

  Widget _previewFab() {
    return M3EFab(
      tooltip: '预览「$_kPreviewText」',
      color: M3EFabColor.primary,
      size: M3EFabSize.medium,
      onPressed: (_loading || _error != null || _voice == null || _previewing)
          ? null
          : () => unawaited(_preview()),
      icon: _previewing
          ? const SizedBox(
              width: 22,
              height: 22,
              child: CircularProgressIndicator(strokeWidth: 2.5),
            )
          : const Icon(Icons.play_arrow_rounded),
    );
  }

  Widget _editorBody(ColorScheme cs) {
    if (_loading) {
      return const Padding(
        padding: EdgeInsets.symmetric(vertical: 48),
        child: Center(child: M3ELoadingIndicator()),
      );
    }
    if (_error != null) {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 48),
        child: Center(child: Text(_error!)),
      );
    }
    final voice = _voice;
    if (voice == null) {
      return const Padding(
        padding: EdgeInsets.symmetric(vertical: 48),
        child: Center(child: Text('音色不存在')),
      );
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (widget.embedded) ...[
          PageCardTitle(
            _voiceName,
            leading: widget.onBack == null
                ? null
                : IconButton(
                    tooltip: '返回角色',
                    onPressed: () => unawaited(_handleBack()),
                    visualDensity: VisualDensity.compact,
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(
                      minWidth: 40,
                      minHeight: 40,
                    ),
                    icon: const Icon(Icons.arrow_back_rounded),
                  ),
          ),
        ],
        VoiceParamEditor(
          voice: voice,
          onChanged: _onVoiceChanged,
          onSelectedVersionChanged: (version) {
            if (_selectedVersion == version) return;
            setState(() => _selectedVersion = version);
          },
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final body = _editorBody(cs);

    if (widget.embedded) {
      return SizedBox.expand(
        child: Stack(
          children: [
            Positioned.fill(
              child: SingleChildScrollView(
                padding: const EdgeInsets.only(bottom: 80),
                child: body,
              ),
            ),
            Positioned(
              right: 4,
              bottom: 4,
              child: _previewFab(),
            ),
          ],
        ),
      );
    }

    return PopScope(
      onPopInvokedWithResult: (didPop, _) {
        if (didPop) unawaited(_flushSave());
      },
      child: Scaffold(
        backgroundColor: cs.surface,
        appBar: AppBar(
          title: Text(_loading ? '编辑音色' : _voiceName),
        ),
        floatingActionButton: _previewFab(),
        body: _loading || _error != null || _voice == null
            ? body
            : SingleChildScrollView(
                padding: const EdgeInsets.only(top: 8, bottom: 96),
                child: ContentFrame(child: body),
              ),
      ),
    );
  }
}
