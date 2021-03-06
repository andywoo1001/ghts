package internal

// ghts의 bin디렉토리에 있는 sync_ctype.bat에서
// go tool cgo -godefs 를 실행시켜서
// wmca_type.h에 있는 C언어 구조체를 자동으로 Go언어 구조체로 변환시킴.
// 생성된 결과물은 서로 직접 변환(cast)되어도 안전함.
//go:generate sync_ctype.bat

// #include <stdlib.h>
// #include "./c_func.h"
import "C"

import (
	공용 "github.com/ghts/ghts/common"

	"strings"
	"time"
	"unsafe"
)

type s콜백_대기 struct {
	식별번호  int
	질의_종류 int
	질의_코드 string // TR코드
	질의    공용.I질의_가변형
	유효기간  time.Time
}

func (this s콜백_대기) G식별번호() int       { return this.식별번호 }
func (this s콜백_대기) G질의_종류() int      { return this.질의_종류 }
func (this s콜백_대기) G질의_코드() string   { return this.질의_코드 }
func (this s콜백_대기) G질의() 공용.I질의_가변형  { return this.질의 }
func (this s콜백_대기) G유효기간() time.Time { return this.유효기간 }

func New콜백_대기(질의_종류 int, TR코드 string, 질의 공용.I질의_가변형) s콜백_대기 {
	return s콜백_대기{
		식별번호:  질의_식별번호.G값(),
		질의_종류: 질의_종류,
		질의_코드: TR코드,
		질의:    질의,
		유효기간:  time.Now().Add(30 * time.Second)}
}

//----------------------------------------------------------------------//
// WMCA_CONNECTED 로그인 구조체
//----------------------------------------------------------------------//
type S로그인_정보 struct {
	접속_시각 time.Time
	접속_서버 string
	접속_ID string
	계좌_목록 []S계좌_정보
}

func (this S로그인_정보) G접속_시각() time.Time { return this.접속_시각 }
func (this S로그인_정보) G접속_서버() string    { return this.접속_서버 }
func (this S로그인_정보) G접속_ID() string    { return this.접속_ID }
func (this S로그인_정보) G계좌_목록() []S계좌_정보  { return this.계좌_목록 }

func New로그인_정보(c블록 *C.LOGINBLOCK) *S로그인_정보 {
	defer C.free(unsafe.Pointer(c블록))

	c := (*C.LOGININFO)(unsafe.Pointer(c블록.LoginInfo))
	g := (*LoginInfo)(unsafe.Pointer(c블록.LoginInfo))

	시각, 에러 := 공용.F2포맷된_시각("20060102150405", g.Date)
	공용.F에러_패닉(에러)
	
	계좌_수량 := 공용.F2정수(공용.F2문자열(g.AccountCount))

	계좌_목록 := make([]S계좌_정보, 계좌_수량)
	for i, c계좌_정보 := range c.Accountlist[:계좌_수량] {
		계좌_목록[i] = new계좌_정보(&c계좌_정보)
	}

	s := new(S로그인_정보)
	s.접속_시각 = 시각
	s.접속_서버 = 공용.F2문자열(g.ServerName)
	s.접속_ID = 공용.F2문자열(g.UserID)
	s.계좌_목록 = 계좌_목록

	return s
}

type S계좌_정보 struct {
	계좌_번호     string
	계좌명       string
	상품_코드     string
	관리점_코드    string
	위임_만기일    time.Time
	일괄주문_허용계좌 bool // ('G' = 허용)
	주석        string
}

func (this S계좌_정보) G계좌_번호() string     { return this.계좌_번호 }
func (this S계좌_정보) G계좌명() string       { return this.계좌명 }
func (this S계좌_정보) G상품_코드() string     { return this.상품_코드 }
func (this S계좌_정보) G관리점_코드() string    { return this.관리점_코드 }
func (this S계좌_정보) G위임_만기일() time.Time { return this.위임_만기일 }
func (this S계좌_정보) G일괄주문_허용계좌() bool   { return this.일괄주문_허용계좌 }
func (this S계좌_정보) G주석() string        { return this.주석 }

func new계좌_정보(c *C.ACCOUNTINFO) S계좌_정보 {
	g := (*AccountInfo)(unsafe.Pointer(c))

	위임_만기일, 에러 := time.Parse("20060102", 공용.F2문자열(g.ExpirationDate))
	if 에러 != nil {
		if 공용.F테스트_모드_실행_중() {
			// 테스트 서버에서는 빈 문자열만 나옴.
			위임_만기일 = time.Now()
		} else {
			panic(에러)
		}
	}

	일괄주문_허용계좌 := false
	if 공용.F2문자열(g.Granted) == "G" {
		일괄주문_허용계좌 = true
	}

	s := new(S계좌_정보)
	s.계좌_번호 = 공용.F2문자열(g.AccountNo)
	s.계좌명 = 공용.F2문자열(g.AccountName)
	s.상품_코드 = 공용.F2문자열(g.AccountProductCode)
	s.관리점_코드 = 공용.F2문자열(g.AmnTabCode)
	s.위임_만기일 = 위임_만기일
	s.일괄주문_허용계좌 = 일괄주문_허용계좌
	s.주석 = 공용.F2문자열(g.Filler)

	return *s
}

//----------------------------------------------------------------------//
// WMCA 문자 message 구조체
//----------------------------------------------------------------------//
type S수신_메시지_블록 struct {
	식별번호   int
	메시지_코드 string //00000:정상, 기타:비정상(코드값은 언제든지 변경될 수 있음.)
	메시지_내용 string
}

func (this S수신_메시지_블록) G식별번호() int      { return this.식별번호 }
func (this S수신_메시지_블록) G메시지_코드() string { return this.메시지_코드 }
func (this S수신_메시지_블록) G메시지_내용() string { return this.메시지_내용 }

func New수신_메시지_블록(c블록 *C.OUTDATABLOCK) S수신_메시지_블록 {
	defer C.free(unsafe.Pointer(c블록))

	const 메시지_최대_길이 = 40

	c := (*C.MSGHEADER)(unsafe.Pointer(c블록.DataStruct.DataString))
	g := (*MsgHeader)(unsafe.Pointer(c블록.DataStruct.DataString))

	코드 := 공용.F2문자열(g.MsgCode)

	공용.F메모("여기에서 계속 메모리 오류가 발생한다. 원인이 뭘까?")
	바이트_모음 := C.GoBytes(unsafe.Pointer(&c.UsrMsg), C.int(len(c.UsrMsg[:])))
	//바이트_모음 := C.GoBytes(unsafe.Pointer(&c.UsrMsg), C.int(len(c.UsrMsg)))
	//if len(바이트_모음) > len(c.UsrMsg[:]) {
	//	바이트_모음 = 바이트_모음[:]
	//}
	
	
	
	메시지 := 공용.F2문자열_CP949(바이트_모음)

	s := new(S수신_메시지_블록)
	s.식별번호 = int(c블록.TrIdNo)
	s.메시지_코드 = 코드
	s.메시지_내용 = 메시지

	return *s
}

//----------------------------------------------------------------------//
// WMCA TR 응답 구조체
//----------------------------------------------------------------------//
type S수신_데이터_블록 struct {
	식별번호   int
	수신_데이터 S수신_데이터
}

func (this S수신_데이터_블록) G식별번호() int { return this.식별번호 }
func (this S수신_데이터_블록) G수신_데이터() S수신_데이터 {
	return this.수신_데이터
}

func New수신_데이터_블록(c *C.OUTDATABLOCK) *S수신_데이터_블록 {
	defer C.free(unsafe.Pointer(c))

	s := new(S수신_데이터_블록)
	s.식별번호 = int(c.TrIdNo)
	s.수신_데이터 = *(new수신_데이터(c.DataStruct))

	return s
}

type S수신_데이터 struct {
	블록_이름 string
	데이터   interface{}
	길이    int
}

func (this S수신_데이터) G블록_이름() string    { return this.블록_이름 }
func (this S수신_데이터) G데이터() interface{} { return this.데이터 }
func (this S수신_데이터) G길이() int          { return this.길이 }

func new수신_데이터(c *C.RECEIVED) *S수신_데이터 {
	블록_이름 := C.GoString(c.BlockName)
	데이터 := f_Go구조체로_변환(c)

	s := S수신_데이터{블록_이름: 블록_이름, 데이터: 데이터, 길이: int(c.Length)}

	return &s
}

//----------------------------------------------------------------------//
// 주식 현재가 조회 (c1101)
//----------------------------------------------------------------------//
func NewC1101InBlock(종목_코드 string) *C.char {
	언어_구분 := C.CString("k")
	공백_문자 := C.CString(" ")

	defer func() {
		C.free(unsafe.Pointer(언어_구분))
		C.free(unsafe.Pointer(공백_문자))
	}()

	c := new(C.Tc1101InBlock)
	c.Lang[0] = *언어_구분

	// 종목코드 검사
	switch {
	case 종목_코드 == "":
		에러 := 공용.F에러("종목 코드 내용이 비어있음.")
		panic(에러)
	case len(종목_코드) > len(c.Code):
		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함.")
		panic(에러)
	}

	for i := 0; i < len(c.Code); i++ {
		if i < len(종목_코드) {
			c.Code[i] = C.char(종목_코드[i])
		} else {
			c.Code[i] = *공백_문자
		}
	}

	return (*C.char)(unsafe.Pointer(c))
}

type S주식_현재가_조회_질의 struct {
	언어    string
	종목_코드 string
}

func (this S주식_현재가_조회_질의) G언어() string    { return this.언어 }
func (this S주식_현재가_조회_질의) G종목_코드() string { return this.종목_코드 }

func New주식_현재가_조회_질의(c *C.Tc1101InBlock) *S주식_현재가_조회_질의 {
	s := new(S주식_현재가_조회_질의)
	s.언어 = C.GoString(&(c.Lang[0]))
	s.종목_코드 = C.GoString(&(c.Code[0]))

	return s
}

