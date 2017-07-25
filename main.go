// Copyright Andrius Sutas BitfinexLendingBot [at] motoko [dot] sutas [dot] eu

package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/eAndrius/bitfinex-go"
)

var (
	configFile  = flag.String("conf", "default.conf", "Configuration file")
	updateLends = flag.Bool("updatelends", false, "Update lend offerings")
	dryRun      = flag.Bool("dryrun", false, "Output strategy decisions without placing orders")
	logToFile   = flag.Bool("logtofile", false, "Log lend history in coin's Log. Steps outputted to stdout")
)

// BotConfig ...
type BotConfig struct {
	Bitfinex BitfinexConf
	Strategy StrategyConf

	API *bitfinex.API
}

// BotConfigs ...
type BotConfigs []BotConfig

// BitfinexConf ...
type BitfinexConf struct {
	APIKey          string
	APISecret       string
	ActiveWallet    string
	MaxActiveAmount float64
	MinLoanUSD      float64
}

func init() {
	flag.Parse()

	if *logToFile {
		dsh, err := os.OpenFile("DASH.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("error opening file: " + err.Error())
		}
		eth, err := os.OpenFile("ETH.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("error opening file: " + err.Error())
		}
		ltc, err := os.OpenFile("LTC.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("error opening file: " + err.Error())
		}
		var file_names map[string]*ast.File
		file_names = make(map[string]*ast.File)
		file_name["dsh"] = dsh
		file_name["eth"] = eth
		file_name["ltc"] = ltc
	}
	
					   
}
}

func main() {
	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatal("Failed to open config file: " + err.Error())
	}

	decoder := json.NewDecoder(file)
	confs := BotConfigs{}
	err = decoder.Decode(&confs)
	if err != nil {
		log.Fatal("Failed to parse config file:" + err.Error())
	}

	for _, conf := range confs {
		log.Println("Using Bitfinex user API key: " + conf.Bitfinex.APIKey)
		conf.API = bitfinex.New(conf.Bitfinex.APIKey, conf.Bitfinex.APISecret)

		balance, err := conf.API.WalletBalances()
		if err != nil {
			log.Println("WARNING: Failed to get wallet funds, skipping: " + err.Error())
			continue
		}

		activeWallet := strings.ToLower(conf.Bitfinex.ActiveWallet)
		log.Println("\tDeposit wallet: " +
			strconv.FormatFloat(balance[bitfinex.WalletKey{"deposit", activeWallet}].Amount, 'f', -1, 64) +
			" " + activeWallet + " (swappable: " +
			strconv.FormatFloat(balance[bitfinex.WalletKey{"deposit", activeWallet}].Available, 'f', -1, 64) +
			" " + activeWallet + ")")

		if *updateLends {
			err = executeStrategy(conf, *dryRun, file_name[activeWallet])
			if err != nil {
				log.Println("WARNING: Failed to execute strategy: " + err.Error())
				continue
			}
		}
	}
}
