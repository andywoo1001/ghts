// WMCALOADERDlg.cpp : implementation file
//
#include "stdafx.h"
#include "WMCALOADER.h"
#include "WMCALOADERDlg.h"

#ifdef _DEBUG
#define new DEBUG_NEW
#undef THIS_FILE
static char THIS_FILE[] = __FILE__;
#endif

#include "trio_inv.h"
#include "trio_ord.h"
#include "math.h"

const	int	TRID_c1101	=1;		//<--���ø� ���̱� �����̸� �������� �ʿ�� �����ϴ�
const	int	TRID_c8201	=2;		//<--���ø� ���̱� �����̸� �������� �ʿ�� �����ϴ�
const	int	TRID_c8102	=3;		//<--���ø� ���̱� �����̸� �������� �ʿ�� �����ϴ�

/************************************************************************************

							!!! ���� !!!

	�� ������ OPEN API library(wmca.dll)�� ������� �����ϱ� ���� ������ �Ұ��մϴ�.
	�� ���������δ� ����ڰ� �ʿ���ϴ� ��� ������ Ȯ���� ���� �����Ƿ�
	�⺻���� ������� ���� �� �� ����ں��� �ʿ��� ����� ������ �����Ǵ� 
	�ڷ�(*.doc)�� ���� ���� �ۼ��Ͻñ� �ٶ��ϴ�.

	����, �� �������� �ǵ����� ���� ���α׷� ������ ���Ե� �� ������ 
	��ȯ�濡 �ٷ� ����Ͻô� ���� �������� �ʽ��ϴ�.
	�뵵�� �°� ������Ʈ�� ���� �����ؼ� �ۼ��Ͻñ⸦ �����մϴ�.

	* ��翡���� �� ������ ���� ���α׷� �����ۿ� ���ؼ� ��� å�ӵ� ���� �ʽ��ϴ�.
	* �� �ҽ� �ڵ�� ����� ��û �Ǵ� ��� ���� ������ ���� �������� ����� �� �ֽ��ϴ�.
	* ������ ���� ���� I/O(TR �� �ǽð��ü� ��Ŷ)�� �������� ����� ���ɼ��� �ֽ��ϴ�.
	* �� ������ Microsoft Visual C++ 6.0, Microsoft Visual C++ 2008 ȯ�濡�� �ۼ��Ǿ����ϴ�.
	* wmca.dll �� ���� �ڵ�� Ansi code��� 32bit �������� �����Ǹ�
	* 64bit �Ǵ� Unicode��������� ������ ������� �ʽ��ϴ�.

************************************************************************************/





/************************************************************************************

	wmca.dll �� wrapping �� CWmcaIntf class�� �̿��Ͽ� ���� ��ɵ��� �׽�Ʈ�մϴ�.
	
	0. ����(�������+��������)
	1. �������� ��ȸ(���簡)		c1101
	2. �������� ��ȸ(�ܰ���ȸ)	c8102
	3. �ǽð� �ü� ����			j8(�ڽ��� �ֽ� ���簡)
	4. �ǽð� ü���뺸			d2


************************************************************************************/



