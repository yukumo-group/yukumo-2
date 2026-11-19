import 'dart:async';

import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../../widgets/content_frame.dart';
import '../shell/bar_padding.dart';
import 'character_edit_page.dart';
import 'name_util.dart';
import 'swipe_actions.dart';

class CharactersPage extends StatefulWidget {
  const CharactersPage({super.key});

  @override
  State<CharactersPage> createState() => _CharactersPageState();
}

class _CharactersPageState extends State<CharactersPage> {
  final _menuController = M3EMenuController();
  List<Character> _items = const [];
  bool _loading = true;
  String? _error;
  int? _menuIndex;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_reload());
    });
  }

  @override
  void dispose() {
    _menuController.dispose();
    super.dispose();
  }

  Future<void> _reload() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final res = await AppScope.of(context).characters.list();
      if (!mounted) return;
      setState(() {
        _items = List.of(res.characters ?? const <Character>[]);
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

  Future<void> _openEdit(String name) async {
    await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => CharacterEditPage(characterName: name),
      ),
    );
    if (mounted) await _reload();
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
    final c = _items[index];
    final client = AppScope.of(context);
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('删除角色'),
        content: Text('确定删除「${c.name}」？'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('删除')),
        ],
      ),
    );
    if (ok != true) return false;
    try {
      await client.characters.delete(CharacterDeleteReq(name: c.name));
      if (!mounted) return false;
      setState(() => _items = List.of(_items)..removeAt(index));
      return true;
    } catch (e) {
      if (!mounted) return false;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('删除失败: $e')));
      return false;
    }
  }

  Future<void> _showAddDialog() async {
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
      await _openEdit(name);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('创建失败: $e')));
    }
  }

  void _openContextMenu(Offset globalPosition, int index) {
    setState(() => _menuIndex = index);
    _menuController.open(position: globalPosition);
  }

  List<Widget> _menuChildren(ColorScheme cs) {
    final index = _menuIndex;
    final name = (index != null && index >= 0 && index < _items.length)
        ? _items[index].name
        : '';
    return [
      M3EMenuGroup(
        shapes: MenuGroupShapes.leading(),
        children: [
          const M3EMenuGroupLabel(child: Text('角色')),
          M3EMenuItemButton(
            onPressed: () {
              _menuController.close();
              if (index != null) unawaited(_openEdit(_items[index].name));
            },
            leadingIcon: const Icon(Icons.edit_rounded),
            child: Text(name.isEmpty ? '编辑' : '编辑「$name」'),
          ),
          M3EMenuItemButton(
            onPressed: () {
              _menuController.close();
              if (index != null) unawaited(_cloneAt(index));
            },
            leadingIcon: const Icon(Icons.copy_rounded),
            child: const Text('克隆'),
          ),
          M3EMenuItemButton(
            onPressed: () {
              _menuController.close();
              unawaited(_reload());
            },
            leadingIcon: const Icon(Icons.refresh_rounded),
            child: const Text('刷新列表'),
          ),
        ],
      ),
      M3EMenuGroup(
        shapes: MenuGroupShapes.trailing(),
        children: [
          const M3EMenuGroupLabel(child: Text('危险操作')),
          M3EMenuItemButton(
            onPressed: () {
              _menuController.close();
              if (index != null) unawaited(_deleteAt(index));
            },
            leadingIcon: Icon(Icons.delete_forever_rounded, color: cs.error),
            itemStyle: M3EMenuItemStyle(foregroundColor: cs.error),
            child: const Text('删除'),
          ),
        ],
      ),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return withBarBottomPadding(context, (bottom) {
      return CustomScrollView(
        slivers: [
          SliverAppBar.large(
            title: const Text('角色'),
            backgroundColor: cs.surface,
            actions: [
              IconButton(
                tooltip: '新建',
                onPressed: _showAddDialog,
                icon: const FaIcon(FontAwesomeIcons.plus, size: 18),
              ),
              IconButton(
                tooltip: '刷新',
                onPressed: _reload,
                icon: const FaIcon(FontAwesomeIcons.arrowsRotate, size: 18),
              ),
            ],
          ),
          if (_loading)
            const SliverFillRemaining(
              child: Center(child: M3ELoadingIndicator()),
            )
          else if (_error != null)
            SliverFillRemaining(
              child: Center(child: Text(_error!)),
            )
          else
            SliverPadding(
              padding: EdgeInsets.only(bottom: bottom),
              sliver: SliverToBoxAdapter(
                child: ContentFrame(
                  child: M3EMenuAnchor(
                    controller: _menuController,
                    colorSchemeVariant: M3EMenuColorSchemeVariant.standard,
                    consumeOutsideTaps: true,
                    anchorTapClosesMenu: true,
                    style: const M3EMenuStyle(
                      minWidth: 180,
                      maxWidth: 300,
                      elevation: 3,
                      spatialMotion: M3EMotion.expressiveSpatialFast,
                      animationStyle: M3EMenuAnimationStyle.popup,
                    ),
                    menuChildren: _menuChildren(cs),
                    builder: (context, controller, child) => child!,
                    child: M3EDismissibleCardColumn(
                      itemCount: _items.length,
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
                      style: characterSwipeStyle(cs),
                      onDismiss: (index, direction) async {
                        if (direction == DismissDirection.startToEnd) {
                          return _cloneAt(index);
                        }
                        if (direction == DismissDirection.endToStart) {
                          return _deleteAt(index);
                        }
                        return false;
                      },
                      onTap: (index) => unawaited(_openEdit(_items[index].name)),
                      itemBuilder: (context, index) {
                        final c = _items[index];
                        final voiceCount = c.voices?.length ?? 0;
                        return GestureDetector(
                          behavior: HitTestBehavior.opaque,
                          onSecondaryTapDown: (details) {
                            _openContextMenu(details.globalPosition, index);
                          },
                          child: Row(
                            children: [
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      c.name,
                                      style: Theme.of(context).textTheme.titleMedium,
                                    ),
                                    if (c.description.isNotEmpty) ...[
                                      const SizedBox(height: 4),
                                      Text(
                                        c.description,
                                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                          color: cs.onSurfaceVariant,
                                        ),
                                      ),
                                    ],
                                    const SizedBox(height: 4),
                                    Text(
                                      '$voiceCount 个音色',
                                      style: Theme.of(context).textTheme.labelMedium?.copyWith(
                                        color: cs.onSurfaceVariant,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              Icon(
                                Icons.chevron_right_rounded,
                                color: cs.onSurfaceVariant,
                              ),
                            ],
                          ),
                        );
                      },
                    ),
                  ),
                ),
              ),
            ),
        ],
      );
    });
  }
}
