package proxanne_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/waaaaargh/proxanne"
)

var _ = Describe("Route", func() {
	var s *ghttp.Server
	var r proxanne.Router
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		s = ghttp.NewServer()
		s.AllowUnhandledRequests = true

		request, _ = http.NewRequest("GET", "/foobar", nil)
		recorder = httptest.NewRecorder()
	})

	When("a matching Route is found", func() {
		It("should forward a request to the target", func() {
			s_url, _ := url.Parse(s.URL())

			r = proxanne.Router{
				&proxanne.Route{
					Matches: regexp.MustCompile("^.*$"),
					Target:  httputil.NewSingleHostReverseProxy(s_url),
				},
			}

			r.ServeHTTP(recorder, request)

			Expect(s.ReceivedRequests()).ToNot(BeEmpty())
		})
	})

	When("no matching Route is found", func() {
		It("Should return a 404 error to the client", func() {
			r = proxanne.Router{}

			r.ServeHTTP(recorder, request)

			Expect(s.ReceivedRequests()).To(BeEmpty())
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	AfterEach(func() {
		s.Close()
	})
})

func TestRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proxanne Suite")
}
