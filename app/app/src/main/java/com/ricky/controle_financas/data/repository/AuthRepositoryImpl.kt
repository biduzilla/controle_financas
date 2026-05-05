package com.ricky.controle_financas.data.repository

import com.ricky.controle_financas.data.api.AuthApi
import com.ricky.controle_financas.data.local.DataStoreUtil
import com.ricky.controle_financas.domain.model.LoginRequest
import com.ricky.controle_financas.domain.model.LoginResponse
import com.ricky.controle_financas.domain.model.RefreshTokenRequest
import com.ricky.controle_financas.domain.model.RefreshTokenResponse
import com.ricky.controle_financas.domain.repository.AuthRepository
import com.ricky.controle_financas.utils.Resource
import javax.inject.Inject

class AuthRepositoryImpl @Inject constructor(
    private val api: AuthApi,
    private val dataStore: DataStoreUtil
) : AuthRepository {
    override suspend fun login(request: LoginRequest): Resource<LoginResponse> {
        val result: Resource<LoginResponse> = safeApiCall {
            api.login(request)
        }

        if (result is Resource.Success) {
            result.data?.let {
                dataStore.saveToken(result.data.accessToken)
                dataStore.saveRefreshToken(result.data.refreshToken)
            }
        }

        return result
    }

    override suspend fun refreshToken(refreshToken: String): Resource<RefreshTokenResponse> {
        val request = RefreshTokenRequest(refreshToken)
        val result: Resource<RefreshTokenResponse> = safeApiCall {
            api.refreshToken(request)
        }

        if (result is Resource.Success) {
            result.data?.let { response ->
                dataStore.saveToken(response.accessToken)
            }
        }

        return result
    }
}