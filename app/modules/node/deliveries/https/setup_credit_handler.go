package https

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/velo-labs/cen/app/modules/node/deliveries/models"
	"gitlab.com/velo-labs/cen/app/utils"
	"net/http"
)

func (h *nodeHandler) SetupCreditHandler(c *gin.Context) {
	var setupCreditRequest models.SetupCreditRequest
	err := c.BindJSON(&setupCreditRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(err.Error()))
		return
	}

	setupCreditEntity, err := h.NodeUsecase.Setup(
		setupCreditRequest.SignedIssuerCreationTx,
		setupCreditRequest.PeggedValue,
		setupCreditRequest.PeggedCurrency,
		setupCreditRequest.AssetName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err.Error()))
		return
	}

	resp, err := new(models.SetupCreditResponse).Parse(&setupCreditEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(resp))
}
