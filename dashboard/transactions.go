package dashboard

import (
	"encoding/json"
	"github.com/adonese/noebs/ebs_fields"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	ebs_fields.GenericEBSResponseFields
}

type PurchaseModel struct {
	gorm.Model
	ebs_fields.PurchaseFields
}
type Env struct {
	Db *gorm.DB
}

func (e *Env) GetTransactionbyID(c *gin.Context) {
	var tran Transaction
	//id := c.Params.ByName("id")
	err := e.Db.Find(&tran).Error
	if err != nil {
		c.AbortWithStatus(404)
	}
	c.JSON(200, gin.H{"result": tran.ID})

	defer e.Db.Close()
}

type MerchantTransactions struct {
	PurchaseNumber         *int     `json:"purchase_number"`
	PurchaseAmount         *float32 `json:"purchase_amount"`
	AllTransactions        *int     `json:"all_transactions"`
	SuccessfulTransactions *int     `json:"successful_transactions"`
	FailedTransactions     *int     `json:"failed_transactions"`
}

// To allow Redis to use this struct directly in marshaling
func (p *MerchantTransactions) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)

}

// To allow Redis to use this struct directly in marshaling
func (p *MerchantTransactions) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func purchaseSum(tran []string) float32 {
	var trans []MerchantTransactions
	var mtran MerchantTransactions
	for _, k := range tran {
		json.Unmarshal([]byte(k), &mtran)
		trans = append(trans, mtran)
	}
	var sum float32
	for _, k := range trans {
		sum += *k.PurchaseAmount
	}
	return sum
}

func ToPurchase(f ebs_fields.PurchaseFields) *MerchantTransactions {
	amount := f.TranAmount
	m := new(MerchantTransactions)
	*m.PurchaseAmount = amount
	return m
}
