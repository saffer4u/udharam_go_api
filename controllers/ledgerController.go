package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saffer4u/udharam/v2/initializers"
	"github.com/saffer4u/udharam/v2/models"
)

type Ledger struct {
	LedgerName string `json:"ledgerName"`
	UserRefer  uint   `json:"userRefer"`
	ID         uint   `json:"id"`
}

func AddLedger(c *gin.Context) {

	// Get the email/pass off req body
	var body struct {
		LedgerName string
		UserRefer  uint
	}
	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	if err := c.Bind(&body); err != nil {

		models.CreateResponse(false, http.StatusBadRequest, "Invailid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	ledger := models.Ledger{
		LedgerName: body.LedgerName,
		UserRefer:  userID,
	}

	// Todo: Make ledger Unique
	result := initializers.DB.Create(&ledger)

	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to create ledger", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusOK, "Ledger created successfully", CreateLedgerResponse(ledger))
	c.JSON(models.Response.Code, models.Response)
}

func CreateLedgerResponse(ledger models.Ledger) Ledger {
	return Ledger{
		LedgerName: ledger.LedgerName,
		UserRefer:  ledger.UserRefer,
		ID:         ledger.ID,
	}
}

func GetLedgers(c *gin.Context) {
	ledgers := []models.Ledger{}

	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	result := initializers.DB.Find(&ledgers, "user_refer = ?", userID)
	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to fetched ledgers", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	responseLedger := &[]Ledger{}

	for i := 0; i < len(ledgers); i++ {
		*responseLedger = append(*responseLedger, CreateLedgerResponse(ledgers[i]))
	}

	models.CreateResponse(true, http.StatusOK, "All ledger are fetched successfully", responseLedger)
	c.JSON(models.Response.Code, models.Response)

}
func GetLedger(c *gin.Context) {

	id := c.Param("id")

	ledger := models.Ledger{}

	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	result := initializers.DB.First(&ledger, id)
	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to fetched ledgers", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	if ledger.UserRefer != userID {
		models.CreateResponse(false, http.StatusNotFound, "Ledger not found", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusOK, "Ledger fetched successfully", CreateLedgerResponse(ledger))
	c.JSON(models.Response.Code, models.Response)

}

func DeleteLedger(c *gin.Context) {
	ledgerId := c.Param("id")

	ledger := models.Ledger{}
	trans := models.Transaction{}

	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	result := initializers.DB.Delete(&ledger, "user_refer = ? AND id = ?", userID, ledgerId)
	if result.Error != nil || result.RowsAffected == 0 {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to delete Ledger", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	result = initializers.DB.Where("user_refer = ? AND ledger_refer = ?", userID, ledgerId).Delete(&trans)
	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to delete Transactions inside ledger", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusOK, "Ledger deleted successfully", nil)
	c.JSON(models.Response.Code, models.Response)
}

func UpdateLedger(c *gin.Context) {
	ledgerId := c.Param("id")

	var body struct {
		LedgerName string
	}

	if len(body.LedgerName) == 0 {
		models.CreateResponse(false, http.StatusBadRequest, "Invailid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	if err := c.Bind(&body); err != nil {

		models.CreateResponse(false, http.StatusBadRequest, "Invailid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	ledger := models.Ledger{}

	result := initializers.DB.Model(&ledger).Where("user_refer = ? AND id = ?", userID, ledgerId).Update("ledger_name", body.LedgerName)

	if result.Error != nil || result.RowsAffected == 0 {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to update Ledger", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusOK, "Ledger updated successfully", CreateLedgerResponse(ledger))
	c.JSON(models.Response.Code, models.Response)
}
