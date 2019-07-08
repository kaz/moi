package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	// #cgo LDFLAGS: -lssl
	// #include "moi.c"
	"C"
)
import "bytes"

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

	resp := []byte(C.GoString(C.start_game()))
	idCh <- string(resp[7:39])
	queCh <- []string{string(resp[53:78]), string(resp[81 : 81+bytes.Index(resp[81:], []byte("\""))])}

	wg.Done()
}

func answerQueryer() {
	buf := make([]byte, 12)

	redis, err := net.Dial("unix", "/var/run/redis/redis.sock")
	if err != nil {
		panic(err)
	}
	if _, err := redis.Write([]byte("*3\r\n$4\r\nHGET\r\n$25\r\n")); err != nil {
		panic(err)
	}

	que := <-queCh
	data := fmt.Sprintf("%s\r\n$%d\r\n%s\r\n", que[0], len(que[1]), que[1])
	if _, err := redis.Write([]byte(data)); err != nil {
		panic(err)
	}

	if n, err := redis.Read(buf); err != nil {
		panic(err)
	} else if n != 12 {
		deleteGame()
		panic(fmt.Errorf("error:%v, query:%v", err, que))
	}

	ansCh <- string(buf[4:10])

	wg.Done()
}

func answerSender() {
	w := bytes.NewBuffer(make([]byte, 0, 4096))
	C.prepare()

	fmt.Println("Answer sending standby ...")

	w.WriteString("POST /internships/2019/games/" + <-idCh + " HTTP/1.1\r\n")
	w.WriteString("Host: apiv2.twitcasting.tv\r\n")
	w.WriteString("Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIn0.eyJhdWQiOiIxODIyMjQ5MzguMjNhNzJmNDA2NzI4M2I0OWY5NjZmOTMyMzViMTg2NDQzN2VjNWY2YTlmY2M5NjVlOGIzOTM5MGRmNWQ2YWE5NCIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIiwiaWF0IjoxNTYyMzMyNzcyLCJuYmYiOjE1NjIzMzI3NzIsImV4cCI6MTU3Nzg4NDc3MSwic3ViIjoiYzphbGljZV9nIiwic2NvcGVzIjpbInJlYWQiXX0.GoWeDo-ZswQ1sp_ejnj9rKT2MwPaZwYpqC_w9GS5r1_bJ-aiPPd8rwUPUY3VFphkSVkVpZRwyGq3bxc1Rx1CQMh5_xBavoaMtCr1iik4YZPTIJJJBjXfpgQTOqRGgKA5HEvz84d4XnQGeFNBCw7zbpfOiBmiByHDfulh_SjkI-AwCUPJ4vaXHVjcHHtqlQfZ5jxYwLH2Zv0duwlsDMxR-tWU70TGeFV71yByE53fL4s6Heg607BeFDFhIvNmkMOULiv5xOrnJlmGrxySfflZn4KbStRydypvfgAc2Kbkz1YUQatQazN4hvCrr-otpyPccdKBLF8cf4UVpwSpCSClzQ\r\n")

	ans := "{\"answer\":\"" + <-ansCh + "\"}"
	w.WriteString("Content-Length: " + strconv.Itoa(len(ans)) + "\r\n\r\n")
	w.WriteString(ans)

	resp := C.GoString(C.answer(C.CString(string(w.Bytes()))))
	fmt.Println(resp)

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
