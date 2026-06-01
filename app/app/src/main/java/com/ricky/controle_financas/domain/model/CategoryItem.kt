package com.ricky.controle_financas.domain.model

import java.util.UUID

data class CategoryItem(
    val id: UUID,
    val name: String,
    val type: CategoryType,
    val color: String, // Hex color
)

enum class CategoryType {
    RECEITA,
    DESPESA,
    UNKNOWN;

    companion object {
        fun fromString(type: String): CategoryType {
            return try {
                valueOf(type.uppercase())
            } catch (e: IllegalArgumentException) {
                UNKNOWN
            }
        }
    }
}
