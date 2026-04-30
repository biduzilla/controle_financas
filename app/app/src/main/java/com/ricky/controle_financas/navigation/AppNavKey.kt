package com.ricky.controle_financas.navigation

import androidx.navigation3.runtime.NavKey
import kotlinx.serialization.Serializable

interface AppNavKey : NavKey {
    @Serializable
    data object Splash : AppNavKey

    @Serializable
    data object Login : AppNavKey

    @Serializable
    data object Register : AppNavKey

    @Serializable
    data object ForgotPassword : AppNavKey

    @Serializable
    data object Home : AppNavKey
}