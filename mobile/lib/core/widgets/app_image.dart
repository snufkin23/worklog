import 'package:flutter/widgets.dart';
import 'package:flutter_svg/flutter_svg.dart';

class AppImage extends StatelessWidget {
  const AppImage(
    this.path, {
    super.key,
    this.width,
    this.height,
    this.fit,
    this.alignment = Alignment.center,
    this.filterQuality = FilterQuality.medium,
    this.semanticLabel,
    this.color,
    this.colorBlendMode,
    this.cacheWidth,
    this.cacheHeight,
  });

  final String path;
  final double? width;
  final double? height;
  final BoxFit? fit;
  final AlignmentGeometry alignment;
  final FilterQuality filterQuality;
  final String? semanticLabel;
  final Color? color;
  final BlendMode? colorBlendMode;
  final int? cacheWidth;
  final int? cacheHeight;

  @override
  Widget build(BuildContext context) => Image.asset(
    path,
    width: width,
    height: height,
    fit: fit,
    alignment: alignment,
    filterQuality: filterQuality,
    semanticLabel: semanticLabel,
    color: color,
    colorBlendMode: colorBlendMode,
    cacheWidth: cacheWidth,
    cacheHeight: cacheHeight,
  );
}

class SvgImage extends StatelessWidget {
  const SvgImage(
    this.path, {
    super.key,
    this.width,
    this.height,
    this.fit = BoxFit.contain,
    this.alignment = Alignment.center,
    this.colorFilter,
    this.semanticsLabel,
  });

  final String path;
  final double? width;
  final double? height;
  final BoxFit fit;
  final AlignmentGeometry alignment;
  final ColorFilter? colorFilter;
  final String? semanticsLabel;

  @override
  Widget build(BuildContext context) => SvgPicture.asset(
    path,
    width: width,
    height: height,
    fit: fit,
    alignment: alignment,
    colorFilter: colorFilter,
    semanticsLabel: semanticsLabel,
  );
}
