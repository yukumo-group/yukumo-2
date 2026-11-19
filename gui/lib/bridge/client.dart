import 'dart:async';

import 'package:json_rpc_2/json_rpc_2.dart';

import '../gen/rpc.dart';
import 'transport.dart';

class YukumoClient {
  YukumoClient(this._transport);

  final RpcTransport _transport;
  Client? _client;

  late final CharacterApi characters;
  late final VoiceApi voices;
  late final ConfigApi config;
  late final SynthApi synth;

  Future<void> connect() async {
    final channel = await _transport.connect();
    final client = Client(channel);
    _client = client;
    characters = CharacterApi(client);
    voices = VoiceApi(client);
    config = ConfigApi(client);
    synth = SynthApi(client);
    unawaited(client.listen());
  }

  Future<void> close() async {
    await _client?.close();
    await _transport.close();
  }
}
