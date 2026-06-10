package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "strconv"
    "time"
)

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"unique"`
    Name     string
    Password string
    Role     string `gorm:"default:user"`
}

type Task struct {
    ID          uint
    Title       string
    Description string
    Status      string `gorm:"default:pending"`
    AssignedTo  *uint
    CreatedBy   uint
}

var db *gorm.DB
var jwtKey = []byte("secret")

func main() {
    // init db
    dsn := "host=localhost user=postgres dbname=taskdb password=secret port=5432"
    db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    db.AutoMigrate(&User{}, &Task{})

    r := gin.Default()

    r.POST("/auth/register", register)
    r.POST("/auth/login", login)
    authorized := r.Group("/")
    authorized.Use(authMiddleware())
    authorized.GET("/tasks", getTasks)

    r.Run(":8080")
}

func register(c *gin.Context) {
    var input struct{ Email, Name, Password string }
    c.BindJSON(&input)
    hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    user := User{Email: input.Email, Name: input.Name, Password: string(hashed)}
    db.Create(&user)
    c.JSON(200, user)
}

func login(c *gin.Context) {
    var input struct{ Email, Password string }
    c.BindJSON(&input)
    var user User
    db.Where("email = ?", input.Email).First(&user)
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(401, gin.H{"error": "invalid"})
        return
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "role": user.Role,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })
    tokenString, _ := token.SignedString(jwtKey)
    c.JSON(200, gin.H{"token": tokenString})
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) { return jwtKey, nil })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("user_id", uint(claims["sub"].(float64)))
        c.Set("role", claims["role"].(string))
        c.Next()
    }
}

func getTasks(c *gin.Context) {
    userID := c.MustGet("user_id").(uint)
    role := c.MustGet("role").(string)
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    status := c.Query("status")
    query := db.Model(&Task{})
    if role != "admin" {
        query = query.Where("assigned_to = ?", userID)
    }
    if status != "" {
        query = query.Where("status = ?", status)
    }
    var tasks []Task
    query.Offset((page - 1) * limit).Limit(limit).Find(&tasks)
    c.JSON(200, tasks)
}
