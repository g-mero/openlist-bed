package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestViper(t *testing.T) {
	err := os.Setenv("MYAPP_SERVER_PORT_S", "9000")
	if err != nil {
		panic(err)
	}

	// prefix (optional)
	viper.SetEnvPrefix("MYAPP")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	port := viper.GetInt("server.port_s")
	if port == 0 {
		fmt.Println("未设置环境变量，使用默认端口: 8080")
		port = 8080
	}

	fmt.Printf("服务端口: %d\n", port)

}
