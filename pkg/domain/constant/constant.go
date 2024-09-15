package constant

import (
	"github.com/labstack/echo/v4"
	"regexp"
	"time"
)

const (
	// Percent 百分率
	Percent = 100
	// Permil 千分率
	Permil = 1000
	// PermyriadDenominator 万分率分母
	PermyriadDenominator = 10000
	// TenPermyriadDenominator 十万分率分母
	TenPermyriadDenominator = 100000
	// ClockSkewDuration 署名等に使うタイムスタンプの許容誤差
	ClockSkewDuration = 5 * time.Minute

	/* Resettable */

	// DefaultResetHour リセット時間
	DefaultResetHour = 5
	// DefaultResetMinute リセット分
	DefaultResetMinute = 0
	// DefaultResetWeek リセット曜日
	DefaultResetWeek = time.Monday
	// DefaultResetDay リセット日
	DefaultResetDay = 1

	/* headers */

	// HeaderAcceptLanguage accept-language
	HeaderAcceptLanguage = "Accept-Language"
	// HeaderAuthorization authorization
	HeaderAuthorization = "Authorization"
	// HeaderValuePrefixBearer Bearer
	HeaderValuePrefixBearer = "Bearer "

	// HeaderCacheControl cache-control
	HeaderCacheControl = "Cache-Control"
	// HeaderContentType content-type
	HeaderContentType = echo.HeaderContentType
	// HeaderContentLength content-length
	HeaderContentLength = echo.HeaderContentLength
	// HeaderReferer Referer
	HeaderReferer = "Referer"
	// HeaderUserAgent user-agent
	HeaderUserAgent = "User-Agent"
	// HeaderXAppCache x-app-cache
	HeaderXAppCache = "X-App-Cache"
	// HeaderXAppRequestID x-app-request-id
	HeaderXAppRequestID = "X-App-Request-Id"
	// HeaderXAppRequestSignature x-app-request-signature
	HeaderXAppRequestSignature = "X-App-Request-Signature"
	// HeaderXAppRequestTimestamp x-app-request-timestamp
	HeaderXAppRequestTimestamp = "X-App-Request-Timestamp"
	// HeaderXAppVersion x-app-version
	HeaderXAppVersion = "X-App-Version"
	// HeaderXAuthToken x-auth-token
	//nolint:gosec // 認証情報と誤検知しているので無視
	HeaderXAuthToken = "X-Auth-Token"
	// HeaderXCloudTraceContext x-cloud-trace-context
	HeaderXCloudTraceContext = "X-Cloud-Trace-Context"
	// HeaderXDebugTime x-debug-time
	HeaderXDebugTime = "X-Debug-Time"
	// HeaderXErrorAction x-error-action
	HeaderXErrorAction = "X-Error-Action"
	// HeaderXErrorCode x-error-code
	HeaderXErrorCode = "X-Error-Code"
	// HeaderXErrorMessageBin x-error-message-bin
	HeaderXErrorMessageBin = "X-Error-Message-Bin"
	// HeaderXForwardedFor x-forwarded-for
	HeaderXForwardedFor = echo.HeaderXForwardedFor
	// HeaderXForwardedProto x-forwarded-proto
	HeaderXForwardedProto = echo.HeaderXForwardedProto
	// HeaderXGoogAuthenticatedUserEmail x-goog-authenticated-user-email
	HeaderXGoogAuthenticatedUserEmail = "X-Goog-Authenticated-User-Email"
	// HeaderXGoogAuthenticatedUserID x-goog-authenticated-user-id
	HeaderXGoogAuthenticatedUserID = "X-Goog-Authenticated-User-Id"
	// HeaderXGoogIapJwtAssertion x-goog-iap-jwt-assertion
	HeaderXGoogIapJwtAssertion = "X-Goog-Iap-Jwt-Assertion"
	// HeaderXMasterVersion x-master-version
	HeaderXMasterVersion = "X-Master-Version"
	// HeaderXPlatform x-platform
	HeaderXPlatform = "X-Platform"
	// HeaderXServerTime x-server-time
	HeaderXServerTime = "X-Server-Time"
	// HeaderXAdID x-ad-id
	HeaderXAdID = "X-Ad-Id"
	// HeaderXDeviceName x-device-name
	HeaderXDeviceName = "X-Device-Name"
	// HeaderXOSVersion x-os-version
	HeaderXOSVersion = "X-Os-Version"
	// HeaderXIOSIsJailBreak x-i-isj
	HeaderXIOSIsJailBreak = "X-I-Isj"
	// HeaderXIOSIsGameGuardian x-i-isg
	HeaderXIOSIsGameGuardian = "X-I-Isg"
	// HeaderXIOSIsDebuggerAttached x-i-isda
	HeaderXIOSIsDebuggerAttached = "X-I-Isda"
	// HeaderXIOSInjectionList x-i-il
	HeaderXIOSInjectionList = "X-I-Il"
	// HeaderXIOSAppSignature x-i-s
	HeaderXIOSAppSignature = "X-I-S"
	// HeaderXAndroidIsRoot x-a-isr
	HeaderXAndroidIsRoot = "X-A-Isr"
	// HeaderXAndroidIsEmulator x-a-ise
	HeaderXAndroidIsEmulator = "X-A-Ise"
	// HeaderXAndroidIsDebuggerAttached x-a-isda
	HeaderXAndroidIsDebuggerAttached = "X-A-Isda"
	// HeaderXAndroidIsOverlayPermit x-a-isop
	HeaderXAndroidIsOverlayPermit = "X-A-Isop"
	// HeaderXAndroidAppSignature x-a-s
	HeaderXAndroidAppSignature = "X-A-S"
	// HeaderXPlatformUserID platform-user-id
	HeaderXPlatformUserID = "X-Platform-User-Id"
	// HeaderXProduceClearTimeValidationEnable produce-clear-time-validation-enable
	HeaderXProduceClearTimeValidationEnable = "X-Produce-Clear-Time-Validation-Enable"

	/* content types */

	// ContentTypeApplicationGZIP application/gzip
	ContentTypeApplicationGZIP = "application/gzip"
	// ContentTypeApplicationJavaScript application/javascript
	ContentTypeApplicationJavaScript = "application/javascript"
	// ContentTypeApplicationJSON application/json
	ContentTypeApplicationJSON = "application/json"
	// ContentTypeOctetStream application/octet-stream
	ContentTypeOctetStream = "application/octet-stream"
	// ContentTypeZSTD application/zstd
	ContentTypeZSTD = "application/zstd"
	// ContentTypeImageJpeg image/jpeg
	ContentTypeImageJpeg = "image/jpeg"
	// ContentTypeImagePng image/png
	ContentTypeImagePng = "image/png"
	// ContentTypeTextCSS text/css
	ContentTypeTextCSS = "text/css"
	// ContentTypeTextHTML text/html
	ContentTypeTextHTML = "text/html"

	/* extension */

	// ExtensionJpg JPEGファイル
	ExtensionJpg = ".jpg"
	// ExtensionPng PNGファイル
	ExtensionPng = ".png"
	// ExtensionMP4 MP4ファイル
	ExtensionMP4 = ".mp4"
	// ExtensionGob GOBファイル
	ExtensionGob = ".gob"
	// ExtensionZSTD ZSTDファイル
	ExtensionZSTD = ".zst"
	// ExtensionPb pbファイル
	ExtensionPb = ".pb"

	/* time format */

	// TimeFormatNotice 時間フォーマット
	TimeFormatNotice = "2006-01-02 15:04"

	/* hash salt */

	// HashSalt ハッシュソルト値
	HashSalt = "team-d"

	/* public user id */

	// PublicUserIDLength 公開ユーザIDの長さ
	PublicUserIDLength = 8
	// PublicUserIDLetterBytes 公開ユーザIDの候補となる文字（AEROから拝借）
	PublicUserIDLetterBytes = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	// PublicUserIDCandidateNum 1回に生成する公開ユーザーID候補の数
	PublicUserIDCandidateNum = 5
	// PublicUserIDRetryLimit 公開ユーザーID生成の試行回数上限
	PublicUserIDRetryLimit = 10
	// NPCUserIDPrefix NPCユーザのID prefix
	NPCUserIDPrefix = "npc#"

	/* locker */

	// LockTTLRequest リクエストロックのTTL
	LockTTLRequest = 3 * time.Second

	/* ratelimit */

	// RateLimitMaxCount レートリミットの最大カウント
	RateLimitMaxCount = 1000000

	/* Auth */
	// AuthCreateRequestsLimit 制限期間あたりの作成リクエストの回数
	AuthCreateRequestsLimit = -1 // TODO: 制限を有効化する
	// AuthCreateRequestsLimitDuration 作成リクエストの回数を制限する期間
	AuthCreateRequestsLimitDuration = 1 * time.Minute
	// AuthLoginFailureEventName ログイン失敗イベント名
	AuthLoginFailureEventName = "login_failure"
	// AuthLoginFailureLimit 制限期間あたりのログイン失敗回数
	AuthLoginFailureLimit = -1 // TODO: 制限を有効化する
	// AuthLoginFailureLimitDuration ログイン失敗回数を制限する期間
	AuthLoginFailureLimitDuration = 10 * time.Minute
	// IPRequestsLimit 制限期間で許可するIPあたりのリクエスト回数
	IPRequestsLimit = -1 // TODO: 制限を有効化する
	// IPRequestsLimitDuration IPあたりのリクエスト回数を制限する期間
	IPRequestsLimitDuration = 1 * time.Minute
	// GooglePlayIntegrityMaxDecodeRetryCount GooglePlayIntegrityのデコードのリトライ回数上限
	GooglePlayIntegrityMaxDecodeRetryCount = 2

	/* Profile */

	// ProfileNameLengthLimit プロフィールの名前文字数上限
	ProfileNameLengthLimit = 140
)

