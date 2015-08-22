// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs D:\Go\src\github.com\ghts\ghts\api\NH\wmca_ctype_orig.go

package NH

type AccountInfo struct {
	AccountNo	[11]int8
	AccountName	[40]int8
	Act_pdt_cdz3	[3]int8
	Amn_tab_cdz4	[4]int8
	ExpirationDate8	[8]int8
	Granted		int8
	Filler		[189]int8
}
type LoginInfo struct {
	Date		[14]int8
	ServerName	[15]int8
	UserID		[8]int8
	AccountCount	[3]int8
	Accountlist	[999]AccountInfo
}
type LoginBlock struct {
	TrIndex		int32
	LoginInfo	*LoginInfo
}

type MsgHeader struct {
	MsgCode	[5]int8
	UsrMsg	[80]int8
}

type Received struct {
	BlockName	*int8
	DataString	*int8
	Length		int32
}
type OutDataBlock struct {
	TrIndex		int32
	DataStruct	*Received
}

type Tc1101OutBlock struct {
	Code			[6]int8
	X_code			int8
	Hname			[13]int8
	X_hname			int8
	Price			[7]int8
	X_price			int8
	Sign			[1]int8
	X_sign			int8
	Change			[6]int8
	X_change		int8
	Chrate			[5]int8
	X_chrate		int8
	Offer			[7]int8
	X_offer			int8
	Bid			[7]int8
	X_bid			int8
	Volume			[9]int8
	X_volume		int8
	Volrate			[6]int8
	X_volrate		int8
	Yurate			[5]int8
	X_yurate		int8
	Value			[9]int8
	X_value			int8
	Uplmtprice		[7]int8
	X_uplmtprice		int8
	High			[7]int8
	X_high			int8
	Open			[7]int8
	X_open			int8
	Opensign		[1]int8
	X_opensign		int8
	Openchange		[6]int8
	X_openchange		int8
	Low			[7]int8
	X_low			int8
	Dnlmtprice		[7]int8
	X_dnlmtprice		int8
	Hotime			[8]int8
	X_hotime		int8
	Offerho			[7]int8
	X_offerho		int8
	P_offer			[7]int8
	X_P_offer		int8
	S_offer			[7]int8
	X_S_offer		int8
	S4_offer		[7]int8
	X_S4_offer		int8
	S5_offer		[7]int8
	X_S5_offer		int8
	S6_offer		[7]int8
	X_S6_offer		int8
	S7_offer		[7]int8
	X_S7_offer		int8
	S8_offer		[7]int8
	X_S8_offer		int8
	S9_offer		[7]int8
	X_S9_offer		int8
	S10_offer		[7]int8
	X_S10_offer		int8
	Bidho			[7]int8
	X_bidho			int8
	P_bid			[7]int8
	X_P_bid			int8
	S_bid			[7]int8
	X_S_bid			int8
	S4_bid			[7]int8
	X_S4_bid		int8
	S5_bid			[7]int8
	X_S5_bid		int8
	S6_bid			[7]int8
	X_S6_bid		int8
	S7_bid			[7]int8
	X_S7_bid		int8
	S8_bid			[7]int8
	X_S8_bid		int8
	S9_bid			[7]int8
	X_S9_bid		int8
	S10_bid			[7]int8
	X_S10_bid		int8
	Offerrem		[9]int8
	X_offerrem		int8
	P_offerrem		[9]int8
	X_P_offerrem		int8
	S_offerrem		[9]int8
	X_S_offerrem		int8
	S4_offerrem		[9]int8
	X_S4_offerrem		int8
	S5_offerrem		[9]int8
	X_S5_offerrem		int8
	S6_offerrem		[9]int8
	X_S6_offerrem		int8
	S7_offerrem		[9]int8
	X_S7_offerrem		int8
	S8_offerrem		[9]int8
	X_S8_offerrem		int8
	S9_offerrem		[9]int8
	X_S9_offerrem		int8
	S10_offerrem		[9]int8
	X_S10_offerrem		int8
	Bidrem			[9]int8
	X_bidrem		int8
	P_bidrem		[9]int8
	X_P_bidrem		int8
	S_bidrem		[9]int8
	X_S_bidrem		int8
	S4_bidrem		[9]int8
	X_S4_bidrem		int8
	S5_bidrem		[9]int8
	X_S5_bidrem		int8
	S6_bidrem		[9]int8
	X_S6_bidrem		int8
	S7_bidrem		[9]int8
	X_S7_bidrem		int8
	S8_bidrem		[9]int8
	X_S8_bidrem		int8
	S9_bidrem		[9]int8
	X_S9_bidrem		int8
	S10_bidrem		[9]int8
	X_S10_bidrem		int8
	T_offerrem		[9]int8
	X_T_offerrem		int8
	T_bidrem		[9]int8
	X_T_bidrem		int8
	O_offerrem		[9]int8
	X_O_offerrem		int8
	O_bidrem		[9]int8
	X_O_bidrem		int8
	Pivot2upz7		[7]int8
	X_pivot2upz7		int8
	Pivot1upz7		[7]int8
	X_pivot1upz7		int8
	Pivotz7			[7]int8
	X_pivotz7		int8
	Pivot1dnz7		[7]int8
	X_pivot1dnz7		int8
	Pivot2dnz7		[7]int8
	X_pivot2dnz7		int8
	Sosokz6			[6]int8
	X_sosokz6		int8
	Jisunamez18		[18]int8
	X_jisunamez18		int8
	Capsizez6		[6]int8
	X_capsizez6		int8
	Output1z16		[16]int8
	X_output1z16		int8
	Marcket1z16		[16]int8
	X_marcket1z16		int8
	Marcket2z16		[16]int8
	X_marcket2z16		int8
	Marcket3z16		[16]int8
	X_marcket3z16		int8
	Marcket4z16		[16]int8
	X_marcket4z16		int8
	Marcket5z16		[16]int8
	X_marcket5z16		int8
	Marcket6z16		[16]int8
	X_marcket6z16		int8
	Cbtext			[6]int8
	X_cbtext		int8
	Parvalue		[7]int8
	X_parvalue		int8
	Prepricetitlez12	[12]int8
	X_prepricetitlez12	int8
	Prepricez7		[7]int8
	X_prepricez7		int8
	Subprice		[7]int8
	X_subprice		int8
	Gongpricez7		[7]int8
	X_gongpricez7		int8
	High5			[7]int8
	X_high5			int8
	Low5			[7]int8
	X_low5			int8
	High20			[7]int8
	X_high20		int8
	Low20			[7]int8
	X_low20			int8
	Yhigh			[7]int8
	X_yhigh			int8
	Yhighdate		[4]int8
	X_yhighdate		int8
	Ylow			[7]int8
	X_ylow			int8
	Ylowdate		[4]int8
	X_ylowdate		int8
	Movlistingz8		[8]int8
	X_movlistingz8		int8
	Listing			[12]int8
	X_listing		int8
	Totpricez9		[9]int8
	X_totpricez9		int8
	Tratimez5		[5]int8
	X_tratimez5		int8
	Off_tra1		[6]int8
	X_off_tra1		int8
	Bid_tra1		[6]int8
	X_bid_tra1		int8
	N_offvolume1		[9]int8
	X_N_offvolume1		int8
	N_bidvolume1		[9]int8
	X_N_bidvolume1		int8
	Off_tra2		[6]int8
	X_off_tra2		int8
	Bid_tra2		[6]int8
	X_bid_tra2		int8
	N_offvolume2		[9]int8
	X_N_offvolume2		int8
	N_bidvolume2		[9]int8
	X_N_bidvolume2		int8
	Off_tra3		[6]int8
	X_off_tra3		int8
	Bid_tra3		[6]int8
	X_bid_tra3		int8
	N_offvolume3		[9]int8
	X_N_offvolume3		int8
	N_bidvolume3		[9]int8
	X_N_bidvolume3		int8
	Off_tra4		[6]int8
	X_off_tra4		int8
	Bid_tra4		[6]int8
	X_bid_tra4		int8
	N_offvolume4		[9]int8
	X_N_offvolume4		int8
	N_bidvolume4		[9]int8
	X_N_bidvolume4		int8
	Off_tra5		[6]int8
	X_off_tra5		int8
	Bid_tra5		[6]int8
	X_bid_tra5		int8
	N_offvolume5		[9]int8
	X_N_offvolume5		int8
	N_bidvolume5		[9]int8
	X_N_bidvolume5		int8
	N_offvolall		[9]int8
	X_N_offvolall		int8
	N_bidvolall		[9]int8
	X_N_bidvolall		int8
	Fortimez6		[6]int8
	X_fortimez6		int8
	Forratez5		[5]int8
	X_forratez5		int8
	Settdatez4		[4]int8
	X_settdatez4		int8
	Cratez5			[5]int8
	X_cratez5		int8
	Yudatez4		[4]int8
	X_yudatez4		int8
	Mudatez4		[4]int8
	X_mudatez4		int8
	Yuratez5		[5]int8
	X_yuratez5		int8
	Muratez5		[5]int8
	X_muratez5		int8
	Formovolz10		[10]int8
	X_formovolz10		int8
	Jasa			[1]int8
	X_jasa			int8
	Listdatez8		[8]int8
	X_listdatez8		int8
	Daeratez5		[5]int8
	X_daeratez5		int8
	Daedatez6		[6]int8
	X_daedatez6		int8
	Clovergb		[1]int8
	X_clovergb		int8
	Depositgb		[1]int8
	X_depositgb		int8
	Capital			[9]int8
	X_capital		int8
	N_alloffvol		[9]int8
	X_N_alloffvol		int8
	N_allbidvol		[9]int8
	X_N_allbidvol		int8
	Hnamez21		[21]int8
	X_hnamez21		int8
	Detourgb		[1]int8
	X_detourgb		int8
	Yuratez6		[6]int8
	X_yuratez6		int8
	Sosokz6_1		[6]int8
	X_sosokz6_1		int8
	Maedatez4		[4]int8
	X_maedatez4		int8
	Lratez5			[5]int8
	X_lratez5		int8
	Perz5			[5]int8
	X_perz5			int8
	Handogb			[1]int8
	X_handogb		int8
	Avgprice		[7]int8
	X_avgprice		int8
	Listing2		[12]int8
	X_listing2		int8
	Addlisting		[12]int8
	X_addlisting		int8
	Gicomment		[100]int8
	X_gicomment		int8
	Prevolume		[9]int8
	X_prevolume		int8
	Presign			[1]int8
	X_presign		int8
	Prechange		[6]int8
	X_prechange		int8
	Yhigh2			[7]int8
	X_yhigh2		int8
	Yhighdate2		[4]int8
	X_yhighdate2		int8
	Ylow2			[7]int8
	X_ylow2			int8
	Ylowdate2		[4]int8
	X_ylowdate2		int8
	Forstock		[15]int8
	X_forstock		int8
	Forlmtz5		[5]int8
	X_forlmtz5		int8
	Maeunit			[5]int8
	X_maeunit		int8
	Mass_opt		[1]int8
	X_mass_opt		int8
	Largemgb		[1]int8
	X_largemgb		int8
}
type Tc1101OutBlock2 struct {
	Time		[8]int8
	X_time		int8
	Price		[7]int8
	X_price		int8
	Sign		[1]int8
	X_sign		int8
	Change		[6]int8
	X_change	int8
	Offer		[7]int8
	X_offer		int8
	Bid		[7]int8
	X_bid		int8
	Movolume	[8]int8
	X_movolume	int8
	Volume		[9]int8
	X_volume	int8
}
