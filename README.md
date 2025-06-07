# Service System REST API

A simple REST API for user management and authentication using JWT, built with **Golang** using the [Gin](https://gin-gonic.com/) web framework, [GORM](https://gorm.io/) ORM, and PostgreSQL database.

---

## Project Structure

```
.
├── controllers/
│   ├── auth_controller.go
│   └── user_controller.go
├── database/
│   └── database.go
├── models/
│   ├── user.go
│   └── user_session.go
├── router/
│   └── router.go
├── services/
│   ├── auth_service.go
│   └── user_service.go
├── utils/
│   └── jwt.go
├── main.go
├── go.mod
├── go.sum
└── .env
```

---

## main.go Explanation

```go
package main

import (
    "service_system/controllers"
    "service_system/database"
    "service_system/router"
)

func main() {
    database.Init()                
    controllers.SetDB(database.DB)  

    r := router.SetupRouter()           
    r.Run(":8800")               
}
```

---

## Getting Started

1. **Clone the repo & install dependencies**
    ```bash
    git clone <repo-url>
    cd service_system
    go mod tidy
    ```

2. **Create a `.env` file in the project root:**
    ```
    DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
    JWT_SECRET_KEY=your_jwt_secret_key
    ```

3. **Run the application**
    ```bash
    go run main.go
    ```
    The server will run at `http://localhost:8800`

---

## API Documentation

### 1. Register User

- **POST** `/signup`
- **Request Body:**
    ```json
    {
      "username": "admin",
      "password": "admin123",
      "email": "admin@example.com"
    }
    ```
- **Response:**
    ```json
    {
      "message": "User registered successfully",
      "user": { ... }
    }
    ```

---

### 2. Login User

- **POST** `/login`
- **Request Body:**
    ```json
    {
      "username": "admin",
      "password": "admin123"
    }
    ```
- **Response:**
    ```json
    {
      "token": "<jwt_token>"
    }
    ```
- **Note:**  
  If already logged in, you cannot login again before logging out. Session is valid for 24 hours.

---

### 3. Logout User

- **POST** `/logout`
- **Header:**  
  `Authorization: Bearer <jwt_token>`
- **Or Request Body:**
    ```json
    {
      "token": "<jwt_token>"
    }
    ```
- **Response:**
    ```json
    {
      "message": "Logout successful"
    }
    ```

---

### 4. Get All Users

- **GET** `/users`
- **Response:**
    ```json
    [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        ...
      },
      ...
    ]
    ```

---

### 5. Soft Delete User

- **DELETE** `/users/:id`
- **Example:**  
  `DELETE http://localhost:8800/users/1`
- **Response:**
    ```json
    {
      "message": "User soft deleted successfully"
    }
    ```

---

## Testing with Postman

- **Register:**  
  Method: POST, URL: `/signup`, Body: raw JSON
- **Login:**  
  Method: POST, URL: `/login`, Body: raw JSON
- **Logout:**  
  Method: POST, URL: `/logout`, Header: Authorization: Bearer `<token>`
- **Get Users:**  
  Method: GET, URL: `/users`
- **Delete User:**  
  Method: DELETE, URL: `/users/1`

---

## License

MIT

---

**Happy coding!**  
If you have any questions, feel free to open an issue in this repository.

