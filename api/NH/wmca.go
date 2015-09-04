package NH

// ghts의 bin디렉토리에 있는 sync_ctype.bat에서
// go tool cgo -godefs 를 실행시켜서
// wmca_type.h에 있는 C언어 구조체를 자동으로 Go언어 구조체로 변환시킴.
// 생성된 결과물은 서로 직접 변환(cast)되어도 안전함.
//
//go:generate sync_ctype.bat

// #cgo CFLAGS: -m32 -Wall
// #include <stdlib.h>
// #include "./wmca_func.h"
import "C"

import (
	공용 "github.com/ghts/ghts/common"
	"golang.org/x/sys/windows"

	"time"
	"unsafe"
)

const wmca_dll = "wmca.dll"
const 실행_성공 = "completed successfully"

// SDK함수 종류 혹은 TR구분
const (
	p접속하기 = "접속하기"
	p접속끊기 = "접속끊기"
	p실시간_서비스_모두_해제 = "실시간 서비스 모두 해제"
)

var nh_OpenAPI_Go루틴_실행_중 = 공용.New안전한_bool(false)

var ch실시간_데이터_수신 = make(chan *S수신_데이터_블록, 1000000)
var ch조회 = make(chan 공용.I질의_가변형)
var ch실시간_서비스_등록 = make(chan 공용.I질의_가변형)
var ch실시간_서비스_해제 = make(chan 공용.I질의_가변형)
var ch실시간_서비스_모두_해제 = make(chan 공용.I질의_가변형)
var ch접속됨 = make(chan 공용.I질의_가변형)
var ch접속 = make(chan 공용.I질의_가변형)
var ch접속_해제 = make(chan 공용.I질의_가변형)
var ch자원_정리 = make(chan 공용.I질의_가변형)
var ch종료 = make(chan 공용.S비어있는_구조체, 1)

var ch콜백_로그인 = make(chan *S로그인_정보_블록)
var ch콜백_조회_데이터 = make(chan *S수신_데이터_블록)
var ch콜백_메시지 = make(chan *S수신_메시지_블록)
var ch콜백_완료 = make(chan *S수신_데이터_블록)
var ch콜백_에러 = make(chan *S수신_데이터_블록)
var ch콜백_소켓_에러 = make(chan int)
var ch콜백_접속_해제 = make(chan 공용.S비어있는_구조체)

// NH OpenAPI는 thread-safe 하다고 명시되어 있지 않으므로,
// 다수 스레드(thread) 혹은 Go루틴에서 API를 호출하는 경우, 
// 문제가 생긴다고 봐야 함.
// 1번에 1개의 호출만 처리하도록 하기 위하여, Go루틴을 사용함.

func new대기중_질의(TR코드 string, 질의 공용.I질의_가변형) s대기중_질의 {
	return s대기중_질의 {
		TR식별번호: f_TR식별번호(),
		TR코드: TR코드,
		M질의: 질의,
		M유효기간: time.Now().Add(30 * time.Second)}	
}

type s대기중_질의 struct {
	TR식별번호 uint32
	TR코드 string
	M질의 공용.I질의_가변형
	M유효기간 time.Time
}

