import 'package:material_ui/material_ui.dart';

const double kPageBreakpoint = 720;
const double kContentMaxWidth = 980;

class ContentFrame extends StatelessWidget {
  const ContentFrame({
    super.key,
    required this.child,
    this.maxWidth = kContentMaxWidth,
    this.padding = const EdgeInsets.symmetric(horizontal: 20),
    this.cardOnWide = true,
  });

  final Widget child;
  final double maxWidth;
  final EdgeInsetsGeometry padding;
  final bool cardOnWide;

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final wide = width > maxWidth + 48;
    final body = ConstrainedBox(
      constraints: BoxConstraints(maxWidth: maxWidth),
      child: child,
    );
    final framed = wide && cardOnWide
        ? Card(
            elevation: 0,
            color: Theme.of(context).colorScheme.surfaceContainerLow,
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(28)),
            clipBehavior: Clip.antiAlias,
            child: Padding(
              padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
              child: body,
            ),
          )
        : body;
    return Align(
      alignment: Alignment.topCenter,
      child: Padding(padding: padding, child: framed),
    );
  }
}

bool isWideLayout(BuildContext context, {double breakpoint = kPageBreakpoint}) {
  return MediaQuery.sizeOf(context).width >= breakpoint;
}

/// Space between [PageCardTitle] and the card body beneath it.
const double kPageCardTitleGap = 24;

/// Left-aligned page title for root tab content (replaces large AppBar title).
class PageCardTitle extends StatelessWidget {
  const PageCardTitle(
    this.title, {
    super.key,
    this.trailing,
    this.bottomGap = kPageCardTitleGap,
  });

  final String title;
  final Widget? trailing;

  /// Gap below the title row before card content. Use `0` when spacing is
  /// applied by the parent instead.
  final double bottomGap;

  @override
  Widget build(BuildContext context) {
    final style = Theme.of(context).textTheme.headlineMedium?.copyWith(
      fontWeight: FontWeight.w600,
    );
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Expanded(child: Text(title, style: style)),
            if (trailing != null) trailing!,
          ],
        ),
        if (bottomGap > 0) SizedBox(height: bottomGap),
      ],
    );
  }
}
