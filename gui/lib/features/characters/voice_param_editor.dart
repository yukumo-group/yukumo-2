import 'dart:async';

import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/int_slider_field.dart';

const _kVersions = [1, 2, 10];
const _kBaseNames = ['F1E', 'F2E', 'M1E'];
const _kConnectedGap = 2.0;

SynthParam defaultSynthParam(int version) {
  switch (version) {
    case 10:
      return const SynthParam(
        version: 10,
        preset: 0,
        base: 0,
        pitch: 100,
        accent: 100,
        lmd: 100,
      );
    case 2:
      return const SynthParam(
        version: 2,
        preset: 0,
        base: 0,
        pitch: 0,
        accent: 0,
        lmd: 0,
      );
    default:
      return const SynthParam(
        version: 1,
        preset: 0,
        base: 0,
        pitch: 0,
        accent: 0,
        lmd: 0,
      );
  }
}

class VoiceParamEditor extends StatefulWidget {
  const VoiceParamEditor({
    super.key,
    required this.voice,
    required this.onChanged,
    this.onSelectedVersionChanged,
  });

  final Voice voice;
  final ValueChanged<Voice> onChanged;
  final ValueChanged<int>? onSelectedVersionChanged;

  @override
  State<VoiceParamEditor> createState() => _VoiceParamEditorState();
}

class _VoiceParamEditorState extends State<VoiceParamEditor> {
  late TextEditingController _name;
  final _presetDropdown = M3EDropdownController<int>();
  int _tabIndex = 0;
  List<String> _presetNames = const [];
  bool _syncingPreset = false;

