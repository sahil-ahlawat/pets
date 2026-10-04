class ApiConstants {
  static const String baseUrl = "http://10.0.2.2:8080/api/v1"; // Standard Android Emulator loopback
  static const String liveBaseUrl = "https://api.petidentity.com/api/v1";

  // Endpoints
  static const String register = "/auth/register";
  static const String login = "/auth/login";
  static const String refresh = "/auth/refresh";
  static const String me = "/auth/me";

  static const String pets = "/pets";
  static const String vaccinationTypes = "/vaccination-types";
  static const String reminders = "/reminders";
  static const String vets = "/vets";
  static const String vetsNearby = "/vets/nearby";
  static const String publicPet = "/public/pets";
}
