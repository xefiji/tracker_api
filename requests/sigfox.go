package requests

type Sigfox struct {
	Data         string `json:"data" validate:"required"`
	Time         string `json:"time" validate:"required"`
	SeqNumber    string `json:"seqNumber" validate:"gte=0"`
	Rssi         string `json:"rssi"`
	DeviceTypeId string `json:"deviceTypeId"`
	Id           string `json:"id"`
	Snr          string `json:"snr"`
	Station      string `json:"station"`
}
