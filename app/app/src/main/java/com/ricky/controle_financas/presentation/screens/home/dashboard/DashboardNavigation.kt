package com.ricky.controle_financas.presentation.screens.home.dashboard

import java.util.UUID

sealed interface DashboardNavigation {
    data object ToNewTransaction : DashboardNavigation
    data object ToAddCategory : DashboardNavigation
    data class ToTransactionDetail(val transactionId: UUID) : DashboardNavigation
    data class ToCategoryDetail(val categoryId: UUID) : DashboardNavigation
}