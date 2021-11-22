package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v40/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var rootCmd = &cobra.Command{
	Use:   "gh-token",
	Short: "Create and revoke GitHub Installation tokens",
	Example: heredoc.Doc(`
		# create a token with an installation ID
		$ gh token create --app-id 123 app-private-key-path path/to/pem --installation-id 123

		# create a token without installation ID
		$ gh token create --app-id 123 app-private-key-path path/to/pem --org org-name

		# revoke a token
		$ gh token revoke --token ghs_123
	`),
}

var createTokenCmd = &cobra.Command{
	Use: "create",
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := cmd.Flags().GetInt64("app-id")
		if err != nil {
			return err
		}

		appPrivKeyPath, err := cmd.Flags().GetString("app-private-key-path")
		if err != nil {
			return err
		}

		installationID, err := cmd.Flags().GetInt64("installation-id")
		if err != nil {
			return err
		}

		appToken, err := createAppToken(appID, appPrivKeyPath)
		if err != nil {
			return err
		}

		ctx := context.Background()

		appClient := createClient(ctx, appToken)

		if installationID == 0 {
			org, err := cmd.Flags().GetString("org")
			if err != nil {
				return err
			}

			installationID, err = findOrgInstallationID(ctx, appClient, org)
		}

		installationToken, err := createInstallationToken(ctx, appClient, installationID)
		if err != nil {
			return err
		}

		fmt.Println(*installationToken.Token)

		return nil
	},
}

var revokeTokenCmd = &cobra.Command{
	Use: "revoke",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}

		return revokeInstallationToken(context.Background(), token)
	},
}

func main() {
	rootCmd.AddCommand(createTokenCmd)
	rootCmd.AddCommand(revokeTokenCmd)

	// create command
	createTokenCmd.Flags().Int64P("app-id", "a", 0, "App ID")
	createTokenCmd.Flags().Int64P("installation-id", "i", 0, "Installation ID")
	createTokenCmd.Flags().StringP("org", "o", "", "Organization name")
	createTokenCmd.Flags().StringP("app-private-key-path", "p", "", "Path to the App Private Key")

	createTokenCmd.MarkFlagRequired("app-id")
	createTokenCmd.MarkFlagRequired("app-private-key-path")

	// revoke command
	revokeTokenCmd.Flags().StringP("token", "t", "", "Installation token to revoke")

	revokeTokenCmd.MarkFlagRequired("token")

	rootCmd.Execute()
}

// Create an authenticated GitHub client
func createClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// Create a JWT token to authenticate as a GitHub App
func createAppToken(appId int64, privKeyPath string) (string, error) {
	privKeyContents, err := os.ReadFile(privKeyPath)
	if err != nil {
		return "", err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyContents)
	if err != nil {
		return "", err
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": appId,
	})

	return token.SignedString(privKey)
}

// Create a token for a GitHub App installation
func createInstallationToken(ctx context.Context, client *github.Client, id int64) (*github.InstallationToken, error) {
	installationToken, _, err := client.Apps.CreateInstallationToken(ctx, id, &github.InstallationTokenOptions{})

	return installationToken, err
}

// Revoke an existing GitHub App installation token
func revokeInstallationToken(ctx context.Context, token string) error {
	client := createClient(ctx, token)

	_, err := client.Apps.RevokeInstallationToken(ctx)

	return err
}

// Find a GitHub App installation for an organization
func findOrgInstallationID(ctx context.Context, client *github.Client, org string) (int64, error) {
	installation, _, err := client.Apps.FindOrganizationInstallation(ctx, org)
	if err != nil {
		return 0, err
	}

	return *installation.ID, nil
}
