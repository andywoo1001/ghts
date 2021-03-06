package internal

// #cgo CFLAGS: -m32 -Wall
// #include <stdlib.h>
// #include <windows.h>
// #include "./c_func.h"
import "C"

import (
	공용 "github.com/ghts/ghts/common"
	"golang.org/x/sys/windows"
	"math"
	"strings"
	"testing"
	"unsafe"
)

var 질의_식별번호 = 공용.New안전한_일련_번호()

func f2등락부호(바이트_모음 [1]byte) uint8 {
	값 := uint8(바이트_모음[0])

	switch 값 {
	case P상한, P상승, P보합, P하한, P하락:
		return 값
	// 이하 내용은 게시판에 답변에 따름.
	case 1, 6:
		return P상한
	case 2, 7:
		return P상승
	case 3, 0:
		return P보합
	case 4, 8:
		return P하한
	case 5, 9:
		return P하락
	}

	go문자열 := 공용.F2문자열(바이트_모음)
	//go문자열 := 공용.F2문자열(값)

	switch go문자열 {
	case "1", "6":
		return P상한
	case "2", "7":
		return P상승
	case "3", "0":
		return P보합
	case "4", "8":
		return P하한
	case "5", "9":
		return P하락
	}

	공용.F변수값_확인(바이트_모음)
	공용.F변수값_확인(값)
	공용.F변수값_확인(go문자열)
	에러 := 공용.F에러("예상치 못한 등락부호 값.")
	panic(에러)

	return 0xFF
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
		에러 := 공용.F에러("등락부호가 예상된 값이 아님. %v", 등락부호)
		panic(에러)
	}
}

func f올바른_등락부호(값 uint8) bool {
	switch 값 {
	case P상한, P상승, P보합, P하락, P하한:
		return true
	default:
		return false
	}
}

