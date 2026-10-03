import 'package:auto_route/auto_route.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';
import 'package:worklog/core/app/app_state/app_controller.dart';
import 'package:worklog/l10n/localization.dart';

@RoutePage()
class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.appName)),
      body: SafeArea(
        child: Center(
          child: ElevatedButton(
            onPressed: () {
              ref.read(appControllerProvider.notifier).logout();
            },
            child: const Text('Logout'),
          ),
        ),
      ),
    );
  }
}
