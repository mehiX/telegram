package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestSendTo(t *testing.T) {

	msg := "lasdjf adsflj woq 2iuwqr 3234 2343 aslfja89 877f*&(&&( afjwqeo4r97))vaferq werqerqerdf afj fjljdsfjdalfjlkdfjlwoq4uroeqwpr-0851"
	testToken := "test-token-mihai"
	testPath := fmt.Sprintf("/bot%s/sendMessage", testToken)

	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path != testPath {
			t.Errorf("bad request path. expecting: %s, got: %s", testPath, path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("bad method. expecting: %s, got: %s", http.MethodPost, r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	cli := &TClient{HttpClient: srvr.Client(), Token: testToken}

	telegramApiURL = srvr.URL

	if err := cli.send(message{chatID: "my-chatid-123", txt: msg}); err != nil {
		t.Errorf("unexpected error sending the message: %v", err)
	}
}
