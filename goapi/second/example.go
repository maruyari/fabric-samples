package second

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

func Hello(c *gin.Context){
	c.JSON(http.StatusOK,map[string]string{
		"hello":"Tanya",
	})
}
func Search (c *gin.Context, ID string){
      
	c.JSON(http.StatusOK,map[string]string{
		"hello":"Tanya",
	})
}
func Add()  {
	//call createMarksheet

	
}