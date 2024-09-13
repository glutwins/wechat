package store

import (
	"context"
	"fmt"

	"github.com/silenceper/wechat/v2/util"
)

type EcGetAftersaleListRequest struct {
	BeginCreateTime int64  `json:"begin_create_time"`
	EndCreateTime   int64  `json:"end_create_time"`
	NextKey         string `json:"next_key"`
}

type EcGetAftersaleListResp struct {
	util.CommonError
	AfterSaleOrderIdList []string `json:"after_sale_order_id_list"`
	HasMore              bool     `json:"has_more"`
	NextKey              string   `json:"next_key"`
}

func (s *Store) EcGetAftersaleList(c context.Context, req *EcGetAftersaleListRequest) (*EcGetAftersaleListResp, error) {
	accessToken, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s?access_token=%s", "https://api.weixin.qq.com/channels/ec/aftersale/getaftersalelist", accessToken)
	response, err := util.PostJSONContext(c, uri, req)
	if err != nil {
		return nil, err
	}
	resp := &EcGetAftersaleListResp{}
	err = util.DecodeWithError(response, resp, "EcGetAftersaleList")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type ProductInfo struct {
	ProductID string `json:"product_id"`
	SkuID     string `json:"sku_id"`
	Count     int    `json:"count"`
}
type Details struct {
	Desc           string        `json:"desc"`
	ReceiveProduct bool          `json:"receive_product"`
	CancelTime     int           `json:"cancel_time"`
	MediaIDList    []interface{} `json:"media id list"`
	TelNumber      string        `json:"tel_number"`
}
type RefundInfo struct {
	Amount       int `json:"amount"`
	RefundReason int `json:"refund_reason"`
}
type ReturnInfo struct {
	WaybillID    string `json:"waybill_id"`
	DeliveryID   string `json:"delivery_id"`
	DeliveryName string `json:"delivery_name"`
}
type MerchantUploadInfo struct {
	RejectReason       string        `json:"reject_reason"`
	RefundCertificates []interface{} `json:"refund_certificates"`
}
type RefundResp struct {
	Code    string `json:"code"`
	Ret     int    `json:"ret"`
	Message string `json:"message"`
}
type AfterSaleOrder struct {
	AfterSaleOrderID   string              `json:"after_sale_order_id"`
	Status             string              `json:"status"`
	Openid             string              `json:"openid"`
	UnionId            string              `json:"unionid"`
	OrderID            string              `json:"order_id"`
	ProductInfo        ProductInfo         `json:"product_info"`
	Details            Details             `json:"details"`
	RefundInfo         *RefundInfo         `json:"refund_info"`
	ReturnInfo         *ReturnInfo         `json:"return_info"`
	MerchantUploadInfo *MerchantUploadInfo `json:"merchant_upload_info"`
	CreateTime         int                 `json:"create_time"`
	UpdateTime         int                 `json:"update_time"`
	Reason             string              `json:"reason"`
	RefundResp         *RefundResp         `json:"refund_resp"`
	Type               string              `json:"type"`
}

type GetaftersaleorderResp struct {
	util.CommonError
	AfterSaleOrder AfterSaleOrder `json:"after_sale_order"`
}

func (s *Store) Getaftersaleorder(c context.Context, orderId string) (*GetaftersaleorderResp, error) {
	accessToken, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s?access_token=%s", "https://api.weixin.qq.com/channels/ec/aftersale/getaftersaleorder", accessToken)
	response, err := util.PostJSONContext(c, uri, map[string]interface{}{
		"after_sale_order_id": orderId,
	})
	if err != nil {
		return nil, err
	}
	resp := &GetaftersaleorderResp{}
	err = util.DecodeWithError(response, resp, "Getaftersaleorder")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
