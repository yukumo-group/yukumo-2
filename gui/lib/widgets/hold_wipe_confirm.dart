import 'package:material_ui/material_ui.dart';

/// Hold-to-confirm control: a linear color wipe fills as progress, then fires
/// [onConfirmed]. Releasing early cancels and resets.
class HoldWipeConfirm extends StatefulWidget {
  const HoldWipeConfirm({
    super.key,
    required this.child,
    required this.onConfirmed,
    required this.wipeColor,
    this.backgroundColor,
    this.borderRadius,
    this.holdDuration = const Duration(milliseconds: 750),
    this.enabled = true,
    this.width,
    this.height,
  });

  final Widget child;
  final VoidCallback onConfirmed;
  final Color wipeColor;
  final Color? backgroundColor;
  final BorderRadius? borderRadius;
  final Duration holdDuration;
  final bool enabled;
  final double? width;
  final double? height;

  @override
  State<HoldWipeConfirm> createState() => _HoldWipeConfirmState();
}

class _HoldWipeConfirmState extends State<HoldWipeConfirm>
    with SingleTickerProviderStateMixin {
  late final AnimationController _progress;
  bool _fired = false;

  @override
  void initState() {
    super.initState();
    _progress = AnimationController(vsync: this, duration: widget.holdDuration)
      ..addStatusListener(_onStatus);
  }

  @override
  void didUpdateWidget(covariant HoldWipeConfirm oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.holdDuration != widget.holdDuration) {
      _progress.duration = widget.holdDuration;
    }
  }

  @override
  void dispose() {
    _progress.removeStatusListener(_onStatus);
    _progress.dispose();
    super.dispose();
  }

  void _onStatus(AnimationStatus status) {
    if (status != AnimationStatus.completed || _fired) return;
    _fired = true;
    widget.onConfirmed();
  }

  void _start() {
    if (!widget.enabled || _fired) return;
    _progress.forward(from: 0);
  }

  void _cancel() {
    if (_fired) return;
    _progress.animateBack(0, duration: const Duration(milliseconds: 160));
  }

  @override
  Widget build(BuildContext context) {
    final radius = widget.borderRadius ?? BorderRadius.circular(12);
    return SizedBox(
      width: widget.width,
      height: widget.height,
      child: Material(
        type: MaterialType.transparency,
        child: Listener(
          behavior: HitTestBehavior.opaque,
          onPointerDown: (_) => _start(),
          onPointerUp: (_) => _cancel(),
          onPointerCancel: (_) => _cancel(),
          child: ClipRRect(
            borderRadius: radius,
            child: AnimatedBuilder(
              animation: _progress,
              builder: (context, child) {
                return Stack(
                  fit: StackFit.passthrough,
                  children: [
                    if (widget.backgroundColor != null)
                      Positioned.fill(
                        child: ColoredBox(color: widget.backgroundColor!),
                      ),
                    Positioned.fill(
                      child: Align(
                        alignment: Alignment.centerLeft,
                        child: FractionallySizedBox(
                          widthFactor: _progress.value.clamp(0.0, 1.0),
                          heightFactor: 1,
                          child: ColoredBox(color: widget.wipeColor),
                        ),
                      ),
                    ),
                    child!,
                  ],
                );
              },
              child: widget.child,
            ),
          ),
        ),
      ),
    );
  }
}
