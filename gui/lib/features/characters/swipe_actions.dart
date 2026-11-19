import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import '../../widgets/entity_list.dart';
import '../../widgets/hold_wipe_swipe_action.dart';

const _kSwipeActionWidth = 88.0;

M3EDismissibleCardStyle entitySwipeStyle(
  ColorScheme cs, {
  Color? color,
}) {
  return M3EDismissibleCardStyle(
    direction: DismissDirection.horizontal,
    outerRadius: 24,
    innerRadius: 10,
    gap: 6,
    color: color,
    padding: kEntitySwipePadding,
    autoExecutePrimaryOnFullSwipe: true,
    actions: [
      M3ESwipeAction(
        isPrimary: true,
        dismissOnTap: true,
        width: _kSwipeActionWidth,
        backgroundColor: cs.secondaryContainer,
        foregroundColor: cs.onSecondaryContainer,
        icon: const FaIcon(FontAwesomeIcons.copy, size: 18),
        label: const Text('克隆'),
        onTap: () {},
      ),
    ],
    secondaryActions: [
      HoldWipeSwipeAction(
        width: _kSwipeActionWidth,
        wipeColor: cs.error,
        backgroundColor: cs.errorContainer,
        foregroundColor: cs.onErrorContainer,
        icon: const FaIcon(FontAwesomeIcons.trashCan, size: 18),
        label: const Text('长按'),
        onTap: () {},
      ),
    ],
  );
}
