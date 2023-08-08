/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"fmt"
	"strings"
	"time"

	"github.com/piprate/json-gold/ld"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	ariessigner "github.com/hyperledger/aries-framework-go/pkg/doc/signature/signer"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/bbsblssignature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/jsonwebsignature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"

	"github.com/trustbloc/vcs/pkg/doc/vc"
	"github.com/trustbloc/vcs/pkg/doc/vc/jws"
	vcsverifiable "github.com/trustbloc/vcs/pkg/doc/verifiable"
	"github.com/trustbloc/vcs/pkg/internal/common/diddoc"
)

const (
	// Ed25519VerificationKey2018 ed25119 verification type.
	Ed25519VerificationKey2018 = "Ed25519VerificationKey2018"
	// Ed25519VerificationKey2020 ed25119 verification type.
	Ed25519VerificationKey2020 = "Ed25519VerificationKey2020"
	// JSONWebKey2020 verification type.
	JSONWebKey2020 = "JsonWebKey2020"
	// EcdsaSecp256k1VerificationKey2019 verification type.
	EcdsaSecp256k1VerificationKey2019 = "EcdsaSecp256k1VerificationKey2019"
	// Bls12381G2Key2020 verification type.
	Bls12381G2Key2020 = "Bls12381G2Key2020"
)

const (
	// Ed25519KeyType ed25519 key type.
	Ed25519KeyType = "Ed25519"

	// P256KeyType EC P-256 key type.
	P256KeyType = "P256"
)

const (
	// AssertionMethod assertionMethod.
	AssertionMethod = "assertionMethod"

	// Authentication authentication.
	Authentication = "authentication"

	// CapabilityDelegation capabilityDelegation.
	CapabilityDelegation = "capabilityDelegation"

	// CapabilityInvocation capabilityInvocation.
	CapabilityInvocation = "capabilityInvocation"
)

const (
	// Purpose is the key of verifiable.Proof.
	Purpose = "proofPurpose"

	// VerificationMethod is the key of verifiable.Proof.
	VerificationMethod = "verificationMethod"
)

type keyManager interface {
	NewVCSigner(creator string, signatureType vcsverifiable.SignatureType) (vc.SignerAlgorithm, error)
}

// New return new instance of vc crypto.
func New(vdr vdrapi.Registry, loader ld.DocumentLoader) *Crypto {
	return &Crypto{vdr: vdr, documentLoader: loader}
}

// signingOpts holds options for the signing credential.
type signingOpts struct {
	VerificationMethod string
	Purpose            string
	Representation     string
	SignatureType      vcsverifiable.SignatureType
	Created            *time.Time
	Challenge          string
	Domain             string
	SDJWTTemplateData  *vc.SelectiveDisclosureTemplate
}

// SigningOpts is signing credential option.
type SigningOpts func(opts *signingOpts)

// WithVerificationMethod is an option to pass verification method for signing.
func WithVerificationMethod(verificationMethod string) SigningOpts {
	return func(opts *signingOpts) {
		opts.VerificationMethod = verificationMethod
	}
}

func WithSDJWTTemplateData(template *vc.SelectiveDisclosureTemplate) SigningOpts {
	return func(opts *signingOpts) {
		opts.SDJWTTemplateData = template
	}
}

// WithPurpose is an option to pass proof purpose option for signing.
func WithPurpose(purpose string) SigningOpts {
	return func(opts *signingOpts) {
		opts.Purpose = purpose
	}
}

// WithSigningRepresentation is an option to pass representation for signing.
func WithSigningRepresentation(representation string) SigningOpts {
	return func(opts *signingOpts) {
		opts.Representation = representation
	}
}

// WithSignatureType is an option to pass signature type for signing.
func WithSignatureType(signatureType vcsverifiable.SignatureType) SigningOpts {
	return func(opts *signingOpts) {
		opts.SignatureType = signatureType
	}
}

// WithCreated is an option to pass created time option for signing.
func WithCreated(created *time.Time) SigningOpts {
	return func(opts *signingOpts) {
		opts.Created = created
	}
}

