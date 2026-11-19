import 'package:m3e_core/m3e_core.dart';
import 'package:material_ui/material_ui.dart';

class IntSliderField extends StatefulWidget {
  const IntSliderField({
    super.key,
    required this.label,
    required this.value,
    required this.min,
    required this.max,
    required this.onChanged,
    this.divisions,
  });

  final String label;
  final int value;
  final int min;
  final int max;
  final ValueChanged<int> onChanged;
  final int? divisions;

  @override
  State<IntSliderField> createState() => _IntSliderFieldState();
}

class _IntSliderFieldState extends State<IntSliderField> {
  late TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: '${widget.value}');
  }

  @override
  void didUpdateWidget(covariant IntSliderField oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.value != widget.value &&
        _controller.text != '${widget.value}') {
      _controller.text = '${widget.value}';
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _commitText() {
    final parsed = int.tryParse(_controller.text.trim());
    if (parsed == null) {
      _controller.text = '${widget.value}';
      return;
    }
    final clamped = parsed.clamp(widget.min, widget.max);
    _controller.text = '$clamped';
    if (clamped != widget.value) widget.onChanged(clamped);
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final divisions = widget.divisions ?? (widget.max - widget.min);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text(widget.label, style: Theme.of(context).textTheme.labelLarge),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: M3ESlider(
                value: widget.value.toDouble(),
                min: widget.min.toDouble(),
                max: widget.max.toDouble(),
                // divisions: divisions <= 0 ? null : divisions,
                motion: M3EMotion.expressiveSpatialFast,
                label: '${widget.value}',
                onChanged: (v) => widget.onChanged(v.round().clamp(widget.min, widget.max)),
                decoration: M3ESliderDecoration(
                  haptic: M3EHapticFeedback.light,
                  trackHeight: 16,
                  trackCornerRadius: 8,
                  thumbWidth: 4,
                  thumbHeight: 32,
                  colors: M3ESliderColors(
                    thumbColor: cs.primary,
                    disabledThumbColor: cs.onSurface.withValues(alpha: 0.38),
                    activeTrackColor: cs.primary,
                    inactiveTrackColor: cs.surfaceContainerHighest,
                    disabledActiveTrackColor: cs.onSurface.withValues(alpha: 0.38),
                    disabledInactiveTrackColor: cs.onSurface.withValues(alpha: 0.12),
                    activeTickColor: cs.onPrimary.withValues(alpha: 0.35),
                    inactiveTickColor: cs.primary.withValues(alpha: 0.35),
                    disabledActiveTickColor: cs.onSurface.withValues(alpha: 0.12),
                    disabledInactiveTickColor: cs.onSurface.withValues(alpha: 0.38),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            SizedBox(
              width: 72,
              child: TextField(
                controller: _controller,
                keyboardType: TextInputType.number,
                textAlign: TextAlign.center,
                decoration: InputDecoration(
                  isDense: true,
                  filled: true,
                  fillColor: cs.surfaceContainer,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide.none,
                  ),
                ),
                onEditingComplete: _commitText,
                onSubmitted: (_) => _commitText(),
              ),
            ),
          ],
        ),
      ],
    );
  }
}
