package config

const (
	// Test and Token Validity
	CTO_USER                 = "community_test"
	CEO_USER                 = "ceo_test"
	DefaultTokenValidityTime = "15"
)

const (
	// Email Service Notification
	ResetPasswordTokenStringURI     = "/api/v1/reset"
	EmailVerificationTokenStringURI = "/api/v1/verify"
	EmailSubjectLine                = "Verify your email address for FleetWise Application"
	EmailTextString                 = "Click on the following link to verify your email address:"
)

const (
	// Errors
	ErrorTokenValidation        = "unable to validate token"
	ErrorTokenSubjectValidation = "unable to validate token subject"
	ErrorFetchingCurrentUser    = "unable to retrieve system user"
	ErrorUserIsAlreadyVerified  = "unable to validate user. user is already verified"
)
