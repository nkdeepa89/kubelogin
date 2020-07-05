package authentication

import (
	"context"

	"github.com/int128/kubelogin/pkg/adaptors/logger"
	"github.com/int128/kubelogin/pkg/adaptors/oidcclient"
	"github.com/int128/kubelogin/pkg/adaptors/reader"
	"golang.org/x/xerrors"
)

// ROPC provides the resource owner password credentials flow.
type ROPC struct {
	Reader reader.Interface
	Logger logger.Interface
}

func (u *ROPC) Do(ctx context.Context, in *ROPCOption, client oidcclient.Interface) (*Output, error) {
	u.Logger.V(1).Infof("performing the resource owner password credentials flow")
	if in.Username == "" {
		var err error
		in.Username, err = u.Reader.ReadString(usernamePrompt)
		if err != nil {
			return nil, xerrors.Errorf("could not read a username: %w", err)
		}
	}
	if in.Password == "" {
		var err error
		in.Password, err = u.Reader.ReadPassword(passwordPrompt)
		if err != nil {
			return nil, xerrors.Errorf("could not read a password: %w", err)
		}
	}
	tokenSet, err := client.GetTokenByROPC(ctx, in.Username, in.Password)
	if err != nil {
		return nil, xerrors.Errorf("resource owner password credentials flow error: %w", err)
	}
	u.Logger.V(1).Infof("resource owner password credentials flow completed")
	return &Output{
		IDToken:       tokenSet.IDToken,
		IDTokenClaims: tokenSet.IDTokenClaims,
		RefreshToken:  tokenSet.RefreshToken,
	}, nil
}
