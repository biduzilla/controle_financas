package com.ricky.controle_financas.domain.model

import kotlinx.serialization.json.Json

data class ErrorResponse(
    val message: String,
    val path: String,
    val status: String,
    val timestamp: String
) {
    companion object {
        private val json = Json { ignoreUnknownKeys = true }

        fun fromJson(jsonString: String?): ErrorResponse? {
            if (jsonString.isNullOrBlank()) return null
            return try {
                json.decodeFromString<ErrorResponse>(jsonString)
            } catch (e: Exception) {
                null
            }
        }
    }
}
