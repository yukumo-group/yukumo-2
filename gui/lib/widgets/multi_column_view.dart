import 'package:material_ui/material_ui.dart';

/// Master + detail (≥ [dualBreakpoint]); optional third pane (≥ [tripleBreakpoint]).
const double kMultiColumnDualBreakpoint = 900;
const double kMultiColumnTripleBreakpoint = 1280;

enum MultiColumnMode { narrow, dual, triple }

/// Resolves narrow / dual / triple from [MediaQuery] width.
MultiColumnMode multiColumnModeOf(
  BuildContext context, {
  double dualBreakpoint = kMultiColumnDualBreakpoint,
  double tripleBreakpoint = kMultiColumnTripleBreakpoint,
}) {
  final width = MediaQuery.sizeOf(context).width;
  if (width >= tripleBreakpoint) return MultiColumnMode.triple;
  if (width >= dualBreakpoint) return MultiColumnMode.dual;
  return MultiColumnMode.narrow;
}

/// Surface card used by each column in [MultiColumnView].
class MultiColumnPane extends StatelessWidget {
  const MultiColumnPane({
    super.key,
    required this.child,
    this.padding = const EdgeInsets.fromLTRB(16, 16, 16, 20),
  });

  final Widget child;
  final EdgeInsetsGeometry padding;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Card(
      elevation: 0,
      color: cs.surfaceContainerHigh,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      clipBehavior: Clip.antiAlias,
      margin: EdgeInsets.zero,
      child: Padding(padding: padding, child: child),
    );
  }
}

/// Stretching row of pane cards for dual / triple master–detail layouts.
///
/// Call only when [multiColumnModeOf] is not [MultiColumnMode.narrow].
/// [secondary] appears only at triple width when non-null.
class MultiColumnView extends StatelessWidget {
  const MultiColumnView({
    super.key,
    required this.master,
    required this.detail,
    this.secondary,
    this.padding = const EdgeInsets.fromLTRB(20, 12, 20, 16),
    this.gap = 12,
    this.masterWidth = 340,
    this.masterWidthTriple = 360,
    this.detailFlex = 7,
    this.detailFlexTriple = 6,
    this.secondaryFlex = 6,
    this.dualBreakpoint = kMultiColumnDualBreakpoint,
    this.tripleBreakpoint = kMultiColumnTripleBreakpoint,
    this.scrollDetail = true,
    this.scrollSecondary = true,
  });

  final Widget master;
  final Widget detail;
  final Widget? secondary;

  final EdgeInsetsGeometry padding;
  final double gap;

  final double masterWidth;
  final double masterWidthTriple;
  final int detailFlex;
  final int detailFlexTriple;
  final int secondaryFlex;

  final double dualBreakpoint;
  final double tripleBreakpoint;

  /// Wrap [detail] / [secondary] in [SingleChildScrollView].
  final bool scrollDetail;
  final bool scrollSecondary;

  @override
  Widget build(BuildContext context) {
    final mode = multiColumnModeOf(
      context,
      dualBreakpoint: dualBreakpoint,
      tripleBreakpoint: tripleBreakpoint,
    );
    final triple = mode == MultiColumnMode.triple;
    final showSecondary = secondary != null && triple;

    Widget wrapScroll(Widget child, {required bool enabled}) {
      if (!enabled) return child;
      return SingleChildScrollView(child: child);
    }

    return SafeArea(
      bottom: false,
      child: Padding(
        padding: padding,
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            SizedBox(
              width: triple ? masterWidthTriple : masterWidth,
              child: MultiColumnPane(child: master),
            ),
            SizedBox(width: gap),
            Expanded(
              flex: showSecondary ? detailFlexTriple : detailFlex,
              child: MultiColumnPane(
                child: wrapScroll(detail, enabled: scrollDetail),
              ),
            ),
            if (showSecondary) ...[
              SizedBox(width: gap),
              Expanded(
                flex: secondaryFlex,
                child: MultiColumnPane(
                  child: wrapScroll(secondary!, enabled: scrollSecondary),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
