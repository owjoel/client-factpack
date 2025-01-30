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
        "/auth/confirmForgetPassword": {
            "post": {
                "description": "Submit Cognito OTP sent to user's email to proceed with password reset",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Confirm Forget Password",
                "parameters": [
                    {
                        "description": "OTP Code",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ConfirmForgetPasswordReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    }
                }
            }
        },
        "/auth/createUser": {
            "post": {
                "description": "Admin registers user with Cognito user pool via email",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Create Users",
                "parameters": [
                    {
                        "description": "User email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignUpReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    }
                }
            }
        },
        "/auth/forgetPassword": {
            "post": {
                "description": "Forget password",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Forget Password",
                "parameters": [
                    {
                        "description": "Username",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ForgetPasswordReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Cognito SSO login using username and password",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Username, Password",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
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
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ConfirmForgetPasswordReq": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "ABCDEF"
                },
                "newPassword": {
                    "type": "string",
                    "example": "67890"
                },
                "username": {
                    "type": "string",
                    "example": "joel.ow.2022"
                }
            }
        },
        "models.ForgetPasswordReq": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string",
                    "example": "joel.ow.2022"
                }
            }
        },
        "models.LoginReq": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "12345"
                },
                "username": {
                    "type": "string",
                    "example": "joel.ow.2022"
                }
            }
        },
        "models.SignUpReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@gmail.com"
                }
            }
        },
        "models.StatusRes": {
            "type": "object",
            "properties": {
                "status": {
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
	Title:            "client-factpack/auth",
	Description:      "Authentication service for managing auth flows",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
