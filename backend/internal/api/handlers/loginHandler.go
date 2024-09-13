package handlers

// import (
// 	"net/http"
// 	"encoding/json"
// 	"github.com/CEM-KEA/whoknows/backend/internal/security"
// 	"github.com/CEM-KEA/whoknows/backend/internal/models"
// )

// type LoginRequest struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required"`
// }

// type LoginResponse struct {
// 	Token string `json:"token"`
// }

// STEPS:
// 1. Create a new struct called LoginRequest with the following fields:
//    - Email string
//    - Password string
// 2. Create a new struct called LoginResponse with the following fields:
//    - Token string
// 3. Implement the Login handler function that takes a http.ResponseWriter and http.Request as arguments.
//    - Parse the request body into a LoginRequest struct.
//    - Retrieve the user from the database by email.
//    - Check if the user exists and if the password matches.
//    - Generate a JWT token for the user.
//    - Return the token in the response.
// 4. Add the Login handler to the router in the NewRouter function in router.go.
// 5. Test the login functionality using a REST client like Postman.
// 6. Write unit tests for the Login handler function.
