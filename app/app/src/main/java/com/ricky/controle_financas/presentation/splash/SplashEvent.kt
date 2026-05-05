package com.ricky.controle_financas.presentation.splash

sealed interface SplashEvent {
    data object OnErrorDismissed : SplashEvent
}