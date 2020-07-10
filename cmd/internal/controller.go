package internal

type ClientGateway interface {
	GetAccessToken(credentials Credentials) (string, error)
	CreatePreference(accessToken string, preference NewPreference) (string, error)
	GetTotalPayments(accessToken string, status string) (int, error)
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

func (s *Controller) CreatePreference(accessToken string, preference NewPreference) (string, error) {
	return s.Client.CreatePreference(accessToken, preference)
}

func (s *Controller) GetTotalPayments(accessToken string, status string) (int, error) {
	return s.Client.GetTotalPayments(accessToken, status)
}