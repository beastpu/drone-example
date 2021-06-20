
package hello

import (
        "testing"
    )
func TestHello(t *testing.T){
   ret:=SayHello()
   if ret!="helloworld"{
      t.Error("not return helloworld")
   }else {
            t.Log("第一个测试通过了")
        }
}
