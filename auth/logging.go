package auth

import (
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   AuthService
}

func (mw loggingMiddleware) Close(logger log.Logger) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Close (Shutting down auth service)",
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.Close(logger)
	return
}

func serializeRolesArray(roles []string) string {
	var output string
	for _, role := range roles {
		output += role + ","
	}
	return output
}

func (mw loggingMiddleware) Login(username string, passwordHash string, roles []string) (output string, expiresInSec int, userID string, firstName string, lastName string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Login",
			"username", username,
			"roles", serializeRolesArray(roles),
			"outputTokenExpiresInSec", expiresInSec,
			"outputToken", output,
			"userID", userID,
			"firstName", firstName,
			"lastName", lastName,
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	output, expiresInSec, userID, firstName, lastName, err = mw.next.Login(username, passwordHash, roles)

	return
}

func (mw loggingMiddleware) Register(username string, passwordHash string, firstName string, lastName string, email string, roles []string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Register",
			"username", username,
			"email", email,
			"firstName", firstName,
			"lastName", lastName,
			"roles", serializeRolesArray(roles),
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.Register(username, passwordHash, firstName, lastName, email, roles)
	return
}

func (mw loggingMiddleware) RefreshToken(token string) (refreshedToken string, expiresInSec int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "RefreshToken",
			"newTokenExpiresInSec", expiresInSec,
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	refreshedToken, expiresInSec, err = mw.next.RefreshToken(token)
	return
}

func (mw loggingMiddleware) ValidateToken(token string) (claims AuthClaims, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ValidateToken",
			"err", errorToString(err),
			"claims:userID", claims.UserID,
			"claims:username", claims.Username,
			"claims:roles", claims.Roles,
			"claims:exp", claims.ExpiresAt.Unix(),
			"took", time.Since(begin),
		)
	}(time.Now())
	claims, err = mw.next.ValidateToken(token)
	return
}

func (mw loggingMiddleware) GetUserProfile(userID string) (profile *UserProfile, err error) {
	defer func(begin time.Time) {
		var printedPhoneNumber string
		var printedOrganizationName string
		var printedAddress string

		if profile == nil {
			printedPhoneNumber = "(none)"
			printedOrganizationName = "(none)"
			printedAddress = "(none)"
		} else {
			if profile.PhoneNumber == "" {
				printedPhoneNumber = "(none)"
			} else {
				printedPhoneNumber = profile.PhoneNumber
			}
			if profile.OrganizationName == "" {
				printedOrganizationName = "(none)"
			} else {
				printedOrganizationName = profile.OrganizationName
			}
			printedAddress = serializeAddress(profile.ShippingAddress)
		}
		mw.logger.Log(
			"method", "GetUserProfile",
			"userID", userID,
			"profile:organizationName", printedOrganizationName,
			"profile:phoneNumber", printedPhoneNumber,
			"profile:shippingAddress", printedAddress,
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	profile, err = mw.next.GetUserProfile(userID)
	return
}

func (mw loggingMiddleware) UpdateUserProfile(userID string, profile *UserProfile) (err error) {
	defer func(begin time.Time) {
		var printedPhoneNumber string
		var printedOrganizationName string
		if profile.PhoneNumber == "" {
			printedPhoneNumber = "(none)"
		} else {
			printedPhoneNumber = profile.PhoneNumber
		}
		if profile.OrganizationName == "" {
			printedOrganizationName = "(none)"
		} else {
			printedOrganizationName = profile.OrganizationName
		}
		mw.logger.Log(
			"method", "UpdateUserProfile",
			"userID", userID,
			"profile:organizationName", printedOrganizationName,
			"profile:phoneNumber", printedPhoneNumber,
			"profile:shippingAddress", serializeAddress(profile.ShippingAddress),
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.UpdateUserProfile(userID, profile)
	return
}

func (mw loggingMiddleware) ChangePassword(userID string, oldPassword string, newPassword string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ChangePassword",
			"userID", userID,
			"oldPassword is non-empty", oldPassword != "",
			"newPassword is non-empty", newPassword != "",
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.ChangePassword(userID, oldPassword, newPassword)
	return
}

func (mw loggingMiddleware) ResetPasswordInitiate(email string) (token string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ResetPasswordInitiate",
			"email", email,
			"token is non-empty", token != "",
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	token, err = mw.next.ResetPasswordInitiate(email)
	return
}

func (mw loggingMiddleware) ResetPasswordComplete(resetToken string, username string, newPassword string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ResetPasswordComplete",
			"resetToken is non-empty", resetToken != "",
			"username", username,
			"newPassword is non-empty", newPassword != "",
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.ResetPasswordComplete(resetToken, username, newPassword)
	return
}

func (mw loggingMiddleware) GetUserByUsername(username string) (user *User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUserByUsername",
			"username", username,
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	user, err = mw.next.GetUserByUsername(username)
	return
}

func (mw loggingMiddleware) GetUserByID(userID string) (user *User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUserByID",
			"userID", userID,
			"err", errorToString(err),
			"took", time.Since(begin),
		)
	}(time.Now())
	user, err = mw.next.GetUserByID(userID)
	return
}

func NewLoggingMiddleware(logger log.Logger, s AuthService) AuthService {
	return &loggingMiddleware{logger, s}
}