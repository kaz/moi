package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var (
	wg    = sync.WaitGroup{}
	idCh  = make(chan string)
	queCh = make(chan []string)
	ansCh = make(chan string)
)

func deleteGame() {
	fmt.Println("Deleting game ...")

	c, err := tls.Dial("tcp", "apiv2.twitcasting.tv:443", &tls.Config{
		CipherSuites: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
	})
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriterSize(c, 1024)
	w.WriteString("DELETE /internships/2019/games HTTP/1.1\r\n")
	w.WriteString("Host: apiv2.twitcasting.tv\r\n")
	w.WriteString("Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIn0.eyJhdWQiOiIxODIyMjQ5MzguMjNhNzJmNDA2NzI4M2I0OWY5NjZmOTMyMzViMTg2NDQzN2VjNWY2YTlmY2M5NjVlOGIzOTM5MGRmNWQ2YWE5NCIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIiwiaWF0IjoxNTYyMzMyNzcyLCJuYmYiOjE1NjIzMzI3NzIsImV4cCI6MTU3Nzg4NDc3MSwic3ViIjoiYzphbGljZV9nIiwic2NvcGVzIjpbInJlYWQiXX0.GoWeDo-ZswQ1sp_ejnj9rKT2MwPaZwYpqC_w9GS5r1_bJ-aiPPd8rwUPUY3VFphkSVkVpZRwyGq3bxc1Rx1CQMh5_xBavoaMtCr1iik4YZPTIJJJBjXfpgQTOqRGgKA5HEvz84d4XnQGeFNBCw7zbpfOiBmiByHDfulh_SjkI-AwCUPJ4vaXHVjcHHtqlQfZ5jxYwLH2Zv0duwlsDMxR-tWU70TGeFV71yByE53fL4s6Heg607BeFDFhIvNmkMOULiv5xOrnJlmGrxySfflZn4KbStRydypvfgAc2Kbkz1YUQatQazN4hvCrr-otpyPccdKBLF8cf4UVpwSpCSClzQ\r\n")
	w.WriteString("\r\n")

	if err := w.Flush(); err != nil {
		panic(err)
	}
	if err := c.CloseWrite(); err != nil {
		panic(err)
	}
}

func gameStarter() {
	fmt.Println("Starting game ...")

	c, err := tls.Dial("tcp", "apiv2.twitcasting.tv:443", &tls.Config{
		CipherSuites: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
	})
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriterSize(c, 16*1024)
	r := bufio.NewReaderSize(c, 16*1024)

	w.WriteString("GET /internships/2019/games?level=3 HTTP/1.1\r\n")
	w.WriteString("Host: apiv2.twitcasting.tv\r\n")
	w.WriteString("Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIn0.eyJhdWQiOiIxODIyMjQ5MzguMjNhNzJmNDA2NzI4M2I0OWY5NjZmOTMyMzViMTg2NDQzN2VjNWY2YTlmY2M5NjVlOGIzOTM5MGRmNWQ2YWE5NCIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIiwiaWF0IjoxNTYyMzMyNzcyLCJuYmYiOjE1NjIzMzI3NzIsImV4cCI6MTU3Nzg4NDc3MSwic3ViIjoiYzphbGljZV9nIiwic2NvcGVzIjpbInJlYWQiXX0.GoWeDo-ZswQ1sp_ejnj9rKT2MwPaZwYpqC_w9GS5r1_bJ-aiPPd8rwUPUY3VFphkSVkVpZRwyGq3bxc1Rx1CQMh5_xBavoaMtCr1iik4YZPTIJJJBjXfpgQTOqRGgKA5HEvz84d4XnQGeFNBCw7zbpfOiBmiByHDfulh_SjkI-AwCUPJ4vaXHVjcHHtqlQfZ5jxYwLH2Zv0duwlsDMxR-tWU70TGeFV71yByE53fL4s6Heg607BeFDFhIvNmkMOULiv5xOrnJlmGrxySfflZn4KbStRydypvfgAc2Kbkz1YUQatQazN4hvCrr-otpyPccdKBLF8cf4UVpwSpCSClzQ\r\n")
	w.WriteString("\r\n")

	if err := w.Flush(); err != nil {
		panic(err)
	}
	if err := c.CloseWrite(); err != nil {
		panic(err)
	}

	if _, err = r.ReadBytes('{'); err != nil {
		panic(err)
	}
	if _, err = r.Discard(6); err != nil {
		panic(err)
	}
	id, err := r.ReadBytes('"')
	if err != nil {
		panic(err)
	}

	idCh <- string(id[:len(id)-1])

	if _, err = r.Discard(13); err != nil {
		panic(err)
	}
	que1, err := r.ReadBytes('=')
	if err != nil {
		panic(err)
	}
	que2, err := r.ReadBytes('"')
	if err != nil {
		panic(err)
	}

	queCh <- []string{string(que1[:len(que1)-2]), string(que2[1 : len(que2)-1])}

	wg.Done()
}

