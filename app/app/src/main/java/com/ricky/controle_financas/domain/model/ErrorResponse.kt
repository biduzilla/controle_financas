package com.ricky.controle_financas.domain.model

import com.ricky.controle_financas.utils.FlexibleMessageSerializer
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json

data class ErrorResponse(
    @Serializable(with = FlexibleMessageSerializer::class)
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