type S주식_현재가_조회_기본_자료 struct {
	M종목_코드       string
	M종목명         string // 첫자리는 KOSPI200은 ‘*’, 스타지수종목은 ‘#’
	M현재가         int64
	M등락부호        uint8
	M등락폭         int64
	M등락율         float64
	//M매도_호가       int64
	//M매수_호가       int64
	M거래량         int64
	M전일대비_거래량_비율 float64 // 금일 거래량 / 전일 거래량
	M유동주_회전율     float64
	M거래대금_100만   int64
	M상한가         int64
	M고가          int64
	M시가          int64
	M시가대비_등락부호   uint8
	M시가대비_등락폭    int64
	M저가          int64
	M하한가         int64
	M시각          time.Time
	M매도_호가_모음 []int64
	M매도_잔량_모음 []int64
	M매수_호가_모음 []int64
	M매수_잔량_모음 []int64
	M시간외_매도_잔량   int64
	M시간외_매수_잔량   int64
	M피봇_2차_저항    int64  // 피봇가 + 전일 고가 – 전일 저가
	M피봇_1차_저항    int64  // (피봇가 * 2) – 전일 저가
	M피봇가         int64  // (전일 고가 + 전일 저가 + 전일 종가) / 3
	M피봇_1차_지지    int64  // (피봇가 * 2) – 전일 고가
	M피봇_2차_지지    int64  // 피봇가 – 전일고가 + 전일 저가
	M시장_구분       string // '코스피' , '코스닥'
	M업종명         string
	M자본금_규모      string
	M결산월         string
	M추가_정보_모음 []string
	M서킷_브레이커_구분 string
	M액면가         int64
	//M전일_종가_타이틀          string
	M전일_종가         int64
	M대용가          int64 // 담보가치인 듯
	M공모가          int64
	M5일_고가        int64
	M5일_저가        int64
	M20일_고가       int64
	M20일_저가       int64
	M52주_고가       int64 // 거래일만 따지면 1년이 52주이었던 것으로 기억함.
	M52주_고가_일자    time.Time
	M52주_저가       int64
	M52주_저가_일자    time.Time
	M유동_주식수_1000주 int64
	//M상장_주식수_1000주 int64
	M시가_총액_억       int64
	M거래원_정보_수신_시각  time.Time
	M매도_거래원_모음 []string
	M매도_거래량_모음 []int64
	M매수_거래원_모음 []string
	M매수_거래량_모음 []int64
	M매도_잔량_총합 int64
	M매수_잔량_총합 int64
	M외국인_매도_거래량    int64
	M외국인_매수_거래량    int64
	M외국인_시간        time.Time
	M외국인_지분율       float64
	M신용잔고_기준_결제일 time.Time
	M신용잔고율         float64
	//M유상_기준일     time.Time
	//M무상_기준일     time.Time
	M유상_배정_비율      float64
	M외국인_순매수량    int64
	M무상_배정_비율      float64
	M당일_자사주_신청_여부  bool
	M상장일           time.Time
	M대주주_지분율       float64
	M대주주_지분율_정보_일자 time.Time
	M네잎클로버_종목_여부   bool
	M증거금_비율        float64
	M자본금           int64
	M전체_거래원_매도_합계  int64
	M전체_거래원_매수_합계  int64
	//M종목명2              string
	M우회_상장_여부 bool
	//M유동주_회전율_2         float64
	//M코스피_구분_2          string    // 앞에 나온 '코스피/코스닥 구분'과 중복 아닌가?
	M공여율_기준일    time.Time // 공여율은 '신용거래 관련 비율'이라고 함.
	M공여율        float64   // 공여율(%)
	PER         float64   // PER
	M종목별_신용한도   float64
	M가중_평균_가격   int64
	M상장_주식수     int64
	M추가_상장_주식수  int64
	M종목_코멘트     string
	M전일_거래량     int64
	M전일_등락부호    uint8
	M전일_등락폭     int64
	M연중_최고가     int64
	M연중_최고가_일자  time.Time // 연중 최고가일
	M연중_최저가     int64
	M연중_최저가_일자  time.Time // 연중 최저가일
	M외국인_보유_주식수 int64     // 외국인 보유 주식수
	M외국인_지분_한도  float64   // % 단위
	M매매_수량_단위   int64     // 매매 수량 단위
	M대량_매매_방향   uint8     // 0 = 해당없음 1 = 매도 2 = 매수
	M대량_매매_존재   bool
}

