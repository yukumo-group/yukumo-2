import 'package:material_ui/material_ui.dart';

import 'hold_wipe_confirm.dart';

/// Danger menu row: hold to confirm with a linear wipe fill.
class HoldWipeMenuItem extends StatelessWidget {
  const HoldWipeMenuItem({
    super.key,
    required this.label,
    required this.onConfirmed,
    this.leadingIcon,
    this.holdDuration = const Duration(milliseconds: 750),
  });

  final Widget label;
  final Widget? leadingIcon;
  final VoidCallback onConfirmed;
  final Duration holdDuration;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
      child: HoldWipeConfirm(
        backgroundColor: Colors.transparent,
        wipeColor: cs.error.withValues(alpha: 0.28),
        borderRadius: BorderRadius.circular(16),
        holdDuration: holdDuration,
        onConfirmed: onConfirmed,
        child: ConstrainedBox(
          constraints: const BoxConstraints(minHeight: 48),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
            child: IconTheme(
              data: IconThemeData(size: 24, color: cs.error),
              child: DefaultTextStyle(
                style: (theme.textTheme.labelLarge ?? const TextStyle())
                    .copyWith(color: cs.error),
                child: Row(
                  children: [
                    leadingIcon ??
                        Icon(Icons.delete_forever_rounded, color: cs.error),
                    const SizedBox(width: 12),
                    Expanded(child: label),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
