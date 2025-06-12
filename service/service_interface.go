package service

type IService interface {
	ProcessRequest(any) ([]byte, error)
	DecodeAndValidate([]byte, string, string) (any, error)
}
