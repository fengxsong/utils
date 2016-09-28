```
import (
    "fmt"
    "time"
    "github.com/fengxsong/utils/token"
)

func main() {
    t := token.NewToken("HelloWorld", "a test message", time.Now().Add(60*time.Second).Unix())
    ss, err := t.SigningString()
    if err != nil {
        fmt.Println(err)
        return
    }

    time.Sleep(61 * time.Second)
    fmt.Println(token.Verify(ss))
}



```
