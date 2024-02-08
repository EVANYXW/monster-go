package sm2

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"testing"
)

func TestSm2Sign(t *testing.T) {
	file, err := os.OpenFile("C:/psbc/test.sm2", os.O_RDONLY, 0600)
	fmt.Println(err)
	fileInfo, err := file.Stat()
	if err !=nil{
		panic(err)
	}
	buf :=make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	if err !=nil{
		panic(err)
	}

	curve := P256Sm2()
	curve.ScalarBaseMult(buf)
	da := new(PrivateKey)
	da.Curve = curve
	da.D = new(big.Int).SetBytes(buf)
	da.PublicKey.X, da.PublicKey.Y = curve.ScalarBaseMult(buf)
	//data :=[]byte("3112 1652172683256 20220510165123")
	//res,err :=da.Sign(rand.Reader, data, nil)
	//res2 ,err:=da.Encrypt(res)
	//fmt.Println(res2)
	//sign:=base64.StdEncoding.EncodeToString(res2)
	//fmt.Println(sign,err)
	sign:="MEQCICDtUq2ole1zxANEJ6lYSgNnFckRcSAlbfdJJgYuwWuFAiA7YR+3wwpUBvTxCQA/jbigLYaiBW06UftHFTijNFxS0w=="
	res2,_:=base64.StdEncoding.DecodeString(sign)
	status:=da.PublicKey.Verify([]byte("3112 1652173397007 20220510170317"),res2)
	fmt.Println(status)

}
