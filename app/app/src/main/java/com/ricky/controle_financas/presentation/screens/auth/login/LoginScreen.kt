package com.ricky.controle_financas.presentation.screens.auth.login

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountBalanceWallet
import androidx.compose.material.icons.filled.Email
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusRequester
import androidx.compose.ui.focus.focusRequester
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.platform.LocalSoftwareKeyboardController
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.ricky.controle_financas.R
import com.ricky.controle_financas.presentation.components.CustomButton
import com.ricky.controle_financas.presentation.components.ErrorBanner
import com.ricky.controle_financas.presentation.components.PasswordFieldCustom
import com.ricky.controle_financas.presentation.components.TextFieldCustom
import com.ricky.controle_financas.ui.theme.Controle_financasTheme

@Composable
fun LoginScreen(
    modifier: Modifier = Modifier,
    state: LoginState,
    onEvent: (LoginEvent) -> Unit,
) {
    val focusManager = LocalFocusManager.current
    val keyboardController = LocalSoftwareKeyboardController.current
    val scrollState = rememberScrollState()

    val emailFocus = remember { FocusRequester() }
    val passwordFocus = remember { FocusRequester() }

    var passwordVisible by remember { mutableStateOf(false) }

    val isButtonEnabled by remember(state.email, state.password) {
        derivedStateOf {
            state.email.isNotBlank() && state.password.length >= 6 && !state.isLoading
        }
    }

    LaunchedEffect(Unit) {
        emailFocus.requestFocus()
    }

    LaunchedEffect(state.email, state.password) {
        if (state.error != null) {
            onEvent(LoginEvent.ClearError)
        }
    }

    val onImeAction = {
        focusManager.clearFocus()
        keyboardController?.hide()
        if (isButtonEnabled) {
            onEvent(LoginEvent.OnLogin)
        }
    }

    Scaffold(
        modifier = modifier
            .fillMaxSize()
            .imePadding(),
        containerColor = MaterialTheme.colorScheme.background
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = 16.dp)
                .verticalScroll(scrollState),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            LoginHeader()

            Spacer(modifier = Modifier.height(32.dp))

            TextFieldCustom(
                value = state.email,
                onChange = { onEvent(LoginEvent.OnChangeEmail(it)) },
                modifier = Modifier.focusRequester(emailFocus),
                placeholder = "seu@email.com.br",
                label = R.string.email,
                leadingIcon = Icons.Default.Email,
                keyboardType = KeyboardType.Email,
                onNext = { passwordFocus.requestFocus() },
                onDone = { onImeAction() }
            )

            Spacer(modifier = Modifier.height(16.dp))

            PasswordFieldCustom(
                value = state.password,
                onChange = { onEvent(LoginEvent.OnChangeSenha(it)) },
                modifier = Modifier.focusRequester(passwordFocus),
                onDone = { onImeAction() },
                label = R.string.senha
            )

            state.error?.let { error ->
                Spacer(modifier = Modifier.height(12.dp))
                ErrorBanner(
                    message = error,
                    onDismiss = { onEvent(LoginEvent.ClearError) }
                )
            }

            Spacer(modifier = Modifier.height(24.dp))

            CustomButton(
                isLoading = state.isLoading,
                enabled = isButtonEnabled,
                onClick = {
                    focusManager.clearFocus()
                    keyboardController?.hide()
                    onEvent(LoginEvent.OnLogin)
                },
                modifier = Modifier.fillMaxWidth(),
                label = R.string.entrar
            )

            Spacer(modifier = Modifier.height(24.dp))

            AuxiliaryLinks(
                onForgotPassword = {
                    keyboardController?.hide()
                    onEvent(LoginEvent.OnNavigate(LoginNavigation.NavigateToForgetPassword))
                },
                onSignUp = {
                    keyboardController?.hide()
                    onEvent(LoginEvent.OnNavigate(LoginNavigation.NavigateToRegister))
                }
            )
        }
    }
}

@Composable
private fun LoginHeader(modifier: Modifier = Modifier) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Surface(
            modifier = Modifier.size(80.dp),
            shape = MaterialTheme.shapes.medium,
            color = MaterialTheme.colorScheme.primaryContainer
        ) {
            Icon(
                imageVector = Icons.Default.AccountBalanceWallet,
                contentDescription = null,
                tint = MaterialTheme.colorScheme.onPrimaryContainer,
                modifier = Modifier.padding(16.dp)
            )
        }

        Spacer(modifier = Modifier.height(16.dp))

        Text(
            text = stringResource(R.string.bem_vindo),
            style = MaterialTheme.typography.headlineSmall,
            color = MaterialTheme.colorScheme.onBackground,
            textAlign = TextAlign.Center
        )

        Text(
            text = stringResource(R.string.acesse_sua_conta),
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.7f),
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(top = 4.dp)
        )
    }
}


@Composable
private fun AuxiliaryLinks(
    onForgotPassword: () -> Unit,
    onSignUp: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        TextButton(onClick = onForgotPassword) {
            Text(
                text = stringResource(R.string.esqueceu_senha),
                color = MaterialTheme.colorScheme.primary
            )
        }

        Row(
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.Center
        ) {
            Text(
                text = stringResource(R.string.nao_tem_conta),
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.7f)
            )
            TextButton(
                onClick = onSignUp,
                contentPadding = PaddingValues(horizontal = 4.dp)
            ) {
                Text(
                    text = stringResource(R.string.cadastre_se),
                    style = MaterialTheme.typography.labelMedium,
                    color = MaterialTheme.colorScheme.primary
                )
            }
        }
    }
}

@Preview(showBackground = true, name = "Login - Estado Inicial")
@Composable
private fun LoginScreenPreview_Idle() {
    Controle_financasTheme() {
        LoginScreen(
            state = LoginState(),
            onEvent = {}
        )
    }
}

@Preview(showBackground = true, name = "Login - Com Dados")
@Composable
private fun LoginScreenPreview_Filled() {
    Controle_financasTheme {
        LoginScreen(
            state = LoginState(
                email = "joao@empresa.com",
                password = "senha123"
            ),
            onEvent = {}
        )
    }
}

@Preview(showBackground = true, name = "Login - Loading")
@Composable
private fun LoginScreenPreview_Loading() {
    Controle_financasTheme {
        LoginScreen(
            state = LoginState(
                email = "joao@empresa.com",
                password = "senha123",
                isLoading = true
            ),
            onEvent = {}
        )
    }
}

@Preview(showBackground = true, name = "Login - Com Erro")
@Composable
private fun LoginScreenPreview_Error() {
    Controle_financasTheme {
        LoginScreen(
            state = LoginState(
                email = "joao@empresa.com",
                password = "senha123",
                error = "Credenciais inválidas. Tente novamente."
            ),
            onEvent = {}
        )
    }
}

@Preview(showBackground = true, name = "Login - Dark Mode")
@Composable
private fun LoginScreenPreview_Dark() {
    Controle_financasTheme(darkTheme = true) {
        LoginScreen(
            state = LoginState(error = "Sem conexão com a internet"),
            onEvent = {}
        )
    }
}