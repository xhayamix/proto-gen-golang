package cerrors

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
)

// ActionType エラーアクションタイプ
type ActionType string

const (
	// None 何もしない
	None ActionType = "none"
	// Retry 通信リトライ
	Retry ActionType = "retry"
	// Title タイトルへ遷移
	Title ActionType = "title"
	// ContinueToast 処理継続（ユーザに行動を促す）: 画面遷移せずにエラーメッセージをトースト表示
	ContinueToast ActionType = "continue_toast"
	// ContinueDialog 処理継続（ユーザに行動を促す）: 画面遷移せずにエラーメッセージをダイアログ表示
	ContinueDialog ActionType = "continue_dialog"
)

// ErrorPattern エラーパターン管理
type ErrorPattern struct {
	// エラーコード
	ErrorCode enum.ErrorCode
	// エラーアクションタイプ
	ActionType ActionType
	// エラーロギングスキップ判定
	SkipLogging bool
	// HTTPステータスコード
	HTTPStatusCode int
	// gRPCステータスコード
	GRPCStatusCode codes.Code
}

var (
	// Unknown 予期せぬエラー
	Unknown = ErrorPattern{
		// ErrorCode:      enum.ErrorCode_Unknown,
		ActionType:     Title,
		HTTPStatusCode: http.StatusInternalServerError,
		GRPCStatusCode: codes.Unknown,
	}
	// InvalidArgument 不正なパラメータ
	InvalidArgument = ErrorPattern{
		ErrorCode:      enum.ErrorCode_InvalidArgument,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
	}
	// Internal サーバー内部エラー
	Internal = ErrorPattern{
		ErrorCode:      enum.ErrorCode_Internal,
		ActionType:     Title,
		HTTPStatusCode: http.StatusInternalServerError,
		GRPCStatusCode: codes.Internal,
	}
	// Unauthenticated 認証必要
	Unauthenticated = ErrorPattern{
		ErrorCode:      enum.ErrorCode_Unauthenticated,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusUnauthorized,
		GRPCStatusCode: codes.Unauthenticated,
	}
	// PermissionDenied アクセス拒否
	PermissionDenied = ErrorPattern{
		ErrorCode:      enum.ErrorCode_PermissionDenied,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusForbidden,
		GRPCStatusCode: codes.PermissionDenied,
	}
	// NotFound 見つからない
	NotFound = ErrorPattern{
		ErrorCode:      enum.ErrorCode_NotFound,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
	}
	// UserNotFound ユーザーが見つからない
	UserNotFound = ErrorPattern{
		ErrorCode:      enum.ErrorCode_UserNotFound,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
	}
	// UserDeleted ユーザー削除済み
	UserDeleted = ErrorPattern{
		ErrorCode:      enum.ErrorCode_UserDeleted,
		ActionType:     Retry,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusInternalServerError,
		GRPCStatusCode: codes.Internal,
	}
	// InMaintenance メンテナンス中
	InMaintenance = ErrorPattern{
		ErrorCode:      enum.ErrorCode_InMaintenance,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
	}
	// AccountBan アカウント停止中
	AccountBan = ErrorPattern{
		ErrorCode:      enum.ErrorCode_AccountBan,
		ActionType:     Title,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
	}
	// NgWordContains NGワードが含まれている
	NgWordContains = ErrorPattern{
		ErrorCode:      enum.ErrorCode_NgWordContains,
		ActionType:     ContinueToast,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
	}
	// ShopInvalidBirthday 無効な日付
	ShopInvalidBirthday = ErrorPattern{
		ErrorCode:      enum.ErrorCode_ShopInvalidBirthday,
		ActionType:     ContinueDialog,
		SkipLogging:    true,
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
	}
)

var (
	// PaymentChargeReceiptUsed 課金決済レシートが既に使用済みだった場合のエラー
	PaymentChargeReceiptUsed = Newf(InvalidArgument, "charge receipt is already used")
)

// CampusError サーバ-クライアント間エラーハンドリング用エラー
type CampusError struct {
	// エラーパターン
	ErrorPattern ErrorPattern

	// システムエラーメッセージ
	errorMessage string
	// xerrors拡張用フィールド
	err error
	// それぞれでfmt.Errorf("%w", err)を記述する必要があるためgo1.13でも引き続きxerrors使う。
	frame xerrors.Frame
}

// stackError stacktrace用エラー
type stackError struct {
	*CampusError
}

// New CampusErrorを生成する
func New(errorPattern ErrorPattern) error {
	return newError(nil, errorPattern, "")
}

// Newf CampusErrorを生成する
func Newf(errorPattern ErrorPattern, format string, a ...interface{}) error {
	return newError(nil, errorPattern, fmt.Sprintf(format, a...))
}

// Wrap エラーをCampusエラーでラップする
func Wrap(cause error, errorPattern ErrorPattern) error {
	var message string
	var cerr *CampusError
	if errors.As(cause, &cerr) {
		message = cerr.Message()
	} else {
		message = cause.Error()
	}
	return newError(cause, errorPattern, message)
}

// Wrapf エラーをCampusエラーでラップする
func Wrapf(cause error, errorPattern ErrorPattern, format string, a ...interface{}) error {
	return newError(cause, errorPattern, fmt.Sprintf(format, a...))
}

func As(cause error) (*CampusError, bool) {
	var cerr *CampusError
	if errors.As(cause, &cerr) {
		return cerr, true
	}
	return nil, false
}

func newError(cause error, errorPattern ErrorPattern, errorMessage string) error {
	return &CampusError{
		ErrorPattern: errorPattern,
		errorMessage: errorMessage,
		err:          cause,
		//nolint:gomnd // Magic number
		frame: xerrors.Caller(2),
	}
}

// Stack エラーをStackする
// スタックフレームを明示的に積んでいく必要があるためエラー出力に記録したいエラーハンドリング箇所ではStackを行う
func Stack(err error) error {
	pattern := Unknown
	message := ""
	var campusError *CampusError
	if errors.As(err, &campusError) {
		pattern = campusError.ErrorPattern
		message = campusError.errorMessage
	}
	return &stackError{
		CampusError: &CampusError{
			ErrorPattern: pattern,
			errorMessage: message,
			err:          err,
			frame:        xerrors.Caller(1),
		},
	}
}

// Error エラーメッセージを取得する
func (e *CampusError) Error() string {
	return fmt.Sprintf("error: code = %v, message = %s", e.ErrorPattern.ErrorCode, e.errorMessage)
}

func (e *CampusError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *CampusError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e *CampusError) FormatError(p xerrors.Printer) error {
	p.Print(e.Message())
	e.frame.Format(p)
	return e.Unwrap()
}

func (e *CampusError) Message() string {
	if e == nil {
		return ""
	}
	return e.errorMessage
}

func (e *stackError) FormatError(p xerrors.Printer) error {
	e.frame.Format(p)
	return e.Unwrap()
}
