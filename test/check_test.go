package test

import (
	"regexp"
	"testing"

	"github.com/XZ0730/hertz-scaffold/pkg/constants"
	"github.com/cloudwego/kitex/pkg/klog"
)

func TestCheck(t *testing.T) {

	// rege := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"
	r := regexp.MustCompile(constants.CardIdRegexp)
	b := r.MatchString("35042320030730001X")
	klog.Info("bool:", b)

}
