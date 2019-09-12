package https

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/velo-labs/cen/node/app/modules/node"
)

type nodeHandler struct {
	NodeUsecase node.UseCase
}

func NewEndpointHttpHandler(ginEngine *gin.Engine, nodeUC node.UseCase) {
	handler := &nodeHandler{
		NodeUsecase: nodeUC,
	}

	v1 := ginEngine.Group("/v1")
	{
		v1.POST("/credits.setup", handler.SetupCreditHandler)
	}
}
