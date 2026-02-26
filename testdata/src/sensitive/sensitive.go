package sensitive

import "log/slog"

func badLiterals() {
	val1 := "hunter2"
	val2 := "eyJhbGciOiJIUzI1NiJ9"
	val3 := "sk-1234567890"

	// sensitive keyword found in the string literal
	slog.Info("user password: " + val1) // want `log message may expose sensitive data \(keyword: "password"\)`
	slog.Debug("api_key=" + val3)       // want `log message may expose sensitive data \(keyword: "api_key"\)`
	slog.Info("token: " + val2)         // want `log message may expose sensitive data \(keyword: "token"\)`
	slog.Error("secret value found")    // want `log message may expose sensitive data \(keyword: "secret"\)`
}

func badVariables() {
	password := "hunter2"
	token := "eyJhbGciOiJIUzI1NiJ9"

	// sensitive keyword found in the variable name
	slog.Info("user id: " + password) // want `log message may expose sensitive data via variable "password"`
	slog.Info("value: " + token)      // want `log message may expose sensitive data via variable "token"`
}

func good() {
	userID := "42"
	email := "john@example.com"

	slog.Info("user created: " + userID)
	slog.Info("email sent to " + email)
	slog.Error("connection refused")
	slog.Debug("request timeout")
}
