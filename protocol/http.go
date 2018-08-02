package protocol

type AddProductionRequest struct {
	Type       string `json:"type"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	ValueDate  int64  `json:"valueDate"`
	DueDate    int64  `json:"dueDate"`
	AnnualRate int64  `json:"annualRate"`
}

type AddProductionResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type Production struct {
	Type       string
	Code       string
	Name       string
	ValueDate  int64
	DueDate    int64
	AnnualRate int64
}

type GetAllProductionRequest struct {
}

type GetAllProductionResponse struct {
	Code       int64        `json:"code"`
	Production []Production `json:"production"`
}

type ErrorResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}
