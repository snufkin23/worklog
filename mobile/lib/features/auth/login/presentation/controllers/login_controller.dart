import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:worklog/core/app/app_state/app_controller.dart';

part 'login_controller.g.dart';

@riverpod
class LoginController extends _$LoginController {
  @override
  bool build() {
    return false;
  }

  void login({required String email, required String password}) {
    if (email.isEmpty || password.isEmpty) {
      return;
    }

    ref.read(appControllerProvider.notifier).setAuthenticated();
  }
}
