package com.ricky.controle_financas.presentation.screens.auth.login

sealed interface LoginEvent {
    data class OnChangeEmail(var email: String) : LoginEvent
    data class OnChangeSenha(var senha: String) : LoginEvent
    data object OnLogin : LoginEvent
    data object ClearError : LoginEvent
    data class OnNavigate(var route: LoginNavigation) : LoginEvent
}