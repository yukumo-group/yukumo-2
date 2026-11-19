import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:stream_channel/stream_channel.dart';

import 'transport.dart';

class StdioRpcTransport implements RpcTransport {
  StdioRpcTransport({String? binary, List<String>? args})
      : binary = binary ?? resolveYukumoBin(),
        args = args ?? const ['serve'];

  final String binary;
  final List<String> args;

  Process? _process;
  final _stdinController = StreamController<String>();

  @override
  Future<StreamChannel<String>> connect() async {
    final process = await Process.start(binary, args);
    _process = process;

    unawaited(
      _stdinController.stream.listen((line) {
        process.stdin.writeln(line);
      }).asFuture<void>().catchError((_) {}),
    );

    process.stderr
        .transform(utf8.decoder)
        .transform(const LineSplitter())
        .listen((line) {
      stderr.writeln('[yukumo-serve] $line');
    });

    final stdoutLines = process.stdout
        .transform(utf8.decoder)
        .transform(const LineSplitter());

    return StreamChannel<String>.withGuarantees(
      stdoutLines,
      _stdinController.sink,
    );
  }

  @override
  Future<void> close() async {
    await _stdinController.close();
    final p = _process;
    if (p == null) return;
    p.stdin.close();
    p.kill();
    await p.exitCode;
  }
}

String resolveYukumoBin() {
  final fromEnv = Platform.environment['YUKUMO_RPC_BIN'];
  if (fromEnv != null && fromEnv.isNotEmpty) {
    return fromEnv;
  }
  final candidates = <String>[
    '${Directory.current.path}/yukumo-script',
    '${Directory.current.path}/../yukumo-script',
    'yukumo-script',
  ];
  for (final path in candidates) {
    if (File(path).existsSync()) {
      return path;
    }
  }
  return 'yukumo-script';
}
