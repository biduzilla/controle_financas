package com.ricky.controle_financas.presentation.auth.login

data class LoginState(
    var email: String = "",
    var password: String = "",
    var isLoading: Boolean = false,
    var error: String? = null
)
