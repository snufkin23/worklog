import 'package:riverpod_annotation/riverpod_annotation.dart';

import 'app_state.dart';

part 'app_controller.g.dart';

@Riverpod(keepAlive: true)
class AppController extends _$AppController {
  @override
  AppState build() {
    return const AppState.guest();
  }

  void setGuest() {
    state = const AppState.guest();
  }

  void setUnauthenticated({String? message}) {
    state = AppState.unauthenticated(message: message);
  }

  void setAuthenticated() {
    state = const AppState.authenticated();
  }

  void logout() {
    state = const AppState.unauthenticated();
  }
}
