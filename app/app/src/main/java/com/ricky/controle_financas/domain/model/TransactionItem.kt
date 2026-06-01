package com.ricky.controle_financas.domain.model

import java.util.UUID

data class TransactionItem(
    val id: UUID,
    val description: String,
    val amount: Double,
    val category: CategoryItem,
    val categoryName: String,
    val categoryColor: String,  // RECEITA or DESPESA
    val createdAt: String,
)
