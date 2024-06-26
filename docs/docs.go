// Package docs GENERATED BY SWAG; DO NOT EDIT
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
        "/download/pdf/{origin}/{target}": {
            "get": {
                "description": "通过申报国家确定税金文件路径",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "download"
                ],
                "summary": "下载税金单文件",
                "parameters": [
                    {
                        "type": "string",
                        "description": "下载的源文件名,不带文件后缀",
                        "name": "origin",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "下载文件后，将文件重命名的文件名，没有后缀将自动添加pdf作为后缀",
                        "name": "target",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "申报国家(BE|NL),默认为NL",
                        "name": "dc",
                        "in": "query"
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
        "/download/xml/{dc}/{filename}": {
            "get": {
                "description": "通过文件名前缀确定文件路径",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "download"
                ],
                "summary": "下载export文件",
                "parameters": [
                    {
                        "type": "string",
                        "description": "申报国家(BE|NL)",
                        "name": "dc",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "export文件的文件名",
                        "name": "filename",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "是否下载，1表示直接下载",
                        "name": "download",
                        "in": "query"
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
        "/export/list/{dc}": {
            "get": {
                "description": "通过指定的申报国家获取其当前Export监听路径下的文件列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "export"
                ],
                "summary": "获取Export监听路径下当前文件列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "申报国家(BE|NL)",
                        "name": "dc",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/softpak.ExportFileListDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    }
                }
            }
        },
        "/export/remover/{dc}": {
            "post": {
                "description": "页面选取Export文件并发送文件完整路径，将Export文件重新移入Export监听路径中，触发文件的CREATE监听，从而重新发送",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "export"
                ],
                "summary": "重新发送Export XML 文件",
                "parameters": [
                    {
                        "type": "string",
                        "description": "申报国家(BE|NL)",
                        "name": "dc",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "需要重新发送的文件完整路径",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.FileResendRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/softpak.ExportFileListDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    }
                }
            }
        },
        "/search/file": {
            "post": {
                "description": "可检索税金单文件以及export报关结果文件，可使用文件名部分做模糊匹配。建议使用Job number 进行检索",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "search"
                ],
                "summary": "搜索文件",
                "parameters": [
                    {
                        "description": "检索内容",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.SearchFileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/softpak.SearchFileResult"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/util.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "softpak.ExportFileListDTO": {
            "type": "object",
            "properties": {
                "filename": {
                    "description": "Type TAX_BILL, EXPORT_XML",
                    "type": "string"
                },
                "filepath": {
                    "type": "string"
                },
                "modifiedTime": {
                    "description": "修改时间",
                    "type": "string"
                },
                "size": {
                    "description": "文件大小 bytes",
                    "type": "integer"
                }
            }
        },
        "softpak.SearchFileResult": {
            "type": "object",
            "properties": {
                "filename": {
                    "type": "string"
                },
                "filepath": {
                    "type": "string"
                },
                "searchText": {
                    "type": "string"
                },
                "type": {
                    "description": "Type TAX_BILL, EXPORT_XML",
                    "type": "string"
                }
            }
        },
        "util.ResponseError": {
            "type": "object",
            "properties": {
                "errors": {
                    "description": "Error messages",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "description": "Status success, fail",
                    "type": "string"
                }
            }
        },
        "web.FileResendRequest": {
            "type": "object",
            "required": [
                "filePaths",
                "inListeningPath"
            ],
            "properties": {
                "filePaths": {
                    "description": "FilePaths 需要重新发送的文件名",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "inListeningPath": {
                    "description": "InListeningPath 是否是监听路径中的文件",
                    "type": "boolean"
                }
            }
        },
        "web.SearchFileRequest": {
            "type": "object",
            "required": [
                "declareCountry",
                "filenames",
                "type"
            ],
            "properties": {
                "declareCountry": {
                    "description": "DeclareCountry NL, BE",
                    "type": "string"
                },
                "filenames": {
                    "description": "Filenames Support use Job number",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "month": {
                    "description": "Month exp: 09",
                    "type": "string"
                },
                "type": {
                    "description": "Type TAX_BILL, EXPORT_XML",
                    "type": "string"
                },
                "year": {
                    "description": "Year exp: 2022",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
