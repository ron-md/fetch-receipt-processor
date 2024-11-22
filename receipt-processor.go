// ok hi !
// so with this takehome we are asked to create a webservice with essentially two api calls
// first is the POST endpoint that process the receipts, here i will take in the receipt json and produce an ID
// should i also calc points here???
// second is the GET endpoint that returns the points for each receipt
// one thing i will have to handle is generating the id, it seems like a UUID v4, hope importing is ok 
package main 

import (
    "net/http"

    "github.com/google/uuid"

    "github.com/gin-gonic/gin"

)


type receipt struct {
	
	id		string	`json:"id"`
	retailer	string	`json:"retailer"`
	purchaseDate	string	`json:"purchaseDate"`
	purchaseTime	string	`json:"purchaseTime"`
	items		[]item	`json:"items"`	
	total		float64	`json:"total"`
	points		int	`json:"points"`
}

type item struct {

	shortDescription	string	`json:"shortDescription"`
	price			float64 `json:"price"`

}

var receipts = []receipt{}


func main() {
	router := gin.Default()
	//should just be two 
	router.POST("/receipts/process", processReceipts)
	router.GET("/receipts/:id/points", getPoints)

	router.Run("localhost:8080")
}


//logic below


func processReceipts(c *gin.Context){
	var newReceipt receipt
	
	if err:= c.BindJSON(&newReceipt); err != nil {
		return
	}

	newReceipt.id = uuid.New().String()
	newReceipt.points = 0 
	receipts = append(receipts, newReceipt)

	c.JSON(http.StatusCreated, gin.H{
		"receipt": newReceipt.id,
	})

}


func getPoints(c *gin.Context){

	id:= c.Param("id")
	var pointReceipt receipt

	//first check that its even there
	isThere := false
	for _, a:= range receipts{
		if a.id == id {
			isThere = true
			pointReceipt = a
		}
	}

	if isThere == false{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request - that receipt does not exist",
		})
	}


	totalPoints := 0

	//characters
	for _, x := range pointReceipt.retailer{
		if int(x) >= 48 || int(x) <= 57 || int(x) >= 65 || int(x) <= 90 || int(x) >= 97 || int(x) <=122{
			totalPoints += 1
		}
	}
	

	c.JSON(http.StatusOK, gin.H{
		"points":totalPoints,
		"test":pointReceipt.retailer,
	})
}

