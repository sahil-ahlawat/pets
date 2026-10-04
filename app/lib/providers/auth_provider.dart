import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../core/constants/api_constants.dart';
import '../core/network/api_client.dart';
import '../models/user.dart';

class AuthProvider extends ChangeNotifier {
  final ApiClient _api = ApiClient();
  User? _currentUser;
  bool _isLoading = false;
  bool _isAuthenticated = false;

  User? get currentUser => _currentUser;
  bool get isLoading => _isLoading;
  bool get isAuthenticated => _isAuthenticated;

  Future<void> checkAuth() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('access_token');
    if (token != null) {
      _isAuthenticated = true;
      try {
        final data = await _api.get(ApiConstants.me);
        if (data != null) {
          _currentUser = User.fromJson(data);
        }
      } catch (_) {
        // If token expired, clear
        _isAuthenticated = false;
        await prefs.remove('access_token');
      }
    } else {
      _isAuthenticated = false;
    }
    notifyListeners();
  }

  Future<bool> login(String email, String password) async {
    _isLoading = true;
    notifyListeners();
    try {
      final res = await _api.post(ApiConstants.login, {
        'email': email,
        'password': password,
      });

      final tokens = res['tokens'];
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', tokens['accessToken']);
      await prefs.setString('refresh_token', tokens['refreshToken']);

      _currentUser = User.fromJson(res['user']);
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }

  Future<bool> register(String email, String password, String firstName, String lastName, String? phone) async {
    _isLoading = true;
    notifyListeners();
    try {
      final res = await _api.post(ApiConstants.register, {
        'email': email,
        'password': password,
        'firstName': firstName,
        'lastName': lastName,
        'phoneNumber': phone,
      });

      final tokens = res['tokens'];
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('access_token', tokens['accessToken']);
      await prefs.setString('refresh_token', tokens['refreshToken']);

      _currentUser = User.fromJson(res['user']);
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }

  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('access_token');
    await prefs.remove('refresh_token');
    _currentUser = null;
    _isAuthenticated = false;
    notifyListeners();
  }
}
