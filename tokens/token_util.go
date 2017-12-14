package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

func generateKey() (string, error) {
	n := 128
	secretKey := make([]byte, n)
	if _, err := rand.Read(secretKey); err != nil {
		return "", err
	}
	secretKeyString := base64.StdEncoding.EncodeToString(secretKey)
	secretKeyString = sanitizeKey(secretKeyString)

	if len(secretKeyString) < 40 {
		/* Wow, this is terribly bad luck */
		return "", fmt.Errorf("Failed to create secretKey due to not enough good characters")
	}

	return secretKeyString[0:40], nil
}

func sanitizeKey(key string) string {
	re := regexp.MustCompile("[O0lI+/=]")
	key = re.ReplaceAllString(key, "")
	return strings.Trim(key, "")
}

func getAuthProviderName(externalID string) string {
	parts := strings.Split(externalID, "_")
	return parts[0]
}
