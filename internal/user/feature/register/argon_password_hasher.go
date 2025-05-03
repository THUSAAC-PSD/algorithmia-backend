package register

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"emperror.dev/errors"
	"golang.org/x/crypto/argon2"
)

type ArgonPasswordHasher struct{}

func (h ArgonPasswordHasher) Hash(password string) (string, error) {
	const (
		iterations  = 2
		memory      = 19 * 1024
		parallelism = 1
		keyLength   = 32
	)

	salt, err := h.generateRandomBytes(32)
	if err != nil {
		return "", errors.WrapIf(err, "failed to generate salt")
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		b64Salt,
		b64Hash,
	)
	return encodedHash, nil
}

func (h ArgonPasswordHasher) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
