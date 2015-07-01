package internal

import (
	"encoding/gob"
	"os"
	"os/exec"
	"sync"
	"time"
)

var p파이썬_경로 string = ""
func F파이썬_스크립트_실행(
		에러_반환_채널 chan error, 타임아웃 time.Duration, 
		파이썬_스크립트_경로 string, 실행옵션 ...interface{}) {
	if p파이썬_경로 == "" {
		p파이썬_경로 = F실행화일_검색("python.exe")
	}

	실행옵션_전달 := make([]interface{}, 0)
	실행옵션_전달 = append(실행옵션_전달, 파이썬_스크립트_경로)
	실행옵션_전달 = append(실행옵션_전달, 실행옵션...)

	F외부_프로세스_실행(에러_반환_채널, 타임아웃, p파이썬_경로, 실행옵션_전달...)
}

func F외부_프로세스_실행(
		회신_채널 chan error, 타임아웃 time.Duration, 
		프로그램 string, 실행옵션 ...interface{}) {
	if !F외부_프로세스_관리_Go루틴_실행_중() {
		회신_채널 := make(chan bool)
		go F외부_프로세스_관리_Go루틴(회신_채널)
		<-회신_채널
	}
	
	go f외부_프로세스_실행_도우미_go루틴(회신_채널, 타임아웃, 프로그램, 실행옵션...)
}

func f외부_프로세스_실행_도우미_go루틴(
		에러_반환_채널 chan error, 타임아웃 time.Duration, 
		프로그램 string, 실행옵션 ...interface{}) {		
	실행옵션_문자열 := make([]string, 0)
	
	for i := 0; i < len(실행옵션); i++ {
		실행옵션_문자열 = append(실행옵션_문자열, F포맷된_문자열("%v", 실행옵션[i]))
	}
	
	외부_명령어 := exec.Command(프로그램, 실행옵션_문자열...)
	외부_명령어.Stdin = os.Stdin
	외부_명령어.Stdout = os.Stdout
	외부_명령어.Stderr = os.Stderr
	에러 := 외부_명령어.Start()
	
	if 에러 != nil {
		F에러_출력(에러)
	}
	
	에러_반환_채널 <- 에러
	
	if 에러 != nil {
		return
	}
	
	pid := 외부_명령어.Process.Pid
	
	실행_내역 :=  S외부_프로세스_실행내역 { M_pid: pid, M_기한: time.Now().Add(타임아웃)}
	Ch외부_프로세스_실행 <- 실행_내역
	
	외부_명령어.Wait()
	
	Ch외부_프로세스_정상_종료 <- pid
}

// 'PID'와 '타임아웃 되는 시점' 2가지 정보를 함께 보관하는 구조체.
// 채널에서 값을 1개만 전달할 수 있으므로 생성했음.
// []interface{}를 사용하면 더 간편하겠지만, 
// 구조체를 만들면 형 안전성도 확보되고, 파일로 저장할 때도 쓸모 있음. 
type S외부_프로세스_실행내역 struct {
	M_pid int
	M_기한 time.Time
}

var Ch외부_프로세스_실행 = make(chan S외부_프로세스_실행내역)	// 버퍼가 있으니까 엄청 어렵더라.
var Ch외부_프로세스_정상_종료 = make(chan int)	// 동기식은 간단한 데, 비동기식은 이유를 알 수 없는 에러 발생.

var ch외부_프로세스_관리_루틴_종료 chan (chan int) = make(chan (chan int))
var ch외부_프로세스_테스트용_채널_활성화 chan bool = make(chan bool)

var ch프로세스_강제_종료_테스트용 chan int = nil
var ch프로세스_정상_종료_테스트용 chan int = nil

var 외부_프로세스_관리_Go루틴_실행_중 = false
var 외부_프로세스_관리_Go루틴_잠금 = &sync.RWMutex{}

const P외부_프로세스_실행_내역_맵_파일 string = "spawned_process_list"
var 외부_프로세스_목록_파일_잠금 = &sync.RWMutex{}

