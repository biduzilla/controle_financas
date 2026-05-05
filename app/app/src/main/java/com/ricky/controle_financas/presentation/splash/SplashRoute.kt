package com.ricky.controle_financas.presentation.splash

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Modifier
import com.ricky.controle_financas.navigation.AppNavKey
import kotlinx.coroutines.flow.Flow

@Composable
fun SplashRoute(
    modifier: Modifier = Modifier,
    state: SplashState,
    navigation: Flow<SplashNavigation>,
    onEvent: (SplashEvent) -> Unit,
    onNavigate: (AppNavKey) -> Unit,
) {
    LaunchedEffect(Unit) {
        navigation.collect { nav ->
            when (nav) {
                SplashNavigation.NavigateToAuth -> onNavigate(AppNavKey.Login)
                SplashNavigation.NavigateToHome -> onNavigate(AppNavKey.Home)
            }
        }
    }
    SplashScreen(modifier, state, onEvent)
}