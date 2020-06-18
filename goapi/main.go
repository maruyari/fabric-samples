package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

type Records struct {
	name  string
	year  string
	board string
	mark  string
}

func find(c *gin.Context) {
	roll := c.Param("rollno")
	output, err:= exec.Command("./find.sh", roll).Output()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"err": err})
	}
	c.JSON(http.StatusOK,  string(output))



}
func modify(c *gin.Context) {

	roll := c.Param("rollno")
	newMarks := c.Param("marks")
	output, err:= exec.Command("./update.sh", roll,newMarks).Output()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"err": err})
	}
	c.JSON(http.StatusOK,  string(output))
}
func all(c *gin.Context) {
	output, err := exec.Command("./AllRecords.sh").Output()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err})
	}
	c.JSON(http.StatusOK,  string(output))
}
func add(c *gin.Context) {
	name:=c.Param()
	//todo

	output, err := exec.Command("./Add.sh", name, year,board, mark, roll ).Output()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err})
	}
	c.JSON(http.StatusOK,  string(output))
}
func main() {
	r := gin.Default()
	r.GET("/modify/:rollno/:marks", modify)
	//r.POST("/add",add)
	r.GET("/find/:rollno", find)
	r.GET("/all", all)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
