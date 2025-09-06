package response

type SuccessListResponse[T any] struct {
	BaseStatusErrorResponse
	Data []T `json:"data"`
}

func NewSuccessListResponse[T any](items []T) SuccessListResponse[T] {
	return SuccessListResponse[T]{
		BaseStatusErrorResponse: BaseStatusErrorResponse{
			Status: "ok",
			Error:  nil,
		},
		Data: items,
	}
}

type BaseStatusErrorResponse struct {
	Status string  `json:"status"`
	Error  *string `json:"error"`
} // @name BaseResponse

type BaseResponseData[T any] struct {
	BaseStatusErrorResponse
	Data *T `json:"data"`
} // @name BaseResponseData

func NewSuccessResponse[T any](data T) BaseResponseData[T] {
	return BaseResponseData[T]{
		BaseStatusErrorResponse: BaseStatusErrorResponse{
			Status: "ok",
			Error:  nil,
		},
		Data: &data,
	}
}

func NewErrorResponse(message string) BaseStatusErrorResponse {
	return BaseStatusErrorResponse{
		Status: "error",
		Error:  &message,
	}
}
