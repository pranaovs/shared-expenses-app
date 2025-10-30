class UserExpenseBreakdown {
  final String expenseId;
  final String title;
  final String? description;
  final int createdAt;
  final double totalAmount;
  final double amountPaid;
  final double amountOwed;
  final double netSpending;

  UserExpenseBreakdown({
    required this.expenseId,
    required this.title,
    this.description,
    required this.createdAt,
    required this.totalAmount,
    required this.amountPaid,
    required this.amountOwed,
    required this.netSpending,
  });

  factory UserExpenseBreakdown.fromJson(Map<String, dynamic> json) {
    return UserExpenseBreakdown(
      expenseId: json['expense_id'] as String,
      title: json['title'] as String,
      description: json['description'] as String?,
      createdAt: json['created_at'] as int,
      totalAmount: (json['total_amount'] as num).toDouble(),
      amountPaid: (json['amount_paid'] as num).toDouble(),
      amountOwed: (json['amount_owed'] as num).toDouble(),
      netSpending: (json['net_spending'] as num).toDouble(),
    );
  }

  DateTime get createdAtDateTime => DateTime.fromMillisecondsSinceEpoch(createdAt * 1000);

  bool get isOwing => netSpending < 0;
  bool get isPaidExtra => netSpending > 0;
}
