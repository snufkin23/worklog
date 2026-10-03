import 'package:auto_route/auto_route.dart';
import 'package:material_ui/material_ui.dart';
import 'package:worklog/core/routes/app_router.gr.dart';
import 'package:worklog/l10n/localization.dart';

@RoutePage()
class WelcomePage extends StatelessWidget {
  const WelcomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              Text(context.l10n.welcome),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: () {
                  context.router.push(const LoginRoute());
                },
                child: Text(context.l10n.login),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
