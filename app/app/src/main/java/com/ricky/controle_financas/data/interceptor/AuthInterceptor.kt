package com.ricky.controle_financas.data.interceptor

import com.ricky.controle_financas.data.auth.SessionStateManager
import com.ricky.controle_financas.data.local.DataStoreUtil
import com.ricky.controle_financas.domain.repository.AuthRepository
import com.ricky.controle_financas.utils.Resource
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject

class AuthInterceptor @Inject constructor(
    private val authRepository: AuthRepository,
    private val dataStoreUtil: DataStoreUtil,
    private val sessionStateManager: SessionStateManager
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()

        if (request.url.encodedPath in PUBLIC_PATHS) {
            return chain.proceed(request)
        }

        val token = runBlocking { dataStoreUtil.getToken() }
        val requestWithAuth = if (token != null) {
            request.newBuilder().addHeader("Authorization", "Bearer $token").build()
        } else {
            request
        }

        var response = chain.proceed(requestWithAuth)
        if (response.code == 401) {
            response.close()

            val refreshResult = runBlocking { authRepository.refreshToken() }

            if (refreshResult is Resource.Success && refreshResult.data != null) {
                val newToken = refreshResult.data.accessToken
                val retryRequest = request.newBuilder()
                    .addHeader("Authorization", "Bearer $newToken")
                    .build()
                response = chain.proceed(retryRequest)
            } else {
                runBlocking {
                    dataStoreUtil.clearTokens()
                    sessionStateManager.notifySessionExpired()
                }
            }
        }
        return response
    }

    companion object {
        private val PUBLIC_PATHS = setOf("/auth/login", "/auth/register", "/auth/refresh")
    }
}