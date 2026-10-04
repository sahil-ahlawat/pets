import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/pet_provider.dart';

class AddEditPetScreen extends StatefulWidget {
  const AddEditPetScreen({super.key});

  @override
  State<AddEditPetScreen> createState() => _AddEditPetScreenState();
}

class _AddEditPetScreenState extends State<AddEditPetScreen> {
  final _nameController = TextEditingController();
  final _breedController = TextEditingController();
  final _weightController = TextEditingController();
  final _microchipController = TextEditingController();
  final _allergiesController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  String _selectedSpecies = 'DOG';
  String _selectedGender = 'MALE';
  bool _isNeutered = false;

  @override
  void dispose() {
    _nameController.dispose();
    _breedController.dispose();
    _weightController.dispose();
    _microchipController.dispose();
    _allergiesController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    final petProv = context.read<PetProvider>();
    final data = {
      'name': _nameController.text.trim(),
      'species': _selectedSpecies,
      'breed': _breedController.text.trim(),
      'gender': _selectedGender,
      'isNeutered': _isNeutered,
      'weightKg': double.tryParse(_weightController.text),
      'microchipNumber': _microchipController.text.trim(),
      'allergies': _allergiesController.text.trim(),
    };

    try {
      await petProv.createPet(data);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Companion registered successfully!')),
      );
      Navigator.pop(context);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final petProv = context.watch<PetProvider>();

    return Scaffold(
      appBar: AppBar(title: const Text('Register Companion')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(labelText: 'Pet Name *'),
                validator: (v) => v == null || v.isEmpty ? 'Name required' : null,
              ),
              const SizedBox(height: 16),

              DropdownButtonFormField<String>(
                value: _selectedSpecies,
                decoration: const InputDecoration(labelText: 'Species'),
                items: const [
                  DropdownMenuItem(value: 'DOG', child: Text('🐶 Dog')),
                  DropdownMenuItem(value: 'CAT', child: Text('🐱 Cat')),
                  DropdownMenuItem(value: 'BIRD', child: Text('🦜 Bird')),
                  DropdownMenuItem(value: 'RABBIT', child: Text('🐰 Rabbit')),
                  DropdownMenuItem(value: 'OTHER', child: Text('🐾 Other')),
                ],
                onChanged: (v) => setState(() => _selectedSpecies = v!),
              ),
              const SizedBox(height: 16),

              TextFormField(
                controller: _breedController,
                decoration: const InputDecoration(labelText: 'Breed (e.g. Golden Retriever)'),
              ),
              const SizedBox(height: 16),

              Row(
                children: [
                  Expanded(
                    child: DropdownButtonFormField<String>(
                      value: _selectedGender,
                      decoration: const InputDecoration(labelText: 'Gender'),
                      items: const [
                        DropdownMenuItem(value: 'MALE', child: Text('Male')),
                        DropdownMenuItem(value: 'FEMALE', child: Text('Female')),
                        DropdownMenuItem(value: 'UNKNOWN', child: Text('Unknown')),
                      ],
                      onChanged: (v) => setState(() => _selectedGender = v!),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextFormField(
                      controller: _weightController,
                      decoration: const InputDecoration(labelText: 'Weight (kg)'),
                      keyboardType: const TextInputType.numberWithOptions(decimal: true),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              SwitchListTile(
                title: const Text('Spayed / Neutered'),
                value: _isNeutered,
                onChanged: (v) => setState(() => _isNeutered = v),
              ),
              const SizedBox(height: 16),

              TextFormField(
                controller: _microchipController,
                decoration: const InputDecoration(labelText: '15-Digit Microchip Number'),
              ),
              const SizedBox(height: 16),

              TextFormField(
                controller: _allergiesController,
                decoration: const InputDecoration(labelText: 'Critical Medical Notes / Allergies'),
                maxLines: 3,
              ),
              const SizedBox(height: 28),

              ElevatedButton(
                onPressed: petProv.isLoading ? null : _submit,
                child: petProv.isLoading
                    ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                    : const Text('Save Companion & Generate QR Tag'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
