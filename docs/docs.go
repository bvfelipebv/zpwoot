package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/sessions/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Cria uma nova sessão do WhatsApp com nome e webhook opcional",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Criar nova sessão",
                "parameters": [
                    {
                        "description": "Dados da sessão",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateSessionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.SessionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retorna a lista de todas as sessões criadas",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Listar sessões",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SessionListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/connect": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Inicia a conexão de uma sessão com o WhatsApp",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Conectar sessão",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Remove uma sessão e todos os seus dados associados",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Deletar sessão",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/disconnect": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Desconecta uma sessão ativa do WhatsApp",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Desconectar sessão",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/info": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retorna informações detalhadas de uma sessão específica",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Obter detalhes da sessão",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SessionResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/pair": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Gera um código de pareamento para conectar o WhatsApp usando número de telefone",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Parear com telefone",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Número de telefone",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.PairPhoneRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.PairPhoneResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/status": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retorna informações detalhadas sobre o status de conexão da sessão",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Obter status da sessão",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SessionStatusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{id}/webhook": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Atualiza a URL e eventos do webhook de uma sessão",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Atualizar webhook",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Configurações de webhook",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdateWebhookRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CreateSessionRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "apikey": {
                    "type": "string",
                    "example": "null"
                },
                "name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 3,
                    "example": "sessao-atendimento-1"
                },
                "proxy": {
                    "$ref": "#/definitions/dto.ProxyConfig"
                },
                "webhook": {
                    "$ref": "#/definitions/dto.WebhookConfig"
                }
            }
        },
        "dto.ErrorResponse": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "object"
                },
                "error": {
                    "type": "string",
                    "example": "invalid_request"
                },
                "message": {
                    "type": "string",
                    "example": "Nome da sessão é obrigatório"
                }
            }
        },
        "dto.PairPhoneRequest": {
            "type": "object",
            "required": [
                "phone_number"
            ],
            "properties": {
                "phone_number": {
                    "description": "Formato: +5511999999999",
                    "type": "string",
                    "example": "+5511999999999"
                }
            }
        },
        "dto.PairPhoneResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Enter the pairing code on your phone"
                },
                "pairing_code": {
                    "type": "string",
                    "example": "ABCD-1234"
                },
                "phone_number": {
                    "type": "string",
                    "example": "+5511999999999"
                },
                "session_id": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                }
            }
        },
        "dto.ProxyConfig": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean",
                    "example": true
                },
                "host": {
                    "type": "string",
                    "example": "10.0.0.1"
                },
                "password": {
                    "type": "string",
                    "example": "proxypass"
                },
                "port": {
                    "type": "integer",
                    "maximum": 65535,
                    "minimum": 1,
                    "example": 3128
                },
                "protocol": {
                    "type": "string",
                    "enum": [
                        "http",
                        "https",
                        "socks5"
                    ],
                    "example": "http"
                },
                "username": {
                    "type": "string",
                    "example": "proxyuser"
                }
            }
        },
        "dto.SessionListResponse": {
            "type": "object",
            "properties": {
                "sessions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.SessionResponse"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 3
                }
            }
        },
        "dto.SessionResponse": {
            "type": "object",
            "properties": {
                "business_name": {
                    "type": "string",
                    "example": "Minha Empresa LTDA"
                },
                "created_at": {
                    "type": "string",
                    "example": "2025-11-05T10:00:00Z"
                },
                "id": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "jid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "last_connected": {
                    "type": "string",
                    "example": "2025-11-05T18:30:00Z"
                },
                "last_disconnected": {
                    "type": "string",
                    "example": "2025-11-05T17:00:00Z"
                },
                "name": {
                    "type": "string",
                    "example": "Minha Sessão WhatsApp"
                },
                "platform": {
                    "type": "string",
                    "example": "android"
                },
                "push_name": {
                    "type": "string",
                    "example": "João Silva"
                },
                "status": {
                    "type": "string",
                    "example": "connected"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2025-11-05T18:30:00Z"
                },
                "webhook_events": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "message",
                        "qr",
                        "connected"
                    ]
                },
                "webhook_url": {
                    "type": "string",
                    "example": "https://seu-webhook.com/whatsapp"
                }
            }
        },
        "dto.SessionStatusResponse": {
            "type": "object",
            "properties": {
                "can_connect": {
                    "type": "boolean",
                    "example": true
                },
                "connection_time": {
                    "description": "Duração formatada",
                    "type": "string",
                    "example": "2h 30m 15s"
                },
                "id": {
                    "type": "string",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "is_connected": {
                    "type": "boolean",
                    "example": true
                },
                "is_logged_in": {
                    "type": "boolean",
                    "example": true
                },
                "jid": {
                    "type": "string",
                    "example": "5511999999999@s.whatsapp.net"
                },
                "last_connected": {
                    "type": "string",
                    "example": "2025-11-05T18:30:00Z"
                },
                "needs_pairing": {
                    "type": "boolean",
                    "example": false
                },
                "platform": {
                    "type": "string",
                    "example": "android"
                },
                "push_name": {
                    "type": "string",
                    "example": "João Silva"
                },
                "status": {
                    "type": "string",
                    "example": "connected"
                }
            }
        },
        "dto.SuccessResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string",
                    "example": "Operação realizada com sucesso"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "dto.UpdateWebhookRequest": {
            "type": "object",
            "required": [
                "webhook"
            ],
            "properties": {
                "webhook": {
                    "$ref": "#/definitions/dto.WebhookConfig"
                }
            }
        },
        "dto.WebhookConfig": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean",
                    "example": true
                },
                "events": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "message",
                        "status",
                        "qr"
                    ]
                },
                "token": {
                    "type": "string",
                    "example": "secreto-opcional"
                },
                "url": {
                    "type": "string",
                    "example": "https://hooks.exemplo.com/wuz"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Insira sua API Key (exemplo: sldkfjsldkflskdfjlsd)",
            "type": "apiKey",
            "name": "apikey",
            "in": "header"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "ZPWoot - WhatsApp Multi-Device API",
	Description:      "API REST para gerenciamento de múltiplas sessões do WhatsApp usando whatsmeow\nPermite criar, conectar e gerenciar sessões do WhatsApp via API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
