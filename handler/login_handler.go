package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/marsuc/SawitProTest/pkg/jwtrsa"
	"github.com/marsuc/SawitProTest/pkg/password"
	"github.com/marsuc/SawitProTest/repository"
)

const (
	invalidUsernameOrPassword = "invalid username or password"
	accessTokenDuration       = 24 * time.Hour
	privateKey                = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAtS/EiblnTN0HePO6jCSYPOXaOKN2vl4paxP0dD7NrL3dvJid\nW8iCNuQUyKkKeYFF54/6rzigjHamG3KhmGy2M/GwGLTyaNIEFj8uPW580orJAVcy\nbFXZWg2coXSQwGOrbcsXgk8ZcwviXsIpuXULai5oL7FaQK1l0o7biUt3doGN9k2Y\nSBNDwLFEQplCSE9qoeabKbKqjgzZmTqZV1yvIJqGSTTQHQh9tqdjJnKF3iTmGtp7\nw8yIQ/hjGnKW4585gHMn5y7NQdTBjnlLXbC9dHJoNnI66Fw8NNF/GtOkX1fZhbxD\nxGkdZzv1+OvG4XYfFF7YYpV+1zyZaulit5/w5wIDAQABAoIBAQCdp2fLMtEot3Zm\nDyV2Be2Vp6be+5U6BfLiIiXl9DPAqCDFlsHteCWdn45aH4Rmv05VNBm713kTX7Yf\nUfo8B/PudNF2XhRDkuJNfUI93+KqzGokSXwtefG7AvcUIbpGPTOQQFQ/ZZZOXbvm\nGep2XdrF1IWBYj+W4Yok0XtDFrBnItLYQQFLjJ+7iEXwm1g2LzPx+JCGbTfj+k6q\n78QH4vovbQcBVYhGTK83tMUuM8Am92wvxR2CEOpWLGj+RHvxhF0db0dSlGUFS6B9\nnNelioefCY6mkdpL3tioEDgSOe1/ceR+ZsEb56oW0Yd1AgHxuuET1SaN7p3Yxr5b\nW4KaqgdZAoGBANdz4D6fnvAKaXPVWw2+BYYryUa4N3A8PNfSWVSUaPVo46bXeTlQ\nZx+QmgZtot6stILygGoGaGTk374WdO/mRyOTfg7VMTvT3mlbidhQnVh3sMmSjybL\nxbS8bwvXxdriB0WEzza5IErhGdX8Yk/lCfzHmq5aHibYH60uaSeFJDztAoGBANdJ\nA8KoyFC2ZtZaZC0bq+NMQkiAfFE5D/ZnjHCQlkbXFZb7RI26MEti4t4MckQUnYok\nOGl9lDUUyWUbojSEm8tXi0vRR3MZ+7k/J1N4zuawPr3BMfSYBQDEFUCZQIs7FT1I\nhQKUxgX5IujJ4xastRlJZx0TontlpkVlpUNoQv6jAoGAf4FrJ4SNqh9vUwbkMRjQ\n6huFrZ2d6YUsuMka5sxB5WKiv31rl3i23t5T2RQPPFrXJVvglV6fb35nz3Y41DTi\nyvIhuyN+VJrJWG69AFCNHesPq+tZXqtfoNuXmFmlFSmJBiJYA1nB+66F/La1c/Tn\nWTrDlwVsLK7g6Du8LZBE5u0CgYEAksjJs5N44O88trHy036md9eq6dwQ5yBM7eg0\nLRuoGqzTn5m6aBemjf/iRxudXSXhNCr1+5cP0hFWL4Xj1oMD5mTOKOeMG8J/ixKw\nMY2RJGDOpnpvISH1Z0xKYT0ccNHb7Wjgp53gVnpDfw0HtJIU+CTAFWcpxZDNCUwA\nnjEcXJkCgYA4Ruf7EgnkPOG3BCkWbz82YO88BAwLtf8El9Cj6Qp0Qaoh6RDu3Nmo\nxALCNoZqXZFtigIiU5pkESAC4nHNhKefFevC7A9o1s6EumoaY+0SIeikyqBV914Y\n2plH67qqDVp903qBcQLebaXckh7KvWQplHaj4JYRz1mvUSLdp8LGwQ==\n-----END RSA PRIVATE KEY-----"
	publicKey                 = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtS/EiblnTN0HePO6jCSY\nPOXaOKN2vl4paxP0dD7NrL3dvJidW8iCNuQUyKkKeYFF54/6rzigjHamG3KhmGy2\nM/GwGLTyaNIEFj8uPW580orJAVcybFXZWg2coXSQwGOrbcsXgk8ZcwviXsIpuXUL\nai5oL7FaQK1l0o7biUt3doGN9k2YSBNDwLFEQplCSE9qoeabKbKqjgzZmTqZV1yv\nIJqGSTTQHQh9tqdjJnKF3iTmGtp7w8yIQ/hjGnKW4585gHMn5y7NQdTBjnlLXbC9\ndHJoNnI66Fw8NNF/GtOkX1fZhbxDxGkdZzv1+OvG4XYfFF7YYpV+1zyZaulit5/w\n5wIDAQAB\n-----END PUBLIC KEY-----"
)

func (s *Server) validateLogin(ctx context.Context, phoneNumber string, pwd string) (user repository.User, valid bool, err error) {
	user, err = s.Repository.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}

		return
	}

	if !password.ComparePassword(pwd, user.Password) {
		return
	}

	return user, true, nil
}

func (s *Server) updateLoginCount(ctx context.Context, user repository.User) (err error) {
	now := time.Now()
	user.LastLogin = &now
	user.LoginCount++
	err = s.Repository.UpdateUser(ctx, user)
	if err != nil {
		return
	}

	return
}

func (s *Server) createAccessToken(ctx context.Context, user repository.User) (accessToken string, err error) {
	claims := map[string]interface{}{
		"id":           user.Id,
		"full_name":    user.FullName,
		"phone_number": user.PhoneNumber,
	}
	accessTokenTtl := time.Duration(accessTokenDuration)
	jwtInput := jwtrsa.GenerateJWTInput{
		PrivateKey:   privateKey,
		Claims:       claims,
		TimeToExpire: accessTokenTtl,
	}

	accessToken, _, err = jwtrsa.GenerateJWT(jwtInput)
	return
}
