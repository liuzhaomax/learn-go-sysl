// Copyright 2015, 2018, 2019 Opsmate, Inc. All rights reserved.
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkcs12

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"io"
)

var (
	// see https://tools.ietf.org/html/rfc7292#appendix-D
	oidCertTypeX509Certificate = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 9, 22, 1})
	oidPKCS8ShroundedKeyBag    = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 12, 10, 1, 2})
	oidCertBag                 = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 12, 10, 1, 3})
	oidSecretBag               = asn1.ObjectIdentifier([]int{1, 2, 840, 113549, 1, 12, 10, 1, 5})
)

type certBag struct {
	Id   asn1.ObjectIdentifier
	Data []byte `asn1:"tag:0,explicit"`
}

// https://tools.ietf.org/html/rfc7292#section-4.2.5
// SecretBag ::= SEQUENCE {
//   secretTypeId   BAG-TYPE.&id ({SecretTypes}),
//   secretValue    [0] EXPLICIT BAG-TYPE.&Type ({SecretTypes}
//					   {@secretTypeId})
// }
type secretBag struct {
	SecretTypeID asn1.ObjectIdentifier
	SecretValue  []byte `asn1:"tag:0,explicit"`
}

type secretValue struct {
	AlgorithmIdentifier pkix.AlgorithmIdentifier
	EncryptedContent    []byte
}

func (i secretValue) Algorithm() pkix.AlgorithmIdentifier {
	return i.AlgorithmIdentifier
}

func (i secretValue) Data() []byte { return i.EncryptedContent }

func (i *secretValue) SetData(data []byte) { i.EncryptedContent = data }

// Encode secret key in a pkcs8
// See ftp://ftp.rsasecurity.com/pub/pkcs/pkcs-8/pkcs-8v1_2.asn, RFC 5208,
// https://github.com/openjdk/jdk/blob/jdk8-b120/jdk/src/share/classes/sun/security/pkcs12/PKCS12KeyStore.java#L613,
// https://github.com/openjdk/jdk/blob/jdk9-b94/jdk/src/java.base/share/classes/sun/security/pkcs12/PKCS12KeyStore.java#L624
// and https://github.com/golang/go/blob/master/src/crypto/x509/pkcs8.go
type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
}

func decodePkcs8ShroudedKeyBag(asn1Data, password []byte) (privateKey interface{}, err error) {
	pkinfo := new(encryptedPrivateKeyInfo)
	if err = unmarshal(asn1Data, pkinfo); err != nil {
		return nil, errors.New("pkcs12: error decoding PKCS#8 shrouded key bag: " + err.Error())
	}

	pkData, err := pbDecrypt(pkinfo, password)
	if err != nil {
		return nil, errors.New("pkcs12: error decrypting PKCS#8 shrouded key bag: " + err.Error())
	}

	ret := new(asn1.RawValue)
	if err = unmarshal(pkData, ret); err != nil {
		return nil, errors.New("pkcs12: error unmarshaling decrypted private key: " + err.Error())
	}

	if privateKey, err = x509.ParsePKCS8PrivateKey(pkData); err != nil {
		return nil, errors.New("pkcs12: error parsing PKCS#8 private key: " + err.Error())
	}

	return privateKey, nil
}

func encodePkcs8ShroudedKeyBag(rand io.Reader, privateKey interface{}, password []byte) (asn1Data []byte, err error) {
	var pkData []byte
	if pkData, err = x509.MarshalPKCS8PrivateKey(privateKey); err != nil {
		return nil, errors.New("pkcs12: error encoding PKCS#8 private key: " + err.Error())
	}

	randomSalt := make([]byte, 8)
	if _, err = rand.Read(randomSalt); err != nil {
		return nil, errors.New("pkcs12: error reading random salt: " + err.Error())
	}
	var paramBytes []byte
	if paramBytes, err = asn1.Marshal(pbeParams{Salt: randomSalt, Iterations: 2048}); err != nil {
		return nil, errors.New("pkcs12: error encoding params: " + err.Error())
	}

	var pkinfo encryptedPrivateKeyInfo
	pkinfo.AlgorithmIdentifier.Algorithm = oidPBEWithSHAAnd3KeyTripleDESCBC
	pkinfo.AlgorithmIdentifier.Parameters.FullBytes = paramBytes

	if err = pbEncrypt(&pkinfo, pkData, password); err != nil {
		return nil, errors.New("pkcs12: error encrypting PKCS#8 shrouded key bag: " + err.Error())
	}

	if asn1Data, err = asn1.Marshal(pkinfo); err != nil {
		return nil, errors.New("pkcs12: error encoding PKCS#8 shrouded key bag: " + err.Error())
	}

	return asn1Data, nil
}

func decodeCertBag(asn1Data []byte) (x509Certificates []byte, err error) {
	bag := new(certBag)
	if err := unmarshal(asn1Data, bag); err != nil {
		return nil, errors.New("pkcs12: error decoding cert bag: " + err.Error())
	}
	if !bag.Id.Equal(oidCertTypeX509Certificate) {
		return nil, NotImplementedError("only X509 certificates are supported")
	}
	return bag.Data, nil
}

func encodeCertBag(x509Certificates []byte) (asn1Data []byte, err error) {
	var bag certBag
	bag.Id = oidCertTypeX509Certificate
	bag.Data = x509Certificates
	if asn1Data, err = asn1.Marshal(bag); err != nil {
		return nil, errors.New("pkcs12: error encoding cert bag: " + err.Error())
	}
	return asn1Data, nil
}

func decodeSecretBag(asn1Data []byte, password []byte) (secretKey []byte, err error) {
	bag := new(secretBag)
	if err := unmarshal(asn1Data, bag); err != nil {
		return nil, errors.New("pkcs12: error decoding secret bag: " + err.Error())
	}
	if !bag.SecretTypeID.Equal(oidPKCS8ShroundedKeyBag) {
		return nil, NotImplementedError("only PKCS#8 shrouded key bag secretTypeID are supported")
	}

	val := new(secretValue)
	if err = unmarshal(bag.SecretValue, val); err != nil {
		return nil, errors.New("pkcs12: error decoding secret value: " + err.Error())
	}

	decrypted, err := pbDecrypt(val, password)
	if err != nil {
		return nil, errors.New("pkcs12: error decrypting secret value: " + err.Error())
	}

	s := new(pkcs8)
	if err = unmarshal(decrypted, s); err != nil {
		return nil, errors.New("pkcs12: error unmarshaling decrypted secret key: " + err.Error())
	}
	if s.Version != 0 {
		return nil, NotImplementedError("only secret key v0 are supported")
	}

	return s.PrivateKey, nil
}
