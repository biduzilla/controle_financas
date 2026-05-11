package com.ricky.controle_financas.presentation.screens.splash

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.ricky.controle_financas.data.local.DataStoreUtil
import com.ricky.controle_financas.domain.use_case.RefreshTokenUserCase
import com.ricky.controle_financas.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class SplashViewModel @Inject constructor(
    private val dataStoreUtil: DataStoreUtil,
    private val refreshTokenUserCase: RefreshTokenUserCase
) : ViewModel() {
    private val _state = MutableStateFlow(SplashState())
    val state = _state.asStateFlow()

    private val _navigation = Channel<SplashNavigation>(Channel.BUFFERED)
    val navigation = _navigation.receiveAsFlow()

    init {
        viewModelScope.launch {
            delay(3000)
            deciderNavigation()
        }
    }

    fun onEvent(event: SplashEvent) {
        when (event) {
            SplashEvent.OnErrorDismissed -> _state.update {
                it.copy(error = null)
            }
        }
    }

    private suspend fun deciderNavigation() {
        val refreshToken = dataStoreUtil.getRefreshToken()
        if (refreshToken.isNullOrBlank()) {
            _state.update { it.copy(isLoading = false) }
            navigate(SplashNavigation.NavigateToAuth)
            return
        }
        refreshTokenUserCase(refreshToken).collect { event ->
            when (event) {
                is Resource.Error<*> -> {
                    _state.update {
                        it.copy(
                            isLoading = false,
                            error = event.message ?: "Error"
                        )
                    }
                }

                is Resource.Loading<*> -> {
                    _state.update {
                        it.copy(
                            isLoading = true,
                        )
                    }
                }

                is Resource.Success<*> -> {
                    _state.update {
                        it.copy(
                            isLoading = false,
                            error = null
                        )
                    }
                    navigate(SplashNavigation.NavigateToHome)
                }
            }

        }
    }

    private fun navigate(route: SplashNavigation) {
        viewModelScope.launch {
            _navigation.send(route)
        }
    }
}