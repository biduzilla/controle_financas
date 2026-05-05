package com.ricky.controle_financas.presentation.splash

sealed class SplashNavigation {
    object NavigateToHome : SplashNavigation()
    object NavigateToAuth : SplashNavigation()
}