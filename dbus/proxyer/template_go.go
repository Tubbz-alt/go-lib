package main

var __GLOBAL_TEMPLATE_GoLang = `
package {{PkgName}}
import "dlib/dbus"
var __conn *dbus.Conn = nil
func getBus() *dbus.Conn {
	if __conn  == nil {
		var err error
		__conn, err = dbus.{{BusType}}Bus()
		if err != nil {
			panic(err)
		}
	}
	return __conn
}
`

var __IFC_TEMPLATE_INIT_GoLang = `/*This file is auto generate by dlib/dbus/proxyer. Don't edit it*/
package {{PkgName}}
import "dlib/dbus"
import "dlib/dbus/property"
import "reflect"
import "log"
/*prevent compile error*/
var _ = log.Printf
var _ = reflect.TypeOf
var _ = property.BaseObserver{}
`

var __IFC_TEMPLATE_GoLang = `
type {{ExportName}} struct {
	Path dbus.ObjectPath
	core *dbus.Object
	{{if or .Properties .Signals}}signal_chan chan *dbus.Signal{{end}}
	{{range .Properties}}
	{{.Name}} dbus.Property{{end}}
}
{{$obj_name := .Name}}
{{range .Methods }}
func ({{OBJ_NAME}} {{ExportName }}) {{.Name}} ({{GetParamterInsProto .Args}}) ({{GetParamterOutsProto .Args}}) {
	{{OBJ_NAME}}.core.Call("{{$obj_name}}.{{.Name}}", 0{{GetParamterNames .Args}}).Store({{GetParamterOuts .Args}})
	return
}
{{end}}

{{range .Signals}}
func ({{OBJ_NAME}} {{ExportName}}) Connect{{.Name}}(callback func({{GetParamterOutsProto .Args}})) {
	__conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='"+string({{OBJ_NAME}}.core.Path())+"', interface='{{IfcName}}',sender='{{DestName}}',member='{{.Name}}'")
	go func() {
		signals := make(chan *dbus.Signal)
		__conn.Signal(signals)
		for v := range(signals) {
			if v.Name != "{{IfcName}}.{{.Name}}" || {{len .Args}} != len(v.Body) {
				continue
			}
			{{range $index, $arg := .Args}}if reflect.TypeOf(v.Body[0]) != reflect.TypeOf((*{{TypeFor $arg.Type}})(nil)).Elem() {
				continue
			}
			{{end}}

			callback({{range $index, $arg := .Args}}{{if $index}},{{end}}v.Body[{{$index}}].({{TypeFor $arg.Type}}){{end}})
		}
	}()

}
{{end}}

{{range .Properties}}
type dbusProperty{{ExportName}}{{.Name}} struct{
	*property.BaseObserver
	core *dbus.Object
}
{{if PropWritable .}}func (this *dbusProperty{{ExportName}}{{.Name}}) Set(v interface{}/*{{TypeFor .Type}}*/) {
	if reflect.TypeOf(v) == reflect.TypeOf((*{{TypeFor .Type}})(nil)).Elem() {
		this.core.Call("org.freedesktop.DBus.Properties.Set", 0, "{{IfcName}}", "{{.Name}}", dbus.MakeVariant(v))
	} else {
		log.Println("The property {{.Name}} of {{IfcName}} is an {{TypeFor .Type}} but Set with an ", reflect.TypeOf(v))
	}
}{{else}}
func (this *dbusProperty{{ExportName}}{{.Name}}) Set(notwritable interface{}) {
	log.Printf("{{IfcName}}.{{.Name}} is not writable")
}{{end}}
{{ $convert := TryConvertObjectPath . }}
func (this *dbusProperty{{ExportName}}{{.Name}}) Get() interface{} /*{{GetObjectPathType .}}*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "{{IfcName}}", "{{.Name}}").Store(&r)
	if err == nil && r.Signature().String() == "{{.Type}}" { {{ if $convert }}
		before := r.Value().({{TypeFor .Type}})
		{{$convert}}
		return after{{else}}
		return r.Value().({{TypeFor .Type}}){{end}}
	}  else {
		panic(err)
	}
}
func (this *dbusProperty{{ExportName}}{{.Name}}) GetType() reflect.Type {
	return reflect.TypeOf((*{{TypeFor .Type}})(nil)).Elem()
}
{{end}}

func Get{{ExportName}}(path string) *{{ExportName}} {
	core := getBus().Object("{{DestName}}", dbus.ObjectPath(path))
	obj := &{{ExportName}}{Path:dbus.ObjectPath(path), core:core{{if .Signals}},signal_chan:make(chan *dbus.Signal){{end}}}
	{{range .Properties}}
	obj.{{.Name}} = &dbusProperty{{ExportName}}{{.Name}}{&property.BaseObserver{}, core}{{end}}
	{{with .Properties}}
	getBus().BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='signal',path='"+path+"',interface='org.freedesktop.DBus.Properties',sender='{{DestName}}',member='PropertiesChanged'")
	getBus().BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='signal',path='"+path+"',interface='{{IfcName}}',sender='{{DestName}}',member='PropertiesChanged'")
	go func() {
		typeString := reflect.TypeOf("")
		typeKeyValues := reflect.TypeOf(map[string]dbus.Variant{})
		typeArrayValues := reflect.TypeOf([]string{})
		for v := range(obj.signal_chan) {
			if v.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" &&
				len(v.Body) == 3 &&
				reflect.TypeOf(v.Body[0]) == typeString &&
				reflect.TypeOf(v.Body[1]) == typeKeyValues &&
				reflect.TypeOf(v.Body[2]) == typeArrayValues &&
				v.Body[0].(string) != "{{IfcName}}" {
				props := v.Body[1].(map[string]dbus.Variant)
				for key, _ := range props {
					if false { {{range .}}
					} else if key == "{{.Name}}" {
						obj.{{.Name}}.(*dbusProperty{{ExportName}}{{.Name}}).Notify()
					{{end}} }
				}
			} else if v.Name == "{{IfcName}}.PropertiesChanged" && len(v.Body) == 1 && reflect.TypeOf(v.Body[0]) == typeKeyValues {
				for key, _ := range v.Body[0].(map[string]dbus.Variant) {
					if false { {{range .}}
					} else if key == "{{.Name}}" {
						obj.{{.Name}}.(*dbusProperty{{ExportName}}{{.Name}}).Notify()
					{{end}} }
				}
			}
		}
	}()
	{{end}}
	{{if or .Properties .Signals }}getBus().Signal(obj.signal_chan){{end}}
	return obj
}

`

var __TEST_TEMPLATE = `/*This file is auto generate by dlib/dbus/proxyer. Don't edit it*/
package {{PkgName}}
import "testing"
{{range .Methods}}
func Test{{ObjName}}Method{{.Name}} (t *testing.T) {
	{{/*
	rnd := rand.New(rand.NewSource(99))
	r := Get{{ObjName}}("{{TestPath}}").{{.Name}}({{.Args}})
--*/}}

}
{{end}}

{{range .Properties}}
func Test{{ObjName}}Property{{.Name}} (t *testing.T) {
	t.Log("Get the property {{.Name}} of object {{ObjName}} ===> ",
		Get{{ObjName}}("{{TestPath}}").Get{{.Name}}())
}
{{end}}

{{range .Signals}}
func Test{{ObjName}}Signal{{.Name}} (t *testing.T) {
}
{{end}}
`
