package coinapi

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

/**
  @retry  重试次数
  @method 调用的函数，比如: api.GetTicker ,注意：不是api.GetTicker(...)
  @params 参数,顺序一定要按照实际调用函数入参顺序一样
  @return 返回
*/
func RE(retry int, method interface{}, params ...interface{}) interface{} {

	invokeM := reflect.ValueOf(method)
	if invokeM.Kind() != reflect.Func {
		panic("method not a function")
		return nil
	}

	var value []reflect.Value = make([]reflect.Value, len(params))
	var i int = 0
	for ; i < len(params); i++ {
		value[i] = reflect.ValueOf(params[i])
	}

	var retV interface{}
	var retryC int = 0
_CALL:
	if retryC > 0 {
		log.Println("sleep....", time.Duration(retryC*200*int(time.Millisecond)))
		time.Sleep(time.Duration(retryC * 200 * int(time.Millisecond)))
	}

	retValues := invokeM.Call(value)

	for _, vl := range retValues {
		if vl.Type().String() == "error" {
			if !vl.IsNil() {
				log.Println(vl)
				retryC++
				if retryC <= retry {
					log.Printf("Invoke Method[%s] Error , Begin Retry Call [%d] ...", invokeM.String(), retryC)
					goto _CALL
				} else {
					panic("Invoke Method Fail ???" + invokeM.String())
				}
			}
		} else {
			retV = vl.Interface()
		}
	}

	return retV
}

/**
 * call all unfinished orders
 */
func CancelAllUnfinishedOrders(api API, currencyPair CurrencyPair) int {
	if api == nil {
		log.Println("api instance is nil ??? , please new a api instance")
		return -1
	}

	orders := RE(10, api.GetUnfinishOrders, currencyPair)
	if orders != nil {
		c := 0
		for _, ord := range orders.([]Order) {
			_, err := api.CancelOrder(fmt.Sprintf("%d", ord.OrderID), currencyPair)
			if err != nil {
				log.Println(err)
			}
			c++
			time.Sleep(100 * time.Millisecond) //控制频率
		}

		return c
	}
	return 0
}

/**
 * call all unfinished future orders
 */
func CancelAllUnfinishedFutureOrders(api FutureRestAPI, contractType string, currencyPair CurrencyPair) {
	if api == nil {
		log.Println("api instance is nil ??? , please new a api instance")
		return
	}

	orders := RE(10, api.GetUnfinishFutureOrders, currencyPair, contractType)
	if orders != nil {
		for _, ord := range orders.([]Order) {
			_, err := api.FutureCancelOrder(currencyPair, contractType, fmt.Sprintf("%d", ord.OrderID))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(100 * time.Millisecond) //控制频率
		}
	}
}
