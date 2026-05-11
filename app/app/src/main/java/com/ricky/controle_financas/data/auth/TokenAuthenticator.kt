package com.ricky.controle_financas.data.auth

import com.ricky.controle_financas.data.api.AuthApi
import com.ricky.controle_financas.data.local.DataStoreUtil
import com.ricky.controle_financas.domain.model.RefreshTokenRequest
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import okhttp3.Authenticator
import okhttp3.Request
import okhttp3.Response
import okhttp3.Route
import javax.inject.Inject
import javax.inject.Provider
import javax.inject.Singleton

@Singleton
class TokenAuthenticator @Inject constructor(
    private val authApiProvider: Provider<AuthApi>,
    private val dataStore: DataStoreUtil,
    private val sessionStateManager: SessionStateManager
) : Authenticator {

    private val refreshMutex = Mutex()

    override fun authenticate(route: Route?, response: Response): Request? {
        if (responseCount(response) >= 3) return null

        val refreshToken = runBlocking { dataStore.getRefreshToken() } ?: return null

        return runBlocking {
            refreshMutex.withLock {
                val currentToken = dataStore.getToken()

                if (currentToken != null && currentToken != response.request.header("Authorization")
                        ?.removePrefix("Bearer ")
                ) {
                    response.request.newBuilder()
                        .header("Authorization", "Bearer $currentToken")
                        .build()
                } else {
                    val authApi = authApiProvider.get()
                    val refreshResponse = authApi.refreshToken(RefreshTokenRequest(refreshToken))

                    if (refreshResponse.isSuccessful) {
                        val newToken = refreshResponse.body()?.accessToken
                        if (!newToken.isNullOrBlank()) {
                            dataStore.saveToken(newToken)
                            response.request.newBuilder()
                                .header("Authorization", "Bearer $newToken")
                                .build()
                        } else null
                    } else {
                        dataStore.clearTokens()
                        sessionStateManager.notifySessionExpired()
                        null
                    }
                }
            }

        }
    }

    private fun responseCount(response: Response): Int =
        response.priorResponse?.let { responseCount(it) + 1 } ?: 1
}
