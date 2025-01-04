package enums

type TransactionStatus uint

const (
	Success TransactionStatus = iota + 1
	Failed
)

func (status TransactionStatus) String() string {
	switch status {
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	}
	return ""
}

func GetAllTransactionStatus() []TransactionStatus {
	return []TransactionStatus{
		Success,
		Failed,
	}
}
