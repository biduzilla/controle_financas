package com.ricky.controle_financas.domain.use_case

import com.ricky.controle_financas.domain.model.RefreshTokenResponse
import com.ricky.controle_financas.domain.repository.AuthRepository
import com.ricky.controle_financas.utils.Resource
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject

class RefreshTokenUserCase @Inject constructor(
    private val repository: AuthRepository
) {
    operator fun invoke(refreshToken: String): Flow<Resource<RefreshTokenResponse>> = flow {
        emit(Resource.Loading())
        val result = repository.refreshToken(refreshToken)
        emit(result)
    }
}