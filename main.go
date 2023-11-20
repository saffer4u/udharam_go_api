package main

import (
	"github.com/gin-gonic/gin"
	"github.com/saffer4u/udharam/v2/controllers"
	"github.com/saffer4u/udharam/v2/initializers"
	"github.com/saffer4u/udharam/v2/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	// initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	//? User : ===>
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	//? Ledger : ===>
	r.POST("/ledger", middleware.RequireAuth, controllers.AddLedger)
	r.GET("/ledger", middleware.RequireAuth, controllers.GetLedgers)
	r.GET("/ledger/:id", middleware.RequireAuth, controllers.GetLedger)
	r.DELETE("/ledger/:id", middleware.RequireAuth, controllers.DeleteLedger)
	r.PUT("/ledger/:id", middleware.RequireAuth, controllers.UpdateLedger)

	//? Transaction : ===>
	r.POST("/transaction", middleware.RequireAuth, controllers.AddTransaction)
	r.GET("/transaction/:id", middleware.RequireAuth, controllers.GetTransaction)
	r.DELETE("/transaction/:id", middleware.RequireAuth, controllers.DeleteTransaction)
	r.POST("/transactionsByLedgerId", middleware.RequireAuth, controllers.GetTransactions)

	r.Run()
}
