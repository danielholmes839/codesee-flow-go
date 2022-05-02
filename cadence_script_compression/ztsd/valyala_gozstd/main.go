//package main
//
//import (
//	"fmt"
//	"log"
//
//	csc "github.com/onflow/flow-go/cadence_script_compression"
//	ztsd "github.com/valyala/gozstd"
//)
//
//const (
//	mainnetDir = "../../contracts/mainnet"
//)
//
//func main() {
//	contracts := csc.ReadContracts(mainnetDir)
//
//	sumOfRatios := float64(0)
//	for _, c := range contracts {
//		compData := &csc.CompressionComparison{
//			CompressedData:   make([]byte, 0),
//			UncompressedData: c.Data,
//		}
//
//		compData.CompressedData = ztsd.Compress(compData.CompressedData, c.Data)
//
//		sumOfRatios = sumOfRatios + compData.CompressionRatio()
//		log.Println(fmt.Sprintf("Name: %s, Uncompressed: %d, Compressed: %d Ratio: %f", c.Name, compData.UnCompressedSize(), compData.CompressedSize(), compData.CompressionRatio()))
//	}
//
//	log.Println(fmt.Sprintf("Average compression Ratio: %f", sumOfRatios/float64(len(contracts))))
//}
