package hhttp

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/magisterquis/connectproxy"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
)

type tlsFingerprint struct{ opt *options }

func (tf *tlsFingerprint) JA3(ja3 string) *options {
	tf.opt.ja3DialTLS = tf.ja3DialTLS(ja3)
	return tf.opt
}

func (tf tlsFingerprint) ja3DialTLS(ja3 string) func(network, addr string) (net.Conn, error) {

	spec, err := stringToSpec(ja3)
	if err != nil {
		return func(network, addr string) (net.Conn, error) { return nil, err }
	}

	config := &utls.Config{}

	return func(network, addr string) (net.Conn, error) {

		var dialConn net.Conn

		if tf.opt.proxy != nil {
			var tfProxy string
			switch tf.opt.proxy.(type) {
			case string:
				tfProxy = tf.opt.proxy.(string)
			case []string:
				tfProxy = tf.opt.proxy.([]string)[rand.Intn(len(tf.opt.proxy.([]string)))]
			}

			proxyURL, err := url.Parse(tfProxy)
			if err != nil {
				return nil, err
			}

			var dialer proxy.Dialer
			switch proxyURL.Scheme {
			case "socks5", "socks5h":
				dialer, err = proxy.FromURL(proxyURL, proxy.Direct)
			case "http", "https":
				dialer, err = connectproxy.New(proxyURL, proxy.Direct)
			default:
				return nil, errors.New("proxy: unknown scheme: " + proxyURL.Scheme)
			}
			if err != nil {
				return nil, err
			}

			dialConn, err = dialer.Dial(network, addr)
			if err != nil {
				return nil, err
			}
		} else {
			dialConn, err = net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
		}

		config.ServerName = strings.Split(addr, ":")[0]
		uConn := utls.UClient(dialConn, config, utls.HelloCustom)

		if err := uConn.ApplyPreset(spec); err != nil {
			return nil, err
		}

		if err := uConn.Handshake(); err != nil {
			return nil, err
		}

		return uConn, nil
	}
}

func stringToSpec(ja3 string) (*utls.ClientHelloSpec, error) {

	extMap := map[string]utls.TLSExtension{
		"0": &utls.SNIExtension{},
		"5": &utls.StatusRequestExtension{},
		// These are applied later
		// "10": &tls.SupportedCurvesExtension{...}
		// "11": &tls.SupportedPointsExtension{...}
		"13": &utls.SignatureAlgorithmsExtension{
			SupportedSignatureAlgorithms: []utls.SignatureScheme{
				utls.ECDSAWithP256AndSHA256,
				utls.ECDSAWithP384AndSHA384,
				utls.ECDSAWithP521AndSHA512,
				utls.ECDSAWithSHA1,
				utls.PKCS1WithSHA1,
				utls.PKCS1WithSHA256,
				utls.PKCS1WithSHA384,
				utls.PKCS1WithSHA512,
				utls.PSSWithSHA256,
				utls.PSSWithSHA384,
				utls.PSSWithSHA512,
			},
		},
		"16": &utls.ALPNExtension{
			AlpnProtocols: []string{"h2", "http/1.1"},
		},
		"18": &utls.SCTExtension{},
		"21": &utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle},
		"23": &utls.UtlsExtendedMasterSecretExtension{},
		"27": &utls.FakeCertCompressionAlgsExtension{},
		"28": &utls.FakeRecordSizeLimitExtension{},
		"35": &utls.SessionTicketExtension{},
		"43": &utls.SupportedVersionsExtension{Versions: []uint16{
			utls.GREASE_PLACEHOLDER,
			utls.VersionTLS13,
			utls.VersionTLS12,
			utls.VersionTLS11,
			utls.VersionTLS10,
		}},
		"44": &utls.CookieExtension{},
		"45": &utls.PSKKeyExchangeModesExtension{Modes: []uint8{
			utls.PskModeDHE,
		}},
		"51": &utls.KeyShareExtension{KeyShares: []utls.KeyShare{
			{Group: utls.X25519},
			{Group: utls.CurveP256},
		}},
		"13172": &utls.NPNExtension{},
		"65281": &utls.RenegotiationInfoExtension{
			Renegotiation: utls.RenegotiateOnceAsClient,
		},
	}

	tokens := strings.Split(ja3, ",")

	if len(tokens) < 5 {
		return nil, errors.New("bad JA3 client fingerprint")
	}

	version := tokens[0]
	ciphers := strings.Split(tokens[1], "-")
	extensions := strings.Split(tokens[2], "-")
	curves := strings.Split(tokens[3], "-")
	if len(curves) == 1 && curves[0] == "" {
		curves = []string{}
	}
	pointFormats := strings.Split(tokens[4], "-")
	if len(pointFormats) == 1 && pointFormats[0] == "" {
		pointFormats = []string{}
	}

	// parse curves
	var targetCurves []utls.CurveID
	for _, c := range curves {
		cid, err := strconv.ParseUint(c, 10, 16)
		if err != nil {
			return nil, err
		}
		targetCurves = append(targetCurves, utls.CurveID(cid))
	}
	extMap["10"] = &utls.SupportedCurvesExtension{Curves: targetCurves}

	// parse point formats
	var targetPointFormats []byte
	for _, p := range pointFormats {
		pid, err := strconv.ParseUint(p, 10, 8)
		if err != nil {
			return nil, err
		}
		targetPointFormats = append(targetPointFormats, byte(pid))
	}
	extMap["11"] = &utls.SupportedPointsExtension{SupportedPoints: targetPointFormats}

	// build extenions list
	var exts []utls.TLSExtension
	for _, e := range extensions {
		te, ok := extMap[e]
		if !ok {
			return nil, fmt.Errorf("extension does not exist: %s", string(e))
		}
		exts = append(exts, te)
	}
	// build SSLVersion
	vid64, err := strconv.ParseUint(version, 10, 16)
	if err != nil {
		return nil, err
	}
	vid := uint16(vid64)

	// build CipherSuites
	var suites []uint16
	for _, c := range ciphers {
		cid, err := strconv.ParseUint(c, 10, 16)
		if err != nil {
			return nil, err
		}
		suites = append(suites, uint16(cid))
	}

	return &utls.ClientHelloSpec{
		TLSVersMin:         vid,
		TLSVersMax:         vid,
		CipherSuites:       suites,
		CompressionMethods: []byte{0},
		Extensions:         exts,
		GetSessionID:       sha256.Sum256,
	}, nil
}
