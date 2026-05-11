package com.ricky.controle_financas.domain.use_case

import com.ricky.controle_financas.domain.model.RefreshTokenResponse
import com.ricky.controle_financas.domain.model.User
import com.ricky.controle_financas.domain.repository.AuthRepository
import com.ricky.controle_financas.domain.repository.UserRepository
import com.ricky.controle_financas.utils.Resource
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject

class SaveUserUseCase @Inject constructor(
    private val repository: UserRepository
) {
    operator fun invoke(user: User): Flow<Resource<User>> = flow {
        emit(Resource.Loading())
        val result = repository.createUser(user)
        emit(result)
    }
}