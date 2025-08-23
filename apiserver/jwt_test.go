package apiserver_test

import (
	"asyncapi/apiserver"
	"asyncapi/config"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJwtManager(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	jwtManager := apiserver.NewJwtManager(conf)
	userId := uuid.New()
	tokenPair, err := jwtManager.GenerateTokenPair(userId)
	require.NoError(t, err)

	require.True(t, jwtManager.IsAccessToken(tokenPair.AccessToken))
	require.False(t, jwtManager.IsAccessToken(tokenPair.RefreshToken))

	accessTokenSubject, err := tokenPair.AccessToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, userId.String(), accessTokenSubject)

	accessTokenIssuer, err := tokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://"+conf.ApiServerHost+":"+conf.ApiServerPort, accessTokenIssuer)

	refreshTokenSubject, err := tokenPair.RefreshToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, userId.String(), refreshTokenSubject)

	refreshTokenIssuer, err := tokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://"+conf.ApiServerHost+":"+conf.ApiServerPort, refreshTokenIssuer)

	parsedAccessToken, err := jwtManager.Parse(tokenPair.AccessToken.Raw)
	require.NoError(t, err)
	require.Equal(t, tokenPair.AccessToken, parsedAccessToken)

	parsedRefreshToken, err := jwtManager.Parse(tokenPair.RefreshToken.Raw)
	require.NoError(t, err)
	require.Equal(t, tokenPair.RefreshToken, parsedRefreshToken)
}