func New주식_현재가_조회_기본_자료(c *C.Tc1101OutBlock) (데이터 *S주식_현재가_조회_기본_자료) {
	defer func() {
		if 값 := recover(); 값 != nil {
			공용.F에러(공용.F포맷된_문자열("%v", 값))
			데이터 = nil
		}
	}()
	
	g := (*Tc1101OutBlock)(unsafe.Pointer(c))
	개장일, 에러 := 공용.F한국증시_최근_개장일()
	공용.F에러_패닉(에러)
	
	지금 := time.Now()
	금일_24시 := time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
		23, 59, 59, 0, 지금.Location())

	s := new(S주식_현재가_조회_기본_자료)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	s.M종목명 = 공용.F2문자열_CP949(g.Title)
	s.M현재가,_ = 공용.F2정수64(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭,_ = 공용.F2정수64(g.Diff)
	s.M등락율 = f2실수_소숫점_추가(g.DiffRate, 2)
	//s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	//s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M거래량,_ = 공용.F2정수64(g.Volume)
	s.M전일대비_거래량_비율 = f2실수_소숫점_추가(g.TrVolRate, 2) // 당일 거래량 / 전일 거래량
	s.M유동주_회전율 = f2실수_소숫점_추가(g.FloatRate, 2)
	s.M거래대금_100만,_ = 공용.F2정수64(g.TrAmount)
	s.M상한가,_ = 공용.F2정수64(g.UpLmtPrice)
	s.M고가,_ = 공용.F2정수64(g.High)
	s.M시가,_ = 공용.F2정수64(g.Open)
	s.M시가대비_등락부호 = f2등락부호(g.VsOpenSign)
	s.M시가대비_등락폭,_ = 공용.F2정수64(g.VsOpenDiff)
	s.M저가,_ = 공용.F2정수64(g.Low)
	s.M하한가,_ = 공용.F2정수64(g.LowLmtPrice)

	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	매도_호가_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.OfferPrice1,
			g.OfferPrice2,
			g.OfferPrice3,
			g.OfferPrice4,
			g.OfferPrice5,
			g.OfferPrice6,
			g.OfferPrice7,
			g.OfferPrice8,
			g.OfferPrice9,
			g.OfferPrice10})
	
	매도_잔량_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.OfferVolume1,
			g.OfferVolume2,
			g.OfferVolume3,
			g.OfferVolume4,
			g.OfferVolume5,
			g.OfferVolume6,
			g.OfferVolume7,
			g.OfferVolume8,
			g.OfferVolume9,
			g.OfferVolume10})
		
	s.M매도_호가_모음 = make([]int64, 0)
	s.M매도_잔량_모음 = make([]int64, 0)
	
	for i, 매도_호가 := range 매도_호가_모음 {
		if 매도_호가 >= s.M현재가 && 매도_잔량_모음[i] > 0 {
			s.M매도_호가_모음 = append(s.M매도_호가_모음, 매도_호가)
			s.M매도_잔량_모음 = append(s.M매도_잔량_모음, 매도_잔량_모음[i])
		}
	}
	
	매수_호가_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.BidPrice1,
			g.BidPrice2,
			g.BidPrice3,
			g.BidPrice4,
			g.BidPrice5,
			g.BidPrice6,
			g.BidPrice7,
			g.BidPrice8,
			g.BidPrice9,
			g.BidPrice10})
	
	매수_잔량_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.BidVolume1,
			g.BidVolume2,
			g.BidVolume3,
			g.BidVolume4,
			g.BidVolume5,
			g.BidVolume6,
			g.BidVolume7,
			g.BidVolume8,
			g.BidVolume9,
			g.BidVolume10})
	
	s.M매수_호가_모음 = make([]int64, 0)
	s.M매수_잔량_모음 = make([]int64, 0)
	
	for i, 매수_호가 := range 매수_호가_모음 {
		if 매수_호가 <= s.M현재가 && 매수_잔량_모음[i] > 0 {
			s.M매수_호가_모음 = append(s.M매수_호가_모음, 매수_호가)
			s.M매수_잔량_모음 = append(s.M매수_잔량_모음, 매수_잔량_모음[i])
		}
	}
	
	s.M매도_잔량_총합,_ = 공용.F2정수64(g.OfferVolTot)
	s.M매수_잔량_총합,_ = 공용.F2정수64(g.BidVolTot)
	s.M시간외_매도_잔량,_ = 공용.F2정수64(g.OfferVolAfterHour)
	s.M시간외_매수_잔량,_ = 공용.F2정수64(g.BidVolAfterHour)
	s.M피봇_2차_저항,_ = 공용.F2정수64(g.PivotUp2)
	s.M피봇_1차_저항,_ = 공용.F2정수64(g.PivotUp1)
	s.M피봇가,_ = 공용.F2정수64(g.PivotPrice)
	s.M피봇_1차_지지,_ = 공용.F2정수64(g.PivotDown1)
	s.M피봇_2차_지지,_ = 공용.F2정수64(g.PivotDown2)
	s.M시장_구분 = 공용.F2문자열_CP949(g.Market)
	s.M업종명 = 공용.F2문자열_CP949(g.Sector)
	s.M자본금_규모 = 공용.F2문자열_CP949(g.CapSize)
	s.M결산월 = 공용.F2문자열_CP949(g.SettleMonth)
	
	s.M추가_정보_모음 = []string{
		공용.F2문자열_CP949(g.MarketAction1),
		공용.F2문자열_CP949(g.MarketAction2),
		공용.F2문자열_CP949(g.MarketAction3),
		공용.F2문자열_CP949(g.MarketAction4),
		공용.F2문자열_CP949(g.MarketAction5),
		공용.F2문자열_CP949(g.MarketAction6)}
	
	s.M서킷_브레이커_구분 = strings.TrimSpace(공용.F2문자열_CP949(g.CircuitBreaker))
	s.M액면가,_ = 공용.F2정수64(g.NominalPrice)
	//s.M전일_종가_타이틀 = 공용.F2문자열_CP949(g.PrevPriceTitle)
	s.M전일_종가,_ = 공용.F2정수64(g.PrevPrice)
	s.M대용가,_ = 공용.F2정수64(g.MortgageValue)
	s.M공모가,_ = 공용.F2정수64(g.PublicOfferPrice)
	s.M5일_고가,_ = 공용.F2정수64(g.High5Day)
	s.M5일_저가,_ = 공용.F2정수64(g.Low5Day)
	s.M20일_고가,_ = 공용.F2정수64(g.High20Day)
	s.M20일_저가,_ = 공용.F2정수64(g.Low20Day)
	s.M52주_고가,_ = 공용.F2정수64(g.High1Year)

	시각, 에러 = 공용.F2포맷된_시각("0102", g.High1YearDate)
	if 에러 == nil {
		s.M52주_고가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}
	
	if s.M52주_고가_일자.After(금일_24시) {
		s.M52주_고가_일자 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}

	s.M52주_저가,_ = 공용.F2정수64(g.Low1Year)

	시각, 에러 = 공용.F2포맷된_시각("0102", g.Low1YearDate)
	if 에러 == nil {
		s.M52주_저가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 시각.Location())
	}
	
	if s.M52주_저가_일자.After(금일_24시) {
		s.M52주_저가_일자 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}

	s.M유동_주식수_1000주,_ = 공용.F2정수64(g.FloatVolume)
	//s.M상장_주식수_1000주,_ = 공용.F2정수64(g.ListVolBy1000)
	s.M시가_총액_억,_ = 공용.F2정수64(g.MarketCapital)

	시각, 에러 = 공용.F2포맷된_시각("15:04", g.TraderInfoTime)
	if 에러 == nil {
		// 가장 최근 영업일을 알아내어야 함.
		s.M거래원_정보_수신_시각 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 0, 0, 지금.Location())
	}
	
	s.M매도_거래원_모음 = []string{
		공용.F2문자열_CP949(g.Seller1),
		공용.F2문자열_CP949(g.Seller2),
		공용.F2문자열_CP949(g.Seller3),
		공용.F2문자열_CP949(g.Seller4),
		공용.F2문자열_CP949(g.Seller5)}
	
	s.M매도_거래량_모음, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.Seller1Volume,
			g.Seller2Volume,
			g.Seller3Volume,
			g.Seller4Volume,
			g.Seller5Volume})
	
	s.M매수_거래원_모음 = []string{
		공용.F2문자열_CP949(g.Buyer1),
		공용.F2문자열_CP949(g.Buyer2),
		공용.F2문자열_CP949(g.Buyer3),
		공용.F2문자열_CP949(g.Buyer4),
		공용.F2문자열_CP949(g.Buyer5)}
	
	s.M매수_거래량_모음, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.Buyer1Volume,
			g.Buyer2Volume,
			g.Buyer3Volume,
			g.Buyer4Volume,
			g.Buyer5Volume})

	if strings.TrimSpace(공용.F2문자열(g.ForeignSellVolume)) != "" {
		s.M외국인_매도_거래량,_ = 공용.F2정수64(g.ForeignSellVolume)
	} else {
		s.M외국인_매도_거래량 = 0
	}
	
	if strings.TrimSpace(공용.F2문자열(g.ForeignBuyVolume)) != "" {
		s.M외국인_매수_거래량,_ = 공용.F2정수64(g.ForeignBuyVolume)
	} else {
		s.M외국인_매수_거래량 = 0
	}

	시각, 에러 = 공용.F2포맷된_시각("15:04", g.ForeignTime)
	if 에러 == nil {
		s.M외국인_시간 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 0, 0, 지금.Location())
	}
		
	s.M외국인_지분율 = f2실수_소숫점_추가(g.ForeignHoldingRate, 2)

	시각, 에러 = 공용.F2포맷된_시각("0102", g.SettleDate)
	if 에러 == nil {
		s.M신용잔고_기준_결제일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}
	
	if s.M신용잔고_기준_결제일.After(금일_24시) {
		s.M신용잔고_기준_결제일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}
	
	s.M신용잔고율 = f2실수_소숫점_추가(g.DebtPercent, 2)

	/* if strings.TrimSpace(공용.F2문자열(g.RightsIssueDate)) != "" {
		공용.F문자열_출력("유상 기준일 : '%v'", 공용.F2문자열(g.RightsIssueDate))
	}
	
	시각, 에러 = 공용.F2포맷된_시각("0102", g.RightsIssueDate)
	if 에러 == nil {
		s.M유상_기준일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 지금.Location())
	}
	
	if s.M유상_기준일.After(금일_24시) {
		s.M유상_기준일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 지금.Location())
	}
	
	if strings.TrimSpace(공용.F2문자열(g.BonusIssueDate)) != "" {
		공용.F문자열_출력("무상 기준일 : '%v'", 공용.F2문자열(g.BonusIssueDate))
	}
	
	시각, 에러 = 공용.F2포맷된_시각("0102", g.BonusIssueDate)
	if 에러 == nil {
		s.M무상_기준일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 지금.Location())
	}
	
	if s.M무상_기준일.After(금일_24시) {
		s.M무상_기준일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 지금.Location())
	} */

	s.M유상_배정_비율 = f2실수_소숫점_추가(g.RightsIssueRate, 2)
	s.M무상_배정_비율 = f2실수_소숫점_추가(g.BonusIssueRate, 2)
	s.M외국인_순매수량,_ = 공용.F2정수64(g.ForeignFloatVol)
	s.M당일_자사주_신청_여부 = 공용.F2참거짓(g.TreasuryStock, "1", true)
	s.M상장일 , _ = 공용.F2포맷된_시각("20060102", g.IpoDate)
	s.M대주주_지분율 = f2실수_소숫점_추가(g.MajorHoldRate, 2)
	s.M대주주_지분율_정보_일자 , _ = 공용.F2포맷된_시각("060102", g.MajorHoldInfoDate)
	s.M네잎클로버_종목_여부 = 공용.F2참거짓(g.FourLeafClover, "1", true)

	증거금_비율_위치 := 공용.F2문자열(g.MarginRate)
	증거금_비율_문자열 := ""
	
	switch 증거금_비율_위치 {
	case "1", "2", "3", "4", "5", "6":
		위치 := 공용.F2정수(증거금_비율_위치)
		증거금_비율_문자열 = s.M추가_정보_모음[위치-1] 
	default:
		에러 := 공용.F에러("예상하지 못한 증거금_비율_위치. %v", 증거금_비율_위치)
		panic(에러)
	}

	증거금_비율_문자열 = strings.Replace(증거금_비율_문자열, "증거금", "", -1)
	증거금_비율_문자열 = strings.Replace(증거금_비율_문자열, "%", "", -1)
	증거금_비율_문자열 = strings.TrimSpace(증거금_비율_문자열)
	s.M증거금_비율 = 공용.F2실수(증거금_비율_문자열)

	s.M자본금,_ = 공용.F2정수64(g.Capital)
	s.M전체_거래원_매도_합계,_ = 공용.F2정수64(g.SellTotalSum)
	s.M전체_거래원_매수_합계,_ = 공용.F2정수64(g.BuyTotalSum)
	//s.M종목명2 = 공용.F2문자열_CP949(g.Title2)
	s.M우회_상장_여부 = 공용.F2참거짓(g.BackdoorListing, "1", true)
	//s.M유동주_회전율_2 = f2실수_소숫점_추가(g.FloatRate2, 2)
	//s.M코스피_구분_2 = 공용.F2문자열_CP949(g.Market2)

	시각, 에러 = 공용.F2포맷된_시각("0102", g.DebtTrDate)
	if 에러 == nil {
		s.M공여율_기준일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())	
	}
	
	if s.M공여율_기준일.After(금일_24시) {
		s.M공여율_기준일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}
	
	s.M공여율 = f2실수_소숫점_추가(g.DebtTrPercent, 2)
	s.PER = f2실수_소숫점_추가(g.PER, 2)

	종목별_신용한도_위치 := 공용.F2문자열(g.DebtLimit)
	종목별_신용한도_문자열 := ""

	switch 종목별_신용한도_위치 {
	case "0":
		종목별_신용한도_문자열 = ""
	case "1", "2", "3", "4", "5", "6":
		위치 := 공용.F2정수(종목별_신용한도_위치)
		종목별_신용한도_문자열 = s.M추가_정보_모음[위치-1]
	default:
		에러 := 공용.F에러("시각 %v, 종목코드 %v\n"+
			"예상하지 못한 종목별_신용한도_위치. %v\n"+
			" 1 : %v\n2 : %v\n3 : %v\n"+
			" 4 : %v\n5 : %v\n6 : %v\n",
			s.M시각, s.M종목_코드, 종목별_신용한도_위치,
			s.M추가_정보_모음[0], s.M추가_정보_모음[1], s.M추가_정보_모음[2],
			s.M추가_정보_모음[3], s.M추가_정보_모음[4], s.M추가_정보_모음[5])
		panic(에러)
	}

	switch {
	case 종목별_신용한도_문자열 == "":
		s.M종목별_신용한도 = 0
	case strings.Contains(종목별_신용한도_문자열, "없음"):
		s.M종목별_신용한도 = 100
	default:
		종목별_신용한도_문자열 = strings.Replace(종목별_신용한도_문자열, "신용", "", -1)
		종목별_신용한도_문자열 = strings.Replace(종목별_신용한도_문자열, "한도", "", -1)
		종목별_신용한도_문자열 = strings.Replace(종목별_신용한도_문자열, "제한", "", -1)
		종목별_신용한도_문자열 = strings.Replace(종목별_신용한도_문자열, "%", "", -1)
		종목별_신용한도_문자열 = strings.TrimSpace(종목별_신용한도_문자열)

		s.M종목별_신용한도 = 공용.F2실수(종목별_신용한도_문자열)
	}

	s.M가중_평균_가격,_ = 공용.F2정수64(g.WeightAvgPrice)
	s.M상장_주식수,_ = 공용.F2정수64(g.ListedVolume)
	s.M추가_상장_주식수,_ = 공용.F2정수64(g.AddListing)
	s.M종목_코멘트 = 공용.F2문자열_CP949(g.Comment)
	s.M전일_거래량,_ = 공용.F2정수64(g.PrevVolume)
	s.M전일_등락부호 = f2등락부호(g.VsPrevSign)
	s.M전일_등락폭,_ = 공용.F2정수64(g.VsPrevDiff)
	s.M연중_최고가,_ = 공용.F2정수64(g.High1Year2)

	시각, 에러 = 공용.F2포맷된_시각("0102", g.High1YearDate2)
	if 에러 == nil {
		s.M연중_최고가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}

	s.M연중_최저가,_ = 공용.F2정수64(g.Low1Year2)

	시각, 에러 = 공용.F2포맷된_시각("0102", g.Low1YearDate2)
	if 에러 == nil {
		s.M연중_최저가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(), 
			0, 0, 0, 0, 지금.Location())
	}
	
	s.M외국인_보유_주식수,_ = 공용.F2정수64(g.ForeignHoldQty)
	s.M외국인_지분_한도 = f2실수_소숫점_추가(g.ForeignLmtPercent, 2)
	s.M매매_수량_단위,_ = 공용.F2정수64(g.TrUnitVolume)
	s.M대량_매매_방향 = uint8(g.DarkPoolOfferBid[0]) // 0 = 해당없음 1 = 매도 2 = 매수
	s.M대량_매매_존재 = 공용.F2참거짓(g.DarkPoolExist, "1", true)

	return s
}

