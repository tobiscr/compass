package gql

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	gcli "github.com/machinebox/graphql"
)

func NewAuthorizedGraphQLClient(bearerToken string) *gcli.Client {
	return NewAuthorizedGraphQLClientWithCustomURL(bearerToken, getDirectorGraphQLURL())
}

func NewAuthorizedGraphQLClientWithCustomURL(bearerToken, url string) *gcli.Client {
	authorizedClient := newAuthorizedHTTPClient(bearerToken)
	return gcli.NewClient(url, gcli.WithHTTPClient(authorizedClient))
}

func getDirectorGraphQLURL() string {
	url := os.Getenv("DIRECTOR_URL")
	if url == "" {
		url = "http://127.0.0.1:3000"
	}
	url = url + "/graphql"
	return url
}

type authenticatedTransport struct {
	http.Transport
	token string
}

func newAuthorizedHTTPClient(bearerToken string) *http.Client {
	transport := &authenticatedTransport{
		Transport: http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		token: bearerToken,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}
}

func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImM1YjNhMWEyNTA0YzBjMjAwODUxMzg2YjZiNTZkOWQ3NjllYmUwY2QifQ.eyJpc3MiOiJodHRwczovL2RleC5reW1hLmxvY2FsIiwic3ViIjoiQ2lCMFoyZzRaWEEwTnpocVozSjFObVIyTmpoMWRXVnJkR0p4Y21oME5XWndOUklGYkc5allXdyIsImF1ZCI6ImNvbXBhc3MtdWkiLCJleHAiOjE2MDYxOTYwNDMsImlhdCI6MTYwNjE2NzI0MywiYXpwIjoiY29tcGFzcy11aSIsIm5vbmNlIjoiYTYxNWQ4N2MyZGQxNDFlNjkxMGM0ODIwNjVlOWY0NzciLCJhdF9oYXNoIjoiVXlZenNUaE1EZGlCeGY3aTFVQmJzUSIsImVtYWlsIjoiYWRtaW5Aa3ltYS5jeCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiYWRtaW4ifQ.KEz_akD1pR_3ZINRtyvw92d7C2vJLPDOdcOyXrX8WchyU1t79N_2FwvRxJYm7HO-KsMRChRk175nU4tMpfLqqUuFxD36CQ-yJQmdXATGdjsu4a9C061IA2yTQ1RsUufoLi8ncMfeeAHVVD1y79zJv3PmgzWF-D9eQJ8P1HR2xv69fuej-1kzcfJIXytvKkeYqJ1X2JHZwuIDORvuZ8h_19hNz5Er6djVMIMYab9CDxquDuZXavgXE2oYXdwSCBzxgDUYca1z2V6Cof97lh1OpbxqQi-3Krtfo4AbIHa6xsTo3s1zHH_-TbpMwE7qV5JGqARLOv6cKQ4EhvUxOys7dg"))
	req.Header.Set("Tenant", "3e64ebae-38b5-46a0-b1ed-9ccee153a0ae")
	return t.Transport.RoundTrip(req)
}
