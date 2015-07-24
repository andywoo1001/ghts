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

package shared_data

import (
	공용 "github.com/ghts/ghts/shared"

	"strconv"
)

var Ch주소 = make(chan 공용.I질의, 10)
var Ch종목 = make(chan 공용.I질의, 100)
var Ch종목별_보유량 = make(chan 공용.I질의, 100)
var Ch문자열_캐시 = make(chan 공용.I질의, 100)
var ch제어_공용_데이터_Go루틴 = make(chan 공용.I질의, 1)

var 공용_데이터_Go루틴_실행_중 = 공용.New안전한_bool(false)

func F공용_데이터_Go루틴_실행_중() bool {
	return 공용_데이터_Go루틴_실행_중.G값()
}

func F공용_데이터_Go루틴(ch초기화 chan bool) {
	초기화_에러 := 공용.F_nil에러()

	defer func() {
		if 초기화_에러 != nil {
			ch초기화 <- false
		}
	}()

	초기화_에러 = 공용_데이터_Go루틴_실행_중.S값(true)
	if 초기화_에러 != nil {
		return
	}

	주소_맵, 초기화_에러 := f주소_맵_초기화()
	if 초기화_에러 != nil {
		return
	}

	종목_맵, 초기화_에러 := f종목_맵_초기화()
	if 초기화_에러 != nil {
		return
	}
	
	보유량_맵, 초기화_에러 := f종목별_보유량_맵_초기화()
	if 초기화_에러 != nil {
		return
	}
	//make(map[string]공용.I종목별_보유량)

	문자열_캐시_맵 := make(map[string][]string)

	공통_종료_채널 := 공용.F공통_종료_채널()

	// 초기화 완료
	ch초기화 <- true

	// 변수를 재활용 하여 조금이라도 GC 부담을 덜자.
	질의 := 공용.New질의(공용.P메시지_GET) //비어있는 내용
	에러 := 공용.F_nil에러()

	for {
		select {
		case 질의 = <-Ch주소:
			if 에러 = 질의.G검사(공용.P메시지_GET, 1); 에러 != nil {
				공용.F에러_출력(에러)
				panic(에러)
			}

			주소, 존재함 := 주소_맵[질의.G내용(0)]

			if !존재함 {
				질의.S회신(공용.F에러_생성("잘못된 주소 질의값 '%v'", 질의.G내용(0)))
			}
			
			질의.S회신(nil, 주소)
		case 질의 = <-Ch종목:
			if 에러 = 질의.G검사(공용.P메시지_GET, 1); 에러 != nil {
				공용.F에러_출력(에러)
				panic(에러)
			}

			종목, 존재함 := 종목_맵[질의.G내용(0)]

			if !존재함 {
				질의.S회신(공용.F에러_생성("잘못된 종목 질의값 '%v'", 질의.G내용(0)))
				continue
			}
			
			질의.S회신(nil, 종목.G코드(), 종목.G이름())
		case 질의 = <-Ch문자열_캐시:
			// shared_data_string_cache.go 참조
			f문자열_캐시_질의_처리(문자열_캐시_맵, 질의)
		case 질의 = <-ch제어_공용_데이터_Go루틴:
			switch 질의.G구분() {
			case 공용.P메시지_초기화:
				주소_맵, 에러 = f주소_맵_초기화()
				if 에러 != nil {
					공용.F에러_출력(에러)
				}
				
				종목_맵, 에러 = f종목_맵_초기화()
				if 에러 != nil {
					공용.F에러_출력(에러)
				}
				
				보유량_맵, 에러 = f종목별_보유량_맵_초기화()
				if 에러 != nil {
					공용.F에러_출력(에러)
				}
				
				문자열_캐시_맵 = make(map[string][]string)
				
				질의.S회신(nil)
			case 공용.P메시지_종료:
				공용_데이터_Go루틴_실행_중.S값(false)
				return
			default:
				에러 = 공용.F에러_생성("예상치 못한 질의 구분 %v", 질의.G구분())
				공용.F에러_출력(에러)
				panic(에러)
			}
		case 질의 = <-Ch종목별_보유량:
			switch 질의.G구분() {
			case 공용.P메시지_GET:
				if 에러 = 질의.G검사(공용.P메시지_GET, 2); 에러 != nil {
					공용.F에러_출력(에러)
					panic(에러)
				}
				
				종목코드 := 질의.G내용(0)
				
				종목별_보유량, 존재함 := 보유량_맵[종목코드]
				
				if !존재함 {
					질의.S회신(공용.F에러_생성("해당 종목을 보유하고 있지 않음. %v", 질의.G내용(0)))
				}
				
				switch 질의.G내용(1) {
				case 공용.P포지션_롱:
					질의.S회신(nil, 종목별_보유량.G롱포지션())
				case 공용.P포지션_숏:
					질의.S회신(nil, 종목별_보유량.G숏포지션())
				default:
					에러 = 공용.F에러_생성("잘못된 포지션 종류. %v", 질의.G내용(1))
					공용.F에러_출력(에러)
					panic(에러)
				}
			case 공용.P메시지_SET:
				공용.F메모("종목별 보유량은 어떻게 저장할 것인가? Gob파일? DB?")
			
				if 에러 = 질의.G검사(공용.P메시지_SET, 3); 에러 != nil {
					공용.F에러_출력(에러)
					panic(에러)
				}
				
				종목코드 := 질의.G내용(0)
				종류 := 질의.G내용(1)
				수량, 에러 := 공용.F2정수(질의.G내용(2))
				
				if 에러 != nil {
					공용.F에러_출력(에러)
					panic(에러)
				}
				
				종목별_보유량, 존재함 := 보유량_맵[종목코드]
				
				if !존재함 {
					종목별_보유량 = 공용.New종목별_보유량(종목코드, 0, 0)
					보유량_맵[종목코드] = 종목별_보유량
				}
				
				switch 종류 {
				case 공용.P포지션_롱:
					종목별_보유량.S더하기_롱포지션(수량)
				case 공용.P포지션_숏:
					종목별_보유량.S더하기_숏포지션(수량)
				default:
					에러 = 공용.F에러_생성("예상치 못한 종류 %v", 종류)
					공용.F에러_출력(에러)
					panic(에러)
				}
			case 공용.P메시지_DEL:
				if 에러 = 질의.G검사(공용.P메시지_DEL, 1); 에러 != nil {
					공용.F에러_출력(에러)
					panic(에러)
				}
				
				delete(보유량_맵, 질의.G내용(0))
			default:
				에러 = 공용.F에러_생성("예상치 못한 질의 구분 %v", 질의.G구분())
				panic(에러)
			}
		case <-공통_종료_채널:
			공용_데이터_Go루틴_실행_중.S값(false)
			return
		}
	}
}

