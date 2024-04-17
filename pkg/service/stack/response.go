package stack

import (
	"time"

	"github.com/emiago/sipgo/sip"
)

const TIME_LAYOUT = "2024-01-01T00:00:00"
const EXPIRES_TIME = 3600

func NewRegisterResponse(req *sip.Request, code sip.StatusCode, reason string) *sip.Response {
	resp := sip.NewResponseFromRequest(req, code, reason, nil)

	newTo := &sip.ToHeader{Address: resp.To().Address, Params: sip.NewParams()}
	newTo.Params.Add("tag", sip.GenerateTagN(10))

	resp.ReplaceHeader(newTo)
	resp.RemoveHeader("Allow")
	expires := sip.ExpiresHeader(EXPIRES_TIME)
	resp.AppendHeader(&expires)
	resp.AppendHeader(sip.NewHeader("Date", time.Now().Format(TIME_LAYOUT)))

	return resp
}
