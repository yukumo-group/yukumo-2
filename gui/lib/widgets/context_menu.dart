import 'package:material_ui/material_ui.dart';

/// [M3EMenuController.open] expects a local offset relative to the MenuAnchor.
Offset menuLocalBelowLeft(BuildContext anchorContext, BuildContext itemContext) {
  final anchorBox = anchorContext.findRenderObject() as RenderBox?;
  final itemBox = itemContext.findRenderObject() as RenderBox?;
  if (anchorBox == null || itemBox == null || !anchorBox.hasSize || !itemBox.hasSize) {
    return Offset.zero;
  }
  final belowLeftGlobal = itemBox.localToGlobal(Offset(0, itemBox.size.height));
  return anchorBox.globalToLocal(belowLeftGlobal);
}
