package scope

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type workloadIdentityCredential struct {
	assertion string
	file      string
	cred      *azidentity.ClientAssertionCredential
	lastRead  time.Time
}

type workloadIdentityCredentialOptions struct {
	azcore.ClientOptions
}

func newWorkloadIdentityCredential(tenantID, clientID, file string, options *workloadIdentityCredentialOptions) (*workloadIdentityCredential, error) {
	w := &workloadIdentityCredential{file: file}
	cred, err := azidentity.NewClientAssertionCredential(tenantID, clientID, w.getAssertion, &azidentity.ClientAssertionCredentialOptions{ClientOptions: options.ClientOptions})
	if err != nil {
		return nil, err
	}
	w.cred = cred
	return w, nil
}

func (w *workloadIdentityCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return w.cred.GetToken(ctx, opts)
}

func (w *workloadIdentityCredential) getAssertion(context.Context) (string, error) {
	if now := time.Now(); w.lastRead.Add(5 * time.Minute).Before(now) {
		content, err := os.ReadFile(w.file)
		if err != nil {
			return "", err
		}
		w.assertion = string(content)
		w.lastRead = now
	}
	return w.assertion, nil
}
