package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/noah-blockchain/Auto-rewards/models"
	"github.com/noah-blockchain/Auto-rewards/utils"
)

func GetValidators() (*models.ValidatorList, error) {
	responseStatus, err := http.Get(fmt.Sprintf("%s/status", os.Getenv("NODE_API_URL")))
	if err != nil {
		return nil, err
	}

	defer responseStatus.Body.Close()
	contentsStatus, err := ioutil.ReadAll(responseStatus.Body)
	if err != nil {
		return nil, err
	}

	nodeStatus := models.NodeStatus{}
	if err = json.Unmarshal(contentsStatus, &nodeStatus); err != nil {
		return nil, err
	}

	fmt.Println(nodeStatus.Result.LatestBlockHeight)

	responseValidators, err := http.Get(fmt.Sprintf("%s/validators?height=%s", os.Getenv("NODE_API_URL"), nodeStatus.Result.LatestBlockHeight))
	if err != nil {
		return nil, err
	}

	defer responseValidators.Body.Close()
	contentsValidators, err := ioutil.ReadAll(responseValidators.Body)
	if err != nil {
		return nil, err
	}

	validatorList := models.ValidatorList{}
	if err = json.Unmarshal(contentsValidators, &validatorList); err != nil {
		return nil, err
	}

	return &validatorList, nil
}

func GetDelegatorsListByNode(pubKey string) (map[string]float64, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/validators/%s", os.Getenv("EXPLORER_API_URL"), pubKey))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	validatorInfo := models.ValidatorInfo{}
	if err = json.Unmarshal(contents, &validatorInfo); err != nil {
		return nil, err
	}

	values := make(map[string]float64, validatorInfo.Result.DelegatorCount)
	for _, delegator := range validatorInfo.Result.DelegatorList {
		if delegator.Coin == os.Getenv("TOKEN") {
			if value, err := strconv.ParseFloat(delegator.Value, 64); err == nil {
				values[delegator.Address] = value
			}
		}
	}

	return values, nil
}

func GetAllDelegators() (map[string]float64, error) {
	allValidators, err := GetValidators()
	if err != nil {
		return nil, err
	}

	allDelegators := make(map[string]float64)

	for _, validator := range allValidators.Validators {
		delegators, err := GetDelegatorsListByNode(validator.PubKey)
		if err != nil {
			continue
		}

		if len(allDelegators) > 0 {
			for address, amount := range delegators {
				if val, ok := allDelegators[address]; ok {
					allDelegators[address] += val
				} else {
					allDelegators[address] = amount
				}
			}
		} else {
			allDelegators = delegators
		}
	}

	return allDelegators, nil
}

func getAllPayedDelegators() (map[string]float64, error) {
	allDelegators, err := GetAllDelegators()
	if err != nil {
		return nil, err
	}

	allPayedDelegators := make(map[string]float64)

	for address, _ := range allDelegators {
		amounts := allDelegators[address]

		minCoinsDelegate, err := strconv.ParseFloat(os.Getenv("MIN_COINS_DELEGATED"), 64)
		if err != nil {
			continue
		}

		if amounts >= minCoinsDelegate {
			allPayedDelegators[address] = amounts
		}
	}
	return allPayedDelegators, nil
}

func getTotalDelegatedCoins() (float64, error) {

	payedDelegatorsList, err := getAllPayedDelegators()
	if err != nil {
		return 0.0, err
	}

	totalDelegatedCoins := 0.0

	for _, amount := range payedDelegatorsList {
		totalDelegatedCoins += amount
	}

	return totalDelegatedCoins, nil
}

func getWalletBalances(address string) (*models.AddressInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/v1/addresses/%s", os.Getenv("EXPLORER_API_URL"), address))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	addressInfo := models.AddressInfo{}
	if err = json.Unmarshal(contents, &contents); err != nil {
		return nil, err
	}

	return &addressInfo, nil
}

func getNoahBalance(address string) (float64, error) {
	balances, err := getWalletBalances(address)
	if err != nil {
		return 0.0, err
	}

	for _, value := range balances.Data.Balances {
		if value.Coin == "NOAH" {
			if value, err := strconv.ParseFloat(value.Amount, 64); err == nil {
				return value, nil
			}
		}
	}

	return 0.0, nil
}

func CreateMultiSendList(walletFrom string, payCoinName string) (*[]models.MultiSendItem, error) {
	totalDelegatedCoins, err := getTotalDelegatedCoins()
	if err != nil || totalDelegatedCoins == 0 {
		return nil, err
	}

	payedDelegatedList, err := getAllPayedDelegators()
	if err != nil {
		return nil, err
	}

	balance, err := getNoahBalance(walletFrom)
	if err != nil {
		return nil, err
	}

	toBePayed := balance * 0.95

	multiSendList := make([]models.MultiSendItem, len(payedDelegatedList))
	for address, amount := range payedDelegatedList {
		percent := amount * 100 / totalDelegatedCoins
		value := utils.FloatToBigInt(toBePayed * percent * 0.01)

		multiSendList = append(multiSendList, models.MultiSendItem{
			Coin:  payCoinName,
			To:    address,
			Value: value,
		})
	}

	return &multiSendList, nil
}
