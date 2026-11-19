import 'package:flutter/material.dart';

import 'bridge/client.dart';
import 'bridge/stdio_transport.dart';
import 'gen/rpc.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final client = YukumoClient(StdioRpcTransport());
  await client.connect();
  runApp(YukumoApp(client: client));
}

class YukumoApp extends StatefulWidget {
  const YukumoApp({super.key, required this.client});

  final YukumoClient client;

  @override
  State<YukumoApp> createState() => _YukumoAppState();
}

class _YukumoAppState extends State<YukumoApp> {
  late Future<CharacterListRes> _characters;

  @override
  void initState() {
    super.initState();
    _characters = widget.client.characters.list();
  }

  @override
  void dispose() {
    widget.client.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Yukumo',
      home: Scaffold(
        appBar: AppBar(title: const Text('Yukumo RPC smoke')),
        body: FutureBuilder<CharacterListRes>(
          future: _characters,
          builder: (context, snap) {
            if (snap.hasError) {
              return Center(child: Text('error: ${snap.error}'));
            }
            if (!snap.hasData) {
              return const Center(child: CircularProgressIndicator());
            }
            final list = snap.data!.characters ?? const <Character>[];
            if (list.isEmpty) {
              return const Center(child: Text('no characters'));
            }
            return ListView.builder(
              itemCount: list.length,
              itemBuilder: (context, i) {
                final c = list[i];
                return ListTile(
                  title: Text(c.name),
                  subtitle: Text(c.description),
                );
              },
            );
          },
        ),
      ),
    );
  }
}
