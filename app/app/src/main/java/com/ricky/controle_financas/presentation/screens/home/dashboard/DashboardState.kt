package com.ricky.controle_financas.presentation.screens.home.dashboard

import com.ricky.controle_financas.domain.model.BalanceSummary
import com.ricky.controle_financas.domain.model.CategoryItem
import com.ricky.controle_financas.domain.model.TransactionItem

data class DashboardState(
    val balanceSummary: BalanceSummary = BalanceSummary(),
    val recentTransactions: List<TransactionItem> = emptyList(),
    val categories: List<CategoryItem> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
)

