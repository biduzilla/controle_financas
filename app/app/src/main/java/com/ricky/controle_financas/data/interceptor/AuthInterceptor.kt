package com.ricky.controle_financas.data.interceptor

import com.ricky.controle_financas.data.local.DataStoreUtil
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Request
import okhttp3.Response
import javax.inject.Inject

class AuthInterceptor @Inject constructor(
    private val dataStoreUtil: DataStoreUtil,
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        if (isPublicPath(request)) {
            return chain.proceed(request)
        }

        val token = runBlocking { dataStoreUtil.getToken() }
        val newRequest = if (!token.isNullOrBlank()) {
            request.newBuilder().addHeader("Authorization", "Bearer $token").build()
        } else request

        return chain.proceed(newRequest)
    }

    companion object {
        private val PUBLIC_PATHS = mapOf(
            "/v1/auth/login"    to listOf("POST"),
            "/v1/auth/register" to listOf("POST"),
            "/v1/auth/refresh"  to listOf("POST"),
            "/v1/user"          to listOf("POST")
        )
    }

    private fun isPublicPath(request: Request): Boolean {
        val allowedMethods = PUBLIC_PATHS[request.url.encodedPath] ?: return false
        return request.method in allowedMethods
    }
}