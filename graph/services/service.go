package services

type Services interface {}

type services struct {}

func New() Services {
	return &services{}
}