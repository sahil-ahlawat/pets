class User {
  final String id;
  final String email;
  final String firstName;
  final String lastName;
  final String? phoneNumber;
  final String role;
  final String? avatarUrl;
  final bool allowDirectCalls;

  User({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    this.phoneNumber,
    required this.role,
    this.avatarUrl,
    required this.allowDirectCalls,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      email: json['email'] ?? '',
      firstName: json['firstName'] ?? '',
      lastName: json['lastName'] ?? '',
      phoneNumber: json['phoneNumber'],
      role: json['role'] ?? 'PET_OWNER',
      avatarUrl: json['avatarUrl'],
      allowDirectCalls: json['allowDirectCalls'] ?? false,
    );
  }
}
