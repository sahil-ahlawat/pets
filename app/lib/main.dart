import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'core/theme/app_theme.dart';
import 'providers/auth_provider.dart';
import 'providers/pet_provider.dart';
import 'providers/reminder_provider.dart';
import 'screens/splash_screen.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProvider(create: (_) => PetProvider()),
        ChangeNotifierProvider(create: (_) => ReminderProvider()),
      ],
      child: const PetIdentityApp(),
    ),
  );
}

class PetIdentityApp extends StatelessWidget {
  const PetIdentityApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'PetIdentity',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: const SplashScreen(),
      onGenerateRoute: (settings) {
        // Deep linking support: /p/:token
        final uri = Uri.parse(settings.name ?? '');
        if (uri.pathSegments.length == 2 && uri.pathSegments[0] == 'p') {
          final token = uri.pathSegments[1];
          // Open finder preview
          return MaterialPageRoute(
            builder: (_) => Scaffold(
              appBar: AppBar(title: const Text('Found Companion Tag')),
              body: Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Text('🚨', style: TextStyle(fontSize: 48)),
                      const SizedBox(height: 12),
                      const Text(
                        'Scanning Tag Token',
                        style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(height: 8),
                      Text(token, style: const TextStyle(fontFamily: 'monospace')),
                      const SizedBox(height: 24),
                      const CircularProgressIndicator(),
                    ],
                  ),
                ),
              ),
            ),
          );
        }
        return null;
      },
    );
  }
}
