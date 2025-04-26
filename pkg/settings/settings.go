package settings

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/envvar"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	ErrorLogFilename string = "./tmp/log/error.log"

	AccessLogFilename string = "./tmp/log/access.log"

	// Database timeout
	Timeout time.Duration = 15 * time.Second

	// JWT authentication scheme
	BearerAuthScheme string = "Bearer"

	// JWT expiration
	LongExpiration time.Duration = 1 * time.Hour

	ShortExpiration time.Duration = 50 * time.Minute

	// Cache
	CacheControlNoStore string = "no-store"

	// Prefer
	PreferTotalOnly string = "total-only"

	// Verification Code
	VerificationCodeLength int = 6

	VerificationCodeExpiration time.Duration = 5 * time.Minute

	// Token type
	TokenTypeVerification int = 1

	// JSON
	VerificationJSONFilename string = "./json/verification.json"
)

var (
	FiberConfig fiber.Config = fiber.Config{
		Prefork:       false,
		ServerHeader:  "",
		StrictRouting: true, // false,
		CaseSensitive: false,
		Immutable:     false,
		UnescapePath:  true, // false,
		// ETag: false,
		BodyLimit:                    fiber.DefaultBodyLimit,
		Concurrency:                  fiber.DefaultConcurrency,
		Views:                        nil,
		ViewsLayout:                  "",
		PassLocalsToViews:            false,
		ReadTimeout:                  0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		ReadBufferSize:               fiber.DefaultReadBufferSize,
		WriteBufferSize:              fiber.DefaultWriteBufferSize,
		CompressedFileSuffix:         fiber.DefaultCompressedFileSuffix,
		ProxyHeader:                  fiber.HeaderXForwardedFor, // "",
		GETOnly:                      false,
		ErrorHandler:                 fiber.DefaultErrorHandler,
		DisableKeepalive:             false,
		DisableDefaultDate:           false,
		DisableDefaultContentType:    false,
		DisableHeaderNormalizing:     false,
		DisableStartupMessage:        false,
		AppName:                      "",
		StreamRequestBody:            false,
		DisablePreParseMultipartForm: false,
		ReduceMemoryUsage:            false,
		JSONEncoder:                  json.Marshal,
		JSONDecoder:                  json.Unmarshal,
		XMLEncoder:                   xml.Marshal,
		Network:                      fiber.NetworkTCP4,
		EnableTrustedProxyCheck:      true,                  // false,
		TrustedProxies:               []string{"127.0.0.1"}, // []string{},
		EnableIPValidation:           true,                  // false,
		EnablePrintRoutes:            false,
		ColorScheme:                  fiber.DefaultColors,
	}

	CacheConfig cache.Config = cache.Config{
		Next:         nil,
		Expiration:   1 * time.Minute,
		CacheHeader:  "X-Cache",
		CacheControl: true, // false,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
		ExpirationGenerator: nil,
		// Storage: fiber.Storage,
		StoreResponseHeaders: false,
		MaxBytes:             0,
		Methods:              []string{fiber.MethodGet, fiber.MethodHead},
	}

	CompressConfig compress.Config = compress.Config{
		Next:  nil,
		Level: compress.LevelBestSpeed, // compress.LevelDefault,
	}

	CORSConfig cors.Config = cors.Config{
		Next:         nil,
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "",
		AllowCredentials: false,
		// AllowCredentials: true, // Enable third-party cookies
		ExposeHeaders: "",
		MaxAge:        0,
	}

	CSRFConfig csrf.Config = csrf.Config{
		Next:         nil,
		KeyLookup:    "cookie:csrf_", // "header:" + csrf.HeaderName,
		CookieName:   "csrf_",
		CookieDomain: "",
		CookiePath:   "",
		CookieSecure: false,
		// CookieSecure: true, // Enable third-party cookies
		CookieHTTPOnly: true, // false,
		CookieSameSite: fiber.CookieSameSiteLaxMode,
		// CookieSameSite: fiber.CookieSameSiteNoneMode, // Enable third-party cookies
		CookieSessionOnly: false,
		Expiration:        1 * time.Minute, // Should be the same as JWT ShortExpiration // 1 * time.Hour,
		// Storage: fiber.Storage,
		ContextKey:   "csrf_", // "",
		KeyGenerator: utils.UUID,
		// ErrorHandler: fiber.DefaultErrorHandler,
		Extractor: csrf.CsrfFromCookie("csrf_"), // csrf.CsrfFromHeader(csrf.HeaderName),
	}

	EncryptCookieConfig encryptcookie.Config = encryptcookie.Config{
		Next:      nil,
		Except:    []string{"csrf_"},
		Key:       encryptcookie.GenerateKey(),
		Encryptor: encryptcookie.EncryptCookie,
		Decryptor: encryptcookie.DecryptCookie,
	}

	EnvVarConfig envvar.Config = envvar.Config{
		ExportVars:  map[string]string{},
		ExcludeVars: map[string]string{},
	}

	ETagConfig etag.Config = etag.Config{
		Weak: false,
		Next: nil,
	}

	ExpvarConfig expvar.Config = expvar.Config{
		Next: nil,
	}

	LimiterConfig limiter.Config = limiter.Config{
		Next: nil,
		Max:  10, // 5,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		// Storage: fiber.Storage,
		LimiterMiddleware: limiter.FixedWindow{},
	}

	LoggerConfig logger.Config = logger.Config{
		Next:         nil,
		Format:       "${time} ${pid} ${locals:requestid} ${status} ${latency} ${ip}:${port} ${ips} ${method} ${protocol} ${host} ${path} ${queryParams} ${url} ${route} ${error} ${referer} ${ua}\n", // "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "2006/01/02 15:04:05.000000",                                                                                                                                                    // "15:04:05",
		TimeZone:     "UTC",                                                                                                                                                                           // "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	}

	MonitorConfig monitor.Config = monitor.Config{
		Title:      "Fiber Monitor",
		Refresh:    3 * time.Second,
		APIOnly:    false,
		Next:       nil,
		CustomHead: "",
		FontURL:    "https://fonts.googleapis.com/css2?family=Roboto:wght@400;900&display=swap",
		ChartJsURL: "https://cdn.jsdelivr.net/npm/chart.js@2.9/dist/Chart.bundle.min.js",
	}

	PprofConfig pprof.Config = pprof.Config{
		Next: nil,
	}

	RecoverConfig recover.Config = recover.Config{
		Next:             nil,
		EnableStackTrace: true, // false,
	}

	RequestIDConfig requestid.Config = requestid.Config{
		Next:       nil,
		Header:     fiber.HeaderXRequestID,
		Generator:  utils.UUID,
		ContextKey: "requestid",
	}

	SessionConfig session.Config = session.Config{
		Expiration: 5 * time.Minute, // Should be the same as JWT LongExpiration // 24 * time.Hour,
		// Storage: fiber.Storage,
		KeyLookup:    "cookie:session_id",
		CookieDomain: "",
		CookiePath:   "",
		CookieSecure: false,
		// CookieSecure: true, // Enable third-party cookies
		CookieHTTPOnly: true, // false,
		CookieSameSite: fiber.CookieSameSiteLaxMode,
		// CookieSameSite: fiber.CookieSameSiteNoneMode, // Enable third-party cookies
		KeyGenerator: utils.UUIDv4,
	}
)
