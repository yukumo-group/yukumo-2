import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import 'hold_wipe_confirm.dart';

/// [M3ESwipeAction] that confirms via hold + linear wipe instead of a tap.
class HoldWipeSwipeAction extends M3ESwipeAction {
  const HoldWipeSwipeAction({
    required super.icon,
    required this.wipeColor,
    super.label,
    super.backgroundColor,
    super.foregroundColor,
    super.onTap,
    super.isPrimary = true,
    super.dismissOnTap = true,
    super.width = 56,
    super.height,
    super.borderRadius,
    super.haptic = M3EHapticFeedback.medium,
    this.holdDuration = const Duration(milliseconds: 750),
  });

  final Color wipeColor;
  final Duration holdDuration;

  @override
  Widget buildButton(
    BuildContext context, {
    required VoidCallback? onTriggered,
    VoidCallback? onDismissTriggered,
    FocusNode? focusNode,
    KeyEventResult Function(FocusNode node, KeyEvent event)? onKeyEvent,
  }) {
    final theme = Theme.of(context);
    final cs = theme.colorScheme;
    final effectiveBg =
        backgroundColor ?? (isPrimary ? cs.primary : cs.secondaryContainer);
    final effectiveFg =
        foregroundColor ?? (isPrimary ? cs.onPrimary : cs.onSecondaryContainer);
    final effectiveRadius =
        (borderRadius as BorderRadius?) ?? BorderRadius.circular(100);

    if (focusNode != null && onKeyEvent != null) {
      focusNode.onKeyEvent = onKeyEvent;
    }

    void handleConfirm() {
      haptic.apply();
      if (dismissOnTap && onDismissTriggered != null) {
        onDismissTriggered();
        onTap?.call();
      } else {
        onTap?.call();
        onTriggered?.call();
      }
    }

    final content = Center(
      child: IconTheme(
        data: IconThemeData(color: effectiveFg, size: 24),
        child: FittedBox(
          fit: BoxFit.scaleDown,
          child: label != null
              ? Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    icon,
                    const SizedBox(height: 2),
                    DefaultTextStyle(
                      style: (theme.textTheme.labelSmall ?? const TextStyle())
                          .copyWith(
                            color: effectiveFg,
                            fontSize: 10,
                            fontWeight: FontWeight.w600,
                          ),
                      child: label!,
                    ),
                  ],
                )
              : icon,
        ),
      ),
    );

    final button = HoldWipeConfirm(
      width: width,
      height: height,
      backgroundColor: effectiveBg,
      wipeColor: wipeColor,
      borderRadius: effectiveRadius,
      holdDuration: holdDuration,
      onConfirmed: handleConfirm,
      child: content,
    );

    if (focusNode == null) return button;

    return ListenableBuilder(
      listenable: focusNode,
      builder: (context, child) {
        return Focus(
          focusNode: focusNode,
          child: child!,
        );
      },
      child: button,
    );
  }
}
