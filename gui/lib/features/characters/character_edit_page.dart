import 'dart:async';

import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import '../../widgets/entity_context_menu.dart';
import '../../widgets/entity_list.dart';
import 'name_util.dart';
import 'swipe_actions.dart';
import 'voice_edit_page.dart';
import 'voice_param_editor.dart';

class CharacterEditPage extends StatefulWidget {
  const CharacterEditPage({
    super.key,
    required this.characterName,
    this.embedded = false,
    this.selectedVoiceName,
    this.onOpenVoice,
    this.onDeselectVoice,
    this.onRenamed,
    this.onUpdated,
  });

  final String characterName;

  /// When true, render only the editor body (no Scaffold / AppBar).
  final bool embedded;

  /// Highlighted voice in the list (wide 3-col selection).
  final String? selectedVoiceName;

  /// If set, called instead of pushing [VoiceEditPage].
  final ValueChanged<String>? onOpenVoice;

  /// Called when the currently selected voice is deleted (wide layout).
  final VoidCallback? onDeselectVoice;

  final ValueChanged<String>? onRenamed;
  final VoidCallback? onUpdated;

  @override
  State<CharacterEditPage> createState() => _CharacterEditPageState();
}

class _CharacterEditPageState extends State<CharacterEditPage> {
  final _name = TextEditingController();
  final _description = TextEditingController();
  Timer? _saveTimer;
  String _originalName = '';
  List<Voice> _voices = const [];
  bool _loading = true;
  bool _ready = false;
  bool _saving = false;
  bool _saveQueued = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _originalName = widget.characterName;
    _name.addListener(_onMetaChanged);
    _description.addListener(_onMetaChanged);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_load());
    });
  }

  @override
  void didUpdateWidget(covariant CharacterEditPage oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.characterName != widget.characterName) {
      _flushSave();
      _originalName = widget.characterName;
      unawaited(_load());
    }
  }

  @override
  void dispose() {
    _saveTimer?.cancel();
    _name.removeListener(_onMetaChanged);
    _description.removeListener(_onMetaChanged);
    _name.dispose();
    _description.dispose();
    super.dispose();
  }

  void _onMetaChanged() {
    if (!_ready) return;
    _scheduleSave();
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
      _originalName = c.name;
      _name.text = c.name;
      _description.text = c.description;
      setState(() {
        _voices = List.of(c.voices ?? const <Voice>[]);
        _loading = false;
        _ready = true;
      });
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
    final name = _name.text.trim();
    if (name.isEmpty) return;
    if (_saving) {
      _saveQueued = true;
      return;
    }
    _saving = true;
    final voices = List<Voice>.of(_voices);
    final description = _description.text.trim();
    final oldName = _originalName;
    final client = AppScope.of(context);
    try {
      await client.characters.update(
        CharacterUpdateReq(
          oldName: oldName,
          character: Character(
            name: name,
            description: description,
            voices: voices,
          ),
        ),
      );
      if (!mounted) return;
      final renamed = oldName != name;
      _originalName = name;
      if (renamed) widget.onRenamed?.call(name);
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

  Future<void> _openVoice(String voiceName) async {
    _flushSave();
    final open = widget.onOpenVoice;
    if (open != null) {
      open(voiceName);
      return;
    }
    await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => VoiceEditPage(
          characterName: _originalName,
          voiceName: voiceName,
        ),
      ),
    );
    if (mounted) await _load();
  }

  Future<void> _addVoice() async {
    final nameCtrl = TextEditingController(
      text: uniqueCopyName('voice', _voices.map((v) => v.name)),
    );
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('添加音色'),
        content: TextField(
          controller: nameCtrl,
          decoration: const InputDecoration(labelText: '名称'),
          autofocus: true,
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('添加')),
        ],
      ),
    );
    if (ok != true) return;
    final name = nameCtrl.text.trim();
    if (name.isEmpty) return;
    if (_voices.any((v) => v.name == name)) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('音色名已存在')),
      );
      return;
    }
    setState(() {
      _voices = [
        ..._voices,
        Voice(name: name, params: [defaultSynthParam(1)]),
      ];
    });
    _scheduleSave(delay: Duration.zero);
    await _openVoice(name);
  }

  Future<bool> _cloneVoiceAt(int index) async {
    final src = _voices[index];
    final newName = uniqueCopyName(src.name, _voices.map((v) => v.name));
    setState(() {
      _voices = [..._voices, cloneVoice(src, newName)];
    });
    _scheduleSave(delay: Duration.zero);
    widget.onUpdated?.call();
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('已克隆为「$newName」')),
      );
    }
    return false;
  }

  Future<bool> _deleteVoiceAt(int index) async {
    if (index < 0 || index >= _voices.length) return false;
    final deleted = _voices[index].name;
    setState(() {
      _voices = List.of(_voices)..removeAt(index);
    });
    _scheduleSave(delay: Duration.zero);
    if (widget.selectedVoiceName == deleted) {
      widget.onDeselectVoice?.call();
    }
    widget.onUpdated?.call();
    return true;
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
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (widget.embedded) ...[
          PageCardTitle(_originalName),
          const SizedBox(height: 16),
        ],
        TextField(
          controller: _name,
          decoration: InputDecoration(
            labelText: '名称',
            filled: true,
            fillColor: cs.surfaceContainer,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(16),
              borderSide: BorderSide.none,
            ),
          ),
        ),
        const SizedBox(height: 12),
        TextField(
          controller: _description,
          minLines: 2,
          maxLines: 4,
          decoration: InputDecoration(
            labelText: '描述',
            filled: true,
            fillColor: cs.surfaceContainer,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(16),
              borderSide: BorderSide.none,
            ),
          ),
        ),
        const SizedBox(height: 20),
        Row(
          children: [
            Text('音色', style: Theme.of(context).textTheme.titleMedium),
            const Spacer(),
            IconButton(
              tooltip: '添加音色',
              onPressed: _addVoice,
              icon: const FaIcon(FontAwesomeIcons.plus, size: 16),
            ),
          ],
        ),
        const SizedBox(height: 8),
        EntityDismissibleList(
          itemCount: _voices.length,
          style: entitySwipeStyle(
            cs,
            color: cs.surfaceContainerHighest,
          ),
          emptyBuilder: (_) => Padding(
            padding: const EdgeInsets.all(24),
            child: Text(
              '暂无音色',
              textAlign: TextAlign.center,
              style: TextStyle(color: cs.onSurfaceVariant),
            ),
          ),
          onDismiss: (index, direction) async {
            if (direction == DismissDirection.startToEnd) {
              return _cloneVoiceAt(index);
            }
            if (direction == DismissDirection.endToStart) {
              return _deleteVoiceAt(index);
            }
            return false;
          },
          onTap: (index) => unawaited(_openVoice(_voices[index].name)),
          menuChildren: (context, menuIndex, close) {
            final name =
                (menuIndex != null &&
                    menuIndex >= 0 &&
                    menuIndex < _voices.length)
                ? _voices[menuIndex].name
                : null;
            return buildEntityContextMenu(
              groupLabel: '音色',
              name: name,
              close: close,
              onEdit: () {
                if (menuIndex != null) {
                  unawaited(_openVoice(_voices[menuIndex].name));
                }
              },
              onClone: () {
                if (menuIndex != null) unawaited(_cloneVoiceAt(menuIndex));
              },
              onDelete: () {
                if (menuIndex != null) unawaited(_deleteVoiceAt(menuIndex));
              },
            );
          },
          itemBuilder: (context, index) {
            final v = _voices[index];
            return EntityListRow(
              leading: const EntityListAvatar(Icons.graphic_eq_rounded),
              selected: widget.selectedVoiceName == v.name,
              title: v.name,
              subtitle: EntityListMeta(voiceSubtitle(v)),
            );
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
      return body;
    }

    return PopScope(
      onPopInvokedWithResult: (didPop, _) {
        if (didPop) _flushSave();
      },
      child: Scaffold(
        backgroundColor: cs.surface,
        appBar: AppBar(
          title: Text(_loading ? '编辑角色' : _originalName),
        ),
        body: _loading || _error != null
            ? body
            : SingleChildScrollView(
                padding: const EdgeInsets.only(top: 8, bottom: 32),
                child: ContentFrame(child: body),
              ),
      ),
    );
  }
}
