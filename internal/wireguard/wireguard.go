package wireguard

import (
	"fmt"
	"regexp"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GenerateKey() (wgtypes.Key, error) {
	key, keyErr := wgtypes.GeneratePrivateKey()

	if keyErr != nil {
		return key, &WireguardGenerateKeyError{Err: keyErr}
	}
	return key, nil
}

// FIXME: Instead of doing a regex replace, decide if we should use a parser
func ConfigAddKey(config string, key wgtypes.Key) string {
	interface_section := "[Interface]"
	interface_section_escaped := regexp.QuoteMeta(interface_section)

	// (?m) enables multi line mode
	// ^ match from beginning of line
	// $ match till end of line
	// So it matches [Interface] section exactly
	interface_re := regexp.MustCompile(fmt.Sprintf("(?m)^%s$", interface_section_escaped))
	to_replace := fmt.Sprintf("%s\nPrivateKey = %s", interface_section, key.String())
	return interface_re.ReplaceAllString(config, to_replace)
}

type WireguardGenerateKeyError struct {
	Err error
}

func (e *WireguardGenerateKeyError) Error() string {
	return fmt.Sprintf("failed generating Wireguard key with error: %v", e.Err)
}