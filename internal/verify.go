package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/jedisct1/go-minisign"
)

// getKeys returns keys taken from https://git.sr.ht/~eduvpn/disco.eduvpn.org#public-keys.
func getKeys() []string {
	return []string{
		"RWRtBSX1alxyGX+Xn3LuZnWUT0w//B6EmTJvgaAxBMYzlQeI+jdrO6KF", // fkooman@tuxed.net, kolla@uninett.no
		"RWQKqtqvd0R7rUDp0rWzbtYPA3towPWcLDCl7eY9pBMMI/ohCmrS0WiM", // RoSp
	}
}

// Verify verifies the signature (.minisig file format) on signedJson.
//
// expectedFileName must be set to the file type to be verified, either "server_list.json" or "organization_list.json".
// minSign must be set to the minimum UNIX timestamp (without milliseconds) for the file version.
// This value should not be smaller than the time on the previous document verified.
// forcePrehash indicates whether or not we want to force the use of prehashed signatures
// In the future we want to remove this parameter and only allow prehashed signatures
//
// The return value will either be (true, nil) for a valid signature or (false, VerifyError) otherwise.
//
// Verify is a wrapper around verifyWithKeys where allowedPublicKeys is set to the list from https://git.sr.ht/~eduvpn/disco.eduvpn.org#public-keys.
func Verify(signatureFileContent string, signedJson []byte, expectedFileName string, minSignTime uint64, forcePrehash bool) (bool, error) {
	keyStrs := getKeys()
	if extraKey != "" {
		keyStrs = append(keyStrs, extraKey)
		_, err := fmt.Fprintf(os.Stderr, "INSECURE TEST MODE ENABLED WITH KEY %q\n", extraKey)
		if err != nil {
			panic(err)
		}
	}
	valid, err := verifyWithKeys(signatureFileContent, signedJson, expectedFileName, minSignTime, keyStrs, forcePrehash)
	if err != nil {
		var verifyCreatePublickeyError *VerifyCreatePublicKeyError
		if errors.As(err, &verifyCreatePublickeyError) {
			panic(err) // This should not happen unless keyStrs has an invalid key
		}
		return valid, err
	}
	return valid, nil
}

// extraKey is an extra allowed key for testing.
var extraKey = ""

// InsecureTestingSetExtraKey adds an extra allowed key for verification with Verify.
// ONLY USE FOR TESTING. Applies to all threads. Probably not thread-safe. Do not call in parallel to Verify.
//
// keyString must be a Base64-encoded Minisign key, or empty to reset.
func InsecureTestingSetExtraKey(keyString string) {
	extraKey = keyString
}

type VerifyUnknownExpectedFilenameError struct {
	Filename string
	Expected string
}

func (e *VerifyUnknownExpectedFilenameError) Error() string {
	return fmt.Sprintf("invalid filename %s, expected %s", e.Filename, e.Expected)
}

type VerifyInvalidSignatureFormatError struct {
	Err error
}

func (e *VerifyInvalidSignatureFormatError) Error() string {
	return fmt.Sprintf("invalid signature format, error %v", e.Err)
}

type VerifyInvalidSignatureAlgorithmError struct {
	Algorithm       string
	WantedAlgorithm string
}

func (e *VerifyInvalidSignatureAlgorithmError) Error() string {
	return fmt.Sprintf("invalid signature algorithm %s, wanted %s", e.Algorithm, e.WantedAlgorithm)
}

type VerifyCreatePublicKeyError struct {
	PublicKey string
	Err       error
}

func (e *VerifyCreatePublicKeyError) Error() string {
	return fmt.Sprintf("failed to create public key %s with error %v", e.PublicKey, e.Err)
}

type VerifyInvalidSignatureError struct {
	Err error
}

func (e *VerifyInvalidSignatureError) Error() string {
	return fmt.Sprintf("invalid signature with error %v", e.Err)
}

type VerifyInvalidTrustedCommentError struct {
	TrustedComment string
	Err            error
}