func F_NH_OpenAPI_Go루틴(ch초기화 chan bool) {
	// 이미 실행 중인 경우에는 새로 생성하지 않음.
	에러 := nh_OpenAPI_Go루틴_실행_중.S값(true)
	if 에러 != nil {
		ch초기화 <- false; return
	}
	
	defer nh_OpenAPI_Go루틴_실행_중.S값(false)
	
	// 콜백을 받지 못한 호출이 대기하는 임시 저장소
	// 키는 TR식별번호 임.
	대기중_질의_맵 := make(map[uint32]s대기중_질의)
	점검_주기 := time.NewTicker(time.Second)
	
	// 초기화 완료.
	ch초기화 <- true	
	
	for {
		select {
		case 실시간_데이터 := <-ch실시간_데이터_수신:
			f실시간_데이터_수신_처리(실시간_데이터)
		case 질의 := <-ch조회:
			대기중_질의 := f조회_질의_처리(질의)
			if 대기중_질의 != nil {
				대기중_질의_맵[대기중_질의.TR식별번호] = *대기중_질의
			}
		case 질의 := <-ch실시간_서비스_등록:
			대기중_질의 := f실시간_서비스_등록_질의_처리(질의)
			if 대기중_질의 != nil {
				대기중_질의_맵[대기중_질의.TR식별번호] = *대기중_질의
			}
		case 질의 := <-ch실시간_서비스_해제:
			대기중_질의 := f실시간_서비스_해제_질의_처리(질의)
			if 대기중_질의 != nil {
				대기중_질의_맵[대기중_질의.TR식별번호] = *대기중_질의
			}
		case 질의 := <-ch실시간_서비스_모두_해제:
			대기중_질의 := f실시간_서비스_모두_해제_질의_처리(질의)
			if 대기중_질의 != nil {
				대기중_질의_맵[대기중_질의.TR식별번호] = *대기중_질의
			}
		case 질의 := <-ch접속됨:
			질의.S회신(nil, f접속됨())
		case 질의 := <-ch접속:
			대기중_질의 := f접속_질의_처리(질의)
			if 대기중_질의 != nil {
				대기중_질의_맵[대기중_질의.TR식별번호] = *대기중_질의
			}
		case 질의 := <-ch접속_해제:
			대기중_질의 := f접속_해제_질의_처리(질의)
			if 대기중_질의 != nil {
				대기중_질의_맵[대기중_질의.TR식별번호] = *대기중_질의
			}
		case 질의 := <-ch자원_정리:
			f자원_정리()
			질의.S회신(nil, nil)
		case 로그인_데이터 := <-ch콜백_로그인:	
			대기중_질의, 존재함 := 대기중_질의_맵[로그인_데이터.TR식별번호]
			
			if !존재함 {
				에러 := 공용.F에러_생성("로그인 데이터를 수신하였으나, 대기 중인 질의가 존재하지 않음.")
				공용.F에러_출력(에러)
				panic(에러)
			}
			
			대기중_질의.M질의.S회신(nil, 로그인_데이터.M로그인_정보)
			delete(대기중_질의_맵, 로그인_데이터.TR식별번호)
		case 데이터 := <-ch콜백_조회_데이터:
			대기중_질의, 존재함 := 대기중_질의_맵[데이터.TR식별번호]
			
			switch {
			case !존재함:
				에러 := 공용.F에러_생성("조회 데이터를 수신하였으나, 이에 해당되는 질의가 존재하지 않음.")
				공용.F에러_출력(에러)
				panic(에러)
			default:
				대기중_질의.M질의.S회신(nil, 데이터.M수신_데이터)
			}
			
			delete(대기중_질의_맵, 데이터.TR식별번호)
		case 데이터 := <-ch콜백_메시지:
			메시지 := 공용.F포맷된_문자열("코드 : %v\n%v", 데이터.M메시지_코드, 데이터.M메시지_내용)
			공용.F문자열_출력(메시지)
			
			대기중_질의, 존재함 := 대기중_질의_맵[데이터.TR식별번호]
			
			switch {
			case !존재함:
				에러 := 공용.F에러_생성("메시지를 수신하였으나, 대기 중인 질의가 존재하지 않음.")
				공용.F에러_출력(에러)
				panic(에러)
			case 대기중_질의.TR코드 == p접속하기:
				대기중_질의.M질의.S회신(공용.F에러_생성(메시지), nil)
			default:
				에러 := 공용.F에러_생성("메시지를 수신하였으나, 예상치 못한 경우임.")
				공용.F에러_출력(에러)
				panic(에러)
			}
			
			delete(대기중_질의_맵, 데이터.TR식별번호)
		case 데이터 := <-ch콜백_완료:
			대기중_질의, 존재함 := 대기중_질의_맵[데이터.TR식별번호]
			
			switch {
			case !존재함:
				에러 := 공용.F에러_생성("완료 메시지를 수신하였으나, 이에 해당되는 질의가 존재하지 않음.")
				공용.F에러_출력(에러)
				panic(에러)
			default:
				대기중_질의.M질의.S회신(nil, 데이터.M수신_데이터)
			}
			
			delete(대기중_질의_맵, 데이터.TR식별번호)
		case 데이터 := <-ch콜백_에러:	
			수신_에러 := 공용.F에러_생성("에러 발생\n%v", 데이터.M수신_데이터.M데이터)
			공용.F에러_출력(수신_에러)
			
			대기중_질의, 존재함 := 대기중_질의_맵[데이터.TR식별번호]
			
			switch {
			case !존재함:
				에러 := 공용.F에러_생성("에러가 발생하였으나, 이에 해당되는 질의가 존재하지 않음.")
				공용.F에러_출력(에러)
				panic(에러)
			default:
				대기중_질의.M질의.S회신(수신_에러, nil)
			}
			
			delete(대기중_질의_맵, 데이터.TR식별번호)
		case 에러_코드 := <-ch콜백_소켓_에러:
			소켓_에러 := 공용.F에러_생성("소켓 에러 발생. 에러코드 : %v", 에러_코드)
			공용.F에러_출력(에러)
			
			// 모든 대기 중 질의에 대해서 에러 회신?
			for 식별번호, 대기중_질의 := range 대기중_질의_맵 {
				대기중_질의.M질의.S회신(소켓_에러, nil)
				delete(대기중_질의_맵, 식별번호)
			}
		case <-ch콜백_접속_해제:
			접속_해제_에러 := 공용.F에러_생성("접속 해제됨. %v", time.Now())
			공용.F에러_출력(에러)
			
			for 식별번호, 대기중_질의 := range 대기중_질의_맵 {
				switch {
				case 대기중_질의.TR코드 == p접속끊기:
					// 접속 해제 요청이 성공했으므로, 에러가 아님
					대기중_질의.M질의.S회신(nil, 접속_해제_에러.Error())
				default:
					// 나머지 모든 경우에 대해서 에러 회신.
					대기중_질의.M질의.S회신(접속_해제_에러, nil)
				}
				
				delete(대기중_질의_맵, 식별번호)
			}
		case <-점검_주기.C: // 유효기간이 지난 대기중 질의 삭제			
			지금 := time.Now()
			
			for tr식별번호, 대기중_질의 := range 대기중_질의_맵 {
				if 대기중_질의.M유효기간.Before(지금) {
					delete(대기중_질의_맵, tr식별번호)
				}
			}
		case <-공용.F공통_종료_채널():
			ch종료 <- 공용.S비어있는_구조체{}
		case <-ch종료:
			go func() {
				공용.New질의_가변형(공용.P메시지_GET).G회신(ch실시간_서비스_모두_해제)
				공용.New질의_가변형(공용.P메시지_GET).G회신(ch접속_해제)
				공용.New질의_가변형(공용.P메시지_GET).G회신(ch자원_정리)
			}
		default:
			// 얼마나 대기해야 하나?
			if len(대기중_질의_맵) == 0 {
				time.Sleep(50 * time.Millisecond)
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func fDLL존재함() bool {
	에러 := windows.NewLazyDLL(wmca_dll).Load()

	if 에러 != nil {
		return false
	} else {
		return true
	}
}

func f실시간_데이터_수신_처리(실시간_데이터 *S수신_데이터_블록) {
	블록_이름 := s.M수신_데이터.M블록_이름
	
	공용.F메모("실시간 데이터 종류에 따라서 적절히 처리할 것.")
	
	switch 블록_이름 {
	default:
		에러 := 공용.F에러_생성("예상치 못한 블록 이름 %v", 블록_이름)
		공용.F에러_출력(에러)
		panic(에러)
	}
}

func f조회_질의_처리(질의 공용.I질의_가변형) *s대기중_질의 {
	여기
}

func f조회(TR식별번호 int, TR코드 string, 데이터_포인터 unsafe.Pointer, 길이 int, 계좌_인덱스 int) bool {
	cTR식별번호 := C.int(TR식별번호)
	cTR코드 := C.CString(TR코드)
	c데이터 := (*C.char)(데이터_포인터)
	c길이 := C.int(길이)
	c계좌_인덱스 := C.int(계좌_인덱스)

	defer func() {
		C.free(unsafe.Pointer(cTR코드))
		C.free(unsafe.Pointer(c데이터)) // C언어 구조체로 변환된 후에는 직접 free 해 줘야 하는 듯.
	}()

	반환값 := C.wmcaQuery(cTR식별번호, cTR코드, c데이터, c길이, c계좌_인덱스)

	return bool(반환값)
}

func f실시간_서비스_등록(타입 string, 코드_모음 string, 코드_길이 int, 전체_길이 int) bool {
	c타입 := C.CString(타입)
	c코드_모음 := C.CString(코드_모음)
	c코드_길이 := C.int(코드_길이)
	c전체_길이 := C.int(전체_길이)

	defer func() {
		C.free(unsafe.Pointer(c타입))
		C.free(unsafe.Pointer(c코드_모음))
	}()

	반환값 := C.wmcaAttach(c타입, c코드_모음, c코드_길이, c전체_길이)

	return bool(반환값)
}

func f실시간_서비스_해제(타입 string, 코드_모음 string, 코드_길이 int, 전체_길이 int) bool {
	c타입 := C.CString(타입)
	c코드_모음 := C.CString(코드_모음)
	c코드_길이 := C.int(코드_길이)
	c전체_길이 := C.int(전체_길이)

	defer func() {
		C.free(unsafe.Pointer(c타입))
		C.free(unsafe.Pointer(c코드_모음))
	}()

	반환값 := C.wmcaDetach(c타입, c코드_모음, c코드_길이, c전체_길이)

	return bool(반환값)
}

func f접속_질의_처리(질의 공용.I질의_가변형) *s대기중_질의 {
	질의_에러 := 질의.G검사(공용.P메시지_GET, 3)
	이미_접속됨 := f접속됨()
	
	switch {
	case 질의_에러 != nil:
		질의.S회신(질의_에러, nil)
		return nil
	case 이미_접속됨:
		에러 := 공용.F에러_생성("이미 접속되어 있음.")
		공용.F에러_출력(에러)
		질의.S회신(nil, nil)
		return nil
	case !f접속(질의.G내용(0).(string), 질의.G내용(1).(string), 질의.G내용(2).(string)):
		// 아이디 - 질의.G내용(0).(string), 암호 - 질의.G내용(1).(string), 공인인증 암호 - 질의.G내용(2).(string) 	  
		에러 := 공용.F에러_생성("접속 실패.")
		공용.F에러_출력(에러)
		질의.S회신(에러, nil)
		return nil
	}
	
	대기중_질의 := new대기중_질의(p접속하기, 질의)
	
	return &대기중_질의
} 

func f접속(아이디, 암호, 공인인증서_암호 string) bool {
	c아이디 := C.CString(아이디)
	c암호 := C.CString(암호)
	c공인인증서_암호 := C.CString(공인인증서_암호)

	defer func() {
		C.free(unsafe.Pointer(c아이디))
		C.free(unsafe.Pointer(c암호))
		C.free(unsafe.Pointer(c공인인증서_암호))
	}()

	return bool(C.wmcaConnect(c아이디, c암호, c공인인증서_암호))
}

func f접속됨() bool {
	return f호출("wmcaIsConnected")
}

func f실시간_서비스_모두_해제() bool {
	return f호출("wmcaDetachAll")
}

func f자원_정리() {
	// cgo의 버그로 인해서 인수가 없으면 '사용하지 않는 변수' 컴파일 경고 발생.
	// 컴파일 경고를 없애기 위해서 사용하지 않는 인수를 추가함.
	C.wmcaFreeResource(C.int(1))
}

func f접속끊기() bool {
	return f호출("wmcaDisconnect")
}