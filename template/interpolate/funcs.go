package interpolate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/packer/common/uuid"
	"github.com/hashicorp/packer/envibin"
	"github.com/hashicorp/packer/version"
)

// InitTime is the UTC time when this package was initialized. It is
// used as the timestamp for all configuration templates so that they
// match for a single build.
var InitTime time.Time

func init() {
	InitTime = time.Now().UTC()
}

// Funcs are the interpolation funcs that are available within interpolations.
var FuncGens = map[string]FuncGenerator{
	"build_name":     funcGenBuildName,
	"build_type":     funcGenBuildType,
	"env":            funcGenEnv,
	"envibin":        funcGenEnvibin,
	"isotime":        funcGenIsotime,
	"pwd":            funcGenPwd,
	"template_dir":   funcGenTemplateDir,
	"timestamp":      funcGenTimestamp,
	"uuid":           funcGenUuid,
	"user":           funcGenUser,
	"packer_version": funcGenPackerVersion,
	"upper":          funcGenPrimitive(strings.ToUpper),
	"lower":          funcGenPrimitive(strings.ToLower),
}

// FuncGenerator is a function that given a context generates a template
// function for the template.
type FuncGenerator func(*Context) interface{}

// Funcs returns the functions that can be used for interpolation given
// a context.
func Funcs(ctx *Context) template.FuncMap {
	result := make(map[string]interface{})
	for k, v := range FuncGens {
		result[k] = v(ctx)
	}
	if ctx != nil {
		for k, v := range ctx.Funcs {
			result[k] = v
		}
	}

	return template.FuncMap(result)
}

func funcGenBuildName(ctx *Context) interface{} {
	return func() (string, error) {
		if ctx == nil || ctx.BuildName == "" {
			return "", errors.New("build_name not available")
		}

		return ctx.BuildName, nil
	}
}

func funcGenBuildType(ctx *Context) interface{} {
	return func() (string, error) {
		if ctx == nil || ctx.BuildType == "" {
			return "", errors.New("build_type not available")
		}

		return ctx.BuildType, nil
	}
}

func funcGenEnv(ctx *Context) interface{} {
	return func(k string) (string, error) {
		if !ctx.EnableEnv {
			// The error message doesn't have to be that detailed since
			// semantic checks should catch this.
			return "", errors.New("env vars are not allowed here")
		}

		return os.Getenv(k), nil
	}
}

func funcGenEnvibin(ctx *Context) interface{} {
	return func(args ...string) (string, error) {
		if len(args) != 3 {
		}

		if len(args) == 3 {
			repo := args[0]
			image := args[1]
			tag := args[2]

			url, err := envibin.Lookup(repo, image, tag)
			if err != nil {
				return "", err
			}

			return url, nil
		} else if len(args) == 2 {
			repo := ""
			image := args[1]
			tag := args[2]

			url, err := envibin.Lookup(repo, image, tag)
			if err != nil {
				return "", err
			}

			return url, nil
		} else {
			return "", fmt.Errorf("unallowed number of arguments, 2 or 3 arguments required: %v", args)
		}
	}
}

func funcGenIsotime(ctx *Context) interface{} {
	return func(format ...string) (string, error) {
		if len(format) == 0 {
			return InitTime.Format(time.RFC3339), nil
		}

		if len(format) > 1 {
			return "", fmt.Errorf("too many values, 1 needed: %v", format)
		}

		return InitTime.Format(format[0]), nil
	}
}

func funcGenPrimitive(value interface{}) FuncGenerator {
	return func(ctx *Context) interface{} {
		return value
	}
}

func funcGenPwd(ctx *Context) interface{} {
	return func() (string, error) {
		return os.Getwd()
	}
}

func funcGenTemplateDir(ctx *Context) interface{} {
	return func() (string, error) {
		if ctx == nil || ctx.TemplatePath == "" {
			return "", errors.New("template path not available")
		}

		path, err := filepath.Abs(filepath.Dir(ctx.TemplatePath))
		if err != nil {
			return "", err
		}

		return path, nil
	}
}

func funcGenTimestamp(ctx *Context) interface{} {
	return func() string {
		return strconv.FormatInt(InitTime.Unix(), 10)
	}
}

func funcGenUser(ctx *Context) interface{} {
	return func(k string) (string, error) {
		if ctx == nil || ctx.UserVariables == nil {
			return "", errors.New("test")
		}

		return ctx.UserVariables[k], nil
	}
}

func funcGenUuid(ctx *Context) interface{} {
	return func() string {
		return uuid.TimeOrderedUUID()
	}
}

func funcGenPackerVersion(ctx *Context) interface{} {
	return func() string {
		return version.FormattedVersion()
	}
}
