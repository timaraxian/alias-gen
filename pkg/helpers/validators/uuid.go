package validators

import (
	"regexp"

	"github.com/timaraxian/hotel-gen/pkg/errors"
)

var uuidPattern = regexp.MustCompile(
	`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func UUID(uuid string) error {
	if uuidPattern.MatchString(uuid) {
		return nil
	}

	return errors.InvalidUUID
}
