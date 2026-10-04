import 'package:flutter/material.dart';
import '../core/constants/api_constants.dart';
import '../core/network/api_client.dart';
import '../models/reminder.dart';

class ReminderProvider extends ChangeNotifier {
  final ApiClient _api = ApiClient();
  List<PetReminder> _reminders = [];
  bool _isLoading = false;

  List<PetReminder> get reminders => _reminders;
  bool get isLoading => _isLoading;

  Future<void> loadReminders() async {
    _isLoading = true;
    notifyListeners();
    try {
      final res = await _api.get(ApiConstants.reminders);
      final list = (res as List?) ?? [];
      _reminders = list.map((json) => PetReminder.fromJson(json)).toList();
      _isLoading = false;
      notifyListeners();
    } catch (_) {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> completeReminder(String id) async {
    try {
      await _api.post('${ApiConstants.reminders}/$id/complete', {});
      await loadReminders();
    } catch (_) {}
  }
}
