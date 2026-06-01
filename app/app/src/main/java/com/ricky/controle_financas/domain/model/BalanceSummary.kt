package com.ricky.controle_financas.domain.model

data class BalanceSummary(
    val totalIncome: Double = 0.0,
    val totalExpense: Double = 0.0,
    val currentBalance: Double = 0.0,
)
