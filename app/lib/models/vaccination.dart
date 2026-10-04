class Vaccination {
  final String id;
  final String petId;
  final String vaccineName;
  final String administeredAt;
  final String? nextDueDate;
  final String? vetName;
  final String? clinicName;
  final String? batchLotNumber;

  Vaccination({
    required this.id,
    required this.petId,
    required this.vaccineName,
    required this.administeredAt,
    this.nextDueDate,
    this.vetName,
    this.clinicName,
    this.batchLotNumber,
  });

  factory Vaccination.fromJson(Map<String, dynamic> json) {
    return Vaccination(
      id: json['id'] ?? '',
      petId: json['petId'] ?? '',
      vaccineName: json['vaccineName'] ?? '',
      administeredAt: json['administeredAt'] ?? '',
      nextDueDate: json['nextDueDate'],
      vetName: json['vetName'],
      clinicName: json['clinicName'],
      batchLotNumber: json['batchLotNumber'],
    );
  }
}
