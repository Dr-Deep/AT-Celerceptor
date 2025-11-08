package atfram

import (
	"testing"
)

func TestProofOfWork(t *testing.T) {
	var (
		uuid       = "2a26f5f0-b6f1-46b3-b7f1-369d2e2756a7"
		difficulty = 3 // len of leading zeros
	)

	t.Logf("UUID: %s", uuid)
	t.Logf("Difficulty: %v", difficulty)

	var (
		nonce = _proofOfWork(uuid, difficulty)
		hash  = generateHash(uuid, nonce)
	)

	if err := verifyProofOfWork(hash, difficulty); err != nil {
		t.Error(err)
	}

	t.Logf("Hash: %s", hash)
	t.Logf("Nonce: %v", nonce)
}
