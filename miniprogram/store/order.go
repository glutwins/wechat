package store

import (
	"context"
	"fmt"

	"github.com/silenceper/wechat/v2/util"
)

const (
	OrderStatusWaitPay     = 10
	OrderStatusWaitDeliver = 20
	OrderStatusPartDeliver = 21
	OrderStatusWaitReceive = 30
	OrderStatusFinish      = 100
	OrderStatusCancel      = 250
)

type TimeRange struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type EcOrderListGetRequest struct {
	CreateTimeRange *TimeRange `json:"create_time_range,omitempty"`
	UpdateTimeRange *TimeRange `json:"update_time_range,omitempty"`
	Status          int        `json:"status,omitempty"`
	OpenId          string     `json:"openid,omitempty"`
	NextKey         string     `json:"next_key,omitempty"`
	PageSize        int        `json:"page_size"`
}

type EcOrderListGetResp struct {
	util.CommonError
	OrderIdList []string `json:"order_id_list"`
	NextKey     string   `json:"next_key,omitempty"`
	HasMore     bool     `json:"has_more"`
}

func (s *Store) EcOrderListGet(c context.Context, req *EcOrderListGetRequest) (*EcOrderListGetResp, error) {
	accessToken, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s?access_token=%s", "https://api.weixin.qq.com/channels/ec/order/list/get", accessToken)
	response, err := util.PostJSONContext(c, uri, req)
	if err != nil {
		return nil, err
	}
	resp := &EcOrderListGetResp{}
	err = util.DecodeWithError(response, resp, "EcOrderListGet")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type EcOrderGetResp struct {
	util.CommonError
	Order Order `json:"order"`
}
type SkuAttrs struct {
	AttrKey   string `json:"attr_key"`
	AttrValue string `json:"attr_value"`
}
type ProductInfos struct {
	ProductID             string     `json:"product_id"`
	SkuID                 string     `json:"sku_id"`
	SkuCnt                int        `json:"sku_cnt"`
	OnAftersaleSkuCnt     int        `json:"on_aftersale_sku_cnt"`
	FinishAftersaleSkuCnt int        `json:"finish_aftersale_sku_cnt"`
	Title                 string     `json:"title"`
	ThumbImg              string     `json:"thumb_img"`
	SalePrice             int        `json:"sale_price"`
	MarketPrice           int        `json:"market_price"`
	SkuAttrs              []SkuAttrs `json:"sku_attrs"`
}
type PayInfo struct {
	PrepayID      string `json:"prepay_id"`
	TransactionID string `json:"transaction_id"`
	PrepayTime    int    `json:"prepay_time"`
	PayTime       int    `json:"pay_time"`
	PaymentMethod int    `json:"payment_method"`
}
type PriceInfo struct {
	ProductPrice            int  `json:"product_price,omitempty"`
	OrderPrice              int  `json:"order_price,omitempty"`
	Freight                 int  `json:"freight,omitempty"`
	DiscountedPrice         int  `json:"discounted_price,omitempty"`
	IsDiscounted            bool `json:"is_discounted,omitempty"`
	OriginalOrderPrice      int  `json:"original_order_price,omitempty"`
	EstimateProductPrice    int  `json:"estimate_product_price,omitempty"`
	ChangeDownPrice         int  `json:"change_down_price,omitempty"`
	ChangeFreight           int  `json:"change_freight,omitempty"`
	IsChangeFreight         bool `json:"is_change_freight,omitempty"`
	UseDeduction            bool `json:"use_deduction,omitempty"`
	DeductionPrice          int  `json:"deduction_price,omitempty"`
	MerchantReceievePrice   int  `json:"merchant_receieve_price,omitempty"`
	MerchantDiscountedPrice int  `json:"merchant_discounted_price,omitempty"`
	FinderDiscountedPrice   int  `json:"finder_discounted_price,omitempty"`
}
type AddressInfo struct {
	UserName     string `json:"user_name"`
	PostalCode   string `json:"postal_code"`
	ProvinceName string `json:"province_name"`
	CityName     string `json:"city_name"`
	CountyName   string `json:"county_name"`
	DetailInfo   string `json:"detail_info"`
	TelNumber    string `json:"tel_number"`
}
type DeliveryProduct struct {
	ProductID  string `json:"product_id"`
	SkuID      string `json:"sku_id"`
	ProductCnt int    `json:"product_cnt"`
}
type DeliveryProductInfo struct {
	WaybillID    string            `json:"waybill_id"`
	DeliveryID   string            `json:"delivery_id"`
	DeliveryTime int               `json:"delivery_time"`
	DeliverType  int               `json:"deliver_type"`
	ProductInfos []DeliveryProduct `json:"product_infos"`
}
type DeliveryInfo struct {
	AddressInfo         AddressInfo           `json:"address_info"`
	DeliveryProductInfo []DeliveryProductInfo `json:"delivery_product_info"`
	ShipDoneTime        int                   `json:"ship_done_time"`
	DeliverMethod       int                   `json:"deliver_method"`
}
type CouponInfo struct {
	UserCouponID string `json:"user_coupon_id"`
}
type ExtInfo struct {
	CustomerNotes string `json:"customer_notes"`
	MerchantNotes string `json:"merchant_notes"`
	FinderID      string `json:"finder_id"`
	LiveID        string `json:"live_id"`
	OrderScene    int    `json:"order_scene"`
}
type SharerInfo struct {
	SharerOpenid     string `json:"sharer_openid"`
	SharerUnionid    string `json:"sharer_unionid"`
	SharerType       int    `json:"sharer_type"`
	ShareScene       int    `json:"share_scene"`
	HandlingProgress int    `json:"handling_progress"`
}
type SettleInfo struct {
	CommissionFee           int   `json:"commission_fee"`
	PredictCommissionFee    int   `json:"predict_commission_fee"`
	PredictWecoinCommission int   `json:"predict_wecoin_commission"`
	WecoinCommission        int   `json:"wecoin_commission"`
	SettleTime              int64 `json:"settle_time"`
}
type SkuSharerInfos struct {
	SharerOpenid  string `json:"sharer_openid"`
	SharerUnionid string `json:"sharer_unionid"`
	SharerType    int    `json:"sharer_type"`
	ShareScene    int    `json:"share_scene"`
	SkuID         string `json:"sku_id"`
}
type CommissionInfo struct {
	SkuID        string `json:"sku_id"`
	NickName     string `json:"nickname"`
	Type         int    `json:"type"`
	Status       int    `json:"status"`
	Amount       int    `json:"amount"`
	FinderId     string `json:"finderid"`
	OpenFinderId string `json:"openfinderid"`
}
type AgentInfo struct {
	AgentFinderId       string `json:"agent_finder_id"`
	AgentFinderNickname string `json:"agent_finder_nickname"`
}
type OrderDetail struct {
	ProductInfos    []ProductInfos   `json:"product_infos"`
	PayInfo         PayInfo          `json:"pay_info"`
	PriceInfo       PriceInfo        `json:"price_info"`
	DeliveryInfo    DeliveryInfo     `json:"delivery_info"`
	CouponInfo      CouponInfo       `json:"coupon_info"`
	ExtInfo         ExtInfo          `json:"ext_info"`
	CommissionInfos []CommissionInfo `json:"commission_infos"`
	SharerInfo      SharerInfo       `json:"sharer_info"`
	SettleInfo      SettleInfo       `json:"settle_info"`
	SkuSharerInfos  []SkuSharerInfos `json:"sku_sharer_infos"`
	AgentInfo       AgentInfo        `json:"agent_info"`
}
type AftersaleOrderList struct {
	AftersaleOrderID string `json:"aftersale_order_id"`
	Status           int    `json:"status"`
}
type AftersaleDetail struct {
	AftersaleOrderList  []AftersaleOrderList `json:"aftersale_order_list"`
	OnAftersaleOrderCnt int                  `json:"on_aftersale_order_cnt"`
}
type Order struct {
	OrderID         string          `json:"order_id"`
	Status          int             `json:"status"`
	CreateTime      int             `json:"create_time"`
	UpdateTime      int             `json:"update_time"`
	OrderDetail     OrderDetail     `json:"order_detail"`
	AftersaleDetail AftersaleDetail `json:"aftersale_detail"`
	Openid          string          `json:"openid"`
	UnionId         string          `json:"unionid"`
}

func (s *Store) EcOrderGet(c context.Context, orderId string) (*EcOrderGetResp, error) {
	accessToken, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s?access_token=%s", "https://api.weixin.qq.com/channels/ec/order/get", accessToken)
	response, err := util.PostJSONContext(c, uri, map[string]interface{}{
		"order_id":              orderId,
		"encode_sensitive_info": false,
	})
	if err != nil {
		return nil, err
	}
	resp := &EcOrderGetResp{}
	err = util.DecodeWithError(response, resp, "EcOrderGet")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
