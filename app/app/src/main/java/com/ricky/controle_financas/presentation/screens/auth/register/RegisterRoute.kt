package com.ricky.controle_financas.presentation.screens.auth.register

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Modifier
import com.ricky.controle_financas.navigation.AppNavKey
import kotlinx.coroutines.flow.Flow

@Composable
fun RegisterRoute(
    modifier: Modifier = Modifier,
    state: RegisterState,
    navigation: Flow<RegisterNavigation>,
    onEvent: (RegisterEvent) -> Unit,
    onNavigate: (AppNavKey) -> Unit,
    onBack: () -> Unit = {}
) {
    LaunchedEffect(Unit) {
        navigation.collect { event ->
            when (event) {
                RegisterNavigation.NavigateToLogin -> onNavigate(AppNavKey.Login)
                RegisterNavigation.NavigateToHome -> onNavigate(AppNavKey.Home)
            }
        }
    }

    RegisterScreen(modifier, state, onEvent, onBack)
}