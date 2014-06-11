package sshkey_test

import (
	"github.com/nitrous-io/nitrous-common/sshkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "sshkey")
}

var _ = Describe("sshkey", func() {
	Describe("Md5Fingerprint()", func() {
		Context("when a given ssh public key is valid", func() {
			It("returns md5 fingerprint for a given ssh public key", func() {
				fingerprint, err := sshkey.Md5Fingerprint("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCxCeR/jJna6xXELDV8HF51QQPbj++JNDE4W/nqLyhzRTEDDnAfyp8gXW/EXXos8w5XXLak6eUDXE/QI6o8blX39RYoT47u45v4b2scN0ZiS5McmESo193KwWH3gXHbnHWywzPy7+GrRF3nHm/uidvp0tqh4ywqIZujaMjST1sIOQ== test@example.com\n")
				Expect(err).To(BeNil())
				Expect(fingerprint).To(Equal("f7:de:4d:6f:54:6b:b9:d5:78:a9:62:15:e2:ad:ec:bb"))

				fingerprint, err = sshkey.Md5Fingerprint("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzRcsDraIFmeuZ+q6Ic9epAzXPy8WnkqZwLnauH9YBMPhrVQH6nnr/YSGxhYivnG1Qu0AirSUBAjFL7uD9YLeUhwwdKkhHa+bPPMpycmOopyafZyYAEIuZUwyjnK/2mniNJZ+YZhfBJ1z9Xh1hMIFtls/LP1Q7C25rIDg7l1khZ1cEV0mAzJjb4QBvqAt6TjC/YbUZ4/45rHRuQsu/RBhYpuffFwqlLHfjSO2i6TuJtNsEJ2kDDLooM/THf3J2TxObomusErPTgwAo8ya1GmEB5NuwysZx1posjPZIQwg0Q3RqLzMOPlnZpG2l/lvVfEU+iCcc246z4D4Ej9WnuWFR")
				Expect(err).To(BeNil())
				Expect(fingerprint).To(Equal("c1:fc:6d:29:ab:b8:03:57:c4:30:a6:b7:58:40:74:af"))
			})
		})

		Context("when a given ssh publickey is invalid", func() {
			It("returns error", func() {
				fingerprint, err := sshkey.Md5Fingerprint("ssh-rsa lolbbqwtf@&*#^&!@ test@example.com")
				Expect(err).To(MatchError(sshkey.ErrInvalidPublicKey))
				Expect(fingerprint).To(Equal(""))
			})
		})
	})
})
