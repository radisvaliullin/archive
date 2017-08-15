package genhandcli

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	lock "github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
)

// CliConf - client configs.
type CliConf struct {
	Addr       string
	GetErrMode bool
}

// Client - client. Work as message generator and handler.
type Client struct {
	conf *CliConf

	rcln *redis.Client

	genLockKey string
	genMsgKey  string
	errMsgKey  string
	// stop msg handler
	stopHndl chan struct{}
	isStoped chan struct{}
}

// NewClient - init new client.
func NewClient(conf *CliConf) *Client {
	cli := &Client{
		conf:       conf,
		genLockKey: "genLockKey",
		genMsgKey:  "genMsgKey",
		errMsgKey:  "errMsgKey",
		stopHndl:   make(chan struct{}),
		isStoped:   make(chan struct{}),
	}
	return cli
}

// Start - start client app.
func (c *Client) Start() error {

	// redis client setup
	opt := &redis.Options{
		Addr:     c.conf.Addr,
		Password: "",
		DB:       0,
	}
	c.rcln = redis.NewClient(opt)

	pong, err := c.rcln.Ping().Result()
	if err != nil {
		log.Print("redis ping err ", err)
		return err
	}
	log.Print("redis ping res ", pong)

	// get error mode
	if c.conf.GetErrMode {
		c.getAllErrMsg()
		return nil
	}

	// run
	go c.run()
	for {
		log.Println("Heart beat")
		time.Sleep(time.Second * 10)
	}
}

//
func (c *Client) run() {

	var lck *lock.Lock
	var err error
	var isMsgHandler bool

	// try lock
	for {
		// lock redis
		lck, err = lock.ObtainLock(c.rcln, "mylockkey", nil)
		if err != nil {
			log.Print("try lock err: ", err)
			return
		} else if lck == nil {
			log.Print("try lock: could not obtain lock")

			// start msg handler
			if !isMsgHandler {
				isMsgHandler = true
				c.startMsgHandler()
			}

			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	// unlock befor exit
	defer lck.Unlock()

	// if isMsgHandler, stop msg handler
	if isMsgHandler {
		c.stopHndl <- struct{}{}
		// wait while stop
		<-c.isStoped
	}

	// start generator mode and renew lock
	// start generator
	c.startGenerator()
	for {
		time.Sleep(time.Second * 2)

		// Renew your lock
		ok, err := lck.Lock()
		if err != nil {
			log.Print("lock renew err ", err)
			return
		} else if !ok {
			log.Print("lock renew: could not renew lock")
			return
		}
		log.Println("lock isrenewed")
	}
}

// startGenerator - start client as generator. Only one instance can be generator.
func (c *Client) startGenerator() {
	go c.runGenerator()
}

//
func (c *Client) runGenerator() {
	log.Print("run generator")
	cnt := 0
	for {
		cnt++

		// push next generate value
		_, err := c.rcln.LPush(c.genMsgKey, fmt.Sprintf("%v", cnt)).Result()
		if err != nil {
			log.Print("generator: push value err ", err)
		}
		log.Print("generator: pushed value - ", cnt)

		// Next message delay
		time.Sleep(time.Millisecond * 500)
	}
}

//
func (c *Client) startMsgHandler() {
	go c.runMsgHandler()
}

//
func (c *Client) runMsgHandler() {
	log.Print("run message handler")
	for {
		select {
		case <-c.stopHndl:
			c.isStoped <- struct{}{}
			return
		default:
			// get next message from redis, and handle
			popres, err := c.rcln.RPop(c.genMsgKey).Result()
			if err == redis.Nil {
				log.Print("msgHandler: pop value, empty list ", err)
				time.Sleep(time.Millisecond * 100)
			} else if err != nil {
				log.Print("msgHandler: pop value err ", err)
			} else {
				log.Print("msgHandler: pop value res ", popres)

				// if message is invalid push to err message list.
				if !c.isMsgValid(popres) {
					_, err := c.rcln.RPush(c.errMsgKey, popres).Result()
					if err != nil {
						log.Print("msgHandler: push err value err ", err)
					}
				}
			}
		}
	}
}

// isMsgValid - 5 percent of message is not valid.
func (c *Client) isMsgValid(msg string) bool {
	rsrc := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rsrc)
	rf := r.Float64()
	b := false
	if rf <= 0.95 {
		b = true
	}
	return b
}

// getAllErrMsg - read all error messages
func (c *Client) getAllErrMsg() {

	//
	msgres, err := c.rcln.LRange(c.errMsgKey, 0, -1).Result()
	if err != nil {
		log.Print("getAllErrMsg: lrange err ", err)
		return
	}
	log.Print("all error messages - ", msgres)
}
