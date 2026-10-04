import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../constants/api_constants.dart';

class ApiClient {
  final http.Client _client = http.Client();

  Future<String?> get _token async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('access_token');
  }

  Future<Map<String, String>> _headers() async {
    final token = await _token;
    return {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  Future<dynamic> get(String endpoint) async {
    final url = Uri.parse('${ApiConstants.baseUrl}$endpoint');
    final response = await _client.get(url, headers: await _headers());
    return _handleResponse(response);
  }

  Future<dynamic> post(String endpoint, Map<String, dynamic> body) async {
    final url = Uri.parse('${ApiConstants.baseUrl}$endpoint');
    final response = await _client.post(url, headers: await _headers(), body: jsonEncode(body));
    return _handleResponse(response);
  }

  Future<dynamic> patch(String endpoint, Map<String, dynamic> body) async {
    final url = Uri.parse('${ApiConstants.baseUrl}$endpoint');
    final response = await _client.patch(url, headers: await _headers(), body: jsonEncode(body));
    return _handleResponse(response);
  }

  Future<dynamic> delete(String endpoint) async {
    final url = Uri.parse('${ApiConstants.baseUrl}$endpoint');
    final response = await _client.delete(url, headers: await _headers());
    return _handleResponse(response);
  }

  dynamic _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      if (response.body.isEmpty) return null;
      return jsonDecode(response.body);
    } else {
      try {
        final errorJson = jsonDecode(response.body);
        final message = errorJson['error']?['message'] ?? errorJson['message'] ?? 'Network Error';
        throw Exception(message);
      } catch (e) {
        throw Exception('Server Error: ${response.statusCode}');
      }
    }
  }
}
