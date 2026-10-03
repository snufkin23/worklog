import 'package:auto_route/auto_route.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:worklog/core/app/app_state/app_controller.dart';
import 'package:worklog/core/app/app_state/app_state.dart';
import 'package:worklog/core/routes/app_router.gr.dart';

class AppAuthGuard extends AutoRouteGuard {
  AppAuthGuard(this.container);

  final ProviderContainer container;

  @override
  void onNavigation(NavigationResolver resolver, StackRouter router) {
    container
        .read(appControllerProvider)
        .when(
          guest: () => resolver.next(),
          unauthenticated: (_) {
            resolver.next(false);
            router.replace(const LoginRoute());
          },
          authenticated: () => resolver.next(),
        );
  }
}
