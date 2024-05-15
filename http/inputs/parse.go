package inputs

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/mdwhatcott/calcy-apps/ext/shuttle"
)

func parseInteger(query url.Values, key string) (parsed int, err error) {
	raw := query.Get(key)
	parsed, err = strconv.Atoi(raw)
	if err != nil {
		return 0, shuttle.InputError{
			Fields:  []string{fmt.Sprintf("query:%s", key)},
			Message: fmt.Sprintf("failed to parse '%s' parameter as integer: [%s]", key, raw),
		}
	}
	return parsed, nil
}
