package internal

type ClientGateway interface {
	GetAccessToken(credentials Credentials) (string, error)
	CreatePreference(preference NewPreference) (string, error)
}

type Controller struct {
	Client ClientGateway
}

func NewController(client ClientGateway) *Controller {
	return &Controller{
		Client: client,
	}
}

func (s *Controller) GetAccessToken(clientID string, clientSecret string) (string, error) {
	return s.Client.GetAccessToken(Credentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
}

func (s *Controller) CreatePreference(preference NewPreference) (string, error) {
	return s.Client.CreatePreference(preference)
}