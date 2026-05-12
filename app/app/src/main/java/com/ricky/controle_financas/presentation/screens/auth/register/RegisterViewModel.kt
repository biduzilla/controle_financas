package com.ricky.controle_financas.presentation.screens.auth.register

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.ricky.controle_financas.domain.model.User
import com.ricky.controle_financas.domain.use_case.LoginUserCase
import com.ricky.controle_financas.domain.use_case.SaveUserUseCase
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
class RegisterViewModel @Inject constructor(
    private val saveUserUseCase: SaveUserUseCase,
    private val loginUserCase: LoginUserCase,
) : ViewModel() {
    private val _state = MutableStateFlow(RegisterState())
    val state = _state.asStateFlow()

    private val _navigation = Channel<RegisterNavigation>(Channel.BUFFERED)
    val navigation = _navigation.receiveAsFlow()

    fun onEvent(event: RegisterEvent) {
        when (event) {
            RegisterEvent.ClearError -> _state.update { it.copy(error = null) }
            is RegisterEvent.OnChangeName -> _state.update { it.copy(name = event.name) }
            is RegisterEvent.OnChangeEmail -> _state.update { it.copy(email = event.email) }
            is RegisterEvent.OnChangePassword -> _state.update { it.copy(password = event.password) }
            is RegisterEvent.OnChangeConfirmPassword -> _state.update { it.copy(confirmPassword = event.confirmPassword) }
            RegisterEvent.OnRegister -> register()
            is RegisterEvent.OnNavigate -> navigate(event.route)
            is RegisterEvent.OnChangePhone -> _state.update { it.copy(phone = event.phone) }
        }
    }

    private fun navigate(route: RegisterNavigation) {
        viewModelScope.launch {
            _navigation.send(route)
        }
    }

    private fun register() {
        val currentState = _state.value
        if (currentState.isLoading) return

        when {
            currentState.name.isBlank() -> {
                _state.update { it.copy(error = "Nome é obrigatório") }
                return
            }

            currentState.email.isBlank() -> {
                _state.update { it.copy(error = "Email é obrigatório") }
                return
            }

            currentState.phone.isBlank() -> {
                _state.update { it.copy(error = "Telefone é obrigatório") }
                return
            }

            currentState.password != currentState.confirmPassword -> {
                _state.update { it.copy(error = "Senhas não conferem") }
                return
            }
        }

        viewModelScope.launch {
            val user = User(
                nome = currentState.name,
                telefone = currentState.phone,
                email = currentState.email,
                senha = currentState.password
            )

            saveUserUseCase(user)
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
                        is Resource.Loading -> {
                            _state.update { it.copy(isLoading = true) }
                        }

                        is Resource.Success -> {
                            _state.update { it.copy(isLoading = false, error = null) }
                            login()
                        }

                        is Resource.Error -> {
                            _state.update {
                                it.copy(
                                    isLoading = false,
                                    error = resource.message ?: "Erro no cadastro"
                                )
                            }
                        }
                    }
                }
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
                            _navigation.send(RegisterNavigation.NavigateToHome)
                        }
                    }
                }
        }
    }

}