type S주식_현재가_조회_변동_거래량_자료 struct { // 변동거래량자료[반복]
	M시각     time.Time
	M현재가    int64
	M등락부호   uint8
	M등락폭    int64
	M매도_호가  int64
	M매수_호가  int64
	M변동_거래량 int64
	M거래량    int64
}

func New주식_현재가_조회_변동_거래량_자료(c *C.Tc1101OutBlock2) (데이터 *S주식_현재가_조회_변동_거래량_자료) {
	defer func() {
		if 값 := recover(); 값 != nil {
			공용.F에러(공용.F포맷된_문자열("%v", 값))
			데이터 = nil
		}
	}()

	g := (*Tc1101OutBlock2)(unsafe.Pointer(c))
	s := new(S주식_현재가_조회_변동_거래량_자료)
	
	지금 := time.Now()
	개장일, 에러 := 공용.F한국증시_최근_개장일()
	공용.F에러_패닉(에러)
	
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M현재가,_ = 공용.F2정수64(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭,_ = 공용.F2정수64(g.Diff)
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M변동_거래량,_ = 공용.F2정수64(g.DiffVolume)
	s.M거래량,_ = 공용.F2정수64(g.Volume)
	
	if s.M매도_호가 < s.M현재가 && s.M매도_호가 > 0 {
		s.M매도_호가 = s.M현재가
	}
	
	if s.M매수_호가 > s.M현재가 && s.M현재가 > 0 {
		s.M매수_호가 = s.M현재가
	}

	return s
}

type S주식_현재가_조회_동시호가 struct { // 종목지표
	M동시호가_구분 uint8 // 0:동시호가 아님 1:동시호가 2:동시호가연장 3:시가범위연장 4:종가범위연장 5:배분개시 6:변동성 완화장치 발동
	M예상_체결가   int64
	M예상_체결_부호  uint8
	M예상_등락폭   int64
	M예상_등락율   float64
	M예상_체결_수량  int64
	//ECN정보_유무        bool // 우리나라에는 아직 ECN이 없는 것으로 알고 있음.
	//ECN전일_종가        int64
	//ECN등락부호        uint8
	//ECN등락폭          int64
	//ECN등락율          float64
	//ECN체결_수량        int64
	//ECN대비_예상_체결_부호  uint8
	//ECN대비_예상_체결_등락폭 int64
	//ECN대비_예상_체결_등락율 float64
}

func New주식_현재가_조회_동시호가(c *C.Tc1101OutBlock3) *S주식_현재가_조회_동시호가 {
	g := (*Tc1101OutBlock3)(unsafe.Pointer(c))

	s := new(S주식_현재가_조회_동시호가)
	s.M동시호가_구분 = uint8(공용.F2정수(g.SyncOfferBid))
	s.M예상_체결가,_ = 공용.F2정수64(g.EstmPrice)
	s.M예상_체결_부호 = f2등락부호(g.EstmSign)
	s.M예상_등락폭,_ = 공용.F2정수64(g.EstmDiff)
	s.M예상_등락율 = f2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결_수량,_ = 공용.F2정수64(g.EstmVolume)

	return s
}

//----------------------------------------------------------------------//
// ETF 현재가 조회 (c1151)
//----------------------------------------------------------------------//
func NewC1151InBlock(종목_코드 string) *C.char {
	언어_구분 := C.CString("k")
	공백_문자 := C.CString(" ")

	defer func() {
		C.free(unsafe.Pointer(언어_구분))
		C.free(unsafe.Pointer(공백_문자))
	}()

	c := new(C.Tc1151InBlock)
	c.Lang[0] = *언어_구분

	// 종목코드 검사
	switch {
	case 종목_코드 == "":
		에러 := 공용.F에러("종목 코드 내용이 비어있음.")
		panic(에러)
	case len(종목_코드) > len(c.Code):
		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함.")
		panic(에러)
	}

	for i := 0; i < len(c.Code); i++ {
		if i < len(종목_코드) {
			c.Code[i] = C.char(종목_코드[i])
		} else {
			c.Code[i] = *공백_문자
		}
	}

	return (*C.char)(unsafe.Pointer(c))
}

type S_ETF_현재가_조회_기본_자료 struct { // 종목마스타기본자료
	M종목_코드      string
	M종목명        string
	M현재가        int64
	M등락부호       uint8
	M등락폭        int64
	M등락율        float64
	M거래량        int64
	M전일대비_거래량_비율 float64
	M유동주_회전율    float64
	M거래대금_100만  int64
	M상한가        int64
	M고가         int64
	M시가         int64
	M시가대비_등락부호  uint8
	M시가대비_등락폭   int64
	M저가         int64
	M하한가        int64
	M시각         time.Time
	M매도_호가_모음 []int64
	M매도_잔량_모음 []int64
	M매수_호가_모음 []int64
	M매수_잔량_모음 []int64
	M시간외_매도_잔량  int64
	M시간외_매수_잔량  int64
	M피봇_2차_저항   int64
	M피봇_1차_저항   int64
	M피봇_가격      int64
	M피봇_1차_지지   int64
	M피봇_2차_지지   int64
	M시장_구분      string
	M업종명        string
	M자본금_규모     string
	M결산월        string
	M추가_정보_모음 []string
	M서킷_브레이커_구분 string
	M액면가        int64
	//M전일_종가_타이틀         string // GUI화면에 쓸 자료이니 별 필요 없을 듯...
	M전일_종가        int64
	M대용가          int64
	M공모가          int64
	M5일_고가        int64
	M5일_저가        int64
	M20일_고가       int64
	M20일_저가       int64
	M52주_고가       int64
	M52주_저가       int64
	M유동_주식수_1000주 int64
	//M상장_주식수_1000주 int64
	M시가_총액_억     int64
	M거래원_정보_수신_시각 time.Time
	M매도_거래원_모음 []string
	M매도_거래량_모음 []int64
	M매수_거래원_모음 []string
	M매수_거래량_모음 []int64
	M외국인_매도_거래량   int64
	M외국인_매수_거래량   int64
	M외국인_시간       time.Time
	M외국인_지분율      float64
	M신용잔고율        float64
	M유상_기준일    time.Time
	M무상_기준일    time.Time
	M유상_배정_비율     float64
	M무상_배정_비율     float64
	M상장일          time.Time
	M상장_주식수       int64
	M전체_거래원_매도_합계 int64
	M전체_거래원_매수_합계 int64
}

func New_ETF_현재가_조회_기본_자료(c *C.char) *S_ETF_현재가_조회_기본_자료 {
	개장일, 에러 := 공용.F한국증시_최근_개장일()
	공용.F에러_패닉(에러)
	
	지금 := time.Now()
	
	g := (*Tc1151OutBlock)(unsafe.Pointer(c))
	s := new(S_ETF_현재가_조회_기본_자료)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	s.M종목명 = 공용.F2문자열_CP949(g.Title)
	s.M현재가,_ = 공용.F2정수64(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭,_ = 공용.F2정수64(g.Diff)
	s.M등락율 = f2실수_소숫점_추가(g.DiffRate, 2)
	//s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	//s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M거래량,_ = 공용.F2정수64(g.Volume)
	s.M전일대비_거래량_비율 = f2실수_소숫점_추가(g.TrVolRate, 2)
	s.M유동주_회전율 = f2실수_소숫점_추가(g.FloatVolRate, 2)
	s.M거래대금_100만,_ = 공용.F2정수64(g.TrAmount)
	s.M상한가,_ = 공용.F2정수64(g.UpLmtPrice)
	s.M고가,_ = 공용.F2정수64(g.High)
	s.M시가,_ = 공용.F2정수64(g.Open)
	s.M시가대비_등락부호 = f2등락부호(g.VsOpenSign)
	s.M시가대비_등락폭,_ = 공용.F2정수64(g.VsOpenDiff)
	s.M저가,_ = 공용.F2정수64(g.Low)
	s.M하한가,_ = 공용.F2정수64(g.LowLmtPrice)

	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	매도_호가_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.OfferPrice1,
			g.OfferPrice2,
			g.OfferPrice3,
			g.OfferPrice4,
			g.OfferPrice5,
			g.OfferPrice6,
			g.OfferPrice7,
			g.OfferPrice8,
			g.OfferPrice9,
			g.OfferPrice10})
		
	매도_잔량_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.OfferVolume1,
			g.OfferVolume2,
			g.OfferVolume3,
			g.OfferVolume4,
			g.OfferVolume5,
			g.OfferVolume6,
			g.OfferVolume7,
			g.OfferVolume8,
			g.OfferVolume9,
			g.OfferVolume10})
		
	s.M매도_호가_모음 = make([]int64, 0)
	s.M매도_잔량_모음 = make([]int64, 0)

	for i, 매도_호가 := range 매도_호가_모음 {		
		if 매도_호가 >= s.M현재가 && 매도_잔량_모음[i] > 0 {
			s.M매도_호가_모음 = append(s.M매도_호가_모음, 매도_호가)
			s.M매도_잔량_모음 = append(s.M매도_잔량_모음, 매도_잔량_모음[i])
		}
	}
	
	매수_호가_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.BidPrice1,
			g.BidPrice2,
			g.BidPrice3,
			g.BidPrice4,
			g.BidPrice5,
			g.BidPrice6,
			g.BidPrice7,
			g.BidPrice8,
			g.BidPrice9,
			g.BidPrice10})
	
	매수_잔량_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.BidVolume1,
			g.BidVolume2,
			g.BidVolume3,
			g.BidVolume4,
			g.BidVolume5,
			g.BidVolume6,
			g.BidVolume7,
			g.BidVolume8,
			g.BidVolume9,
			g.BidVolume10})
	
	s.M매수_호가_모음 = make([]int64, 0)
	s.M매수_잔량_모음 = make([]int64, 0)
		
	for i, 매수_호가 := range 매수_호가_모음 {
		if 매수_호가 >= s.M현재가 && 매수_잔량_모음[i] > 0 {
			s.M매수_호가_모음 = append(s.M매수_호가_모음, 매수_호가)
			s.M매수_잔량_모음 = append(s.M매수_잔량_모음, 매수_잔량_모음[i])
		}
	}
		
