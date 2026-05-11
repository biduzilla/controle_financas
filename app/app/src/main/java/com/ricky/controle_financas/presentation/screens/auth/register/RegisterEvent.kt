package com.ricky.controle_financas.presentation.screens.auth.register

sealed interface RegisterEvent {
    data class OnChangeName(val name: String) : RegisterEvent
    data class OnChangeEmail(val email: String) : RegisterEvent
    data class OnChangePhone(val phone: String) : RegisterEvent
    data class OnChangePassword(val password: String) : RegisterEvent
    data class OnChangeConfirmPassword(val confirmPassword: String) : RegisterEvent
    data object OnRegister : RegisterEvent
    data object ClearError : RegisterEvent
    data class OnNavigate(val route: RegisterNavigation) : RegisterEvent
}