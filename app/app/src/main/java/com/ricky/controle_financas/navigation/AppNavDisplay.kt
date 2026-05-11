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
import com.ricky.controle_financas.presentation.screens.auth.login.LoginRoute
import com.ricky.controle_financas.presentation.screens.auth.login.LoginViewModel
import com.ricky.controle_financas.presentation.screens.auth.register.RegisterRoute
import com.ricky.controle_financas.presentation.screens.auth.register.RegisterViewModel
import com.ricky.controle_financas.presentation.screens.splash.SplashRoute
import com.ricky.controle_financas.presentation.screens.splash.SplashViewModel

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
                    onNavigate = { destination ->
                        when (destination) {
                            AppNavKey.Register -> navigationViewModel.navigate(destination)
                            else -> navigationViewModel.navigateAndClear(destination)
                        }
                    }
                )
            }

            entry<AppNavKey.Register> {
                val viewModel = hiltViewModel<RegisterViewModel>()
                val state by viewModel.state.collectAsState()
                RegisterRoute(
                    state = state,
                    navigation = viewModel.navigation,
                    onEvent = viewModel::onEvent,
                    onNavigate = { key ->
                        navigationViewModel.navigateAndClear(key)
                    },
                    onBack = {
                        navigationViewModel.popBackStack()
                    }
                )
            }
        }
    )
}