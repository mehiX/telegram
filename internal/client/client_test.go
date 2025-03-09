package client

import (
	"strings"
	"testing"
)

func TestSplitLongMessage(t *testing.T) {

	msg := "lasdkf dsajfaljf ssad fpewl;af afjaf w.erasd f88sfa == f;ldaf;jl a;f; lfka; fk;ds kf 23445;1$1 ldaskfj&*()(DPASDP;k paeipfeiqp r)"

	parts := splitLongMessage(msg, 12)

	reconstructed := strings.Join(parts, "")

	if msg != reconstructed {
		t.Errorf("reconstructed message does not match the original. parts: %d\noriginal: %s\nreconstr: %s\n", len(parts), msg, reconstructed)
	}
}
