import 'package:flutter/material.dart';
import '../core/constants/api_constants.dart';
import '../core/network/api_client.dart';
import '../models/pet.dart';
import '../models/vaccination.dart';

class PetProvider extends ChangeNotifier {
  final ApiClient _api = ApiClient();
  List<Pet> _pets = [];
  Pet? _activePet;
  List<Vaccination> _vaccinations = [];
  bool _isLoading = false;

  List<Pet> get pets => _pets;
  Pet? get activePet => _activePet;
  List<Vaccination> get vaccinations => _vaccinations;
  bool get isLoading => _isLoading;

  Future<void> loadPets() async {
    _isLoading = true;
    notifyListeners();
    try {
      final res = await _api.get(ApiConstants.pets);
      final list = (res['pets'] as List?) ?? [];
      _pets = list.map((json) => Pet.fromJson(json)).toList();

      if (_pets.isNotEmpty) {
        if (_activePet == null || !_pets.any((p) => p.id == _activePet!.id)) {
          _activePet = _pets.first;
        } else {
          _activePet = _pets.firstWhere((p) => p.id == _activePet!.id);
        }
        await loadActivePetDetails();
      } else {
        _activePet = null;
      }
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }

  void setActivePet(Pet pet) {
    _activePet = pet;
    notifyListeners();
    loadActivePetDetails();
  }

  Future<void> loadActivePetDetails() async {
    if (_activePet == null) return;
    try {
      final res = await _api.get('${ApiConstants.pets}/${_activePet!.id}');
      _activePet = Pet.fromJson(res);

      // Load vaccinations
      final vacRes = await _api.get('${ApiConstants.pets}/${_activePet!.id}/vaccinations');
      final vacList = (vacRes as List?) ?? [];
      _vaccinations = vacList.map((j) => Vaccination.fromJson(j)).toList();

      notifyListeners();
    } catch (_) {}
  }

  Future<Pet> createPet(Map<String, dynamic> data) async {
    _isLoading = true;
    notifyListeners();
    try {
      final res = await _api.post(ApiConstants.pets, data);
      final newPet = Pet.fromJson(res);
      await loadPets();
      _activePet = newPet;
      notifyListeners();
      return newPet;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> markLost(String address) async {
    if (_activePet == null) return;
    await _api.post('${ApiConstants.pets}/${_activePet!.id}/lost', {
      'lastSeenAddress': address,
    });
    await loadActivePetDetails();
  }

  Future<void> markReunited() async {
    if (_activePet == null) return;
    await _api.post('${ApiConstants.pets}/${_activePet!.id}/reunited', {});
    await loadActivePetDetails();
  }
}
