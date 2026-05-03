package com.ricky.controle_financas.data.repository

import com.ricky.controle_financas.domain.model.ErrorResponse
import com.ricky.controle_financas.utils.Resource
import retrofit2.HttpException
import java.io.IOException

inline fun <T> safeApiCall(call: () -> retrofit2.Response<T>): Resource<T> {
    return try {
        val response = call()
        if (response.isSuccessful) {
            response.body()?.let { Resource.Success(it) }
                ?: Resource.Error("Resposta vazia do servidor")
        } else {
            val errorBody = response.errorBody()?.string()
            val errorResponse = ErrorResponse.fromJson(errorBody)
            Resource.Error(errorResponse?.message ?: "Erro desconhecido (${response.code()})")
        }
    } catch (e: HttpException) {
        Resource.Error(message = "Erro de conexão: ${e.localizedMessage}")
    } catch (e: IOException) {
        Resource.Error(message = "Sem internet. Verifique sua conexão.")
    } catch (e: Exception) {
        Resource.Error(message = "Erro inesperado: ${e.localizedMessage}")
    }
}