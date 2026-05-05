package com.ricky.controle_financas.navigation

import androidx.compose.animation.core.tween
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.togetherWith
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import androidx.lifecycle.viewmodel.navigation3.rememberViewModelStoreNavEntryDecorator
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.runtime.rememberSaveableStateHolderNavEntryDecorator
import androidx.navigation3.ui.NavDisplay
import com.ricky.controle_financas.presentation.auth.login.LoginRoute
import com.ricky.controle_financas.presentation.auth.login.LoginViewModel
import com.ricky.controle_financas.presentation.splash.SplashRoute
import com.ricky.controle_financas.presentation.splash.SplashViewModel

@Composable
fun AppNavDisplay(
    modifier: Modifier = Modifier,
    navigationViewModel: NavigationViewModel
) {
    NavDisplay(
        backStack = navigationViewModel.backStack,
        modifier = modifier,
        transitionSpec = {
            fadeIn(tween(300)) togetherWith fadeOut(tween(300))
        },
        entryDecorators = listOf(
            rememberSaveableStateHolderNavEntryDecorator(),
            rememberViewModelStoreNavEntryDecorator()
        ),
        entryProvider = entryProvider {
            entry<AppNavKey.Splash> {
                val viewModel = hiltViewModel<SplashViewModel>()
                val state by viewModel.state.collectAsState()

                SplashRoute(
                    state = state,
                    navigation = viewModel.navigation,
                    onEvent = viewModel::onEvent,
                    onNavigate = { key ->
                        navigationViewModel.navigateAndClear(key)
                    })
            }
            entry<AppNavKey.Login> {
                val viewModel = hiltViewModel<LoginViewModel>()
                val state by viewModel.state.collectAsState()
                LoginRoute(
                    state = state,
                    navigation = viewModel.navigation,
                    onEvent = viewModel::onEvent,
                    onNavigate = { key ->
                        navigationViewModel.navigateAndClear(key)
                    }
                )
            }

        }
    )
}