func f테스트_등락부호(테스트 *testing.T, 등락부호 uint8, 값, 비교대상, 상한, 하한 int64) {
	switch 등락부호 {
	case P상한:
		공용.F테스트_같음(테스트, 값, 상한)
	case P상승:
		공용.F테스트_참임(테스트, 값 > 비교대상, 값, 비교대상)
	case P보합:
		if 값 != 0 && 비교대상 != 0 {
			공용.F테스트_같음(테스트, 값, 비교대상)
		}
	case P하락:
		공용.F테스트_참임(테스트, 값 < 비교대상, 값, 비교대상)
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

func f2실수_소숫점_추가(값 interface{}, 소숫점_이하_자릿수 int) float64 {
	문자열 := strings.Replace(공용.F2문자열(값), " ", "0", -1)

	if len(문자열) < 소숫점_이하_자릿수 {
		에러 := 공용.F에러("문자열 길이가 소숫점_이하_자릿수 보다 짧습니다.")
		panic(에러)
	}

	소숫점_추가_문자열 := ""
	
	switch {
	case strings.Contains(문자열, "INF"):
		return math.Inf(1)
	case strings.Contains(문자열, "."):
		소숫점_추가_문자열 = 문자열
	default:
		소숫점_추가_문자열 = 문자열[:len(문자열)-소숫점_이하_자릿수] + "." + 
			문자열[len(문자열)-소숫점_이하_자릿수:]
	}

	return 공용.F2실수(소숫점_추가_문자열)
}

func f바이트2참거짓(값 []byte, 조건 string, 결과 bool) bool {
	if string(값) == 조건 {
		return 결과
	}

	return !결과
}

func f_Go구조체로_변환(c *C.RECEIVED) interface{} {
	// 반대로 변환할 때는 (*C.char)(unsafe.Pointer(&b[0]))

	공용.F문자열_출력("블록 이름 : %v", C.GoString(c.BlockName))

	g := (*Received)(unsafe.Pointer(c))

	블록_이름 := C.GoString(c.BlockName)
	//전체_길이 := int(c.Length)
	전체_길이 := int(g.Length)
	데이터 := c.DataString

	if 전체_길이 == 0 {
		return nil
	}

	switch 블록_이름 {
	case "c1101OutBlock":
		c := (*C.Tc1101OutBlock)(unsafe.Pointer(데이터))
		return New주식_현재가_조회_기본_자료(c)
	case "c1101OutBlock2":
		수량 := 전체_길이 / int(unsafe.Sizeof(C.Tc1101OutBlock2{}))

		// 큰 배열로 캐스팅 한 다음에 슬라이스를 취함.
		// 충분히 큰 숫자이면 아무 것이나 상관없으며, 반드시 반드시 10000이어야 하는 것은 아님.
		// Go위키에서는 '1 << 30'을 사용하지만, 너무 큰 수를 사용하니까 메모리 범위를 벗어난다고 에러 발생.
		슬라이스 := (*[10000]C.Tc1101OutBlock2)(unsafe.Pointer(데이터))[:수량:수량]
		go슬라이스 := make([]S주식_현재가_조회_변동_거래량_자료, 수량)

		for i := 0; i < 수량; i++ {
			c := 슬라이스[i]
			g := New주식_현재가_조회_변동_거래량_자료(&c)
			
			go슬라이스[i] = *g
			//C.free(unsafe.Pointer(&c))
		}

		return go슬라이스
	case "c1101OutBlock3":
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tc1101OutBlock3{}))
		c := (*C.Tc1101OutBlock3)(unsafe.Pointer(데이터))
		return New주식_현재가_조회_동시호가(c)
	case "c1151OutBlock":
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tc1151OutBlock{}))
		return New_ETF_현재가_조회_기본_자료(데이터)
	case "c1151OutBlock2":
		수량 := 전체_길이 / int(unsafe.Sizeof(C.Tc1151OutBlock2{}))

		// 큰 배열로 캐스팅 한 다음에 슬라이스를 취함.
		// 충분히 큰 숫자이면 아무 것이나 상관없으며, 반드시 반드시 10000이어야 하는 것은 아님.
		// Go위키에서는 '1 << 30'을 사용하지만, 너무 큰 수를 사용하니까 메모리 범위를 벗어난다고 에러 발생.
		슬라이스 := (*[10000]C.Tc1151OutBlock2)(unsafe.Pointer(데이터))[:수량:수량]
		go슬라이스 := make([]S_ETF_현재가_조회_변동_거래_자료, 수량)

		for i := 0; i < 수량; i++ {
			c := 슬라이스[i]
			g := New_ETF_현재가_조회_변동_거래_자료(&c)
			
			go슬라이스[i] = *g
			//C.free(unsafe.Pointer(&c))
		}

		return go슬라이스
	case "c1151OutBlock3":
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tc1151OutBlock3{}))
		return New_ETF_현재가_조회_동시호가(데이터)
	case "c1151OutBlock4":
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tc1151OutBlock4{}))
		return New_ETF_현재가_조회_ETF자료(데이터)
	case "c1151OutBlock5":
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tc1151OutBlock5{}))
		return New_ETF_현재가_조회_지수_자료(데이터)
	case "h1OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Th1OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Th1OutBlock{}))
		return New코스피_호가_잔량(데이터)
	case "k3OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tk3OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tk3OutBlock{}))
		return New코스닥_호가_잔량(데이터)
	case "h2OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Th2OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Th2OutBlock{}))
		return New코스피_시간외_호가_잔량(데이터)
	case "k4OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tk4OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tk4OutBlock{}))
		return New코스닥_시간외_호가_잔량(데이터)
	case "h3OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Th3OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Th3OutBlock{}))
		return New코스피_예상_호가_잔량(데이터)
	case "k5OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tk5OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tk5OutBlock{}))
		return New코스닥_예상_호가_잔량(데이터)
	case "j8OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tj8OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tj8OutBlock{}))
		return New코스피_체결(데이터)
	case "k8OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tk8OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tk8OutBlock{}))
		return New코스닥_체결(데이터)
	case "j1OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tj1OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tj1OutBlock{}))
		return New코스피_ETF(데이터)
	case "j0OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tj0OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tj0OutBlock{}))
		return New코스닥_ETF(데이터)
	case "u1OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tu1OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tu1OutBlock{}))
		return New코스피_업종_지수(데이터)
	case "k1OutBlock":
		공용.F문자열_출력("전체 길이, sizeof : %v, %v", 전체_길이, unsafe.Sizeof(C.Tk1OutBlock{}))
		f반복되면_패닉(블록_이름, 전체_길이, unsafe.Sizeof(C.Tk1OutBlock{}))
		return New코스닥_업종_지수(데이터)
	default:
		에러 := 공용.F에러("예상치 못한 블록 이름 %v", 블록_이름)
		panic(에러)
	}

	return nil
}

func f반복되면_패닉(블록_이름 string, 전체_길이 int, 구조체_길이 uintptr) {
	if 전체_길이 == 0 && 구조체_길이 > 0 {
		공용.F문자열_출력("데이터 길이가 0임. 데이터 구조체 형식이 잘못됨.")
		// '전체_길이' 값이 제대로 수신되지 않음.
		return
	}

	수량 := 전체_길이 / int(구조체_길이)

	if 수량 != 1 {
		에러 := 공용.F에러("반복되는 구조체임. %v", 블록_이름)
		panic(에러)
	}
}

func fHWND() C.HWND {
	return C.getHWND()
}

