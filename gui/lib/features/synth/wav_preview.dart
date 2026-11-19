import 'dart:io';

Future<void> playWavBytes(List<int> bytes) async {
  final dir = await Directory.systemTemp.createTemp('yukumo_preview_');
  final file = File('${dir.path}/preview.wav');
  await file.writeAsBytes(bytes, flush: true);
  try {
    await playWavFile(file.path);
  } finally {
    try {
      await dir.delete(recursive: true);
    } catch (_) {}
  }
}

Future<void> playWavFile(String path) async {
  final candidates = <List<String>>[];
  if (Platform.isMacOS) {
    candidates.add(['afplay', path]);
  } else if (Platform.isWindows) {
    candidates.add([
      'powershell',
      '-NoProfile',
      '-Command',
      '(New-Object Media.SoundPlayer \'$path\').PlaySync()',
    ]);
  } else {
    // Prefer absolute paths: Flutter desktop launches often have a thin PATH.
    candidates.addAll([
      ['/usr/bin/ffplay', '-nodisp', '-autoexit', '-loglevel', 'quiet', path],
      ['/usr/bin/mpv', '--no-video', '--really-quiet', path],
      ['/usr/bin/paplay', path],
      ['/usr/bin/aplay', path],
      ['ffplay', '-nodisp', '-autoexit', '-loglevel', 'quiet', path],
      ['mpv', '--no-video', '--really-quiet', path],
      ['paplay', path],
      ['aplay', path],
      ['play', '-q', path],
    ]);
  }

  Object? last;
  for (final cmd in candidates) {
    try {
      final result = await Process.run(
        cmd.first,
        cmd.sublist(1),
        runInShell: false,
      );
      if (result.exitCode == 0) return;
      last = 'exit ${result.exitCode}: ${result.stderr}'.toString().trim();
    } catch (e) {
      last = e;
    }
  }
  throw StateError('no audio player worked: $last');
}
