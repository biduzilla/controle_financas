package com.ricky.controle_financas.presentation.screens.auth.register

sealed class RegisterNavigation {
    object NavigateToLogin : RegisterNavigation()
    object NavigateToHome : RegisterNavigation()
}