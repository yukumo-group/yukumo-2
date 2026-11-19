import 'dart:io';

import 'package:yukumo_gui_next/bridge/client.dart';
import 'package:yukumo_gui_next/bridge/stdio_transport.dart';
import 'package:yukumo_gui_next/gen/rpc.dart';

Future<void> main() async {
  final client = YukumoClient(StdioRpcTransport());
  await client.connect();
  try {
    final list = await client.characters.list();
    final names = (list.characters ?? const <Character>[])
        .map((c) => c.name)
        .join(', ');
    stdout.writeln('characters: ${names.isEmpty ? "(none)" : names}');

    final cfg = await client.config.get();
    stdout.writeln('outputFormat: ${cfg.synth.outputFormat}');
  } finally {
    await client.close();
  }
}
