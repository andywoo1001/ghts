package internal

import (
	공용 "github.com/ghts/ghts/common"

	"strings"
	"testing"
)

func f접속_확인() {
	if !f접속됨() {
		질의 := 공용.New질의_가변형(P30초, 공용.P메시지_GET, 테스트용_ID, 테스트용_암호, 테스트용_공인인증_암호)
		질의.S질의(Ch접속)
		질의.G회신()	// 메시지
		질의.G회신()	// 로그인 정보
		
		if !f접속됨() {
			panic("접속 시도 실패")
		}
	}
}

func f등락부호2정수(등락부호 uint8) int64 {
	switch 등락부호 {
	case P상한, P상승:
		return int64(1)
	case P보합:
		return int64(0)
	case P하락, P하한:
		return int64(-1)
	default:
		에러 := 공용.F에러_생성("등락부호가 예상된 값이 아님. %v", 등락부호)
		공용.F에러_출력(에러)
		panic(에러)
	}
}

func f테스트_등락부호(테스트 *testing.T, 등락부호 uint8, 값, 비교대상, 상한, 하한 int64) {
	switch 등락부호 {
	case P상한:
		공용.F테스트_같음(테스트, 값, 상한)
	case P상승:
		공용.F테스트_참임(테스트, 값 > 비교대상)
	case P보합:
		공용.F테스트_참임(테스트, 값 == 비교대상)
	case P하락:
		공용.F테스트_참임(테스트, 값 < 비교대상)
	case P하한:
		공용.F테스트_같음(테스트, 값, 하한)
	default:
		공용.F문자열_출력("등락부호가 예상된 값이 아님. %v", 등락부호)
		테스트.FailNow()
	}
}

func f테스트_등락율(테스트 *testing.T, 부호 uint8, 등락율 float64) {
	switch 부호 {
	case P상한, P상승:
		공용.F테스트_참임(테스트, 등락율 > 0)
	case P보합:
		// 이게 구체적으로 무슨 의미?? 일단은 임의로 변동폭 10% 이내라고 가정함. 
		공용.F테스트_참임(테스트, 등락율 < 10 && 등락율 > -10)
	case P하락, P하한:
		공용.F테스트_참임(테스트, 등락율 < 0)
	default:
		공용.F문자열_출력("등락부호가 예상된 값이 아님. %v", 부호)
		테스트.FailNow()
	}
}

func Test_Ch접속(테스트 *testing.T) {
	if f접속됨() {
		공용.New질의_가변형(P30초, 공용.P메시지_GET).S질의(Ch접속_해제).G회신()
	}

	질의 := 공용.New질의_가변형(P30초, 공용.P메시지_GET, 테스트용_ID, 테스트용_암호, 테스트용_공인인증_암호)
	질의.S질의(Ch접속)
	
	for 질의.G회신_종료() {
		회신 := 질의.G회신()
		공용.F테스트_에러없음(테스트, 회신.G에러())
		
		switch 회신.G구분() {
		case P회신_접속:
			공용.F테스트_같음(테스트, 회신.G길이(), 1)
			
			로그인_정보, ok := 회신.G내용(0).(S로그인_정보)
			공용.F테스트_참임(테스트, ok)
			공용.F테스트_같음(테스트, 로그인_정보.G접속_ID(), 테스트용_ID)
			공용.F테스트_참임(테스트, strings.HasPrefix(로그인_정보.G접속_서버(), "mt")) //??
			공용.F테스트_참임(테스트, len(로그인_정보.G계좌_목록()) > 0)
		
			테스트용_계좌_인덱스 := -1
			for 인덱스, 계좌_정보 := range 로그인_정보.G계좌_목록() {
				if 계좌_정보.G계좌_번호() == 테스트용_계좌_번호 {
					테스트용_계좌_인덱스 = 인덱스 + 1 //계좌번호 인덱스는 '1'부터 시작됩니다.
					break
				}
			}
		
			공용.F테스트_참임(테스트, 테스트용_계좌_인덱스 > 0)
			
			// 접속 되었는 지 확인
			공용.F테스트_참임(테스트, f접속됨())
			
			회신 := 공용.New질의_가변형(P30초, 공용.P메시지_GET).S질의(Ch접속됨).G회신()
			공용.F테스트_에러없음(테스트, 회신.G에러())
			공용.F테스트_참임(테스트, 회신.G내용(0).(bool))
		case P회신_메시지:
			공용.F테스트_같음(테스트, 회신.G길이(), 2)
			
			_, ok := 회신.G내용(0).(string)	// 코드
			공용.F테스트_참임(테스트, ok)
			
			메시지, ok := 회신.G내용(1).(string)	// 메시지
			공용.F테스트_참임(테스트, ok)
			
			공용.F테스트_참임(테스트, strings.Contains(메시지, "로그인"))
			공용.F테스트_참임(테스트, strings.Contains(메시지, "성공"))
		default:
			공용.F문자열_출력("예상하지 못한 회신 종류 %v", 회신.G구분())
			공용.F변수값_확인(회신)
			테스트.FailNow()
		}
	}	
}

func TestCh접속_해제(테스트 *testing.T) {
	f접속_확인()
	
	// 접속 해제
	회신 := 공용.New질의_가변형(P30초, 공용.P메시지_GET).S질의(Ch접속_해제).G회신()
	공용.F테스트_에러없음(테스트, 회신.G에러())
	공용.F테스트_같음(테스트, 회신.G구분(), P회신_접속_해제)
	공용.F테스트_참임(테스트, 회신.G내용(0).(bool))
	
	// 접속 해제 확인
	회신 = 공용.New질의_가변형(P30초, 공용.P메시지_GET).S질의(Ch접속됨).G회신()
	공용.F테스트_에러없음(테스트, 회신.G에러())
	공용.F테스트_같음(테스트, 회신.G구분(), P회신_접속됨)
	공용.F테스트_거짓임(테스트, 회신.G내용(0).(bool), false)
}
