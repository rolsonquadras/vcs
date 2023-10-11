/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc4vci

import "github.com/trustbloc/vcs/pkg/doc/verifiable"

type JWTProofClaims struct {
	Issuer   string `json:"iss,omitempty"`
	Audience string `json:"aud,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
}

type CredentialRequest struct {
	Format string   `json:"format,omitempty"`
	Types  []string `json:"types"`
	Proof  JWTProof `json:"proof,omitempty"`
}

type JWTProof struct {
	JWT       string `json:"jwt"`
	ProofType string `json:"proof_type"`
}

type CredentialResponse struct {
	AcceptanceToken string                `json:"acceptance_token,omitempty"`
	CNonce          string                `json:"c_nonce,omitempty"`
	CNonceExpiresIn int                   `json:"c_nonce_expires_in,omitempty"`
	Credential      interface{}           `json:"credential"`
	Format          verifiable.OIDCFormat `json:"format"`
}