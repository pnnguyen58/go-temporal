package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretSuite struct {
	SecretName string
	Region     string
}

func getSecret(ctx context.Context, secret SecretSuite) ([]byte, error) {
	//Create a Secrets Manager client
	var (
		err       error
		sess      *session.Session
		result    *secretsmanager.GetSecretValueOutput
		resultLen int
	)
	sess, err = session.NewSession()
	if err != nil {
		// Handle session creation error
		fmt.Println(err.Error())
		return nil, err
	}
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(secret.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret.SecretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err = svc.GetSecretValueWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		return []byte(*result.SecretString), nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		resultLen, err = base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return nil, err
		}
		return decodedBinarySecretBytes[:resultLen], err
	}
}