//	s.M매도_잔량_총합,_ = 공용.F2정수64(g.OfferVolTot)
//	s.M매수_잔량_총합,_ = 공용.F2정수64(g.BidVolTot)
	s.M시간외_매도_잔량,_ = 공용.F2정수64(g.OfferVolAfterHour)
	s.M시간외_매수_잔량,_ = 공용.F2정수64(g.BidVolAfterHour)
	s.M피봇_2차_저항,_ = 공용.F2정수64(g.PivotUp2)
	s.M피봇_1차_저항,_ = 공용.F2정수64(g.PivotUp1)
	s.M피봇_가격,_ = 공용.F2정수64(g.PivotPrice)
	s.M피봇_1차_지지,_ = 공용.F2정수64(g.PivotDown1)
	s.M피봇_2차_지지,_ = 공용.F2정수64(g.PivotDown2)
	s.M시장_구분 = 공용.F2문자열_CP949(g.Market)
	s.M업종명 = 공용.F2문자열_CP949(g.Sector)
	s.M자본금_규모 = strings.TrimSpace(공용.F2문자열_CP949(g.CapSize))
	s.M결산월 = 공용.F2문자열_CP949(g.SettleMonth)
	s.M추가_정보_모음 = []string{
		공용.F2문자열_CP949(g.MarketAction1),
		공용.F2문자열_CP949(g.MarketAction2),
		공용.F2문자열_CP949(g.MarketAction3),
		공용.F2문자열_CP949(g.MarketAction4),
		공용.F2문자열_CP949(g.MarketAction5),
		공용.F2문자열_CP949(g.MarketAction6)}	
	s.M서킷_브레이커_구분 = strings.TrimSpace(공용.F2문자열_CP949(g.CircuitBreaker))
	s.M액면가,_ = 공용.F2정수64(g.NominalPrice)
	//s.M전일_종가_타이틀 = 공용.F2문자열_CP949(g.PrevPriceTitle)
	s.M전일_종가,_ = 공용.F2정수64(g.PrevPrice)
	s.M대용가,_ = 공용.F2정수64(g.MortgageValue)
	s.M공모가,_ = 공용.F2정수64(g.PublicOfferPrice)
	s.M5일_고가,_ = 공용.F2정수64(g.High5Day)
	s.M5일_저가,_ = 공용.F2정수64(g.Low5Day)
	s.M20일_고가,_ = 공용.F2정수64(g.High20Day)
	s.M20일_저가,_ = 공용.F2정수64(g.Low20Day)
	s.M52주_고가,_ = 공용.F2정수64(g.High1Year)
	s.M52주_저가,_ = 공용.F2정수64(g.Low1Year)
	s.M유동_주식수_1000주,_ = 공용.F2정수64(g.FloatVolume)
	//s.M상장_주식수_1000주,_ = 공용.F2정수64(g.ListVolBy1000)
	s.M시가_총액_억,_ = 공용.F2정수64(g.MarketCapital)
	
	시각, 에러 = 공용.F2포맷된_시각("15:04", g.TraderInfoTime)
	if 에러 == nil {
		s.M거래원_정보_수신_시각 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 0, 0, 지금.Location())
	}
	
	s.M매도_거래원_모음 = []string{
		공용.F2문자열_CP949(g.Seller1),
		공용.F2문자열_CP949(g.Seller2),
		공용.F2문자열_CP949(g.Seller3),
		공용.F2문자열_CP949(g.Seller4),
		공용.F2문자열_CP949(g.Seller5)}
	
	s.M매도_거래량_모음, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.Seller1Volume,
			g.Seller2Volume,
			g.Seller3Volume,
			g.Seller4Volume,
			g.Seller5Volume})
	
	s.M매수_거래원_모음 = []string{
		공용.F2문자열_CP949(g.Buyer1),
		공용.F2문자열_CP949(g.Buyer2),
		공용.F2문자열_CP949(g.Buyer3),
		공용.F2문자열_CP949(g.Buyer4),
		공용.F2문자열_CP949(g.Buyer5)}

	s.M매수_거래량_모음, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.Buyer1Volume,
			g.Buyer2Volume,
			g.Buyer3Volume,
			g.Buyer4Volume,
			g.Buyer5Volume})
		
	s.M외국인_매도_거래량,_ = 공용.F2정수64(g.ForeignSellVolume)
	s.M외국인_매수_거래량,_ = 공용.F2정수64(g.ForeignBuyVolume)
	
	시각, 에러 = 공용.F2포맷된_시각("15:04", g.ForeignTime)
	if 에러 == nil {
		s.M외국인_시간 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(), 
			시각.Hour(), 시각.Minute(), 0, 0, 지금.Location())
	}
	
	s.M외국인_지분율 = f2실수_소숫점_추가(g.ForeignHoldingRate, 2)
	s.M신용잔고율 = f2실수_소숫점_추가(g.DebtPercent, 2)
	
	/* 공용.F메모("유상 기준일에 왜 연도가 없는 거지? 매년 하는 것이란 말인가? 그럼 유상배정일이 아니잖아!!!")
	if strings.TrimSpace(공용.F2문자열(g.RightsIssueDate)) != "" {
		공용.F문자열_출력("유상 기준일 : '%v'", 공용.F2문자열(g.RightsIssueDate))
	}
	
	일자, 에러 := 공용.F2포맷된_시각("0102", g.RightsIssueDate)
	if 에러 == nil {
		s.M유상_기준일 = time.Date(지금.Year(), 일자.Month(), 일자.Day(), 
			0, 0, 0, 0, 지금.Location()) 
	}
	
	if strings.TrimSpace(공용.F2문자열(g.BonusIssueDate)) != "" {
		공용.F문자열_출력("무상 기준일 : '%v'", 공용.F2문자열(g.BonusIssueDate))
	}
	
	일자, 에러 = 공용.F2포맷된_시각("0102", g.BonusIssueDate)
	if 에러 == nil {
		s.M무상_기준일 = time.Date(지금.Year(), 일자.Month(), 일자.Day(),
			0, 0, 0, 0, 지금.Location())
	} */
	
	s.M유상_배정_비율 = f2실수_소숫점_추가(g.RightsIssueRate, 2)
	s.M무상_배정_비율 = f2실수_소숫점_추가(g.BonusIssueRate, 2)
	s.M상장일, _ = 공용.F2포맷된_시각("20060102", g.IpoDate)
	s.M상장_주식수,_ = 공용.F2정수64(g.ListedVolume)
	s.M전체_거래원_매도_합계,_ = 공용.F2정수64(g.SellTotalSum)
	s.M전체_거래원_매수_합계,_ = 공용.F2정수64(g.BuyTotalSum)

	return s
}

type S_ETF_현재가_조회_변동_거래_자료 struct {
	M시각     time.Time
	M현재가    int64
	M등락부호   uint8
	M등락폭    int64
	M매도_호가  int64
	M매수_호가  int64
	M변동_거래량 int64
	M거래량    int64
}

func New_ETF_현재가_조회_변동_거래_자료(c *C.Tc1151OutBlock2) (
	데이터 *S_ETF_현재가_조회_변동_거래_자료) {
	defer func() {
		if 값 := recover(); 값 != nil {
			공용.F에러(공용.F포맷된_문자열("%v", 값))
			데이터 = nil
		}
	}()
	
	지금 := time.Now()
	개장일, 에러 := 공용.F한국증시_최근_개장일()
	공용.F에러_패닉(에러)
	
	g := (*Tc1151OutBlock2)(unsafe.Pointer(c))
	s := new(S_ETF_현재가_조회_변동_거래_자료)
	
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(개장일.Year(), 개장일.Month(), 개장일.Day(),
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M현재가,_ = 공용.F2정수64(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭,_ = 공용.F2정수64(g.Diff)
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M변동_거래량,_ = 공용.F2정수64(g.DiffVolume)
	s.M거래량,_ = 공용.F2정수64(g.Volume)
	
	if s.M매도_호가 < s.M현재가 && s.M매도_호가 > 0 {
		s.M매도_호가 = s.M현재가
	}
	
	if s.M매수_호가 > s.M현재가 && s.M현재가 > 0 {
		s.M매수_호가 = s.M현재가
	}

	return s
}

type S_ETF_현재가_조회_동시호가 struct {
	M동시호가_구분 uint8
	M예상_체결가    int64
	M예상_체결_부호  uint8
	M예상_등락폭 int64
	M예상_등락율 float64
	M예상_체결_수량  int64
}

func New_ETF_현재가_조회_동시호가(c *C.char) *S_ETF_현재가_조회_동시호가 {
	g := (*Tc1151OutBlock3)(unsafe.Pointer(c))

	s := new(S_ETF_현재가_조회_동시호가)
	s.M동시호가_구분 = uint8(공용.F2정수(g.SyncOfferBid))
	s.M예상_체결가,_ = 공용.F2정수64(g.EstmPrice)
	s.M예상_체결_부호 = f2등락부호(g.EstmSign)
	s.M예상_등락폭,_ = 공용.F2정수64(g.EstmDiff)
	s.M예상_등락율 = f2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결_수량,_ = 공용.F2정수64(g.EstmVolume)

	return s
}

type S_ETF_현재가_조회_ETF자료 struct {
	ETF구분           string
	NAV             float64
	NAV등락부호        uint8
	NAV등락폭         float64
	M전일NAV          float64
	M괴리율            float64
	M괴리율_부호         uint8
	M설정단위_당_현금_배당액  float64
	M구성_종목수         int64
	M순자산_총액_억      int64
	M추적_오차율         float64
	LP_매도_잔량_모음 []int64
	LP_매수_잔량_모음 []int64
	ETF_복제_방법_구분_코드 string
	ETF_상품_유형_코드    string
}

func New_ETF_현재가_조회_ETF자료(c *C.char) *S_ETF_현재가_조회_ETF자료 {
	g := (*Tc1151OutBlock4)(unsafe.Pointer(c))

	s := new(S_ETF_현재가_조회_ETF자료)
	
	switch 공용.F2문자열(g.ETF) {
	case "4":
		s.ETF구분 = P코스닥
	case "8":
		s.ETF구분 = P코스피
	default:
		공용.F에러("예상치 못한 구분값. %v", 공용.F2문자열(g.ETF))
	}
	
	s.NAV = f2실수_소숫점_추가(g.NAV, 2)
	s.NAV등락부호 = f2등락부호(g.DiffSign)
	s.NAV등락폭 = f2실수_소숫점_추가(g.Diff, 2)
	s.M전일NAV = f2실수_소숫점_추가(g.PrevNAV, 2)
	s.M괴리율 = f2실수_소숫점_추가(g.DivergeRate, 2)
	s.M괴리율_부호 = f2등락부호(g.DivergeSign)
	s.M설정단위_당_현금_배당액 = 공용.F2실수(g.DividendPerCU)
	s.M구성_종목수,_ = 공용.F2정수64(g.ConstituentNo)
	s.M순자산_총액_억,_ = 공용.F2정수64(g.NAVBy100Million)
	s.M추적_오차율 = f2실수_소숫점_추가(g.TrackingErrRate, 2)
	
	매도_잔량_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.LP_OfferVolume1,
			g.LP_OfferVolume2,
			g.LP_OfferVolume3,
			g.LP_OfferVolume4,
			g.LP_OfferVolume5,
			g.LP_OfferVolume6,
			g.LP_OfferVolume7,
			g.LP_OfferVolume8,
			g.LP_OfferVolume9,
			g.LP_OfferVolume10})
	
	s.LP_매도_잔량_모음 = make([]int64, 0)
	
	for _, 매도_잔량 := range 매도_잔량_모음 {
		// 수량이 0인 주문은 없는 것이나 마찬가지이니 걸러낸다.
		if 매도_잔량 == 0 {
			continue
		}
		
		s.LP_매도_잔량_모음 = append(s.LP_매도_잔량_모음, 매도_잔량)
	}
	
	매수_잔량_모음, _ := 공용.F2정수64_모음(
		[]interface{}{
			g.LP_BidVolume1,
			g.LP_BidVolume2,
			g.LP_BidVolume3,
			g.LP_BidVolume4,
			g.LP_BidVolume5,
			g.LP_BidVolume6,
			g.LP_BidVolume7,
			g.LP_BidVolume8,
			g.LP_BidVolume9,
			g.LP_BidVolume10})
	
	s.LP_매수_잔량_모음 = make([]int64, 0)
	
	for _, 매수_잔량 := range 매수_잔량_모음 {
		// 수량이 0인 주문은 없는 것이나 마찬가지이니 걸러낸다.
		if 매수_잔량 == 0 {
			continue
		}
		
		s.LP_매수_잔량_모음 = append(s.LP_매수_잔량_모음, 매수_잔량)
	}
	
	s.ETF_복제_방법_구분_코드 = 공용.F2문자열_CP949(g.TrackingMethod)
	s.ETF_상품_유형_코드 = 공용.F2문자열_CP949(g.ETF_Type)

	return s
}

