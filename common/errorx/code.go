package errorx

const (
	Success int64 = 0

	InvalidParam int64 = 10001
	Unauthorized int64 = 10002
	Forbidden    int64 = 10003

	UserNotFound        int64 = 20001
	UserInvalidPassword int64 = 20002
	UserDisabled        int64 = 20003
	UserAlreadyExists   int64 = 20004

	ProductNotFound int64 = 30001
	ProductOffShelf int64 = 30002
	StockNotEnough  int64 = 30003

	OrderNotFound      int64 = 40001
	OrderStatusInvalid int64 = 40002
	OrderCreateFailed  int64 = 40003

	PaymentNotFound       int64 = 50001
	PaymentAlreadyExists  int64 = 50002
	PaymentAlreadySuccess int64 = 50003
	PaymentStatusInvalid  int64 = 50004

	ServerError       int64 = 90001
	RpcError          int64 = 90002
	DownstreamFailure int64 = 90003
)
