/* This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>.

@author: UnHa Kim <unha.kim.ghts@gmail.com> */

package common

import (
	"runtime"
	"testing"
	"time"
)

func TestF문자열_검색_복수_정규식(테스트 *testing.T) {
	검색_대상 := "aabbcc <span>xxx2006.01.02xxx</span> ddeeff"
	정규식_문자열_모음 := []string{
		`<span>.*[0-9]{4}.[0-9]{1,2}.[0-9]{1,2}.*</span>`,
		`[0-9]{4}.[0-9]{1,2}.[0-9]{1,2}`}

	검색_결과 := F문자열_검색_복수_정규식(검색_대상, 정규식_문자열_모음)

	F테스트_같음(테스트, 검색_결과, "2006.01.02")
}

func TestF절대값(테스트 *testing.T) {
	F테스트_같음(테스트, F절대값(-1), 1.0)
	F테스트_같음(테스트, F절대값(1), 1.0)
	F테스트_같음(테스트, F절대값(int64(-1)), 1.0)
	F테스트_같음(테스트, F절대값(int64(1)), 1.0)
	F테스트_같음(테스트, F절대값(float32(-1.0)), 1.0)
	F테스트_같음(테스트, F절대값(float32(1.0)), 1.0)
	F테스트_같음(테스트, F절대값(float64(-1.0)), 1.0)
	F테스트_같음(테스트, F절대값(float64(1.0)), 1.0)
}

func TestF2포맷된_시각(테스트 *testing.T) {
	시각, 에러 := F2포맷된_시각("2006-01-02 15:04:05", "2001-02-03 04:05:06")
	F테스트_에러없음(테스트, 에러)
	F테스트_같음(테스트, 시각.Year(), 2001)
	F테스트_같음(테스트, 시각.Month(), time.February)
	F테스트_같음(테스트, 시각.Day(), 3)
	F테스트_같음(테스트, 시각.Hour(), 4)
	F테스트_같음(테스트, 시각.Minute(), 5)
	F테스트_같음(테스트, 시각.Second(), 6)

	시각, 에러 = F2포맷된_시각("2006-01-02 15:04:05", "에러 발생 유발 문자열")
	F테스트_에러발생(테스트, 에러)
	F테스트_참임(테스트, 시각.IsZero(), 시각)
}

func TestF문자열_모음2인터페이스_모음(테스트 *testing.T) {
	문자열_모음 := []string{"테스트1", "2", "3.0"}
	인터페이스_모음 := F2인터페이스_모음(문자열_모음)

	F테스트_같음(테스트, len(인터페이스_모음), len(문자열_모음))

	for i := 0; i < len(문자열_모음); i++ {
		F테스트_같음(테스트, 인터페이스_모음[i].(string), 문자열_모음[i])
	}
}

func TestF인터페이스_모음2문자열_모음(테스트 *testing.T) {
	인터페이스_모음 := []interface{}{"테스트", 1, time.Now()}
	문자열_모음 := F2문자열_모음(인터페이스_모음)

	F테스트_같음(테스트, len(인터페이스_모음), len(문자열_모음))

	for i := 0; i < len(문자열_모음); i++ {
		F테스트_같음(테스트, F2문자열(인터페이스_모음[i]), 문자열_모음[i])
	}
}

func TestF문자열_복사(테스트 *testing.T) {
	테스트.Parallel()

	F테스트_같음(테스트, F문자열_복사("12 34 "), "12 34 ")
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
