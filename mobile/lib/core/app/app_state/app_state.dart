import 'package:freezed_annotation/freezed_annotation.dart';

part 'app_state.freezed.dart';

@freezed
sealed class AppState with _$AppState {
  const factory AppState.guest() = _Guest;

  const factory AppState.unauthenticated({String? message}) = _Unauthenticated;

  const factory AppState.authenticated() = _Authenticated;
}