// WithChallenge proof challenge.
func WithChallenge(challenge string) SigningOpts {
	return func(opts *signingOpts) {
		opts.Challenge = challenge
	}
}

// WithDomain proof domain.
func WithDomain(domain string) SigningOpts {
	return func(opts *signingOpts) {
		opts.Domain = domain
	}
}

// Crypto to sign credential.
type Crypto struct {
	vdr            vdrapi.Registry
	documentLoader ld.DocumentLoader
}

func (c *Crypto) SignCredential(
	signerData *vc.Signer, vc *verifiable.Credential, opts ...SigningOpts) (*verifiable.Credential, error) {
	switch signerData.Format {
	case vcsverifiable.Jwt:
		return c.signCredentialJWT(signerData, vc, opts...)
	case vcsverifiable.Ldp:
		return c.signCredentialLDP(signerData, vc, opts...)
	default:
		return nil, fmt.Errorf("unknown signature format %s", signerData.Format)
	}
}

// signCredentialLDP adds verifiable.LinkedDataProofContext to the VC.
func (c *Crypto) signCredentialLDP(
	signerData *vc.Signer, vc *verifiable.Credential, opts ...SigningOpts) (*verifiable.Credential, error) {
	signOpts := &signingOpts{}
	// apply opts
	for _, opt := range opts {
		opt(signOpts)
	}

	signatureType := signerData.SignatureType
	if signOpts.SignatureType != "" {
		signatureType = signOpts.SignatureType
	}

	if signOpts.Purpose == "" {
		signOpts.Purpose = Authentication
	}

	signingCtx, err := c.getLinkedDataProofContext(signerData, signerData.KMS, signatureType, Authentication,
		signerData.SignatureRepresentation, signOpts)
	if err != nil {
		return nil, err
	}

	err = vc.AddLinkedDataProof(signingCtx, jsonld.WithDocumentLoader(c.documentLoader))
	if err != nil {
		return nil, fmt.Errorf("failed to sign vc: %w", err)
	}

	return vc, nil
}

// signCredentialJWT returns vc in JWT format including the signature section.
func (c *Crypto) signCredentialJWT(
	signerData *vc.Signer,
	credential *verifiable.Credential,
	opts ...SigningOpts) (*verifiable.Credential, error) {
	signOpts := &signingOpts{}
	// apply opts
	for _, opt := range opts {
		opt(signOpts)
	}

	signatureType := signerData.SignatureType
	if signOpts.SignatureType != "" {
		signatureType = signOpts.SignatureType
	}

	s, _, err := c.getSigner(signerData.KMSKeyID, signerData.KMS, signOpts, signatureType)
	if err != nil {
		return nil, fmt.Errorf("getting signer for JWS: %w", err)
	}

	method := signerData.Creator

	didDoc, err := diddoc.GetDIDDocFromVerificationMethod(method, c.vdr)
	if err != nil {
		return nil, fmt.Errorf("unable to get did doc from verification method %w", err)
	}

	proofPurpose := Authentication
	if signOpts.Purpose != "" {
		proofPurpose = signOpts.Purpose
	}

	err = ValidateProofPurpose(proofPurpose, method, didDoc)
	if err != nil {
		return nil, fmt.Errorf("ValidateProofPurpose error: %w", err)
	}

	jwsAlgo, err := verifiable.KeyTypeToJWSAlgo(signerData.KeyType)
	if err != nil {
		return nil, fmt.Errorf("getting JWS algo based on signature type: %w", err)
	}

	if signerData.SDJWT.Enable {
		options := []verifiable.MakeSDJWTOption{
			verifiable.MakeSDJWTWithHash(signerData.SDJWT.HashAlg),
		}

		if signOpts.SDJWTTemplateData != nil {
			options = append(options,
				verifiable.MakeSDJWTWithVersion(signOpts.SDJWTTemplateData.Version),
				verifiable.MakeSDJWTWithRecursiveClaimsObjects(signOpts.SDJWTTemplateData.RecursiveClaims),
				verifiable.MakeSDJWTWithAlwaysIncludeObjects(signOpts.SDJWTTemplateData.AlwaysInclude),
				verifiable.MakeSDJWTWithNonSelectivelyDisclosableClaims(
					signOpts.SDJWTTemplateData.NonSelectivelyDisclosable),
			)
		}

		return c.getSDJWTSignedCredential(credential, s, jwsAlgo, method, options...)
	}

	return c.getJWTSignedCredential(credential, s, jwsAlgo, method)
}

