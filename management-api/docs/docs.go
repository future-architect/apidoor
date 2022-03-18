// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api": {
            "get": {
                "description": "Get list of APIs and its information",
                "produces": [
                    "application/json"
                ],
                "summary": "Get list of API info.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.APIInfoList"
                        }
                    }
                }
            },
            "post": {
                "description": "Get list of APIs and its information",
                "produces": [
                    "application/json"
                ],
                "summary": "Get list of API information",
                "parameters": [
                    {
                        "description": "api information",
                        "name": "api_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostAPIInfoReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/search": {
            "get": {
                "description": "Get list of APIs and its information",
                "produces": [
                    "application/json"
                ],
                "summary": "search for api info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search query words (split words by '.', ex: 'foo.bar'). If containing multiple words, items which match the all search words return",
                        "name": "q",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "all",
                        "description": "search target fields. You can choose field(s) from 'all' (represents searching all fields), 'name', 'description', or 'source'. (if there are multiple target fields, split target by '.', ex: 'name.source')",
                        "name": "target_fields",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "exact",
                            "partial"
                        ],
                        "type": "string",
                        "default": "partial",
                        "description": "pattern match, chosen from 'exact' or 'partial'",
                        "name": "pattern_match",
                        "in": "query"
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "default": 50,
                        "description": "the maximum number of results",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "the starting point for the result set",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.SearchAPIInfoResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/token": {
            "post": {
                "description": "post api tokens for calling external api",
                "produces": [
                    "application/json"
                ],
                "summary": "post api tokens for call external api",
                "parameters": [
                    {
                        "description": "api token description",
                        "name": "tokens",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostAPITokenReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete api tokens for calling external api",
                "summary": "delete api tokens for call external api",
                "parameters": [
                    {
                        "description": "target api_key",
                        "name": "api_key",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "target api_key",
                        "name": "path",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "$ref": "#/definitions/model.EmptyResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/contract": {
            "post": {
                "description": "Post an API product",
                "produces": [
                    "application/json"
                ],
                "summary": "Post a product",
                "parameters": [
                    {
                        "description": "contract definition",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostContractReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "checks whether this API works correctly or not",
                "produces": [
                    "text/plain"
                ],
                "summary": "checks if API works",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/keys": {
            "post": {
                "description": "post api key used for authentication in apidoor gateway",
                "produces": [
                    "application/json"
                ],
                "summary": "post api key",
                "parameters": [
                    {
                        "description": "api key owner",
                        "name": "api_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostAPIKeyReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.PostAPIKeyResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/keys/products": {
            "post": {
                "description": "Post relationship between api key and authorized products linked to the key",
                "produces": [
                    "application/json"
                ],
                "summary": "Post relationship between api key and authorized products linked to the key",
                "parameters": [
                    {
                        "description": "relationship between apikey and products linked to the apikey",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostAPIKeyProductsReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/products": {
            "post": {
                "description": "Post an API product",
                "produces": [
                    "application/json"
                ],
                "summary": "Post a product",
                "parameters": [
                    {
                        "description": "product definition",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostProductReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/routing": {
            "post": {
                "description": "Post a new API routing",
                "produces": [
                    "application/json"
                ],
                "summary": "Post an API routing",
                "parameters": [
                    {
                        "description": "routing parameters",
                        "name": "api_routing",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostAPIRoutingReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Create a user",
                "produces": [
                    "application/json"
                ],
                "summary": "Create a user",
                "parameters": [
                    {
                        "description": "user description",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PostUserReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/validator.BadRequestResp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.APIContent": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.APIInfo": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "swagger_url": {
                    "type": "string"
                },
                "thumbnail": {
                    "type": "string"
                }
            }
        },
        "model.APIInfoList": {
            "type": "object",
            "properties": {
                "api_info_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.APIInfo"
                    }
                }
            }
        },
        "model.AccessToken": {
            "type": "object",
            "required": [
                "key",
                "param_type",
                "value"
            ],
            "properties": {
                "key": {
                    "type": "string"
                },
                "param_type": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "model.AuthorizedContractProducts": {
            "type": "object",
            "required": [
                "contract_id"
            ],
            "properties": {
                "contract_id": {
                    "type": "integer"
                },
                "product_ids": {
                    "description": "ProductIDs is the list of product which is linked to the key\nif this field is nil or empty, all products the contract contains will be linked",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "model.ContractProducts": {
            "type": "object",
            "required": [
                "product_name"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "product_name": {
                    "type": "string"
                }
            }
        },
        "model.EmptyResp": {
            "type": "object"
        },
        "model.PostAPIInfoReq": {
            "type": "object",
            "required": [
                "description",
                "name",
                "source",
                "swagger_url",
                "thumbnail"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "swagger_url": {
                    "type": "string"
                },
                "thumbnail": {
                    "type": "string"
                }
            }
        },
        "model.PostAPIKeyProductsReq": {
            "type": "object",
            "required": [
                "apikey_id",
                "contracts"
            ],
            "properties": {
                "apikey_id": {
                    "type": "integer"
                },
                "contracts": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/model.AuthorizedContractProducts"
                    }
                }
            }
        },
        "model.PostAPIKeyReq": {
            "type": "object",
            "required": [
                "user_account_id"
            ],
            "properties": {
                "user_account_id": {
                    "type": "string"
                }
            }
        },
        "model.PostAPIKeyResp": {
            "type": "object",
            "properties": {
                "access_key": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_account_id": {
                    "type": "string"
                }
            }
        },
        "model.PostAPIRoutingReq": {
            "type": "object",
            "required": [
                "api_key",
                "forward_url",
                "path"
            ],
            "properties": {
                "api_key": {
                    "type": "string"
                },
                "forward_url": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "model.PostAPITokenReq": {
            "type": "object",
            "required": [
                "api_key",
                "path",
                "tokens"
            ],
            "properties": {
                "api_key": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "tokens": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.AccessToken"
                    }
                }
            }
        },
        "model.PostContractReq": {
            "type": "object",
            "required": [
                "products",
                "user_id"
            ],
            "properties": {
                "products": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/model.ContractProducts"
                    }
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "model.PostProductReq": {
            "type": "object",
            "required": [
                "description",
                "name",
                "source",
                "thumbnail"
            ],
            "properties": {
                "api_contents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.APIContent"
                    }
                },
                "description": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string"
                },
                "is_available": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "thumbnail": {
                    "type": "string"
                }
            }
        },
        "model.PostUserReq": {
            "type": "object",
            "required": [
                "account_id",
                "email_address",
                "password"
            ],
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "email_address": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.ResultSet": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                }
            }
        },
        "model.SearchAPIInfoMetaData": {
            "type": "object",
            "properties": {
                "result_set": {
                    "$ref": "#/definitions/model.ResultSet"
                }
            }
        },
        "model.SearchAPIInfoResp": {
            "type": "object",
            "properties": {
                "api_info_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.APIInfo"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/model.SearchAPIInfoMetaData"
                }
            }
        },
        "validator.BadRequestResp": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "validation_errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/validator.ValidationError"
                    }
                }
            }
        },
        "validator.ValidationError": {
            "type": "object",
            "properties": {
                "constraint_type": {
                    "type": "string"
                },
                "enum": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "field": {
                    "type": "string"
                },
                "got": {},
                "gte": {
                    "type": "string"
                },
                "lte": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "ne": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/mgmt",
	Schemes:          []string{},
	Title:            "Management API",
	Description:      "This is an API that manages products.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
