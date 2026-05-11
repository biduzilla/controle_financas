package com.ricky.controle_financas.presentation.screens.splash

sealed class SplashNavigation {
    object NavigateToHome : SplashNavigation()
    object NavigateToAuth : SplashNavigation()
}