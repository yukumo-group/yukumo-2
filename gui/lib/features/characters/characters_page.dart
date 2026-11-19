import 'dart:async';

import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import '../../widgets/entity_context_menu.dart';
import '../../widgets/entity_list.dart';
import '../../widgets/multi_column_view.dart';
import '../shell/bar_padding.dart';
import 'character_edit_page.dart';
import 'name_util.dart';
import 'swipe_actions.dart';
import 'voice_edit_page.dart';

class CharactersPage extends StatefulWidget {
  const CharactersPage({super.key});

  @override
  State<CharactersPage> createState() => _CharactersPageState();
}

class _CharactersPageState extends State<CharactersPage> {
  List<Character> _items = const [];
  bool _loading = true;
  String? _error;
  String? _selectedCharacter;
  String? _selectedVoice;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_reload());
    });
  }

  Future<void> _reload({bool keepSelection = true}) async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final res = await AppScope.of(context).characters.list();
      if (!mounted) return;
      final list = List.of(res.characters ?? const <Character>[]);
      String? selected = keepSelection ? _selectedCharacter : null;
      if (selected != null && !list.any((c) => c.name == selected)) {
        selected = null;
      }
      selected ??= list.isNotEmpty ? list.first.name : null;
      String? voice = keepSelection ? _selectedVoice : null;
      if (selected == null) {
        voice = null;
      } else if (voice != null) {
        Character? match;
        for (final c in list) {
          if (c.name == selected) {
            match = c;
            break;
          }
        }
        final voices = match?.voices ?? const <Voice>[];
        if (!voices.any((v) => v.name == voice)) voice = null;
      }
      setState(() {
        _items = list;
        _selectedCharacter = selected;
        _selectedVoice = voice;
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

  Future<void> _openEditMobile(String name) async {
    await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => CharacterEditPage(characterName: name),
      ),
    );
    if (mounted) await _reload();
  }

  void _selectCharacter(String name, {required bool dual}) {
    if (dual) {
      setState(() {
        _selectedCharacter = name;
        _selectedVoice = null;
      });
      return;
    }
    unawaited(_openEditMobile(name));
  }

  Future<bool> _cloneAt(int index) async {
    final src = _items[index];
    final client = AppScope.of(context);
    final newName = uniqueCopyName(src.name, _items.map((c) => c.name));
    try {
      await client.characters.add(
        CharacterAddReq(character: cloneCharacter(src, newName)),
      );
      if (!mounted) return false;
      await _reload();
      if (mounted) {
        setState(() {
          _selectedCharacter = newName;
          _selectedVoice = null;
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('已克隆为「$newName」')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('克隆失败: $e')),
        );
      }
    }
    return false;
  }

  Future<bool> _deleteAt(int index) async {
    if (index < 0 || index >= _items.length) return false;
    final c = _items[index];
    final client = AppScope.of(context);
    try {
      await client.characters.delete(CharacterDeleteReq(name: c.name));
      if (!mounted) return false;
      final removingSelected = _selectedCharacter == c.name;
      setState(() {
        _items = List.of(_items)..removeAt(index);
        if (removingSelected) {
          // Close the detail pane / page for the deleted character.
          _selectedCharacter = null;
          _selectedVoice = null;
        }
      });
      return true;
    } catch (e) {
      if (!mounted) return false;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('删除失败: $e')),
      );
      return false;
    }
  }

  Future<void> _showAddDialog({required bool dual}) async {
    final client = AppScope.of(context);
    final nameCtrl = TextEditingController();
    final descCtrl = TextEditingController();
    final added = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('新建角色'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameCtrl,
              decoration: const InputDecoration(labelText: '名称'),
              autofocus: true,
            ),
            const SizedBox(height: 12),
            TextField(
              controller: descCtrl,
              decoration: const InputDecoration(labelText: '描述'),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('创建')),
        ],
      ),
    );
    if (added != true) return;
    final name = nameCtrl.text.trim();
    if (name.isEmpty) return;
    try {
      await client.characters.add(
        CharacterAddReq(
          character: Character(name: name, description: descCtrl.text.trim()),
        ),
      );
      if (!mounted) return;
      await _reload();
      if (!mounted) return;
      _selectCharacter(name, dual: dual);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('创建失败: $e')),
      );
    }
  }

  String _characterMeta(Character c) {
    final voiceCount = c.voices?.length ?? 0;
    final desc = c.description.trim();
    return desc.isEmpty ? '$voiceCount' : '$voiceCount · $desc';
  }

  Widget _listPane({
    required ColorScheme cs,
    required bool dual,
  }) {
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
    return EntityDismissibleList(
      itemCount: _items.length,
      style: entitySwipeStyle(cs),
      emptyBuilder: (_) => Padding(
        padding: const EdgeInsets.all(32),
        child: Text(
          '暂无角色',
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.bodyLarge?.copyWith(
            color: cs.onSurfaceVariant,
          ),
        ),
      ),
      onDismiss: (index, direction) async {
        if (direction == DismissDirection.startToEnd) {
          return _cloneAt(index);
        }
        if (direction == DismissDirection.endToStart) {
          return _deleteAt(index);
        }
        return false;
      },
      onTap: (index) =>
          _selectCharacter(_items[index].name, dual: dual),
      menuChildren: (context, menuIndex, close) {
        final name =
            (menuIndex != null &&
                menuIndex >= 0 &&
                menuIndex < _items.length)
            ? _items[menuIndex].name
            : null;
        return buildEntityContextMenu(
          groupLabel: '角色',
          name: name,
          close: close,
          onEdit: () {
            if (menuIndex != null) {
              _selectCharacter(_items[menuIndex].name, dual: dual);
            }
          },
          onClone: () {
            if (menuIndex != null) unawaited(_cloneAt(menuIndex));
          },
          onDelete: () {
            if (menuIndex != null) unawaited(_deleteAt(menuIndex));
          },
        );
      },
      itemBuilder: (context, index) {
        final c = _items[index];
        return EntityListRow(
          leading: const EntityListAvatar(Icons.person_rounded),
          title: c.name,
          subtitle: EntityListMeta(_characterMeta(c)),
          selected: dual && _selectedCharacter == c.name,
        );
      },
    );
  }

  Widget _emptyDetail(String message, ColorScheme cs) {
    return Center(
      key: ValueKey('empty-$message'),
      child: Text(
        message,
        style: Theme.of(context).textTheme.bodyLarge?.copyWith(
          color: cs.onSurfaceVariant,
        ),
      ),
    );
  }

  Widget _detailPane({required bool triple}) {
    final cs = Theme.of(context).colorScheme;
    final name = _selectedCharacter;
    if (name == null) {
      return _emptyDetail('选择一个角色开始编辑', cs);
    }
    // Dual: open voice in the right column (covers character edit).
    if (!triple && _selectedVoice != null) {
      return VoiceEditPage(
        key: ValueKey('voice-$name-$_selectedVoice'),
        characterName: name,
        voiceName: _selectedVoice!,
        embedded: true,
        onBack: () => setState(() => _selectedVoice = null),
        onVoiceRenamed: (next) => setState(() => _selectedVoice = next),
        onUpdated: () => unawaited(_reload()),
      );
    }
    return CharacterEditPage(
      key: ValueKey('char-$name'),
      characterName: name,
      embedded: true,
      selectedVoiceName: _selectedVoice,
      onOpenVoice: (voice) => setState(() => _selectedVoice = voice),
      onDeselectVoice: () => setState(() => _selectedVoice = null),
      onRenamed: (next) {
        setState(() => _selectedCharacter = next);
        unawaited(_reload());
      },
      onUpdated: () => unawaited(_reload()),
    );
  }

  Widget _voicePane() {
    final cs = Theme.of(context).colorScheme;
    final character = _selectedCharacter;
    final voice = _selectedVoice;
    if (character == null || voice == null) {
      return _emptyDetail('选择一个音色进行调参', cs);
    }
    return VoiceEditPage(
      key: ValueKey('voice-$character-$voice'),
      characterName: character,
      voiceName: voice,
      embedded: true,
      onVoiceRenamed: (next) => setState(() => _selectedVoice = next),
      onUpdated: () => unawaited(_reload()),
    );
  }

  Widget _listTitleActions() {
    return M3EToggleButtonGroup(
      type: M3EButtonGroupType.connected,
      style: M3EButtonStyle.tonal,
      size: M3EButtonSize.xs,
      shape: M3EButtonShape.round,
      density: M3EButtonGroupDensity.compact,
      overflow: M3EButtonGroupOverflow.none,
      selectedIndex: null,
      onSelectedIndexChanged: (idx) {
        if (idx == 0) unawaited(_reload());
        if (idx == 1) unawaited(_showAddDialog(dual: true));
      },
      actions: const [
        M3EToggleButtonGroupAction(
          icon: Icon(Icons.refresh_rounded, size: 18),
          tooltip: '刷新',
        ),
        M3EToggleButtonGroupAction(
          icon: Icon(Icons.person_add_alt_1_rounded, size: 18),
          tooltip: '添加角色',
        ),
      ],
    );
  }

  Widget _wideLayout({
    required ColorScheme cs,
    required bool triple,
    required double bottom,
  }) {
    final voiceCoversDetail = !triple && _selectedVoice != null;
    return MultiColumnView(
      padding: EdgeInsets.fromLTRB(20, 12, 20, bottom + 16),
      scrollDetail: !voiceCoversDetail,
      scrollSecondary: false,
      master: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          PageCardTitle('角色', trailing: _listTitleActions()),
          Expanded(
            child: SingleChildScrollView(
              child: _listPane(cs: cs, dual: true),
            ),
          ),
        ],
      ),
      detail: _detailPane(triple: triple),
      secondary: triple ? _voicePane() : null,
    );
  }

  Widget _narrowLayout({
    required ColorScheme cs,
    required double top,
    required double bottom,
  }) {
    return CustomScrollView(
      slivers: [
        SliverPadding(
          padding: EdgeInsets.fromLTRB(20, top + 12, 20, bottom + 88),
          sliver: SliverToBoxAdapter(
            child: MultiColumnPane(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const PageCardTitle('角色'),
                  _listPane(cs: cs, dual: false),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final top = MediaQuery.paddingOf(context).top;
    final mode = multiColumnModeOf(context);
    final dual = mode != MultiColumnMode.narrow;
    final triple = mode == MultiColumnMode.triple;

    return withBarBottomPadding(context, (bottom) {
      if (dual) {
        return _wideLayout(cs: cs, triple: triple, bottom: bottom);
      }
      return Stack(
        children: [
          Positioned.fill(
            child: _narrowLayout(cs: cs, top: top, bottom: bottom),
          ),
          Positioned(
            right: 20,
            bottom: bottom + 8,
            child: M3EFabMenu(
              position: M3EFabMenuPosition.right,
              color: M3EFabColor.primary,
              size: M3EFabSize.medium,
              decoration: const M3EFabMenuDecoration(),
              items: [
                M3EFabMenuItem(
                  icon: const Icon(Icons.refresh_rounded),
                  label: '刷新',
                  onPressed: () => unawaited(_reload()),
                ),
                M3EFabMenuItem(
                  icon: const Icon(Icons.person_add_alt_1_rounded),
                  label: '添加角色',
                  onPressed: () => unawaited(_showAddDialog(dual: false)),
                ),
              ],
            ),
          ),
        ],
      );
    });
  }
}
