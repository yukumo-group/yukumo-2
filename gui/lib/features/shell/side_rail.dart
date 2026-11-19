import 'package:material_ui/material_ui.dart';

/// Width at which [MainShell] switches from floating bottom dock to a rail.
const double kShellRailBreakpoint = 840;

/// Very wide layouts use an extended rail with inline labels.
const double kShellRailExtendedBreakpoint = 1280;

class SideNavRail extends StatelessWidget {
  const SideNavRail({
    super.key,
    required this.selectedIndex,
    required this.onDestinationSelected,
    this.extended = false,
  });

  final int selectedIndex;
  final ValueChanged<int> onDestinationSelected;
  final bool extended;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return NavigationRail(
      selectedIndex: selectedIndex,
      onDestinationSelected: onDestinationSelected,
      extended: extended,
      labelType: extended
          ? NavigationRailLabelType.none
          : NavigationRailLabelType.all,
      backgroundColor: cs.surface,
      groupAlignment: -1,
      destinations: const [
        NavigationRailDestination(
          icon: Icon(Icons.music_note_outlined),
          selectedIcon: Icon(Icons.music_note_rounded),
          label: Text('合成'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.groups_outlined),
          selectedIcon: Icon(Icons.groups_rounded),
          label: Text('角色'),
        ),
      ],
    );
  }
}
