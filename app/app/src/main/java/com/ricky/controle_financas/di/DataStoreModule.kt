package com.ricky.controle_financas.di

import android.content.Context
import com.ricky.controle_financas.data.local.DataStoreUtil
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DataStoreModule {
    @Provides
    @Singleton
    fun provideDataStoreUtil(
        @ApplicationContext context: Context,
    ): DataStoreUtil {
        return DataStoreUtil(context)
    }
}