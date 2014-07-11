package httpclient_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/nitrous-io/nitro/pkg/httpclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "httpclient")
}

var _ = Describe("httpclient", func() {
	var (
		srv    *ghttp.Server
		c      *httpclient.Client
		token  = "correct_token"
		params = &url.Values{
			"containers": {"arch", "ubuntu"},
			"say":        {"hello world!"},
		}
		check func()
	)

	BeforeEach(func() {
		srv = ghttp.NewServer()
		c = httpclient.New(token)
	})

	ItHasAuthHeader := func() {
		It("has auth header", func() {
			srv.AppendHandlers(
				ghttp.VerifyHeader(http.Header{
					"Authorization": []string{"Bearer " + token},
				}))
			check()
		})
	}

	ItHasNoAuthHeader := func() {
		It("has no auth header", func() {
			srv.AppendHandlers(
				ghttp.VerifyHeader(http.Header{
					"Authorization": nil,
				}))
			check()
		})
	}

	ItUsesHTTPVerb := func(verb string) {
		It("is a "+verb+" request", func() {
			srv.AppendHandlers(
				func(w http.ResponseWriter, r *http.Request) {
					Expect(r.Method).To(Equal(verb), "Verb mismatch")
				},
			)
			check()
		})
	}

	ItIncludesQueryString := func() {
		It("includes query string", func() {
			srv.AppendHandlers(
				func(w http.ResponseWriter, r *http.Request) {
					Expect(strings.HasSuffix(r.URL.String(), "?containers=arch&containers=ubuntu&say=hello+world%21")).To(BeTrue())
				},
			)
			check()
		})
	}

	ItIncludesFormData := func() {
		It("includes form data", func() {
			srv.AppendHandlers(
				func(w http.ResponseWriter, r *http.Request) {
					body, err := ioutil.ReadAll(r.Body)
					Expect(err).To(BeNil())
					Expect(body).To(Equal([]byte("containers=arch&containers=ubuntu&say=hello+world%21")))
				},
			)
			check()
		})
	}

	Describe("Get()", func() {
		BeforeEach(func() {
			check = func() { c.Get(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("GET")
		ItIncludesQueryString()
	})

	Describe("Delete()", func() {
		BeforeEach(func() {
			check = func() { c.Delete(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("DELETE")
		ItIncludesQueryString()
	})

	Describe("PostForm()", func() {
		BeforeEach(func() {
			check = func() { c.PostForm(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("POST")
		ItIncludesFormData()
	})

	Describe("PutForm()", func() {
		BeforeEach(func() {
			check = func() { c.PutForm(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("PUT")
		ItIncludesFormData()
	})

	Describe("GetAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.GetAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("GET")
		ItIncludesQueryString()
	})

	Describe("DeleteAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.DeleteAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("DELETE")
		ItIncludesQueryString()
	})

	Describe("PostFormAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.PostFormAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("POST")
		ItIncludesFormData()
	})

	Describe("PutFormAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.PutFormAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("PUT")
		ItIncludesFormData()
	})
})
