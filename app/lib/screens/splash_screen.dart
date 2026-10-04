import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../providers/pet_provider.dart';
import 'auth/login_screen.dart';
import 'home/dashboard_screen.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _checkInit();
  }

  Future<void> _checkInit() async {
    await Future.delayed(const Duration(milliseconds: 600));
    if (!mounted) return;

    final auth = context.read<AuthProvider>();
    await auth.checkAuth();

    if (!mounted) return;
    if (auth.isAuthenticated) {
      await context.read<PetProvider>().loadPets();
      if (!mounted) return;
      Navigator.pushReplacement(context, MaterialPageRoute(builder: (_) => const DashboardScreen()));
    } else {
      Navigator.pushReplacement(context, MaterialPageRoute(builder: (_) => const LoginScreen()));
    }
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('🐾', style: TextStyle(fontSize: 64)),
            SizedBox(height: 16),
            Text(
              'PetIdentity',
              style: TextStyle(fontSize: 28, fontWeight: FontWeight.w800, letterSpacing: -0.5),
            ),
            SizedBox(height: 8),
            Text(
              'Smart Digital Companion Platform',
              style: TextStyle(color: Colors.grey),
            ),
            SizedBox(height: 32),
            CircularProgressIndicator(strokeWidth: 2),
          ],
        ),
      ),
    );
  }
}
