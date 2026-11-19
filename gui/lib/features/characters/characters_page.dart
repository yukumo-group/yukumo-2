import 'dart:async';

import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../app_scope.dart';
import '../../gen/rpc.dart';
import '../shell/bar_padding.dart';

class CharactersPage extends StatefulWidget {
  const CharactersPage({super.key});

  @override
  State<CharactersPage> createState() => _CharactersPageState();
}

class _CharactersPageState extends State<CharactersPage> {
  List<Character> _items = const [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      unawaited(_reload());
    });
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
      await _reload();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('创建失败: $e')));
    }
  }

  void _openContextMenu(Offset globalPosition, int index) {
    final c = _items[index];
    final overlay = Overlay.of(context).context.findRenderObject() as RenderBox;
    showMenu<String>(
      context: context,
      position: RelativeRect.fromRect(
        globalPosition & const Size(1, 1),
        Offset.zero & overlay.size,
      ),
      items: [
        PopupMenuItem(value: 'delete', child: Text('删除「${c.name}」')),
        const PopupMenuItem(value: 'refresh', child: Text('刷新列表')),
      ],
    ).then((value) async {
      if (value == 'delete') {
        await _deleteAt(index);
      } else if (value == 'refresh') {
        await _reload();
      }
    });
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
              padding: EdgeInsets.fromLTRB(16, 0, 16, bottom),
              sliver: SliverToBoxAdapter(
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
                  style: M3EDismissibleCardStyle(
                    direction: DismissDirection.startToEnd,
                    outerRadius: 24,
                    innerRadius: 10,
                    gap: 6,
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                    actions: [
                      M3ESwipeAction(
                        isPrimary: true,
                        dismissOnTap: true,
                        backgroundColor: cs.error,
                        foregroundColor: cs.onError,
                        icon: const FaIcon(FontAwesomeIcons.trashCan, size: 18),
                        label: const Text('删除'),
                        onTap: () {},
                      ),
                    ],
                    autoExecutePrimaryOnFullSwipe: true,
                  ),
                  onDismiss: (index, direction) => _deleteAt(index),
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
        ],
      );
    });
  }
}
