package service

import "github.com/webjohny/cashflow-go/repository"

type LobbyService interface {
	CreateLobby(username string) error
	Join(ID uint64, username string) error
}

type lobbyService struct {
	lobbyRepository repository.LobbyRepository
}

func NewLobbyService(lobbyRepository repository.LobbyRepository) LobbyService {
	return &lobbyService{
		lobbyRepository: lobbyRepository,
	}
}

func (service *lobbyService) CreateLobby(username string) error {
	return nil
}

func (service *lobbyService) Join(ID uint64, username string) error {
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby != nil {
		lobby.AddPlayer(username)
	}

	return nil
}
