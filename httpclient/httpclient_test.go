package httpclient_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/nitrous-io/nitrous-common/httpclient"

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
		headers = http.Header{
			"One": []string{"two"},
			"A":   []string{"bc"},
		}
		check func()
	)

	BeforeEach(func() {
		srv = ghttp.NewServer()
		c = httpclient.New(token)
	})

	AfterEach(func() {
		srv.Close()
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

	ItIncludesHeaders := func() {
		It("includes headers", func() {
			srv.AppendHandlers(
				func(w http.ResponseWriter, r *http.Request) {
					for k, _ := range headers {
						Expect(r.Header.Get(k)).To(headers.Get(k))
					}
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

	Describe("GetWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.GetWithHeaders(srv.URL(), params, headers) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("GET")
		ItIncludesQueryString()
		ItIncludesHeaders()
	})

	Describe("Delete()", func() {
		BeforeEach(func() {
			check = func() { c.Delete(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("DELETE")
		ItIncludesQueryString()
	})

	Describe("DeleteWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.DeleteWithHeaders(srv.URL(), params, headers) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("DELETE")
		ItIncludesQueryString()
		ItIncludesHeaders()
	})

	Describe("PostForm()", func() {
		BeforeEach(func() {
			check = func() { c.PostForm(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("POST")
		ItIncludesFormData()
	})

	Describe("PostFormWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.PostFormWithHeaders(srv.URL(), params, headers) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("POST")
		ItIncludesFormData()
		ItIncludesHeaders()
	})

	Describe("PutForm()", func() {
		BeforeEach(func() {
			check = func() { c.PutForm(srv.URL(), params) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("PUT")
		ItIncludesFormData()
	})

	Describe("PutFormWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.PutFormWithHeaders(srv.URL(), params, headers) }
		})

		ItHasNoAuthHeader()
		ItUsesHTTPVerb("PUT")
		ItIncludesFormData()
		ItIncludesHeaders()
	})

	Describe("GetAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.GetAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("GET")
		ItIncludesQueryString()
	})

	Describe("GetAuthedWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.GetAuthedWithHeaders(srv.URL(), params, headers) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("GET")
		ItIncludesQueryString()
		ItIncludesHeaders()
	})

	Describe("DeleteAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.DeleteAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("DELETE")
		ItIncludesQueryString()
	})

	Describe("DeleteAuthedWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.DeleteAuthedWithHeaders(srv.URL(), params, headers) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("DELETE")
		ItIncludesQueryString()
		ItIncludesHeaders()
	})

	Describe("PostFormAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.PostFormAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("POST")
		ItIncludesFormData()
	})

	Describe("PostFormAuthedWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.PostFormAuthedWithHeaders(srv.URL(), params, headers) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("POST")
		ItIncludesFormData()
		ItIncludesHeaders()
	})

	Describe("PutFormAuthed()", func() {
		BeforeEach(func() {
			check = func() { c.PutFormAuthed(srv.URL(), params) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("PUT")
		ItIncludesFormData()
	})

	Describe("PutFormAuthedWithHeaders()", func() {
		BeforeEach(func() {
			check = func() { c.PutFormAuthedWithHeaders(srv.URL(), params, headers) }
		})

		ItHasAuthHeader()
		ItUsesHTTPVerb("PUT")
		ItIncludesFormData()
		ItIncludesHeaders()
	})

	Describe("CheckResponse()", func() {
		Context("when expected response is received", func() {
			BeforeEach(func() {
				srv.AppendHandlers(
					ghttp.RespondWith(http.StatusOK, ""),
				)
			})

			It("returns nil", func() {
				resp, err := c.Get(srv.URL(), nil)
				Expect(err).To(BeNil())

				errResponse := httpclient.CheckResponse(resp, http.StatusOK)
				Expect(errResponse).To(BeNil())
			})
		})

		Context("when an error response with a description is received", func() {
			BeforeEach(func() {
				srv.AppendHandlers(
					ghttp.RespondWith(http.StatusInternalServerError,
						`{ "error": "internal_server_error", "error_description": "some error occurred" }`),
				)
			})

			It("returns ErrResponse with description", func() {
				resp, err := c.Get(srv.URL(), nil)
				Expect(err).To(BeNil())

				errResponse := httpclient.CheckResponse(resp, http.StatusOK).(*httpclient.ErrResponse)
				Expect(errResponse.Error()).To(Equal("some error occurred"))
			})
		})

		Context("when an error response with a hash is received", func() {
			BeforeEach(func() {
				srv.AppendHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, `{ "error": "internal_server_error", "error_hash": "1234" }`),
				)
			})

			It("returns ErrResponse with hash", func() {
				resp, err := c.Get(srv.URL(), nil)
				Expect(err).To(BeNil())

				errResponse := httpclient.CheckResponse(resp, http.StatusOK).(*httpclient.ErrResponse)
				Expect(errResponse.Error()).To(ContainSubstring("something went wrong (contact support with this error code: 1234)"))
			})
		})

		Context("when an error response without a description or a hash is received", func() {
			BeforeEach(func() {
				srv.AppendHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, `{ "error": "internal_server_error" }`),
				)
			})

			It("returns ErrResponse with hash", func() {
				resp, err := c.Get(srv.URL(), nil)
				Expect(err).To(BeNil())

				errResponse := httpclient.CheckResponse(resp, http.StatusOK).(*httpclient.ErrResponse)
				Expect(errResponse.Error()).To(ContainSubstring("unknown error occurred"))
			})
		})
	})
})
