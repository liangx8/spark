package spark

// Distribute is options middleware, to distribe actions by Given Name
// distrubute handler by name string
type Distribute struct{
	name string
	pairs map[string]Handler
}

func (d *Distribute)Bind(name string,h Handler){
	d.pairs[name]=h

}
func NewDistribute(name string) *Distribute{

	return &Distribute{name:name,pairs:make(map[string]Handler)}
}
func (d *Distribute)Handler(c Context,p Params){
	actionName,ok1 := p[d.name]
	if ok1 {
		handler,ok2 := d.pairs[actionName]

		if ok2 {
			vals:=c.Invoke(handler)
			if len(vals)>0{
				c.OnReturn(vals)
			}
			return
		}
	}
	c.Invoke(NotFound)
}
