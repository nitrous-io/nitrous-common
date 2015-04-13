package hmactoken_test

import (
	"github.com/nitrous-io/nitrous-common/hmactoken"
	"github.com/nitrous-io/nitrous-common/msgverifier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "hmactoken")
}

var _ = Describe("hmactoken", func() {
	const secret = "s3cr3t"

	Describe("Generate", func() {
		It("generates a new hmac token string with given permissions", func() {
			tok, err := hmactoken.Generate(19, []string{"read", "write", "execute"}, secret)
			Expect(err).To(BeNil())
			Expect(tok).NotTo(BeEmpty())

			var t *hmactoken.Token
			err = msgverifier.New(secret).Verify(tok, &t)
			Expect(err).To(BeNil())
			Expect(t.Permissions).To(Equal([]string{"read", "write", "execute"}))
			Expect(t.CreatedAt).NotTo(BeNil())
			Expect(t.ExpiresAt).NotTo(BeNil())
			Expect(t.ExpiresAt.Sub(*t.CreatedAt)).To(Equal(hmactoken.DefaultValidityPeriod))
			Expect(t.UserID).To(Equal(int64(19)))
		})
	})

	Describe("Verify", func() {
		Context("when a given token is invalid", func() {
			It("returns an error", func() {
				tok, err := hmactoken.Verify([]string{"read", "execute"}, "foobar123", secret)
				Expect(err).To(Equal(hmactoken.ErrInvalidToken))
				Expect(tok).To(BeNil())
			})
		})

		Context("when a given token is generated with a different secret", func() {
			It("returns an error", func() {
				tok, err := hmactoken.Generate(19, []string{"read", "write", "execute"}, "different-secret")
				Expect(err).To(BeNil())
				Expect(tok).NotTo(BeEmpty())
				t, err := hmactoken.Verify([]string{"read", "execute"}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrInvalidToken))
				Expect(t).To(BeNil())
			})
		})

		Context("when a given token is valid, but does not grant all of the requested permissions", func() {
			It("returns an error", func() {
				tok, err := hmactoken.Generate(19, []string{"read", "write", "execute"}, secret)
				Expect(err).To(BeNil())
				Expect(tok).NotTo(BeEmpty())

				t, err := hmactoken.Verify([]string{"delete"}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrPermissionDenied))
				Expect(t).To(BeNil())
				t, err = hmactoken.Verify([]string{"read", "write", "execute", "delete"}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrPermissionDenied))
				Expect(t).To(BeNil())
			})
		})

		Context("when a given token is valid, and has all of the requested permissions", func() {
			It("returns an error", func() {
				tok, err := hmactoken.Generate(19, []string{"read", "write", "execute"}, secret)
				Expect(err).To(BeNil())
				Expect(tok).NotTo(BeEmpty())

				t, err := hmactoken.Verify([]string{}, tok, secret)
				Expect(err).To(BeNil())
				Expect(t.UserID).To(Equal(int64(19)))
				t, err = hmactoken.Verify([]string{"read"}, tok, secret)
				Expect(err).To(BeNil())
				Expect(t.UserID).To(Equal(int64(19)))
				t, err = hmactoken.Verify([]string{"read", "execute"}, tok, secret)
				Expect(err).To(BeNil())
				Expect(t.UserID).To(Equal(int64(19)))
				t, err = hmactoken.Verify([]string{"read", "write", "execute"}, tok, secret)
				Expect(err).To(BeNil())
				Expect(t.UserID).To(Equal(int64(19)))
			})
		})

		Context("when a given token is valid, and has all of the requested permissions, but is expired", func() {
			It("returns an error", func() {
				yesterday := time.Now().Add(-24 * time.Hour)
				tok, err := msgverifier.New(secret).Generate(&hmactoken.Token{
					Permissions: []string{"read", "write", "execute"},
					CreatedAt:   &yesterday,
					ExpiresAt:   &yesterday,
				})
				Expect(err).To(BeNil())
				Expect(tok).NotTo(BeEmpty())

				t, err := hmactoken.Verify([]string{}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrTokenExpired))
				Expect(t).To(BeNil())
				t, err = hmactoken.Verify([]string{"read"}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrTokenExpired))
				Expect(t).To(BeNil())
				t, err = hmactoken.Verify([]string{"read", "execute"}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrTokenExpired))
				Expect(t).To(BeNil())
				t, err = hmactoken.Verify([]string{"read", "write", "execute"}, tok, secret)
				Expect(err).To(Equal(hmactoken.ErrTokenExpired))
				Expect(t).To(BeNil())
			})
		})
	})
})
