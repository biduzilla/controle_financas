package com.ricky.controle_financas.presentation.screens.home.dashboard

import java.util.UUID

sealed interface DashboardEvent {
    data object OnUserClick : DashboardEvent
    data object OnNewTransaction : DashboardEvent
    data object OnAddCategory : DashboardEvent
    data class OnTransactionClick(val transactionId: UUID) : DashboardEvent
    data class OnCategorySelected(val categoryId: UUID) : DashboardEvent
    data object OnLogout : DashboardEvent
    data object OnRetry : DashboardEvent
    data object OnErrorDismissed : DashboardEvent
}
