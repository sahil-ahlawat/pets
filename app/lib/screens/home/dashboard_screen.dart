import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../core/theme/app_theme.dart';
import '../../providers/auth_provider.dart';
import '../../providers/pet_provider.dart';
import '../../providers/reminder_provider.dart';
import '../pet/add_edit_pet_screen.dart';
import '../pet/pet_detail_screen.dart';
import '../pet/qr_display_screen.dart';
import '../reminders/reminders_screen.dart';
import '../vets/vet_discovery_screen.dart';
import '../settings/profile_screen.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  int _currentIndex = 0;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<PetProvider>().loadPets();
      context.read<ReminderProvider>().loadReminders();
    });
  }

  @override
  Widget build(BuildContext context) {
    final petProv = context.watch<PetProvider>();
    final pets = petProv.pets;
    final activePet = petProv.activePet;

    final pages = [
      _buildDashboardHome(context, pets, activePet),
      const RemindersScreen(),
      const VetDiscoveryScreen(),
      const ProfileScreen(),
    ];

    return Scaffold(
      appBar: AppBar(
        title: const Text('🐾 PetIdentity'),
        actions: [
          IconButton(
            icon: const Icon(Icons.qr_code_scanner),
            tooltip: 'Tag Details',
            onPressed: () {
              if (activePet != null && activePet.activeQrToken != null) {
                Navigator.push(context, MaterialPageRoute(builder: (_) => QRDisplayScreen(pet: activePet)));
              }
            },
          ),
          IconButton(
            icon: const Icon(Icons.add),
            tooltip: 'Add Companion',
            onPressed: () {
              Navigator.push(context, MaterialPageRoute(builder: (_) => const AddEditPetScreen()));
            },
          ),
        ],
      ),
      body: pages[_currentIndex],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _currentIndex,
        onDestinationSelected: (i) => setState(() => _currentIndex = i),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.home_outlined), selectedIcon: Icon(Icons.home), label: 'Home'),
          NavigationDestination(icon: Icon(Icons.alarm_outlined), selectedIcon: Icon(Icons.alarm), label: 'Reminders'),
          NavigationDestination(icon: Icon(Icons.local_hospital_outlined), selectedIcon: Icon(Icons.local_hospital), label: 'Vets'),
          NavigationDestination(icon: Icon(Icons.person_outline), selectedIcon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }

  Widget _buildDashboardHome(BuildContext context, List<dynamic> pets, dynamic activePet) {
    final user = context.watch<AuthProvider>().currentUser;

    if (pets.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('🐶', style: TextStyle(fontSize: 64)),
              const SizedBox(height: 16),
              Text('No Companions Yet', style: Theme.of(context).textTheme.headlineMedium),
              const SizedBox(height: 8),
              const Text(
                'Register your companion to generate their lifetime QR digital collar tag.',
                textAlign: TextAlign.center,
                style: TextStyle(color: AppTheme.textMuted),
              ),
              const SizedBox(height: 24),
              ElevatedButton.icon(
                icon: const Icon(Icons.add),
                label: const Text('Register First Pet'),
                onPressed: () {
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const AddEditPetScreen()));
                },
              ),
            ],
          ),
        ),
      );
    }

    return ListView(
      padding: const EdgeInsets.all(20),
      children: [
        Text(
          'Hello, ${user?.firstName ?? 'Friend'} 👋',
          style: Theme.of(context).textTheme.headlineMedium,
        ),
        const SizedBox(height: 4),
        const Text('Here is your companion health & security overview', style: TextStyle(color: AppTheme.textMuted)),
        const SizedBox(height: 24),

        // Active Pet Card
        if (activePet != null) ...[
          Card(
            clipBehavior: Clip.antiAlias,
            child: InkWell(
              onTap: () {
                Navigator.push(context, MaterialPageRoute(builder: (_) => PetDetailScreen(pet: activePet)));
              },
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  children: [
                    CircleAvatar(
                      radius: 32,
                      backgroundColor: AppTheme.primaryLight,
                      child: Text(
                        activePet.species == 'CAT' ? '🐱' : '🐶',
                        style: const TextStyle(fontSize: 28),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Text(activePet.name, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w700)),
                              const SizedBox(width: 8),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                decoration: BoxDecoration(
                                  color: activePet.status == 'LOST' ? AppTheme.danger : AppTheme.primaryLight,
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Text(
                                  activePet.status,
                                  style: TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.w700,
                                    color: activePet.status == 'LOST' ? Colors.white : AppTheme.primaryDark,
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 4),
                          Text('${activePet.breed ?? activePet.species} • ${activePet.weightKg ?? '--'} kg',
                              style: const TextStyle(color: AppTheme.textMuted, fontSize: 13)),
                        ],
                      ),
                    ),
                    const Icon(Icons.arrow_forward_ios, size: 16, color: AppTheme.textSubtle),
                  ],
                ),
              ),
            ),
          ),
          const SizedBox(height: 20),

          // Fast Actions Row
          Row(
            children: [
              Expanded(
                child: OutlinedButton.icon(
                  icon: const Icon(Icons.qr_code),
                  label: const Text('Collar QR'),
                  onPressed: () {
                    Navigator.push(context, MaterialPageRoute(builder: (_) => QRDisplayScreen(pet: activePet)));
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: ElevatedButton.icon(
                  icon: const Icon(Icons.medical_services_outlined),
                  label: const Text('Records'),
                  onPressed: () {
                    Navigator.push(context, MaterialPageRoute(builder: (_) => PetDetailScreen(pet: activePet)));
                  },
                ),
              ),
            ],
          ),
        ],

        const SizedBox(height: 28),
        Text('All Companions (${pets.length})', style: Theme.of(context).textTheme.titleLarge),
        const SizedBox(height: 12),

        ...pets.map((p) => ListTile(
              leading: Text(p.species == 'CAT' ? '🐱' : '🐶', style: const TextStyle(fontSize: 24)),
              title: Text(p.name, style: const TextStyle(fontWeight: FontWeight.w600)),
              subtitle: Text('${p.breed ?? p.species} • Status: ${p.status}'),
              trailing: const Icon(Icons.chevron_right),
              onTap: () {
                context.read<PetProvider>().setActivePet(p);
                Navigator.push(context, MaterialPageRoute(builder: (_) => PetDetailScreen(pet: p)));
              },
            )),
      ],
    );
  }
}
