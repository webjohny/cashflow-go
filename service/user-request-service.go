package service

import (
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"gopkg.in/errgo.v2/errors"
	"log"
)

type UserRequestService interface {
	GetOneById(userRequestId uint64) entity.UserRequest
	Update(userRequest entity.UserRequest) entity.UserRequest
	HandleUserRequest(data dto.HandleUserRequestBodyDto) (error, entity.UserRequest)
	GetAllByRaceId(raceId uint64) []entity.UserRequest
}

type userRequestService struct {
	userRequestRepository repository.UserRequestRepository
}

func NewUserRequestService(userReqRepo repository.UserRequestRepository) UserRequestService {
	return &userRequestService{
		userRequestRepository: userReqRepo,
	}
}

func (service *userRequestService) Update(userRequest entity.UserRequest) entity.UserRequest {
	err, updatedUserRequest := service.userRequestRepository.Update(&userRequest)
	if err != nil {
		log.Fatalf("Failed to map: %v", err)
	}
	return updatedUserRequest
}

func (service *userRequestService) GetOneById(userRequestId uint64) entity.UserRequest {
	return service.userRequestRepository.FindOneById(userRequestId)
}

func (service *userRequestService) GetAllByRaceId(raceId uint64) []entity.UserRequest {
	requests := service.userRequestRepository.All("race_id = ? AND status = 0", []interface{}{raceId})
	return requests
}

func (service *userRequestService) HandleUserRequest(data dto.HandleUserRequestBodyDto) (error, entity.UserRequest) {
	var err error
	userRequest := service.GetOneById(data.UserRequestId)

	if userRequest.ID > 0 && userRequest.Status == 0 {
		userRequest.Status = data.Status
		userRequest.RejectMessage = data.Message
		return nil, service.Update(userRequest)
	} else if userRequest.Status > 0 {
		err = errors.New(storage.ErrorUserRequestHasBeenAlreadyApproved)
	} else {
		err = errors.New(storage.ErrorUndefinedUserRequest)
	}

	return err, entity.UserRequest{}
}
