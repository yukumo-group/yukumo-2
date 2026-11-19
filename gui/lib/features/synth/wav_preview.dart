import 'dart:async';
import 'dart:io';

Process? _activePlayer;
Directory? _activeDir;

/// Writes [bytes] to a temp WAV and starts a system player without blocking
/// the Dart event loop on playback completion.
Future<void> playWavBytes(List<int> bytes) async {
  await stopPreview();
  final dir = await Directory.systemTemp.createTemp('yukumo_preview_');
  final file = File('${dir.path}/preview.wav');
  await file.writeAsBytes(bytes, flush: true);
  try {
    final process = await startWavPlayer(file.path);
    _activePlayer = process;
    _activeDir = dir;
    unawaited(
      process.exitCode.whenComplete(() {
        if (identical(_activePlayer, process)) {
          _activePlayer = null;
        }
        if (identical(_activeDir, dir)) {
          _activeDir = null;
        }
        try {
          dir.deleteSync(recursive: true);
        } catch (_) {}
      }),
    );
  } catch (e) {
    try {
      await dir.delete(recursive: true);
    } catch (_) {}
    rethrow;
  }
}

Future<void> stopPreview() async {
  final process = _activePlayer;
  final dir = _activeDir;
  _activePlayer = null;
  _activeDir = null;
  if (process != null) {
    try {
      process.kill();
    } catch (_) {}
    try {
      await process.exitCode;
    } catch (_) {}
  }
  if (dir != null) {
    try {
      await dir.delete(recursive: true);
    } catch (_) {}
  }
}

Future<Process> startWavPlayer(String path) async {
  final candidates = <List<String>>[];
  if (Platform.isMacOS) {
    candidates.add(['afplay', path]);
  } else if (Platform.isWindows) {
    // PlaySync blocks the PowerShell process only; we do not await exitCode.
    candidates.add([
      'powershell',
      '-NoProfile',
      '-Command',
      "(New-Object Media.SoundPlayer '$path').PlaySync()",
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
    final exe = cmd.first;
    if (exe.startsWith('/') && !File(exe).existsSync()) {
      continue;
    }
    try {
      return await Process.start(exe, cmd.sublist(1));
    } catch (e) {
      last = e;
    }
  }
  throw StateError('no audio player worked: $last');
}
