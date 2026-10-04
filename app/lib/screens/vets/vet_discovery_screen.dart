import 'package:flutter/material.dart';
import '../../core/constants/api_constants.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/app_theme.dart';
import '../../models/vet.dart';

class VetDiscoveryScreen extends StatefulWidget {
  const VetDiscoveryScreen({super.key});

  @override
  State<VetDiscoveryScreen> createState() => _VetDiscoveryScreenState();
}

class _VetDiscoveryScreenState extends State<VetDiscoveryScreen> {
  final ApiClient _api = ApiClient();
  List<Vet> _vets = [];
  bool _isLoading = true;
  bool _emergencyOnly = false;

  @override
  void initState() {
    super.initState();
    _fetchVets();
  }

  Future<void> _fetchVets() async {
    setState(() => _isLoading = true);
    try {
      final query = _emergencyOnly ? '?emergencyOnly=true' : '';
      final res = await _api.get('${ApiConstants.vets}$query');
      final list = (res as List?) ?? [];
      setState(() {
        _vets = list.map((j) => Vet.fromJson(j)).toList();
        _isLoading = false;
      });
    } catch (_) {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('🏥 Nearby Vets & Hospitals'),
        actions: [
          IconButton(
            icon: Icon(
              _emergencyOnly ? Icons.emergency : Icons.emergency_outlined,
              color: _emergencyOnly ? AppTheme.danger : null,
            ),
            tooltip: '24H Emergency Only',
            onPressed: () {
              setState(() => _emergencyOnly = !_emergencyOnly);
              _fetchVets();
            },
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _vets.isEmpty
              ? const Center(child: Text('No veterinary clinics found.', style: TextStyle(color: AppTheme.textMuted)))
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: _vets.length,
                  itemBuilder: (context, i) {
                    final v = _vets[i];
                    return Card(
                      margin: const EdgeInsets.only(bottom: 12),
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Expanded(
                                  child: Text(
                                    v.clinicName,
                                    style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                                  ),
                                ),
                                if (v.is24hEmergency)
                                  Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                    decoration: BoxDecoration(
                                      color: AppTheme.danger.withOpacity(0.1),
                                      borderRadius: BorderRadius.circular(12),
                                    ),
                                    child: const Text(
                                      '24H EMERGENCY',
                                      style: TextStyle(color: AppTheme.danger, fontSize: 10, fontWeight: FontWeight.bold),
                                    ),
                                  ),
                              ],
                            ),
                            const SizedBox(height: 6),
                            if (v.doctorName != null)
                              Text('👨‍⚕️ ${v.doctorName}', style: const TextStyle(color: AppTheme.textMuted, fontSize: 13)),
                            Text('📍 ${v.addressLine1}, ${v.city}', style: const TextStyle(color: AppTheme.textMuted, fontSize: 13)),
                            const SizedBox(height: 12),

                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Row(
                                  children: [
                                    const Icon(Icons.star, color: Colors.amber, size: 16),
                                    const SizedBox(width: 4),
                                    Text('${v.rating} (${v.reviewCount})', style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
                                  ],
                                ),
                                Text(
                                  '📞 ${v.emergencyPhone ?? v.phoneNumber}',
                                  style: const TextStyle(color: AppTheme.primary, fontWeight: FontWeight.bold, fontSize: 13),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    );
                  },
                ),
    );
  }
}
