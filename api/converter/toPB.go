package converter

import (
	v1 "github.com/Teeam-Sync/Sync-Server/api/proto/v1"
)

func (t JWTToken) ToPB() (pb *v1.Token) {
	pb = &v1.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
	return pb
}
