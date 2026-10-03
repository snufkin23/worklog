import 'package:flutter/widgets.dart';

import 'core/app/app.dart';
import 'core/di/injection.dart';

void main() {
  configureDependencies();

  runApp(App());
}
