import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import 'hold_wipe_menu_item.dart';

/// Standard edit / clone / hold-delete menu used by entity lists.
List<Widget> buildEntityContextMenu({
  required String groupLabel,
  required String? name,
  required VoidCallback close,
  required VoidCallback onEdit,
  required VoidCallback onClone,
  required VoidCallback onDelete,
  List<Widget> extraItems = const [],
}) {
  return [
    M3EMenuGroup(
      shapes: MenuGroupShapes.leading(),
      children: [
        M3EMenuGroupLabel(child: Text(groupLabel)),
        M3EMenuItemButton(
          onPressed: () {
            close();
            onEdit();
          },
          leadingIcon: const Icon(Icons.edit_rounded),
          child: Text(
            (name == null || name.isEmpty) ? '编辑' : '编辑「$name」',
          ),
        ),
        M3EMenuItemButton(
          onPressed: () {
            close();
            onClone();
          },
          leadingIcon: const Icon(Icons.copy_rounded),
          child: const Text('克隆'),
        ),
        ...extraItems,
      ],
    ),
    M3EMenuGroup(
      shapes: MenuGroupShapes.trailing(),
      children: [
        const M3EMenuGroupLabel(child: Text('危险操作')),
        HoldWipeMenuItem(
          label: const Text('长按删除'),
          onConfirmed: () {
            close();
            onDelete();
          },
        ),
      ],
    ),
  ];
}
