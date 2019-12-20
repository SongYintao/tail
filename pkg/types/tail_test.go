package types

import (
	"fmt"
	"strconv"
	"testing"
)

func TestTailFile(t *testing.T) {
	config:=Config{
		Follow:      true,
		MaxLineSize: 100,
	}
	fileName:="/Users/nali/Documents/go/tail/pkg/types/tail.go"
	tf,err:=TailFile(fileName,config)
	if err!=nil{
		fmt.Println(err)
	}
	count:=0
	for line:=range tf.Lines{
		fmt.Println(line.Text)
		count++
		fmt.Println("==========="+strconv.Itoa(count)+"===========")
		if count>10{
			break

		}
	}

	tf.Done()
	tf.close()
	tf.Cleanup()
	//tf.close()
	//err=tf.Wait()
	//if err!=nil{
	//	fmt.Println(err)
	//}
}
