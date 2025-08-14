package token

import (
	"lemfi/simplebank/config"
)

func SetTokenMaker() {
	maker, err := NewJWTMaker(config.Get().TokenSymmetricKey)
	if err != nil {
		config.Logger.Error("cannot create token maker: ", "error", err)
		panic(err)
	}
	tokenMaker = maker
}
