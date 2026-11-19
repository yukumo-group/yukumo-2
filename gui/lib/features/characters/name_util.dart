import '../../gen/rpc.dart';

String uniqueCopyName(String base, Iterable<String> existing) {
  final taken = existing.toSet();
  var name = '$base - copy';
  var n = 2;
  while (taken.contains(name)) {
    name = '$base - copy $n';
    n++;
  }
  return name;
}

Character cloneCharacter(Character source, String newName) {
  return Character(
    name: newName,
    description: source.description,
    voices: [
      for (final v in source.voices ?? const <Voice>[])
        Voice(
          name: v.name,
          params: [
            for (final p in v.params ?? const <SynthParam>[])
              SynthParam(
                version: p.version,
                preset: p.preset,
                speed: p.speed,
                phontData: p.phontData,
                base: p.base,
                volume: p.volume,
                pitch: p.pitch,
                accent: p.accent,
                lmd: p.lmd,
                fsc: p.fsc,
              ),
          ],
        ),
    ],
  );
}
