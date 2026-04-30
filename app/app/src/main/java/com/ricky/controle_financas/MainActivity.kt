package com.ricky.controle_financas

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.ui.Modifier
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import com.ricky.controle_financas.navigation.AppNavDisplay
import com.ricky.controle_financas.navigation.NavigationViewModel
import com.ricky.controle_financas.ui.theme.Controle_financasTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            Controle_financasTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    val viewModel = hiltViewModel<NavigationViewModel>()
                    AppNavDisplay(
                        modifier = Modifier.padding(innerPadding),
                        navigationViewModel = viewModel
                    )
                }
            }
        }
    }
}
