package com.ricky.controle_financas.presentation.auth.login

sealed class LoginNavigation {
    object NavigateToHome : LoginNavigation()
    object NavigateToRegister : LoginNavigation()
    object NavigateToForgetPassword : LoginNavigation()
}