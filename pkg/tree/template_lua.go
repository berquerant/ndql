package tree

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/berquerant/ndql/pkg/cachex"
	"github.com/berquerant/ndql/pkg/util"
	lua "github.com/yuin/gopher-lua"
)

//
// lua gen template
//
// lua(script, entrypoint)
//
// The function specified in entrypoint must be defined within the script,
// take exactly one argument, and return a string.
// The argument passed to this function is a node.
// Additionally, a global table E is defined, which contains the environment variables from os.Environ.

type LuaGenTemplate struct {
	script     string
	entrypoint string
}

func NewLuaGenTemplate(script, entrypoint string) *LuaGenTemplate {
	return &LuaGenTemplate{
		script:     script,
		entrypoint: entrypoint,
	}
}

var _ GenTemplate = &LuaGenTemplate{}

const (
	luaGenTemplateGlobalEnvTable = "E"
	// luaGenTemplateGlobalNodeTable = "N"
)

var luaGenTemplateCache = util.Must(cachex.NewLuaCache())

func (g LuaGenTemplate) Generate(ctx context.Context, n *N) ([]byte, error) {
	proto, err := luaGenTemplateCache.Get(g.script)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to compile lua", errors.Join(ErrGenTemplate, err))
	}

	state := lua.NewState()
	defer state.Close()
	state.SetContext(ctx)
	state.OpenLibs()
	g.registerGlobalTables(state, n)

	state.Push(state.NewFunctionFromProto(proto))
	if err := state.PCall(0, 0, nil); err != nil {
		return nil, fmt.Errorf("%w: failed to call lua proto", errors.Join(ErrGenTemplate, err))
	}

	if err := state.CallByParam(lua.P{
		Fn:      state.GetGlobal(g.entrypoint),
		NRet:    1,
		Protect: true,
	}, g.mapToLTable(state, NodeAsStructuredMap(n))); err != nil {
		return nil, fmt.Errorf("%w: failed to call lua script", errors.Join(ErrGenTemplate, err))
	}

	lRet := state.Get(-1)
	state.Pop(1)
	lStr, ok := lRet.(lua.LString)
	if !ok {
		return nil, fmt.Errorf("%w: failed to retrieve lua return value", ErrGenTemplate)
	}

	return []byte(lStr.String()), nil
}

func (LuaGenTemplate) registerGlobalTables(state *lua.LState, n *N) {
	//
	// register node table
	//
	// nt := state.NewTypeMetatable(luaGenTemplateGlobalNodeTable)
	// state.SetGlobal(luaGenTemplateGlobalNodeTable, nt)
	// state.SetField(nt, "get", state.NewFunction(func(state *lua.LState) int {
	// 	key := state.CheckString(1)
	// 	v := state.OptString(2, "")
	// 	value := genTemplateGetOr(n, key, v)
	// 	state.Push(lua.LString(value))
	// 	return 1
	// }))

	//
	// register env table
	//
	et := state.NewTypeMetatable(luaGenTemplateGlobalEnvTable)
	state.SetGlobal(luaGenTemplateGlobalEnvTable, et)
	state.SetField(et, "get", state.NewFunction(func(state *lua.LState) int {
		key := state.CheckString(1)
		v := state.OptString(2, "")
		value := genTemplateEnvOr(key, v)
		state.Push(lua.LString(value))
		return 1
	}))
}

func (g LuaGenTemplate) mapToLTable(state *lua.LState, d map[string]any) *lua.LTable {
	t := state.NewTable()
	for k, v := range d {
		var lVal lua.LValue
		switch v := v.(type) {
		case float64:
			lVal = lua.LNumber(v)
		case int64:
			lVal = lua.LNumber(v)
		case bool:
			lVal = lua.LBool(v)
		case string:
			lVal = lua.LString(v)
		case time.Time:
			lVal = lua.LString(v.Format(time.DateTime))
		case time.Duration:
			lVal = lua.LString(v.String())
		case map[string]any:
			lVal = g.mapToLTable(state, v)
		default:
			continue
		}
		state.SetField(t, k, lVal)
	}
	return t
}