func F외부_프로세스_관리_Go루틴(실행_성공_회신_채널 chan bool) {
	if F외부_프로세스_관리_Go루틴_실행_중() {
		실행_성공_회신_채널 <- false
		return
	}
	
	외부_프로세스_관리_Go루틴_잠금.Lock()
	
	// Go루틴 시작하기 전에 마지막으로 1번 더 확인.
	if 외부_프로세스_관리_Go루틴_실행_중 {
		// 그 짧은 시간 사이에 Go 루틴이 생성되었다면 바로 종료.
		외부_프로세스_관리_Go루틴_잠금.Unlock()
		
		실행_성공_회신_채널 <- false
		return
	}
	
	외부_프로세스_관리_Go루틴_실행_중 = true
	외부_프로세스_관리_Go루틴_잠금.Unlock()
	
	// 채널 버퍼 비우기.
	if len(Ch외부_프로세스_실행) > 0 {
		F문자열_출력("Ch외부_프로세스_실행 버퍼 비우기 실행")
		
		for len(Ch외부_프로세스_실행) > 0 {
			<-Ch외부_프로세스_실행
		}
	}
	
	if len(Ch외부_프로세스_정상_종료) > 0 {
		F문자열_출력("Ch외부_프로세스_정상_종료 버퍼 비우기 실행")
		
		for len(Ch외부_프로세스_정상_종료) > 0 {
			<-Ch외부_프로세스_정상_종료
		}
	}
	
	테스트용_채널_활성화 := false
	
	// 시작할 때와 종료할 때, 남은 외부 프로세스 정리.
	시작_전에_파일로_정리된_수량, 에러 := f외부_프로세스_정리_by_파일()
	if 에러 != nil {
		F에러_출력(에러)
		
		panic(에러)
	}
	
	if 시작_전에_파일로_정리된_수량 > 0 {
		F문자열_출력("시작 전에 파일로 정리된 외부 프로세스 수량 : %v.", 시작_전에_파일로_정리된_수량)
	}
	
	ticker := time.NewTicker(time.Second)	// 일정 주기마다 신호 생성.
	실행_내역_맵 := make(map[int]S외부_프로세스_실행내역)
	
	// 준비완료.
	실행_성공_회신_채널 <- true
	
	for {
		select {
		case 실행_내역 := <-Ch외부_프로세스_실행:
			실행_내역_맵[실행_내역.M_pid] = 실행_내역
			
			에러 = f실행_내역_맵_파일에_저장(실행_내역_맵)
			F에러_체크(에러)
		case pid := <-Ch외부_프로세스_정상_종료:
			_, 존재함 := 실행_내역_맵[pid]
			
			if !존재함 {
				// 유효 기한이 지나서 자동으로 강제 종료된 경우,
				// 관리 목록에서 이미 지워진 상태임.
				// 또 다시 처리해 줄 필요가 없다.
				break
			}
			
			// 유효 기간 이전에 정상적으로 실행을 마친 경우임.
			delete(실행_내역_맵, pid)
			
			에러 = f실행_내역_맵_파일에_저장(실행_내역_맵)
			F에러_체크(에러)
			
			if ch프로세스_정상_종료_테스트용 != nil {
				ch프로세스_정상_종료_테스트용 <- pid			
			}			
		case <-ticker.C:
			// 추가로 실행된 (혹은 관리해야 할) 외부 프로세스가 
			//  존재하지 않으면 굳이 처리할 필요 없음.
			if len(실행_내역_맵) == 0 { break }
			
			정리된_pid_모음, 에러 := f외부_프로세스_정리_by_맵(실행_내역_맵)
			
			// 강제 종료된 외부 프로세스가 존재하지 않으면 굳이 더 이상 처리할 필요 없음.
			if len(정리된_pid_모음) == 0 { break }
			
			for _, pid := range 정리된_pid_모음 {
				delete(실행_내역_맵, pid)
			}
			
			에러 = f실행_내역_맵_파일에_저장(실행_내역_맵)
			F에러_체크(에러)
			
			if 테스트용_채널_활성화 && len(정리된_pid_모음) > 0 {
				ch프로세스_강제_종료_테스트용 <- len(정리된_pid_모음)
			}
		case 회신_채널 := <-ch외부_프로세스_관리_루틴_종료:
			파일_기반_강제_종료_수량, 에러 := f외부_프로세스_정리_by_파일()
			F에러_체크(에러)
			
			if 테스트용_채널_활성화 {
				테스트용_채널_활성화 = false
				f테스트용_채널_활성화(테스트용_채널_활성화)
			}
			
			외부_프로세스_관리_Go루틴_잠금.Lock()
			외부_프로세스_관리_Go루틴_실행_중 = false
			외부_프로세스_관리_Go루틴_잠금.Unlock()
			
			회신_채널 <- 파일_기반_강제_종료_수량
			
			return
		case 테스트용_채널_활성화 = <-ch외부_프로세스_테스트용_채널_활성화:
			f테스트용_채널_활성화(테스트용_채널_활성화)
		}
	}
}

func F외부_프로세스_관리_Go루틴_실행_중() bool {
	외부_프로세스_관리_Go루틴_잠금.RLock()
	defer 외부_프로세스_관리_Go루틴_잠금.RUnlock()
	
	return 외부_프로세스_관리_Go루틴_실행_중
}

func F외부_프로세스_관리_Go루틴_종료() (bool, int) {
	if !F외부_프로세스_관리_Go루틴_실행_중() {
		return false, 0
	}
	
	회신_채널 := make(chan int)
	ch외부_프로세스_관리_루틴_종료 <- 회신_채널
	
	파일_기반_프로세스_강제_종료_수량 := <-회신_채널
	
	return true, 파일_기반_프로세스_강제_종료_수량  
}

