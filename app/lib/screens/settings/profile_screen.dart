import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../core/theme/app_theme.dart';
import '../../providers/auth_provider.dart';
import '../auth/login_screen.dart';

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final user = context.watch<AuthProvider>().currentUser;

    return Scaffold(
      appBar: AppBar(title: const Text('Account & Privacy')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  CircleAvatar(
                    radius: 28,
                    backgroundColor: AppTheme.primaryLight,
                    child: Text(
                      user != null && user.firstName.isNotEmpty ? user.firstName[0].toUpperCase() : 'U',
                      style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: AppTheme.primaryDark),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '${user?.firstName ?? ''} ${user?.lastName ?? ''}',
                          style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                        ),
                        Text(user?.email ?? '', style: const TextStyle(color: AppTheme.textMuted, fontSize: 13)),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 20),

          Card(
            child: Column(
              children: [
                SwitchListTile(
                  title: const Text('Allow Direct Calls on Tag'),
                  subtitle: const Text('When enabled, finders scanning your tag can call your mobile directly.'),
                  value: user?.allowDirectCalls ?? false,
                  onChanged: (val) {
                    // Update state
                  },
                ),
                const Divider(),
                ListTile(
                  leading: const Icon(Icons.notifications_active_outlined),
                  title: const Text('Push Notifications (FCM / APNs)'),
                  subtitle: const Text('Emergency scan alerts & booster reminders'),
                  trailing: const Text('Active', style: TextStyle(color: AppTheme.primary, fontWeight: FontWeight.bold)),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          OutlinedButton.icon(
            icon: const Icon(Icons.logout, color: AppTheme.danger),
            label: const Text('Sign Out', style: TextStyle(color: AppTheme.danger)),
            onPressed: () async {
              await context.read<AuthProvider>().logout();
              if (context.mounted) {
                Navigator.pushAndRemoveUntil(context, MaterialPageRoute(builder: (_) => const LoginScreen()), (_) => false);
              }
            },
          ),
        ],
      ),
    );
  }
}
