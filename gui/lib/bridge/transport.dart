import 'package:stream_channel/stream_channel.dart';

abstract class RpcTransport {
  Future<StreamChannel<String>> connect();
  Future<void> close();
}
