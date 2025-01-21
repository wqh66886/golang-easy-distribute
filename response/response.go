package response

/**
* description:
* author: wqh
* date: 2025/1/21
 */

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
