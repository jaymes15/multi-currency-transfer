package transfers

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	requests "lemfi/simplebank/internal/apps/transfers/requests"
	transferValidation "lemfi/simplebank/internal/apps/transfers/validationMessages"

	"lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/requestHandler"
	"lemfi/simplebank/pkg/responseHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (transferController *TransferController) MakeTransferController(c *gin.Context) {
	config.Logger.Info("Making transfer", "method", "POST", "endpoint", "/transfers")

	var req requests.MakeTransferRequest

	err := requestHandler.ReadJSONGin(c, &req, transferValidation.MakeTransferValidationMessages)
	if err != nil {
		config.Logger.Error("Failed to read transfer request", "error", err.Error())
		errorResponse.BadRequestResponse(c, err)
		return
	}

	config.Logger.Info("Transfer request validated successfully", "fromAccountID", req.FromAccountID, "toAccountID", req.ToAccountID, "amount", req.Amount, "fromCurrency", req.FromCurrency, "toCurrency", req.ToCurrency)

	transfer, err := transferController.transferService.MakeTransfer(req)
	if err != nil {
		config.Logger.Error("Failed to make transfer", "error", err.Error(), "fromAccountID", req.FromAccountID, "toAccountID", req.ToAccountID, "amount", req.Amount, "fromCurrency", req.FromCurrency, "toCurrency", req.ToCurrency)

		// Check if it's a client error (400) or server error (500)
		if clientErr, isClient := core.IsClientError(err); isClient {
			errorResponse.BadRequestResponse(c, clientErr)
		} else {
			errorResponse.ServerErrorResponse(c, err)
		}
		return
	}

	config.Logger.Info("Transfer made successfully",
		"transferID", transfer.Transfer.ID,
		"fromAccountID", transfer.Transfer.FromAccountID,
		"toAccountID", transfer.Transfer.ToAccountID,
		"amount", transfer.Transfer.Amount,
		"fromCurrency", transfer.Transfer.FromCurrency,
		"toCurrency", transfer.Transfer.ToCurrency,
	)

	response := responseHandler.Envelope{
		"transfer": transfer,
	}

	err = responseHandler.WriteJSON(c.Writer, http.StatusCreated, response, nil)
	if err != nil {
		config.Logger.Error("Failed to write JSON response", "error", err.Error())
		errorResponse.ServerErrorResponse(c, err)
		return
	}

	config.Logger.Info("Transfer response written successfully",
		"transferID", transfer.Transfer.ID,
		"fromAccountID", transfer.Transfer.FromAccountID,
		"toAccountID", transfer.Transfer.ToAccountID,
		"amount", transfer.Transfer.Amount,
		"fromCurrency", transfer.Transfer.FromCurrency,
		"toCurrency", transfer.Transfer.ToCurrency,
	)
}
