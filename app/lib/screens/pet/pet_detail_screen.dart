import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../core/theme/app_theme.dart';
import '../../models/pet.dart';
import '../../providers/pet_provider.dart';
import 'qr_display_screen.dart';

class PetDetailScreen extends StatefulWidget {
  final Pet pet;
  const PetDetailScreen({super.key, required this.pet});

  @override
  State<PetDetailScreen> createState() => _PetDetailScreenState();
}

class _PetDetailScreenState extends State<PetDetailScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final petProv = context.watch<PetProvider>();
    final pet = petProv.activePet ?? widget.pet;
    final vacs = petProv.vaccinations;

    return Scaffold(
      appBar: AppBar(
        title: Text(pet.name),
        actions: [
          IconButton(
            icon: const Icon(Icons.qr_code),
            onPressed: () {
              Navigator.push(context, MaterialPageRoute(builder: (_) => QRDisplayScreen(pet: pet)));
            },
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppTheme.primary,
          unselectedLabelColor: AppTheme.textMuted,
          indicatorColor: AppTheme.primary,
          tabs: const [
            Tab(text: 'Overview'),
            Tab(text: 'Vaccines'),
            Tab(text: 'Emergency'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          // Tab 1: Overview
          ListView(
            padding: const EdgeInsets.all(20),
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      CircleAvatar(
                        radius: 40,
                        backgroundColor: AppTheme.primaryLight,
                        child: Text(pet.species == 'CAT' ? '🐱' : '🐶', style: const TextStyle(fontSize: 40)),
                      ),
                      const SizedBox(height: 12),
                      Text(pet.name, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold)),
                      Text('${pet.breed ?? pet.species} • Status: ${pet.status}', style: const TextStyle(color: AppTheme.textMuted)),
                      const SizedBox(height: 16),

                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          _buildStatColumn('Weight', '${pet.weightKg ?? '--'} kg'),
                          _buildStatColumn('Neutered', pet.isNeutered ? 'Yes' : 'No'),
                          _buildStatColumn('Microchip', pet.microchipNumber != null ? 'Yes' : 'No'),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              if (pet.allergies != null && pet.allergies!.isNotEmpty)
                Card(
                  color: AppTheme.danger.withOpacity(0.08),
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Row(
                          children: [
                            Icon(Icons.warning_amber_rounded, color: AppTheme.danger),
                            SizedBox(width: 8),
                            Text('Medical Alert / Allergies', style: TextStyle(fontWeight: FontWeight.bold, color: AppTheme.danger)),
                          ],
                        ),
                        const SizedBox(height: 8),
                        Text(pet.allergies!),
                      ],
                    ),
                  ),
                ),

              const SizedBox(height: 20),
              ElevatedButton.icon(
                icon: const Icon(Icons.qr_code_2),
                label: const Text('View Smart Collar QR Tag'),
                onPressed: () {
                  Navigator.push(context, MaterialPageRoute(builder: (_) => QRDisplayScreen(pet: pet)));
                },
              ),
            ],
          ),

          // Tab 2: Vaccines
          ListView(
            padding: const EdgeInsets.all(20),
            children: [
              if (vacs.isEmpty)
                const Center(
                  child: Padding(
                    padding: EdgeInsets.all(32),
                    child: Text('No vaccination records found.', style: TextStyle(color: AppTheme.textMuted)),
                  ),
                )
              else
                ...vacs.map((v) => Card(
                      margin: const EdgeInsets.only(bottom: 12),
                      child: ListTile(
                        leading: const Icon(Icons.shield_outlined, color: AppTheme.primary),
                        title: Text(v.vaccineName, style: const TextStyle(fontWeight: FontWeight.w600)),
                        subtitle: Text('Administered: ${v.administeredAt}\nNext due: ${v.nextDueDate ?? 'N/A'}'),
                        isThreeLine: true,
                      ),
                    )),
            ],
          ),

          // Tab 3: Emergency Contacts
          ListView(
            padding: const EdgeInsets.all(20),
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Secondary Emergency Contact', style: TextStyle(fontWeight: FontWeight.w700)),
                      const SizedBox(height: 12),
                      Text('Name: ${pet.emergencyContactName ?? 'None recorded'}'),
                      const SizedBox(height: 6),
                      Text('Phone: ${pet.emergencyContactPhone ?? 'None'}'),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStatColumn(String label, String value) {
    return Column(
      children: [
        Text(value, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 2),
        Text(label, style: const TextStyle(fontSize: 12, color: AppTheme.textMuted)),
      ],
    );
  }
}
