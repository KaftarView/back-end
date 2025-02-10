package localization

var Persian = map[string]interface{}{
	"errors": map[string]interface{}{
		"generic":                 "خطایی رخ داده است، لطفا دوباره تلاش کنید.",
		"numeric":                 "`{0}` باید عدد باشد.",
		"alreadyExist":            "این `{0}` قبلا ثبت شده است.",
		"minimumLength":           "باید از 7 کاراکتر بیشتر باشد.",
		"containsLowercase":       "باید حتما دارای حرف کوچک باشد.",
		"containsUppercase":       "باید حتما دارای حرف بزرگ باشد.",
		"containsNumber":          "باید حتما دارای عدد باشد.",
		"containsSpecialChar":     "باید حتما دارای حروف خاص باشد.",
		"notMatchConfirmPAssword": "رمزعیور با تکرار رمزعبور باید یکسان باشد.",
		"alreadyVerified":         "کاربر قبلا احراز هویت شده است.",
		"expiredToken":            "توکن شما منقضی شده است.",
		"invalidToken":            "توکن شما اشتباه است.",
		"loginFailed":             "ورود ناموفق بود. لطفاً نام کاربری و رمز عبور خود را بررسی کنید و اطمینان حاصل کنید که حساب شما فعال شده است.",
		"emailNotExist":           "ابتدا باید ثبت نام کنید.",
		"rateLimitExceed":         "شما از حد مجاز درخواستها فراتر رفته اید.",
		"NewsNotFound":            "خبر پیدا نشد",
		"forbidden":               "شما به این بخش دسترسی ندارید.",
		"unauthorized":            "مشکلی برای احراز هویت شما پیش آمده است. مجدد تلاش کنید.",
		"fileRequired":            "وارد کردن فایل اجباری است.",
		"locationAlreadyTaken":    "این محل در این زمان در حال برگزاری رویداد دیگری است.",
		"notFoundError":           "`{0}` یافت نشد.",
		"alreadySubscribed":       "`{0}` قبلا دنبال شده است.",
		"notSubscribed":           "`{0}` هنوز دنبال نشده است.",
		"notAvailable":            "`{0}` در حال حاضر در دسترس نیست.",
		"purchaseFailed":          "خرید با شکست همراه شد. لطغا مجدد تلاش کنید.",
		"AlreadyPurchased":        "بلیت قبلا خریداری شده است.",
	},
	"successMessage": map[string]interface{}{
		"userRegistration":       "با موفقیت انجام شد. لطفا ایمیل خود را تایید کنید تا حساب شما فعال شود.",
		"NewsCreation":           "خبر با موفقیت ایجاد شد.",
		"NewsFound":              "خبر یافت شد",
		"emailVerification":      "ایمیل شما تایید شد.",
		"login":                  "با موفقیت وارد شدید.",
		"forgotPassword":         "لطفاً ایمیل خود را بررسی کنید تا رمزعبور یکبار مصرف خود را وارد کنید.",
		"resetPassword":          "رمزعبور شما با موفقیت تغییر کرد.",
		"changeUsername":         "نام کاربری با موفقیت تغییر کرد.",
		"refreshToken":           "نشسست شما تمدید شد.",
		"uploadObjectToBucket":   "فایل مورد نظر شما با موفقیت به صندوقچه اضافه شد.",
		"deleteObjectFromBucket": "فایل مورد نظر شما با موفقیت از صندوقچه حذف شد.",
		"updateUserRole":         "نقش مورد نظر با موفقیت به کاربر اضافه شد.",
		"createRole":             "نقش با موفقیت ساخته شد.",
		"updateRole":             "نقش با موفقیت به روز شد.",
		"deleteRole":             "نقش با موفقیت حذف شد.",
		"deleteRolePermission":   "دسترسی نقش مورد نظر با موفقیت محدود شد.",
		"deleteUserRole":         "نقش برای فرد مورد نظر حذف شد.",
		"NewsUpdated":            "خبر با موفقیت ویرایش شد.",
		"NewsDeleted":            "خبر با موفقیت حذف شد.",
		"createEvent":            "رویداد با موفقیت ساخته شد.",
		"addTicket":              "بلیت با موفقیت اضافه شد.",
		"addDiscount":            "تخفیف با موفقیت اضافه شد.",
		"updateTicket":           "بلیت با موفقیت بروزرسانی شد",
		"updateDiscount":         "تخفیف با موفقیت بروزرسانی شد",
		"getEvent":               "رویداد با موفقیت یافت شد",
		"updateEvent":            "رویداد با موفقیت بروزرسانی شد",
		"deleteEvent":            "رویداد با موفقیت حذف شد.",
		"uploadMedia":            "منابع مورد نظر با موفقیت اضافه شد.",
		"deleteMedia":            "منابع مورد نظر با موفقیت حذف شد.",
		"updateMedia":            "منابع مورد نظر با موفقیت به روز شد.",
		"publishEvent":           "رویداد با موفقیت منتشر شد.",
		"unpublishEvent":         "رویداد با موفقیت از لیست انتشار حذف شد.",
		"deleteTicket":           "بلیت با موفقیت حذف شد.",
		"deleteDiscount":         "کدتخفیف با موفقیت حذف شد.",
		"organizerRegistration":  "با موفقیت انجام شد. پس از تایید ایمیل توسط برگزارکننده، عملیات شما فعال خواهد شد.",
		"organizerActivated":     "اطلاعات شما با موفقیت به عنوان برگزارکننده رویداد ثبت شد.",
		"addComment":             "نظر شما با موفقیت ثبت شد.",
		"editComment":            "نظر شما با موفقیت ویرایش شد.",
		"deleteComment":          "نظر شما با موفقیت حذف شد.",
		"deleteOrganizer":        "برگزارکننده با موفقیت حذف شد.",
		"createPodcast":          "پادکست با موفقیت ساخته شد.",
		"updatePodcast":          "پادکست با موفقیت به روزرسانی شد.",
		"createPodcastEpisode":   "قسمت پادکست مورد نظر شما با موفقیت ساخته شد.",
		"updatePodcastEpisode":   "قسمت پادکست مورد نظر شما با موفقیت به روزرسانی شد.",
		"deletePodcastEpisode":   "قسمت پادکست مورد نظر با موفقیت حذف شد.",
		"subscribePodcast":       "با موفقیت به دنبال کنندگان پادکست اضافه شدید.",
		"unSubscribePodcast":     "شما از لیست دنبال کنندگان پادکست خارج شدید.",
		"deletePodcast":          "پادکست با موفقیت حذف شد.",
		"createJournal":          "نشریه با موفقیت ساخته شد.",
		"updateJournal":          "نشریه با موفقیت به روزرسانی شد.",
		"deleteJournal":          "نشریه با موفقیت حذف شد.",
		"createCouncilor":        "عضو انجمن با موفقیت ساخته شد.",
		"reserveTicket":          "بلیت های شما با موفقیت رزرو شدند.",
		"purchaseTicket":         "بلیت های شما با موفقیت خریداری شدند.",
	},
}
