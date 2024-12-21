package application_payment

import (
	"context"
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
	"first-project/src/pkg/zarinpal"
)

type PaymentService struct {
	PayInfo *bootstrap.PayInfo
}

func NewPaymentService(PayInfo *bootstrap.PayInfo) *PaymentService {
	return &PaymentService{
		PayInfo: PayInfo,
	}
}

func (ps *PaymentService) ZarinPay(amount uint, callBackUrl string, description string, email string, mobile string) string {
	var merch = ps.PayInfo.ZarinMerchantID
	service, err := zarinpal.NewService(merch)
	if err != nil {
		var PaymentServerError = exceptions.NewPaymentServerError()
		panic(PaymentServerError)
	}

	request := &zarinpal.PaymentRequestDto{
		Amount:      int(amount),
		Description: description,
		Email:       email,
		Mobile:      mobile,
		Currency:    "IRR",
		CallbackURL: callBackUrl,
	}

	ctx := context.Background()
	paymentURL, err := service.Request(ctx, request)
	if err != nil {
		var PaymentError = exceptions.NewPaymentError()
		panic(PaymentError)
	}

	return paymentURL
}
