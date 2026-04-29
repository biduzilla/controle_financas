package com.ricky.controle_financas.domain.repository

import com.ricky.controle_financas.domain.model.User
import com.ricky.controle_financas.utils.Resource

interface UserRepository {
    suspend fun createUser(user: User): Resource<User>
}