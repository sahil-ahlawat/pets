class PetReminder {
  final String id;
  final String petId;
  final String? petName;
  final String reminderType;
  final String title;
  final String? description;
  final DateTime dueDate;
  final String recurrence;
  final bool isCompleted;

  PetReminder({
    required this.id,
    required this.petId,
    this.petName,
    required this.reminderType,
    required this.title,
    this.description,
    required this.dueDate,
    required this.recurrence,
    required this.isCompleted,
  });

  factory PetReminder.fromJson(Map<String, dynamic> json) {
    return PetReminder(
      id: json['id'] ?? '',
      petId: json['petId'] ?? '',
      petName: json['petName'],
      reminderType: json['reminderType'] ?? 'CUSTOM',
      title: json['title'] ?? '',
      description: json['description'],
      dueDate: DateTime.tryParse(json['dueDate'] ?? '') ?? DateTime.now(),
      recurrence: json['recurrence'] ?? 'ONCE',
      isCompleted: json['isCompleted'] ?? false,
    );
  }
}