func f호출(함수명 string, 인수 ...uintptr) bool {
	if !fDLL존재함() {
		return false
	}

	// Call()의 2번째 반환값은 '윈도우 + C언어'조합에서는 필요없는 듯함.
	// 인터넷에서 찾은 예제 코드들은 모두 2번째 반환값을 '_' 처리함.
	반환값, _, 에러 := windows.NewLazyDLL(wmca_dll).NewProc(함수명).Call(인수...)

	// C언어에서 BOOL의 정의는 0이면 false,그 이외의 값은 true임.
	// 일반적인 프로그래밍 언어는 true부터 먼저 확인해도 되지만
	// C언어의 BOOL은 0인지 (즉, false인지)부터 확인해야 함. (순서에 유의)
	switch {
	case !strings.Contains(에러.Error(), 실행_성공):
		공용.F에러_출력(에러)
		return false
	case 반환값 == 0:
		return false
	default:
		return true
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

func f접속_안_되어_있으면_에러(질의 공용.I질의_가변형) error {
	if !f접속됨() {
		에러 := 공용.F에러("접속되지 않음")
		질의.S회신(에러, 공용.P메시지_에러)

		return 에러
	}

	return nil
}

func f조회(TR식별번호 int, TR코드 string, c데이터 *C.char, c길이 C.int, 계좌_인덱스 int) bool {
	cTR식별번호 := C.int(TR식별번호)
	cTR코드 := C.CString(TR코드)
	c계좌_인덱스 := C.int(계좌_인덱스)

	defer C.free(unsafe.Pointer(cTR코드))
	
	// C언어의 원칙에 따라서 변수를 생성한 곳에서 free()하도록 한다.
	//defer C.free(unsafe.Pointer(c데이터))
	
	반환값 := C.wmcaQuery(cTR식별번호, cTR코드, c데이터, c길이, c계좌_인덱스)

	return bool(반환값)
}

func f실시간_서비스_등록(타입 string, 코드_모음 string, 단위_길이 int, 전체_길이 int) bool {
	c타입 := C.CString(타입)
	c코드_모음 := C.CString(코드_모음)
	c단위_길이 := C.int(단위_길이)
	c전체_길이 := C.int(전체_길이)

	defer func() {
		C.free(unsafe.Pointer(c타입))
		C.free(unsafe.Pointer(c코드_모음))
	}()

	반환값 := C.wmcaAttach(c타입, c코드_모음, c단위_길이, c전체_길이)

	return bool(반환값)
}

func f실시간_서비스_해제(타입 string, 코드_모음 string, 단위_길이 int, 전체_길이 int) bool {
	c타입 := C.CString(타입)
	c코드_모음 := C.CString(코드_모음)
	c단위_길이 := C.int(단위_길이)
	c전체_길이 := C.int(전체_길이)

	defer func() {
		C.free(unsafe.Pointer(c타입))
		C.free(unsafe.Pointer(c코드_모음))
	}()

	반환값 := C.wmcaDetach(c타입, c코드_모음, c단위_길이, c전체_길이)

	return bool(반환값)
}

func f접속(아이디, 암호, 공인인증서_암호 string) bool {
	f자원_정리()

	c아이디 := C.CString(아이디)
	c암호 := C.CString(암호)
	c공인인증서_암호 := C.CString(공인인증서_암호)

	defer func() {
		C.free(unsafe.Pointer(c아이디))
		C.free(unsafe.Pointer(c암호))
		C.free(unsafe.Pointer(c공인인증서_암호))
	}()

	서버_이름 := (*C.char)(unsafe.Pointer(nil))
	포트_번호 := 0

	if 공용.F테스트_모드_실행_중() {
		//공용.F문자열_출력("테스트용 모의 서버")
		서버_이름 = C.CString("newmt.wontrading.com")
		포트_번호 = 8400
	} else {
		공용.F문자열_출력("거래 서버")
		서버_이름 = C.CString("wmca.wontrading.com")
		포트_번호 = 8200
	}

	defer C.free(unsafe.Pointer(서버_이름))

	로드_성공 := bool(C.wmcaLoad())
	if !로드_성공 {
		공용.F문자열_출력("로드 실패")
		return false
	}

	서버_설정_성공 := bool(C.wmcaSetServer(서버_이름))
	if !서버_설정_성공 {
		공용.F문자열_출력("서버 설정 실패")
		return false
	}

	포트_설정_성공 := bool(C.wmcaSetPort(C.int(포트_번호)))
	if !포트_설정_성공 {
		공용.F문자열_출력("포트 설정 실패")
		return false
	}

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

func f접속_해제() bool {
	return f호출("wmcaDisconnect")
}
