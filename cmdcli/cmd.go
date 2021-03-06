/*
@Time : 2018/7/27 15:36 
@Author : 星空之钥丶
@File : cmd
@Software: GoLang
*/
package cmdcli

import(
	cli "gopkg.in/urfave/cli.v1"
	"os"
	"github.com/bottos-project/bottos/config"
	"strconv"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
	"strings"
)
var Conf  config.Parameter
var GenConf config.GenesisConfig
var KeyPair config.KeyPair

func Init() (config.Parameter, error) {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config",
			Value:"./chainconfig.json",
			Usage: "config file path the greeting,If without this path, the bottos process will boot up with default config in hardcode",
		},
		cli.StringFlag{
			Name: "genesis",
			Usage: "genesis config file path the greeting",
		},
		cli.StringFlag{
			Name: "datadir",
			Usage: "datadir's path",
		},
		cli.StringFlag{
			Name: "disable-api",
			Value:"0",
			Usage: "disable restful api's requests",
		},
		cli.StringFlag{
			Name: "apiport",
			Value: "8090",
			Usage: "api service port for the greeting",
		},
		cli.StringFlag{
			Name: "disable-rpc",
			Usage: "disable rpc requests",
		},
		cli.StringFlag{
			Name: "rpcport",
			Value: "8080",
			Usage: "rpc port for the greeting",
		},
		cli.StringFlag{
			Name: "p2pport",
			Value: "8096",
			Usage: "local listen on this p2p port to receive remote p2p messages",
		},
		cli.StringFlag{
			Name: "servaddr",
			Usage: "for p2p sync / reply local server ip& port info",
		},
		cli.StringFlag{
			Name: "peerlist",
			Usage: "for p2p add pne / add neighbour. Example: 192.168.1.2:9868, 192.168.1.3:9868, 192.168.1.4:9868",
		},
		cli.StringFlag{
			Name: "delegate-signkey",
			Usage: "--delegate-signkey=<pubkey>,<private key>.Param struct needs be modified ,public and private key for native contract, external contracts' accounts",
		},
		cli.StringFlag{
			Name: "delegate",
			Usage: "Assign one producer. Later this section will no more be used.\n Only one delegate is allowed in one node(other than bottos account).",
		},
		cli.StringFlag{
			Name: "enable-stale-report",
			Usage: "",
		},
		cli.StringFlag{
			Name: "enable-mongodb",
			Usage: "",
		},
		cli.StringFlag{
			Name: "mongodb",
			Usage: "db inst for load mongodb",
		},
		cli.StringFlag{
			Name: "logconfig",
			Usage: "for seelog config",
		},
	}

	app.Action = func(c *cli.Context) error {
		if(len(c.String("config")) > 0){
			file, e := loadConfigJson(c.String("config"))
			if e != nil {
				fmt.Println("Read config file error: ", e)
				return e
			}

			e = json.Unmarshal(file, &Conf)
			if e != nil {
				fmt.Println("Unmarshal config file error: ", e)
				return e
			}
		}

		if(len(c.String("genesis")) > 0){
			file, e := loadConfigJson(c.String("genesis"))
			if e != nil {
				fmt.Println("Read genesis config file error: ", e)
				return e
			}

			e = json.Unmarshal(file, &GenConf)
			if e != nil {
				fmt.Println("Unmarshal config file error: ", e)
				return e
			}
		}

		if(len(c.String("datadir")) > 0){
			Conf.DataDir = c.String("datadir")
		}

		if(len(c.String("disable-api")) > 0){
			//TODO

		}

		if(len(c.String("apiport")) > 0){
			api_port,e:=strconv.Atoi(c.String("apiport"))
			if e != nil {
				fmt.Println(e.Error())
				return e
			}
			Conf.APIPort = api_port
		}

		if(len(c.String("p2pport")) > 0){
			/*p2p_port_,e:=strconv.Atoi(c.String("p2pport"))
			if e != nil {
				fmt.Println(e.Error())
				return e
			}*/
			Conf.P2PPort = c.String("p2pport")//p2p_port
		}

		if(len(c.String("peerlist")) > 0){
			peer_list := c.String("peerlist")
			Conf.PeerList = strings.Split(peer_list, ",")
		}

		if(len(c.String("delegate-signkey")) > 0){
			key := strings.Split(c.String("delegate-signkey"), ",")
			if(len(key) != 2){
				return fmt.Errorf("delegate-signkey params exception");
			}
			KeyPair.PrivateKey = key[0];
			KeyPair.PublicKey = key[1];
			Conf.KeyPairs[0] = KeyPair;
		}

		if(len(c.String("delegate")) > 0){
			//TODO
		}

		if(len(c.String("enable-stale-report")) > 0){
			fmt.Println(c.String("enable-stale-report"))
		}

		if(len(c.String("api_service_version")) > 0){
			Conf.ApiServiceVersion = c.String("api_service_version")
		}

		if(len(c.String("enable-mongodb")) > 0){
			//TODO
		}

		if(len(c.String("mongodb")) > 0){
			Conf.OptionDb = c.String("mongodb")
		}

		if(len(c.String("logconfig")) > 0){
			Conf.LogConfig = c.String("logconfig")
		}

		return nil
	}
	err := app.Run(os.Args)
	return Conf, err
}


func loadConfigJson(fn string) ([]byte, error) {
	file, e := ioutil.ReadFile(fn)
	if e != nil {
		return nil, e
	}

	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	return file, nil
}
