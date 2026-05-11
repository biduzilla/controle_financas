package com.ricky.controle_financas.presentation.screens.auth.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.ricky.controle_financas.domain.use_case.LoginUserCase
import com.ricky.controle_financas.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class LoginViewModel @Inject constructor(
    private val loginUserCase: LoginUserCase
) : ViewModel() {
    private val _state = MutableStateFlow(LoginState())
    val state = _state.asStateFlow()

    private val _navigation = Channel<LoginNavigation>(Channel.BUFFERED)
    val navigation = _navigation.receiveAsFlow()

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
            is LoginEvent.OnNavigate -> navigateEvent(event.route)
        }
    }

    private fun navigateEvent(route: LoginNavigation){
        viewModelScope.launch {
            _navigation.send(route)
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
                            _navigation.send(LoginNavigation.NavigateToHome)
                        }
                    }
                }
        }
    }
}

