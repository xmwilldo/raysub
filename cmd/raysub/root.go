/*
Copyright © 2020 Will

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	base64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xmwilldo/v2ray-sub/cmd/raysub/config"
	"github.com/xmwilldo/v2ray-sub/cmd/raysub/validate"
	"github.com/xmwilldo/v2ray-sub/version"
)

var (
	cfgFile string

	description = ``

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:     "raysub",
		Short:   "update v2ray's config by subscription url",
		Long:    description,
		RunE:    validate.SubCommandExists,
		Version: version.Version,
	}
)

const (
	DefautConfigFile  = "/etc/raysub/config.yml"
	DefaultConfigPath = "/etc/raysub"
)

func init() {
	cobra.OnInitialize(
		initConfig,
	)

	//rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "raysub config path.[default: /etc/raysub/config.yml]")
	//rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "subscription url")
	//rootCmd.MarkPersistentFlagRequired("url")
	//rootCmd.Flags().StringVarP(&configPath, "config", "c", "/etc/v2ray/config.json", "v2ray config path")
}

func initConfig() {
	if cfgFile == "" {
		cfgFile = DefautConfigFile
	}

	configPath := DefaultConfigPath
	if isExist := exists(configPath); !isExist {
		if err := os.Mkdir(configPath, 0644); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	if isExist := exists(cfgFile); !isExist {
		if _, err := os.Create(DefautConfigFile); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func runE(cmd *cobra.Command, args []string) error {
	log.Println("get all proxy configs by subscription url...")
	subscriptionUrl := viper.GetString("subscriptionUrl")
	v2rayConfigPath := viper.GetString("v2rayConfigPath")
	log.Println("subscriptionUrl:", subscriptionUrl, "v2rayConfigPath:", v2rayConfigPath)
	proxyConfigs, err := getAllProxyConfigs(subscriptionUrl)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("select a fastest proxy config...")
	// todo get a fastest proxy
	wantProxyConfig := getFastestProxy(proxyConfigs)

	log.Println("modify v2ray config...")
	contentB, err := readFile(v2rayConfigPath)
	if err != nil {
		log.Println(err)
		return err
	}

	v2rayCoreConfig := config.V2rayCoreConfig{}
	if err := json.Unmarshal(contentB, &v2rayCoreConfig); err != nil {
		log.Println(err)
		return err
	}

	// modify the v2ray config
	for index, outbound := range v2rayCoreConfig.Outbounds {
		if outbound.Tag == "proxy" {
			alterId, err := strconv.Atoi(wantProxyConfig.Aid)
			if err != nil {
				log.Println(err)
				return err
			}

			port, err := strconv.Atoi(wantProxyConfig.Port)
			if err != nil {
				log.Println(err)
				return err
			}

			v2rayCoreConfig.Outbounds[index].Settings.Vnext[0].Address = wantProxyConfig.Add
			v2rayCoreConfig.Outbounds[index].Settings.Vnext[0].Users[0].ID = wantProxyConfig.ID
			v2rayCoreConfig.Outbounds[index].Settings.Vnext[0].Users[0].AlterID = alterId
			v2rayCoreConfig.Outbounds[index].Settings.Vnext[0].Port = port
		} else {
			continue
		}
	}

	outputContent, err := json.Marshal(v2rayCoreConfig)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := writeToFile(v2rayConfigPath, outputContent); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("restart v2ray service...")
		// restart v2ray
		//code, err := restartService()
		//if err != nil {
		//	log.Printf("restart v2ray service err: %v, error code: %d\n", err, code)
		//	return err
		//}
	}

	log.Println("done.")
	return nil
}

func getFastestProxy(configs []config.ProxyConfig) config.ProxyConfig {
	return configs[0]
}

func getAllProxyConfigs(url string) ([]config.ProxyConfig, error) {
	proxyConfigs := make([]config.ProxyConfig, 0)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return proxyConfigs, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return proxyConfigs, err
	}

	respContentB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return proxyConfigs, err
	}

	proxyConfigsEncodeB, err := base64.StdEncoding.DecodeString(string(respContentB))
	if err != nil {
		return proxyConfigs, err
	}

	proxyConfigsStrs := strings.Split(string(proxyConfigsEncodeB), "vmess://")
	// delete header and tail
	proxyConfigsStrs = proxyConfigsStrs[1 : len(proxyConfigsStrs)-1]

	proxyConfig := config.ProxyConfig{}
	for _, configStr := range proxyConfigsStrs {
		proxyConfigB, err := base64.StdEncoding.DecodeString(configStr)
		if err != nil {
			return proxyConfigs, err
		}

		if err := json.Unmarshal(proxyConfigB, &proxyConfig); err != nil {
			return proxyConfigs, err
		}

		proxyConfigs = append(proxyConfigs, proxyConfig)
	}

	return proxyConfigs, nil
}

func readFile(filePath string) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	output := make([]byte, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return output, nil
			}
			return nil, err
		}
		output = append(output, line...)
		output = append(output, []byte("\n")...)
	}
	return output, nil
}

func writeToFile(filePath string, outPut []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func restartService() (int, error) {
	cmd := exec.Command("/bin/bash", "-c", "systemctl restart v2ray") //不加第一个第二个参数会报错

	//cmd.Stdout = os.Stdout // cmd.Stdout -> stdout  重定向到标准输出，逐行实时打印
	//cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
	//也可以重定向文件 cmd.Stderr= fd (文件打开的描述符即可)

	stdout, _ := cmd.StdoutPipe() //创建输出管道
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	//cmdPid := cmd.Process.Pid //查看命令pid

	//result, _ := ioutil.ReadAll(stdout) // 读取输出结果
	//resdata := string(result)

	var res int
	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			log.Println("cmd exit status")
			res = ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
		}
	}

	return res, nil
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
