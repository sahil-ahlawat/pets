import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import '../../core/theme/app_theme.dart';
import '../../models/pet.dart';

class QRDisplayScreen extends StatelessWidget {
  final Pet pet;
  const QRDisplayScreen({super.key, required this.pet});

  @override
  Widget build(BuildContext context) {
    final token = pet.activeQrToken ?? 'demo-token';
    final qrUrl = 'https://petidentity.com/p/$token';

    return Scaffold(
      appBar: AppBar(title: Text('${pet.name}\'s QR Tag')),
      body: Center(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Card(
                elevation: 4,
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    children: [
                      Text(
                        pet.name,
                        style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        'Smart Digital Collar Tag',
                        style: TextStyle(color: AppTheme.textMuted, fontSize: 13),
                      ),
                      const SizedBox(height: 20),

                      QrImageView(
                        data: qrUrl,
                        version: QrVersions.auto,
                        size: 220.0,
                        backgroundColor: Colors.white,
                      ),
                      const SizedBox(height: 20),

                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                        decoration: BoxDecoration(
                          color: AppTheme.bgSubtle,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          'TOKEN: $token',
                          style: const TextStyle(fontFamily: 'monospace', fontSize: 12, fontWeight: FontWeight.bold),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),

              const Text(
                'Attach this QR code to your companion\'s collar.\nAnyone scanning with a smartphone camera will see the safe emergency profile.',
                textAlign: TextAlign.center,
                style: TextStyle(color: AppTheme.textMuted, fontSize: 13),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
