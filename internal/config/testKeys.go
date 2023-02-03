package config

// example keys to be used with jwt that can be used locally or for testing
// generated with
//   $ openssl genrsa -out id_rsa 4096
//   $ openssl rsa -in id_rsa -pubout -out id_rsa.pub

// TestPrivateKey private rsa usable for local testing
const TestPrivateKey = `-----BEGIN PRIVATE KEY-----
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
Naleq1zHimsxXew0ha8kaYP559mDlXE=
-----END PRIVATE KEY-----`

// note the public key is not actually password protected, just the correlating private key

// TestPublicKey public rsa usable for local testing that matches TestPrivateKey
const TestPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2RAgaYZrsSH6o4gT4S/v
JbXjeliofK5yzoTyMficwV7j3O68Q+ImoCC3sj/bmZSToAiPBOpXTIw/A6Vixuv9
69KK1f9lBWcBQe2jrdPLITFtc65mMXY2KbhJaStMo2GC62fJRDY1nbJVyjL83WY/
psUx6CjOh4lSiiRgut3RekbDJL1yvbpvu7vxQvcaSEh1QTe2M4N9YCymGQpdyi59
dHi2vlUCZuJ1EHp6bIDEHGuiYPmyfnHL7/okKJD8Va9h1Xqb/Ca310SY9kboWZ3X
jKe8aVk0lH0CyDP7+gxxEiriJDP6XEKRtyK7iD7SGxZtPrO+o/vIMqy4r87dO7pK
CFDjHBN0anBRnbtdmwTi9GYxRhbhvXx+kFJ/2RYJdBKjKXzf3aLTlxBXYSKJ5PyH
2N+JH5KsAV8TV+mj7tgrG/RNmtYYR2jFoODvHcQE2ClytWa47c3V8yWBBIGdKu07
ajf3MZO1g+gxM4vXhzxYMmcZCcT0ma2YUF/InQab1IpKAWkv9FwuYxJOqK1SevKR
YCpHe9tNpcp3hYqhEEyEjlxvIjcY8T8GN2hkVE0NBJfSZ+wI+J3iW0Pyo72QyrPZ
2WnDyXVM+jnOy5JB0w30+VGSffpFwFw5bIcQeayOKf4d2Uj2VtclguV96q1TRtlm
AngA+35tv3Bs6PNEcej5oOECAwEAAQ==
-----END PUBLIC KEY-----`
