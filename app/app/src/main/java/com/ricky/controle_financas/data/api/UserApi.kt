package com.ricky.controle_financas.data.api

import com.ricky.controle_financas.domain.model.User
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

interface UserApi {
    @POST("user")
    suspend fun createUser(@Body request: User): Response<User>

    @GET("user/{id}")
    suspend fun getUser(@Path("id") id: String): Response<User>
}