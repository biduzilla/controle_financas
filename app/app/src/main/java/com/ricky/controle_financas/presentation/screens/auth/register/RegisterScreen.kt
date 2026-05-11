package com.ricky.controle_financas.presentation.screens.auth.register

import androidx.activity.compose.BackHandler
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountBalanceWallet
import androidx.compose.material.icons.filled.Email
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Phone
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
import androidx.compose.runtime.remember
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
fun RegisterScreen(
    modifier: Modifier = Modifier,
    state: RegisterState,
    onEvent: (RegisterEvent) -> Unit,
    onBack: () -> Unit
) {
    BackHandler(onBack = onBack)

    val focusManager = LocalFocusManager.current
    val keyboardController = LocalSoftwareKeyboardController.current
    val scrollState = rememberScrollState()

    val nameFocus = remember { FocusRequester() }
    val emailFocus = remember { FocusRequester() }
    val phoneFocus = remember { FocusRequester() }
    val passwordFocus = remember { FocusRequester() }
    val confirmPasswordFocus = remember { FocusRequester() }

    val isButtonEnabled by remember(
        state.name,
        state.email,
        state.phone,
        state.password,
        state.confirmPassword
    ) {
        derivedStateOf {
            state.name.isNotBlank() && state.email.isNotBlank() &&
                    state.phone.isNotBlank() &&
                    state.password.length >= 6 && state.password == state.confirmPassword &&
                    !state.isLoading
        }
    }

    LaunchedEffect(Unit) {
        nameFocus.requestFocus()
    }

    LaunchedEffect(
        state.name,
        state.email,
        state.phone,
        state.password,
        state.confirmPassword,
    ) {
        if (state.error != null) {
            onEvent(RegisterEvent.ClearError)
        }
    }

    val onImeAction: () -> Unit = {
        focusManager.clearFocus()
        keyboardController?.hide()
        if (isButtonEnabled) {
            onEvent(RegisterEvent.OnRegister)
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
            RegisterHeader()

            TextFieldCustom(
                value = state.name,
                onChange = { onEvent(RegisterEvent.OnChangeName(it)) },
                label = R.string.nome,
                placeholder = "seu nome",
                leadingIcon = Icons.Default.Person,
                onNext = { emailFocus.requestFocus() },
                modifier = Modifier.focusRequester(nameFocus)
            )

            Spacer(modifier = Modifier.height(16.dp))

            TextFieldCustom(
                value = state.email,
                onChange = { onEvent(RegisterEvent.OnChangeEmail(it)) },
                label = R.string.email,
                placeholder = "seu@email.com.br",
                leadingIcon = Icons.Default.Email,
                onNext = { phoneFocus.requestFocus() },
                modifier = Modifier.focusRequester(emailFocus),
                keyboardType = KeyboardType.Email
            )

            Spacer(modifier = Modifier.height(16.dp))

            TextFieldCustom(
                value = state.phone,
                onChange = { onEvent(RegisterEvent.OnChangeEmail(it)) },
                label = R.string.phone,
                placeholder = "99999999999",
                leadingIcon = Icons.Default.Phone,
                onNext = { passwordFocus.requestFocus() },
                modifier = Modifier.focusRequester(phoneFocus),
                keyboardType = KeyboardType.Phone
            )

            Spacer(modifier = Modifier.height(16.dp))

            PasswordFieldCustom(
                modifier = Modifier.focusRequester(passwordFocus),
                value = state.password,
                label = R.string.senha,
                onNext = { confirmPasswordFocus.requestFocus() },
            )

            Spacer(modifier = Modifier.height(16.dp))

            PasswordFieldCustom(
                modifier = Modifier.focusRequester(confirmPasswordFocus),
                value = state.password,
                label = R.string.senha,
                onDone = { onImeAction() },
            )

            state.error?.let { error ->
                Spacer(modifier = Modifier.height(16.dp))
                ErrorBanner(
                    message = error,
                    onDismiss = { onEvent(RegisterEvent.ClearError) }
                )
            }

            Spacer(modifier = Modifier.height(24.dp))

            CustomButton(
                enabled = isButtonEnabled,
                label = R.string.cadastrar,
                isLoading = state.isLoading
            ) {
                focusManager.clearFocus()
                keyboardController?.hide()
                onEvent(RegisterEvent.OnRegister)
            }

            Spacer(modifier = Modifier.height(24.dp))

            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.Center
            ) {
                Text(
                    text = stringResource(R.string.ja_tem_conta),
                    style = MaterialTheme.typography.bodyMedium,
                    color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.7f)
                )
                TextButton(
                    onClick = { onEvent(RegisterEvent.OnNavigate(RegisterNavigation.NavigateToLogin)) },
                    contentPadding = PaddingValues(horizontal = 4.dp)
                ) {
                    Text(
                        text = stringResource(R.string.fazer_login),
                        style = MaterialTheme.typography.labelMedium,
                        color = MaterialTheme.colorScheme.primary
                    )
                }
            }
        }
    }
}

@Composable
private fun RegisterHeader(modifier: Modifier = Modifier) {
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
            text = stringResource(R.string.crie_sua_conta),
            style = MaterialTheme.typography.headlineSmall,
            color = MaterialTheme.colorScheme.onBackground,
            textAlign = TextAlign.Center
        )

        Text(
            text = stringResource(R.string.preencha_dados_cadastro),
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onBackground.copy(alpha = 0.7f),
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(top = 4.dp)
        )
    }
}

@Preview(showBackground = true, name = "Register - Estado Inicial")
@Composable
private fun RegisterScreenPreview_Idle() {
    Controle_financasTheme {
        RegisterScreen(state = RegisterState(), onEvent = {})
    }
}

@Preview(showBackground = true, name = "Register - Preenchido")
@Composable
private fun RegisterScreenPreview_Filled() {
    Controle_financasTheme {
        RegisterScreen(
            state = RegisterState(
                name = "Maria Silva",
                email = "maria@email.com",
                phone = "9999999999999",
                password = "123456",
                confirmPassword = "123456"
            ),
            onEvent = {}
        )
    }
}

@Preview(showBackground = true, name = "Register - Loading")
@Composable
private fun RegisterScreenPreview_Loading() {
    Controle_financasTheme {
        RegisterScreen(
            state = RegisterState(
                name = "Maria Silva",
                email = "maria@email.com",
                password = "123456",
                phone = "9999999999999",
                confirmPassword = "123456",
                isLoading = true
            ),
            onEvent = {}
        )
    }
}

@Preview(showBackground = true, name = "Register - Erro")
@Composable
private fun RegisterScreenPreview_Error() {
    Controle_financasTheme {
        RegisterScreen(
            state = RegisterState(
                name = "Maria Silva",
                email = "maria@email.com",
                phone = "9999999999",
                password = "123456",
                confirmPassword = "1234567",
                error = "Senhas não conferem"
            ),
            onEvent = {}
        )
    }
}