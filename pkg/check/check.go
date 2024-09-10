package check

import (
	"context"
	"fmt"
	"unicode"
	"it-tanlov/storage"
)

func PhoneNumber(phone string) bool {
	for _, r := range phone {
		if r == '+' {
			continue
		} else if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func IPhoneExist(phone string, partnerStorage storage.IPartnerStorage) (bool, error) {
	exists, err := partnerStorage.PhoneExist(context.Background(), phone)
	if err != nil {
		return false, fmt.Errorf("error while checking phone existence: %w", err)
	}
	return exists, nil
}

func IEmailExist(email string, partnerStorage storage.IPartnerStorage) (bool, error) {
	exists, err := partnerStorage.IEmailExist(context.Background(), email)
	if err != nil {
		return false, fmt.Errorf("error while checking email existence: %w", err)
	}
	return exists, nil
}

func IVideoLinkExist(video_link string, partnerStorage storage.IPartnerStorage) (bool, error) {
	exists, err := partnerStorage.IVideoLinkExist(context.Background(), video_link)
	if err != nil {
		return false, fmt.Errorf("error while checking video_link existence: %w", err)
	}
	return exists, nil
}


func IUserEmailExist(phone string, userStorage storage.IUserStorage) (bool, error) {
	exists, err := userStorage.IUserEmailExist(context.Background(), phone)
	if err != nil {
		return false, fmt.Errorf("error while checking phone existence: %w", err)
	}
	return exists, nil
}