package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saffer4u/udharam/v2/initializers"
	"github.com/saffer4u/udharam/v2/models"
)

type Transaction struct {
	Title       string  `json:"title"`
	Amount      float32 `json:"amount"`
	LedgerRefer uint    `json:"ledgerRefer"`
	ID          uint    `json:"id"`
}

type Transactions struct {
	TotalAmount float32 `json:"total_amount"`
	Transaction []Transaction
}

func CreateTransactionResponse(trans []Transaction, TotalAmount float32) Transactions {
	return Transactions{
		TotalAmount: TotalAmount,
		Transaction: trans,
	}
}

func CreateTransaction(trans models.Transaction) Transaction {
	return Transaction{
		Title:       trans.Title,
		Amount:      trans.Amount,
		LedgerRefer: trans.LedgerRefer,
		ID:          trans.ID,
	}
}

func AddTransaction(c *gin.Context) {

	// Get the email/pass off req body
	var body struct {
		Title       string
		LedgerRefer uint
		Amount      float32
	}
	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	if err := c.Bind(&body); err != nil {

		models.CreateResponse(false, http.StatusBadRequest, "Invailid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	ledger := models.Ledger{}

	result := initializers.DB.Find(&ledger, "user_refer = ? AND id = ?", userID, body.LedgerRefer)

	if result.RowsAffected == 0 {
		models.CreateResponse(false, http.StatusBadRequest, "Ledger not found", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	transaction := models.Transaction{
		Title:       body.Title,
		LedgerRefer: body.LedgerRefer,
		UserRefer:   userID,
		Amount:      body.Amount,
	}

	result = initializers.DB.Create(&transaction)

	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to Add transaction", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusCreated, "Transaction added successfully", CreateTransaction(transaction))
	c.JSON(models.Response.Code, models.Response)
}

func GetTransactions(c *gin.Context) {
	var body struct {
		LedgerId uint
	}

	if err := c.Bind(&body); err != nil {

		models.CreateResponse(false, http.StatusBadRequest, "Invailid request body", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	transactions := []models.Transaction{}

	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	result := initializers.DB.Find(&transactions, "user_refer = ? AND ledger_refer = ?", userID, body.LedgerId)
	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to fetched Transactions", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	s := strconv.FormatUint(uint64(body.LedgerId), 7)
	if result.RowsAffected == 0 {
		models.CreateResponse(false, http.StatusBadRequest, ("Ledger with id " + s + " not found"), nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}
	transactionList := &[]Transaction{}

	var totalAmount float32

	for i := 0; i < len(transactions); i++ {
		totalAmount += transactions[i].Amount
		*transactionList = append(*transactionList, CreateTransaction(transactions[i]))
	}

	models.CreateResponse(true, http.StatusOK, "All transactions are fetched successfully", CreateTransactionResponse(*transactionList, totalAmount))
	c.JSON(models.Response.Code, models.Response)

}

func GetTransaction(c *gin.Context) {

	transId := c.Param("id")

	transaction := models.Transaction{}

	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	result := initializers.DB.First(&transaction, "user_refer = ? AND id = ?", userID, transId)
	if result.Error != nil {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to fetched Transaction", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusOK, "Transaction fetched successfully", CreateTransaction(transaction))
	c.JSON(models.Response.Code, models.Response)

}

func DeleteTransaction(c *gin.Context) {
	transId := c.Param("id")

	transaction := models.Transaction{}

	userObj, _ := c.Get("user")
	userID := userObj.(models.User).ID

	result := initializers.DB.Delete(&transaction, "user_refer = ? AND id = ?", userID, transId)
	if result.Error != nil || result.RowsAffected == 0 {
		models.CreateResponse(false, http.StatusBadRequest, "Unable to delete Transaction", nil)
		c.JSON(models.Response.Code, models.Response)
		return
	}

	models.CreateResponse(true, http.StatusOK, "Transaction deleted successfully", nil)
	c.JSON(models.Response.Code, models.Response)
}
