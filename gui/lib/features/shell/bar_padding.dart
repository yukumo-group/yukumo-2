import 'package:flutter_floating_bottom_bar/flutter_floating_bottom_bar.dart';
import 'package:material_ui/material_ui.dart';

Widget withBarBottomPadding(
  BuildContext context,
  Widget Function(double bottomPadding) builder,
) {
  final barHeight = BottomBarScope.maybeOf(context)?.barHeight;
  if (barHeight == null) {
    return builder(MediaQuery.paddingOf(context).bottom + 24);
  }
  return ValueListenableBuilder<double>(
    valueListenable: barHeight,
    builder: (context, value, _) => builder(value + 24),
  );
}
