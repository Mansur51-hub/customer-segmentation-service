// Code generated by swaggo/swag. DO NOT EDIT.

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
            "name": "Mansur Mamedov",
            "email": "mansyr001mamedov@mail.ru"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/segments": {
            "post": {
                "description": "Create new segment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "segments"
                ],
                "summary": "Create new segment",
                "parameters": [
                    {
                        "description": "Segment slug",
                        "name": "slug",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SegmentInputData"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "segment",
                        "schema": {
                            "$ref": "#/definitions/model.Segment"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "409": {
                        "description": "Conflict"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Delete segment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "segments"
                ],
                "summary": "Delete segment",
                "parameters": [
                    {
                        "description": "Segment slug",
                        "name": "slug",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SegmentInputData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Create new user segments",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create new user segments",
                "parameters": [
                    {
                        "description": "segments info",
                        "name": "segment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.UserSegmentsInputData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "409": {
                        "description": "Conflict"
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.SegmentInfo": {
            "type": "object",
            "required": [
                "slug"
            ],
            "properties": {
                "slug": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 1
                },
                "ttl": {
                    "type": "string"
                }
            }
        },
        "handler.SegmentInputData": {
            "type": "object",
            "required": [
                "slug"
            ],
            "properties": {
                "slug": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 1
                }
            }
        },
        "handler.UserSegmentsInputData": {
            "type": "object",
            "required": [
                "segments_to_add",
                "segments_to_delete"
            ],
            "properties": {
                "segments_to_add": {
                    "type": "array",
                    "maxItems": 10,
                    "items": {
                        "$ref": "#/definitions/handler.SegmentInfo"
                    }
                },
                "segments_to_delete": {
                    "type": "array",
                    "maxItems": 10,
                    "items": {
                        "type": "string"
                    }
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "model.Segment": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "slug": {
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
	BasePath:         "",
	Schemes:          []string{},
	Title:            "customer segmentation service",
	Description:      "swagger",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}