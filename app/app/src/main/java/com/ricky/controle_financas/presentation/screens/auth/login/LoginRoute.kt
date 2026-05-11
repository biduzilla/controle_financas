package com.ricky.controle_financas.presentation.screens.auth.login

    import androidx.compose.runtime.Composable
    import androidx.compose.runtime.LaunchedEffect
    import androidx.compose.ui.Modifier
    import com.ricky.controle_financas.navigation.AppNavKey
    import kotlinx.coroutines.flow.Flow

    @Composable
    fun LoginRoute(
        modifier: Modifier = Modifier,
        state: LoginState,
        navigation: Flow<LoginNavigation>,
        onEvent: (LoginEvent) -> Unit,
        onNavigate: (AppNavKey) -> Unit,
    ) {
        LaunchedEffect(Unit) {
            navigation.collect { event ->
                when (event) {
                    LoginNavigation.NavigateToForgetPassword -> onNavigate(AppNavKey.ForgotPassword)
                    LoginNavigation.NavigateToHome -> onNavigate(AppNavKey.Home)
                    LoginNavigation.NavigateToRegister -> onNavigate(AppNavKey.Register)
                }
            }
        }

        LoginScreen(modifier, state, onEvent)
    }