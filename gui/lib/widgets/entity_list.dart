import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

import 'context_menu.dart';

/// Must stay in sync with dismissible card padding in swipe styles.
const kEntitySwipePadding = EdgeInsets.symmetric(horizontal: 14, vertical: 10);

/// Fixed-size leading glyph shared by character / voice rows.
class EntityListAvatar extends StatelessWidget {
  const EntityListAvatar(
    this.icon, {
    super.key,
  });

  final IconData icon;

  static const double size = 40;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return SizedBox(
      width: size,
      height: size,
      child: Center(
        child: Icon(icon, size: 32, color: cs.primary),
      ),
    );
  }
}

/// One-line meta under an [EntityListRow] title (shared character / voice height).
class EntityListMeta extends StatelessWidget {
  const EntityListMeta(this.text, {super.key});

  final String text;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Text(
      text,
      maxLines: 1,
      overflow: TextOverflow.ellipsis,
      style: Theme.of(context).textTheme.bodySmall?.copyWith(
        color: cs.onSurfaceVariant,
      ),
    );
  }
}

/// Shared title / subtitle / chevron row used by character & voice lists.
class EntityListRow extends StatelessWidget {
  const EntityListRow({
    super.key,
    required this.title,
    this.subtitle,
    this.leading,
    this.showChevron = true,
    this.selected = false,
    this.selectedColor,
  });

  final String title;
  final Widget? subtitle;
  final Widget? leading;
  final bool showChevron;
  final bool selected;
  final Color? selectedColor;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final row = Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        if (leading != null) ...[
          leading!,
          const SizedBox(width: 12),
        ],
        Expanded(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                title,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 2),
              SizedBox(
                height: 18,
                child: subtitle,
              ),
              const SizedBox(height: 2),
            ],
          ),
        ),
        if (showChevron)
          Icon(Icons.chevron_right_rounded, color: cs.onSurfaceVariant),
      ],
    );
    if (!selected) return row;
    // Bleed into card padding so only the card background color changes.
    return Stack(
      clipBehavior: Clip.none,
      children: [
        Positioned(
          left: -kEntitySwipePadding.left,
          right: -kEntitySwipePadding.right,
          top: -kEntitySwipePadding.top,
          bottom: -kEntitySwipePadding.bottom,
          child: ColoredBox(
            color: selectedColor ?? cs.secondaryContainer,
          ),
        ),
        row,
      ],
    );
  }
}

typedef EntityMenuChildrenBuilder = List<Widget> Function(
  BuildContext context,
  int? menuIndex,
  VoidCallback closeMenu,
);

/// Dismissible card list with a shared context menu (secondary click + long press).
class EntityDismissibleList extends StatefulWidget {
  const EntityDismissibleList({
    super.key,
    required this.itemCount,
    required this.itemBuilder,
    required this.menuChildren,
    required this.style,
    this.onDismiss,
    this.onTap,
    this.emptyBuilder,
  });

  final int itemCount;
  final IndexedWidgetBuilder itemBuilder;
  final EntityMenuChildrenBuilder menuChildren;
  final M3EDismissibleCardStyle style;
  final Future<bool> Function(int index, DismissDirection direction)? onDismiss;
  final void Function(int index)? onTap;
  final WidgetBuilder? emptyBuilder;

  @override
  State<EntityDismissibleList> createState() => _EntityDismissibleListState();
}

class _EntityDismissibleListState extends State<EntityDismissibleList> {
  final _menuController = M3EMenuController();
  BuildContext? _anchorContext;
  int? _menuIndex;

  @override
  void dispose() {
    _menuController.dispose();
    super.dispose();
  }

  void _closeMenu() => _menuController.close();

  void _openMenu(BuildContext itemContext, int index) {
    final anchor = _anchorContext;
    if (anchor == null) return;
    setState(() => _menuIndex = index);
    _menuController.open(
      position: menuLocalBelowLeft(anchor, itemContext),
    );
  }

  @override
  Widget build(BuildContext context) {
    return M3EMenuAnchor(
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
        position: M3EMenuPosition.under,
      ),
      menuChildren: widget.menuChildren(context, _menuIndex, _closeMenu),
      builder: (context, _, child) {
        _anchorContext = context;
        return child!;
      },
      child: M3EDismissibleCardColumn(
        itemCount: widget.itemCount,
        emptyBuilder: widget.emptyBuilder,
        style: widget.style,
        onDismiss: widget.onDismiss,
        onTap: widget.onTap,
        itemBuilder: (context, index) {
          return GestureDetector(
            behavior: HitTestBehavior.opaque,
            onSecondaryTapDown: (_) => _openMenu(context, index),
            onLongPress: () => _openMenu(context, index),
            child: widget.itemBuilder(context, index),
          );
        },
      ),
    );
  }
}
