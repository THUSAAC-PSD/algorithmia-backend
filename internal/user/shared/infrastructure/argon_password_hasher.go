package infrastructure

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"golang.org/x/crypto/argon2"
)

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var defaultArgonParams = argonParams{
	iterations:  2,
	memory:      19 * 1024,
	parallelism: 1,
	keyLength:   32,
}

type ArgonPasswordHasher struct{}

func (h ArgonPasswordHasher) Hash(password string) (string, error) {
	salt, err := h.generateRandomBytes(32)
	if err != nil {
		return "", errors.WrapIf(err, "failed to generate salt")
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultArgonParams.iterations,
		defaultArgonParams.memory,
		defaultArgonParams.parallelism,
		defaultArgonParams.keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultArgonParams.memory,
		defaultArgonParams.iterations,
		defaultArgonParams.parallelism,
		b64Salt,
		b64Hash,
	)
	return encodedHash, nil
}

func (h ArgonPasswordHasher) Check(hashedPassword, plainPassword string) (bool, error) {
	p, salt, hash, err := h.decodeHash(hashedPassword)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(plainPassword), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (h ArgonPasswordHasher) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h ArgonPasswordHasher) decodeHash(encodedHash string) (p *argonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("unsupported argon2 version")
	}

	p = &argonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
