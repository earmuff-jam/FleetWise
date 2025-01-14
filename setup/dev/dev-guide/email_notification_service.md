# Email Notification Service

To help us reach our audience during registration process, we use `SendGrid` Email Servers to help us communicate with our clients. Sending our email via a third party server requires us
to have a dedicated API key for SendGrid Server. This api key is introduced in the `.env` file.

All username and email_address are also configured in the `.env` file. This feature is behind a flag of feature flags. Enable feature flag with `_SENDGRID_EMAIL_SERVICE="true"`

Since the token that is used for the keys cannot be added into the system, the env variable scripts do not contain the send_grid email api keys. To add this, you must configure this
from the production server itself. This is done so to prevent misuse of the key.
