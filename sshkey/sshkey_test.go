package sshkey_test

import (
	"github.com/nitrous-io/nitrous-common/sshkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "sshkey")
}

var _ = Describe("sshkey", func() {
	Describe("Md5Fingerprint()", func() {
		testCases := []struct {
			fingerprint string
			keydata     string
		}{
			{ // rsa with comment
				"f7:de:4d:6f:54:6b:b9:d5:78:a9:62:15:e2:ad:ec:bb",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCxCeR/jJna6xXELDV8HF51QQPbj++JNDE4W/nqLyhzRTEDDnAfyp8gXW/EXXos8w5XXLak6eUDXE/QI6o8blX39RYoT47u45v4b2scN0ZiS5McmESo193KwWH3gXHbnHWywzPy7+GrRF3nHm/uidvp0tqh4ywqIZujaMjST1sIOQ== test@example.com\n",
			},
			{ // rsa
				"c1:fc:6d:29:ab:b8:03:57:c4:30:a6:b7:58:40:74:af",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzRcsDraIFmeuZ+q6Ic9epAzXPy8WnkqZwLnauH9YBMPhrVQH6nnr/YSGxhYivnG1Qu0AirSUBAjFL7uD9YLeUhwwdKkhHa+bPPMpycmOopyafZyYAEIuZUwyjnK/2mniNJZ+YZhfBJ1z9Xh1hMIFtls/LP1Q7C25rIDg7l1khZ1cEV0mAzJjb4QBvqAt6TjC/YbUZ4/45rHRuQsu/RBhYpuffFwqlLHfjSO2i6TuJtNsEJ2kDDLooM/THf3J2TxObomusErPTgwAo8ya1GmEB5NuwysZx1posjPZIQwg0Q3RqLzMOPlnZpG2l/lvVfEU+iCcc246z4D4Ej9WnuWFR",
			},
			{ // ecdsa
				"c8:ce:12:46:4b:37:b5:12:5f:f5:81:a8:b2:d8:0a:29",
				"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBNjtS9yPn7Kiowq/SkNC3BwA7UEugtnLV0Wt6JwWy041fmg0TOMESsgMzV2bIwz6ArWKLR7vhcZOL/jP4OCM20g=",
			},
			{ // dsa
				"0a:3f:9f:66:c2:17:cf:15:5f:4a:26:2a:60:48:9c:a3",
				"ssh-dss AAAAB3NzaC1kc3MAAACBAOnV2UjLF88/tMZi6stfie6qbQNHVtAbzk4vUPkBnnXV3mL27NI3Dxapr8R9vQ8urd4kG+B9j4KM/ycEmgqM24X0RIEZCfDFqKEflzZDZuJWxHSulGqsBXQpuBah6bs5Su6WoeYMAEffVUpuOrSrZ0Y6dtt19j58/MXn4V0Dxmc9AAAAFQCCEaPBHd78Ehu9MI+Bxt+NGAMwFQAAAIEA2WQOwSDJR+6FFSDRCCx+7k1sZIfTMg5XOhYfYtZHU6NQGhCCdrtlAXvy6mXMQpFyVWjIA/Pp8wj2Vt2+akKWeQ9FXuTcdCKfVQsdvL8WjsEyIMqJHPHE0uSPJMoFyWo7VaDelsbeMc9Vn4Isof4y6/uvxtbJeXyluYaBWIDDpWQAAACBANLvKHPjgV50faooAy5YfB+n4gljM0ngmBJ7jJOj8Sws3TWVb15LzYL++z5iCGUq7mss4keHu4zDSHA6SC2DxMSJUFHM38WL/09GocxgqyLkqo9EZvUOItgUXt5l3Nf06GnuFEG/4Tk68GeS0mbqxG/wrE8R+FJj3dRuGME+9K7l",
			},
			{
				"invalid",
				"AAAAB3NzaC1yc2EAAAADAQABAAAAgQCxCeR/jJnnopenotvalidatall",
			},
			{
				"invalid",
				"ssh-nsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCxCeR/jJnnopenotvalidatall",
			},
			{
				"invalid",
				"ssh-rsa lolbbqwtf@&*#^&!@ test@example.com",
			},
		}

		for _, tc2 := range testCases {
			tc := tc2
			if tc.fingerprint == "invalid" {
				It("returns invalid error", func() {
					fingerprint, err := sshkey.Md5Fingerprint(tc.keydata)
					Expect(err).To(MatchError(sshkey.ErrInvalidPublicKey))
					Expect(fingerprint).To(Equal(""))
				})
			} else {
				It("returns fingerprint for given ssh public key: "+tc.fingerprint, func() {
					fingerprint, err := sshkey.Md5Fingerprint(tc.keydata)
					Expect(err).To(BeNil())
					Expect(fingerprint).To(Equal(tc.fingerprint))
				})
			}
		}
	})
})