var (
	// NormalDatetimeRegExp datetime format
	NormalDatetimeRegExp = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{1,2}:\d{2}:\d{2}`)
	// HyphenDatetimeRegExp datetime format
	HyphenDatetimeRegExp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{1,2}:\d{2}:\d{2}`)

	// UnlimitedDatetime 期限なし設定の定数
	UnlimitedDatetime = time.Date(2199, 1, 1, 0, 0, 0, 0, time.Local)
)

type Weekday string

var (
	WeekdaySun Weekday = "Sun"
	WeekdayMon Weekday = "Mon"
	WeekdayThe Weekday = "Tue"
	WeekdayWed Weekday = "Wed"
	WeekdayThu Weekday = "Thu"
	WeekdayFri Weekday = "Fri"
	WeekdaySat Weekday = "Sat"
)

func (w Weekday) Validate() bool {
	switch w {
	case WeekdaySun:
		return true
	case WeekdayMon:
		return true
	case WeekdayThe:
		return true
	case WeekdayWed:
		return true
	case WeekdayThu:
		return true
	case WeekdayFri:
		return true
	case WeekdaySat:
		return true
	}

	return false
}

func (w Weekday) String() string {
	switch w {
	case WeekdaySun:
		return "日"
	case WeekdayMon:
		return "月"
	case WeekdayThe:
		return "火"
	case WeekdayWed:
		return "水"
	case WeekdayThu:
		return "木"
	case WeekdayFri:
		return "金"
	case WeekdaySat:
		return "土"
	}

	return ""
}
