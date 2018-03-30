package api

import (
	"github.com/9chain/nbcapid/apikey"
	"github.com/9chain/nbcapid/primitives"
	"github.com/gin-gonic/gin"
)

type JSONError = primitives.JSONError
type JSON2Request = primitives.JSON2Request
type JSON2Response = primitives.JSON2Response

var (
	handlers = make(map[string]func(ctx *gin.Context, params interface{}) (interface{}, *JSONError))
	//notifyChans = userNotify{channels: make(map[string]chan interface{}), rwLock: sync.RWMutex{}}
)

/*
type userNotify struct {
	channels map[string]chan interface{}
	rwLock   sync.RWMutex
}

func (n *userNotify) Notify(username string, msg interface{}) {
	n.rwLock.RLock()
	defer n.rwLock.RUnlock()

	ch, ok := n.channels[username]
	if !ok {
		fmt.Println("ignore to notify", username, msg)
		return
	}

	ch <- msg
}

func (n *userNotify) getChannel(username string) chan interface{} {
	lock, channels := n.rwLock, n.channels
	lock.Lock()
	defer lock.Unlock()

	lastChan, ok := channels[username]
	if ok {
		fmt.Println("last chan", lastChan)
		close(lastChan)
	}

	ch := make(chan interface{}, 10)
	channels[username] = ch
	return ch
}

func (n *userNotify) close(username string) {
	lock, channels := n.rwLock, n.channels
	lock.Lock()
	defer lock.Unlock()

	ch, ok := channels[username]
	if !ok {
		return
	}

	delete(channels, username)
	close(ch)
}
*/
func Init(r *gin.RouterGroup) {
	// 标准json2rpc处理
	r.POST("v1", func(ctx *gin.Context) {
		j, err := parseJSON2Request(ctx)
		if err != nil {
			handleV1Error(ctx, nil, primitives.NewInvalidRequestError())
			return
		}

		jsonResp, jsonError := handleV1Request(ctx, j)
		if jsonError != nil {
			handleV1Error(ctx, j, jsonError)
			return
		}

		ctx.JSON(200, jsonResp)
	})

	//r.GET("v1/notify", handleNotify)

	//go func() { // TODO
	//	for {
	//		time.Sleep(time.Second * 5)
	//		notifyChans.Notify("kitty", time.Now())
	//	}
	//}()
}

func checkApiKey(ctx *gin.Context) bool {
	k := ctx.GetHeader("X-Api-Key")
	if len(k) == 0 {
		return false
	}

	if apikey.CheckApiKey(k) {
		ctx.Set("apiKey", k)
		return true
	}

	return false
}

func handleV1Request(ctx *gin.Context, j *JSON2Request) (*JSON2Response, *JSONError) {
	var resp interface{}
	var jsonError *JSONError

	if !checkApiKey(ctx) {
		return nil, primitives.NewCustomInternalError("invalid apikey")
	}

	//　查找、调用　处理函数
	if f, ok := handlers[j.Method]; ok {
		resp, jsonError = f(ctx, j.Params)
	} else {
		jsonError = primitives.NewMethodNotFoundError()
	}

	if jsonError != nil {
		return nil, jsonError
	}

	jsonResp := primitives.NewJSON2Response()
	jsonResp.ID = j.ID
	jsonResp.Result = resp

	return jsonResp, nil
}

/*
func handleNotify(ctx *gin.Context) {
	//if !checkApiKey(ctx) {
	//	ctx.AbortWithError(400, errors.New("invalid username/apikey"))
	//	return
	//}
	_ = errors.New("1")
	var upgrader = websocket.Upgrader{} // use default options

	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	username := ctx.GetHeader("X-Username")

	ch := notifyChans.getChannel(username)

	for {
		msg, ok := <-ch
		if msg != nil {
			fmt.Printf("recv: %+v\n", msg)
			err = c.WriteJSON(msg)
			if err != nil {
				fmt.Println("close chan . write error", err, ch)
				notifyChans.close(username)
				break
			}
		}

		if !ok {
			fmt.Println("already close ", username, ok, ch)
			break
		}
	}
}
*/