type S_ETF_현재가_조회_지수_자료 struct {
	M업종_코드       string
	KRX지수_코드       string
	M지수_이름       string
	M지수      float64
	M지수_등락부호 uint8
	M지수_등락폭  float64
	M채권_지수       float64
	M채권_지수_등락부호  uint8
	M채권_지수_등락폭   float64
	M해외_지수_코드    string
	M기타_업종_코드    string
	M채권_지수_코드    string
	M채권_지수_세부_코드 string
}

func New_ETF_현재가_조회_지수_자료(c *C.char) *S_ETF_현재가_조회_지수_자료 {
	g := (*Tc1151OutBlock5)(unsafe.Pointer(c))

	s := new(S_ETF_현재가_조회_지수_자료)
	s.M업종_코드 = strings.TrimSpace(공용.F2문자열(g.SectorCode))
	s.KRX지수_코드 = strings.TrimSpace(공용.F2문자열(g.IndexCode))
	s.M지수_이름 = strings.TrimSpace(공용.F2문자열_CP949(g.IndexName))
	s.M지수 = f2실수_소숫점_추가(g.KP200Index, 2)
	s.M지수_등락부호 = f2등락부호(g.KP200Sign)
	s.M지수_등락폭 = f2실수_소숫점_추가(g.KP200Diff, 2)
	s.M채권_지수 = f2실수_소숫점_추가(g.BondIndex, 2)
	s.M채권_지수_등락부호 = f2등락부호(g.BondSign)
	s.M채권_지수_등락폭 = f2실수_소숫점_추가(g.BondDiff, 2)
	s.M해외_지수_코드 = strings.TrimSpace(공용.F2문자열(g.ForeignIndexSymbol))
	s.M기타_업종_코드 = strings.TrimSpace(공용.F2문자열(g.EtcSectorCode))
	s.M채권_지수_코드 = strings.TrimSpace(공용.F2문자열(g.BondIndexCode))
	s.M채권_지수_세부_코드 = strings.TrimSpace(공용.F2문자열(g.BondDetailCode))

	return s
}

//----------------------------------------------------------------------//
// 코스피 호가 잔량 (h1)
//----------------------------------------------------------------------//
func New코스피_호가_잔량_질의(종목_코드 string) C.Th1InBlock {
	c := Th1InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.", len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Th1InBlock)(unsafe.Pointer(&c)))
}

type S코스피_호가_잔량 struct {
	M종목_코드         string
	M시각            time.Time
	M매도_호가         int64
	M매수_호가         int64
	M매도_호가_잔량      int64
	M매수_호가_잔량      int64
	M차선_매도_호가      int64
	M차선_매수_호가      int64
	M차선_매도_호가_잔량   int64
	M차선_매수_호가_잔량   int64
	M차차선_매도_호가     int64
	M차차선_매수_호가     int64
	M차차선_매도_호가_잔량  int64
	M차차선_매수_호가_잔량  int64
	M4차선_매도_호가     int64
	M4차선_매수_호가     int64
	M4차선_매도_호가_잔량  int64
	M4차선_매수_호가_잔량  int64
	M5차선_매도_호가     int64
	M5차선_매수_호가     int64
	M5차선_매도_호가_잔량  int64
	M5차선_매수_호가_잔량  int64
	M6차선_매도_호가     int64
	M6차선_매수_호가     int64
	M6차선_매도_호가_잔량  int64
	M6차선_매수_호가_잔량  int64
	M7차선_매도_호가     int64
	M7차선_매수_호가     int64
	M7차선_매도_호가_잔량  int64
	M7차선_매수_호가_잔량  int64
	M8차선_매도_호가     int64
	M8차선_매수_호가     int64
	M8차선_매도_호가_잔량  int64
	M8차선_매수_호가_잔량  int64
	M9차선_매도_호가     int64
	M9차선_매수_호가     int64
	M9차선_매도_호가_잔량  int64
	M9차선_매수_호가_잔량  int64
	M10차선_매도_호가    int64
	M10차선_매수_호가    int64
	M10차선_매도_호가_잔량 int64
	M10차선_매수_호가_잔량 int64
	M누적_거래량        int64
}

func New코스피_호가_잔량(c *C.char) *S코스피_호가_잔량 {
	g := (*Th1OutBlock)(unsafe.Pointer(c))
	
	지금 := time.Now()

	s := new(S코스피_호가_잔량)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
		
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice1)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice1)
	s.M매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume1)
	s.M매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume1)
	s.M차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice2)
	s.M차선_매수_호가,_ = 공용.F2정수64(g.BidPrice2)
	s.M차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume2)
	s.M차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume2)
	s.M차차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice3)
	s.M차차선_매수_호가,_ = 공용.F2정수64(g.BidPrice3)
	s.M차차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume3)
	s.M차차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume3)
	s.M4차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice4)
	s.M4차선_매수_호가,_ = 공용.F2정수64(g.BidPrice4)
	s.M4차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume4)
	s.M4차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume4)
	s.M5차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice5)
	s.M5차선_매수_호가,_ = 공용.F2정수64(g.BidPrice5)
	s.M5차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume5)
	s.M5차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume5)
	s.M6차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice6)
	s.M6차선_매수_호가,_ = 공용.F2정수64(g.BidPrice6)
	s.M6차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume6)
	s.M6차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume6)
	s.M7차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice7)
	s.M7차선_매수_호가,_ = 공용.F2정수64(g.BidPrice7)
	s.M7차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume7)
	s.M7차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume7)
	s.M8차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice8)
	s.M8차선_매수_호가,_ = 공용.F2정수64(g.BidPrice8)
	s.M8차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume8)
	s.M8차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume8)
	s.M9차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice9)
	s.M9차선_매수_호가,_ = 공용.F2정수64(g.BidPrice9)
	s.M9차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume9)
	s.M9차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume9)
	s.M10차선_매도_호가,_ = 공용.F2정수64(g.OfferPrice10)
	s.M10차선_매수_호가,_ = 공용.F2정수64(g.BidPrice10)
	s.M10차선_매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume10)
	s.M10차선_매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume10)
	s.M누적_거래량,_ = 공용.F2정수64(g.Volume)

	return s
}

//----------------------------------------------------------------------//
// 코스닥 호가 잔량 (k3)
//----------------------------------------------------------------------//
func New코스닥_호가_잔량_질의(종목_코드 string) C.Tk3InBlock {
	c := Tk3InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.", len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Tk3InBlock)(unsafe.Pointer(&c)))
}

type S코스닥_호가_잔량 struct {
	M종목_코드         string
	M시각            time.Time
	M매도_호가_모음 []int64
	M매수_호가_모음 []int64
	M매도_호가_잔량 []int64
	M매수_호가_잔량 []int64
	M누적_거래량        int64
}

func New코스닥_호가_잔량(c *C.char) *S코스닥_호가_잔량 {
	g := (*Tk3OutBlock)(unsafe.Pointer(c))
	
	지금 := time.Now()

	s := new(S코스닥_호가_잔량)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M매도_호가_모음, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.OfferPrice1,
			g.OfferPrice2,
			g.OfferPrice3,
			g.OfferPrice4,
			g.OfferPrice5,
			g.OfferPrice6,
			g.OfferPrice7,
			g.OfferPrice8,
			g.OfferPrice9,
			g.OfferPrice10})
	
	s.M매도_호가_잔량, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.OfferVolume1,
			g.OfferVolume2,
			g.OfferVolume3,
			g.OfferVolume4,
			g.OfferVolume5,
			g.OfferVolume6,
			g.OfferVolume7,
			g.OfferVolume8,
			g.OfferVolume9,
			g.OfferVolume10})
	
	s.M매수_호가_모음, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.BidPrice1,
			g.BidPrice2,
			g.BidPrice3,
			g.BidPrice4,
			g.BidPrice5,
			g.BidPrice6,
			g.BidPrice7,
			g.BidPrice8,
			g.BidPrice9,
			g.BidPrice10})
	
	s.M매수_호가_잔량, _ = 공용.F2정수64_모음(
		[]interface{}{
			g.BidVolume1,
			g.BidVolume2,
			g.BidVolume3,
			g.BidVolume4,
			g.BidVolume5,
			g.BidVolume6,
			g.BidVolume7,
			g.BidVolume8,
			g.BidVolume9,
			g.BidVolume10})
	
	s.M누적_거래량,_ = 공용.F2정수64(g.Volume)

	return s
}

