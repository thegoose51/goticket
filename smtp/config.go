package smtp

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Config defines the SMTP configuration used in sending email notifications.
type Config struct {
	Enabled  bool   `json:"enabled" toml:"enabled" yaml:"enabled"`
	Host     string `json:"host" toml:"host" yaml:"host"`
	Port     int    `json:"port" toml:"port" yaml:"port"`
	Username string `json:"username" toml:"username" yaml:"username"`
	Password string `json:"password" toml:"password" yaml:"password"`

	// Whether to skip TLS verify.
	NoVerify bool `json:"no-verify" toml:"no-verify" yaml:"no-verify"`

	// Whether all alerts should trigger an email.
	Global bool `json:"global" toml:"global" yaml:"global"`

	// Whether all alerts should automatically use stateChangesOnly mode.
	// Only applies if global is also set.
	StateChangesOnly bool `json:"state-changes-only" toml:"state-changes-only" yaml:"state-changes-only"`

	// From address
	From string `json:"from" toml:"from" yaml:"from"`

	// Default To addresses
	To []string `json:"to" toml:"to" yaml:"to"`

	// Close connection to SMTP server after idle timeout has elapsed
	IdleTimeout time.Duration `json:"idle-timeout" toml:"idle-timeout" yaml:"idle-timeout"`
}

// NewConfig creates a new config with default values.
//    return Config{
//        Host: "localhost",
//        Port: 25,
//        IdleTimeout: time.Duraction(time.Second * 30),
//    }
func NewConfig() Config {
	return Config{
		Host:        "localhost",
		Port:        25,
		IdleTimeout: time.Duration(time.Second * 30),
	}
}

// Validate the configuration with a few simple checks.
//     1. Host is not empty.
//     2. Port is greater than 0 and less than 65536.
//     3. IdleTimeout is not negative.
//     4. From field is not empty.
//     5. From contains '@'.
//     6. To(s) contains '@'.
func (c Config) Validate() error {
	if c.Host == "" {
		return errors.New("host cannot be empty")
	}
	if c.Port <= 0 || c.Port >= 65536 {
		return errors.Errorf("invalid port %d", c.Port)
	}
	if c.IdleTimeout < 0 {
		return errors.New("idle timeout must be positive")
	}
	if c.Enabled && c.From == "" {
		return errors.New("must provide a 'from' address")
	}
	// Poor mans email validation, but since emails have a very large domain this is probably good enough
	// to catch user error.
	if c.From != "" && !strings.ContainsRune(c.From, '@') {
		return errors.Errorf("invalid from email address: %q", c.From)
	}
	for _, t := range c.To {
		if t != "" && !strings.ContainsRune(t, '@') {
			return errors.Errorf("invalid to email address: %q", t)
		}
	}
	return nil
}
