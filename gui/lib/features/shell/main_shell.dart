import 'package:flutter_floating_bottom_bar/flutter_floating_bottom_bar.dart';
import 'package:material_ui/material_ui.dart';

import '../characters/characters_page.dart';
import '../synth/synth_page.dart';
import 'bottom_dock.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final width = MediaQuery.sizeOf(context).width;
    final barWidth = width < 520 ? width - 32 : 360.0;

    return Scaffold(
      backgroundColor: cs.surface,
      body: BottomBar(
        layout: BottomBarLayout(
          width: barWidth,
          maxWidth: 420,
          offset: 16,
          borderRadius: const BorderRadius.all(Radius.circular(32)),
          clip: Clip.none,
        ),
        motion: const BottomBarMotion.cupertino(
          preset: BottomBarCupertinoMotion.snappy,
          duration: Duration(milliseconds: 460),
          extraBounce: 0.03,
        ),
        scrollBehavior: const BottomBarScrollBehavior(
          hideOnScroll: true,
          showAtStart: true,
        ),
        theme: BottomBarThemeData(
          barDecoration: BoxDecoration(
            color: Colors.transparent,
            borderRadius: BorderRadius.circular(16),
            boxShadow: [
              BoxShadow(
                color: cs.shadow.withValues(alpha: 0.18),
                blurRadius: 28,
                offset: const Offset(0, 14),
              ),
            ],
          ),
        ),
        body: IndexedStack(
          index: _index,
          children: const [
            SynthPage(),
            CharactersPage(),
          ],
        ),
        child: FloatingBottomDock(
          selectedIndex: _index,
          onDestinationSelected: (i) => setState(() => _index = i),
        ),
      ),
    );
  }
}