//----------------------------------------------------------------------//
// 코스피 시간외 호가 잔량 (h2)
//----------------------------------------------------------------------//

func New코스피_시간외_호가_잔량_질의(종목_코드 string) C.Th2InBlock {
	c := Th2InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.", len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Th2InBlock)(unsafe.Pointer(&c)))
}

type S코스피_시간외_호가_잔량 struct {
	M종목_코드    string
	M시각       time.Time
	M매도_호가_잔량 int64
	M매수_호가_잔량 int64
}

func New코스피_시간외_호가_잔량(c *C.char) *S코스피_시간외_호가_잔량 {
	g := (*Th2OutBlock)(unsafe.Pointer(c))

	s := new(S코스피_시간외_호가_잔량)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	지금 := time.Now()
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume)
	s.M매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume)

	return s
}

//----------------------------------------------------------------------//
// 코스닥 시간외 호가 잔량 (k4)
//----------------------------------------------------------------------//

func New코스닥_시간외_호가_잔량_질의(종목_코드 string) C.Tk4InBlock {
	c := Tk4InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.", len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Tk4InBlock)(unsafe.Pointer(&c)))
}

type S코스닥_시간외_호가_잔량 struct {
	M종목_코드    string
	M시각       time.Time
	M매도_호가_잔량 int64
	M매수_호가_잔량 int64
}

func New코스닥_시간외_호가_잔량(c *C.char) *S코스닥_시간외_호가_잔량 {
	g := (*Tk4OutBlock)(unsafe.Pointer(c))
	
	지금 := time.Now()

	s := new(S코스닥_시간외_호가_잔량)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume)
	s.M매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume)

	return s
}

//----------------------------------------------------------------------//
// 코스피 예상 호가 잔량 (h3)
//----------------------------------------------------------------------//

func New코스피_예상_호가_잔량_질의(종목_코드 string) C.Th3InBlock {
	c := Th3InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.", len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Th3InBlock)(unsafe.Pointer(&c)))
}

type S코스피_예상_호가_잔량 struct {
	M종목_코드    string
	M시각       time.Time
	M동시호가_구분 string
	M예상_체결가   int64
	M예상_등락부호  uint8
	M예상_등락폭   int64
	M예상_등락율   float64
	M예상_체결_수량 int64
	M매도_호가    int64
	M매수_호가    int64
	M매도_호가_잔량 int64
	M매수_호가_잔량 int64
}

func New코스피_예상_호가_잔량(c *C.char) *S코스피_예상_호가_잔량 {
	g := (*Th3OutBlock)(unsafe.Pointer(c))
	
	지금 := time.Now()

	s := new(S코스피_예상_호가_잔량)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M동시호가_구분 = 공용.F2문자열(g.SyncOfferBid)
	s.M예상_체결가,_ = 공용.F2정수64(g.EstmPrice)
	s.M예상_등락부호 = f2등락부호(g.EstmDiffSign)
	s.M예상_등락폭,_ = 공용.F2정수64(g.EstmDiff)
	s.M예상_등락율 = f2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결_수량,_ = 공용.F2정수64(g.EstmVolume)
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume)
	s.M매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume)

	return s
}

//----------------------------------------------------------------------//
// 코스닥 예상 호가 잔량 (k5)
//----------------------------------------------------------------------//

func New코스닥_예상_호가_잔량_질의(종목_코드 string) C.Tk5InBlock {
	c := Tk5InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.", len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Tk5InBlock)(unsafe.Pointer(&c)))
}

type S코스닥_예상_호가_잔량 struct {
	M종목_코드    string
	M시각       time.Time
	M동시호가_구분 string
	M예상_체결가   int64
	M예상_등락부호  uint8
	M예상_등락폭   int64
	M예상_등락율   float64
	M예상_체결_수량 int64
	M매도_호가    int64
	M매수_호가    int64
	M매도_호가_잔량 int64
	M매수_호가_잔량 int64
}

func New코스닥_예상_호가_잔량(c *C.char) *S코스닥_예상_호가_잔량 {
	g := (*Tk5OutBlock)(unsafe.Pointer(c))
	
	지금 := time.Now()

	s := new(S코스닥_예상_호가_잔량)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M동시호가_구분 = 공용.F2문자열(g.SyncOfferBid)
	s.M예상_체결가,_ = 공용.F2정수64(g.EstmPrice)
	s.M예상_등락부호 = f2등락부호(g.EstmDiffSign)
	s.M예상_등락폭,_ = 공용.F2정수64(g.EstmDiff)
	s.M예상_등락율 = f2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결_수량,_ = 공용.F2정수64(g.EstmVolume)
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M매도_호가_잔량,_ = 공용.F2정수64(g.OfferVolume)
	s.M매수_호가_잔량,_ = 공용.F2정수64(g.BidVolume)

	return s
}

//----------------------------------------------------------------------//
// 코스피 체결 (j8)
//----------------------------------------------------------------------//

func New코스피_체결_질의(종목_코드 string) C.Tj8InBlock {
	c := Tj8InBlock{}

	if len(종목_코드) > len(c.Code) {
		종목_코드 = 종목_코드[:len(c.Code)]

		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.",len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)
	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Tj8InBlock)(unsafe.Pointer(&c)))
}

type S코스피_체결 struct {
	M종목_코드       string
	M시각          time.Time
	M등락부호        uint8
	M등락폭         int64
	M현재가         int64
	M등락율         float64
	M고가          int64
	M저가          int64
	M매도_호가       int64
	M매수_호가       int64
	M누적_거래량      int64
	M전일대비_거래량_비율 float64
	M변동_거래량      int64
	M거래_대금       int64
	M시가          int64
	M가중_평균_가격    int64
	M장구분         string
}

func New코스피_체결(c *C.char) *S코스피_체결 {
	g := (*Tj8OutBlock)(unsafe.Pointer(c))
	지금 := time.Now()

	s := new(S코스피_체결)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭,_ = 공용.F2정수64(g.Diff)
	s.M현재가,_ = 공용.F2정수64(g.MarketPrice)
	s.M등락율 = f2실수_소숫점_추가(g.DiffRate, 2)
	s.M고가,_ = 공용.F2정수64(g.High)
	s.M저가,_ = 공용.F2정수64(g.Low)
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M누적_거래량,_ = 공용.F2정수64(g.Volume)
	s.M전일대비_거래량_비율 = f2실수_소숫점_추가(g.VsPrevVolRate, 2)
	s.M변동_거래량,_ = 공용.F2정수64(g.DiffVolume)
	s.M거래_대금,_ = 공용.F2정수64(g.TrAmount)
	s.M시가,_ = 공용.F2정수64(g.Open)
	s.M가중_평균_가격,_ = 공용.F2정수64(g.WeightAvgPrice)
	s.M장구분 = 공용.F2문자열(g.Market)

	return s
}

//----------------------------------------------------------------------//
// 코스닥 체결 (k8)
//----------------------------------------------------------------------//

func New코스닥_체결_질의(종목_코드 string) C.Tk8InBlock {
	c := Tk8InBlock{}

	if len(종목_코드) > len(c.Code) {
		에러 := 공용.F에러("종목 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 종목_코드 길이  = %v.",len(c.Code), len(종목_코드))
		panic(에러)
	}

	바이트_모음 := []byte(종목_코드)

	for i := 0; i < len(바이트_모음); i++ {
		c.Code[i] = 바이트_모음[i]
	}

	return *((*C.Tk8InBlock)(unsafe.Pointer(&c)))
}

type S코스닥_체결 struct {
	M종목_코드       string
	M시각          time.Time
	M등락부호        uint8
	M등락폭         int64
	M현재가         int64
	M등락율         float64
	M고가          int64
	M저가          int64
	M매도_호가       int64
	M매수_호가       int64
	M누적_거래량      int64
	M전일대비_거래량_비율 float64
	M변동_거래량      int64
	M거래_대금       int64
	M시가          int64
	M가중_평균_가격    int64
	M장구분         string
}

func New코스닥_체결(c *C.char) *S코스닥_체결 {
	g := (*Tj8OutBlock)(unsafe.Pointer(c))
	지금 := time.Now()

	s := new(S코스닥_체결)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
 	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭,_ = 공용.F2정수64(g.Diff)
	s.M현재가,_ = 공용.F2정수64(g.MarketPrice)
	s.M등락율 = f2실수_소숫점_추가(g.DiffRate, 2)
	s.M고가,_ = 공용.F2정수64(g.High)
	s.M저가,_ = 공용.F2정수64(g.Low)
	s.M매도_호가,_ = 공용.F2정수64(g.OfferPrice)
	s.M매수_호가,_ = 공용.F2정수64(g.BidPrice)
	s.M누적_거래량,_ = 공용.F2정수64(g.Volume)
	s.M전일대비_거래량_비율 = f2실수_소숫점_추가(g.VsPrevVolRate, 2)
	s.M변동_거래량,_ = 공용.F2정수64(g.DiffVolume)
	s.M거래_대금,_ = 공용.F2정수64(g.TrAmount)
	s.M시가,_ = 공용.F2정수64(g.Open)
	s.M가중_평균_가격,_ = 공용.F2정수64(g.WeightAvgPrice)
	s.M장구분 = 공용.F2문자열(g.Market)

	return s
}

//----------------------------------------------------------------------//
// 코스피 ETF NAV (j1) (예제코드가 없음. 패딩 필드가 필요한 지 추가 확인 필요함.)
//----------------------------------------------------------------------//

