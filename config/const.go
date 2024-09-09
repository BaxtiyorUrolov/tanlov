package config

import "time"

var SignKey = []byte("uhfu8eh8743r348fyu44#$#@#@54fuwf")

const (
	AccessTokenExpireTime = time.Minute * 2
	RefreshTokenExpireTime = time.Hour * 7 * 2
)