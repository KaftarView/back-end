package controller

import (
	application_payment "first-project/src/application/payment"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService *application_payment.PaymentService
}

func NewPaymentController(paymentService *application_payment.PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (pc *PaymentController) ZarinPayTest(c *gin.Context) {
	res, err := pc.paymentService.ZarinPay(100000, "http://localhost:8080/v1/events/event-details/14", "Dr", "alos@gmail.com")
	if err != nil {
		Response(c, 400, "Fail", nil)
	} else {
		Response(c, http.StatusOK, "success", res)
	}

}

func (pc *PaymentController) Verify(c *gin.Context) {
	Response(c, http.StatusOK, "success", nil)
}
