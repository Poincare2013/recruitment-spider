package crawl

import (
	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/streadway/amqp"
	"go-recruitment-spider/config"
	"go-recruitment-spider/models"
	"go-recruitment-spider/modules/mq"
	"go-recruitment-spider/modules/mymgo"
	"go-recruitment-spider/utils"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func FetchCookie() (res string) {
	f, err := ioutil.ReadFile("./config/liepin.json")
	if err != nil {
		log.Fatal("read fail", err)
	}
	var cookies []Cookie
	err = json.Unmarshal(f, &cookies)
	if err != nil {
		log.Fatal("unmarshal fail", err)
	}
	for _, v := range cookies {
		res += v.Name + "=" + v.Value + "; "
	}
	if len(res) >= 3 {
		res = res[0 : len(res)-2]
	}
	return
}

type liepinCrawl struct {
	mdbSession *mymgo.MdbSession
	//同时爬取线程数
	threads   int
	toCrawlCh <-chan amqp.Delivery
	//推回mq重爬
	reCrawlCh chan string
	//记录重试url
	retryMap *utils.ConMap
	//插入mongo成功数
	successSum uint64
	//失败数
	failSum    uint64
	collector  *colly.Collector
	errCh      chan *models.ErrItem
	wg         *sync.WaitGroup
	queueName  string
	interrupt  chan os.Signal
	interrupt2 chan os.Signal
	isOver     bool
}

const liepinQueue string = "liepin_list_page"
const liepinCollection string = "liepinCrawl"
const maxRetryTimes int = 5
const selectorResIdEncode = "body > div.wrap.relative > aside.board > section.title-info > div:nth-child(1) > h6.float-left > small:nth-child(1)"
const errResIdEncodeCollection = "errResIdEncode"

func NewLiepinCrawl(interrupt chan os.Signal) *liepinCrawl {
	m := mq.NewMyReceiver(liepinQueue)
	msgs, err := m.GetDelivery()
	if err != nil {
		log.Fatal(err)
	}
	mdbSession := mymgo.GetDb()
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
	return &liepinCrawl{
		mdbSession: mdbSession,
		threads:    config.GetTomlConfig().Concurrency.Num,
		toCrawlCh:  msgs,
		reCrawlCh:  make(chan string, 1000),
		retryMap:   &utils.ConMap{Data: make(map[string]int), Mutex: &sync.Mutex{}},
		collector:  c,
		errCh:      make(chan *models.ErrItem, 100),
		wg:         &sync.WaitGroup{},
		queueName:  liepinQueue,
		interrupt:  interrupt,
		interrupt2: make(chan os.Signal, 100), //用来结束第二层协程
		isOver:     false,
	}
}

func (t *liepinCrawl) Crawl(ctx context.Context) {
	cookie := FetchCookie()
	t.handle()
	//一个协程把失败url推回队列
	go t.reCrawl(ctx)
	ctx2, cancel2 := context.WithCancel(ctx)
	go func() {
		<-t.interrupt2
		cancel2()
	}()
	for i := 0; i < t.threads; i++ {
		t.wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for {
				select {
				case <-ctx2.Done():
					wg.Done()
					log.Println("退出爬虫并发协程.....")
					return
				case toCrawl := <-t.toCrawlCh:
					//url要转成wap去爬
					resIdEncode := string(toCrawl.Body)
					toCrawl.Ack(false)
					u, _ := url.Parse(buildUrl(resIdEncode))
					collyCtx := colly.NewContext()
					collyCtx.Put("resIdEncode", resIdEncode)
					h := &http.Header{}
					h.Set("cookie", cookie)
					req := &colly.Request{
						URL:     u,
						Method:  "GET",
						Ctx:     collyCtx,
						Headers: h,
					}
					rb, err := req.Marshal()
					if err != nil {
						log.Println("request序列化错误 ", buildUrl(resIdEncode))
						continue
					}
					r, err := t.collector.UnmarshalRequest(rb)
					if err != nil {
						log.Println("creq失败", buildUrl(resIdEncode))
						continue
					}
					r.Do()
				}
			}
		}(t.wg)
	}
	t.wg.Wait()
	return
}

func (t *liepinCrawl) handle() {
	t.collector.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
		resIdEncode := r.Ctx.Get("resIdEncode")
		if err != nil {
			log.Println(err)
			t.errCh <- &models.ErrItem{
				Id:  resIdEncode,
				Err: err.Error(),
			}
			return
		}
		if doc.Find(selectorResIdEncode) != nil && doc.Find(selectorResIdEncode).First() != nil &&
			len(doc.Find(selectorResIdEncode).First().Nodes) >= 1 &&
			doc.Find(selectorResIdEncode).First().Nodes[0] != nil &&
			doc.Find(selectorResIdEncode).First().Nodes[0].FirstChild != nil {
			data := doc.Find(selectorResIdEncode).First().Nodes[0].FirstChild.Data
			//如果没有简历编号,说明没被发爬
			if !strings.Contains(data, "简历编号") {
				t.reCrawlCh <- resIdEncode
				log.Println("账号已被封禁,程序终止")
				t.finish()
				return
			}
			err = t.mdbSession.UpsertId(liepinCollection, resIdEncode, bson.M{"$set": bson.M{"responseBody": string(r.Body)}})
			if err == nil {
				//成功了就把map中的key删掉
				t.retryMap.Delete(resIdEncode)
				atomic.AddUint64(&t.successSum, 1)
				log.Println("简历:", resIdEncode, " 爬取成功...")
			}
		} else {
			t.reCrawlCh <- resIdEncode
			log.Println("账号已被封禁,程序终止")
			t.finish2()
			return
		}
	})

	t.collector.OnError(func(r *colly.Response, e error) {
		resIdEncode := r.Ctx.Get("resIdEncode")
		//有些删掉的url也要存到mongo
		log.Println(e.Error(), " resIdEncode:", resIdEncode)
		//这种错误往错误队列推了一次
		t.errCh <- &models.ErrItem{
			Id:  resIdEncode,
			Err: e.Error(),
		}
	})
}

func (t *liepinCrawl) reCrawl(ctx context.Context) {
	go utils.CronTask(func() {
		if t.isOver && len(t.reCrawlCh) == 0 {
			t.finish()
		}
	}, time.Second)
	s := mq.NewSession(t.queueName)
	t.wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			t.wg.Done()
			log.Println("退出reCrawl协程")
			return
		case msg := <-t.reCrawlCh:
			err := s.Ch.Publish(
				"",           // exchange
				s.Queue.Name, // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
				})
			//log.Printf(" [x] Sent %s", msg)
			if err != nil {
				log.Println("Failed to publish message: ", msg)
			}
		}
	}
}

func (t *liepinCrawl) finish() {
	t.interrupt <- syscall.SIGINT
}

func (t *liepinCrawl) finish2() {
	t.isOver = true
	t.interrupt2 <- syscall.SIGINT
}

//简历id唯一决定url
func buildUrl(resIdEncode string) string {
	return "https://lpt.liepin.com/cvview/showresumedetail?resIdEncode=" + resIdEncode
}
