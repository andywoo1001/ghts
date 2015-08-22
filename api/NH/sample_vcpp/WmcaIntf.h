// WmcaIntf.h: interface for the CWmcaIntf class.
//
//////////////////////////////////////////////////////////////////////

#if !defined(AFX_WMCAINTF_H__2A474430_7EB7_47E8_950B_40C7DEB33352__INCLUDED_)
#define AFX_WMCAINTF_H__2A474430_7EB7_47E8_950B_40C7DEB33352__INCLUDED_

#if _MSC_VER > 1000
#pragma once
#endif // _MSC_VER > 1000

const	DWORD	CA_WMCAEVENT		=WM_USER+8400;

const	DWORD	CA_CONNECTED		=WM_USER+110;
const	DWORD	CA_DISCONNECTED		=WM_USER+120;
const	DWORD	CA_SOCKETERROR		=WM_USER+130;
const	DWORD	CA_RECEIVEDATA		=WM_USER+210;
const	DWORD	CA_RECEIVESISE		=WM_USER+220;
const	DWORD	CA_RECEIVEMESSAGE	=WM_USER+230;
const	DWORD	CA_RECEIVECOMPLETE	=WM_USER+240;
const	DWORD	CA_RECEIVEERROR		=WM_USER+250;

typedef	BOOL	(__stdcall TLoad				)();
typedef	BOOL	(__stdcall TFree				)();
typedef	BOOL	(__stdcall TSetServer			)(const char* szServer);
typedef	BOOL	(__stdcall TSetPort				)(const int nPort);
typedef	BOOL	(__stdcall TIsConnected			)();
typedef	BOOL	(__stdcall TConnect				)(HWND hWnd,DWORD dwMsg,char cMediaType,char cUserType,const char* pszID,const char* pszPassword,const char* pszSignPassword);
typedef	BOOL	(__stdcall TDisconnect			)();
typedef	BOOL	(__stdcall TTransact			)(HWND hWnd,int nTransactionID,const char* pszTrCode,const char* pszInputData,int nInputDataSize,int nHeaderType,int nAccountIndex);
typedef	BOOL	(__stdcall TQuery   			)(HWND hWnd,int nTransactionID,const char* pszTrCode,const char* pszInputData,int nInputDataSize,int nAccountIndex);
typedef	BOOL	(__stdcall TRequest				)(HWND hWnd,int nTransactionID,const char* pszTrCode,const char* pszInputData,int nInputDataSize,const char* pszOpenBranchCode);
typedef	BOOL	(__stdcall TAttach				)(HWND hWnd,const char* pszSiseName,const char* pszInputCode,int nInputCodeSize,int nInputCodeTotalSize);
typedef	BOOL	(__stdcall TDetach				)(HWND hWnd,const char* pszSiseName,const char* pszInputCode,int nInputCodeSize,int nInputCodeTotalSize);
typedef	BOOL	(__stdcall TDetachWindow		)(HWND hWnd);
typedef	BOOL	(__stdcall TDetachAll			)();
typedef BOOL	(__stdcall TSetOption			)(const char* szKey,const char* szValue);
typedef BOOL	(__stdcall TSetAccountIndexPwd	)(const char* pszHashOut,int nAccountIndex,const char* pszPassword); 
typedef BOOL	(__stdcall TSetOrderPwd			)(const char* pszHashOut,const char* pszPassword);
typedef BOOL	(__stdcall TSetHashPwd			)(const char* pszHashOut,const char* pszKey,const char* pszPassword);
typedef BOOL	(__stdcall TSetAccountNoPwd		)(const char* pszHashOut,const char* pszAccountNo,const char* pszPassword);
typedef BOOL	(__stdcall TSetAccountNoByIndex	)(const char* pszHashOut,int nAccountIndex);

//----------------------------------------------------------------------//
// WMCA_CONNECTED �α��� ����ü
//----------------------------------------------------------------------//
typedef	struct {
	char 	szAccountNo[11];		//���¹�ȣ
	char	szAccountName[40];		//���¸�
    char	act_pdt_cdz3[3];		//��ǰ�ڵ�
    char	amn_tab_cdz4[4];		//�������ڵ�
    char	expr_datez8[8];			//���Ӹ�����
	char	granted;				//�ϰ��ֹ� ������(G:���)
    char	filler[189];			//filler
}ACCOUNTINFO;

typedef struct {
	char    szDate			[14];	// ���ӽð�
	char	szServerName	[15];	// ���Ӽ���
	char	szUserID		[8];	// ������ID
	char    szAccountCount	[3];	// ���¼�
	ACCOUNTINFO	accountlist	[999];	// ���¸��
}LOGININFO;

