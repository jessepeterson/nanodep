package apinext

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/micromdm/nanodep/albc"

	"github.com/micromdm/nanolib/log"
	"github.com/micromdm/nanolib/log/ctxlog"
)

// NewBypassCodeHandler returns a utility HTTP handler for working with Apple Activation Lock Bypass Codes.
func NewBypassCodeHandler(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bc albc.BypassCode

		code := r.URL.Query().Get("code")
		raw := r.URL.Query().Get("raw")

		var err error

		logger := ctxlog.Logger(r.Context(), logger)

		if raw != "" && code != "" {
			logAndWriteJSONError(logger, w, "validating input", errors.New("raw or code but not both"), http.StatusBadRequest)
			return
		} else if raw == "" && code == "" {
			// no raw or code provided, make a new random code
			bc, err = albc.New()
			if err != nil {
				logAndWriteJSONError(logger, w, "new bypass code", err, http.StatusInternalServerError)
				return
			}
		} else if raw != "" {
			// decode and use raw value
			b, err := hex.DecodeString(raw)
			if err != nil {
				logAndWriteJSONError(logger, w, "decode raw", err, http.StatusBadRequest)
				return
			}
			bc, err = albc.NewFromBytes(b)
			if err != nil {
				logAndWriteJSONError(logger, w, "new from raw", err, http.StatusBadRequest)
				return
			}
		} else if code != "" {
			// decode the dash-separated "human readable" form
			bc, err = albc.NewFromCode(code)
			if err != nil {
				logAndWriteJSONError(logger, w, "decode code", err, http.StatusBadRequest)
				return
			}
		}

		out := &BypassCodeResponseJson{Raw: hex.EncodeToString(bc[:])}

		out.Code, err = bc.Code()
		if err != nil {
			logAndWriteJSONError(logger, w, "create code", err, http.StatusInternalServerError)
			return
		}

		out.Hash, err = bc.Hash()
		if err != nil {
			logAndWriteJSONError(logger, w, "create hash", err, http.StatusInternalServerError)
			return
		}

		writeJSON(w, out, http.StatusOK, logger)
	}
}
