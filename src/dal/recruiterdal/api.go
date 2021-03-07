package recruiterdal

import (
	"github.com/jobbox-tech/recruiter-api/database/dbmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecruiterDal ...
type RecruiterDal interface {
	Create(txID string, account *dbmodels.Recruiter) (primitive.ObjectID, error)
	GetAccountByEmail(email string) (*dbmodels.Recruiter, error)
}
