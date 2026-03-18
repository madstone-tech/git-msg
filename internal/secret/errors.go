package secret

import "errors"

// ErrNoCredential is returned when no API key can be found via any lookup path.
var ErrNoCredential = errors.New("no credential found")
