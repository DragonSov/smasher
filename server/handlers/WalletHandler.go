package handlers

import (
	"fmt"
	"github.com/DragonSov/smasher/server/domain/Users"
	"github.com/DragonSov/smasher/server/domain/Wallets"
	"github.com/DragonSov/smasher/server/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

type WalletHandlers struct {
	Service services.WalletService
}

func (h WalletHandlers) CreateWallet(c *gin.Context) {
	// Getting the current user and creating a wallet model
	u, _ := c.Get("authUser")
	w := Wallets.WalletModel{}

	// Checking form fields
	w.Currency, w.UserID = c.PostForm("currency"), u.(Users.UserModel).ID
	if w.Currency == "" || w.Currency == " " {
		c.JSON(406, gin.H{
			"status":  "error",
			"message": "The form fields are not filled in",
		})
		return
	}

	// Creating new wallet
	_, err := h.Service.CreateWallet(w)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error has occurred",
		})
		return
	}

	// Return of the currency of the new wallet
	c.JSON(200, gin.H{
		"status":   "success",
		"currency": w.Currency,
	})
}

func (h WalletHandlers) SelectUserWallets(c *gin.Context) {
	// Getting the current user
	u, _ := c.Get("authUser")

	// Selecting user wallets
	wallets, err := h.Service.SelectUserWallets(u.(Users.UserModel).ID)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get the wallets. Try again later",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   wallets,
	})
}

func (h WalletHandlers) SelectWalletByUUID(c *gin.Context) {
	// Getting a Wallet ID
	walletID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get the wallet ID. Try again later",
		})
		return
	}

	// Selecting this wallet
	wallet, err := h.Service.SelectWalletByUUID(walletID)
	if err != nil || wallet == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "Wallet not found",
		})
		return
	}

	// Return information about this wallet
	c.JSON(200, gin.H{
		"status": "success",
		"data":   wallet,
	})
}

func (h WalletHandlers) ReplenishWalletByUUID(c *gin.Context) {
	// Parsing wallet ID
	walletID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get the wallet ID. Try again later",
		})
		return
	}

	moneyString := c.PostForm("money")

	// Parsing a number from a string of amount of money
	money, err := strconv.Atoi(moneyString)
	if err != nil || money <= 0 {
		c.JSON(406, gin.H{
			"status":  "error",
			"message": "Incorrect amount of money to send",
		})
		return
	}

	// Selecting wallet
	wallet, err := h.Service.SelectWalletByUUID(walletID)
	if err != nil || wallet == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "The wallet was not found",
		})
		return
	}

	// Adding money
	wallet.Money += money

	// Updating wallet by UUID
	_, err = h.Service.UpdateWalletByUUID(*wallet)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error occurred while updating the wallet. Try again later",
		})
		return
	}

	// Create a new transaction
	transaction := Wallets.TransactionModel{
		Sender:    wallet.ID,
		Recipient: wallet.ID,
		Money:     money,
	}
	_, err = h.Service.CreateTransaction(transaction)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error occurred while creating transaction. Try again later",
		})
		return
	}

	// Returning wallet information
	c.JSON(200, gin.H{
		"status": "success",
		"data":   wallet,
	})
}

func (h WalletHandlers) TransferWalletByUUID(c *gin.Context) {
	// Getting the current user and sender ID, recipient ID and amount of money
	u, _ := c.Get("authUser")
	sender, recipient, moneyString := c.Param("uuid"), c.PostForm("recipient"), c.PostForm("money")

	// Parsing a number from a string of amount of money
	money, err := strconv.Atoi(moneyString)
	if err != nil || money <= 0 {
		c.JSON(406, gin.H{
			"status":  "error",
			"message": "Incorrect amount of money to send",
		})
		return
	}

	// Parsing the sender ID and recipient ID in a UUID
	senderID, err := uuid.Parse(sender)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get the sender`s wallet ID. Try again later",
		})
		return
	}
	recipientID, err := uuid.Parse(recipient)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get the recipient`s wallet ID. Try again later",
		})
		return
	}

	// Selecting wallets
	senderWallet, err := h.Service.SelectWalletByUUID(senderID)
	if err != nil || senderWallet == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "The sender's wallet was not found",
		})
		return
	}
	recipientWallet, err := h.Service.SelectWalletByUUID(recipientID)
	if err != nil || recipientWallet == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "The recipient's wallet was not found",
		})
		return
	}

	// Check if the sender's wallet is not equal to the recipient's wallet
	if senderWallet.ID == recipientWallet.ID {
		c.JSON(409, gin.H{
			"status":  "error",
			"message": "The sender's wallet and the recipient's wallet are one and the same",
		})
		return
	}

	// Checking user access permissions to the sender's wallet
	if senderWallet.UserID != u.(Users.UserModel).ID {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": "You do not have access to this wallet",
		})
		return
	}

	// Checking the currency of the recipient's wallet and the sender's wallet
	if senderWallet.Currency != recipientWallet.Currency {
		c.JSON(409, gin.H{
			"status":  "error",
			"message": "The currency of the sender's wallet does not match the currency of the recipient's wallet",
		})
		return
	}

	// Checking the amount of money in the sender's wallet
	if senderWallet.Money < money {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": "Not enough money to send",
		})
		return
	}

	// Taking money out of wallets
	senderWallet.Money -= money
	recipientWallet.Money += money

	// Updating wallets
	_, err = h.Service.UpdateWalletByUUID(*senderWallet)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error occurred while updating the sender's wallet. Try again later",
		})
		return
	}
	_, err = h.Service.UpdateWalletByUUID(*recipientWallet)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error occurred while updating the recipient's wallet. Try again later",
		})
		return
	}

	// Create a new transaction
	transaction := Wallets.TransactionModel{
		Sender:    senderWallet.ID,
		Recipient: recipientWallet.ID,
		Money:     money,
	}
	_, err = h.Service.CreateTransaction(transaction)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error occurred while creating transaction. Try again later",
		})
		return
	}

	// Return the balance of the sender's and recipient's wallets
	c.JSON(200, gin.H{
		"status":          "success",
		"sender_money":    senderWallet.Money,
		"recipient_money": recipientWallet.Money,
	})
}

func (h WalletHandlers) DeleteWalletByUUID(c *gin.Context) {
	// Getting the current user and parsing wallet ID
	u, _ := c.Get("authUser")
	walletID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Failed to get the wallet ID. Try again later",
		})
		return
	}

	// Selecting wallet
	wallet, err := h.Service.SelectWalletByUUID(walletID)
	if err != nil || wallet == nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "The wallet was not found",
		})
		return
	}

	// Checking user access permissions to the sender's wallet
	if wallet.UserID != u.(Users.UserModel).ID {
		c.JSON(403, gin.H{
			"status":  "error",
			"message": "You do not have access to this wallet",
		})
		return
	}

	// Deleting wallet by ID
	err = h.Service.DeleteWalletByUUID(walletID)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "An unexpected error occurred while deleting the wallet. Try again later",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"id":     wallet.ID,
	})
}
