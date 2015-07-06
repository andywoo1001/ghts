/* This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>.

@author: UnHa Kim <unha.kim.ghts@gmail.com> */

package shared

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"
)

// 이하 편의함수 모음

func TestF에러_체크(테스트 *testing.T) {
	// 그냥 nil을 넘겨주면 reflect.ValueOf(nil)이 무효한 값이 되어서 에러 발생함.
	// nil 대신에 error형식의 reflect.Zero값을 넘겨줘야 함.

	빈_에러값 := reflect.Zero(reflect.TypeOf(fmt.Errorf("")))
	F테스트_패닉없음(테스트, F에러_체크, 빈_에러값)

	//화면 출력을 안 보이게 하기
	원래_출력장치 := os.Stdout
	_, 출력장치, 에러 := os.Pipe()

	if 에러 != nil {
		테스트.FailNow()
	}

	os.Stdout = 출력장치

	F테스트_패닉발생(테스트, F에러_체크, fmt.Errorf("테스트용 에러"))

	출력장치.Close()
	os.Stdout = 원래_출력장치
}

func TestF문자열_복사(테스트 *testing.T) {
	테스트.Parallel()

	F테스트_같음(테스트, F문자열_복사("12 34 "), "12 34 ")
}

func TestF메시지_송신(테스트 *testing.T) {
	회신_채널 := make(chan bool)

	질의_메시지 := []interface{}{P메시지_구분_일반, "질의_메시지"}
	회신_에러 := fmt.Errorf("회신_에러")

	go f메시지_송신_테스트_REQ(회신_채널, 질의_메시지, 회신_에러)
	go f에러_메시지_송신_테스트_REP(회신_채널, 질의_메시지, 회신_에러)

	for i := 0; i < 2; i++ {
		테스트_결과 := <-회신_채널
		F테스트_참임(테스트, 테스트_결과)
	}
}

func f메시지_송신_테스트_REQ(회신_채널 chan bool, 질의_메시지 []interface{}, 회신_에러 error) {

	var 에러 error = nil

	defer func() {
		if 에러 != nil {
			회신_채널 <- false
		}
	}()

	소켓_REQ, 에러 := zmq.NewSocket(zmq.REQ)
	if 에러 != nil {
		return
	}

	defer 소켓_REQ.Close()

	에러 = 소켓_REQ.Connect(P테스트_결과_주소)
	if 에러 != nil {
		return
	}

	에러 = F메시지_송신(소켓_REQ, 질의_메시지...)
	if 에러 != nil {
		return
	}

	메시지, 에러 := 소켓_REQ.RecvMessage(0)
	if 에러 != nil {
		return
	}

	if len(메시지) != 2 ||
		메시지[0] != P메시지_구분_에러 ||
		메시지[1] != 회신_에러.Error() {
		회신_채널 <- false
		return
	}

	회신_채널 <- true
}

func f에러_메시지_송신_테스트_REP(회신_채널 chan bool, 질의_메시지 []interface{}, 회신_에러 error) {
	var 에러 error = nil

	defer func() {
		if 에러 != nil {
			회신_채널 <- false
		}
	}()

	소켓_REP, 에러 := zmq.NewSocket(zmq.REP)
	if 에러 != nil {
		return
	}

	defer 소켓_REP.Close()

	에러 = 소켓_REP.Bind(P테스트_결과_주소)
	if 에러 != nil {
		return
	}

	메시지, 에러 := 소켓_REP.RecvMessage(0)
	if 에러 != nil {
		return
	}

	if len(메시지) != len(질의_메시지) {
		회신_채널 <- false
		return
	}

	for i := 0; i < len(메시지); i++ {
		if 메시지[i] != 질의_메시지[i] {
			회신_채널 <- false
			return
		}
	}

	에러 = F에러_메시지_송신(소켓_REP, 회신_에러)
	if 에러 != nil {
		return
	}

	회신_채널 <- true
}

// 이하 최대 스레드 수량 관련 함수

func TestF단일_스레드_모드(테스트 *testing.T) {
	최대_스레드_수량_원본 := runtime.GOMAXPROCS(-1)
	defer func() {
		runtime.GOMAXPROCS(최대_스레드_수량_원본)
	}()

	runtime.GOMAXPROCS(2)
	F단일_스레드_모드()

	F테스트_같음(테스트, runtime.GOMAXPROCS(-1), 1)
}

func TestF멀티_스레드_모드(테스트 *testing.T) {
	최대_스레드_수량_원본 := runtime.GOMAXPROCS(-1)
	defer func() {
		runtime.GOMAXPROCS(최대_스레드_수량_원본)
	}()

	runtime.GOMAXPROCS(1)
	F멀티_스레드_모드()

	F테스트_같음(테스트, runtime.GOMAXPROCS(-1), runtime.NumCPU())
}

func TestF단일_스레드_모드임(테스트 *testing.T) {
	최대_스레드_수량_원본 := runtime.GOMAXPROCS(-1)
	defer func() {
		runtime.GOMAXPROCS(최대_스레드_수량_원본)
	}()

	F단일_스레드_모드()
	F테스트_참임(테스트, F단일_스레드_모드임())

	F멀티_스레드_모드()
	F테스트_거짓임(테스트, F단일_스레드_모드임())
}

func TestF멀티_스레드_모드임(테스트 *testing.T) {
	최대_스레드_수량_원본 := runtime.GOMAXPROCS(-1)
	defer func() {
		runtime.GOMAXPROCS(최대_스레드_수량_원본)
	}()

	F단일_스레드_모드()
	F테스트_거짓임(테스트, F멀티_스레드_모드임())

	F멀티_스레드_모드()
	F테스트_참임(테스트, F멀티_스레드_모드임())
}

// Go루틴 정리 관련 기능 테스트

func TestF등록된_Go루틴_종료(테스트 *testing.T) {
	// 공통 종료 채널을 이용하는 다른 테스트에 영향을 주지 않기 위해서
	// 새로운 채널을 이용해서 테스트를 진행함.
	ch공통_종료_채널_원본 := ch공통_종료_채널
	ch공통_종료_채널 = make(chan S비어있는_구조체)

	// 테스트 종료할 때, 공통 종료 채널을 원래대로 되돌려 놓음.
	defer func() {
		ch공통_종료_채널 = ch공통_종료_채널_원본
	}()

	ch입력_모음 := make([](chan int), 10)
	ch출력 := make(chan int)
	ch공통_종료 := F공통_종료_채널()

	// Go루틴 10개 생성
	for i, _ := range ch입력_모음 {
		ch입력 := make(chan int)
		ch입력_모음[i] = ch입력

		go f등록된_Go루틴_종료_테스트_도우미(ch입력, ch출력, ch공통_종료)
	}

	// 모든 Go루틴 존재 확인
	for _, ch입력 := range ch입력_모음 {
		ch입력 <- 1
		F테스트_같음(테스트, <-ch출력, 1)
	}

	F등록된_Go루틴_종료()

	for range ch입력_모음 {
		F테스트_같음(테스트, <-ch출력, 999)
	}
}

func f등록된_Go루틴_종료_테스트_도우미(ch입력, ch출력 chan int, ch공통_종료 chan S비어있는_구조체) {
	for {
		select {
		case _, ok := <-ch입력:
			if !ok {
				ch출력 <- 10
			}

			ch출력 <- 1
		case <-ch공통_종료:
			ch출력 <- 999
			return
		}
	}
}
