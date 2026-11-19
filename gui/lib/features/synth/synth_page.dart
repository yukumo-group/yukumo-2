import 'dart:async';
import 'dart:convert';

import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import '../shell/bar_padding.dart';
import 'wav_preview.dart';

class SynthPage extends StatefulWidget {
  const SynthPage({super.key});

  @override
  State<SynthPage> createState() => _SynthPageState();
}

class _SynthPageState extends State<SynthPage> {
  final _text = TextEditingController();
  final _characterDropdown = M3EDropdownController<String>();
  final _voiceDropdown = M3EDropdownController<String>();

  List<Character> _characters = const [];
  String? _character;
  String? _voice;
  bool _busy = false;
  String? _status;
  bool _syncing = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_reloadCharacters());
    });
  }

  @override
  void dispose() {
    _text.dispose();
    _characterDropdown.dispose();
    _voiceDropdown.dispose();
    super.dispose();
  }

  Future<void> _reloadCharacters() async {
    try {
      final res = await AppScope.of(context).characters.list();
      if (!mounted) return;
      final list = res.characters ?? const <Character>[];
      setState(() {
        _characters = list;
        _character ??= list.isNotEmpty ? list.first.name : null;
        _voice = _voiceFor(_character);
      });
      await _syncDropdowns();
    } catch (e) {
      if (!mounted) return;
      setState(() => _status = '加载失败: $e');
    }
  }

  String? _voiceFor(String? character) {
    if (character == null) return null;
    for (final c in _characters) {
      if (c.name == character) {
        final voices = c.voices ?? const <Voice>[];
        return voices.isNotEmpty ? voices.first.name : null;
      }
    }
    return null;
  }

  List<Voice> _voicesOf(String? character) {
    if (character == null) return const [];
    for (final c in _characters) {
      if (c.name == character) return c.voices ?? const [];
    }
    return const [];
  }

  Future<List<M3EDropdownItem<String>>> _characterItems() async {
    return [
      for (final c in _characters)
        M3EDropdownItem<String>(
          value: c.name,
          label: c.name,
          selected: c.name == _character,
        ),
    ];
  }

  Future<List<M3EDropdownItem<String>>> _voiceItems() async {
    final voices = _voicesOf(_character);
    return [
      for (final v in voices)
        M3EDropdownItem<String>(
          value: v.name,
          label: v.name,
          selected: v.name == _voice,
        ),
    ];
  }

  Future<void> _syncDropdowns() async {
    if (!mounted) return;
    _syncing = true;
    try {
      if (_characterDropdown.initialized) {
        _characterDropdown.setItems(await _characterItems());
      }
      if (_voiceDropdown.initialized) {
        _voiceDropdown.setItems(await _voiceItems());
      }
    } finally {
      _syncing = false;
    }
  }

  int _versionForVoice() {
    final voices = _voicesOf(_character);
    for (final v in voices) {
      if (v.name == _voice) {
        final params = v.params;
        if (params != null && params.isNotEmpty) return params.first.version;
      }
    }
    return 1;
  }

  Future<void> _runSynth({required bool save}) async {
    final text = _text.text.trim();
    if (text.isEmpty) {
      setState(() => _status = '请输入文本');
      return;
    }
    if (_character == null || _voice == null) {
      setState(() => _status = '请选择角色和音色');
      return;
    }
    setState(() {
      _busy = true;
      _status = save ? '保存中…' : '预览中…';
    });
    try {
      final res = await AppScope.of(context).synth.render(
        SynthRequest(
          voiceRef: '$_character:$_voice',
          text: text,
          version: _versionForVoice(),
          speedFactor: 1,
          preview: false,
          noWrite: !save,
        ),
      );
      if (!mounted) return;
      if (!save) {
        final b64 = res.wAV;
        if (b64 == null || b64.isEmpty) {
          setState(() => _status = '预览无音频数据');
          return;
        }
        setState(() => _status = '播放中…');
        // Start playback without awaiting completion so the UI stays responsive.
        unawaited(() async {
          try {
            await playWavBytes(base64Decode(b64));
            if (!mounted) return;
            setState(() => _status = '预览完成');
          } catch (e) {
            if (!mounted) return;
            setState(() => _status = '播放失败: $e');
          }
        }());
      } else {
        setState(() {
          _status = res.path != null && res.path!.isNotEmpty
              ? '已保存: ${res.path}'
              : (res.wrote ? '已保存' : '合成完成（未落盘）');
        });
      }
    } catch (e) {
      if (!mounted) return;
      setState(() => _status = '失败: $e');
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Widget _dropdown({
    required ColorScheme cs,
    required String label,
    required String hint,
    required String searchHint,
    required M3EDropdownController<String> controller,
    required Future<List<M3EDropdownItem<String>>> Function() future,
    required ValueChanged<String?> onChanged,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(label, style: Theme.of(context).textTheme.labelLarge),
        const SizedBox(height: 8),
        NotificationListener<ScrollNotification>(
          onNotification: (_) => true,
          child: M3EDropdownMenu<String>.future(
            future: future,
            controller: controller,
            singleSelect: true,
            searchEnabled: true,
            maxSelections: 1,
            containerRadius: 24,
            fieldStyle: M3EDropdownFieldStyle(
              hintText: hint,
              showClearIcon: false,
              padding: const EdgeInsets.all(14),
            ),
            dropdownStyle: M3EDropdownStyle(
              elevation: 4,
              maxHeight: 320,
              backgroundColor: cs.surfaceContainer,
            ),
            searchStyle: M3ESearchStyle(
              hintText: searchHint,
              showClearIcon: true,
            ),
            itemStyle: M3EDropdownItemStyle(
              outerRadius: 16,
              innerRadius: 6,
              selectedBackgroundColor: cs.secondaryContainer,
            ),
            onSelectionChanged: (items) {
              if (_syncing) return;
              onChanged(items.isEmpty ? null : items.first.value);
            },
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final top = MediaQuery.paddingOf(context).top;
    return withBarBottomPadding(context, (bottom) {
      return CustomScrollView(
        slivers: [
          SliverPadding(
            padding: EdgeInsets.only(top: top + 72, bottom: bottom),
            sliver: SliverToBoxAdapter(
              child: ContentFrame(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const PageCardTitle('合成'),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _text,
                      minLines: 8,
                      maxLines: 16,
                      decoration: InputDecoration(
                        labelText: '文本',
                        hintText: '输入要合成的日文假名或文本',
                        alignLabelWithHint: true,
                        filled: true,
                        fillColor: cs.surfaceContainer,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(20),
                          borderSide: BorderSide.none,
                        ),
                      ),
                    ),
                    const SizedBox(height: 20),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                          child: _dropdown(
                            cs: cs,
                            label: '角色',
                            hint: '选择角色',
                            searchHint: '搜索角色',
                            controller: _characterDropdown,
                            future: _characterItems,
                            onChanged: (name) {
                              setState(() {
                                _character = name;
                                _voice = _voiceFor(name);
                              });
                              unawaited(_syncDropdowns());
                            },
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: _dropdown(
                            cs: cs,
                            label: '音色',
                            hint: '选择音色',
                            searchHint: '搜索音色',
                            controller: _voiceDropdown,
                            future: _voiceItems,
                            onChanged: (name) {
                              setState(() => _voice = name);
                            },
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 24),
                    Row(
                      children: [
                        Expanded(
                          child: M3EButton(
                            onPressed: _busy ? null : () => _runSynth(save: false),
                            style: M3EButtonStyle.tonal,
                            shape: M3EButtonShape.round,
                            child: const Text('预览'),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: M3EButton(
                            onPressed: _busy ? null : () => _runSynth(save: true),
                            style: M3EButtonStyle.filled,
                            shape: M3EButtonShape.round,
                            child: const Text('保存'),
                          ),
                        ),
                      ],
                    ),
                    if (_busy) ...[
                      const SizedBox(height: 16),
                      const Center(child: M3ELoadingIndicator()),
                    ],
                    if (_status != null) ...[
                      const SizedBox(height: 16),
                      Text(
                        _status!,
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: cs.onSurfaceVariant,
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ),
        ],
      );
    });
  }
}
