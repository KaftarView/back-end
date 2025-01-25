package errors

import (
	"errors"
	"fmt"
)

var zarinpalErrorMessages = map[int]string{
	-1:   "اطلاعات ارسال شده ناقص است",
	-2:   "مرچنت کد یا آدرس IP پذيرنده صحيح نيست",
	-3:   "مبلغ تراکنش باید بیشتر از 100 تومان باشد",
	-4:   "سطح تایید پذیرنده پایین تر از سطح نقره ای است",
	-9:   "پارامترهای ورودی نامعتبر است",
	-10:  "ای پی پذیرنده صحیح نیست",
	-11:  "درخواست مورد نظر یافت نشد",
	-12:  "امکان ویرایش درخواست میسر نمی باشد",
	-15:  "زمان ثبت درخواست به پایان رسیده است",
	-21:  "هیچ نوع عملیات مالی برای این تراکنش یافت نشد",
	-22:  "تراکنش ناموفق است",
	-30:  "اجازه دسترسی به این متد وجود ندارد",
	-31:  "حساب بانکی متصل به زرین پال تایید نشده است",
	-32:  "امکان انصراف از درخواست میسر نمی باشد",
	-33:  "رقم تراکنش با رقم پرداخت شده مطابقت ندارد",
	-34:  "سقف تقسیم تراکنش از لحاظ تعداد یا رقم عبور نموده است",
	-35:  "تاریخ ارسال درخواست باید بزرگتر از تاریخ تراکنش باشد",
	-40:  "اجازه دسترسی به متد مربوطه وجود ندارد",
	-41:  "اطلاعات ارسال شده مربوط به AdditionalData نامعتبر است",
	-42:  "طول عمر شناسه پرداخت باید بین مدت زمان ۳۰ دقیقه تا ۴۵ روز باشد",
	-50:  "مبلغ پرداخت شده با مقدار درخواستی مطابقت ندارد",
	-51:  "پرداخت ناموفق بوده و یا توسط کاربر لغو شده است",
	-54:  "درخواست مورد نظر آرشیو شده است",
	-55:  "زمان درخواست تکراری است",
	-101: "عملیات پرداخت موفق بوده ولی PaymentVerification قبلا انجام شده است",
	-102: "تراکنش بسته شده است",
	-103: "تراکنش ناموفق به علت عدم تایید اطلاعات",
	-104: "شناسه درخواست اشتباه است",
	-105: "عدم تایید اطلاعات بانکی",
	-106: "پرداخت تایید نشده است",
	-107: "شناسه پرداخت نامعتبر است",
	101:  "عملیات پرداخت موفق بوده ولی PaymentVerification قبلا انجام شده است",
}

func GetZarinpalError(code int, message interface{}) error {
	if msg, exists := zarinpalErrorMessages[code]; exists {
		return errors.New(msg + fmt.Sprintf(" - راهنمایی: %v", message))
	}

	return errors.New("خطای ناشناخته - " + fmt.Sprintf("راهنمایی: %v", message))
}