func (c *Crypto) getJWTSignedCredential(
	credential *verifiable.Credential,
	signer vc.SignerAlgorithm,
	jwsAlgo verifiable.JWSAlgorithm,
	signingKeyID string) (*verifiable.Credential, error) {
	claims, err := credential.JWTClaims(false)
	if err != nil {
		return nil, fmt.Errorf("creating JWT claims for VC: %w", err)
	}

	jwt, err := claims.MarshalJWS(jwsAlgo, signer, signingKeyID)
	if err != nil {
		return nil, fmt.Errorf("MarshalJWS error: %w", err)
	}

	credential.JWT = jwt

	return credential, nil
}

func (c *Crypto) getSDJWTSignedCredential(
	credential *verifiable.Credential,
	signer vc.SignerAlgorithm,
	jwsAlgo verifiable.JWSAlgorithm,
	signingKeyID string,
	options ...verifiable.MakeSDJWTOption,
) (*verifiable.Credential, error) {
	jwsAlgName, err := jwsAlgo.Name()
	if err != nil {
		return nil, fmt.Errorf("getting JWS algo name error: %w", err)
	}

	joseSigner := jws.NewSigner(signingKeyID, jwsAlgName, signer)

	//
	sdjwt, err := credential.MakeSDJWT(joseSigner, signingKeyID, options...)
	if err != nil {
		return nil, fmt.Errorf("make SDJWT credential error: %w", err)
	}

	credential.JWT = sdjwt

	return credential, nil
}

// SignPresentation signs a presentation.
func (c *Crypto) SignPresentation(signerData *vc.Signer, vp *verifiable.Presentation,
	opts ...SigningOpts) (*verifiable.Presentation, error) {
	signOpts := &signingOpts{}
	// apply opts
	for _, opt := range opts {
		opt(signOpts)
	}

	signatureType := signerData.SignatureType
	if signOpts.SignatureType != "" {
		signatureType = signOpts.SignatureType
	}

	signingCtx, err := c.getLinkedDataProofContext(
		signerData, signerData.KMS, signatureType, Authentication, signerData.SignatureRepresentation, signOpts)
	if err != nil {
		return nil, err
	}

	if signingCtx.Purpose == "" {
		signingCtx.Purpose = Authentication
	}

	err = vp.AddLinkedDataProof(signingCtx, jsonld.WithDocumentLoader(c.documentLoader))
	if err != nil {
		return nil, fmt.Errorf("failed to sign vc: %w", err)
	}

	return vp, nil
}

