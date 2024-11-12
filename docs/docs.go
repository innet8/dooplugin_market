// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "xxyijixx",
            "email": "xxyijixx@gmail.com"
        },
        "license": {
            "name": "AGPL-3.0",
            "url": "https://opensource.org/license/agpl-v3"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/apps": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "获取插件列表",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "page_size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "class",
                        "name": "class",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/dto.PageResult"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "items": {
                                                            "type": "array",
                                                            "items": {
                                                                "$ref": "#/definitions/model.App"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/apps/installed": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "获取已安装插件列表",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "page_size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "分类",
                        "name": "class",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/dto.PageResult"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "items": {
                                                            "type": "array",
                                                            "items": {
                                                                "type": "object"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/apps/installed/{id}/logs": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "获取插件日志信息",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "description": "开始时间(Unix时间戳)",
                        "name": "since",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "结束时间(Unix时间戳)",
                        "name": "until",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1000,
                        "description": "查询条数",
                        "name": "tail",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/dto.Response"
                        }
                    }
                }
            }
        },
        "/apps/installed/{id}/params": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "获取插件参数信息",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.AppInstalledParamsResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "修改插件参数信息",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "RequestBody",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AppInstall"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/dto.Response"
                        }
                    }
                }
            }
        },
        "/apps/tags": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "获取插件分类信息",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.Tag"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/apps/{key}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "app update",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "RequestBody",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AppInstalledOperate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/dto.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "插件安装",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "RequestBody",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AppInstall"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/dto.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "插件卸载",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "RequestBody",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AppUnInstall"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/dto.Response"
                        }
                    }
                }
            }
        },
        "/apps/{key}/detail": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "获取插件详情",
                "parameters": [
                    {
                        "type": "string",
                        "default": "zh",
                        "description": "i18n",
                        "name": "language",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/dto.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.AppDetail"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/public/health": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public"
                ],
                "summary": "health",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.PageResult": {
            "type": "object",
            "properties": {
                "items": {},
                "total": {
                    "type": "integer"
                }
            }
        },
        "dto.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {},
                "msg": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "model.App": {
            "type": "object",
            "properties": {
                "class": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "depends_version": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "github": {
                    "type": "string"
                },
                "icon": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "key": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "sort": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "model.Tag": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "key": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "sort": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "request.AppInstall": {
            "type": "object",
            "required": [
                "cpus",
                "docker_compose",
                "memory_limit",
                "params"
            ],
            "properties": {
                "cpus": {
                    "type": "string"
                },
                "docker_compose": {
                    "type": "string"
                },
                "memory_limit": {
                    "type": "string"
                },
                "params": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "request.AppInstalledOperate": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "params": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "request.AppUnInstall": {
            "type": "object"
        },
        "response.AppDetail": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "depends_version": {
                    "type": "string"
                },
                "docker_compose": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "params": {
                    "$ref": "#/definitions/response.AppParams"
                },
                "repo": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "response.AppInstalledParamsResp": {
            "type": "object",
            "properties": {
                "cpus": {
                    "type": "string"
                },
                "docker_compose": {
                    "type": "string"
                },
                "memory_limit": {
                    "type": "string"
                },
                "params": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.FormField"
                    }
                }
            }
        },
        "response.AppParams": {
            "type": "object",
            "properties": {
                "form_fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.FormField"
                    }
                }
            }
        },
        "response.FormField": {
            "type": "object",
            "properties": {
                "default": {
                    "type": "string"
                },
                "env_key": {
                    "type": "string"
                },
                "key": {
                    "type": "string"
                },
                "label": {
                    "type": "string"
                },
                "required": {
                    "type": "boolean"
                },
                "rule": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                },
                "values": {}
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "token",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Doo Store API Documentation",
	Description:      "Description of Doo Store API documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
