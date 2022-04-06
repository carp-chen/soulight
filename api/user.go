package api

import "github.com/gin-gonic/gin"

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello!"})
}

func UserRegister(c *gin.Context) {

}
