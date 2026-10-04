import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import '../../core/theme/app_theme.dart';
import '../../providers/reminder_provider.dart';

class RemindersScreen extends StatelessWidget {
  const RemindersScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final remProv = context.watch<ReminderProvider>();
    final reminders = remProv.reminders;

    return Scaffold(
      appBar: AppBar(title: const Text('⏰ Preventive Reminders')),
      body: remProv.isLoading
          ? const Center(child: CircularProgressIndicator())
          : reminders.isEmpty
              ? const Center(
                  child: Text('No reminders scheduled.', style: TextStyle(color: AppTheme.textMuted)),
                )
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: reminders.length,
                  itemBuilder: (context, i) {
                    final r = reminders[i];
                    return Card(
                      margin: const EdgeInsets.only(bottom: 12),
                      child: ListTile(
                        leading: CircleAvatar(
                          backgroundColor: r.isCompleted ? AppTheme.primaryLight : AppTheme.accent.withOpacity(0.15),
                          child: Icon(
                            r.isCompleted ? Icons.check : Icons.alarm,
                            color: r.isCompleted ? AppTheme.primaryDark : AppTheme.accent,
                          ),
                        ),
                        title: Text(
                          r.title,
                          style: TextStyle(
                            fontWeight: FontWeight.w600,
                            decoration: r.isCompleted ? TextDecoration.lineThrough : null,
                          ),
                        ),
                        subtitle: Text(
                          'Due: ${DateFormat.yMMMd().format(r.dueDate)} • ${r.recurrence}',
                          style: const TextStyle(fontSize: 12),
                        ),
                        trailing: r.isCompleted
                            ? const Icon(Icons.done_all, color: AppTheme.primary)
                            : IconButton(
                                icon: const Icon(Icons.check_circle_outline),
                                onPressed: () {
                                  context.read<ReminderProvider>().completeReminder(r.id);
                                },
                              ),
                      ),
                    );
                  },
                ),
    );
  }
}
