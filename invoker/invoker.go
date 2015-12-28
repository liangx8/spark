package invoker
import(
	"reflect"
	"fmt"
)
/*
type Mapper interface{
	// Maps interface{} value to type
	Map(interface{})
	// Maps first interface{} value ot type of second interface{}
	// Seconde should be a pointer to interface{} otherwise panic
	MapTo(interface{},interface{})
	Set(reflect.Type,reflect.Value)
	Get(reflect.Type)reflect.Value
}
*/
type Invoker interface{
	// Maps interface{} value to type
	Map(interface{})
	// Maps first interface{} value ot type of second interface{}
	// Seconde should be a pointer to interface{} otherwise panic
	MapTo(interface{},interface{})
	Set(reflect.Type,reflect.Value)
	Get(reflect.Type)reflect.Value
	// interface{} should be kind of reflect.Func or panic
	Invoke(interface{}) []reflect.Value
	SetParent(Invoker)
}
type invoker struct{
	args map[reflect.Type]reflect.Value
	parent Invoker
}

func (ink *invoker)Map(a interface{}){
	ink.args[reflect.TypeOf(a)]=reflect.ValueOf(a)
}
func (ink *invoker)MapTo(a interface{},iPtr interface{}){
	val := reflect.ValueOf(a)
	typ := interfaceOf(iPtr)
	t := val.Type()
	if !t.Implements(typ) {
		panic(fmt.Errorf("value of %T can't implements %T ",a,iPtr))
	}
/*
	if !t.ConvertibleTo(typ){
		return fmt.Errorf("value of %T can't convert to %T",a,iPtr)
	}
*/
	ink.args[typ]=val

}
func (ink *invoker)Set(t reflect.Type,v reflect.Value){
	ink.args[t]=v
}
func (ink *invoker)Get(t reflect.Type) reflect.Value{
	val := ink.args[t]
	if val.IsValid(){
		return val
	}
	// no found, try to find implementors
	if t.Kind() == reflect.Interface {
		for k,v := range ink.args {
			if k.Implements(t){
				val = v
				break
			}
		}
	}
	if val.IsValid(){
		return val
	}
	if ink.parent != nil {
		return ink.parent.Get(t)
	}
	return reflect.Zero(t)
}

func (ink *invoker)Invoke(f interface{})[]reflect.Value{
	typ := reflect.TypeOf(f)
	in := make([]reflect.Value,typ.NumIn()) // Painc if typ is not kind of Func
	for i,_ := range in{
		argType := typ.In(i)
		val := ink.Get(argType)
		in[i]=val
	}
	return reflect.ValueOf(f).Call(in)
}
func (ink *invoker)SetParent(p Invoker){
	ink.parent = p
}
func interfaceOf(interfacePtr interface{}) reflect.Type{
	t := reflect.TypeOf(interfacePtr)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic(fmt.Sprintf("Called inject.InterfaceOf with a value that is not a pointer to an interface. %T",interfacePtr))
	}

	return t
}
// New Invoker
func New() Invoker{
	return &invoker{make(map[reflect.Type]reflect.Value),nil}
}