func f외부_프로세스_정리_by_맵(
		실행_내역_맵 map[int]S외부_프로세스_실행내역) (
			[]int, error) {
	정리된_pid_모음 := make([]int, 0)
	
	현재_시점 := time.Now()
	
	for pid, 실행_내역 := range 실행_내역_맵 {
		// 유효 기한이 지난 외부 프로세스는 정리 대상임.
		if 실행_내역.M_기한.Before(현재_시점) {
			에러 := f프로세스_종료_by_PID(pid)
			
			if 에러 != nil {
				continue
			}
			
			정리된_pid_모음 = append(정리된_pid_모음, pid)
		}
	}
	
	return 정리된_pid_모음, nil
}
 
func f외부_프로세스_정리_by_파일() (int, error) {
	var 정리된_프로세스_수량 = 0
	
	실행_내역_맵, 에러 := f실행_내역_맵_읽기()
	
	if 에러 != nil {
		F에러_출력(에러)
		
		return 0, 에러
	}
	
	for pid, _ := range 실행_내역_맵 {
		에러 := f프로세스_종료_by_PID(pid)
			
		if 에러 != nil {
			continue
		}
		
		정리된_프로세스_수량++
	}
	
	에러 = f실행_내역_맵_파일_초기화()
	
	if 에러 != nil {
		F에러_출력(에러)
		
		return 정리된_프로세스_수량, 에러
	}
	
	return 정리된_프로세스_수량, nil
}

func f프로세스_종료_by_PID(pid int) error {
	프로세스, 에러 := os.FindProcess(pid)
	
	if 에러 != nil {
		return 에러
	}
	
	에러 = 프로세스.Kill()
	
	if 에러 != nil {
		//F에러_출력(에러)
		
		return 에러
	}
	
	return nil
}

// 이하 외부 프로세스 목록 파일 관리 함수 모음
func f실행_내역_맵_파일_존재함() (bool, error) {
	외부_프로세스_목록_파일_잠금.RLock()
	defer 외부_프로세스_목록_파일_잠금.RUnlock()
	
	_, 에러 := os.Stat(P외부_프로세스_실행_내역_맵_파일)
	
	if os.IsNotExist(에러) {
		return false, nil
	}
	
	if 에러 == nil {
		return true, nil
	}
	
	F문자열_출력("예상치 못한 경우.\n%v", 에러)
	
	return false, 에러
}

func f실행_내역_맵_파일_초기화() error {
	비어있는_맵 := make(map[int]S외부_프로세스_실행내역)
	
	return f실행_내역_맵_파일에_저장(비어있는_맵)
}

func f실행_내역_맵_읽기() (map[int]S외부_프로세스_실행내역, error) {
	존재함, 에러	:= f실행_내역_맵_파일_존재함()
	
	if 에러 != nil {
		F에러_출력(에러)
		return nil, 에러
	}
	
	if !존재함 {
		에러 = f실행_내역_맵_파일_초기화()
		
		if 에러 != nil {
			F에러_출력(에러)
			return nil, 에러
		}
	}
	
	외부_프로세스_목록_파일_잠금.RLock()
	defer 외부_프로세스_목록_파일_잠금.RUnlock()
	
	파일, 에러 := os.Open(P외부_프로세스_실행_내역_맵_파일)
	
	if 에러 != nil {
		F에러_출력(에러)
		return nil, 에러
	}
	
	defer 파일.Close()
	
	var 실행_내역_맵 map[int]S외부_프로세스_실행내역
	
	디코더 := gob.NewDecoder(파일)	
	에러 = 디코더.Decode(&실행_내역_맵)
	
	if 에러 != nil {
		if 에러.Error() == "EOF" {
			//F문자열_출력("파일 내용 없음.")
			
			// 파일에 아무 내용이 없으므로 비어있는 맵을 새로 생성해서 반환
			비어있는_맵 := make(map[int]S외부_프로세스_실행내역)
			
			return 비어있는_맵, nil	 
		}
		
		F에러_출력(에러)
		return nil, 에러
	}
	
	return 실행_내역_맵, nil
}

func f실행_내역_맵_파일에_저장(실행_내역_맵 map[int]S외부_프로세스_실행내역) error {
	외부_프로세스_목록_파일_잠금.Lock()
	defer 외부_프로세스_목록_파일_잠금.Unlock()
	
	파일, 에러 := os.Create(P외부_프로세스_실행_내역_맵_파일)
	
	if 에러 != nil {
		F에러_출력(에러)
		return 에러
	}
	
	defer 파일.Close()
	
	인코더 := gob.NewEncoder(파일)
	에러 = 인코더.Encode(실행_내역_맵)
	
	if 에러 != nil {
		F에러_출력(에러)
		return 에러
	}
	
	for i:=0 ; i<10 ; i++ {
		에러 = 파일.Sync()
		
		if 에러 == nil {
			break
		}
	} 
	
	return nil
}

func f테스트용_채널_활성화(활성화 bool) {
	if 활성화 {
		ch프로세스_정상_종료_테스트용 = make(chan int, 1000)
		ch프로세스_강제_종료_테스트용 = make(chan int, 1000)	
	} else {
		ch프로세스_정상_종료_테스트용 = nil
		ch프로세스_강제_종료_테스트용 = nil
	}
}