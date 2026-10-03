import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';
import 'package:worklog/core/app/app_state/app_controller.dart';
import 'package:worklog/core/app/app_state/app_state.dart';
import 'package:worklog/core/di/injection.dart';
import 'package:worklog/core/routes/app_router.dart';

class App extends ConsumerWidget {
  App({super.key});

  final AppRouter appRouter = getIt<AppRouter>();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.listen<AppState>(appControllerProvider, (AppState? previous, AppState next) {
      next.when(
        guest: () {},
        unauthenticated: (_) {
          appRouter.replace(const LoginRoute());
        },
        authenticated: () {
          appRouter.replace(const HomeRoute());
        },
      );
    });

    return MaterialApp.router(routerConfig: appRouter.config());
  }
}
