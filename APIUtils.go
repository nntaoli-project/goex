package goex

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

/**
  本函数只适合，返回两个参数的API重试调用，其中一个参数必须是error
  @retry  重试次数
  @delay  每次重试延迟时间间隔
  @method 调用的函数，比如: api.GetTicker ,注意：不是api.GetTicker(...)
  @params 参数,顺序一定要按照实际调用函数入参顺序一样
  @return 返回
*/
func RE(retry int, delay time.Duration, method interface{}, params ...interface{}) interface{} {

	invokeM := reflect.ValueOf(method)
	if invokeM.Kind() != reflect.Func {
		return errors.New("method not a function")
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
		log.Println("sleep ", delay, " after re call")
		time.Sleep(delay)
	}

	retValues := invokeM.Call(value)

	for _, vl := range retValues {
		if vl.Type().String() == "error" {
			if vl.IsNil() {
				continue
			}

			log.Println("[api error]", vl)
			retryC++
			if retryC <= retry-1 {
				log.Printf("Invoke Method[%s] Error , Begin Retry Call [%d] ...", invokeM.String(), retryC)
				goto _CALL
			} else {
				log.Println("Invoke Method Fail ???" + invokeM.String())
				return vl.Interface()
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

	c := 0

	for {
		ret := RE(2, 200*time.Millisecond, api.GetUnfinishOrders, currencyPair)

		if err, isok := ret.(error); !isok {
			log.Println("[api error]", err)
			break
		}

		if ret == nil {
			break
		}

		orders, isok := ret.([]Order)
		if !isok || len(orders) == 0 {
			break
		}

		for _, ord := range orders {
			_, err := api.CancelOrder(ord.OrderID2, currencyPair)
			if err != nil {
				log.Println(err)
			} else {
				c++
			}
			time.Sleep(120 * time.Millisecond) //控制频率
		}
	}

	return c
}

/**
 * call all unfinished future orders
 * @return c 成功撤单数量
 */
func CancelAllUnfinishedFutureOrders(api FutureRestAPI, contractType string, currencyPair CurrencyPair) int {
	if api == nil {
		log.Println("api instance is nil ??? , please new a api instance")
		return 0
	}

	c := 0

	for {
		ret := RE(10, 200*time.Millisecond, api.GetUnfinishFutureOrders, currencyPair, contractType)
		if err, isok := ret.(error); !isok {
			log.Println("[api error]", err)
			break
		}

		if ret == nil {
			break
		}

		orders, isok := ret.([]Order)
		if !isok || len(orders) == 0 {
			break
		}

		for _, ord := range orders {
			_, err := api.FutureCancelOrder(currencyPair, contractType, fmt.Sprintf("%d", ord.OrderID))
			if err != nil {
				log.Println(err)
			} else {
				c++
			}
			time.Sleep(100 * time.Millisecond) //控制频率
		}
	}

	return c
}
