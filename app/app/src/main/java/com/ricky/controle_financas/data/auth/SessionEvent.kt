package com.ricky.controle_financas.data.auth

sealed class SessionEvent {
    object SessionExpired : SessionEvent()
}