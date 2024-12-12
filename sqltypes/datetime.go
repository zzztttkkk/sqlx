package sqltypes

import (
	"database/sql"
	"time"
	"unsafe"
)

type datetimeTypeBuilder struct {
	typecommonBuilder[time.Time, datetimeTypeBuilder]
}

func Datetime(ptr *time.Time) *datetimeTypeBuilder {
	ins := &datetimeTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype("datetime")
	return ins
}

func NullableDatetime(ptr *sql.Null[time.Time]) *datetimeTypeBuilder {
	return Datetime((*time.Time)(unsafe.Pointer(ptr))).Nullable()
}
