package com.ricky.controle_financas.navigation

import androidx.compose.runtime.mutableStateListOf
import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject

@HiltViewModel
class NavigationViewModel @Inject constructor() : ViewModel() {
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