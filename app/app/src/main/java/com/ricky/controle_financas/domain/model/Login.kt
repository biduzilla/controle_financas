package com.ricky.controle_financas.domain.model

import java.util.UUID

data class LoginRequest(
    val email: String,
    val password: String
)

data class LoginResponse(
    val accessToken: String,
    val refreshToken: String,
    val expiration: Long,
)
