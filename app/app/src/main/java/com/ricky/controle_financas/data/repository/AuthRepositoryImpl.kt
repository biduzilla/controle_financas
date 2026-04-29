package com.ricky.controle_financas.data.repository

import com.ricky.controle_financas.data.api.AuthApi
import com.ricky.controle_financas.domain.model.LoginRequest
import com.ricky.controle_financas.domain.model.LoginResponse
import com.ricky.controle_financas.domain.repository.AuthRepository
import com.ricky.controle_financas.utils.Resource
import javax.inject.Inject

class AuthRepositoryImpl @Inject constructor(
    private val api: AuthApi
) : AuthRepository {
    override suspend fun login(request: LoginRequest): Resource<LoginResponse> {
        return safeApiCall {
            api.login(request)
        }
    }
}