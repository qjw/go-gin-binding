package binding

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// Bind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// 		"application/json" --> JSON binding
// 		"application/xml"  --> XML binding
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func Bind(c * gin.Context,obj interface{}) (error,[]string) {
	b := Default(c.Request.Method, c.ContentType())
	return BindWith(c, obj, b)
}


func parseError(err error, obj interface{})(error,[]string){
	if err == nil{
		return nil,nil
	}else{
		tips := make([]string, 0)
		real_err,ok := err.(validator.ValidationErrors)
		if ok{
			objt := reflect.TypeOf(obj).Elem()
			if objt.Kind() != reflect.Struct{
				return real_err,tips
			}

			for _, v := range real_err {
				elem,ok := objt.FieldByName(v.StructField())
				if !ok{
					continue
				}
				str,ok := elem.Tag.Lookup("error")
				if ok{
					// log.Printf("tag : %s\n",str)
					tips = append(tips,str)
				}
			}
			return real_err,tips
		}else{
			return err,tips
		}
	}
}

// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
func BindJSON(c * gin.Context,obj interface{}) (error,[]string) {
	return BindWith(c, obj, JSON)
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func BindWith(c * gin.Context,obj interface{}, b Binding) (error,[]string) {
	err := b.Bind(c, obj)
	return parseError(err,obj)
}