func answerQueryer() {
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	fmt.Println("Communicate with redis ...")
	fmt.Println(rdb.Ping().Result())

	que := <-queCh

	ans, err := rdb.HGet(que[0], que[1]).Result()
	if err != nil {
		fmt.Println(">>>", que)
		deleteGame()
		panic(err)
	}

	ansCh <- ans

	wg.Done()
}

func answerSender() {
	c, err := tls.Dial("tcp", "apiv2.twitcasting.tv:443", &tls.Config{
		CipherSuites: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
	})
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriterSize(c, 16*1024)

	fmt.Println("Answer sending standby ...")
	w.WriteString("POST /internships/2019/games/" + <-idCh + " HTTP/1.1\r\n")

	if err := w.Flush(); err != nil {
		panic(err)
	}

	w.WriteString("Host: apiv2.twitcasting.tv\r\n")
	w.WriteString("Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIn0.eyJhdWQiOiIxODIyMjQ5MzguMjNhNzJmNDA2NzI4M2I0OWY5NjZmOTMyMzViMTg2NDQzN2VjNWY2YTlmY2M5NjVlOGIzOTM5MGRmNWQ2YWE5NCIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIiwiaWF0IjoxNTYyMzMyNzcyLCJuYmYiOjE1NjIzMzI3NzIsImV4cCI6MTU3Nzg4NDc3MSwic3ViIjoiYzphbGljZV9nIiwic2NvcGVzIjpbInJlYWQiXX0.GoWeDo-ZswQ1sp_ejnj9rKT2MwPaZwYpqC_w9GS5r1_bJ-aiPPd8rwUPUY3VFphkSVkVpZRwyGq3bxc1Rx1CQMh5_xBavoaMtCr1iik4YZPTIJJJBjXfpgQTOqRGgKA5HEvz84d4XnQGeFNBCw7zbpfOiBmiByHDfulh_SjkI-AwCUPJ4vaXHVjcHHtqlQfZ5jxYwLH2Zv0duwlsDMxR-tWU70TGeFV71yByE53fL4s6Heg607BeFDFhIvNmkMOULiv5xOrnJlmGrxySfflZn4KbStRydypvfgAc2Kbkz1YUQatQazN4hvCrr-otpyPccdKBLF8cf4UVpwSpCSClzQ\r\n")

	if err := w.Flush(); err != nil {
		panic(err)
	}

	ans := "{\"answer\":\"" + <-ansCh + "\"}"
	w.WriteString("Content-Length: " + strconv.Itoa(len(ans)) + "\r\n\r\n")
	w.WriteString(ans)

	if err := w.Flush(); err != nil {
		panic(err)
	}
	if err := c.CloseWrite(); err != nil {
		panic(err)
	}

	resp, err := ioutil.ReadAll(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp))

	wg.Done()
}

func main() {
	wg.Add(3)

	go answerQueryer()

	time.Sleep(1 * time.Second)
	go answerSender()

	time.Sleep(1 * time.Second)
	go gameStarter()

	wg.Wait()
}
