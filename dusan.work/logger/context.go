package logger

import (
	"context"
	"reflect"
)

func ExtractLogger(ctx context.Context) Lgr {
	lgrVal := ctx.Value("logger")
	if lgrVal == nil {
		return Nil
	}
	if reflect.ValueOf(lgrVal).CanConvert(reflect.TypeOf(&Logger{})) {
		return lgrVal.(*Logger)
	}
	return Nil
}

func ExtractLoggerWithKV(ctx context.Context, kvs map[string]string) Lgr {
	lgrVal := ctx.Value("logger")
	if lgrVal == nil {
		return Nil
	}
	if reflect.ValueOf(lgrVal).CanConvert(reflect.TypeOf(&Logger{})) {
		lgr := lgrVal.(*Logger)
		for k, v := range kvs {
			lgr = lgr.WithStr(k, v)
		}
		return lgrVal.(*Logger)
	}
	return Nil
}
