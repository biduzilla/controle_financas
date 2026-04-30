package com.ricky.controle_financas.presentation.auth.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.ricky.controle_financas.domain.use_case.LoginUserCase
import com.ricky.controle_financas.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class LoginViewModel @Inject constructor(
    private val loginUserCase: LoginUserCase
) : ViewModel() {
    private val _state = MutableStateFlow(LoginState())
    val state = _state.asStateFlow()

    private val _navigation = MutableSharedFlow<LoginNavigation>(extraBufferCapacity = 1)
    val navigation: SharedFlow<LoginNavigation> = _navigation.asSharedFlow()

    fun onEvent(event: LoginEvent) {
        when (event) {
            LoginEvent.ClearError -> _state.update {
                it.copy(
                    error = null
                )
            }

            is LoginEvent.OnChangeEmail -> _state.update {
                it.copy(
                    email = event.email
                )
            }

            is LoginEvent.OnChangeSenha -> _state.update {
                it.copy(
                    password = event.senha
                )
            }

            LoginEvent.OnLogin -> login()
            LoginEvent.NavigateForgetPassword -> _navigation.tryEmit(LoginNavigation.NavigateToForgetPassword)
            LoginEvent.NavigateRegister -> _navigation.tryEmit(LoginNavigation.NavigateToRegister)
            LoginEvent.NavigateHome -> _navigation.tryEmit(LoginNavigation.NavigateToHome)
        }
    }

    private fun login() {
        val currentState = _state.value

        if (currentState.isLoading) return

        if (currentState.email.isBlank()) {
            _state.update { it.copy(error = "Email é obrigatório") }
            return
        }
        if (currentState.password.length < 6) {
            _state.update { it.copy(error = "Senha deve ter no mínimo 6 caracteres") }
            return
        }
        viewModelScope.launch {
            loginUserCase(currentState.email, currentState.password)
                .catch { e ->
                    _state.update {
                        it.copy(
                            isLoading = false,
                            error = e.localizedMessage ?: "Erro inesperado"
                        )
                    }
                }
                .collect { resource ->
                    when (resource) {
                        is Resource.Error -> {
                            _state.update {
                                it.copy(
                                    isLoading = false,
                                    error = resource.message ?: "Error"
                                )
                            }
                        }

                        is Resource.Loading -> {
                            _state.update {
                                it.copy(
                                    isLoading = true,
                                )
                            }
                        }

                        is Resource.Success -> {
                            _state.update {
                                it.copy(
                                    isLoading = false,
                                    error = null
                                )
                            }
                            _navigation.tryEmit(LoginNavigation.NavigateToHome)
                        }
                    }
                }
        }
    }
}

sealed class LoginNavigation {
    object NavigateToHome : LoginNavigation()
    object NavigateToRegister : LoginNavigation()
    object NavigateToForgetPassword : LoginNavigation()
}