type S코스피_ETF_NAV struct {
	M종목_코드    string
	M시각       time.Time
	M등락부호     uint8
	M등락폭      float64
	NAV_현재가   float64
	NAV_시가    float64
	NAV_고가    float64
	NAV_저가    float64
	M추적_오차_부호 uint8
	M추적_오차    float64
	M괴리율_부호   uint8
	M괴리율      float64
}

func New코스피_ETF(c *C.char) *S코스피_ETF_NAV {
	g := (*Tj1OutBlock)(unsafe.Pointer(c))
	지금 := time.Now()

	s := new(S코스피_ETF_NAV)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = f2실수_소숫점_추가(g.Diff, 2)
	s.NAV_현재가 = f2실수_소숫점_추가(g.Current, 2)
	s.NAV_시가 = f2실수_소숫점_추가(g.Open, 2)
	s.NAV_고가 = f2실수_소숫점_추가(g.High, 2)
	s.NAV_저가 = f2실수_소숫점_추가(g.Low, 2)
	s.M추적_오차_부호 = f2등락부호(g.TrackErrSign)
	s.M추적_오차 = f2실수_소숫점_추가(g.TrackingError, 2)
	s.M괴리율_부호 = f2등락부호(g.DivergeSign)
	s.M괴리율 = f2실수_소숫점_추가(g.DivergeRate, 2)

	return s
}

//----------------------------------------------------------------------//
// 코스닥 ETF NAV (j0) (예제코드가 없음. 패딩 필드가 필요한 지 추가 확인 필요함.)
//----------------------------------------------------------------------//

type S코스닥_ETF_NAV struct {
	M종목_코드    string
	M시각       time.Time
	M등락부호     uint8
	M등락폭      float64
	NAV_현재가   float64
	NAV_시가    float64
	NAV_고가    float64
	NAV_저가    float64
	M추적_오차_부호 uint8
	M추적_오차    float64
	M괴리율_부호   uint8
	M괴리율      float64
}

func New코스닥_ETF(c *C.char) *S코스닥_ETF_NAV {
	g := (*Tj0OutBlock)(unsafe.Pointer(c))
	지금 := time.Now()

	s := new(S코스닥_ETF_NAV)
	s.M종목_코드 = 공용.F2문자열(g.Code)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = f2실수_소숫점_추가(g.Diff, 2)
	s.NAV_현재가 = f2실수_소숫점_추가(g.Current, 2)
	s.NAV_시가 = f2실수_소숫점_추가(g.Open, 2)
	s.NAV_고가 = f2실수_소숫점_추가(g.High, 2)
	s.NAV_저가 = f2실수_소숫점_추가(g.Low, 2)
	s.M추적_오차_부호 = f2등락부호(g.TrackErrSign)
	s.M추적_오차 = f2실수_소숫점_추가(g.TrackingError, 2)
	s.M괴리율_부호 = f2등락부호(g.DivergeSign)
	s.M괴리율 = f2실수_소숫점_추가(g.DivergeRate, 2)

	return s
}

/* 코스피/코스닥 업종코드 참고표
코스피 업종명			코스닥 업종명
00 	KRX 100			01 	코스닥지수
01 	코스피지수			03 	기타서비스
02 	대형주			04 	코스닥 IT
03 	중형주			06 	제조
04 	소형주			07 	건설
05 	음식료품			08 	유통
06 	섬유의복			10 	운송
07 	종이목재			11 	금융
08 	화학				12 	통신방송서비스
09 	의약품			13 	IT S/W & SVC
10 	비금속광물			14 IT H/W
11 	철강금속			15 	음식료담배
12 	기계				16 	섬유의류
13 	전기전자			17 	종이목재
14 	의료정밀			18 	출판매체복제
15 	운수장비			19 	화학
16 	유통업			20 	제약
17 	전기가스업			21 	비금속
18 	건설업			22 	금속
19 	운수창고			23 	기계장비
20 	통신업			24 	일반전기전자
21 	금융업			25 	의료정밀기기
22 	은행				26 	운송장비부품
24 	증권				27 	기타 제조
25 	보험				28 	통신서비스
26 	서비스업			29 	방송서비스
27 	제조업			30 	인터넷
28 	코스피 200		31 	디지털컨텐츠
29 	코스피 100		32 	소프트웨어
30 	코스피 50			33 	컴퓨터서비스
32 	코스피 배당		34 	통신장비
39 	KP200 건설기계		35 	정보기기
40 	KP200 조선운송		36 	반도체
41 	KP200 철강소재		37 	IT부품
42 	KP200 에너지화학	38 	KOSDAQ 100
43 	KP200 정보기술		39 	KOSDAQ MID 300
44 	KP200 금융		40 	KOSDAQ SMALL
45 	KP200 생활소비재	43 	코스닥 스타
46 	KP200 경기소비재	44 	오락문화
47 	동일가중 KP200		45 	프리미어
48 	동일가중 KP100		46 	우량기업부
49 	동일가중 KP50		47 	벤처기업부
　	　				48 	중견기업부
　	　				49 	기술성장기업부  */

//----------------------------------------------------------------------//
// 코스피 업종 지수 (u1)
//----------------------------------------------------------------------//

func New코스피_업종_지수_질의(업종_코드 string) C.Tu1InBlock {
	c := Tu1InBlock{}

	if len(업종_코드) > len(c.SectorCode) {
		에러 := 공용.F에러("업종 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 업종_코드 길이  = %v.", len(c.SectorCode), len(업종_코드))
		panic(에러)
	}

	바이트_모음 := []byte(업종_코드)

	for i := 0; i < len(바이트_모음); i++ {
		c.SectorCode[i] = 바이트_모음[i]
	}

	return *((*C.Tu1InBlock)(unsafe.Pointer(&c)))
}

type S코스피_업종_지수 struct {
	M업종_코드  string
	M시각     time.Time
	M지수값    float64
	M등락부호   uint8
	M등락폭    float64
	M거래량    int64
	M거래_대금  int64
	M개장값    float64
	M최고값    float64
	M최고값_시각 time.Time
	M최저값    float64
	M최저값_시간 time.Time
	M지수_등락율 float64
	M거래_비중  float64
}

func New코스피_업종_지수(c *C.char) *S코스피_업종_지수 {
	g := (*Tu1OutBlock)(unsafe.Pointer(c))
	지금 := time.Now()

	s := new(S코스피_업종_지수)
	s.M업종_코드 = 공용.F2문자열(g.SectorCode)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M지수값 = f2실수_소숫점_추가(g.IndexValue, 2)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = f2실수_소숫점_추가(g.Diff, 2)
	s.M거래량,_ = 공용.F2정수64(g.Volume)
	s.M거래_대금,_ = 공용.F2정수64(g.TrAmount)
	s.M개장값 = f2실수_소숫점_추가(g.Open, 2)
	s.M최고값 = f2실수_소숫점_추가(g.High, 2)
	
	공용.F문자열_출력("시각 %v", 공용.F2문자열(g.Time))
	s.M최고값_시각, 에러 = 공용.F2포맷된_시각("포맷 문자열", g.HighTime)
	s.M최저값 = f2실수_소숫점_추가(g.Low, 2)
	
	공용.F문자열_출력("시각 %v", 공용.F2문자열(g.Time))
	s.M최저값_시간, 에러 = 공용.F2포맷된_시각("포맷 문자열", g.LowTime)
	s.M지수_등락율 = f2실수_소숫점_추가(g.DiffRate, 2)
	s.M거래_비중 = f2실수_소숫점_추가(g.TrVolRate, 2)

	return s
}

//----------------------------------------------------------------------//
// 코스닥 업종 지수 (k1)
//----------------------------------------------------------------------//

func New코스닥_업종_지수_질의(업종_코드 string) C.Tk1InBlock {
	c := Tk1InBlock{}

	if len(업종_코드) > len(c.SectorCode) {
		에러 := 공용.F에러("업종 코드 길이가 허용 한도를 초과함. " +
			"허용 한도  = %v, 업종_코드 길이  = %v.", len(c.SectorCode), len(업종_코드))
		panic(에러)
	}

	바이트_모음 := []byte(업종_코드)

	for i := 0; i < len(바이트_모음); i++ {
		c.SectorCode[i] = 바이트_모음[i]
	}

	return *((*C.Tk1InBlock)(unsafe.Pointer(&c)))
}

type S코스닥_업종_지수 struct {
	M업종_코드  string
	M시각     time.Time
	M지수값    float64
	M등락부호   uint8
	M등락폭    float64
	M거래량    int64
	M거래_대금  int64
	M개장값    float64
	M최고값    float64
	M최고값_시각 time.Time
	M최저값    float64
	M최저값_시간 time.Time
	M지수_등락율 float64
	M거래_비중  float64
}

func New코스닥_업종_지수(c *C.char) *S코스닥_업종_지수 {
	g := (*Tk1OutBlock)(unsafe.Pointer(c))
	지금 := time.Now()

	s := new(S코스닥_업종_지수)
	s.M업종_코드 = 공용.F2문자열(g.SectorCode)
	
	공용.F메모("** 가장 최근 영업일을 알아내어서 시각 값에 적용해야 함. **")
	시각, 에러 := 공용.F2포맷된_시각("15:04:05", g.Time)
	if 에러 == nil {
		s.M시각 = time.Date(지금.Year(), 지금.Month(), 지금.Day(), 
			시각.Hour(), 시각.Minute(), 시각.Second(), 0, 지금.Location())
	}
	
	s.M지수값 = f2실수_소숫점_추가(g.IndexValue, 2)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = f2실수_소숫점_추가(g.Diff, 2)
	s.M거래량,_ = 공용.F2정수64(g.Volume)
	s.M거래_대금,_ = 공용.F2정수64(g.TrAmount)
	s.M개장값 = f2실수_소숫점_추가(g.Open, 2)
	s.M최고값 = f2실수_소숫점_추가(g.High, 2)
	
	공용.F문자열_출력("시각 %v", 공용.F2문자열(g.HighTime))
	s.M최고값_시각, 에러 = 공용.F2포맷된_시각("포맷 문자열", g.HighTime)
	s.M최저값 = f2실수_소숫점_추가(g.Low, 2)
	
	공용.F문자열_출력("시각 %v", 공용.F2문자열(g.LowTime))
	s.M최저값_시간, 에러 = 공용.F2포맷된_시각("포맷 문자열", g.LowTime)
	s.M지수_등락율 = f2실수_소숫점_추가(g.DiffRate, 2)
	s.M거래_비중 = f2실수_소숫점_추가(g.TrVolRate, 2)

	return s
}
