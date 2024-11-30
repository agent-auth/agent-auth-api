package tokendal

import (
	"github.com/agent-auth/agent-auth-api/database/dbmodels"
)

type TokenDal interface {
	Create(txID string, token *dbmodels.Token) (*dbmodels.Token, error)
	GetByUUID(uuid string) (*dbmodels.Token, error)
	Update(token *dbmodels.Token) error
	DeleteByAccessToken(token string) error
}
