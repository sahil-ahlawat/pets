class Vet {
  final String id;
  final String clinicName;
  final String? doctorName;
  final String phoneNumber;
  final String? emergencyPhone;
  final String addressLine1;
  final String city;
  final bool is24hEmergency;
  final double rating;
  final int reviewCount;
  final double? distanceKm;

  Vet({
    required this.id,
    required this.clinicName,
    this.doctorName,
    required this.phoneNumber,
    this.emergencyPhone,
    required this.addressLine1,
    required this.city,
    required this.is24hEmergency,
    required this.rating,
    required this.reviewCount,
    this.distanceKm,
  });

  factory Vet.fromJson(Map<String, dynamic> json) {
    return Vet(
      id: json['id'] ?? '',
      clinicName: json['clinicName'] ?? '',
      doctorName: json['doctorName'],
      phoneNumber: json['phoneNumber'] ?? '',
      emergencyPhone: json['emergencyPhone'],
      addressLine1: json['addressLine1'] ?? '',
      city: json['city'] ?? '',
      is24hEmergency: json['is24hEmergency'] ?? false,
      rating: (json['rating'] as num?)?.toDouble() ?? 5.0,
      reviewCount: json['reviewCount'] ?? 0,
      distanceKm: (json['distanceKm'] as num?)?.toDouble(),
    );
  }
}