func f주소_맵_초기화() (맵 map[string]string, 에러 error) {
	defer func() {
		if r := recover(); r != nil {
			맵 = nil
			에러 = 공용.F에러_생성("주소 맵 초기화 중 패닉 발생.\n%v", r)
		}
	}()

	맵 = make(map[string]string)

	맵[공용.P주소명_주소정보] = 공용.P주소_주소정보
	맵[공용.P주소명_테스트_결과] = 공용.P주소_테스트_결과

	주소_모음 := make([]string, 0)
	주소_모음 = append(주소_모음, 공용.P주소명_종목정보)
	주소_모음 = append(주소_모음, 공용.P주소명_가격정보)
	주소_모음 = append(주소_모음, 공용.P주소명_가격정보_입수)
	주소_모음 = append(주소_모음, 공용.P주소명_가격정보_배포)

	for i := 0; i < len(주소_모음); i++ {
		맵[주소_모음[i]] = "tcp://127.0.0.1:" + strconv.Itoa(3010+i) // 3010번 포트부터 차례대로 배정.
	}

	return 맵, nil
}

func f종목_맵_초기화() (맵 map[string]공용.I종목, 에러 error) {
	공용.F메모("개발할 것. (현재 샘플 데이터만 넣은 상태)")

	defer func() {
		if r := recover(); r != nil {
			맵 = nil
			에러 = 공용.F에러_생성("종목 맵 초기화 중 패닉 발생.\n%v", r)
		}
	}()

	맵 = make(map[string]공용.I종목)

	// 임시로 샘플 데이터만 사용해서 테스트 할 수 있도록 함.
	종목_모음 := 공용.F샘플_종목_모음()

	for i := 0; i < len(종목_모음); i++ {
		맵[종목_모음[i].G코드()] = 종목_모음[i]
	}

	return 맵, nil
}

func f종목별_보유량_맵_초기화() (맵 map[string]공용.I종목별_보유량, 에러 error) {
	공용.F메모("개발할 것. (현재 비어있는 상태로 시작함.)")

	defer func() {
		if r := recover(); r != nil {
			맵 = nil
			에러 = 공용.F에러_생성("종목별 보유량 맵 초기화 중 패닉 발생.\n%v", r)
		}
	}()

	맵 = make(map[string]공용.I종목별_보유량)

	return 맵, nil
}