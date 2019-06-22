# bulrush-mq
```go
// @Summary MQ测试
// @Description MQ测试
// @Tags MQ
// @Accept mpfd
// @Produce json
// @Param accessToken query string true "令牌"
// @Success 200 {string} json "{"message": "ok"}"
// @Failure 400 {string} json "{"message": "failure"}"
// @Router /mgoTest [get]
func mqHello(router *gin.RouterGroup, event events.EventEmmiter, q *mq.MQ) {
	q.Register("mqTest", func(mess mq.Message) error {
		fmt.Println(mess.Body)
		return nil
	})
	router.GET("/mqTest", func(c *gin.Context) {
		q.Push(mq.Message{
			Type: "mqTest",
			Body: map[string]interface{}{
				"message": "ok",
			},
		})
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}
// RegisterMq defined hello routes
func RegisterMq(router *gin.RouterGroup, event events.EventEmmiter, ri *bulrush.ReverseInject) {
	ri.Register(mqHello)
}
```
## MIT License

Copyright (c) 2018-2020 Double

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.