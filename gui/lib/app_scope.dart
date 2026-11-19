import 'package:material_ui/material_ui.dart';

import 'bridge/client.dart';

class AppScope extends InheritedWidget {
  const AppScope({
    super.key,
    required this.client,
    required super.child,
  });

  final YukumoClient client;

  static YukumoClient of(BuildContext context) {
    final scope = context.dependOnInheritedWidgetOfExactType<AppScope>();
    assert(scope != null, 'AppScope not found');
    return scope!.client;
  }

  @override
  bool updateShouldNotify(AppScope oldWidget) => client != oldWidget.client;
}
