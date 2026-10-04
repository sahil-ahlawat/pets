import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../core/theme/app_theme.dart';
import '../../providers/auth_provider.dart';
import '../../providers/pet_provider.dart';
import '../home/dashboard_screen.dart';
import 'register_screen.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    final auth = context.read<AuthProvider>();
    try {
      await auth.login(_emailController.text.trim(), _passwordController.text);
      if (!mounted) return;
      await context.read<PetProvider>().loadPets();
      if (!mounted) return;
      Navigator.pushReplacement(context, MaterialPageRoute(builder: (_) => const DashboardScreen()));
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString().replaceAll('Exception: ', ''))),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final auth = context.watch<AuthProvider>();

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const Center(child: Text('🐾', style: TextStyle(fontSize: 54))),
                  const SizedBox(height: 12),
                  Text(
                    'Welcome Back',
                    style: Theme.of(context).textTheme.headlineMedium,
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 6),
                  const Text(
                    'Sign in to manage your companion digital records',
                    style: TextStyle(color: AppTheme.textMuted),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 36),

                  TextFormField(
                    controller: _emailController,
                    decoration: const InputDecoration(labelText: 'Email Address'),
                    keyboardType: TextInputType.emailAddress,
                    validator: (v) => v == null || v.isEmpty ? 'Email required' : null,
                  ),
                  const SizedBox(height: 16),

                  TextFormField(
                    controller: _passwordController,
                    decoration: const InputDecoration(labelText: 'Password'),
                    obscureText: true,
                    validator: (v) => v == null || v.isEmpty ? 'Password required' : null,
                  ),
                  const SizedBox(height: 24),

                  ElevatedButton(
                    onPressed: auth.isLoading ? null : _submit,
                    child: auth.isLoading
                        ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                        : const Text('Sign In'),
                  ),
                  const SizedBox(height: 16),

                  TextButton(
                    onPressed: () {
                      Navigator.push(context, MaterialPageRoute(builder: (_) => const RegisterScreen()));
                    },
                    child: const Text('Don\'t have an account? Create Free Account'),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