func (e *VerifyInvalidTrustedCommentError) Error() string {
	return fmt.Sprintf("invalid trusted comment %s with error %v", e.TrustedComment, e.Err)
}

type VerifyWrongSigFilenameError struct {
	Filename    string
	SigFilename string
}

func (e *VerifyWrongSigFilenameError) Error() string {
	return fmt.Sprintf("wrong filename %s, expected filename %s for signature", e.Filename, e.SigFilename)
}

type VerifySigTimeEarlierError struct {
	SigTime    uint64
	MinSigTime uint64
}

func (e *VerifySigTimeEarlierError) Error() string {
	return fmt.Sprintf("Sign time %d is earlier than sign time %d", e.SigTime, e.MinSigTime)
}

type VerifyUnknownKeyError struct {
	Filename string
}

func (e *VerifyUnknownKeyError) Error() string {
	return fmt.Sprintf("signature for filename %s was created with an unknown key", e.Filename)
}

// verifyWithKeys verifies the Minisign signature in signatureFileContent (minisig file format) over the server_list/organization_list JSON in signedJson.
//
// Verification is performed using a matching key in allowedPublicKeys.
// The signature is checked to be a Ed25519 Minisign (optionally Ed25519 Blake2b-512 prehashed, see forcePrehash) signature with a valid trusted comment.
// The file type that is verified is indicated by expectedFileName, which must be one of "server_list.json"/"organization_list.json".
// The trusted comment is checked to be of the form "timestamp:<timestamp>\tfile:<expectedFileName>", optionally suffixed by something, e.g. "\thashed".
// The signature is checked to have a timestamp with a value of at least minSignTime, which is a UNIX timestamp without milliseconds.
//
// The return value will either be (true, nil) on success or (false, detailedVerifyError) on failure.
func verifyWithKeys(signatureFileContent string, signedJson []byte, filename string, minSignTime uint64, allowedPublicKeys []string, forcePrehash bool) (bool, error) {
	switch filename {
	case "server_list.json", "organization_list.json":
		break
	default:
		return false, &VerifyUnknownExpectedFilenameError{Filename: filename, Expected: "server_list.json or organization_list.json"}
	}

	sig, err := minisign.DecodeSignature(signatureFileContent)
	if err != nil {
		return false, &VerifyInvalidSignatureFormatError{Err: err}
	}

	// Check if signature is prehashed, see https://jedisct1.github.io/minisign/#signature-format
	if forcePrehash && sig.SignatureAlgorithm != [2]byte{'E', 'D'} {
		return false, &VerifyInvalidSignatureAlgorithmError{Algorithm: string(sig.SignatureAlgorithm[:]), WantedAlgorithm: "ED (BLAKE2b-prehashed EdDSA)"}
	}

	// Find allowed key used for signature
	for _, keyStr := range allowedPublicKeys {
		key, err := minisign.NewPublicKey(keyStr)
		if err != nil {
			// Should only happen if Verify is wrong or extraKey is invalid
			return false, &VerifyCreatePublicKeyError{PublicKey: keyStr, Err: err}
		}

		if sig.KeyId != key.KeyId {
			continue // Wrong key
		}

		valid, err := key.Verify(signedJson, sig)
		if !valid {
			return false, &VerifyInvalidSignatureError{Err: err}
		}

		// Parse trusted comment
		var signTime uint64
		var sigFileName string
		// sigFileName cannot have spaces
		_, err = fmt.Sscanf(sig.TrustedComment, "trusted comment: timestamp:%d\tfile:%s", &signTime, &sigFileName)
		if err != nil {
			return false, &VerifyInvalidTrustedCommentError{TrustedComment: sig.TrustedComment, Err: err}
		}

		if sigFileName != filename {
			return false, &VerifyWrongSigFilenameError{Filename: filename, SigFilename: sigFileName}
		}

		if signTime < minSignTime {
			return false, &VerifySigTimeEarlierError{SigTime: signTime, MinSigTime: minSignTime}
		}

		return true, nil
	}

	// No matching allowed key found
	return false, &VerifyUnknownKeyError{Filename: filename}
}