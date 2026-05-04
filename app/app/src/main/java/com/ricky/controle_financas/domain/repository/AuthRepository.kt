package com.ricky.controle_financas.domain.repository

import com.ricky.controle_financas.domain.model.LoginRequest
import com.ricky.controle_financas.domain.model.LoginResponse
import com.ricky.controle_financas.domain.model.RefreshTokenRequest
import com.ricky.controle_financas.domain.model.RefreshTokenResponse
import com.ricky.controle_financas.utils.Resource

interface AuthRepository {
    suspend fun login(request: LoginRequest): Resource<LoginResponse>
    suspend fun refreshToken(): Resource<RefreshTokenResponse>
}