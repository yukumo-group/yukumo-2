import 'dart:async';

import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import 'name_util.dart';
import 'swipe_actions.dart';

class CharacterEditPage extends StatefulWidget {
  const CharacterEditPage({super.key, required this.characterName});

  final String characterName;

  @override
  State<CharacterEditPage> createState() => _CharacterEditPageState();
}

class _CharacterEditPageState extends State<CharacterEditPage> {
  final _name = TextEditingController();
  final _description = TextEditingController();
  String _originalName = '';
  List<Voice> _voices = const [];
  bool _loading = true;
  bool _saving = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _originalName = widget.characterName;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_load());
    });
  }

  @override
  void dispose() {
    _name.dispose();
    _description.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
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
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = '$e';
        _loading = false;
      });
    }
  }

  Future<void> _save() async {
    final name = _name.text.trim();
    if (name.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('名称不能为空')),
      );
      return;
    }
    setState(() => _saving = true);
    final client = AppScope.of(context);
    try {
      await client.characters.update(
        CharacterUpdateReq(
          oldName: _originalName,
          character: Character(
            name: name,
            description: _description.text.trim(),
            voices: _voices,
          ),
        ),
      );
      if (!mounted) return;
      setState(() {
        _originalName = name;
        _saving = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('已保存')),
      );
    } catch (e) {
      if (!mounted) return;
      setState(() => _saving = false);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('保存失败: $e')),
      );
    }
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
        Voice(
          name: name,
          params: const [
            SynthParam(
              version: 1,
              preset: 0,
              speed: 100,
              base: 0,
              volume: 0,
              pitch: 0,
              accent: 0,
              lmd: 0,
              fsc: 0,
            ),
          ],
        ),
      ];
    });
  }

  Future<bool> _cloneVoiceAt(int index) async {
    final src = _voices[index];
    final newName = uniqueCopyName(src.name, _voices.map((v) => v.name));
    setState(() {
      _voices = [
        ..._voices,
        Voice(
          name: newName,
          params: [
            for (final p in src.params ?? const <SynthParam>[])
              SynthParam(
                version: p.version,
                preset: p.preset,
                speed: p.speed,
                phontData: p.phontData,
                base: p.base,
                volume: p.volume,
                pitch: p.pitch,
                accent: p.accent,
                lmd: p.lmd,
                fsc: p.fsc,
              ),
          ],
        ),
      ];
    });
    return false;
  }

  Future<bool> _deleteVoiceAt(int index) async {
    final v = _voices[index];
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('删除音色'),
        content: Text('确定删除「${v.name}」？'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('删除')),
        ],
      ),
    );
    if (ok != true) return false;
    setState(() => _voices = List.of(_voices)..removeAt(index));
    return true;
  }

  Future<void> _editVoiceAt(int index) async {
    final v = _voices[index];
    final nameCtrl = TextEditingController(text: v.name);
    final speed = v.params?.isNotEmpty == true ? v.params!.first.speed : 100;
    final speedCtrl = TextEditingController(text: '$speed');
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('编辑音色'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameCtrl,
              decoration: const InputDecoration(labelText: '名称'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: speedCtrl,
              decoration: const InputDecoration(labelText: '语速 (50–300)'),
              keyboardType: TextInputType.number,
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('应用')),
        ],
      ),
    );
    if (ok != true) return;
    final newName = nameCtrl.text.trim();
    if (newName.isEmpty) return;
    if (newName != v.name && _voices.any((x) => x.name == newName)) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('音色名已存在')),
      );
      return;
    }
    final parsedSpeed = int.tryParse(speedCtrl.text.trim()) ?? speed;
    final clamped = parsedSpeed.clamp(50, 300);
    setState(() {
      _voices = [
        for (var i = 0; i < _voices.length; i++)
          if (i == index)
            Voice(
              name: newName,
              params: [
                for (var j = 0; j < (v.params?.length ?? 0); j++)
                  SynthParam(
                    version: v.params![j].version,
                    preset: v.params![j].preset,
                    speed: j == 0 ? clamped : v.params![j].speed,
                    phontData: v.params![j].phontData,
                    base: v.params![j].base,
                    volume: v.params![j].volume,
                    pitch: v.params![j].pitch,
                    accent: v.params![j].accent,
                    lmd: v.params![j].lmd,
                    fsc: v.params![j].fsc,
                  ),
              ],
            )
          else
            _voices[i],
      ];
    });
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Scaffold(
      backgroundColor: cs.surface,
      appBar: AppBar(
        title: Text(_loading ? '编辑角色' : _originalName),
        actions: [
          TextButton(
            onPressed: _saving || _loading ? null : _save,
            child: _saving
                ? const SizedBox(
                    width: 18,
                    height: 18,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('保存'),
          ),
        ],
      ),
      body: _loading
          ? const Center(child: M3ELoadingIndicator())
          : _error != null
              ? Center(child: Text(_error!))
              : LayoutBuilder(
                  builder: (context, constraints) {
                    final wide = constraints.maxWidth >= kPageBreakpoint;
                    final form = Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
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
                      ],
                    );
                    final voicesHeader = Row(
                      children: [
                        Text('音色', style: Theme.of(context).textTheme.titleMedium),
                        const Spacer(),
                        IconButton(
                          tooltip: '添加音色',
                          onPressed: _addVoice,
                          icon: const FaIcon(FontAwesomeIcons.plus, size: 16),
                        ),
                      ],
                    );
                    final voicesList = M3EDismissibleCardColumn(
                      itemCount: _voices.length,
                      emptyBuilder: (_) => Padding(
                        padding: const EdgeInsets.all(24),
                        child: Text(
                          '暂无音色',
                          textAlign: TextAlign.center,
                          style: TextStyle(color: cs.onSurfaceVariant),
                        ),
                      ),
                      style: characterSwipeStyle(cs),
                      onDismiss: (index, direction) async {
                        if (direction == DismissDirection.startToEnd) {
                          return _cloneVoiceAt(index);
                        }
                        if (direction == DismissDirection.endToStart) {
                          return _deleteVoiceAt(index);
                        }
                        return false;
                      },
                      onTap: _editVoiceAt,
                      itemBuilder: (context, index) {
                        final v = _voices[index];
                        final params = v.params ?? const <SynthParam>[];
                        final summary = params.isEmpty
                            ? '无参数'
                            : 'v${params.first.version} · speed ${params.first.speed}';
                        return Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(v.name, style: Theme.of(context).textTheme.titleSmall),
                            const SizedBox(height: 4),
                            Text(
                              summary,
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: cs.onSurfaceVariant,
                              ),
                            ),
                          ],
                        );
                      },
                    );

                    final content = wide
                        ? Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Expanded(child: form),
                              const SizedBox(width: 24),
                              Expanded(
                                child: Column(
                                  children: [
                                    voicesHeader,
                                    voicesList,
                                  ],
                                ),
                              ),
                            ],
                          )
                        : Column(
                            crossAxisAlignment: CrossAxisAlignment.stretch,
                            children: [
                              form,
                              const SizedBox(height: 20),
                              voicesHeader,
                              voicesList,
                            ],
                          );

                    return SingleChildScrollView(
                      padding: const EdgeInsets.only(top: 8, bottom: 32),
                      child: ContentFrame(child: content),
                    );
                  },
                ),
    );
  }
}
