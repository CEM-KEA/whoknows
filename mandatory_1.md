## Legacy codebase dependency graphs

### Before upgrading Python and dependencies
![fb59e541-6bbe-4be4-8b1f-47638736bc83](https://github.com/user-attachments/assets/83eb04a8-2354-466c-932e-47d9abba3838)

### After upgrading Pythong and dependencies
![f8be1849-b6e5-49c3-820e-a42d3678abf7](https://github.com/user-attachments/assets/21f32f2c-7587-49f9-9954-c39e331d753a)


## Problems with legacy codebase (listed by severity high->low)
1. Outdated Language and dependencies 

2. admin password in clear text in codebase 

3. SQL injection 

4. Hashing with MD5 

5. Hardcoded configuration values 

6. no CSFR protection 

7. Single file 


## Branching strategy 
- What version control strategy did you choose and how did you actually do it / enforce it? 

  - Trunk based GitHub flow. 

  - When starting to work on a new issue, we create a branch for that issue. When done we create a pull request to main from the issue branch. All branches are based on trunk(main) and merged through pull requests directly to main on completion. When done reviewing and merging a pull request we delete the branch. 

- Why did your group choose the one you did? Why did you not choose others? 

  - We made a pros/cons for each branching strategy and figured that trunk based GitHub flow was the right one for us. As it is relatively low complexity (fewer active branches) and allows for rapid integration/regular deployment which results in getting feedback more often/faster. 

  - Some of the disadvantages to trunk based GitHub flow is that it requires discipline/thorough testing/review of new changes before they are merged into trunk. This is achieved through GitHub flow by using pull requests that require a separate reviewer to accept the changes, in addition to workflows that run tests + static code analysis on pull requests. 

  - The reason we did not choose others, such as feature/Release branching/Git flow, was due to the higher complexity introduced by having multiple long-lived branches, which can provide some good isolation/control for larger teams and projects, but we felt it was overkill for us.  

- What advantages and disadvantages did you run into during the course? 

  - Advantages: close to no merge conflicts. We have a very nice overview over what is currently being worked on and as soon as an issue is done the others can benefit from it. Enforcing only using peer reviewed PRs forces us to dive into what the others have made. 

  - Disadvantages: Using only peer reviewed PRs to merge to trunk, makes making small changes/debugging workflows tedious. 
  Feel free to revise the document even after deadline with new insights during the course. 


## OpenAPI specification
*Note: we couldn't find any Go library to generate OpenAPI specifications newer than Swagger 2.0.*
```json
{
    "swagger": "2.0",
    "info": {
        "description": "This is the API for the WhoKnows application",
        "title": "WhoKnows API",
        "contact": {},
        "version": "1.0"
    },
    "basePath": "/",
    "paths": {
        "/api/login": {
            "post": {
                "description": "Login with username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid username or password",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/logout": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Logs out the user by revoking the jwt token",
                "responses": {
                    "200": {
                        "description": "Logged out successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid Authorization header format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to revoke token",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "User data",
                        "name": "register",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Validation error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to create user",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/search": {
            "get": {
                "description": "Search for pages by content",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search query",
                        "name": "q",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Language filter",
                        "name": "language",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SearchResponse"
                        }
                    },
                    "400": {
                        "description": "Search query (q) is required",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Search query failed",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/validate-login": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Validates the jwt token",
                "responses": {
                    "200": {
                        "description": "valid",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Token expired/revoked",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/weather": {
            "get": {
                "description": "Get weather information",
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.WeatherResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch weather data",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handlers.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "handlers.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "password2": {
                    "description": "Password2 is used to confirm the password, it is optional, so it is omitted if it is not provided or an empty string\nIf it is provided, it must be equal to the Password field",
                    "type": "string"
                },
                "username": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 3
                }
            }
        },
        "handlers.SearchResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "additionalProperties": true
                    }
                }
            }
        },
        "handlers.WeatherResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "JWT",
            "in": "header"
        }
    }
}
```
