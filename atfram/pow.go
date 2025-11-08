package atfram

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
TextOutputCallback: ouput[0]:message:value decoden und in engine
entweder js engine (document=json-raw *AldiTalk_Auth_Challenge)
oder in go neu implementieren
*/

const (
	powMaxGenerationsPerIteration = 5000
	//count := 0; count < powMaxGenerationsPerIteration; count++
)

var (
	/* REGEXP
	   var difficulty = 3;
	   var work = "2a26f5f0-b6f1-46b3-b7f1-369d2e2756a7";
	*/
	aldiTalk_PoW_Difficulty_RE = regexp.MustCompile(`var\s+difficulty\s*=\s*(?P<difficulty>\d+);`)
	aldiTalk_PoW_UUID_RE       = regexp.MustCompile(`var\s+work\s*=\s*"(?P<work>[a-f0-9-]+)"`)
)

func GetProofOfWorkHash(uuid, _difficulty string) (string, error) {
	diff, err := strconv.Atoi(_difficulty)
	if err != nil {
		return "", err
	}

	var (
		nonce = _proofOfWork(uuid, diff)
		hash  = generateHash(uuid, nonce)
	)

	if err := verifyProofOfWork(hash, diff); err != nil {
		return "", err
	}

	return hash, nil

}

func verifyProofOfWork(hash string, difficulty int) error {
	if !strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
		return fmt.Errorf(
			"%w: Hash: %s Difficulty: %v",
			ErrAldiTalkCallbackPoW,
			hash,
			difficulty,
		)
	}

	return nil
}

func _proofOfWork(uuid string, difficulty int) (nonce int) {
	var target = strings.Repeat("0", difficulty)
	nonce = 0

	for {
		hash := generateHash(uuid, nonce)
		if strings.HasPrefix(hash, target) {
			return nonce // nonce found
		}
		nonce++
	}
}

// generateHash -> SHA1(uuid + nonce)
func generateHash(uuid string, nonce int) string {
	msg := uuid + strconv.Itoa(nonce)
	hash := sha1.Sum([]byte(msg))

	return hex.EncodeToString(hash[:])
}