  @override
  void initState() {
    super.initState();
    _name = TextEditingController(text: widget.voice.name);
    _tabIndex = _initialTab();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      widget.onSelectedVersionChanged?.call(_selectedVersion);
      unawaited(_loadPresets());
    });
  }

  @override
  void didUpdateWidget(covariant VoiceParamEditor oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.voice.name != widget.voice.name &&
        _name.text != widget.voice.name) {
      _name.text = widget.voice.name;
    }
    final oldPreset = _paramOf(oldWidget.voice, _selectedVersion)?.preset;
    final nextPreset = _currentParam?.preset;
    if (oldPreset != nextPreset ||
        oldWidget.voice.params != widget.voice.params) {
      unawaited(_syncPresetDropdown());
    }
  }

  @override
  void dispose() {
    _name.dispose();
    _presetDropdown.dispose();
    super.dispose();
  }

  int _initialTab() {
    final params = widget.voice.params ?? const <SynthParam>[];
    if (params.isEmpty) return 0;
    final v = params.first.version;
    final i = _kVersions.indexOf(v);
    return i < 0 ? 0 : i;
  }

  int get _selectedVersion => _kVersions[_tabIndex.clamp(0, _kVersions.length - 1)];

  SynthParam? _paramOf(Voice voice, int version) {
    final params = voice.params ?? const <SynthParam>[];
    for (final p in params) {
      if (p.version == version) return p;
    }
    return null;
  }

  SynthParam? get _currentParam => _paramOf(widget.voice, _selectedVersion);

  Future<void> _loadPresets() async {
    try {
      final res = await AppScope.of(context).voices.presets(
        VoicePresetsReq(version: _selectedVersion),
      );
      if (!mounted) return;
      setState(() => _presetNames = res.names ?? const []);
      await _syncPresetDropdown();
    } catch (_) {
      if (!mounted) return;
      setState(() => _presetNames = const []);
      await _syncPresetDropdown();
    }
  }

  List<M3EDropdownItem<int>> _presetItems(int selected) {
    if (_presetNames.isEmpty) {
      return [
        for (var i = 0; i <= 20; i++)
          M3EDropdownItem<int>(
            value: i,
            label: 'Preset $i',
            selected: i == selected,
          ),
      ];
    }
    return [
      for (var i = 0; i < _presetNames.length; i++)
        M3EDropdownItem<int>(
          value: i,
          label: _presetNames[i],
          selected: i == selected,
        ),
    ];
  }

  Future<List<M3EDropdownItem<int>>> _presetItemsFuture() async {
    return _presetItems(_currentParam?.preset ?? 0);
  }

  Future<void> _syncPresetDropdown() async {
    if (!_presetDropdown.initialized) return;
    _syncingPreset = true;
    try {
      _presetDropdown.setItems(_presetItems(_currentParam?.preset ?? 0));
    } finally {
      _syncingPreset = false;
    }
  }

  void _emit(Voice next) => widget.onChanged(next);

  Voice _voiceWith({
    String? name,
    int? speed,
    List<SynthParam>? params,
  }) {
    return Voice(
      name: name ?? widget.voice.name,
      speed: speed ?? widget.voice.speed,
      params: params ?? widget.voice.params,
    );
  }

  SynthParam _copy(
    SynthParam p, {
    int? preset,
    int? base,
    int? pitch,
    int? accent,
    int? lmd,
  }) {
    return SynthParam(
      version: p.version,
      preset: preset ?? p.preset,
      phontData: p.phontData,
      base: base ?? p.base,
      pitch: pitch ?? p.pitch,
      accent: accent ?? p.accent,
      lmd: lmd ?? p.lmd,
    );
  }

  void _replaceParam(SynthParam next) {
    final params = List<SynthParam>.of(widget.voice.params ?? const []);
    final i = params.indexWhere((p) => p.version == next.version);
    if (i >= 0) {
      params[i] = next;
    } else {
      params.add(next);
    }
    params.sort((a, b) => a.version.compareTo(b.version));
    _emit(_voiceWith(params: params));
  }

  void _setVersionEnabled(bool enabled) {
    if (enabled) {
      if (_currentParam != null) return;
      _replaceParam(defaultSynthParam(_selectedVersion));
      return;
    }
    final params = List<SynthParam>.of(widget.voice.params ?? const []);
    if (params.length <= 1) return;
    params.removeWhere((p) => p.version == _selectedVersion);
    _emit(_voiceWith(params: params));
  }

  void _commitName() {
    final name = _name.text.trim();
    if (name.isEmpty || name == widget.voice.name) {
      _name.text = widget.voice.name;
      return;
    }
    _emit(_voiceWith(name: name));
  }

  Widget _versionTabs() {
    return LayoutBuilder(
      builder: (context, constraints) {
        final n = _kVersions.length;
        final width =
            (constraints.maxWidth - _kConnectedGap * (n - 1)) / n;
        return M3EToggleButtonGroup(
          type: M3EButtonGroupType.connected,
          style: M3EButtonStyle.tonal,
          size: M3EButtonSize.sm,
          shape: M3EButtonShape.round,
          density: M3EButtonGroupDensity.regular,
          overflow: M3EButtonGroupOverflow.none,
          selectedIndex: _tabIndex,
          onSelectedIndexChanged: (idx) {
            if (idx == null) return;
            setState(() => _tabIndex = idx);
            widget.onSelectedVersionChanged?.call(_kVersions[idx]);
            unawaited(_loadPresets());
          },
          actions: [
            for (var i = 0; i < _kVersions.length; i++)
              M3EToggleButtonGroupAction(
                icon: Icon(
                  _paramOf(widget.voice, _kVersions[i]) != null
                      ? Icons.check_circle_rounded
                      : Icons.block_rounded,
                  size: 16,
                ),
                label: Text(const ['AQ1', 'AQ2', 'AQ10'][i]),
                width: width.isFinite ? width : null,
              ),
          ],
        );
      },
    );
  }

  Widget _enableSwitch(ColorScheme cs) {
    final enabled = _currentParam != null;
    final canDisable = (widget.voice.params?.length ?? 0) > 1;
    return Theme(
      data: Theme.of(context).copyWith(
        splashColor: Colors.transparent,
        highlightColor: Colors.transparent,
        hoverColor: Colors.transparent,
        splashFactory: NoSplash.splashFactory,
      ),
      child: SwitchListTile(
        contentPadding: EdgeInsets.zero,
        title: Text('启用 AquesTalk $_selectedVersion'),
        value: enabled,
        onChanged: (!enabled || canDisable) ? _setVersionEnabled : null,
      ),
    );
  }

  Widget _presetDropdownField(SynthParam p, ColorScheme cs) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text('Preset', style: Theme.of(context).textTheme.labelLarge),
        const SizedBox(height: 8),
        NotificationListener<ScrollNotification>(
          onNotification: (_) => true,
          child: M3EDropdownMenu<int>.future(
            key: ValueKey('preset-$_selectedVersion-${_presetNames.length}'),
            future: _presetItemsFuture,
            controller: _presetDropdown,
            singleSelect: true,
            searchEnabled: true,
            maxSelections: 1,
            showChipAnimation: false,
            containerRadius: 24,
            fieldStyle: M3EDropdownFieldStyle(
              hintText: '选择 Preset',
              showClearIcon: false,
              padding: const EdgeInsets.all(14),
            ),
            dropdownStyle: M3EDropdownStyle(
              elevation: 4,
              maxHeight: 320,
              backgroundColor: cs.surfaceContainer,
            ),
            searchStyle: const M3ESearchStyle(
              hintText: '搜索 Preset',
              showClearIcon: true,
            ),
            itemStyle: M3EDropdownItemStyle(
              outerRadius: 16,
              innerRadius: 6,
              selectedBackgroundColor: cs.secondaryContainer,
            ),
            onSelectionChanged: (items) {
              if (_syncingPreset || items.isEmpty) return;
              final next = items.first.value;
              if (next == p.preset) return;
              _replaceParam(_copy(p, preset: next));
            },
          ),
        ),
      ],
    );
  }

  Widget _paramFields(SynthParam p, ColorScheme cs) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _presetDropdownField(p, cs),
        // AQ10_CUSTOM == 0; named presets ignore base/pitch/accent/lmd.
        if (p.version == 10 && p.preset == 0) ...[
          const SizedBox(height: 16),
          Text('基本素片', style: Theme.of(context).textTheme.labelLarge),
          const SizedBox(height: 8),
          LayoutBuilder(
            builder: (context, constraints) {
              final n = _kBaseNames.length;
              final width =
                  (constraints.maxWidth - _kConnectedGap * (n - 1)) / n;
              return M3EToggleButtonGroup(
                type: M3EButtonGroupType.connected,
                style: M3EButtonStyle.tonal,
                size: M3EButtonSize.xs,
                shape: M3EButtonShape.round,
                overflow: M3EButtonGroupOverflow.none,
                selectedIndex: p.base.clamp(0, 2),
                onSelectedIndexChanged: (idx) {
                  if (idx == null) return;
                  _replaceParam(_copy(p, base: idx));
                },
                actions: [
                  for (final name in _kBaseNames)
                    M3EToggleButtonGroupAction(
                      label: Text(name),
                      width: width.isFinite ? width : null,
                    ),
                ],
              );
            },
          ),
          const SizedBox(height: 16),
          IntSliderField(
            label: '音高',
            value: p.pitch.clamp(20, 200),
            min: 20,
            max: 200,
            onChanged: (v) => _replaceParam(_copy(p, pitch: v)),
          ),
          const SizedBox(height: 16),
          IntSliderField(
            label: '重音',
            value: p.accent.clamp(0, 200),
            min: 0,
            max: 200,
            onChanged: (v) => _replaceParam(_copy(p, accent: v)),
          ),
          const SizedBox(height: 16),
          IntSliderField(
            label: 'LMD',
            value: p.lmd.clamp(0, 200),
            min: 0,
            max: 200,
            onChanged: (v) => _replaceParam(_copy(p, lmd: v)),
          ),
        ],
        if (p.version == 2) ...[
          const SizedBox(height: 12),
          Text(
            'V2 phont 数据请在后续版本通过文件导入；当前可先调 preset。',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
              color: cs.onSurfaceVariant,
            ),
          ),
        ],
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final param = _currentParam;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(
          controller: _name,
          decoration: InputDecoration(
            labelText: '音色名称',
            filled: true,
            fillColor: cs.surfaceContainer,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide.none,
            ),
          ),
          onEditingComplete: _commitName,
          onSubmitted: (_) => _commitName(),
        ),
        const SizedBox(height: 16),
        IntSliderField(
          label: '语速',
          value: widget.voice.speed.clamp(50, 300),
          min: 50,
          max: 300,
          onChanged: (v) => _emit(_voiceWith(speed: v)),
        ),
        const SizedBox(height: 16),
        _versionTabs(),
        const SizedBox(height: 8),
        _enableSwitch(cs),
        const SizedBox(height: 8),
        if (param == null)
          Text(
            '尚未启用 AQ$_selectedVersion 参数，打开上方开关以添加。',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: cs.onSurfaceVariant,
            ),
          )
        else
          _paramFields(param, cs),
      ],
    );
  }
}

String voiceSubtitle(Voice v) {
  final params = v.params ?? const <SynthParam>[];
  if (params.isEmpty) return '无参数';
  return params.map((p) => 'AQ${p.version}').join(' · ');
}
