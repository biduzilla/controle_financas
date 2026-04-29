package com.ricky.controle_financas.data.api

import com.ricky.controle_financas.domain.model.LoginRequest
import com.ricky.controle_financas.domain.model.LoginResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface AuthApi {
    @POST("v1/auth")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>
}