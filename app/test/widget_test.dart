import 'package:flutter_test/flutter_test.dart';
import 'package:pet_identity_app/main.dart';

void main() {
  testWidgets('PetIdentity initial smoke test', (WidgetTester tester) async {
    await tester.pumpWidget(const PetIdentityApp());
    expect(find.text('PetIdentity'), findsOneWidget);
  });
}
