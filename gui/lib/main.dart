import 'package:material_ui/material_ui.dart';

import 'app_scope.dart';
import 'bridge/client.dart';
import 'bridge/stdio_transport.dart';
import 'features/shell/main_shell.dart';
import 'theme.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final client = YukumoClient(StdioRpcTransport());
  Object? connectError;
  try {
    await client.connect();
  } catch (e) {
    connectError = e;
  }
  runApp(YukumoApp(client: client, connectError: connectError));
}

class YukumoApp extends StatelessWidget {
  const YukumoApp({
    super.key,
    required this.client,
    this.connectError,
  });

  final YukumoClient client;
  final Object? connectError;

  @override
  Widget build(BuildContext context) {
    return AppScope(
      client: client,
      child: MaterialApp(
        title: 'Yukumo',
        theme: buildYukumoTheme(),
        darkTheme: buildYukumoTheme(brightness: Brightness.dark),
        home: connectError == null
            ? const MainShell()
            : _ConnectErrorPage(error: connectError!, client: client),
      ),
    );
  }
}

class _ConnectErrorPage extends StatelessWidget {
  const _ConnectErrorPage({required this.error, required this.client});

  final Object error;
  final YukumoClient client;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('无法连接 yukumo-script serve', style: Theme.of(context).textTheme.titleLarge),
              const SizedBox(height: 12),
              Text('$error', textAlign: TextAlign.center),
              const SizedBox(height: 8),
              const Text('设置 YUKUMO_RPC_BIN 指向二进制后重试'),
              const SizedBox(height: 20),
              FilledButton(
                onPressed: () async {
                  try {
                    await client.connect();
                    if (context.mounted) {
                      Navigator.of(context).pushReplacement(
                        MaterialPageRoute(builder: (_) => const MainShell()),
                      );
                    }
                  } catch (e) {
                    if (context.mounted) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('$e')),
                      );
                    }
                  }
                },
                child: const Text('重试'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
