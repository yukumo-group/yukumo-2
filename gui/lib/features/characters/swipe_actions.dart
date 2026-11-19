import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

M3EDismissibleCardStyle characterSwipeStyle(ColorScheme cs) {
  return M3EDismissibleCardStyle(
    direction: DismissDirection.horizontal,
    outerRadius: 24,
    innerRadius: 10,
    gap: 6,
    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
    autoExecutePrimaryOnFullSwipe: true,
    actions: [
      M3ESwipeAction(
        isPrimary: true,
        dismissOnTap: true,
        backgroundColor: cs.secondaryContainer,
        foregroundColor: cs.onSecondaryContainer,
        icon: const FaIcon(FontAwesomeIcons.copy, size: 18),
        onTap: () {},
      ),
    ],
    secondaryActions: [
      M3ESwipeAction(
        isPrimary: true,
        dismissOnTap: true,
        backgroundColor: cs.error,
        foregroundColor: cs.onError,
        icon: const FaIcon(FontAwesomeIcons.trashCan, size: 18),
        onTap: () {},
      ),
    ],
  );
}
