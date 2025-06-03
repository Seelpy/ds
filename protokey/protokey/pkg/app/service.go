package app

import (
	"errors"
	"regexp"
)

type Service interface {
	Set(key string, value int) error
	Get(key string) (int, error)
	Keys(prefix string) ([]string, error)
}

type ProtoKeyService struct {
	commands  chan Command
	responses chan Response
}

var (
	keyPattern           = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{1,1000}$`)
	minInt32             = -2147483648
	maxInt32             = 2147483647
	ErrBadRequest        = errors.New("bad request")
	ErrCantParseResponse = errors.New("cant parse response")
)

func NewProtoKeyService(cmdCh chan Command, respCh chan Response) *ProtoKeyService {
	return &ProtoKeyService{commands: cmdCh, responses: respCh}
}

func (service *ProtoKeyService) Set(key string, value int) error {
	if !keyPattern.MatchString(key) {
		return ErrBadRequest
	}
	err := service.validateValue(value)
	if err != nil {
		return err
	}

	service.commands <- &SetCommand{Key: key, Value: value}
	rawResp := <-service.responses

	resp, ok := rawResp.(*SetResponse)
	if !ok {
		return ErrCantParseResponse
	}

	return resp.Err
}

func (service *ProtoKeyService) Get(key string) (int, error) {
	if !keyPattern.MatchString(key) {
		return 0, ErrBadRequest
	}

	service.commands <- &GetCommand{Key: key}
	rawResp := <-service.responses

	resp, ok := rawResp.(*GetResponse)
	if !ok {
		return 0, ErrCantParseResponse
	}

	return resp.Value, resp.Err
}

func (service *ProtoKeyService) Keys(prefix string) ([]string, error) {
	if !keyPattern.MatchString(prefix) {
		return nil, ErrBadRequest
	}

	service.commands <- &ListKeysCommand{Prefix: prefix}
	rawResp := <-service.responses

	resp, ok := rawResp.(*ListKeysResponse)
	if !ok {
		return nil, ErrCantParseResponse
	}

	return resp.Keys, resp.Err
}

func (service *ProtoKeyService) validateValue(value int) error {
	if value < minInt32 || value > maxInt32 {
		return ErrBadRequest
	}
	return nil
}
