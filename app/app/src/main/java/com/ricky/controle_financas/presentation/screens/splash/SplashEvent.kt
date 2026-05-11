package com.ricky.controle_financas.presentation.screens.splash

sealed interface SplashEvent {
    data object OnErrorDismissed : SplashEvent
}