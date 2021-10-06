// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
            "post": {
                "description": "Post a new API routing",
                "produces": [
                    "application/json"
                ],
                "summary": "Post API routing",
                "parameters": [
                    {
                        "description": "routing parameters",
                        "name": "api_routing",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/managementapi.PostAPIRoutingReq"
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
                            "type": "string"
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
        "/product": {
            "post": {
                "description": "Get list of APIs and its information",
                "produces": [
                    "application/json"
                ],
                "summary": "Get list of products",
                "parameters": [
                    {
                        "description": "api information",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/managementapi.PostProductReq"
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
                            "type": "string"
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
            "get": {
                "description": "Get list of APIs and its information",
                "produces": [
                    "application/json"
                ],
                "summary": "Get list of products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/managementapi.Products"
                        }
                    }
                }
            }
        },
        "/products/search": {
            "get": {
                "description": "Get list of APIs and its information",
                "produces": [
                    "application/json"
                ],
                "summary": "search for products",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search query words (split words by '.', ex: 'foo.bar'). If q contains multiple words, items which contain the all search words return",
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
                            "$ref": "#/definitions/managementapi.SearchProductsResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
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
        "managementapi.PostAPIRoutingReq": {
            "type": "object",
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
        "managementapi.PostProductReq": {
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
        "managementapi.Product": {
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
        "managementapi.Products": {
            "type": "object",
            "properties": {
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/managementapi.Product"
                    }
                }
            }
        },
        "managementapi.ResultSet": {
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
        "managementapi.SearchProductsMetaData": {
            "type": "object",
            "properties": {
                "result_set": {
                    "$ref": "#/definitions/managementapi.ResultSet"
                }
            }
        },
        "managementapi.SearchProductsResp": {
            "type": "object",
            "properties": {
                "metadata": {
                    "$ref": "#/definitions/managementapi.SearchProductsMetaData"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/managementapi.Product"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/mgmt",
	Schemes:     []string{},
	Title:       "Management API",
	Description: "This is an API that manages products.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
