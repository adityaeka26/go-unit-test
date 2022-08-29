package helper

import "crypto/rand"

func GenerateOTP(length int) (*string, error) {
	otpChars := "1234567890"

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	result := string(buffer)

	return &result, nil
}
