import 'dart:async';

import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import 'voice_param_editor.dart';

class VoiceEditPage extends StatefulWidget {
  const VoiceEditPage({
    super.key,
    required this.characterName,
    required this.voiceName,
    this.embedded = false,
    this.onUpdated,
    this.onVoiceRenamed,
  });

  final String characterName;
  final String voiceName;

  /// When true, render only the editor body (no Scaffold / AppBar).
  final bool embedded;

  final VoidCallback? onUpdated;
  final ValueChanged<String>? onVoiceRenamed;

  @override
  State<VoiceEditPage> createState() => _VoiceEditPageState();
}

class _VoiceEditPageState extends State<VoiceEditPage> {
  Timer? _saveTimer;
  String _characterName = '';
  String _voiceName = '';
  String _description = '';
  List<Voice> _voices = const [];
  bool _loading = true;
  bool _ready = false;
  bool _saving = false;
  bool _saveQueued = false;
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
  void didUpdateWidget(covariant VoiceEditPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.characterName != widget.characterName ||
        oldWidget.voiceName != widget.voiceName) {
      _flushSave();
      _characterName = widget.characterName;
      _voiceName = widget.voiceName;
      unawaited(_load());
    }
  }

  @override
  void dispose() {
    _saveTimer?.cancel();
    super.dispose();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _ready = false;
      _error = null;
    });
    try {
      final c = await AppScope.of(context).characters.get(
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
      if (!exists && !widget.embedded) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (mounted) Navigator.of(context).maybePop();
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
    _saveTimer?.cancel();
    _saveTimer = Timer(delay, () {
      unawaited(_save());
    });
  }

  void _flushSave() {
    _saveTimer?.cancel();
    if (!_ready || _loading) return;
    unawaited(_save());
  }

  Future<void> _save() async {
    if (!mounted || !_ready || _loading) return;
    if (_saving) {
      _saveQueued = true;
      return;
    }
    _saving = true;
    final client = AppScope.of(context);
    final voices = List<Voice>.of(_voices);
    final characterName = _characterName;
    final description = _description;
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
      if (!mounted) return;
      widget.onUpdated?.call();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('保存失败: $e')),
      );
    } finally {
      _saving = false;
      if (_saveQueued) {
        _saveQueued = false;
        unawaited(_save());
      }
    }
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
          PageCardTitle(_voiceName),
          const SizedBox(height: 16),
        ],
        VoiceParamEditor(
          voice: voice,
          onChanged: _onVoiceChanged,
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final body = _editorBody(cs);

    if (widget.embedded) {
      return body;
    }

    return PopScope(
      onPopInvokedWithResult: (didPop, _) {
        if (didPop) _flushSave();
      },
      child: Scaffold(
        backgroundColor: cs.surface,
        appBar: AppBar(
          title: Text(_loading ? '编辑音色' : _voiceName),
        ),
        body: _loading || _error != null || _voice == null
            ? body
            : SingleChildScrollView(
                padding: const EdgeInsets.only(top: 8, bottom: 32),
                child: ContentFrame(child: body),
              ),
      ),
    );
  }
}
