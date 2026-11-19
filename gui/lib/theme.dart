import 'package:material_ui/material_ui.dart';

ThemeData buildYukumoTheme({Brightness brightness = Brightness.light}) {
  const seed = Color(0xFF5B6CFF);
  final scheme = ColorScheme.fromSeed(seedColor: seed, brightness: brightness);
  return ThemeData(
    colorScheme: scheme,
    useMaterial3: true,
    popupMenuTheme: PopupMenuThemeData(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
    ),
  );
}
