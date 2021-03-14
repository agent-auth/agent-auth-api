package authservice

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/jobbox-tech/recruiter-api/models/authmodel"
	"github.com/jobbox-tech/recruiter-api/web/interfaces/v1/authinterface"
	"github.com/jobbox-tech/recruiter-api/web/renderers"
	"github.com/mssola/user_agent"
	"github.com/spf13/viper"
)

func (as *authservice) Authenticate(w http.ResponseWriter, r *http.Request) {
	txID := r.Header["transaction_id"][0]
	body := &authinterface.AuthenticateReqInterface{}
	if err := render.Bind(r, body); err != nil {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Failed to read the request body with error %v", err)
		render.Render(w, r, renderers.ErrorUnauthorized(authmodel.ErrLoginToken))
		return
	}

	// reterive token from database
	token, err := as.tokenDal.GetByUUID(body.Token)
	if err != nil {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Failed to reterive token from database with error %v", err)
		render.Render(w, r, renderers.ErrorUnauthorized(authmodel.ErrLoginToken))
		return
	}

	// reterive associated account with token
	acc, err := as.recruiterDal.GetByID(token.AccountID)
	if err != nil {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Failed to recruiter from database with error %v", err)
		render.Render(w, r, renderers.ErrorUnauthorized(authmodel.ErrUnknownLogin))
		return
	}

	// check if token is not expired or already claimed
	if time.Now().UTC().After(token.ExpiryTimestampUTC) || token.Claimed {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Access token is expired", err)
		render.Render(w, r, renderers.ErrorUnauthorized(authmodel.ErrLoginToken))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, renderers.ErrorUnauthorized(authmodel.ErrLoginDisabled))
		return
	}

	access, refresh, err := as.tokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Failed to generate token with error %v", err)
		render.Render(w, r, renderers.ErrorInternalServerError(""))
		return
	}

	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	token.AccessToken = access
	token.ResfreshToken = refresh
	token.Claimed = true
	token.ExpiryTimestampUTC = time.Now().UTC().Add(viper.GetDuration("jwt.auth_jwt_expiry"))
	token.UpdatedTimestampUTC = time.Now().UTC()
	token.UserAgent = fmt.Sprintf("%s on %s", browser, ua.OS())
	token.Mobile = ua.Mobile()

	if err := as.tokenDal.Update(token); err != nil {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Failed to update token details with error %v", err)
		render.Render(w, r, renderers.ErrorInternalServerError(""))
		return
	}

	acc.LastLogin = time.Now().UTC()
	if err := as.recruiterDal.Update(acc); err != nil {
		as.logger.Error(txID, authmodel.FailedToAuthenticateToken).Errorf("Failed to update recruiter details with error %v", err)
	}

	render.Respond(w, r, &authinterface.AuthenticateResInterface{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