func (c *Crypto) getLinkedDataProofContext(signerData *vc.Signer, km keyManager,
	signatureType vcsverifiable.SignatureType, proofPurpose string,
	signRep verifiable.SignatureRepresentation, opts *signingOpts) (*verifiable.LinkedDataProofContext, error) {
	s, _, err := c.getSigner(signerData.KMSKeyID, km, opts, signatureType)
	if err != nil {
		return nil, err
	}

	if opts.Purpose != "" {
		proofPurpose = opts.Purpose
	}

	method := signerData.Creator

	didDoc, err := diddoc.GetDIDDocFromVerificationMethod(method, c.vdr)
	if err != nil {
		return nil, err
	}

	err = ValidateProofPurpose(proofPurpose, method, didDoc)
	if err != nil {
		return nil, err
	}

	var signatureSuite ariessigner.SignatureSuite

	switch signatureType { //nolint: exhaustive
	case vcsverifiable.Ed25519Signature2018:
		signatureSuite = ed25519signature2018.New(suite.WithSigner(s))
	case vcsverifiable.Ed25519Signature2020:
		signatureSuite = ed25519signature2020.New(suite.WithSigner(s))
	case vcsverifiable.JSONWebSignature2020:
		signatureSuite = jsonwebsignature2020.New(suite.WithSigner(s))
	case vcsverifiable.BbsBlsSignature2020:
		signatureSuite = bbsblssignature2020.New(suite.WithSigner(s))
	case vcsverifiable.EcdsaSecp256k1Signature2019:
		signatureSuite = ecdsasecp256k1signature2019.New(suite.WithSigner(s))
	default:
		return nil, fmt.Errorf("signature type unsupported %s", signatureType)
	}

	if opts.Representation != "" {
		signRep, err = getSignatureRepresentation(opts.Representation)
		if err != nil {
			return nil, err
		}
	}

	vm := method

	signingCtx := &verifiable.LinkedDataProofContext{
		VerificationMethod:      vm,
		SignatureRepresentation: signRep,
		SignatureType:           signatureType.Name(),
		Suite:                   signatureSuite,
		Purpose:                 opts.Purpose,
		Created:                 opts.Created,
		Challenge:               opts.Challenge,
		Domain:                  opts.Domain,
	}

	return signingCtx, nil
}

// getSigner returns signer and verification method based on profile and signing opts
// verificationMethod from opts takes priority to create signer and verification method.
//
//nolint:unparam
func (c *Crypto) getSigner(kmsKeyID string, km keyManager, opts *signingOpts,
	signatureType vcsverifiable.SignatureType) (vc.SignerAlgorithm, string, error) {
	s, err := km.NewVCSigner(kmsKeyID, signatureType)

	return s, kmsKeyID, err
}

// ValidateProofPurpose validates the proof purpose.
func ValidateProofPurpose(proofPurpose, method string, didDoc *did.Doc) error {
	// TODO https://github.com/trustbloc/vcs/issues/368 remove check once did:sov returns both
	//  assertionMethod and authentication
	if strings.Contains(method, "did:sov") {
		return nil
	}

	var vmMatched bool

	switch proofPurpose {
	case AssertionMethod:
		assertionMethods := didDoc.VerificationMethods(did.AssertionMethod)[did.AssertionMethod]

		vmMatched = isValidVerificationMethod(method, assertionMethods)
	case Authentication:
		authMethods := didDoc.VerificationMethods(did.Authentication)[did.Authentication]

		vmMatched = isValidVerificationMethod(method, authMethods)
	case CapabilityDelegation:
		capabilityDelegationMethods := didDoc.VerificationMethods(did.CapabilityDelegation)[did.CapabilityDelegation]

		vmMatched = isValidVerificationMethod(method, capabilityDelegationMethods)
	case CapabilityInvocation:
		capabilityInvocationMethods := didDoc.VerificationMethods(did.CapabilityInvocation)[did.CapabilityInvocation]

		vmMatched = isValidVerificationMethod(method, capabilityInvocationMethods)
	default:
		return fmt.Errorf("proof purpose %s not supported", proofPurpose)
	}

	if !vmMatched {
		return fmt.Errorf("unable to find matching %s key IDs for given verification method %s",
			proofPurpose, method)
	}

	return nil
}

func isValidVerificationMethod(method string, vms []did.Verification) bool {
	for _, vm := range vms {
		if method == vm.VerificationMethod.ID {
			return true
		}
	}

	return false
}

// getSignatureRepresentation returns signing repsentation for given representation key.
func getSignatureRepresentation(signRep string) (verifiable.SignatureRepresentation, error) {
	var signatureRepresentation verifiable.SignatureRepresentation

	switch signRep {
	case "jws":
		signatureRepresentation = verifiable.SignatureJWS
	case "proofValue":
		signatureRepresentation = verifiable.SignatureProofValue
	default:
		return -1, fmt.Errorf("invalid proof format : %s", signRep)
	}

	return signatureRepresentation, nil
}
