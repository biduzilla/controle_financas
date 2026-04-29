package com.ricky.controle_financas.presentation.auth.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.ricky.controle_financas.domain.use_case.LoginUserCase
import com.ricky.controle_financas.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class LoginViewModel @Inject constructor(
    private val loginUserCase: LoginUserCase
) : ViewModel() {
    private val _state = MutableStateFlow(LoginState())
    private val state = _state.asStateFlow()

    private fun login(email: String, password: String) {
        if (_state.value.isLoading) return
        viewModelScope.launch {
            loginUserCase(email, password)
                .collect { resource ->
                    when (resource) {
                        is Resource.Error -> {
                            _state.value = _state.value.copy(
                                isLoading = false,
                                error = resource.message ?: "Error"
                            )
                        }

                        is Resource.Loading -> {
                            _state.value = _state.value.copy(
                                isLoading = true,
                            )
                        }

                        is Resource.Success -> {
                            _state.value = _state.value.copy(
                                isLoading = false,
                                error = null
                            )
                        }
                    }
                }
        }
    }
}