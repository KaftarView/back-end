package application_payment

import (
	"context"
	"first-project/src/application/payment/zarinpal"
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
)

type PaymentService struct {
	PayInfo *bootstrap.PayInfo
}

func NewPaymentService(PayInfo *bootstrap.PayInfo) *PaymentService {
	return &PaymentService{
		PayInfo: PayInfo,
	}
}

func (ps *PaymentService) ZarinPay(amount uint, callBackUrl string, description string, email string) (string, error) {
	var merch = ps.PayInfo.ZarinMerchantID
	service, err := zarinpal.NewService(merch)
	if err != nil {
		return "", exceptions.PaymentServerError{Message: "Failed to initialize payment service"}
	}

	request := &zarinpal.PaymentRequestDto{
		Amount:      int(amount),
		Description: description,
		Email:       email,
		Mobile:      "09120000000",
		Currency:    "IRR",
		CallbackURL: callBackUrl,
	}

	ctx := context.Background()
	paymentURL, err := service.Request(ctx, request)
	if err != nil {
		return "", exceptions.PaymentError{Message: "Unknown error occurred during payment", Code: 500}
	}

	return paymentURL, nil
}
