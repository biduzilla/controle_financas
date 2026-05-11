package com.ricky.controle_financas.utils

import kotlinx.serialization.KSerializer
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.descriptors.SerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import kotlinx.serialization.json.*

object FlexibleMessageSerializer : KSerializer<String> {
    override val descriptor: SerialDescriptor =
        PrimitiveSerialDescriptor("FlexibleMessage", PrimitiveKind.STRING)

    override fun deserialize(decoder: Decoder): String {
        val jsonDecoder = decoder as JsonDecoder
        return when (val element = jsonDecoder.decodeJsonElement()) {
            is JsonPrimitive -> element.content
            is JsonObject -> {
                element.entries.joinToString("\n") { (key, value) ->
                    "$key: ${value.jsonPrimitive.content}"
                }
            }
            else -> ""
        }
    }

    override fun serialize(encoder: Encoder, value: String) {
        encoder.encodeString(value)
    }
}