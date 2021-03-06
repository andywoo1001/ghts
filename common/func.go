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
	"C"

	"github.com/suapapa/go_hangul/encoding/cp949"

	"bytes"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func F_HTTP회신_본문(url string) (string, error) {
	응답, 에러 := http.Get(url)
	defer func() {
		if 응답 != nil && 응답.Body != nil {
			응답.Body.Close()
		}
	}()

	if 에러 != nil || 응답.Body == nil {
		return "", 에러
	}

	바이트_모음, 에러 := ioutil.ReadAll(응답.Body)

	if 에러 != nil || 바이트_모음 == nil {
		return "", 에러
	}

	return string(바이트_모음), nil
}

func F_HTTP회신_본문_CP949(url string) (string, error) {
	응답, 에러 := http.Get(url)
	defer func() {
		if 응답 != nil && 응답.Body != nil {
			응답.Body.Close()
		}
	}()

	if 에러 != nil || 응답.Body == nil {
		return "", 에러
	}

	바이트_모음, 에러 := ioutil.ReadAll(응답.Body)

	if 에러 != nil || 바이트_모음 == nil {
		return "", 에러
	}

	return F2문자열_CP949(바이트_모음), nil
}

func F문자열_검색_복수_정규식(검색_대상 string, 정규식_문자열_모음 []string) string {
	검색_결과 := ""

	for _, 정규식_문자열 := range 정규식_문자열_모음 {
		정규식 := regexp.MustCompile(정규식_문자열)
		검색_결과 = 정규식.FindString(검색_대상)

		if 검색_결과 == "" {
			break
		}

		검색_대상 = 검색_결과
	}

	return 검색_결과
}

func F절대값(값 interface{}) float64 {
	실수값 := float64(0.0)
	switch 값.(type) {
	case int:
		실수값 = float64(값.(int))
	case int64:
		실수값 = float64(값.(int64))
	case float32:
		실수값 = float64(값.(float32))
	case float64:
		실수값 = 값.(float64)
	default:
		에러 := F에러("예상치 못한 자료형. %v %v", reflect.TypeOf(값), 값)
		panic(에러)
	}

	return math.Abs(실수값)
}

func F에러_패닉(에러 error) {
	if 에러 != nil {
		F에러_출력(에러)
		panic("")
	}
}

func F2가변형(값 interface{}) interface{} { return 값 }

func F2바이트_모음(값 interface{}) []byte {
	switch 값.(type) {
	case [1]byte:
		배열 := 값.([1]byte)
		return 배열[:]
	case [2]byte:
		배열 := 값.([2]byte)
		return 배열[:]
	case [3]byte:
		배열 := 값.([3]byte)
		return 배열[:]
	case [4]byte:
		배열 := 값.([4]byte)
		return 배열[:]
	case [5]byte:
		배열 := 값.([5]byte)
		return 배열[:]
	case [6]byte:
		배열 := 값.([6]byte)
		return 배열[:]
	case [7]byte:
		배열 := 값.([7]byte)
		return 배열[:]
	case [8]byte:
		배열 := 값.([8]byte)
		return 배열[:]
	case [9]byte:
		배열 := 값.([9]byte)
		return 배열[:]
	case [10]byte:
		배열 := 값.([10]byte)
		return 배열[:]
	case [11]byte:
		배열 := 값.([11]byte)
		return 배열[:]
	case [12]byte:
		배열 := 값.([12]byte)
		return 배열[:]
	case [13]byte:
		배열 := 값.([13]byte)
		return 배열[:]
	case [14]byte:
		배열 := 값.([14]byte)
		return 배열[:]
	case [15]byte:
		배열 := 값.([15]byte)
		return 배열[:]
	case [16]byte:
		배열 := 값.([16]byte)
		return 배열[:]
	case [17]byte:
		배열 := 값.([17]byte)
		return 배열[:]
	case [18]byte:
		배열 := 값.([18]byte)
		return 배열[:]
	case [19]byte:
		배열 := 값.([19]byte)
		return 배열[:]
	case [20]byte:
		배열 := 값.([20]byte)
		return 배열[:]
	case [21]byte:
		배열 := 값.([21]byte)
		return 배열[:]
	case [22]byte:
		배열 := 값.([22]byte)
		return 배열[:]
	case [23]byte:
		배열 := 값.([23]byte)
		return 배열[:]
	case [24]byte:
		배열 := 값.([24]byte)
		return 배열[:]
	case [25]byte:
		배열 := 값.([25]byte)
		return 배열[:]
	case [26]byte:
		배열 := 값.([26]byte)
		return 배열[:]
	case [27]byte:
		배열 := 값.([27]byte)
		return 배열[:]
	case [28]byte:
		배열 := 값.([28]byte)
		return 배열[:]
	case [29]byte:
		배열 := 값.([29]byte)
		return 배열[:]
	case [30]byte:
		배열 := 값.([30]byte)
		return 배열[:]
	case [80]byte:
		배열 := 값.([80]byte)
		return 배열[:]
	case [100]byte:
		배열 := 값.([100]byte)
		return 배열[:]
	default:
		F변수값_확인(값)
		에러 := F에러("예상치 못한 자료형")
		panic(에러)
	}
}

