import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // Brand Colors
  static const Color primary = Color(0xFF0D9488); // Sophisticated Teal
  static const Color primaryDark = Color(0xFF0F766E);
  static const Color primaryLight = Color(0xFFCCFBF1);
  static const Color accent = Color(0xFFF59E0B); // Warm Amber
  static const Color danger = Color(0xFFE11D48); // Rose Emergency
  static const Color success = Color(0xFF10B981); // Emerald

  static const Color bgCanvas = Color(0xFFFCFBF9); // Warm artisan canvas
  static const Color bgSurface = Color(0xFFFFFFFF);
  static const Color bgSubtle = Color(0xFFF8FAFC);

  static const Color textMain = Color(0xFF0F172A);
  static const Color textMuted = Color(0xFF64748B);
  static const Color textSubtle = Color(0xFF94A3B8);
  static const Color border = Color(0xFFE2E8F0);

  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      scaffoldBackgroundColor: bgCanvas,
      colorScheme: ColorScheme.light(
        primary: primary,
        secondary: accent,
        error: danger,
        surface: bgSurface,
        onPrimary: Colors.white,
        onSurface: textMain,
      ),
      textTheme: GoogleFonts.plusJakartaSansTextTheme().copyWith(
        displayLarge: GoogleFonts.outfit(fontWeight: FontWeight.w800, color: textMain),
        displayMedium: GoogleFonts.outfit(fontWeight: FontWeight.w700, color: textMain),
        headlineLarge: GoogleFonts.outfit(fontWeight: FontWeight.w700, color: textMain),
        headlineMedium: GoogleFonts.outfit(fontWeight: FontWeight.w600, color: textMain),
        titleLarge: GoogleFonts.outfit(fontWeight: FontWeight.w600, color: textMain),
        bodyLarge: GoogleFonts.plusJakartaSans(color: textMain, fontSize: 16),
        bodyMedium: GoogleFonts.plusJakartaSans(color: textMuted, fontSize: 14),
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: bgCanvas,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        iconTheme: IconThemeData(color: textMain),
        titleTextStyle: TextStyle(
          color: textMain,
          fontSize: 20,
          fontWeight: FontWeight.w700,
        ),
      ),
      cardTheme: CardTheme(
        color: bgSurface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: const BorderSide(color: border),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: bgSurface,
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: border),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: border),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: primary, width: 2),
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: Colors.white,
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(30),
          ),
          textStyle: GoogleFonts.plusJakartaSans(
            fontWeight: FontWeight.w600,
            fontSize: 15,
          ),
        ),
      ),
    );
  }
}
