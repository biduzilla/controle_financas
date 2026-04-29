package com.ricky.controle_financas.data.repository

import android.util.Log
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
            Log.i("infoteste", "safeApiCall: $errorBody")
            Resource.Error(errorBody ?: "")
        }
    } catch (e: HttpException) {
        Resource.Error(message = "Erro de conexão: ${e.localizedMessage}")
    } catch (e: IOException) {
        Resource.Error(message = "Sem internet. Verifique sua conexão.")
    } catch (e: Exception) {
        Resource.Error(message = "Erro inesperado: ${e.localizedMessage}")
    }
}