func F2문자열_CP949(값 interface{}) string {
	바이트_모음_CP949 := make([]byte, 0)

	switch 값.(type) {
	case []byte:
		바이트_모음_CP949 = 값.([]byte)
	case [1]byte, [2]byte, [3]byte, [4]byte, [5]byte,
		[6]byte, [7]byte, [8]byte, [9]byte, [10]byte,
		[11]byte, [12]byte, [13]byte, [14]byte, [15]byte,
		[16]byte, [17]byte, [18]byte, [19]byte, [20]byte,
		[21]byte, [22]byte, [23]byte, [24]byte, [25]byte,
		[26]byte, [27]byte, [28]byte, [29]byte, [30]byte,
		[80]byte, [100]byte:
		바이트_모음_CP949 = F2바이트_모음(값)
	default:
		에러 := F에러("예상치 못한 자료 형식. %v", reflect.TypeOf(값))
		panic(에러)
	}

	null문자_인덱스 := strings.Index(string(바이트_모음_CP949), "\x00")

	if null문자_인덱스 >= 0 {
		바이트_모음_CP949 = 바이트_모음_CP949[:null문자_인덱스]
	}

	바이트_모음_utf8, 에러 := cp949.From(바이트_모음_CP949)
	F에러_패닉(에러)

	return string(바이트_모음_utf8)
}

