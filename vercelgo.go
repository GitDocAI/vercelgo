package vercelgo

import "sync"

var (
	vercelClient     *VercelClient
	vercelClientOnce sync.Once
)

type VercelClient struct {
	Token string
}

func GetVercelClient(token string) *VercelClient {
	vercelClientOnce.Do(func() {
		vercelClient = &VercelClient{
			Token: token,
		}
	})
	return vercelClient
}

func (c *VercelClient) GetHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.Token,
		"Content-Type":  "application/json",
	}
}
