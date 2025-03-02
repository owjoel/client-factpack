// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/createProfile": {
            "post": {
                "description": "Create new client profile, given the populated json",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "clients"
                ],
                "summary": "Create Clients",
                "parameters": [
                    {
                        "description": "Client data",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Client"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Basic health check",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "ping",
                "responses": {
                    "200": {
                        "description": "Connection status",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            }
        },
        "/retrieveAllProfiles": {
            "get": {
                "description": "Retrieve all client data",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "clients"
                ],
                "summary": "Get All Clients",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handlers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.Client"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            }
        },
        "/retrieveProfile": {
            "get": {
                "description": "Retrieve client data by profile id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "clients"
                ],
                "summary": "Get Client By ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hex id used to identify client",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handlers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.Client"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "status": {
                    "type": "integer"
                },
                "timestamp": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "model.Associate": {
            "type": "object",
            "properties": {
                "associatedCompanies": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "relationship": {
                    "type": "string"
                }
            }
        },
        "model.Client": {
            "type": "object",
            "properties": {
                "associates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Associate"
                    }
                },
                "id": {
                    "description": "gorm.Model",
                    "type": "string",
                    "example": "6e9938xhdfv27bhspbf73jks"
                },
                "investments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Investment"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/model.Metadata"
                },
                "profile": {
                    "$ref": "#/definitions/model.Profile"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "model.Contact": {
            "type": "object",
            "properties": {
                "phone": {
                    "type": "string"
                },
                "workAddress": {
                    "type": "string"
                }
            }
        },
        "model.Investment": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "industry": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "value": {
                    "$ref": "#/definitions/model.InvestmentValue"
                }
            }
        },
        "model.InvestmentValue": {
            "type": "object",
            "properties": {
                "currency": {
                    "type": "string"
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "model.Metadata": {
            "type": "object",
            "properties": {
                "sources": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "model.NetWorth": {
            "type": "object",
            "properties": {
                "currency": {
                    "type": "string"
                },
                "estiamtedValue": {
                    "type": "integer"
                },
                "source": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "model.Profile": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer",
                    "example": 55
                },
                "contact": {
                    "$ref": "#/definitions/model.Contact"
                },
                "currentResidence": {
                    "$ref": "#/definitions/model.Residence"
                },
                "industries": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string",
                    "example": "john doe"
                },
                "nationality": {
                    "type": "string",
                    "example": "chinese"
                },
                "netWorth": {
                    "$ref": "#/definitions/model.NetWorth"
                },
                "occupations": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "socials": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.SocialMedia"
                    }
                }
            }
        },
        "model.Residence": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                }
            }
        },
        "model.SocialMedia": {
            "type": "object",
            "properties": {
                "platform": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "client-factpack/clients",
	Description:      "Client resource module. Manages manually typed and compiled online data of prospective clients",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
