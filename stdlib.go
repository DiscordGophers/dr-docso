package main

import (
	"strings"

	"github.com/hhhapz/doc"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func packageOptions(parts []string, pkg doc.Package) fuzzy.Ranks {
	var opts []string

	for _, f := range pkg.Functions {
		opts = append(opts, f.Name)
	}
	for _, t := range pkg.Types {
		opts = append(opts, t.Name)
		for _, m := range t.Methods {
			opts = append(opts, t.Name+"."+m.Name)
		}
	}

	joined := strings.Join(parts, ".")

	ranks := fuzzy.RankFindFold(joined, opts)
	switch {
	case joined == "":
		ranks = append([]fuzzy.Rank{
			{
				Source: "<package info>",
				Target: "<pkginfo>",
			},
		}, ranks...)

	case len(ranks) == 0:
		ranks = append(ranks, fuzzy.Rank{
			Source: joined,
			Target: joined,
		})
	}
	return ranks
}

func (b *botState) packageCache(query string) fuzzy.Ranks {
	packages := map[string]string{}
	for lib := range stdlib {
		packages[lib] = lib
	}

	for k, v := range b.cfg.Aliases {
		// swap to prevent duplicates with cache
		packages[v] = k
	}

	b.searcher.WithCache(func(cache map[string]*doc.CachedPackage) {
		for k, pkg := range cache {
			if shorthand, ok := packages[k]; ok {
				k = shorthand
			}
			for _, sub := range pkg.Subpackages {
				s := strings.Replace(sub, pkg.URL, k, 1)
				packages[s] = s
			}
		}
	})

	vals := make([]string, 0, len(packages))
	for _, val := range packages {
		if strings.Contains(val, "examples") {
			continue
		}
		vals = append(vals, val)
	}

	ranks := fuzzy.RankFindFold(query, vals)
	if len(ranks) == 0 {
		ranks = append(ranks, fuzzy.Rank{
			Source:        query,
			Target:        query,
			Distance:      0,
			OriginalIndex: -1,
		})
	}
	return ranks
}

var stdlibAliases = map[string]string{
	"tar": "archive/tar",
	"zip": "archive/zip",

	"bzip2": "compress/bzip2",
	"flate": "compress/flate",
	"gzip":  "compress/gzip",
	"lzw":   "compress/lzw",
	"zlib":  "compress/zlib",

	"heap": "container/heap",
	"list": "container/list",
	"ring": "container/ring",

	"aes":      "crypto/aes",
	"cipher":   "crypto/cipher",
	"des":      "crypto/des",
	"dsa":      "crypto/dsa",
	"ecdsa":    "crypto/ecdsa",
	"ed25519":  "crypto/ed25519",
	"elliptic": "crypto/elliptic",
	"hmac":     "crypto/hmac",
	"md5":      "crypto/md5",
	"rc4":      "crypto/rc4",
	"rsa":      "crypto/rsa",
	"sha1":     "crypto/sha1",
	"sha256":   "crypto/sha256",
	"sha512":   "crypto/sha512",
	"subtle":   "crypto/subtle",
	"tls":      "crypto/tls",
	"x509":     "crypto/x509",
	"pkix":     "crypto/x509/pkix",

	"sql": "database/sql",

	"dwarf":    "debug/dwarf",
	"elf":      "debug/elf",
	"gosym":    "debug/gosym",
	"macho":    "debug/macho",
	"pe":       "debug/pe",
	"plan9obj": "debug/plan9obj",

	"ascii85": "encoding/ascii85",
	"asn1":    "encoding/asn1",
	"base32":  "encoding/base32",
	"base64":  "encoding/base64",
	"binary":  "encoding/binary",
	"csv":     "encoding/csv",
	"gob":     "encoding/gob",
	"hex":     "encoding/hex",
	"json":    "encoding/json",
	"pem":     "encoding/pem",
	"xml":     "encoding/xml",

	"ast":           "go/ast",
	"build":         "go/build",
	"constraint":    "go/build/constraint",
	"constant":      "go/constant",
	"docformat":     "go/docformat",
	"importer":      "go/importer",
	"parserprinter": "go/parserprinter",
	"scanner":       "go/scanner",
	"token":         "go/token",
	"types":         "go/types",

	"adler32": "hash/adler32",
	"crc32":   "hash/crc32",
	"crc64":   "hash/crc64",
	"fnv":     "hash/fnv",
	"maphash": "hash/maphash",

	"color":   "image/color",
	"draw":    "image/draw",
	"gif":     "image/gif",
	"jpeg":    "image/jpeg",
	"parsing": "image/parsing",

	"suffixarray": "index/suffixarray",

	"fs":     "io/fs",
	"ioutil": "io/ioutil",

	"big":   "math/big",
	"bits":  "math/bits",
	"cmplx": "math/cmplx",

	"multipart":       "mime/multipart",
	"quotedprintable": "mime/quotedprintable",

	"http":      "net/http",
	"cgi":       "net/http/cgi",
	"cookiejar": "net/http/cookiejar",
	"fcgi":      "net/http/fcgi",
	"httptest":  "net/http/httptest",
	"httptrace": "net/http/httptrace",
	"httputil":  "net/http/httputil",
	"mail":      "net/mail",
	"rpc":       "net/rpc",
	"jsonrpc":   "net/rpc/jsonrpc",
	"smtp":      "net/smtp",
	"textproto": "net/textproto",
	"url":       "net/url",

	"exec":   "os/exec",
	"signal": "os/signal",
	"user":   "os/user",

	"filepath": "path/filepath",

	"syntax": "regexp/syntax",

	"cgo":     "runtime/cgo",
	"metrics": "runtime/metrics",
	"msan":    "runtime/msan",
	"race":    "runtime/race",
	"trace":   "runtime/trace",

	"js": "syscall/js",

	"fstest": "testing/fstest",
	"iotest": "testing/iotest",
	"quick":  "testing/quick",

	"tabwriter": "text/tabwriter",

	"parse": "text/template/parse",

	"tzdata": "time/tzdata",

	"utf16": "unicode/utf16",
	"utf8":  "unicode/utf8",
}

var stdlib = map[string]bool{
	"archive":              true,
	"archive/tar":          true,
	"archive/zip":          true,
	"bufio":                true,
	"builtin":              true,
	"bytes":                true,
	"cmd":                  true,
	"cmd/addr2line":        true,
	"cmd/api":              true,
	"cmd/asm":              true,
	"cmd/buildid":          true,
	"cmd/cgo":              true,
	"cmd/compile":          true,
	"cmd/cover":            true,
	"cmd/dist":             true,
	"cmd/doc":              true,
	"cmd/fix":              true,
	"cmd/go":               true,
	"cmd/gofmt":            true,
	"cmd/link":             true,
	"cmd/nm":               true,
	"cmd/objdump":          true,
	"cmd/pack":             true,
	"cmd/pprof":            true,
	"cmd/test2json":        true,
	"cmd/trace":            true,
	"cmd/vet":              true,
	"compress":             true,
	"compress/bzip2":       true,
	"compress/flate":       true,
	"compress/gzip":        true,
	"compress/lzw":         true,
	"compress/zlib":        true,
	"container":            true,
	"container/heap":       true,
	"container/list":       true,
	"container/ring":       true,
	"context":              true,
	"crypto":               true,
	"crypto/aes":           true,
	"crypto/cipher":        true,
	"crypto/des":           true,
	"crypto/dsa":           true,
	"crypto/ecdsa":         true,
	"crypto/ed25519":       true,
	"crypto/elliptic":      true,
	"crypto/hmac":          true,
	"crypto/md5":           true,
	"crypto/rand":          true,
	"crypto/rc4":           true,
	"crypto/rsa":           true,
	"crypto/sha1":          true,
	"crypto/sha256":        true,
	"crypto/sha512":        true,
	"crypto/subtle":        true,
	"crypto/tls":           true,
	"crypto/x509":          true,
	"crypto/x509/pkix":     true,
	"database":             true,
	"database/sql":         true,
	"database/sql/driver":  true,
	"debug":                true,
	"debug/dwarf":          true,
	"debug/elf":            true,
	"debug/gosym":          true,
	"debug/macho":          true,
	"debug/pe":             true,
	"debug/plan9obj":       true,
	"embed":                true,
	"encoding":             true,
	"encoding/ascii85":     true,
	"encoding/asn1":        true,
	"encoding/base32":      true,
	"encoding/base64":      true,
	"encoding/binary":      true,
	"encoding/csv":         true,
	"encoding/gob":         true,
	"encoding/hex":         true,
	"encoding/json":        true,
	"encoding/pem":         true,
	"encoding/xml":         true,
	"errors":               true,
	"expvar":               true,
	"flag":                 true,
	"fmt":                  true,
	"go":                   true,
	"go/ast":               true,
	"go/build":             true,
	"go/build/constraint":  true,
	"go/constant":          true,
	"go/doc":               true,
	"go/format":            true,
	"go/importer":          true,
	"go/parser":            true,
	"go/printer":           true,
	"go/scanner":           true,
	"go/token":             true,
	"go/types":             true,
	"hash":                 true,
	"hash/adler32":         true,
	"hash/crc32":           true,
	"hash/crc64":           true,
	"hash/fnv":             true,
	"hash/maphash":         true,
	"html":                 true,
	"html/template":        true,
	"image":                true,
	"image/color":          true,
	"image/color/palette":  true,
	"image/draw":           true,
	"image/gif":            true,
	"image/jpeg":           true,
	"image/png":            true,
	"index":                true,
	"index/suffixarray":    true,
	"io":                   true,
	"io/fs":                true,
	"io/ioutil":            true,
	"log":                  true,
	"log/syslog":           true,
	"math":                 true,
	"math/big":             true,
	"math/bits":            true,
	"math/cmplx":           true,
	"math/rand":            true,
	"mime":                 true,
	"mime/multipart":       true,
	"mime/quotedprintable": true,
	"net":                  true,
	"net/http":             true,
	"net/http/cgi":         true,
	"net/http/cookiejar":   true,
	"net/http/fcgi":        true,
	"net/http/httptest":    true,
	"net/http/httptrace":   true,
	"net/http/httputil":    true,
	"net/http/pprof":       true,
	"net/mail":             true,
	"net/rpc":              true,
	"net/rpc/jsonrpc":      true,
	"net/smtp":             true,
	"net/textproto":        true,
	"net/url":              true,
	"os":                   true,
	"os/exec":              true,
	"os/signal":            true,
	"os/user":              true,
	"path":                 true,
	"path/filepath":        true,
	"plugin":               true,
	"reflect":              true,
	"regexp":               true,
	"regexp/syntax":        true,
	"runtime":              true,
	"runtime/cgo":          true,
	"runtime/debug":        true,
	"runtime/metrics":      true,
	"runtime/msan":         true,
	"runtime/pprof":        true,
	"runtime/race":         true,
	"runtime/trace":        true,
	"sort":                 true,
	"strconv":              true,
	"strings":              true,
	"sync":                 true,
	"sync/atomic":          true,
	"syscall":              true,
	"syscall/js":           true,
	"testing":              true,
	"testing/fstest":       true,
	"testing/iotest":       true,
	"testing/quick":        true,
	"text":                 true,
	"text/scanner":         true,
	"text/tabwriter":       true,
	"text/template":        true,
	"text/template/parse":  true,
	"time":                 true,
	"time/tzdata":          true,
	"unicode":              true,
	"unicode/utf16":        true,
	"unicode/utf8":         true,
	"unsafe":               true,

	"x/crypto/acme":              true,
	"x/crypto/acme/autocert":     true,
	"x/crypto/argon2":            true,
	"x/crypto/bcrypt":            true,
	"x/crypto/blake2b":           true,
	"x/crypto/blake2s":           true,
	"x/crypto/blowfish":          true,
	"x/crypto/bn256":             true,
	"x/crypto/cast5":             true,
	"x/crypto/chacha20poly1305":  true,
	"x/crypto/cryptobyte":        true,
	"x/crypto/cryptobyte/asn1":   true,
	"x/crypto/curve25519":        true,
	"x/crypto/ed25519":           true,
	"x/crypto/hkdf":              true,
	"x/crypto/md4":               true,
	"x/crypto/nacl":              true,
	"x/crypto/nacl/auth":         true,
	"x/crypto/nacl/box":          true,
	"x/crypto/nacl/secretbox":    true,
	"x/crypto/nacl/sign":         true,
	"x/crypto/ocsp":              true,
	"x/crypto/openpgp":           true,
	"x/crypto/openpgp/armor":     true,
	"x/crypto/openpgp/clearsign": true,
	"x/crypto/openpgp/elgamal":   true,
	"x/crypto/openpgp/errors":    true,
	"x/crypto/openpgp/packet":    true,
	"x/crypto/openpgp/s2k":       true,
	"x/crypto/otr":               true,
	"x/crypto/pbkdf2":            true,
	"x/crypto/pkcs12":            true,
	"x/crypto/poly1305":          true,
	"x/crypto/ripemd160":         true,
	"x/crypto/salsa20":           true,
	"x/crypto/salsa20/salsa":     true,
	"x/crypto/scrypt":            true,
	"x/crypto/sha3":              true,
	"x/crypto/ssh":               true,
	"x/crypto/ssh/agent":         true,
	"x/crypto/ssh/knownhosts":    true,
	"x/crypto/ssh/terminal":      true,
	"x/crypto/ssh/test":          true,
	"x/crypto/tea":               true,
	"x/crypto/twofish":           true,
	"x/crypto/xtea":              true,
	"x/crypto/xts":               true,

	"x/net/bpf":             true,
	"x/net/context":         true,
	"x/net/context/ctxhttp": true,
	"x/net/dict":            true,
	"x/net/dns":             true,
	"x/net/dns/dnsmessage":  true,
	"x/net/html":            true,
	"x/net/html/atom":       true,
	"x/net/html/charset":    true,
	"x/net/http":            true,
	"x/net/http2":           true,
	"x/net/http2/h2c":       true,
	"x/net/http2/h2i":       true,
	"x/net/http2/hpack":     true,
	"x/net/http/httpguts":   true,
	"x/net/http/httpproxy":  true,
	"x/net/icmp":            true,
	"x/net/idna":            true,
	"x/net/ipv4":            true,
	"x/net/ipv6":            true,
	"x/net/lif":             true,
	"x/net/nettest":         true,
	"x/net/netutil":         true,
	"x/net/proxy":           true,
	"x/net/publicsuffix":    true,
	"x/net/route":           true,
	"x/net/trace":           true,
	"x/net/webdav":          true,
	"x/net/websocket":       true,
	"x/net/xsrftoken":       true,

	"x/oauth2/amazon":            true,
	"x/oauth2/bitbucket":         true,
	"x/oauth2/cern":              true,
	"x/oauth2/clientcredentials": true,
	"x/oauth2/endpoints":         true,
	"x/oauth2/facebook":          true,
	"x/oauth2/fitbit":            true,
	"x/oauth2/foursquare":        true,
	"x/oauth2/github":            true,
	"x/oauth2/gitlab":            true,
	"x/oauth2/google":            true,
	"x/oauth2/heroku":            true,
	"x/oauth2/hipchat":           true,
	"x/oauth2/instagram":         true,
	"x/oauth2/jira":              true,
	"x/oauth2/jws":               true,
	"x/oauth2/jwt":               true,
	"x/oauth2/kakao":             true,
	"x/oauth2/linkedin":          true,
	"x/oauth2/mailchimp":         true,
	"x/oauth2/mailru":            true,
	"x/oauth2/mediamath":         true,
	"x/oauth2/microsoft":         true,
	"x/oauth2/nokiahealth":       true,
	"x/oauth2/odnoklassniki":     true,
	"x/oauth2/paypal":            true,
	"x/oauth2/slack":             true,
	"x/oauth2/spotify":           true,
	"x/oauth2/stackoverflow":     true,
	"x/oauth2/twitch":            true,
	"x/oauth2/uber":              true,
	"x/oauth2/vk":                true,
	"x/oauth2/yahoo":             true,
	"x/oauth2/yandex":            true,

	"x/image/bmp":                           true,
	"x/image/ccitt":                         true,
	"x/image/colornames":                    true,
	"x/image/draw":                          true,
	"x/image/example":                       true,
	"x/image/example/font":                  true,
	"x/image/font":                          true,
	"x/image/font/basicfont":                true,
	"x/image/font/gofont":                   true,
	"x/image/font/gofont/gobold":            true,
	"x/image/font/gofont/gobolditalic":      true,
	"x/image/font/gofont/goitalic":          true,
	"x/image/font/gofont/gomedium":          true,
	"x/image/font/gofont/gomediumitalic":    true,
	"x/image/font/gofont/gomono":            true,
	"x/image/font/gofont/gomonobold":        true,
	"x/image/font/gofont/gomonobolditalic":  true,
	"x/image/font/gofont/gomonoitalic":      true,
	"x/image/font/gofont/goregular":         true,
	"x/image/font/gofont/gosmallcaps":       true,
	"x/image/font/gofont/gosmallcapsitalic": true,
	"x/image/font/gofont/ttfs":              true,
	"x/image/font/inconsolata":              true,
	"x/image/font/opentype":                 true,
	"x/image/font/plan9font":                true,
	"x/image/font/sfnt":                     true,
	"x/image/math":                          true,
	"x/image/math/f32":                      true,
	"x/image/math/f64":                      true,
	"x/image/math/fixed":                    true,
	"x/image/riff":                          true,
	"x/image/tiff":                          true,
	"x/image/tiff/lzw":                      true,
	"x/image/vector":                        true,
	"x/image/vp8":                           true,
	"x/image/vp8l":                          true,
	"x/image/webp":                          true,

	"x/time/rate": true,

	"x/sync/errgroup":     true,
	"x/sync/semaphore":    true,
	"x/sync/singleflight": true,
	"x/sync/syncmap":      true,

	"x/text/cases":                       true,
	"x/text/collate":                     true,
	"x/text/collate/build":               true,
	"x/text/collate/tools":               true,
	"x/text/collate/tools/colcmp":        true,
	"x/text/currency":                    true,
	"x/text/date":                        true,
	"x/text/encoding":                    true,
	"x/text/encoding/charmap":            true,
	"x/text/encoding/htmlindex":          true,
	"x/text/encoding/ianaindex":          true,
	"x/text/encoding/japanese":           true,
	"x/text/encoding/korean":             true,
	"x/text/encoding/simplifiedchinese":  true,
	"x/text/encoding/traditionalchinese": true,
	"x/text/encoding/unicode":            true,
	"x/text/encoding/unicode/utf32":      true,
	"x/text/feature":                     true,
	"x/text/feature/plural":              true,
	"x/text/language":                    true,
	"x/text/language/display":            true,
	"x/text/message":                     true,
	"x/text/message/catalog":             true,
	"x/text/message/pipeline":            true,
	"x/text/number":                      true,
	"x/text/runes":                       true,
	"x/text/search":                      true,
	"x/text/secure":                      true,
	"x/text/secure/bidirule":             true,
	"x/text/secure/precis":               true,
	"x/text/transform":                   true,
	"x/text/unicode":                     true,
	"x/text/unicode/bidi":                true,
	"x/text/unicode/cldr":                true,
	"x/text/unicode/norm":                true,
	"x/text/unicode/rangetable":          true,
	"x/text/unicode/runenames":           true,
	"x/text/width":                       true,

	"x/sys/cpu":                  true,
	"x/sys/plan9":                true,
	"x/sys/unix":                 true,
	"x/sys/unix/linux":           true,
	"x/sys/windows":              true,
	"x/sys/windows/registry":     true,
	"x/sys/windows/svc":          true,
	"x/sys/windows/svc/debug":    true,
	"x/sys/windows/svc/eventlog": true,
	"x/sys/windows/svc/example":  true,
	"x/sys/windows/svc/mgr":      true,
}
