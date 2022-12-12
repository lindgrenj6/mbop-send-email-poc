# `POST /sendEmail` POC

This is a small POC for emulating BOP's /sendEmail endpoint.

Took the openapi spec and created a struct in the `mailer/` package as well as the implementation to send a message from there using the AWS SES service.

### Layout

The base HTTP server/handler are in `main.go`, basically all the handler does is:
- Parse the JSON request body to `mailer.Emails` struct
- Instantiate a `Emailer` to send them, using the configured `MAILER_MODULE` in the environment (currently only `"aws"` is supported)
  - `mailer.NewMailer() (mailer.Mailer, error)`
  - Looks at the AWS Config, returning an error if it wasn't initialized
- Loops through the emails and sends them, logging any failures
  - `mailer.SendEmail(context.Context, *mailler.Email) error`

### Questions

- Bail out if one email fails to send?

### Things TODO if we were to pull this in:
- Register the "from" address in AWS SES, hopefully through terraform?
- Create a `sendEmailRequest` struct to parse the body to -> then create the Email struct from that for safety
- Parse the recipients to accept `"Real Name" email@address.com` since that format is in the examples for BOP
- Lookup recipients by username, for same reason as above