/////////////////////////////////////////////////////////////////////////////
// CWMCALOADERDlg dialog
//������������������������������������������������������������������������������������������������
CWMCALOADERDlg::CWMCALOADERDlg(CWnd* pParent /*=NULL*/)
	: CDialog(CWMCALOADERDlg::IDD, pParent)
{
	//{{AFX_DATA_INIT(CWMCALOADERDlg)
		// NOTE: the ClassWizard will add member initialization here
	//}}AFX_DATA_INIT
	m_hIcon = AfxGetApp()->LoadIcon(IDR_MAINFRAME);
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::DoDataExchange(CDataExchange* pDX)
{
	CDialog::DoDataExchange(pDX);
	//{{AFX_DATA_MAP(CWMCALOADERDlg)
	DDX_Control(pDX, IDC_LIST1, m_listboxTrace);
	DDX_Control(pDX, IDC_COMBOACCOUNTLIST, m_comboAccountList);
	DDX_Control(pDX, IDC_STATICLOGIN, m_staticLoginTime);
	DDX_Control(pDX, IDC_EDIT3, m_editSignPassword);
	DDX_Control(pDX, IDC_EDIT2, m_editPassword);
	DDX_Control(pDX, IDC_EDIT1, m_editID);
	DDX_Control(pDX, IDC_CONNECT, m_buttonConnect);
	DDX_Control(pDX, IDC_DISCONNECT, m_buttonDisconnect);
	DDX_Control(pDX, IDC_BALANCE, m_buttonBalance);
	//}}AFX_DATA_MAP
}
//������������������������������������������������������������������������������������������������
BEGIN_MESSAGE_MAP(CWMCALOADERDlg, CDialog)
	//{{AFX_MSG_MAP(CWMCALOADERDlg)
	ON_WM_PAINT()
	ON_WM_QUERYDRAGICON()
	ON_BN_CLICKED(IDC_CONNECT, OnConnect)
	ON_BN_CLICKED(IDC_DISCONNECT, OnDisconnect)
	ON_BN_CLICKED(IDC_BALANCE, OnBalance)
	ON_BN_CLICKED(IDC_BUTTON_CURRENT, OnButtonCurrent)
	ON_BN_CLICKED(IDC_BUTTON_CLEAR, OnButtonClear)
	ON_BN_CLICKED(IDC_ORDER, OnButtonOrder)
	//}}AFX_MSG_MAP

	ON_MESSAGE(CA_WMCAEVENT, OnWmcaEvent)
END_MESSAGE_MAP()
//������������������������������������������������������������������������������������������������
// If you add a minimize button to your dialog, you will need the code below
//  to draw the icon.  For MFC applications using the document/view model,
//  this is automatically done for you by the framework.
//������������������������������������������������������������������������������������������������

/////////////////////////////////////////////////////////////////////////////
// CWMCALOADERDlg message handlers
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnPaint() 
{
	if (IsIconic())
	{
		CPaintDC dc(this); // device context for painting

		SendMessage(WM_ICONERASEBKGND, (WPARAM) dc.GetSafeHdc(), 0);

		// Center icon in client rectangle
		int cxIcon = GetSystemMetrics(SM_CXICON);
		int cyIcon = GetSystemMetrics(SM_CYICON);
		CRect rect;
		GetClientRect(&rect);
		int x = (rect.Width() - cxIcon + 1) / 2;
		int y = (rect.Height() - cyIcon + 1) / 2;

		// Draw the icon
		dc.DrawIcon(x, y, m_hIcon);
	}
	else
	{
		CDialog::OnPaint();
	}
}
//������������������������������������������������������������������������������������������������
HCURSOR CWMCALOADERDlg::OnQueryDragIcon()
{
	return (HCURSOR) m_hIcon;
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnOK() 
{
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnCancel() 
{
	CDialog::OnCancel();
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnConnect() 
{
	CString	strID;
	CString	strPassword;
	CString	strSignPassword;

	m_editID.GetWindowText(strID);
	m_editPassword.GetWindowText(strPassword);
	m_editSignPassword.GetWindowText(strSignPassword);

	//���� �� �α���
	//��ü�ڵ�� Ư���� ��츦 �����ϰ� �׻� �Ʒ� �⺻���� ����Ͻñ� �ٶ��ϴ�.
	//m_wmca.Connect(GetSafeHwnd(),CA_WMCAEVENT,'P','1',strID,strPassword,strSignPassword);	//mug OpenAPI ����ڿ�
	m_wmca.Connect(GetSafeHwnd(),CA_WMCAEVENT,'T','W',strID,strPassword,strSignPassword);	//TX OpenAPI ����ڿ�
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnDisconnect() 
{
	m_wmca.Disconnect();
}

//������������������������������������������������������������������������������������������������
//	���ڿ��� ����ü �ʵ�� ������ ä��� ��ƿ��Ƽ �Լ��Դϴ�. (��ġ�� �߸�)
//������������������������������������������������������������������������������������������������
void	_smove(char* szTarget,int nSize,CString strSource)
{
	int	nMin	=min(nSize,strSource.GetLength());

	memset(szTarget,' ',nSize);
	strncpy(szTarget,strSource,nMin);
}

#define SMOVE(fTarget,strSource)	_smove(fTarget,sizeof fTarget,strSource)

//������������������������������������������������������������������������������������������������
//	���ڸ� ����ü �ʵ�� ������ ä��� ��ƿ��Ƽ �Լ��Դϴ�.	  (��ġ�� ����)
//������������������������������������������������������������������������������������������������
void	_nmove(char* szTarget,int nTarget,CString strSource)
{
	if(strSource.GetLength()>nTarget)
		return;

	memset(szTarget,'0',nTarget);
	strncpy(szTarget+nTarget-strSource.GetLength(),strSource,strSource.GetLength());
}

#define NMOVE(fTarget,strSource)	_nmove(fTarget,sizeof fTarget,strSource)

//������������������������������������������������������������������������������������������������
//	����ü �ʵ尪�� ���ڿ� ���·� ��ȯ�ϴ� ��ƿ��Ƽ �Լ��Դϴ�.
//������������������������������������������������������������������������������������������������
CString	_scopy(const char* szData,int nSize)
{
	char	szBuf[256];				//�ʵ��� �ִ� ũ��� ��Ȳ�� ���� ������ �ʿ䰡 ����
	memset(szBuf,0,sizeof szBuf);
	strncpy(szBuf,szData,nSize);
	
	return szBuf;
}

#define SCOPY(x)	_scopy(x,sizeof x)

//������������������������������������������������������������������������������������������������
//	����ü �ʵ尪�� ���� ���·� ��ȯ�ϴ� ��ƿ��Ƽ �Լ��Դϴ�.
//������������������������������������������������������������������������������������������������
CString	_ncopy(const char* szData,int nSize)
{
	CString	strInput	=_scopy(szData,nSize);

	for(int i=0;i<strInput.GetLength();i++)
		if(strInput.GetAt(i)=='0')
		{
			//������ '0' �Ǵ� ������ �Ҽ����� ���� ��쿡�� ������ �ʴ´�
			if(i!=strInput.GetLength()-1 && strInput.GetAt(i+1)!='.')	
				strInput.SetAt(i,' ');
		}
		else if(strInput.GetAt(i)=='-')	//minus
			;
		else if(strInput.GetAt(i)==' ')	//�������ӿ��� ���鹮�ڰ� ���� ��찡 ����
			;
		else
			break;

	return	strInput;
}

#define NCOPY(x)	_ncopy(x,sizeof x)

//������������������������������������������������������������������������������������������������
//	 ���ڿ� õ ������ �ĸ� �����ϴ� ��ƿ��Ƽ �Լ��Դϴ�.
//������������������������������������������������������������������������������������������������
CString	_comma(const char* szData,int nSize,int decimal=0)
{
	CString	strInput	=_scopy(szData,nSize);
	strInput.TrimLeft();
	strInput.TrimRight();
	
	double	fInput	=atof(strInput);

	CString	strTemp;
	strTemp.Format("%0.*f",decimal,fInput);

	if(fInput>=1000000000.0F)
		strTemp.Insert(strTemp.GetLength()-9,',');
	if(fInput>=1000000.0F)
		strTemp.Insert(strTemp.GetLength()-6,',');
	if(fInput>=1000.0F)
		strTemp.Insert(strTemp.GetLength()-3,',');

	CString	strOutput;
	strOutput.Format("%*s",nSize,strTemp);

	return strOutput;
}

#define COMMA(x)	_comma(x,sizeof x)
#define COMMAF(x,d)	_comma(x,sizeof x,d)

//������������������������������������������������������������������������������������������������
//	wmca.dll�� ���� ������ ������ �޽����� ���� �� �̺�Ʈ �ڵ鷯�� �б��մϴ�
//������������������������������������������������������������������������������������������������
LRESULT CWMCALOADERDlg::OnWmcaEvent(WPARAM dwMessageType, LPARAM lParam)
{
	switch(dwMessageType) {
	case CA_CONNECTED:			//�α��� ����
		OnWmConnected( (LOGINBLOCK*)lParam );
		break;
	case CA_DISCONNECTED:		//���� ����
		OnWmDisconnected();
		break;
	case CA_SOCKETERROR:		//��� ���� �߻�
		OnWmSocketerror( (int)lParam );
		break;
	case CA_RECEIVEDATA:		//���� ���� ����(TR)
		OnWmReceivedata( (OUTDATABLOCK*)lParam );
		break;
	case CA_RECEIVESISE:		//�ǽð� ������ ����(BC)
		OnWmReceivesise( (OUTDATABLOCK*)lParam );
		break;
	case CA_RECEIVEMESSAGE:		//���� �޽��� ���� (�Է°��� �߸��Ǿ��� ��� ���ڿ����·� ������ ���ŵ�)
		OnWmReceivemessage( (OUTDATABLOCK*)lParam );
		break;
	case CA_RECEIVECOMPLETE:	//���� ó�� �Ϸ�
		OnWmReceivecomplete( (OUTDATABLOCK*)lParam );
		break;
	case CA_RECEIVEERROR:		//���� ó���� ���� �߻� (�Է°� ������)
		OnWmReceiveerror( (OUTDATABLOCK*)lParam );
		break;
	default:
		break;
	}

	return TRUE;
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmConnected( LOGINBLOCK* pLogin )
{
	//�α����� �����ϸ�, ���ӽð� �� ���¹�ȣ ������ �޾� ������ ����/����մϴ�.
	//���¹�ȣ�� ���� ����(�ε���)�� ���°��� ���� ȣ��� ���ǹǷ� �߿��մϴ�.
	
	//��Ʈ�� Ȱ�����º���
	m_editID.EnableWindow(FALSE);
	m_editPassword.EnableWindow(FALSE);
	m_editSignPassword.EnableWindow(FALSE);

	m_buttonConnect.EnableWindow(FALSE);
	m_buttonDisconnect.EnableWindow(TRUE);

	//���ӽð� ���
	char	szText[256]	={0};
	strncpy(szText,pLogin->pLoginInfo->szDate,sizeof pLogin->pLoginInfo->szDate);
	CString	strText	="���ӽð� : ";
	strText	+=szText;
	m_staticLoginTime.SetWindowText(strText);

	//���¹�ȣ ���
	char	szAccountCount[8]	={0};
	strncpy(szAccountCount,pLogin->pLoginInfo->szAccountCount,sizeof pLogin->pLoginInfo->szAccountCount);
	int	nAccountCount	=atoi(szAccountCount);

	m_comboAccountList.ResetContent();
	for(int i=0;i<nAccountCount;i++)
	{
		char	szAccountNo[16]	={0};
		strncpy(szAccountNo,(char*)&pLogin->pLoginInfo->accountlist[i].szAccountNo,sizeof pLogin->pLoginInfo->accountlist[i].szAccountNo);

		CString	strAccountNo	=szAccountNo;
		m_comboAccountList.AddString(strAccountNo);

		//strAccountNo.Insert(3,'-');
		//strAccountNo.Insert(6,'-');

		//CString	strIndex;
		//strIndex.Format("����%3d:%s",i+1,strAccountNo);		//���¹�ȣ�� '1'������ �����մϴ�. 

		//m_comboAccountList.AddString(strIndex);
	}

	//�ڵ����� ù ��° ���¸� ����
	if(nAccountCount)
		m_comboAccountList.SetCurSel(0);

	//
	m_listboxTrace.AddString("���� ���ӵǾ����ϴ�");
	ScrollDown();
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmDisconnected()
{
	//������ ������ ��� ó��
	m_staticLoginTime.SetWindowText("");

	m_editID.EnableWindow(TRUE);
	m_editPassword.EnableWindow(TRUE);
	m_editSignPassword.EnableWindow(TRUE);

	m_buttonConnect.EnableWindow(TRUE);
	m_buttonDisconnect.EnableWindow(FALSE);
	
	//
	m_listboxTrace.AddString("������ ������ϴ�");
	ScrollDown();
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmSocketerror(int socket_error_code)
{
	m_listboxTrace.AddString("��ſ����߻�");
	ScrollDown();
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmReceivedata( OUTDATABLOCK* pOutData )
{
	switch(pOutData->TrIndex)
	{
	case TRID_c1101:
		//////////////////////////////////////////////////////////////////////////
		//�ݺ����� �ʴ� �ܼ���� ó�� ���
		//////////////////////////////////////////////////////////////////////////

		if(strcmp(pOutData->pData->szBlockName,"c1101OutBlock")==0)			
		{
			Tc1101OutBlock*	pc1101outblock	=(Tc1101OutBlock*)pOutData->pData->szData;
			
			m_listboxTrace.AddString(">>  �ֽ����簡��ȸ - ���簡");
			
			CString	strList;
			strList.Format("%10s %6s %20s %15s�� %15s��",
				SCOPY(pc1101outblock->hotime),
				SCOPY(pc1101outblock->code),
				SCOPY(pc1101outblock->hname),
				COMMA(pc1101outblock->price),
				COMMA(pc1101outblock->volume)
				);
			
			m_listboxTrace.AddString(strList);
			ScrollDown();
		}
		//////////////////////////////////////////////////////////////////////////
		//�ݺ������� ��¿� ���� ó�� ���
		//////////////////////////////////////////////////////////////////////////

		//�ֽ� ���簡/�����ŷ����ڷ�	
		if(strcmp(pOutData->pData->szBlockName,"c1101OutBlock2")==0)		
		{
			Tc1101OutBlock2*	pc1101outblock2	=(Tc1101OutBlock2*)pOutData->pData->szData;
			
			m_listboxTrace.AddString(">>  �ֽ����簡��ȸ - �����ŷ���");
			
			//���� �����Ϳ� ���� ���� ������ ���� ���� �������̹Ƿ�
			//�����ڷ� ũ�⸦ ����ü ũ��� ������ �� �� �ݺ��Ǵ��� �����
			int		nOccursCount	=pOutData->pData->nLen/sizeof(Tc1101OutBlock2);			//�ݺ� ȸ���� ����
			
			for(int i=0;i<nOccursCount;i++)	//ȸ����ŭ �ݺ��ϸ鼭 ������ ����
			{
				CString	strList;
				strList.Format("%10s %9s�� %7s�� %7s�� %7s�� %10s�� %12s��",
					SCOPY(pc1101outblock2->time),
					COMMA(pc1101outblock2->price),
					COMMA(pc1101outblock2->change),
					COMMA(pc1101outblock2->offer),
					COMMA(pc1101outblock2->bid),
					COMMA(pc1101outblock2->movolume),
					COMMA(pc1101outblock2->volume)
					);
				
				m_listboxTrace.AddString(strList);
				
				pc1101outblock2++;	//���������� ������ �̵�
			}

			ScrollDown();
		}
		break;
	case TRID_c8201:
		//////////////////////////////////////////////////////////////////////////
		//�ݺ����� �ʴ� �ܼ���� ó�� ���
		//////////////////////////////////////////////////////////////////////////
		if(strcmp(pOutData->pData->szBlockName,"c8201OutBlock")==0)			
		{
			Tc8201OutBlock*	pc8201outblock	=(Tc8201OutBlock*)pOutData->pData->szData;
			
			m_listboxTrace.AddString("*** �ܰ���ȸ ���� ***");
			m_listboxTrace.AddString("  �����Ѿ�    ="+COMMA(pc8201outblock->dpsit_amtz16));
			m_listboxTrace.AddString("  ���ű��Ѿ�    ="+COMMA(pc8201outblock->chgm_pos_amtz16));
			m_listboxTrace.AddString("  �ֹ����ɱݾ�  ="+COMMA(pc8201outblock->pos_csamt4z16));

			ScrollDown();
		}
		//////////////////////////////////////////////////////////////////////////
		//�ݺ������� ��¿� ���� ó�� ���
		//////////////////////////////////////////////////////////////////////////
		else if(strcmp(pOutData->pData->szBlockName,"c8201OutBlock1")==0)			
		{
			Tc8201OutBlock1*	pc8201outblock1	=(Tc8201OutBlock1*)pOutData->pData->szData;

			int		nOccursCount	=pOutData->pData->nLen/sizeof(Tc8201OutBlock1);			//���� ũ�⸦ ���ڵ� ũ��� ����� ���� �ݺ� ȸ���� �� �� ����
			m_listboxTrace.AddString("*** ����������ȸ ���� ***");

			for(int i=0;i<nOccursCount;i++)	//ȸ����ŭ �ݺ��ϸ鼭 ������ ����
			{
				//�Ʒ� �׸���� �ܼ��� ���ø� ���� ���� �����Դϴ�. �ʿ信 ���� ÷���Ͻñ� �ٶ��ϴ�.
				//��� HTS�� ��ȸ�Ǵ� ����� �����ϸ� ���� ������ �����ٿ� ���� ǥ��� ��� ������ �����ڵ尡 �����Ǳ⵵ �մϴ�.
				m_listboxTrace.AddString("  �����ڵ�      ="	+SCOPY(pc8201outblock1->issue_codez6));
				m_listboxTrace.AddString("  �����        ="	+SCOPY(pc8201outblock1->issue_namez40));
				m_listboxTrace.AddString("  ���űݷ�(%)   ="	+COMMA(pc8201outblock1->issue_mgamt_ratez6));
				m_listboxTrace.AddString("  �ܰ�����      ="	+SCOPY(pc8201outblock1->bal_typez6));
				m_listboxTrace.AddString("  �̰���        ="	+COMMA(pc8201outblock1->unstl_qtyz16));
				m_listboxTrace.AddString("  �����ܰ�      ="	+COMMA(pc8201outblock1->jan_qtyz16));
				m_listboxTrace.AddString("  ��ո��԰�    ="	+COMMA(pc8201outblock1->slby_amtz16));
				m_listboxTrace.AddString("  �򰡼���(õ��)="	+COMMA(pc8201outblock1->lsnpf_amtz16));
				m_listboxTrace.AddString("  ���ͷ�        ="	+COMMAF(pc8201outblock1->earn_ratez9,2));
				m_listboxTrace.AddString("  �򰡱ݾ�      ="	+COMMA(pc8201outblock1->ass_amtz16));
				m_listboxTrace.AddString(" ");

				pc8201outblock1++;	//���������� ������ �̵�
			}

			ScrollDown();
		}
		break;
	case TRID_c8102:
		//////////////////////////////////////////////////////////////////////////
		//�ݺ����� �ʴ� �ܼ���� ó�� ���
		//////////////////////////////////////////////////////////////////////////
		if(strcmp(pOutData->pData->szBlockName,"c8102OutBlock")==0)			
		{
			Tc8102OutBlock*	pc8102OutBlock	=(Tc8102OutBlock*)pOutData->pData->szData;
			
			CString	strOrderNum		=SCOPY(pc8102OutBlock->order_noz10);
			strOrderNum.TrimLeft();
			strOrderNum.TrimRight();

			if(strOrderNum.IsEmpty())
			{
				m_listboxTrace.AddString(">>  �ֹ� ����");
				break;
			}
			
			m_listboxTrace.AddString(">>  �ֹ� ����");
			m_listboxTrace.AddString("�ֹ���ȣ ="+strOrderNum);
			
			ScrollDown();
		}
		break;
	}
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::ScrollDown()
{
//	m_listboxTrace.AddString("");
	m_listboxTrace.SetTopIndex(	m_listboxTrace.GetCount()-1	);
	m_listboxTrace.SetSel(m_listboxTrace.GetCount()-1);
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmReceivesise( OUTDATABLOCK* pSiseData )
{
	//�ڽ��� �ֽ����簡 �ǽð� �ü� ����
	if(strncmp(pSiseData->pData->szBlockName,"j8",2)==0)		
	{
		Tj8OutBlock*	pj8	=(Tj8OutBlock*)(pSiseData->pData->szData+3);	//���� 3����Ʈ�� ��Ŷ������ ���౸���̹Ƿ� skip
		
		CString	strOut;
		strOut.Format("�ڽ��� �ֽ� ���簡(�ǽð�):  %10s %6s %15s %15s",
				SCOPY(pj8->time),
				SCOPY(pj8->code),
				COMMA(pj8->price),
				COMMA(pj8->volume)
		);

		m_listboxTrace.AddString(strOut);
		ScrollDown();
	}
	
	//����) �ǽð� ü���뺸	-	�ǽð� ü���뺸�� ������ Attach()�Լ��� ȣ������ �ʾƵ� �ڵ� ���ŵ˴ϴ�
	//d2 ����ü ������ ���� �ʿ�� �ϴ� ������ ������ �� �ֽ��ϴ�.
	else if(strncmp(pSiseData->pData->szBlockName,"d2",2)==0)	
	{
		Td2OutBlock*	pd2	=(Td2OutBlock*)(pSiseData->pData->szData+3);	//���� 3����Ʈ�� ��Ŷ������ ���౸���̹Ƿ� skip

		//
		CString	strResult;
		strResult.Format("�ǽð� ü���뺸- �ð�(%s) ���¹�ȣ(%s) ����(%15s) ����(%15s)",SCOPY(pd2->conctime),
			SCOPY(pd2->accountno),COMMA(pd2->concgty),COMMA(pd2->concprc));
		m_listboxTrace.AddString(strResult);
		ScrollDown();
	}
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmReceivemessage( OUTDATABLOCK* pMessage )
{
	//���� ���¸� ���ڿ� ���·� �����
	MSGHEADER*	pMsgHeader		=(MSGHEADER*)pMessage->pData->szData;

	CString	strMsg;
	strMsg.Format("[%10d]  %5s : %s",
			pMessage->TrIndex,
			SCOPY(pMsgHeader->msg_cd),
			SCOPY(pMsgHeader->user_msg)	
	);

	m_listboxTrace.AddString(strMsg);
	ScrollDown();
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmReceivecomplete( OUTDATABLOCK* pOutData )
{
	//TR ȣ�� �Ϸ������ ���� �ؾ��Ѵٸ� �Ʒ��� ���� ó���� �� �ֽ��ϴ�
	switch(pOutData->TrIndex)
	{
	case TRID_c1101:	;
		m_listboxTrace.AddString("�ֽ� ���簡 ��ȸ �Ϸ�");
		ScrollDown();
		break;
	case TRID_c8201:	;
		m_listboxTrace.AddString("���� �ܰ� ��ȸ �Ϸ�");
		ScrollDown();
		break;
	}
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnWmReceiveerror( OUTDATABLOCK* pError )
{
	//���� ȣ���� ������ ��� ���ŵ�
	m_listboxTrace.AddString("*** ERROR ***");
	m_listboxTrace.AddString(pError->pData->szData);
	ScrollDown();
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnBalance()					//�ܰ� ��ȸ(Query�Լ� ��� ����)
{
    //�ܰ���ȸ�� ���� ����� ���¹�ȣ �ε����� ���մϴ�.
	//���¹�ȣ �ε����� '1'���� ���۵˴ϴ�.
	int	nAccountIndex	=m_comboAccountList.GetCurSel()+1;

	//�ܰ���ȸ ���񽺿��� �䱸�ϴ� �Է°��� ������ ����ü �����Դϴ�.
	//�Ӽ�����Ʈ ������ �����ϱ� ���� ������ �Է±���ü�� ����ϴ� ���� �����մϴ�.
	//��� OpenAPI�� ó�� ����Ͻô� �е鲲�� ���� ���� �Ǽ��ϴ� �κ��Դϴ�. 
	//
	//ex)'�Է°��� ����� �־��ٰ� �����ϴµ� ���� ���� �ʴٰ� ��� �źε˴ϴ�. ���� �߸��� �ɱ��?'

	Tc8201InBlock	c8201inblock;
	memset(&c8201inblock,0x20,sizeof Tc8201InBlock);	//���鹮��(space)�� �ʱ�ȭ�մϴ�.

	//ȭ����� ���� �� ��й�ȣ�� OpenAPI���� �����ǰ� �ִ� hash�Լ��� �Է��ϸ�
	//������ �Է��ʵ忡 hash ��й�ȣ�� ä�����ϴ�.
	//��� hash ��й�ȣ�� ���̴� 44����Ʈ �����Դϴ�.
	CString	strPassword;
	GetDlgItem(IDC_EDIT_PASSWORD)->GetWindowText(strPassword);
	m_wmca.SetAccountIndexPwd(c8201inblock.pswd_noz44,nAccountIndex,strPassword);

	//�� �Է� �ʵ忡�� �䱸�ϴ� ���鿡 ���� ���Ǵ� *.doc ������ ���� Ȯ���� �� �ֽ��ϴ�.
	SMOVE(c8201inblock.bnc_bse_cdz1,"1");			//����	

	m_wmca.Query(
			GetSafeHwnd(),			//�� ������� ���� �޽����� �ްڽ��ϴ�.
			TRID_c8201,				//�� ���񽺿� ���ؼ� TRID_c8201(3) �ĺ��ڸ� ���̰ڴٴ� �ǹ��̸� �ݵ�� ����� �ʿ�� �����ϴ�.
			"c8201",				//ȣ���Ϸ��� ���� �ڵ�� 'c8201' �Դϴ�.
			(char*)&c8201inblock,	//c8201���� �䱸�ϴ� �Է� ����ü �����͸� �����մϴ�
			sizeof Tc8201InBlock,	//�Է� ����ü ũ���Դϴ�
			nAccountIndex			//���¹�ȣ �ε����� '1'���� ���۵˴ϴ�.
	);
}
//������������������������������������������������������������������������������������������������
void CWMCALOADERDlg::OnButtonCurrent()	//���簡 ��ȸ
{
	//�ֽ����簡��ȸ ���񽺿��� �䱸�ϴ� �Է°��� ������ ����ü �����Դϴ�.
	Tc1101InBlock	c1101inblock;
	memset(&c1101inblock,0x20,sizeof Tc1101InBlock);

	//ȭ����� �����ڵ带 ���մϴ�
	CString	strJCode;
	GetDlgItem(IDC_EDIT_JCODE)->GetWindowText(strJCode);
	
	//�� �Է� �ʵ忡�� �䱸�ϴ� ���鿡 ���� ���Ǵ� *.doc ������ ���� Ȯ���� �� �ֽ��ϴ�.
	SMOVE(c1101inblock.formlang,"k");
	SMOVE(c1101inblock.code,strJCode);
	
	//�ֽ� ���簡 ��ȸ
	m_wmca.Query(
			GetSafeHwnd(),			//�� ������� ���� �޽����� �ްڽ��ϴ�.
			TRID_c1101,				//�� ���񽺿� ���ؼ� TRID_c1101(5) �ĺ��ڸ� ���̰ڴٴ� �ǹ��̸� �ݵ�� ����� �ʿ�� �����ϴ�.
			"c1101",				//ȣ���Ϸ��� ���� �ڵ�� 'c1101' �Դϴ�.
			(char*)&c1101inblock,	//c1101���� �䱸�ϴ� �Է� ����ü �����͸� �����մϴ�
			sizeof Tc1101InBlock 	//�Է� ����ü ũ���Դϴ�
									//���簡�� ������ �������� ��ȸ�� ���¹�ȣ�� �����ϹǷ� ���¹�ȣ �ε����� �������� �ʽ��ϴ�.
	);

	//������ ��û�� �ǽð� �ü��� �ִٸ� ��� ���
	m_wmca.DetachAll();

	//�ֽ� �ǽð� ���簡 ��Ŷ ��û ����
	m_wmca.Attach(
			GetSafeHwnd(),			//�ǽð� ������ ���� �޽����� �� ������� �ްڽ��ϴ�.
			"j8",					//��û�ϴ� �ǽð� ��Ŷ�� '�ڽ����ֽ����簡'(j8)�Դϴ�
			strJCode,				//j8 ���񽺿��� �䱸�ϴ� �����ڵ��Է°��� �����մϴ�.
			6/*�����ڵ� ũ��*/,	//���������ڵ��� ���̴� 6����Ʈ�̸�
			6/*�Է°� ��ü ũ��*/	//��ü�����ڵ��� ���̵� 6����Ʈ�Դϴ�. ���� ������ ��쿡�� (�����ڵ����*�����)�� �����մϴ�.
	);			

	//�ֽ� �ǽð� ���簡 ��Ŷ ��� ����
	//���ʿ��� �ǽð� ��Ŷ�� ������� ���� ��� �ش� ��Ŷ�� ��� ���ŵǹǷ� �����Ͻñ� �ٶ��ϴ�!
	//
	//m_wmca.Detach(
	//	GetSafeHwnd(),			//�ǽð� ������ ���� �޽����� �� ������� �ްڽ��ϴ�.
	//	"j8",					//����ϴ� �ǽð� ��Ŷ�� '�ڽ����ֽ����簡'(j8)�Դϴ�
	//	strJCode,				//j8 ���񽺿��� �䱸�ϴ� �����ڵ��Է°��� �����մϴ�.
	//	6/*�����ڵ� ũ��*/,	//���������ڵ��� ���̴� 6����Ʈ�̸�
	//	6/*�Է°� ��ü ũ��*/	//��ü�����ڵ��� ���̵� 6����Ʈ�Դϴ�.
	//	);			
	
	//3����(000660,005940,005930)�� ���� �ǽð� ���簡 ������ ��û ����
	//m_wmca.Attach(GetSafeHwnd(),"j8","000660005940005930",	6/*�����ڵ� ����*/,18/*=�����ڵ����(6)*�����(3)*/);		
	
	//���� ���簡, ȣ�� �ǽð� ������ ��û
	//m_wmca.Attach(GetSafeHwnd(),"f8","1206",	4,4);	//2012�� 6���� ���� ���簡
	//m_wmca.Attach(GetSafeHwnd(),"f1","1206",	4,4);	//2012�� 6���� ���� ȣ��
	
	//�ɼ� ���簡, ȣ�� �ǽð� ������ ��û
	//m_wmca.Attach(GetSafeHwnd(),"o2","20906157",8,8);	//2009�� 06���� �ݿɼ� ��簡 157 ü�ᰡ
	//m_wmca.Attach(GetSafeHwnd(),"o1","20906157",8,8);	//2009�� 06���� �ݿɼ� ��簡 157 ȣ��
}

void CWMCALOADERDlg::OnButtonClear() 
{
	m_listboxTrace.ResetContent();
}

void CWMCALOADERDlg::OnButtonOrder()
{
/*
	AfxMessageBox(	"�ֹ� �Ǽ� ������ ���� ������ �����Ͽ����ϴ�.\n\n�ҽ� �ڵ带 Ȯ���Ͻ� �� ������ �����Ͽ� ����Ͻñ� �ٶ��ϴ�",MB_ICONSTOP);
	return;
*/

	//���¹�ȣ �ε����� ���մϴ�.	(���¹�ȣ�� '1'������ ����)
	int	nAccountIndex	=m_comboAccountList.GetCurSel()+1;

	//�ֹ� test�� ���ǰ� �䱸�˴ϴ�
	//�ֽ��ֹ� ���񽺿��� �䱸�ϴ� �Է°��� ������ ����ü �����Դϴ�.
	Tc8102InBlock	c8102inblock	={0};
	memset(&c8102inblock,0x20,sizeof Tc8102InBlock);

	//ȭ����� ���� ��й�ȣ�� ���մϴ�
	CString	strPassword;
	GetDlgItem(IDC_EDIT_PASSWORD)->GetWindowText(strPassword);

	//�����ڵ� Ȯ��
	CString	strItemCode;
	GetDlgItem(IDC_EDIT_JCODE)->GetWindowText(strItemCode);
	strItemCode.TrimRight();

	//���� Ȯ��
	CString	strAmount;
	GetDlgItem(IDC_EDIT_AMOUNT)->GetWindowText(strAmount);
	strAmount.TrimRight();

	//�ܰ� Ȯ��
	CString	strPrice;
	GetDlgItem(IDC_EDIT_PRICE)->GetWindowText(strPrice);
	strPrice.TrimRight();

	//�� �Է� �ʵ忡�� �䱸�ϴ� ���鿡 ���� ���Ǵ� *.doc ������ ���� Ȯ���� �� �ֽ��ϴ�.
	//���º�й�ȣ�� ��������ȯ�濡���� Ȯ������ ������ ��ȯ�濡���� Ȯ���ϹǷ� ��Ȯ�ϰ� �Է��Ͻñ� �ٶ��ϴ�
	//���º�й�ȣ�� �ŷ�(�ֹ�)��й�ȣ�ʹ� �ٸ��ϴ�. �ŷ���й�ȣ�� ���� �ʵ��� �����Ͻñ� �ٶ��ϴ�.
	m_wmca.SetAccountIndexPwd(c8102inblock.pswd_noz8,nAccountIndex,strPassword);	

	SMOVE(c8102inblock.issue_codez6, strItemCode);		//�����ڵ�
	NMOVE(c8102inblock.order_qtyz12, strAmount);		//����
	NMOVE(c8102inblock.order_unit_pricez10,strPrice);	//����
	SMOVE(c8102inblock.trade_typez2, "00");				//������(00),���尡(03)

	//�Ʒ� �Է��ϴ� �ŷ���й�ȣ1, 2������ ������� �ŷ���й�ȣ�� �Է��ؾ��մϴ�.
	//�ŷ�(�ֹ�)��й�ȣ�� ���º�й�ȣ�ʹ� �ٸ��� ���º�й�ȣ�� ���� �ʵ��� �����Ͻñ� �ٶ��ϴ�.
	m_wmca.SetOrderPwd(c8102inblock.trad_pswd_no_1z8,"--------");		
	m_wmca.SetOrderPwd(c8102inblock.trad_pswd_no_2z8,"--------");		
	
	//�ֽ� �ֹ�
	m_wmca.Query(
			GetSafeHwnd(),			//�� ������� ���� �޽����� �ްڽ��ϴ�.
			TRID_c8102,				//�� ���񽺿� ���ؼ� TRID_c8102 �ĺ��ڸ� ���̰ڽ��ϴ�. (������ ���̸�,�ٸ� ���� �־ ����ϼŵ� �˴ϴ�)
			"c8102",				//ȣ���Ϸ��� ���� �ڵ�� c8102 �Դϴ�.
			(char*)&c8102inblock,	//c8102���� �䱸�ϴ� �Է� ����ü �����͸� �����մϴ�.
			sizeof Tc8102InBlock,	//�Է� ����ü ũ���Դϴ�
			nAccountIndex			//���ϴ� ���¹�ȣ �ε����� �����մϴ�.		(���¹�ȣ�� '1'������ ����)
	);

}

//������������������������������������������������������������������������������������������������
BOOL CWMCALOADERDlg::OnInitDialog()
{
	CDialog::OnInitDialog();
	
	SetIcon(m_hIcon, TRUE);			// Set big icon
	SetIcon(m_hIcon, FALSE);		// Set small icon

	//
	GetDlgItem(IDC_EDIT_JCODE)->SetWindowText("003160");

	GetDlgItem(IDC_EDIT1)->SetWindowText("");
	GetDlgItem(IDC_EDIT2)->SetWindowText("");
	GetDlgItem(IDC_EDIT3)->SetWindowText("");
	//////////////////////////////////////////////////////////////////////////
	
	return TRUE;  // return TRUE  unless you set the focus to a control
}

