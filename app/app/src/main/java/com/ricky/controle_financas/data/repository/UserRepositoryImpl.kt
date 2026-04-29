package com.ricky.controle_financas.data.repository

import com.ricky.controle_financas.data.api.UserApi
import com.ricky.controle_financas.domain.model.User
import com.ricky.controle_financas.domain.repository.UserRepository
import com.ricky.controle_financas.utils.Resource
import javax.inject.Inject

class UserRepositoryImpl @Inject constructor(
    private val api: UserApi
) : UserRepository {
    override suspend fun createUser(user: User): Resource<User> {
        return safeApiCall {
            api.createUser(user)
        }
    }
}