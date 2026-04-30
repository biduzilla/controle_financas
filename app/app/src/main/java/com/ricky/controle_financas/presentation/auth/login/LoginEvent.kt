package com.ricky.controle_financas.presentation.auth.login

sealed interface LoginEvent {
    data class OnChangeEmail(var email: String) : LoginEvent
    data class OnChangeSenha(var senha: String) : LoginEvent
    data object OnLogin : LoginEvent
    data object ClearError : LoginEvent
    data object NavigateHome : LoginEvent
    data object NavigateForgetPassword : LoginEvent
    data object NavigateRegister : LoginEvent
}