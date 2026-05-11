package com.ricky.controle_financas.data.interceptor

import com.ricky.controle_financas.data.auth.SessionStateManager
import com.ricky.controle_financas.data.local.DataStoreUtil
import com.ricky.controle_financas.domain.use_case.RefreshTokenUserCase
import com.ricky.controle_financas.utils.Resource
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject

class AuthInterceptor @Inject constructor(
    private val dataStoreUtil: DataStoreUtil,
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        if (request.url.encodedPath in PUBLIC_PATHS) {
            return chain.proceed(request)
        }

        val token = runBlocking { dataStoreUtil.getToken() }
        val newRequest = if (!token.isNullOrBlank()) {
            request.newBuilder().addHeader("Authorization", "Bearer $token").build()
        } else request

        return chain.proceed(newRequest)
    }

    companion object {
        private val PUBLIC_PATHS = setOf("/auth/login", "/auth/register", "/auth/refresh")
    }
}