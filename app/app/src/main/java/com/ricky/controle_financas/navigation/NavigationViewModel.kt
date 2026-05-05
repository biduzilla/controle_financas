package com.ricky.controle_financas.navigation

import androidx.compose.runtime.mutableStateListOf
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.ricky.controle_financas.data.auth.SessionEvent
import com.ricky.controle_financas.data.auth.SessionStateManager
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class NavigationViewModel @Inject constructor(
    private val sessionStateManager: SessionStateManager
) : ViewModel() {

    init {
        viewModelScope.launch {
            sessionStateManager.events.collect { event ->
                when (event) {
                    SessionEvent.SessionExpired -> {
                        navigateAndClear(AppNavKey.Login)
                    }
                }
            }
        }
    }

    val backStack = mutableStateListOf<AppNavKey>(AppNavKey.Splash)

    fun navigate(key: AppNavKey) {
        backStack.add(key)
    }

    fun navigateAndClear(key: AppNavKey) {
        backStack.clear()
        backStack.add(key)
    }

    fun popBackStack() {
        if (backStack.size > 1) {
            backStack.removeAt(backStack.lastIndex)
        }
    }

    fun popToRoot() {
        if (backStack.size > 1) {
            backStack.subList(1, backStack.size).clear()
        }
    }

    fun replace(key: AppNavKey) {
        if (backStack.isNotEmpty()) {
            backStack[backStack.lastIndex] = key
        }
    }
}