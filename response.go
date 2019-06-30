package happy

import "encoding/json"

func (c *Context) WJson(data interface{})  {
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if paramJson, err := json.Marshal(data); err != nil {
		panic("WJson: Json parsing failed ")
		c.ResponseWriter.WriteHeader(500)
	} else {
		c.ResponseWriter.Write([]byte(paramJson))
	}
}

func (c *Context) WriterText(str string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.ResponseWriter.Write([]byte(str))
}