typedef struct{
	int       TrIndex;
	LOGININFO *pLoginInfo;
}LOGINBLOCK;

//----------------------------------------------------------------------//
// WMCA ����message ����ü
//----------------------------------------------------------------------//
typedef struct  {
	char	msg_cd		[5];	//00000:����, ��Ÿ:������(�ش� �ڵ尪�� �̿��Ͽ� �ڵ����� ������. �ڵ尪�� �������� ����� �� �ֽ��ϴ�.)
	char	user_msg	[80];
} MSGHEADER;

		 
//----------------------------------------------------------------------//
// WMCA TR ���� ����ü
//----------------------------------------------------------------------//

typedef struct {
	char*	szBlockName;
	char*	szData;
	int	nLen;
} RECEIVED;

typedef struct {
	int		  TrIndex;
	RECEIVED* pData;
} OUTDATABLOCK;


//----------------------------------------------------------------------//
// wmca.dll wrapping functions
//----------------------------------------------------------------------//

class CWmcaIntf  
{
private:
	HINSTANCE		m_hDll;

	TLoad					*m_pLoad;
	TFree					*m_pFree;
	TSetServer				*m_pSetServer;
	TSetPort				*m_pSetPort;
	TIsConnected			*m_pIsConnected;
	TConnect				*m_pConnect;
	TDisconnect				*m_pDisconnect;
	TTransact				*m_pTransact;
	TQuery   				*m_pQuery;   
	TRequest				*m_pRequest;
	TAttach					*m_pAttach;
	TDetach					*m_pDetach;
	TDetachWindow			*m_pDetachWindow;
	TDetachAll				*m_pDetachAll;
	TSetOption				*m_pSetOption;
	TSetAccountIndexPwd     *m_pSetAccountIndexPwd;
	TSetOrderPwd     		*m_pSetOrderPwd;
	TSetHashPwd     		*m_pSetHashPwd;
	TSetAccountNoPwd     	*m_pSetAccountNoPwd;
	TSetAccountNoByIndex	*m_pSetAccountNoByIndex;
	
public:
	CWmcaIntf();
	virtual ~CWmcaIntf();
public:
	BOOL Load				();
	BOOL Free				();
	BOOL Connect			(HWND hWnd, DWORD msg, char MediaType,char UserType,const char* szID,const char* szPW, const char* szCertPW);				//����� ���� �� ����
	BOOL Disconnect			();
	BOOL SetServer			(const char* szServer);																										//���Ӽ��� ����(�ʿ��)
	BOOL SetPort			(int port);																													//������Ʈ ����(�ʿ��)
	BOOL IsConnected		();
	BOOL Transact			(HWND hWnd, int nTRID, const char* szTRCode, const char* szInput, int nInputLen, int nHeadType=0, int nAccountIndex=0);		//�Ǽ��� �����Ϸ��� Transact()���
	BOOL Query				(HWND hWnd, int nTRID, const char* szTRCode, const char* szInput, int nInputLen, int nAccountIndex=0);						//Query()�Լ��� ����ϼ���.
	BOOL Request			(HWND hWnd, int nTRID, const char* szTRCode, const char* szInput, int nInputLen, const char* szOpenBranchCode=NULL);
	BOOL Attach				(HWND hWnd, const char* szBCType, const char* szInput, int nCodeLen, int nInputLen);	//�ǽð� �ü� ��û
	BOOL Detach				(HWND hWnd, const char* szBCType, const char* szInput, int nCodeLen, int nInputLen);	//�ǽð� �ü� ���
	BOOL DetachWindow		(HWND hWnd);																			//�ǽð� �ü� �ϰ����(���������)
	BOOL DetachAll			();																						//�ǽð� �ü� �ϰ����(��ü)
	BOOL SetOption			(const char* szKey,const char* szValue);
    BOOL SetAccountIndexPwd	(const char* pszHashOut,int nAccountIndex,const char* pszPassword);						//�����ε��� �� ��й�ȣ �Է�
    BOOL SetOrderPwd		(const char* pszHashOut,const char* pszPassword);										//�ŷ���й�ȣ �Է�
    BOOL SetAccountNoPwd	(const char* pszHashOut,const char* pszAccountNo,const char* pszPassword);				//���¹�ȣ �� ��й�ȣ �Է�
    BOOL SetHashPwd			(const char* pszHashOut,const char* pszKey,const char* pszPassword);		
    BOOL SetAccountNoByIndex(const char* pszHashOut,int nAccountIndex);						//�����ε����� �ش��ϴ� ���¹�ȣ
};

#endif // !defined(AFX_WMCAINTF_H__2A474430_7EB7_47E8_950B_40C7DEB33352__INCLUDED_)