func F2문자열(값 interface{}) string {
	switch 값.(type) {
	case string:
		return 값.(string)
	case uint, uint8, uint16, uint32, uint64,
		int, int8, int16, int32, int64, bool:
		// 많이 쓰는 형식들은 이 단계에서 바로 처리해서 속도 향상 도모.
		return F포맷된_문자열("%v", 값)
	case float32:
		return strconv.FormatFloat(float64(값.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(값.(float64), 'f', -1, 64)
	case time.Time:
		return 값.(time.Time).Format(P시간_형식)
	case []byte:
		return string(값.([]byte))
	case [1]byte, [2]byte, [3]byte, [4]byte, [5]byte,
		[6]byte, [7]byte, [8]byte, [9]byte, [10]byte,
		[11]byte, [12]byte, [13]byte, [14]byte, [15]byte,
		[16]byte, [17]byte, [18]byte, [19]byte, [20]byte,
		[21]byte, [22]byte, [23]byte, [24]byte, [25]byte,
		[26]byte, [27]byte, [28]byte, [29]byte, [30]byte,
		[80]byte, [100]byte:
		바이트_모음 := F2바이트_모음(값)

		바이트_모음 = bytes.TrimPrefix(바이트_모음, []byte("\x00"))
		null문자_인덱스 := strings.Index(string(바이트_모음), "\x00")

		if null문자_인덱스 >= 0 {
			바이트_모음 = 바이트_모음[:null문자_인덱스]
		}

		return string(바이트_모음)
	default:
		if 값 != nil {
			자료형 := reflect.TypeOf(값)

			if 자료형.Kind() == reflect.Array &&
				strings.HasSuffix(자료형.String(), "_Ctype_char") {
				에러 := F에러("C.char 배열")
				panic(에러)
			}
		}

		return F포맷된_문자열("%v", 값)
	}
}

func F2문자열_모음(인터페이스_모음 []interface{}) []string {
	if 인터페이스_모음 == nil {
		return nil
	}

	문자열_모음 := make([]string, len(인터페이스_모음))

	for i := 0; i < len(인터페이스_모음); i++ {
		문자열_모음[i] = F2문자열(인터페이스_모음[i])
	}

	return 문자열_모음
}

func F2정수(값 interface{}) int {
	문자열 := ""

	switch 값.(type) {
	case string:
		문자열 = 값.(string)
	default:
		문자열 = F2문자열(값)
	}

	문자열 = strings.TrimSpace(문자열)

	반환값, 에러 := strconv.Atoi(문자열)
	F에러_패닉(에러)

	return 반환값
}

func F2정수64(값 interface{}) (int64, error) {
	문자열 := ""

	switch 값.(type) {
	case string:
		문자열 = 값.(string)
	default:
		문자열 = F2문자열(값)
	}

	문자열 = strings.TrimSpace(문자열)

	반환값, 에러 := strconv.ParseInt(문자열, 10, 64)

	if 에러 != nil {
		반환값 = 0
	}

	return 반환값, 에러
}

func F2정수64_모음(값_모음 []interface{}) ([]int64, error) {
	정수64_모음 := make([]int64, 0)

	for _, 값 := range 값_모음 {
		정수64, 에러 := F2정수64(값)

		if 에러 != nil {
			F에러("정수 변환 에러 발생. '%v'\n%v", 값, 에러)
			return make([]int64, 0), 에러
		}

		정수64_모음 = append(정수64_모음, 정수64)
	}

	return 정수64_모음, nil
}

func F2실수(값 interface{}) float64 {
	문자열 := ""

	switch 값.(type) {
	case string:
		문자열 = 값.(string)
	default:
		문자열 = F2문자열(값)
	}

	문자열 = strings.TrimSpace(문자열)

	반환값, 에러 := strconv.ParseFloat(문자열, 64)
	F에러_패닉(에러)

	return 반환값
}

func F2시각(값 interface{}) time.Time {
	문자열 := ""

	switch 값.(type) {
	case string:
		문자열 = 값.(string)
	default:
		문자열 = F2문자열(값)
	}

	문자열 = strings.TrimSpace(문자열)

	반환값, 에러 := time.Parse(P시간_형식, 문자열)
	F에러_패닉(에러)

	return 반환값
}

func F2포맷된_시각(포맷 string, 값 interface{}) (time.Time, error) {
	문자열 := ""

	switch 값.(type) {
	case string:
		문자열 = 값.(string)
	default:
		문자열 = F2문자열(값)
	}

	문자열 = strings.TrimSpace(문자열)

	시각, 에러 := time.Parse(포맷, 문자열)

	if 에러 != nil {
		return time.Time{}, 에러
	}

	if strings.Contains(포맷, "MST") {
		시각 = 시각.Local() // 현지 시간으로 변환
	} else {
		// 포멧에 시간대가 없으면 UTC임. 현지 시간대로 바꿈.
		시각 = time.Date(시각.Year(), 시각.Month(), 시각.Day(),
			시각.Hour(), 시각.Minute(), 시각.Second(), 시각.Nanosecond(),
			time.Now().Location())
	}

	return 시각, 에러
}

func F2참거짓(값 interface{}, 조건 interface{}, 결과 bool) bool {
	if F2문자열(값) == F2문자열(조건) {
		return 결과
	} else {
		return !결과
	}
}

func F2인터페이스_모음(문자열_모음 []string) []interface{} {
	if 문자열_모음 == nil {
		return nil
	}

	인터페이스_모음 := make([]interface{}, len(문자열_모음))

	for i := 0; i < len(문자열_모음); i++ {
		인터페이스_모음[i] = 문자열_모음[i]
	}

	return 인터페이스_모음
}

/*
func F바이트_모음_늘리기(바이트_모음 []byte, 길이 int) []byte {
	if len(바이트_모음) > 길이 {
		에러 := F에러("지정된 길이가 더 짧음.")
		panic(에러)
	}

	반환값 := make([]byte, 길이)

	for i := 0; i < len(바이트_모음); i++ {
		반환값[i] = 바이트_모음[i]
	}

	return 반환값
}
*/

func F타입_이름(i interface{}) string {
	return reflect.TypeOf(i).Name()
}

func F문자열_복사(문자열 string) string {
	return (문자열 + " ")[:len(문자열)]
}

// 이하 최대 스레드 수량 관련 함수

func F단일_스레드_모드() { runtime.GOMAXPROCS(1) }
func F멀티_스레드_모드() { runtime.GOMAXPROCS(runtime.NumCPU()) }

func F단일_스레드_모드임() bool {
	if runtime.GOMAXPROCS(-1) == 1 {
		return true
	} else {
		return false
	}
}

func F멀티_스레드_모드임() bool { return !F단일_스레드_모드임() }

// 이하 종료 시 존재하는 모든 Go루틴 정리(혹은 종료) 관련 함수 모음
var ch공통_종료_채널 = make(chan S비어있는_구조체)

func F공통_종료_채널() chan S비어있는_구조체 {
	return ch공통_종료_채널
}

func F공통_종료_채널_재설정() {
	ch공통_종료_채널 = make(chan S비어있는_구조체)
}

func F등록된_Go루틴_종료() {
	close(ch공통_종료_채널)
}

func F_nil에러() error { return nil }
