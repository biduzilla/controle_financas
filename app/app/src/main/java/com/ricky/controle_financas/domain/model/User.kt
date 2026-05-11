package com.ricky.controle_financas.domain.model

import java.util.UUID

data class User(
    val id: UUID? = null,
    val nome: String? = null,
    val telefone: String? = null,
    val email: String? = null,
    val senha: String? = null,
    val version: Int? = null
)

data class RegisterUser(
    val id: UUID? = null,
    val nome: String? = null,
    val telefone: String? = null,
    val email: String? = null,
    val senha: String? = null,
    val version: Int? = null
)
