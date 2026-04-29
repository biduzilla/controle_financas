package com.ricky.controle_financas.di

import com.ricky.controle_financas.data.repository.AuthRepositoryImpl
import com.ricky.controle_financas.data.repository.UserRepositoryImpl
import com.ricky.controle_financas.domain.repository.AuthRepository
import com.ricky.controle_financas.domain.repository.UserRepository
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {
    @Binds
    @Singleton
    abstract fun bindUserRepository(impl: UserRepositoryImpl): UserRepository

    @Binds
    @Singleton
    abstract fun bindAuthRepository(impl: AuthRepositoryImpl): AuthRepository
}