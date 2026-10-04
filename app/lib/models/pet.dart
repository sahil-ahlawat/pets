class Pet {
  final String id;
  final String name;
  final String species;
  final String? breed;
  final String? gender;
  final bool isNeutered;
  final double? weightKg;
  final String? avatarUrl;
  final String status;
  final String? activeQrToken;
  final String? microchipNumber;
  final String? allergies;
  final String? emergencyContactName;
  final String? emergencyContactPhone;

  Pet({
    required this.id,
    required this.name,
    required this.species,
    this.breed,
    this.gender,
    required this.isNeutered,
    this.weightKg,
    this.avatarUrl,
    required this.status,
    this.activeQrToken,
    this.microchipNumber,
    this.allergies,
    this.emergencyContactName,
    this.emergencyContactPhone,
  });

  factory Pet.fromJson(Map<String, dynamic> json) {
    return Pet(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      species: json['species'] ?? 'DOG',
      breed: json['breed'],
      gender: json['gender'],
      isNeutered: json['isNeutered'] ?? false,
      weightKg: json['weightKg'] != null ? (json['weightKg'] as num).toDouble() : null,
      avatarUrl: json['avatarUrl'],
      status: json['status'] ?? 'ACTIVE',
      activeQrToken: json['activeQrToken'],
      microchipNumber: json['microchipNumber'],
      allergies: json['allergies'],
      emergencyContactName: json['emergencyContactName'],
      emergencyContactPhone: json['emergencyContactPhone'],
    );
  }
}
