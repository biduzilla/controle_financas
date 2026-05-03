package com.ricky.controle_financas.data.local

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import com.google.gson.Gson
import com.ricky.controle_financas.domain.model.User
import com.ricky.controle_financas.utils.IS_DARK_MODE
import com.ricky.controle_financas.utils.KEY_REFRESH_TOKEN
import com.ricky.controle_financas.utils.KEY_TOKEN
import com.ricky.controle_financas.utils.KEY_USER
import com.ricky.controle_financas.utils.SETTINGS
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map


class DataStoreUtil(private val context: Context) {

    companion object {
        private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(SETTINGS)

        val THEME_KEY = booleanPreferencesKey(IS_DARK_MODE)
        val TOKEN = stringPreferencesKey(KEY_TOKEN)
        val REFRESH_TOKEN = stringPreferencesKey(KEY_REFRESH_TOKEN)
        val USER = stringPreferencesKey(KEY_USER)
    }

    suspend fun saveTheme(isDark: Boolean) {
        context.dataStore.edit { preferences ->
            preferences[THEME_KEY] = isDark
        }
    }

    suspend fun saveToken(token: String) {
        context.dataStore.edit { p -> p[TOKEN] = token }
    }

    suspend fun saveRefreshToken(refreshToken: String) {
        context.dataStore.edit { p -> p[TOKEN] = refreshToken }
    }

    suspend fun saveUser(user: User) {
        val json = Gson().toJson(user)
        context.dataStore.edit { p -> p[USER] = json }
    }

    suspend fun getToken(): String? {
        return context.dataStore.data.first()[TOKEN]
    }

    suspend fun getRefreshToken(): String? {
        return context.dataStore.data.first()[REFRESH_TOKEN]
    }

    fun getTokenAsFlow(): Flow<String?> {
        return context.dataStore.data.map { it[TOKEN] }
    }

    fun getRefreshTokenAsFlow(): Flow<String?> {
        return context.dataStore.data.map { it[REFRESH_TOKEN] }
    }

    fun getTheme(): Flow<Boolean> {
        return context.dataStore.data.map { preferences ->
            preferences[THEME_KEY] ?: false
        }
    }

    fun getUser(): Flow<User?> {
        return context.dataStore.data.map { preferences ->
            val json = preferences[USER] ?: ""
            if (json.isNotEmpty()) {
                Gson().fromJson(json, User::class.java)
            } else {
                null
            }
        }
    }
}