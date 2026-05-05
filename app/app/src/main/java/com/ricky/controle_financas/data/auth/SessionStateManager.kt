package com.ricky.controle_financas.data.auth

import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.asSharedFlow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class SessionStateManager @Inject constructor() {
    private val _events = MutableSharedFlow<SessionEvent>(extraBufferCapacity = 1)
    val events = _events.asSharedFlow()

    suspend fun emitEvent(event: SessionEvent) {
        _events.emit(event)
    }

    suspend fun notifySessionExpired() = emitEvent(SessionEvent.SessionExpired)
}