import 'package:stream_channel/stream_channel.dart';

import 'transport.dart';

class WasmRpcTransport implements RpcTransport {
  @override
  Future<StreamChannel<String>> connect() {
    throw UnimplementedError(
      'WasmRpcTransport: wasip1 core.wasm + host_synth is deferred; see docs/refactorings/4-flutter-gui.md §9',
    );
  }

  @override
  Future<void> close() async {}
}
