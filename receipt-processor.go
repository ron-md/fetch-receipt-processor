// ok hi !
// so with this takehome we are asked to create a webservice with essentially two api calls
// first is the POST endpoint that process the receipts, here i will take in the receipt json and produce an ID
// should i also calc points here???
// second is the GET endpoint that returns the points for each receipt
// one thing i will have to handle is generating the Id, it seems like a UUID v4, hope importing is ok 
package main 

import (
    "net/http"

    "github.com/google/uuid"

    "github.com/gin-gonic/gin"

    "strconv"

    "math"

    "strings"
)


type Receipt struct {
	
	Id		string	`json:"Id"` 
	Retailer	string	`json:"Retailer" binding:"required"`
	PurchaseDate	string	`json:"PurchaseDate" binding:"required"`
	PurchaseTime	string	`json:"PurchaseTime" binding:"required"`
	Items		[]Item	`json:"Items" binding:"required,dive"`	
	Total		string	`json:"Total" binding:"required"`
	Points		int	`json:"Points"`
}

type Item struct {

	ShortDescription	string	`json:"ShortDescription" binding:"required"`
	Price			string 	`json:"Price" binding:"required"`

}

//storing in memory
var Receipts = []Receipt{}


func main() {
	router := gin.Default()
	//should just be two 
	router.POST("/receipts/process", processReceipts)
	router.GET("/receipts/:Id/points", getPoints)

	router.Run("localhost:8080")
}


//logic below


func processReceipts(c *gin.Context){
	var newReceipt Receipt
	

	if err:= c.ShouldBindJSON(&newReceipt); err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{
            	"error": err.Error(),
        	})
        	return
	}

	newReceipt.Id = uuid.New().String()
	newReceipt.Points = 0 
	Receipts = append(Receipts, newReceipt)

	c.JSON(http.StatusCreated, gin.H{
		"id": newReceipt.Id,
	})

}


func getPoints(c *gin.Context){

	Id:= c.Param("Id")
	var pointReceipt Receipt

	//first check that its even there
	isThere := false
	for _, a:= range Receipts{
		if a.Id == Id {
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
	for _, x := range pointReceipt.Retailer{
		if int(x) >= 48 && int(x) <= 57 || int(x) >= 65 && int(x) <= 90 || int(x) >= 97 && int(x) <=122{
			totalPoints += 1
		}
	}
	

	//round total 
	for i, x := range pointReceipt.Total{
		if int(x) == 46{ //period
			if pointReceipt.Total[i+1] == 48 && pointReceipt.Total[i+2] == 48 {
				totalPoints += 50
				break
				}
			}
		}

	// multiple of 0.25
	if multOfQuart, err := strconv.ParseFloat(pointReceipt.Total, 64); err == nil {
		
		if math.Mod(multOfQuart, 0.25) == 0 {
			totalPoints += 25
		}


	}

	//5 for every 2 items
	itemLength := len(pointReceipt.Items)
	if (itemLength % 2) == 1 {
		itemLength -= 1

	}

	totalPoints += (itemLength / 2) * 5
		
	

	//trimmed descriptions
	for _, x := range pointReceipt.Items{
		trimmedDesc := strings.TrimSpace(x.ShortDescription)
		itemPrice, _  := strconv.ParseFloat(x.Price, 64)
		if len(trimmedDesc) % 3 == 0 {
			itemPrice = itemPrice * 0.2 
			totalPoints += int(math.Ceil(itemPrice))
		}
	}
		

	//if purchase date is odd
	getDate := pointReceipt.PurchaseDate[len(pointReceipt.PurchaseDate)-2:len(pointReceipt.PurchaseDate)]
	
	getDateInt, _ := strconv.Atoi(getDate)

	if getDateInt % 2 == 1 {
		
		totalPoints += 6 		
	}

	//purchase time between 2pm and 4pm
	getHour := pointReceipt.PurchaseTime[:2]

	getHourInt, _ := strconv.Atoi(getHour)

	//AFTER 2:00pm, i'll be thorough at expense of being annoying lol
	getMinute := pointReceipt.PurchaseTime[len(pointReceipt.PurchaseTime)-1:len(pointReceipt.PurchaseTime)]
	
	getMinuteInt, _ := strconv.Atoi(getMinute)

	if (getHourInt >=14 && getMinuteInt > 0) && getHourInt < 16{

		totalPoints += 10
		}


	c.JSON(http.StatusOK, gin.H{
		"points": totalPoints,
	})

	
}

