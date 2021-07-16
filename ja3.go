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
	tls "github.com/refraction-networking/utls"
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

	config := &tls.Config{}

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
		uConn := tls.UClient(dialConn, config, tls.HelloCustom)

		if err := uConn.ApplyPreset(spec); err != nil {
			return nil, err
		}

		if err := uConn.Handshake(); err != nil {
			return nil, err
		}

		return uConn, nil
	}
}

func stringToSpec(ja3 string) (*tls.ClientHelloSpec, error) {

	var extMap = map[string]tls.TLSExtension{
		"0": &tls.SNIExtension{},
		"5": &tls.StatusRequestExtension{},
		// These are applied later
		// "10": &tls.SupportedCurvesExtension{...}
		// "11": &tls.SupportedPointsExtension{...}
		"13": &tls.SignatureAlgorithmsExtension{
			SupportedSignatureAlgorithms: []tls.SignatureScheme{
				tls.ECDSAWithP256AndSHA256,
				tls.PSSWithSHA256,
				tls.PKCS1WithSHA256,
				tls.ECDSAWithP384AndSHA384,
				tls.PSSWithSHA384,
				tls.PKCS1WithSHA384,
				tls.PSSWithSHA512,
				tls.PKCS1WithSHA512,
				tls.PKCS1WithSHA1,
			},
		},
		"16": &tls.ALPNExtension{
			AlpnProtocols: []string{"h2", "http/1.1"},
		},
		"18": &tls.SCTExtension{},
		// "21": &tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
		"21": &tls.UtlsPaddingExtension{WillPad: true},
		"23": &tls.UtlsExtendedMasterSecretExtension{},
		"27": &tls.FakeCertCompressionAlgsExtension{},
		"28": &tls.FakeRecordSizeLimitExtension{},
		"35": &tls.SessionTicketExtension{},
		"43": &tls.SupportedVersionsExtension{Versions: []uint16{
			tls.GREASE_PLACEHOLDER,
			tls.VersionTLS13,
			tls.VersionTLS12,
			tls.VersionTLS11,
			tls.VersionTLS10}},
		"44": &tls.CookieExtension{},
		"45": &tls.PSKKeyExchangeModesExtension{
			Modes: []uint8{
				tls.PskModeDHE,
			}},
		"51":    &tls.KeyShareExtension{KeyShares: []tls.KeyShare{}},
		"13172": &tls.NPNExtension{},
		"65281": &tls.RenegotiationInfoExtension{
			Renegotiation: tls.RenegotiateOnceAsClient,
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
	var targetCurves []tls.CurveID
	for _, c := range curves {
		cid, err := strconv.ParseUint(c, 10, 16)
		if err != nil {
			return nil, err
		}
		targetCurves = append(targetCurves, tls.CurveID(cid))
	}
	extMap["10"] = &tls.SupportedCurvesExtension{Curves: targetCurves}

	// parse point formats
	var targetPointFormats []byte
	for _, p := range pointFormats {
		pid, err := strconv.ParseUint(p, 10, 8)
		if err != nil {
			return nil, err
		}
		targetPointFormats = append(targetPointFormats, byte(pid))
	}
	extMap["11"] = &tls.SupportedPointsExtension{SupportedPoints: targetPointFormats}

	// build extenions list
	var exts []tls.TLSExtension
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

	return &tls.ClientHelloSpec{
		TLSVersMin:         vid,
		TLSVersMax:         vid,
		CipherSuites:       suites,
		CompressionMethods: []byte{0},
		Extensions:         exts,
		GetSessionID:       sha256.Sum256,
	}, nil
}
