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
        "/auth/changePassword": {
            "post": {
                "description": "Users are required to change password on first time login, using their username and password sent via email.\nSubmit The user's username and new password to respond to this auth challenge.\nRequest must contain \"session\" cookie containing the session token to respond to the challenge\nOn success, responds with next auth challenge, which should be to set up MFA",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Change Password for first-time Login",
                "parameters": [
                    {
                        "type": "string",
                        "example": "ABCDEF",
                        "name": "newPassword",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "joel.ow.2022",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AuthChallengeRes"
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
                    }
                }
            }
        },
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
                        "type": "string",
                        "example": "ABCDEF",
                        "name": "code",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "67890",
                        "name": "newPassword",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "joel.ow.2022",
                        "name": "username",
                        "in": "formData",
                        "required": true
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
                "description": "Admin registers user with Cognito user pool via email. Cognito sends an email with a temporary password to the user.",
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
                        "type": "string",
                        "description": "User's email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
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
                        "type": "string",
                        "example": "joel.ow.2022",
                        "name": "username",
                        "in": "formData"
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
                "description": "Cognito SSO login using username and password, returns the next auth challenge, either",
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
                        "type": "string",
                        "example": "12345",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "joel.ow.2022",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AuthChallengeRes"
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
        "/auth/loginMFA": {
            "post": {
                "description": "Responds to Cognito auth challenge after successful credential sign in\nRequest must contain \"session\" cookie containing the session token to respond to the challenge",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Submit user TOTP code from authenticator app for all subsequent log ins.",
                "parameters": [
                    {
                        "type": "string",
                        "name": "code",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "username",
                        "in": "formData",
                        "required": true
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
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Clears the session by expiring the cookies containing the JWT tokens",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout User",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.StatusRes"
                        }
                    }
                }
            }
        },
        "/auth/setupMFA": {
            "get": {
                "description": "Submit GET query to cognito to obtain an OTP token.\nThe user can use this token to set up their authenticator app, either through QR code or by manual keying in of the token.\nRequest must contain \"session\" cookie containing the session token to respond to the challenge\nOn success, the token is returned, and the cookie is updated for the next auth step",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Get OTP Token for setting up TOTP authenticator",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SetupMFARes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
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
        "/auth/verifyMFA": {
            "post": {
                "description": "User submits the code from their authenticator app to verify the TOTP setup\nRequest must contain \"session\" cookie containing the session token to respond to the challenge\nOn success, the user can proceed to sign in again",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify initial code from authenticator app",
                "parameters": [
                    {
                        "type": "string",
                        "name": "code",
                        "in": "formData",
                        "required": true
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
                    "500": {
                        "description": "Internal Server Error",
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
        "models.AuthChallengeRes": {
            "type": "object",
            "properties": {
                "challenge": {
                    "type": "string",
                    "enum": [
                        "NEW_PASSWORD_REQUIRED",
                        "MFA_SETUP",
                        "SOFTWARE_TOKEN_MFA"
                    ],
                    "example": "SOFTWARE_TOKEN_MFA"
                }
            }
        },
        "models.SetupMFARes": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
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
