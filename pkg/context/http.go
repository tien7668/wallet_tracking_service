package context

import "net/http"

func (c *Context) RespondWith(code int, msg string, data interface{}) {
	c.AbortWithStatusJSON(code, data)
}

func (c *Context) AbortWith400(message string) {
	c.RespondWith(http.StatusBadRequest, message, nil)
}

func (c *Context) AbortWith500(message string) {
	c.RespondWith(http.StatusInternalServerError, message, nil)
}
