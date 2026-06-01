package com.ricky.controle_financas.presentation.components

import android.content.res.Configuration
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TopBar(
    modifier: Modifier = Modifier,
    userName: String,
    userIcon: @Composable () -> Unit = {
        Icon(
            imageVector = Icons.Default.AccountCircle,
            contentDescription = null,
            modifier = Modifier.size(24.dp)
        )
    },
    onUserClick: (() -> Unit)? = null,
    rightIcon: @Composable () -> Unit = {},
    rightIconContentDescription: String = "",
    onRightIconClick: (() -> Unit)? = null,
    actions: @Composable RowScope.() -> Unit = {},
) {
    TopAppBar(
        modifier = modifier.fillMaxWidth(),
        title = {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                modifier = Modifier.clickable(
                    enabled = onUserClick != null,
                    onClick = { onUserClick?.invoke() }
                )
            ) {
                userIcon()
                Text(
                    text = userName,
                    style = MaterialTheme.typography.titleMedium,
                    maxLines = 1,
                    overflow = TextOverflow.Ellipsis
                )
            }
        },
        actions = {
            if (onRightIconClick != null) {
                IconButton(
                    onClick = onRightIconClick,
                    modifier = Modifier.semantics { this.contentDescription = rightIconContentDescription }
                ) {
                    rightIcon()
                }
            }
            actions()
        }
    )
}

@Preview(showBackground = true, name = "Light - Padrão")
@Preview(showBackground = true, name = "Dark - Padrão", uiMode = Configuration.UI_MODE_NIGHT_YES)
@Composable
private fun TopBarPreview() {
    MaterialTheme {
        TopBar(
            userName = "Ricky Silva",
            onUserClick = {},
            rightIcon = { Icon(Icons.Default.Settings, null) },
            rightIconContentDescription = "Configurações",
            onRightIconClick = {}
        )
    }
}

@Preview(name = "Nome Longo + Sem Ação Direita")
@Composable
private fun TopBarLongNamePreview() {
    MaterialTheme {
        TopBar(
            userName = "Maria Joana de Oliveira e Silva Santos",
            onUserClick = {},
            rightIconContentDescription = ""
        )
    }
}